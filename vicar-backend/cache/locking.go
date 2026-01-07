package cache

import (
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

var rs *redsync.Redsync

func initializeLocking() {
	rs = redsync.New(goredis.NewPool(Redis))
}

func NewLock(name string, expiry *time.Duration, tries *int) *redsync.Mutex {
	options := []redsync.Option{}

	if expiry != nil {
		options = append(options, redsync.WithExpiry(*expiry))
	}

	if tries != nil {
		options = append(options, redsync.WithTries(*tries))
	}

	return rs.NewMutex(name, options...)
}

func LockProject(projecId string) (*redsync.Mutex, error) {
	lock := NewLock("lock:project:"+projecId, nil, nil)
	if err := lock.Lock(); err != nil {
		return nil, err
	}

	return lock, nil
}
