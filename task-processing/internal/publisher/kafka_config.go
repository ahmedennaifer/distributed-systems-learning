package publisher

import "context"

type KafkaConfig struct {
	Ctx     context.Context
	Network string
	Addr    string
}

func NewKafkaConfig(network string, addr string) KafkaConfig {
	return KafkaConfig{
		Ctx:     context.TODO(),
		Network: network,
		Addr:    addr,
	}
}
