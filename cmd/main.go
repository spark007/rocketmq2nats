package main

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"rocketmq2nats/pkg/config"
	"rocketmq2nats/pkg/core"
	"rocketmq2nats/pkg/log"
	"syscall"
)

func main() {
	log.InitLogger()
	config.InitConfig()

	globalCfg := config.NewGlobalConfig()

	p := core.NewNatsProducer(globalCfg.NatsCfg)
	defer p.Shutdown()

	callback := func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		//zap.L().Info("consume message", zap.Int("message count", len(msgs)))
		for _, msg := range msgs {
			//zap.L().Info("consume message", zap.String("message", string(msg.Body)))
			go p.Publish(msg.Body)
		}
		return consumer.ConsumeSuccess, nil
	}

	c := core.NewRockerMQConsumer(globalCfg.RocketMqCfg, callback)
	c.Subscribe()
	c.Start()
	defer c.Shutdown()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM) // 监听 Ctrl+C 和终止信号
	<-sigChan                                               // 阻塞直到收到信号
	zap.L().Info("received interrupt signal, exiting...")
}
