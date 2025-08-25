package main

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"rocketmq2nats/pkg/config"
	"rocketmq2nats/pkg/core"
	"rocketmq2nats/pkg/log"
	"time"
)

func main() {
	log.InitLogger()
	config.InitConfig()

	ctx, _ := context.WithTimeout(context.Background(), 300*time.Second)

	p := core.NewPiper(ctx)
	callback := func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		//zap.L().Info("consume message", zap.Int("message count", len(msgs)))
		for _, msg := range msgs {
			p.ChMsg <- msg.Body
		}
		return consumer.ConsumeSuccess, nil
	}

	go p.Source(callback)
	go p.Sink()
	defer p.Shutdown()

	<-p.Quit
}
