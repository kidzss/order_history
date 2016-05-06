package main

import (
	"github.com/Shopify/sarama"
	"log"
)

type TopicInfo struct {
	Msgs         chan *sarama.ConsumerMessage
	PartConsumer []sarama.PartitionConsumer
}

type KafkaClient struct {
	Consumer sarama.Consumer
	Topics   map[string]*TopicInfo
	Server   []string
}

func FindOffset(where string) int64 {
	var ret int64
	switch where {
	case "begin":
		ret = sarama.OffsetOldest
	case "now":
		ret = sarama.OffsetNewest
	default:
		ret = sarama.OffsetOldest
	}
	return ret
}

func NewKafkaClient() *KafkaClient {
	return &KafkaClient{Topics: make(map[string]*TopicInfo, 0)}
}

func (c *KafkaClient) NewConsumer(conf *Configure) error {
	hostports := conf.Kafka.Hosts
	consumer, err := sarama.NewConsumer(hostports, nil)
	if err != nil {
		log.Printf("[kafka] new a consumer %+v error, %s\n", hostports, err)
	} else {
		log.Printf("[kafka] new a consumer %+v success.\n", hostports)
	}
	c.Consumer = consumer
	c.Server = hostports
	return err
}

func (c *KafkaClient) Close() {
	c.Consumer.Close()
	for _, v := range c.Topics {
		for i, _ := range v.PartConsumer {
			v.PartConsumer[i].Close()
		}
	}
	log.Printf("[kafka] close a connect success.\n")
}

func (c *KafkaClient) GetTopicMsg(topic string, where string) {
	partitions, err := c.Consumer.Partitions(topic)
	log.Printf("[kafka] topic:%s, partitions:%+v\n", topic, partitions)
	if err == nil {
		if _, ok := c.Topics[topic]; ok == false {
			topic_m := &TopicInfo{Msgs: make(chan *sarama.ConsumerMessage, 30)}
			c.Topics[topic] = topic_m
		}
		for _, partition := range partitions {
			//每一个topic_partition创建一个go
			go c.OpenParttionConsumer(topic, partition, where)
		}
	}
}

func (c *KafkaClient) OpenParttionConsumer(topic string, partition int32, where string) {
	partitionConsumer, err := c.Consumer.ConsumePartition(topic, partition, FindOffset(where))
	if err == nil {
		log.Printf("[kafka] topic[%s] at partition[%d] open a partitionconsumer, read from %s\n", topic, partition, where)
		c.Topics[topic].PartConsumer = append(c.Topics[topic].PartConsumer, partitionConsumer)
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				c.Topics[topic].Msgs <- msg
			}
		}
	}
}
