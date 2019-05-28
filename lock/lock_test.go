package lock

import (
	"log"
	"testing"
)

func ExampleNew() {
	l := New("testme")
	if l.Try() {
		// do work
	}
	// release lock
	defer l.Unlock()
}

func TestNew(t *testing.T) {
	l := New("testme")
	log.Println("Path to file", l.name)
	if l.Try() {
		log.Println("we have the lock lets do some work")
		defer l.Unlock()
	}
}