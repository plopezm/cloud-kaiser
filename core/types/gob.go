package types

import (
	"encoding/gob"
)

func RegisterCoreTypes() {
	gob.Register(Task{})
	gob.Register(JobTask{})
	gob.Register(Job{})
	gob.Register(JobActivation{})
	gob.Register(JobArgs{})
}
