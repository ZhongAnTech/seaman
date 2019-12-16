package core

import (
	"fmt"
	"github.com/golang/glog"
	"path/filepath"
	"seaman/utils"
	"strings"
)

type Git struct {
	commandName string
	dir         string
	url         string
	branch      string
	token       string
}

func NewGit(dir, url, branch, token string) *Git {
	return &Git{
		commandName: "git",
		dir:         dir,
		url:         url,
		branch:      branch,
		token:       token,
	}
}

func (g *Git) configUserEmail() error {
	glog.Infof(`git config user.name "seaman"`)
	args := []string{"config", "user.email", `"seaman@zhongan.io"`}
	out, err := utils.ExecCommand(g.dir, g.commandName, args...)
	if err != nil {
		glog.Error(out)
		return fmt.Errorf(out)
	}
	glog.Infof(`git config user.name "seaman"`)
	args = []string{"config", "user.name", `"seaman"`}
	out, err = utils.ExecCommand(g.dir, "git", args...)
	if err != nil {
		glog.Error(out)
		return fmt.Errorf(out)
	}
	return nil
}

func (g *Git) Clone() error {
	currentDir, _ := filepath.Abs(`.`)
	urlWithToken := strings.ReplaceAll(g.url, "//", fmt.Sprintf("//:%s@", g.token))
	glog.Infof("git clone %s -b %s", g.url, g.branch)
	args := []string{"clone", urlWithToken, "-b", g.branch, g.dir}
	out, err := utils.ExecCommand(currentDir, g.commandName, args...)
	if err != nil && !strings.Contains(out, "already exists and is not an empty directory") {
		glog.Error(out)
		return fmt.Errorf(out)
	}
	glog.Infoln(out)
	if err := g.configUserEmail(); err != nil {
		return fmt.Errorf(out)
	}
	return nil
}

func (g *Git) RevParseHEAD() (string, error) {
	//glog.Infof("git rev-parse HEAD")
	args := []string{"rev-parse", "HEAD"}
	out, err := utils.ExecCommand(g.dir, g.commandName, args...)
	if err != nil {
		glog.Infof("git rev-parse HEAD")
		glog.Error(out)
		return "", fmt.Errorf(out)
	}
	//glog.Infoln(out)
	return out, nil
}

func (g *Git) RmCommit(files []string, msg string) error {
	for _, f := range files {
		glog.Infof("git rm %s", f)
		args := []string{"rm", f}
		out, err := utils.ExecCommand(g.dir, g.commandName, args...)
		if err != nil {
			glog.Error(out)
			return fmt.Errorf(out)
		}
	}
	args := []string{"commit", "-m", fmt.Sprintf(`"%s"`, msg)}
	out, err := utils.ExecCommand(g.dir, g.commandName, args...)
	if err != nil && !strings.Contains(out, "nothing to commit, working tree clean") {
		glog.Errorf("git commit -m  %s error: %s", fmt.Sprintf(`"%s"`, msg), out)
		return fmt.Errorf(out)
	}
	glog.Infof("git commit -m %s", fmt.Sprintf(`"%s"`, msg))
	glog.Infoln(out)
	return nil
}
func (g *Git) Pull() error {
	glog.Infof("git pull origin %s", g.branch)
	args := []string{"pull", "origin", g.branch}
	out, err := utils.ExecCommand(g.dir, g.commandName, args...)
	if err != nil {
		glog.Error(out)
		return fmt.Errorf(out)
	}
	glog.Infoln(out)
	return nil
}

func (g *Git) Push() error {
	glog.Infof("git push origin %s", g.branch)
	args := []string{"push", "origin", g.branch}
	out, err := utils.ExecCommand(g.dir, g.commandName, args...)
	if err != nil {
		glog.Error(out)
		return fmt.Errorf(out)
	}
	glog.Infoln(out)
	return nil
}

func (g *Git) Diff(from, to string) ([]string, error) {
	from = strings.Replace(from, "\n", "", -1)
	to = strings.Replace(to, "\n", "", -1)
	glog.Infof("git diff --name-status %s %s", from, to)
	args := []string{"diff", "--name-status", from, to}
	out, err := utils.ExecCommand(g.dir, g.commandName, args...)
	if err != nil {
		glog.Error(out)
		return nil, fmt.Errorf(out)
	}
	glog.Infoln(out)

	diffFiles := []string{}
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		change := strings.Split(line, "\t")
		if len(change) != 2 {
			glog.Errorf("invaild line: %s", line)
			continue
		}
		m := change[0]
		f := change[1]
		if m == "A" || m == "M" {
			diffFiles = append(diffFiles, f)
		}
	}
	diffFiles = utils.RemoveDuplicatesAndEmpty(diffFiles)
	return diffFiles, err
}
