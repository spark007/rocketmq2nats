package core

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"os"
	"rocketmq2nats/pkg/config"
	"strconv"
	"time"
)

type IProducer interface {
	Publish(data []byte)
	Shutdown()
}

type NatsProducer struct {
	conn *nats.Conn
	cfg  *config.NatsConfig
}

func (p *NatsProducer) Publish(data []byte) {
	if p.conn != nil {
		err := p.conn.Publish(p.cfg.Subject, data)
		if err != nil {
			zap.L().Error("nats publish error", zap.Error(err))
		}
	}
}

func (p *NatsProducer) Shutdown() {
	if p.conn != nil {
		p.conn.Close()
	}
}

func NewNatsProducer(cfg *config.NatsConfig) *NatsProducer {
	p := &NatsProducer{
		cfg: cfg,
	}
	var err error

	p.conn, err = nats.Connect(cfg.Host+":"+strconv.Itoa(cfg.Port),
		nats.MaxReconnects(3),
		nats.ReconnectWait(15*time.Second),
		nats.ReconnectHandler(func(_ *nats.Conn) {

		}),
	)
	if err != nil {
		zap.L().Error("connect nats error", zap.String("error", err.Error()))
		os.Exit(-1)
	}

	return p
}
