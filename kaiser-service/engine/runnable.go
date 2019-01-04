package engine

type RunnableStatus string

const (
	RunnableStatusRunning = "RUNNING"
	RunnableStatusStopped = "STOPPED"
)

type Runnable interface {
	GetIdentifier() string
	SetStartTime(ts int64)
	GetStartTime() int64
	GetStatus() RunnableStatus
	SetStatus(status RunnableStatus)
	SetParameters(map[string]interface{})
	SetResultStatus(result bool)
	GetResultStatus() bool
	Run()
}
