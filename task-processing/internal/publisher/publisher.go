package publisher

import (
	"fmt"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	Config         KafkaConfig
	ControllerConn *kafka.Conn
	Topics         []kafka.TopicConfig
	Err            error
}

func NewKafkaClient() *KafkaClient {
	return &KafkaClient{}
}

func (k *KafkaClient) SetupController() (*KafkaClient, error) {
	conn, err := kafka.Dial(k.Config.Network, k.Config.Addr)
	if err != nil {
		return k, err
	}
	controller, err := conn.Controller()
	if err != nil {
		return k, err
	}
	if k.Err != nil {
		return k, k.Err
	}

	controllerConn, err := kafka.Dial(k.Config.Network, net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return k, err
	}
	k.ControllerConn = controllerConn
	return k, nil
}

func (k *KafkaClient) WithConfig(config KafkaConfig) *KafkaClient {
	k.Config = config
	return k
}

func (k *KafkaClient) WithTopics(topics []kafka.TopicConfig) *KafkaClient {
	// maybe validate each topic
	k.Topics = topics
	return k
}

func CreateTopics() error {
	cfg := NewKafkaConfig("tcp", "broker:9092")
	kafkaClient, err := NewKafkaClient().
		WithConfig(cfg).
		WithTopics([]kafka.TopicConfig{
			{
				Topic:             "tasks",
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
		}).
		SetupController()
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}

	for _, topic := range kafkaClient.Topics {
		if err := kafkaClient.ControllerConn.CreateTopics(topic); err != nil {
			fmt.Printf("failed to create topics: %v\n", err)
			return err
		}
		fmt.Printf("Successfully created topic: %v\n", topic.Topic)
	}

	fmt.Printf("Created %v topic(s)\n", len(kafkaClient.Topics))
	return nil
}
