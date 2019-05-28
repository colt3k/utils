package asciipb

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Requires replacing token {n} with a number to move that many locations
type CursorNav string

const CursorNavPrefix = "\u001b["
const (
	UP    = "A"
	DOWN  = "B"
	RIGHT = "C"
	LEFT  = "D"
	NEXTLINE = "E"
	PREVLINE = "F"
	SETCOLUMN = "G"
	SETPOSITION = "H" //uses row and column \u001b[{n};{m}H
	CLRSCREEN = "J"		// 0..cursor until end of screen, 1.. cursor to beginning of screen, 2..entire screen
	CLEARLINE = "K"		// int on this 0..cursor to end of line, 1..cursor to start of line, 2.. entire line
)

func ProgressIndicator(size int) {

	for i := 0; i < size; i++ {
		time.Sleep(5 * time.Millisecond)
		t := CursorNavPrefix + strconv.Itoa(1000) + LEFT
		os.Stdout.Write([]byte(fmt.Sprintf("%s%s%%", t, strconv.Itoa(i+1))))
	}
	os.Stdout.Write([]byte("\n"))
}

func ProgressBarIndicator(size int) {

	for i := 0; i < size; i++ {
		time.Sleep(5 * time.Millisecond)
		width := (i + 1) / 4
		// printout created here
		hash := strings.Repeat("#", width)
		spcs := strings.Repeat(" ", (25 - width))
		bar := "[" + hash + spcs + "]"

		t := CursorNavPrefix + strconv.Itoa(1000) + LEFT
		os.Stdout.Write([]byte(fmt.Sprintf("%s%s", t, bar)))
	}
	os.Stdout.Write([]byte("\n"))
}

type multi struct {
	ttl  int
	done bool
}

func MultiProgressBarIndicator(size int) {

	count := size
	// setup
	dar := make([]multi, count)
	for i := range dar {
		t := multi{ttl: 0, done: false}
		dar[i] = t
	}

	allProgress := count
	nl := strings.Repeat("\n", allProgress)
	// ensure we have space to draw our bars
	os.Stdout.Write([]byte(nl))

	for {
		allComplete := true
		for _, d := range dar {
			if !d.done {
				allComplete = false
			}
		}
		if allComplete {
			break
		}
		time.Sleep(5 * time.Millisecond)
		idx := pickIdx(dar)
		dar[idx].ttl += 1
		if dar[idx].ttl == 100 {
			dar[idx].done = true
		}
		// move left 1k
		t := CursorNavPrefix + strconv.Itoa(1000) + LEFT
		os.Stdout.Write([]byte(t))
		// move up count i.e. 2
		os.Stdout.Write([]byte(CursorNavPrefix + strconv.Itoa(count) + UP))
		for _, d := range dar {
			width := d.ttl / 4
			hash := strings.Repeat("#", width)
			spcs := strings.Repeat(" ", 25-width)
			bar := "[" + hash + spcs + "]"
			os.Stdout.Write([]byte(fmt.Sprintf("%s", bar)))
			os.Stdout.Write([]byte("\n"))
		}
	}

	os.Stdout.Write([]byte("\n"))
}

func pickIdx(dar []multi) int {
	rand.Seed(time.Now().Unix())
	idx := rand.Intn(len(dar))
	d := dar[idx]
	// if done pick another until we found an incomplete one
	if d.done {
		return pickIdx(dar)
	}
	return idx
}
