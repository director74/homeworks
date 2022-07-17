package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	eCom := exec.Command(cmd[0], cmd[1:]...)
	eCom.Stdin = os.Stdin
	eCom.Stdout = os.Stdout
	eCom.Stderr = os.Stderr

	for key, envItem := range env {
		if envItem.NeedRemove {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, envItem.Value)
		}
	}

	eCom.Env = append(os.Environ())

	if err := eCom.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
	}

	return
}
