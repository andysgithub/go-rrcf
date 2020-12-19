package logging

import (
	"bytes"
	"fmt"
	"os"
)

// Logger - Logging utility
type Logger struct {
	logMain    *os.File
	bufferMain bytes.Buffer
}

// NewLogger returns a new logger
func NewLogger(fileName string) *Logger {
	logger := Logger{}

	os.Remove(fileName)
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		logger = Logger{
			logMain: file,
		}
	}
	return &logger
}

// Clear -
func (log *Logger) Clear() {
	log.bufferMain.Reset()
}

// Section -
func (log *Logger) Section(text string) {
	log.bufferMain.Reset()
	log.Add(text)
}

// Add -
func (log *Logger) Add(text string) {
	log.bufferMain.WriteString(text)
}

// AddLast -
func (log *Logger) AddLast(text string) {
	log.Add(text)
	log.Write()
}

// Write -
func (log *Logger) Write() {
	fmt.Fprintln(log.logMain, log.bufferMain.String())
}

// Close -
func (log *Logger) Close() {
	log.logMain.Close()
}
