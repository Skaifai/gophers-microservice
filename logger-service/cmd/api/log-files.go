package main

import (
	"fmt"
	"os"
	"time"
)

func WriteToFile(line string, file *os.File) {
	newLine := fmt.Sprintf("%v::%s\n", time.Now().Format(logHeaderLayout), line)
	_, err := file.WriteString(newLine)
	failOnError(err, "could not write to file")
}

func CreateLogFile() *os.File {
	fileName := fmt.Sprintf("log-files/logs-%v.txt", time.Now().Format(fileNameLayout))
	f, err := os.Create(fileName)
	failOnError(err, "Could not create file")
	return f
}
