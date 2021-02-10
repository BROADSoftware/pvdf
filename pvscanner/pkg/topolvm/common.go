package topolvm

import (
	"github.com/BROADSoftware/pvdf/shared/pkg/logging"
	"github.com/sirupsen/logrus"
)

var log = logging.Log.WithFields(logrus.Fields{})

var Containerized bool
var Nsenter string
var Lvm string
