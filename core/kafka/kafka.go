package kafka

import (
	"backend-rms/core/config"
	"backend-rms/core/kafka/handle"
	"backend-rms/core/logger"
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
)

type consumerGroupHandler struct {
	topic string
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		logger.Infof("[Topic: %s] Partition: %d, Offset: %d, Key: %s, Value: %s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		switch msg.Topic {
		case config.Conf.Kafka.Topic.Order:
			handle.HandleOrderEvent(msg)
			break
		case config.Conf.Kafka.Topic.Group:
			handle.HandleGroupEvent(msg)
			break
		default:
			logger.Infof("Unknown topic: %s, skipping message", msg.Topic)
		}
		sess.MarkMessage(msg, "") // 提交 offset
	}
	return nil
}

func Start() {
	configSarama := sarama.NewConfig()
	configSarama.Version = sarama.V2_1_0_0
	configSarama.Consumer.Return.Errors = true
	configSarama.Consumer.Offsets.Initial = sarama.OffsetNewest // 从最新的开始消费

	// 多个 topic
	topics := []string{config.Conf.Kafka.Topic.Group, config.Conf.Kafka.Topic.Order}

	// Kafka broker 地址
	brokers := strings.Split(config.Conf.Kafka.Addr, ",")

	// 创建 consumer group
	groupID := config.Conf.Kafka.Consumer
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, configSarama)
	if err != nil {
		logger.Errorf("Failed to start consumer group: %v\n", err)
	}
	defer consumerGroup.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 信号监听
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	var wg sync.WaitGroup
	// 持续读取错误通道，避免阻塞
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case err, ok := <-consumerGroup.Errors():
				if !ok {
					return
				}
				if err != nil {
					logger.Errorf("Consumer group error: %v\n", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// 启动消费
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// 注意：这里传入所有 topics
			err := consumerGroup.Consume(ctx, topics, &consumerGroupHandler{})
			if err != nil {
				logger.Errorf("Error from consumer: %v\n", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-sigchan
	logger.Info("Interrupt signal received, shutting down...")
	cancel()
	wg.Wait()
}
