package analog

type Pedal interface {
	GetInputJack() []InputJack
	GetOutputJack() []OutputJack
	Toggle()
	Run(shouldRun *bool)
}

func NewPedal(effect Effect) Pedal {
	return &BasePedal{
		effect:  effect,
		inputs:  []InputJack{NewInputJack()},
		outputs: []OutputJack{NewOutputJack()},
	}
}

func NewDummyPedal() Pedal {
	return &BasePedal{
		effect:  &DummyEffect{},
		inputs:  []InputJack{NewInputJack()},
		outputs: []OutputJack{NewOutputJack()},
	}
}

type BasePedal struct {
	effect  Effect
	inputs  []InputJack
	outputs []OutputJack
	isOn    bool
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

// TODO: define strategy for multiple input and outputs "channels" (a.k.a. jacks)
func (p *BasePedal) Run(shouldRun *bool) {
	for *shouldRun {
		bufferSize := cap(p.inputs[0].GetWire())

		signals := p.inputs[0].BufferedReceiveSignal(bufferSize)
		if p.isOn {
			signals = p.effect.Process(signals)
		}

		p.outputs[0].BufferedSendSignal(signals)
	}
}
