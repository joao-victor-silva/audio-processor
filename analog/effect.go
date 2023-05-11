package analog

type Effect interface {
	Process([]Signal) []Signal
}

type DummyEffect struct {
}

func (e *DummyEffect) Process(signals []Signal) []Signal {
	return signals
}
