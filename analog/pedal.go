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

func (p *BasePedal) GetInputJack() []InputJack {
	return p.inputs
}

func (p *BasePedal) GetOutputJack() []OutputJack {
	return p.outputs
}

func (p *BasePedal) Toggle() {
	p.isOn = !p.isOn
}

func (p *BasePedal) Run() {
	if (p.isOn) {
		// call effect
	} else {
		// passthrough
	}
}
