package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var exitError *exec.ExitError

	eCom := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
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

	eCom.Env = os.Environ()

	if err := eCom.Run(); err != nil {
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
	}

	return 0
}
