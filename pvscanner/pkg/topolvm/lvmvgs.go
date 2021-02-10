package topolvm

import (
	"encoding/json"
	"github.com/BROADSoftware/pvdf/pvscanner/pkg/lib"
	"io/ioutil"
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


// wrapExecCommand calls cmd with args but wrapped to run on the host
func wrapExecCommand(cmd string, args ...string) *exec.Cmd {
	if Containerized {
		args = append([]string{"--root=" + lib.RootfsPath, "-t", "1", cmd}, args...)
		cmd = Nsenter
	}
	c := exec.Command(cmd, args...)
	return c
}

func getLvmVgReport() (*LvmVgReport, error) {
	args := []string { "vgs", "--unit", "b", "--reportformat", "json", "--unbuffered"}
	c := wrapExecCommand(Lvm, args...)
	stderr, err := c.StderrPipe()
	if err != nil {
		return nil, err
	}
	//c.Stderr = os.Stderr
	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := c.Start(); err != nil {
		return nil, err
	}
	stderrData, err := ioutil.ReadAll(stderr)
	if err != nil {
		return nil, err
	}
	out, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	if err := c.Wait(); err != nil {
		return nil, err
	}
	if len(stderrData) > 0 {
		log.Debugf("Stderr:%s", stderrData)
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

