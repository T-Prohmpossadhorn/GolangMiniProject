package db

import (
	"bufio"
	"log"
	"os"
)

type DB struct {
	filepath string
}

func New(path string) *DB {
	return &DB{
		filepath: path,
	}
}

func (d *DB) Getfulllist() []string {
	readFile, err := os.Open(d.filepath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer readFile.Close()

	var ret []string

	scanner := bufio.NewScanner(readFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil
	}

	return ret
}

func (d *DB) Getfirstfilter(filter string) string {
	return ""
}
