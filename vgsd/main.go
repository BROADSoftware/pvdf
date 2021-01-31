package main

import (
	"flag"
	"fmt"
	"github.com/BROADSoftware/pvdf/shared/pkg/logging"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
)

var log = logging.Log.WithFields(logrus.Fields{})

func main() {
	var socketName string
	flag.StringVar(&logging.Level, "logLevel", "INFO", "Log message verbosity")
	flag.BoolVar(&logging.LogJson, "logJson", false, "logs in JSON")
	flag.StringVar(&socketName, "socketName", "/run/vgsd/vgsd.sock", "Socket name")

	logging.ConfigLogger()

	// UNIX domain socket file should be removed before listening.
	err := os.Remove(socketName)
	if err != nil && !os.IsNotExist(err) {
		panic(fmt.Sprintf("Unable to remove %s: %v", socketName, err))
	}
	log.Infof("vgsd starts. logLevel is '%s'", logging.Level)
	listener, err := net.Listen("unix", socketName)
	if err != nil {
		panic(fmt.Sprintf("Unable to listen on %s: %s", socketName, err))
	}
	mux := http.NewServeMux()
	mux.Handle("/vgs", &myHandler{ content: "Hello"})
	server := http.Server{
		Handler: mux,
	}
	log.Infof("Listen on socket '%s'", socketName)
	err = server.Serve(listener)
	log.Errorf("Ended with error:%v", err)


}

type myHandler struct {
	content string

}


func (this *myHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	err := getLvmVgReport(response)
	if err != nil && err != io.EOF {
		log.Errorf("%v", err)
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}


func getLvmVgReport(writer io.Writer) (error) {
	args := []string{"vgs", "--unit", "b", "--reportformat", "json", "--unbuffered"}
	c := exec.Command("/sbin/lvm", args...)
	c.Stderr = os.Stderr
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	if err := c.Start(); err != nil {
		return err
	}
	_, err = io.Copy(writer, stdout)
	return err
}

