module github.com/BROADSoftware/pvdf/pvscanner

go 1.15

replace github.com/BROADSoftware/pvdf/shared v0.1.0 => ../shared

require (
	github.com/BROADSoftware/pvdf/shared v0.1.0
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.3
)
