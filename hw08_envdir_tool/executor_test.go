package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("check return code", func(t *testing.T) {
		commands := make([]string, 0)
		commands = append(commands, "ls")
		commands = append(commands, "unknown_*")
		envs := make(map[string]EnvValue)

		retCode := RunCmd(commands, envs)

		require.Equal(t, 2, retCode)
	})

	t.Run("check installed env", func(t *testing.T) {
		rescueStdout := os.Stdout
		reader, writer, _ := os.Pipe()
		os.Stdout = writer

		commands := make([]string, 0)
		commands = append(commands, "env")
		envs := make(map[string]EnvValue)

		envs["HM"] = EnvValue{Value: "Hello", NeedRemove: false}

		retCode := RunCmd(commands, envs)

		writer.Close()
		out, _ := io.ReadAll(reader)
		os.Stdout = rescueStdout

		require.True(t, strings.Contains(string(out), "HM=Hello"))
		require.Equal(t, 0, retCode)
	})
}
