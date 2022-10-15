package kafka

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"github.com/x-foby/kakafka/pkg/kaql"
)

type Broker struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	ID   int    `json:"id"`
	Rack string `json:"rack"`
}

func newBrokers(list []kafka.Broker) []Broker {
	result := make([]Broker, len(list))

	for i, b := range list {
		result[i] = Broker(b)
	}

	return result
}

type Partition struct {
	Topic    string   `json:"topic"`
	ID       int      `json:"id"`
	Leader   Broker   `json:"leader"`
	Replicas []Broker `json:"replicas"`
	Isr      []Broker `json:"-"`
	Error    error    `json:"-"`
}

func newPartiotions(list []kafka.Partition) []Partition {
	result := make([]Partition, len(list))

	for i, p := range list {
		result[i] = Partition{
			Topic:    p.Topic,
			ID:       p.ID,
			Leader:   Broker(p.Leader),
			Replicas: newBrokers(p.Replicas),
			Isr:      newBrokers(p.Isr),
			Error:    p.Error,
		}
	}

	return result
}

type Topic struct {
	Name       string      `json:"name"`
	Internal   bool        `json:"internal"`
	Partitions []Partition `json:"partitions"`
	Error      error       `json:"-"`
}

func newTopic(t kafka.Topic) Topic {
	return Topic{
		Name:       t.Name,
		Internal:   t.Internal,
		Partitions: newPartiotions(t.Partitions),
		Error:      t.Error,
	}
}

type ReplicaAssignment struct {
	Partition int   `json:"partition"`
	Replicas  []int `json:"replicas"`
}

type ReplicaAssignments []ReplicaAssignment

func (r ReplicaAssignments) toKafka() []kafka.ReplicaAssignment {
	result := make([]kafka.ReplicaAssignment, len(r))

	for i, a := range r {
		result[i] = kafka.ReplicaAssignment(a)
	}

	return result
}

type ConfigEntry struct {
	ConfigName  string `json:"name"`
	ConfigValue string `json:"value"`
}

type ConfigEntries []ConfigEntry

func (r ConfigEntries) toKafka() []kafka.ConfigEntry {
	result := make([]kafka.ConfigEntry, len(r))

	for i, e := range r {
		result[i] = kafka.ConfigEntry(e)
	}

	return result
}

type TopicConfig struct {
	Topic              string              `json:"topic"`
	NumPartitions      int                 `json:"numPartitions"`
	ReplicationFactor  int                 `json:"replicationFactor"`
	ReplicaAssignments []ReplicaAssignment `json:"replicaAssignments"`
	ConfigEntries      []ConfigEntry       `json:"configEntries"`
}

func (c TopicConfig) toKafka() kafka.TopicConfig {
	return kafka.TopicConfig{
		Topic:              c.Topic,
		NumPartitions:      c.NumPartitions,
		ReplicationFactor:  c.ReplicationFactor,
		ReplicaAssignments: ReplicaAssignments(c.ReplicaAssignments).toKafka(),
		ConfigEntries:      ConfigEntries(c.ConfigEntries).toKafka(),
	}
}

type Conn struct {
	brokers []string
	client  *kafka.Client
	topics  []Topic
}

var ErrNoBrokers = errors.New("no brokers")

func Connect(ctx context.Context, brokers ...string) (*Conn, error) {
	if len(brokers) == 0 {
		return nil, ErrNoBrokers
	}

	client := kafka.Client{
		Addr: kafka.TCP(brokers...),
	}

	conn := Conn{
		brokers: brokers,
		client:  &client,
	}

	if err := conn.loadTopics(ctx); err != nil {
		return nil, err
	}

	return &conn, nil
}

func (c *Conn) loadTopics(ctx context.Context) error {
	metadata, err := c.client.Metadata(ctx, &kafka.MetadataRequest{
		Addr: kafka.TCP(c.brokers...),
	})
	if err != nil {
		return err
	}

	c.topics = make([]Topic, len(metadata.Topics))

	for i, t := range metadata.Topics {
		c.topics[i] = newTopic(t)
	}

	return nil
}

func (c *Conn) GetTopics(ctx context.Context, resresh bool) ([]Topic, error) {
	if !resresh {
		return c.topics, nil
	}

	if err := c.loadTopics(ctx); err != nil {
		return nil, err
	}

	return c.topics, nil
}

