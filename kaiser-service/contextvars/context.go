package contextvars

// Context variables used for plugins
const (
	JobName     = "$JOB_NAME"
	JobVersion  = "$JOB_VERSION"
	TaskName    = "$TASK_NAME"
	TaskVersion = "$TASK_VERSION"
)

const (
	//DefaultLogFolder Cointains the default name of the log folder used by kaiser-service
	DefaultLogFolder = "logs"
)
