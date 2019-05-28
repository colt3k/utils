package store

import (
	"testing"

	log "github.com/colt3k/nglog/ng"
)

func TestNewMVKeySet(t *testing.T) {

	mks := NewMVKeySet()

	mks.Add("1", "test")
	if mks.ContainsKey("1") {
		log.Logf(log.INFO, "Contains Key 1")
	}
	if mks.ContainsVal("1", "test") {
		log.Logf(log.INFO, "Contains Val: test")
	}
	if !mks.ContainsVal("1", "test2") {
		log.Logf(log.INFO, "DOES NOT Contain Val: test2")
	}

	mks.Add("1", "test2")
	mks.Add("1", "test2")
	if mks.ContainsVal("1", "test2") {
		log.Logf(log.INFO, "Contains Val: test2")
	}

	mks.Add("1", "test3")
	mks.Add("2", "yo")
	mks.Add("2", "yo")
	mks.Add("2", "dude")

	mks.RemoveKey("1")
	log.Logf(log.INFO, "Data %v", mks.Sets)
	mks.RemoveVal("2", "yo")
	log.Logf(log.INFO, "Data %v", mks.Sets)

	log.Logf(log.INFO, "size keys: %v", mks.SizeKeys())
	log.Logf(log.INFO, "size vals: %v", mks.SizeVals("2"))
	mks.Add(1, 890)

}
