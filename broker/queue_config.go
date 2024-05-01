package broker

type QueueType string

const QueueMemory = QueueType("memory")
const QueueHdd = QueueType("hdd")

type QueueConfig struct {
	Name      string             `json:"name"`
	QueueType QueueType          `json:"type"`
	Bind      []QueueBindOptions `json:"bind"`
}

type QueueBindOptions struct {
	RouterName string `json:"routerName"`
	RoutingKey string `json:"routingKey"`
}
