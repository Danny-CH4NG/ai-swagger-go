package apiv1

import (
	"sync"
)

var (
	co     *controller
	coLock sync.Once
)

type controller struct{}

func Controller() *controller {
	coLock.Do(func() {
		co = &controller{}
	})
	return co
}
