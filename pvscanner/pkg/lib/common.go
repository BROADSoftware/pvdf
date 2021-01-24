package lib

import (
	"github.com/BROADSoftware/pvdf/shared/pkg/logging"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"time"
)

var log = logging.Log.WithFields(logrus.Fields{})


//  Cli Parameters. Check main.go
var ProcPath string
var RootfsPath string
var StatfsTimeout time.Duration
var Period time.Duration



func procFilePath(name string) string {
	return filepath.Join(ProcPath, name)
}

func rootfsFilePath(name string) string {
	return filepath.Join(RootfsPath, name)
}

