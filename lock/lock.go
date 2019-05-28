package lock

import (
	"os"
	"path/filepath"

	log "github.com/colt3k/nglog/ng"
)

type Lock struct {
	name string
}

// New create a lock file
func New(appName string) *Lock {
	name := filepath.Join(os.TempDir(), appName+".lck")
	t := new(Lock)
	t.name = name

	return t
}

// Try try to create lock file
func (l *Lock) Try() bool {
	f, err := os.OpenFile(l.name, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return false
	}
	defer f.Close()

	return true
}
func (l *Lock) Unlock() {
	err := os.Remove(l.name)
	if err != nil {
		log.Logf(log.ERROR, "issue removing lock file %+v", err)
	}
}
