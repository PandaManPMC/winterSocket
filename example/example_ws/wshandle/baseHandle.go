package wshandle

import (
	"sync"
)

type baseHandle struct {
	lock *sync.Mutex
}

func (that *baseHandle) locked() {
	that.lock.Lock()
}

func (that *baseHandle) unLocked() {
	that.lock.Unlock()
}
