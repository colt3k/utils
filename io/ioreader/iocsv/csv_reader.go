package iocsv

import (
	"bufio"
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"

	log "github.com/colt3k/nglog/ng"
	"github.com/iancoleman/orderedmap"

	"github.com/colt3k/utils/io"
	"github.com/colt3k/utils/io/data"
)

//ReadOrderedKV reads in as an ordered Map
func ReadOrderedKV(filePathStr string) (*orderedmap.OrderedMap, error) {

	kvMap := orderedmap.New()

	if !filepath.IsAbs(filePathStr) {
		tmppath, _ := filepath.Abs(filePathStr)
		filePathStr = tmppath
	}
	var err error
	if _, err = os.Stat(filePathStr); err == nil {
		repData, _ := ReadCSV(filePathStr, false, []string{"key", "value"})
		for _, row := range repData.Rows {
			var k string
			var v string
			for _, col := range row.Cols {
				if col.Name == "key" {
					k = col.Data
				} else if col.Name == "value" {
					v = col.Data
				}
			}
			kvMap.Set(k, v)
		}

		return kvMap, nil
	}
	return nil, err
}

//ReadKV read in the file as a key/value pair and return a map of string/string
func ReadKV(filePathStr string) *map[string]string {

	var kvMap map[string]string
	kvMap = make(map[string]string)

	if !filepath.IsAbs(filePathStr) {
		tmppath, _ := filepath.Abs(filePathStr)
		filePathStr = tmppath
	}
	if _, err := os.Stat(filePathStr); err == nil {
		repData, _ := ReadCSV(filePathStr, false, []string{"key", "value"})
		for _, row := range repData.Rows {
			var k string
			var v string
			for _, col := range row.Cols {
				if col.Name == "key" {
					k = col.Data
				} else if col.Name == "value" {
					v = col.Data
				}
			}
			kvMap[k] = v
		}

		return &kvMap
	}
	return nil
}

/*
ReadCSV read a file as a CSV and return a Table
	skipHeader (true) will skip the first line for data but use it for column names
	skipHeader (false) will include all lines for data and use headAr for column names
*/
func ReadCSVFromFile(filePath string, skipHeader bool, headAr []string) (*data.Table, error) {
	if !filepath.IsAbs(filePath) {
		tmpdir, _ := filepath.Abs(filePath)
		filePath = tmpdir
	}
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rdr := bufio.NewReader(f)
	str, err := io.ReadLine(rdr)
	if err != nil {
		return nil, err
	}

	return ReadCSV(str, skipHeader, headAr)
}
func ReadCSV(fileData string, skipHeader bool, headAr []string) (*data.Table, error) {

	sr := strings.NewReader(fileData)
	r := csv.NewReader(sr)
	lines, err := r.ReadAll()
	if err != nil {
		log.Logf(log.FATAL, "error reading all lines\n%+v", err)
	}

	//Create a Row, Add Colums to it then add it to our Table
	var rows []data.Row
	if skipHeader {
		rows = make([]data.Row, len(lines)-1)
	} else {
		rows = make([]data.Row, len(lines))
	}

	var header []string
	if !skipHeader {
		header = headAr
	}

	for i, line := range lines {
		if skipHeader && i == 0 {
			//continue
			header = make([]string, len(line))
		}

		//Create our Column splice
		cols := make([]data.Col, len(line))
		//Populate Columns
		for j, seg := range line {
			if skipHeader && i == 0 {
				//Capture header fields for column names
				header[j] = seg
				continue
			}

			col := data.Col{Data: strings.TrimSpace(string(seg)), Name: strings.TrimSpace(header[j])}

			cols[j] = col
		}
		//Add Columns to Row if skipHeader is false
		if !skipHeader || i != 0 {
			row := data.Row{Cols: cols}
			if !skipHeader {
				rows[i] = row
			} else {
				rows[i-1] = row
			}
		}
	}
	dataSet := &data.Table{Rows: rows}

	return dataSet, nil
}
