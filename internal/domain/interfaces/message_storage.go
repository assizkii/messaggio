package interfaces

type Message struct {
	Phone string
	Text  string
}

type MessageStorage interface {
	Add(m *Message) (int, error)
}
