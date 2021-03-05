package qrand

import (
	"testing"
)

func TestGenerateSeedData(t *testing.T) {
	i, err := GenerateSeedData(300)
	if err != nil {
		t.Errorf("issue %v", err)
		t.Failed()
	}
	t.Logf("vals: %v", i)

}