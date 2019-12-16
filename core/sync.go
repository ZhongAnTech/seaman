package core

import (
	"fmt"
	"seaman/utils/config"
	"strings"

	"github.com/golang/glog"
	"seaman/services/kubecloud"
)

func Sync(g *Git) error {
	// git pull
	if err := g.Pull(); err != nil {
		glog.Errorf("git pull error: %s", err.Error())
		return err
	}

	// git diff
	toCommitId, err := g.RevParseHEAD()
	if err != nil {
		glog.Errorf("get repo latest commit id error: %s", err.Error())
		return err
	}

	companyId := config.GetConfig().Kubecloud.CompanyId
	cluster := config.GetConfig().Kubecloud.Cluster
	kc := kubecloud.NewKubecloud()
	clusterRsp, err := kc.GetClusterInfo(cluster, companyId)
	if err != nil {
		glog.Errorf("get cluster latest commit id error: %s", err.Error())
		return err
	}
	fromCommitId := clusterRsp.Data.LastCommitId

	if fromCommitId == toCommitId {
		return nil
	}

	glog.Infof("git diff from %s to %s", fromCommitId, toCommitId)
	diffFiles, err := g.Diff(fromCommitId, toCommitId)
	if err != nil {
		glog.Errorf("git diff filed, from %s to %s, error: %s", fromCommitId, toCommitId, err.Error())
		return err
	}
	glog.Infof("git diff success, file list of changes between %s and %s: %s", fromCommitId, toCommitId, diffFiles)

	// kubectl
	for _, f := range diffFiles {
		if !strings.HasSuffix(f, ".yaml") {
			continue

		}
		filePath := fmt.Sprintf("%s%s", g.dir, f)
		kubectl := NewKubectl(filePath)
		del, _ := kubectl.isDelete()
		if del {
			kubectl.Delete()
			//delFiles = append(delFiles, f)
		} else {
			kubectl.Apply()
		}
	}

	clusterRep := clusterRsp.Data
	clusterRep.LastCommitId = toCommitId
	if err := kc.UpdateCluster(cluster, companyId, clusterRep); err != nil {
		return err
	}
	glog.Infof("sync success")
	return nil
}
