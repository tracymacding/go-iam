package db

import (
	"fmt"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
	gService  Service
)

type Driver interface {
	Open(args ...interface{}) (Service, error)
}

func RegisterDriver(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()

	if driver == nil {
		panic("db: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("db: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func Open(driverName string, args ...interface{}) error {
	driversMu.RLock()
	driveri, ok := drivers[driverName]
	driversMu.RUnlock()

	if !ok {
		return fmt.Errorf("db: unknown driver %q (forgotten import?)", driverName)
	}

	var err error
	gService, err = driveri.Open(args...)
	return err
}

func ActiveService() Service {
	return gService
}
