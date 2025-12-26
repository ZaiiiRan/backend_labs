package messages

type Message interface {
	RoutingKey() string
}
