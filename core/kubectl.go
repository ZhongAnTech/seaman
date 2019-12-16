package core

import (
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"path/filepath"
	"seaman/utils"
	"strings"
)

type Kubectl struct {
	command  string
	filePath string
}

func NewKubectl(filePath string) *Kubectl {
	return &Kubectl{
		command:  "kubectl",
		filePath: filePath,
	}
}

func (k *Kubectl) ApplyOrDelete() error {
	del, err := k.isDelete()
	if err != nil {
		return err
	}
	if del {
		return k.Delete()
	}
	return k.Apply()
}

func (k *Kubectl) Apply() error {
	dir, _ := filepath.Abs(`.`)
	args := []string{"--kubeconfig=kubeconfig/config", "apply", "-f", k.filePath}
	out, err := utils.ExecCommand(dir, k.command, args...)
	if err != nil {
		glog.Errorf("%s %s, exec failed: %s", k.command, args, out)
		return err
	}
	glog.Infof("%s %s, exec success: %s", k.command, args, out)
	return nil
}

func (k *Kubectl) Delete() error {
	dir, _ := filepath.Abs(`.`)
	args := []string{"--kubeconfig=kubeconfig/config", "delete", "-f", k.filePath}
	out, err := utils.ExecCommand(dir, k.command, args...)
	if err != nil {
		glog.Errorf("%s %s, exec failed: %s", k.command, args, out)
		return err
	}
	glog.Infof("%s %s, exec success: %s", k.command, args, out)
	return nil
}

func (k *Kubectl) isDelete() (bool, error) {
	yamlFile, err := os.OpenFile(k.filePath, os.O_RDONLY, 0644)
	if err != nil {
		glog.Errorf("Failed to open the file: %s", err.Error())
		return false, err
	}
	defer yamlFile.Close()
	contents, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		glog.Errorf("Failed to read the file: %s", err.Error())
		return false, err
	}
	del := strings.Contains(string(contents), `kubecloud/delete: "true"`)
	return del, nil
}
