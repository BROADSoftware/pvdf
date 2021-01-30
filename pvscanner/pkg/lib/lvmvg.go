package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
)

type LvmVg struct {
	VgName    string `json:"vg_name"`
	PvCount   string `json:"pv_count"`
	LvCount   string `json:"lv_count"`
	SnapCount string `json:"snap_count"`
	VgAttr    string `json:"vg_attr"`
	VgSize    string `json:"vg_size"`
	VgFree    string `json:"vg_free"`
}

type LvmVgReport struct {
	Report []struct {
		Vg []LvmVg `json:"vg"`
	} `json:"report"`
}

func getLvmVgReport() (*LvmVgReport, error) {
	args := []string { "vgs", "--unit", "b", "--reportformat", "json", "--unbuffered"}
	c:= exec.Command("/sbin/lvm", args...)
	c.Stderr = os.Stderr
	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := c.Start(); err != nil {
		return nil, err
	}
	out, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	if err := c.Wait(); err != nil {
		return nil, err
	}
	//fmt.Printf("out:%s\n", string(out))
	var vgs LvmVgReport
	err = json.Unmarshal(out, &vgs)
	if err != nil {
		return nil, err
	}
	return &vgs, nil
}

func GetVgByName() (map[string]LvmVg, error) {
	report, err := getLvmVgReport()
	if err != nil {
		return nil, err
	}
	vgByName := make(map[string]LvmVg)
	for i := 0; i < len(report.Report); i++ {
		for j := 0; j < len(report.Report[i].Vg); j++ {
			log.Debugf("LVM VolumeGroup:%s  size:%s  free:%s\n", report.Report[i].Vg[j].VgName, report.Report[i].Vg[j].VgSize, report.Report[i].Vg[j].VgFree)
			vgByName[report.Report[i].Vg[j].VgName] = report.Report[i].Vg[j]
		}
	}
	return vgByName, nil
}