var errNoTopic = errors.New("no topic")

func (c *Conn) CreateTopic(ctx context.Context, topic TopicConfig) (*Topic, error) {
	if _, err := c.client.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Addr:   kafka.TCP(c.brokers...),
		Topics: []kafka.TopicConfig{topic.toKafka()},
	}); err != nil {
		return nil, err
	}

	if err := c.loadTopics(ctx); err != nil {
		return nil, err
	}

	if t, ok := c.topic(topic.Topic); ok {
		return (*Topic)(&t), nil
	}

	return nil, errNoTopic
}

func (c *Conn) DeleteTopic(ctx context.Context, topic string) error {
	if _, err := c.client.DeleteTopics(ctx, &kafka.DeleteTopicsRequest{
		Addr:   kafka.TCP(c.brokers...),
		Topics: []string{topic},
	}); err != nil {
		return err
	}

	return c.loadTopics(ctx)
}

type PartitionOffset struct {
	Partition       int    `json:"partition"`
	CommittedOffset int64  `json:"committedOffset"`
	Metadata        string `json:"metadata"`
	Error           error  `json:"-"`
}

func newPartitionOffset(partitions []kafka.OffsetFetchPartition) []PartitionOffset {
	result := make([]PartitionOffset, len(partitions))

	for i, offsets := range partitions {
		result[i] = PartitionOffset(offsets)
	}

	return result
}

type ConsumerOffset struct {
	Consumer string            `json:"consumer"`
	Offsets  []PartitionOffset `json:"offsets"`
}

func (c *Conn) ConsumerOffsets(ctx context.Context, topicName string) ([]ConsumerOffset, error) {
	groups, err := c.client.ListGroups(ctx, &kafka.ListGroupsRequest{
		Addr: kafka.TCP(c.brokers...),
	})
	if err != nil {
		return nil, err
	}

	topic, ok := c.topic(topicName)
	if !ok {
		return nil, errNoTopic
	}

	partitions := make([]int, len(topic.Partitions))
	for i, p := range topic.Partitions {
		partitions[i] = p.ID
	}

	offsets := make([]ConsumerOffset, len(groups.Groups))

	for i, group := range groups.Groups {
		resp, err := c.client.OffsetFetch(ctx, &kafka.OffsetFetchRequest{
			Addr:    kafka.TCP(c.brokers...),
			Topics:  map[string][]int{topic.Name: partitions},
			GroupID: group.GroupID,
		})
		if err != nil {
			return nil, err
		}

		offsets[i] = ConsumerOffset{
			Consumer: group.GroupID,
			Offsets:  newPartitionOffset(resp.Topics[topic.Name]),
		}
	}

	return offsets, nil
}

type Message struct {
	Key    string    `json:"key"`
	Value  string    `json:"value"`
	Offset int64     `json:"offset"`
	Time   time.Time `json:"time"`
}

func newMessageFromRecord(rec *protocol.Record) (Message, error) {
	key, err := readString(rec.Key)
	if err != nil {
		return Message{}, err
	}

	value, err := readString(rec.Value)
	if err != nil {
		return Message{}, err
	}

	return Message{
		Offset: rec.Offset,
		Time:   rec.Time,
		Key:    key,
		Value:  value,
	}, nil
}

func readString(data protocol.Bytes) (string, error) {
	buf, err := io.ReadAll(data)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (c *Conn) GetMessages(ctx context.Context, topicName, query string) ([]Message, error) {
	resp, err := c.client.Fetch(ctx, &kafka.FetchRequest{
		Addr:  kafka.TCP(c.brokers...),
		Topic: topicName,
	})
	if err != nil {
		return nil, err
	}

	kaQL, err := kaql.New(query)
	if err != nil {
		return nil, err
	}

	var messages []Message

	for i := 0; i <= 50; {
		rec, err := resp.Records.ReadRecord()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		msg, err := newMessageFromRecord(rec)
		if err != nil {
			return nil, err
		}

		if kaQL.Match(msg) {
			messages = append(messages, msg)
			i++

			continue
		}
	}

	return messages, nil
}

func (c *Conn) topic(name string) (Topic, bool) {
	for _, t := range c.topics {
		if t.Name == name {
			return t, true
		}
	}

	return Topic{}, false
}
