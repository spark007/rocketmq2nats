package core

import (
	"context"
	"go.uber.org/zap"
	"rocketmq2nats/pkg/config"
)

type Piper struct {
	ctx      context.Context
	consumer IConsumer
	producer IProducer
	ChMsg    chan []byte
	Quit     chan bool
}

func NewPiper(ctx context.Context) *Piper {
	globalCfg := config.NewGlobalConfig()
	p := NewNatsProducer(globalCfg.NatsCfg)

	/*	callback := func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		//zap.L().Info("consume message", zap.Int("message count", len(msgs)))
		for _, ChMsg := range msgs {
			//zap.L().Info("consume message", zap.String("message", string(ChMsg.Body)))
			go p.Publish(ChMsg.Body)
		}
		return consumer.ConsumeSuccess, nil
	}*/

	c := NewRockerMQConsumer(globalCfg.RocketMqCfg)

	return &Piper{
		ctx:      ctx,
		consumer: c,
		producer: p,
		ChMsg:    make(chan []byte, 100),
		Quit:     make(chan bool),
	}
}

func (p *Piper) Source(callback Callback) {
	p.consumer.Subscribe(callback)
	p.consumer.Start()
}

func (p *Piper) Sink() {
	for {
		select {
		case <-p.ctx.Done():
			zap.L().Info("sink quit")
			p.Quit <- true
			return
		case msg := <-p.ChMsg:
			p.producer.Publish(msg)
		}
	}
}

func (p *Piper) Shutdown() {
	p.consumer.Shutdown()
	close(p.ChMsg)
	p.producer.Shutdown()
	close(p.Quit)
}
