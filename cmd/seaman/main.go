package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"runtime"
	"time"

	"seaman/api"
	"seaman/core"
	"seaman/utils/config"
)

const (
	version = "v1.0.0"
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	defer glog.Flush()

	glog.Infof("Seaman version: %s", version)
	glog.Infof("Golang version: %s", runtime.Version())
	glog.Infof("Gin version: %s", gin.Version)

	go func() {
		dir := config.GetConfig().Git.Dir
		url := config.GetConfig().Git.URL
		branch := config.GetConfig().Git.Branch
		token := config.GetConfig().Git.Token
		git := core.NewGit(dir, url, branch, token)
		if err := git.Clone(); err != nil {
			glog.Fatal(err.Error())
		}

		syncTimeSecond := config.GetConfig().Sync.Second
		for {
			if err := core.Sync(git); err != nil {
				glog.Errorf("sync failed: %s", err.Error())
			}
			time.Sleep(syncTimeSecond * time.Second)
		}
	}()

	gin.SetMode(gin.DebugMode)
	glog.Fatal(api.Router().Run("0.0.0.0:8080"))
}
