package publisher

import (
	"context"
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

func Init() (*KafkaClient, error) {
	cfg := NewKafkaConfig("tcp", "broker:9092")
	kafkaClient, err := NewKafkaClient().
		WithConfig(cfg).
		WithTopics([]kafka.TopicConfig{
			{
				Topic:             "tasks",
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
			{
				Topic:             "failed_tasks", // dead letter queue?
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
		}).
		SetupController()
	if err != nil {
		fmt.Printf("%v", err)
		return &KafkaClient{}, err
	}

	for _, topic := range kafkaClient.Topics {
		if err := kafkaClient.ControllerConn.CreateTopics(topic); err != nil {
			fmt.Printf("failed to create topics: %v\n", err)
			return &KafkaClient{}, err
		}
		fmt.Printf("Successfully created topic: %v\n", topic.Topic)
	}

	fmt.Printf("Created %v topic(s)\n", len(kafkaClient.Topics))
	return kafkaClient, nil
}

func (kClient *KafkaClient) Publish(topic string, key, message []byte) error {
	brokers := []string{kClient.Config.Addr}
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	defer w.Close()
	kMsg := kafka.Message{
		Key:   key,
		Value: message,
	}
	if err := w.WriteMessages(context.TODO(), kMsg); err != nil {
		return fmt.Errorf("error sending message : %v to topic %v:", string(message), topic, err)
	}
	fmt.Println("sent message %v to topic %v", string(message), topic)
	return nil
}
