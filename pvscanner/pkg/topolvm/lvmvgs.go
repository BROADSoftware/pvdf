package topolvm

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
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


func NewVgsClient(socketName string) http.Client {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketName, )
			},
		},
	}
	return httpc
}

func getLvmVgReport(vgsClient http.Client) (*LvmVgReport, error) {
	response, err := vgsClient.Get("http://localhost/vgs")
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Received error %d from vgsd daemon", response.StatusCode)
	}
	var vgs LvmVgReport
	err = json.NewDecoder(response.Body).Decode(&vgs)
	if err != nil {
		return nil, err
	}
	return &vgs, nil
}

func GetVgByName(vgsClient http.Client) (map[string]LvmVg, error) {
	report, err := getLvmVgReport(vgsClient)
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

