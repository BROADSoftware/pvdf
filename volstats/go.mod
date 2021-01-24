module github.com/BROADSoftware/pvdf/volstats

go 1.15

replace github.com/BROADSoftware/pvdf/shared v0.1.0 => ../shared

require (
	github.com/BROADSoftware/pvdf/shared v0.1.0
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.1.1
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.3
)
