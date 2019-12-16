package utils

import (
	"os/exec"
	"sort"
)

func ExecCommand(dir, name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func RemoveDuplicatesAndEmpty(a []string) []string {
	ret := []string{}
	a_len := len(a)
	sort.Strings(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return ret
}
