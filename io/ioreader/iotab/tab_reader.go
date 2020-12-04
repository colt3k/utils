package iotab

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/colt3k/utils/io/data"
)

func ReadTabDelim(path string) (*data.Table, error) {
	if !filepath.IsAbs(path) {
		tmpdir, _ := filepath.Abs(path)
		path = tmpdir
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lines := make([]string, 0)
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	buf := make([]byte, 12)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	rows := make([]data.Row, len(lines))

	for i, line := range lines {
		// split line into parts
		parts := strings.Fields(line)
		cols := make([]data.Col, len(parts))

		for j, seg := range parts {
			col := data.Col{Data: strings.TrimSpace(seg)}
			cols[j] = col
		}
		row := data.Row{Cols: cols}
		rows[i] = row
	}
	dataSet := &data.Table{Rows: rows}
	return dataSet, nil
}
