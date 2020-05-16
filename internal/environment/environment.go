package environment

import (
	"fmt"
	"os"
	"strings"
)

////////////////////////////////////////////////////////////////
//// Environment
////

type Environment map[string]string

func (env Environment) Slice() (sl []string) {
	for k, v := range env {
		sl = append(sl, fmt.Sprintf("%s=%s", k, v))
	}
	return sl
}

func (env Environment) Set(k, v string) Environment {
	env[k] = v
	return env
}

func Environ() Environment {
	env := Environment{}

	for _, e := range os.Environ() {
		i := strings.Index(e, "=")
		if i >= 0 {
			env[e[:i]] = e[i+1:]
		}
	}

	return env
}
