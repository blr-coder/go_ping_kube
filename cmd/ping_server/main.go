package main

import (
	"github.com/sirupsen/logrus"
	"go_ping_kube/pkg/log"
)

func main() {
	log.InitLogger()
	logrus.Info("!!! START !!!")

	if err := runPingApp(); err != nil {
		logrus.Fatal(err)
	}
}
