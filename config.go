package main

import (
	"github.com/gookit/config/v2"
	"os"
	"strings"
)

func loadConfig() {
	envs := map[string]string{}
	for _, env := range os.Environ() {
		s := strings.SplitN(env, "=", 2)
		envName := s[0]

		envs[envName] = strings.ReplaceAll(strings.ToLower(envName), "_", ".")
	}

	config.LoadOSEnvs(envs)
}
