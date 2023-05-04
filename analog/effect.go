package analog

type Effect interface {
	Process([]Signal) Signal
}
