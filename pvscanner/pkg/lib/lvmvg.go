package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
)


type LvmVgs struct {
	Report []struct {
		Vg []struct {
			VgName    string `json:"vg_name"`
			PvCount   string `json:"pv_count"`
			LvCount   string `json:"lv_count"`
			SnapCount string `json:"snap_count"`
			VgAttr    string `json:"vg_attr"`
			VgSize    string `json:"vg_size"`
			VgFree    string `json:"vg_free"`
		} `json:"vg"`
	} `json:"report"`
}

func GetLvmVg() (*LvmVgs, error) {
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
	var vgs LvmVgs
	err = json.Unmarshal(out, &vgs)
	if err != nil {
		return nil, err
	}
	return &vgs, nil
}