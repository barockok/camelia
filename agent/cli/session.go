package main

import (
	"os"

	"sync"

	"fmt"

	"github.com/barockok/camelia/agent"
	"github.com/barockok/camelia/misc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type backupJOB struct {
	state        string
	log          *logrus.Entry
	restartCount int64
}

func (bj *backupJOB) run() {
	defer func() {
		if r := recover(); r != nil {
			bj.log.WithError(r.(error)).Error("Panic while running")
		}

		if err := bj.restart(); err != nil {
			panic("Could not restart backup-job")
		}
	}()

	bjConcurrent := viper.GetInt("concurrency.backup_job")
	bj.log.Infof("About to start backup job with %d concurrencies", bjConcurrent)

	uploader := getUploader()
	kfPool := []agent.KFAgent{}

	var wg sync.WaitGroup
	for i := 0; i < bjConcurrent; i++ {
		wg.Add(1)
		go func() {
			agent.NewKFAgent()
		}()
	}

	bj.state = "RUNNING"
	wg.Wait()
}

func (bj *backupJOB) stop() {
	bj.state = "STOP"
}

func (bj *backupJOB) restart() error {
	bj.state = "RESTARTING"

	bj.log.Infof("About to restarting")

	return nil
}

func (bj *backupJOB) isRunning() bool {
	return bj.state != "STOP"
}

func newBackupJob(restartCount int64) *backupJOB {
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	log := logrus.WithFields(logrus.Fields{"host": hostName, "job": "background-job"})
	return &backupJOB{log: log, restartCount: restartCount}
}

func getUploader() misc.Uploader {
	switch driver := viper.GetString("backup.uploader"); driver {
	case "S3":
		u := agent.UploaderS3{}
		u.Configure(viper.GetStringMapString("uploader.s3"))
		return &u
	default:
		panic(fmt.Errorf("Unknow Uploader driver %s", driver))
	}

}
