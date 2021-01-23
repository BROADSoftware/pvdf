package lib

import (
	"github.com/BROADSoftware/pvdf/pvscanner/pkg/logging"
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


var FreeAnnotation = "pvscanner.pvdf.broadsoftware.com/free_mib"
var SizeAnnotation = "pvscanner.pvdf.broadsoftware.com/size_mib"


func procFilePath(name string) string {
	return filepath.Join(ProcPath, name)
}

func rootfsFilePath(name string) string {
	return filepath.Join(RootfsPath, name)
}

