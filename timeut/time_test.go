package timeut

import (
	"fmt"
	"testing"
	"time"
)

func TestMyTime_Add(t *testing.T) {
	hours := 4

	mt := GMTTime()
	fmt.Println("GMT : ", mt.Time)

	mt.Add(hours, time.Hour) // add X hours to it
	fmt.Printf("GMT Ahead %d hours: %s\n", hours, mt.Time)

	mt.Sub(hours, time.Hour) // subtract X hours to it
	fmt.Printf("GMT Subtract %d hours: %s\n", hours, mt.Time)

}
