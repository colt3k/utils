package qrand

import (
	"fmt"
	"testing"
)

func TestGenerateSeedData(t *testing.T) {
	i, err := GenerateSeedData(10)
	if err != nil {
		t.Errorf("issue %v", err)
	}
	fmt.Printf("vals: %v", i)
}