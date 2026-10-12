package handle

import (
	"backend-rms/core/logger"
	"backend-rms/model"
	"backend-rms/service"
	"backend-rms/tr/init_data"

	"encoding/json"

	"github.com/Shopify/sarama"
)

type GroupEvent struct {
	GroupEventType int   `json:"Group_event_type"`
	Login          int64 `json:"Login"`
	OldGroupName   string
	NewGroupName   string
	Platform       int
	ModifyAt       int64
}

func HandleGroupEvent(msg *sarama.ConsumerMessage) {
	var groupEvent GroupEvent
	if err := json.Unmarshal(msg.Value, &groupEvent); err != nil {
		logger.Errorf("Failed to parse group message: %v, raw: %s", err, string(msg.Value))
		return
	}

	// 判断是否是新增用户，而非转组
	if groupEvent.GroupEventType == 1 || groupEvent.OldGroupName == "" {
		return
	}

	var err error
	var oldBook string
	if _, ok := service.InitializedData.Groups[groupEvent.OldGroupName]; !ok {
		// 如果组不在当前初始化的组内，去数据库查询，如果查询到更新到初始化的组内，如果没有查询到记录错误日志
		oldBook, err = service.GroupService.GetBookByGroup(groupEvent.OldGroupName)
		if err != nil {
			logger.Errorf("Select old group err:%v", err)
			return
		}
		service.InitializedData.Groups[groupEvent.OldGroupName] = &model.Group{
			Name: groupEvent.OldGroupName,
			Book: oldBook,
		}
	} else {
		oldBook = service.InitializedData.Groups[groupEvent.OldGroupName].Book
	}

	var newBook string
	if _, ok := service.InitializedData.Groups[groupEvent.NewGroupName]; !ok {
		// 如果组不在当前初始化的组内，去数据库查询，如果查询到更新到初始化的组内，如果没有查询到记录错误日志
		newBook, err = service.GroupService.GetBookByGroup(groupEvent.NewGroupName)
		if err != nil {
			logger.Errorf("Select new group err:%v", err)
			return
		}
		service.InitializedData.Groups[groupEvent.NewGroupName] = &model.Group{
			Name: groupEvent.NewGroupName,
			Book: newBook,
		}
	} else {
		newBook = service.InitializedData.Groups[groupEvent.NewGroupName].Book
	}

	// 同book转组无需操作
	if oldBook == newBook {
		return
	}

	// 处理账户持仓订单
	// 更新初始化的账号的group前处理订单
	for _, trade := range service.InitializedData.Trades {
		if trade.Login == groupEvent.Login {
			service.SummaryService.HandleForOrder(trade, "close")
			service.ExposureService.HandleForOrder(trade, "close")
		}
	}

	// 更新初始化的账号的group后处理订单
	service.InitializedData.Accounts[groupEvent.Login].Group = groupEvent.NewGroupName
	for _, trade := range service.InitializedData.Trades {
		if trade.Login == groupEvent.Login {
			service.SummaryService.HandleForOrder(trade, "open")
			service.ExposureService.HandleForOrder(trade, "open")
		}
	}

	// 把summary变化更新到redis (不同book转组，A,B组都有变化)
	for _, v := range init_data.Books {
		service.CacheService.SetSummaryBook(v)
	}
}
