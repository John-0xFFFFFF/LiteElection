package election

// not so reasonable
type LiteElection interface {
	IsLeader() bool
	Do()
}
