package analog

type Pedal interface {
	GetInputJack() []InputJack
	GetOutputJack() []OutputJack
	Toggle()
}
