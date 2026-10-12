package grpc

import (
	"backend-rms/core/config"
	"backend-rms/core/logger"
	pb "backend-rms/grpc/et6api"
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var TSPCli pb.ApiServiceClient

type Response interface {
	GetCResp() *pb.CommonResponse
}

func InitTSPClient() {
	conf := config.Conf.GRPC.TSP
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryClientInterceptor("TSP", conf.Timeout)),
	}

	cc, err := grpc.Dial(conf.Address, opts...)
	if err != nil {
		panic(err)
	}

	TSPCli = pb.NewApiServiceClient(cc)
}

func unaryClientInterceptor(name string, timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()

		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			err = errors.WithStack(err)
		} else {
			if res, ok := reply.(Response); ok {
				if !res.GetCResp().GetIsSucc() {
					err = errors.New("rpc call failed")
				}
			}
		}
		writeLog(name, cc.Target(), method, req, reply, start, err)
		return err
	}
}

func writeLog(name, address, method string, req, reply interface{}, start time.Time, err error) {
	fields := []interface{}{
		"name", name,
		"address", address,
		"method", method,
		"req", req,
		"reply", reply,
		"elapsed_time", time.Since(start).String(),
	}

	if err != nil {
		logger.Errorw(err.Error(), fields...)
		return
	}

	logger.Infow("OK", fields...)
}
