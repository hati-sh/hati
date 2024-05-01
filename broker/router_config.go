package broker

type RouterType string

const RouterDirect = RouterType("direct")
const RouterFanout = RouterType("fanout")
const RouterTopic = RouterType("topic")

type RouterConfig struct {
	Name       string     `json:"name"`
	RouterType RouterType `json:"type"`
}
