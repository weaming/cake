package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func PrepareDir(filePath string, force bool) {
	if !strings.HasSuffix(filePath, "/") || force {
		filePath = path.Dir(filePath)
	}
	err := os.MkdirAll(filePath, os.FileMode(0755))
	if err != nil {
		log.Fatal(err)
	}
}

func WriteLog(path, line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	PrepareDir(path, true)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	fmt.Fprintln(f, line)
}
