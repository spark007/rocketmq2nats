package core

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"os"
	"rocketmq2nats/pkg/config"
	"strconv"
	"time"
)

type Callback func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)

type IConsumer interface {
	Start()
	Shutdown()
	Subscribe(callback Callback)
	Unsubscribe()
}

type RocketMQConsumer struct {
	consumer rocketmq.PushConsumer
	cfg      *config.RocketMQConfig
}

func (c *RocketMQConsumer) Start() {
	if err := c.consumer.Start(); err != nil {
		zap.L().Error("start consumer error", zap.String("error", err.Error()))
		os.Exit(-1)
	}
	zap.L().Info("start consumer success")
}

func (c *RocketMQConsumer) Shutdown() {
	if err := c.consumer.Shutdown(); err != nil {
		zap.L().Error("shutdown consumer error", zap.String("error", err.Error()))
	}
	zap.L().Info("shutdown consumer success")
}

func (c *RocketMQConsumer) Subscribe(callback Callback) {
	selector := consumer.MessageSelector{}

	if err := c.consumer.Subscribe("ficc", selector, callback); err != nil {
		zap.L().Error("subscribe topic error", zap.String("error", err.Error()))
		os.Exit(-1)
	}
	zap.L().Info("subscribe topic success")
}

func (c *RocketMQConsumer) Unsubscribe() {
	if err := c.consumer.Unsubscribe(c.cfg.TopicName); err != nil {
		zap.L().Error("unsubscribe topic error", zap.String("error", err.Error()))
	}
}

func NewRockerMQConsumer(cfg *config.RocketMQConfig) *RocketMQConsumer {
	rc := &RocketMQConsumer{
		cfg: cfg,
	}

	nameserverResolver := primitive.NewPassthroughResolver([]string{cfg.Host + ":" + strconv.Itoa(cfg.Port)})

	traceCfg := &primitive.TraceConfig{
		Access:   primitive.Local,
		Resolver: nameserverResolver,
	}
	var err error
	// 创建consumer
	rc.consumer, err = rocketmq.NewPushConsumer(
		// 设置消费者组
		consumer.WithGroupName(cfg.GroupName),
		// 设置服务地址
		consumer.WithNsResolver(nameserverResolver),
		// 设置trace, 用于发送消息时记录消息轨迹,如果不需要，不设置即可
		consumer.WithTrace(traceCfg),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromWhere(cfg.Offset)),
		consumer.WithConsumeTimeout(time.Second*time.Duration(cfg.Timeout)),
		consumer.WithRetry(cfg.MaxRetry),
	)
	if err != nil {
		zap.L().Error("init consumer error", zap.String("error", err.Error()))
		os.Exit(0)
	}
	return rc
}
