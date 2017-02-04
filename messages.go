package webq

type AuthMessage struct {
	Content string
}

type SubscribeMessage struct {
	Topic string
}

type PublishMessage struct {
	Topic   string
	Message []byte
}
