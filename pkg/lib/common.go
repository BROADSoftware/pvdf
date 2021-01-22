package lib

import (
	"github.com/BROADSoftware/pvdf/pkg/logging"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"time"
)

var log = logging.Log.WithFields(logrus.Fields{})


// TODO: As parameters
var procPath = "/proc"
var rootfsPath string = "/"
var statfsTimeout, _ = time.ParseDuration("5s")
var Period, _ = time.ParseDuration("10s")

func procFilePath(name string) string {
	return filepath.Join(procPath, name)
}

func rootfsFilePath(name string) string {
	return filepath.Join(rootfsPath, name)
}

