package analog

type Pedal interface {
	GetInputJack() []InputJack
	GetOutputJack() []OutputJack
	Toggle()
}

func NewPedal(effect Effect) Pedal {
	return &BasePedal{
		effect: effect,
		inputs: make([]InputJack, 1),
		outputs: make([]OutputJack, 1),
	}
}

func NewDummyPedal() Pedal {
	return &BasePedal{
		effect: &DummyEffect{},
		inputs: make([]InputJack, 1),
		outputs: make([]OutputJack, 1),
	}
}

type BasePedal struct {
	effect Effect
	inputs []InputJack
	outputs []OutputJack
	isOn bool
}

