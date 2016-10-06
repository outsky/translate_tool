package log_test

import (
	"testing"
	"trans/log"
)

func Test_example(t *testing.T) {
	defer log.CloseLog()
	log.Info("fc", "test, type: Info, dest: console & file")
	log.Info("f", "test, type: Info, dest: file")
	log.Info("c", "test, type: Info, dest: console")

	log.Error("fc", "test, type: Error, dest: console & file")
	log.Error("f", "test, type: Error, dest: file")
	log.Error("c", "test, type: Error, dest: console")
}
