package jobs

type Job interface {
	Execute() error
}
