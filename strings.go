package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
)

func removeDirSizesEntry() {
	// NOTE: path to file to match on
	fpath := "./file"

	f, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// NOTE: pattern to match string matched lines are removed
	pattern := "string"

	var bs []byte
	buf := bytes.NewBuffer(bs)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// NOTE: the not operator
		if !strings.Contains(scanner.Text(), pattern) {
			_, err := buf.Write(scanner.Bytes())
			if err != nil {
				log.Fatal(err)
			}
			_, err = buf.WriteString("\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(fpath, buf.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
