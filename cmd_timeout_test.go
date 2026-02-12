package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommandTimeout(t *testing.T) {
	// Build the dockerize binary first
	buildCmd := exec.Command("go", "build", "-o", "dockerize-test", ".")
	err := buildCmd.Run()
	assert.NoError(t, err)
	defer os.Remove("dockerize-test")

	// Test that a command times out and shows the timeout message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "./dockerize-test", "-cmd-timeout", "1s", "sleep", "10")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	assert.Error(t, err, "Expected command to fail due to timeout")

	stderrOutput := stderr.String()
	assert.Contains(t, stderrOutput, "Killing command due to timeout", "Expected timeout message in stderr")
	assert.Contains(t, stderrOutput, "1s", "Expected timeout duration in message")
}

func TestCommandWithoutTimeout(t *testing.T) {
	// Build the dockerize binary first
	buildCmd := exec.Command("go", "build", "-o", "dockerize-test", ".")
	err := buildCmd.Run()
	assert.NoError(t, err)
	defer os.Remove("dockerize-test")

	// Test that a command without timeout executes normally
	cmd := exec.Command("./dockerize-test", "echo", "test")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err = cmd.Run()
	assert.NoError(t, err, "Expected command to succeed")

	stdoutOutput := strings.TrimSpace(stdout.String())
	assert.Equal(t, "test", stdoutOutput, "Expected echo output")
}

func TestCommandTimeoutZero(t *testing.T) {
	// Build the dockerize binary first
	buildCmd := exec.Command("go", "build", "-o", "dockerize-test", ".")
	err := buildCmd.Run()
	assert.NoError(t, err)
	defer os.Remove("dockerize-test")

	// Test that -cmd-timeout 0 (default) means no timeout
	cmd := exec.Command("./dockerize-test", "-cmd-timeout", "0", "echo", "no-timeout")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err = cmd.Run()
	assert.NoError(t, err, "Expected command to succeed with zero timeout")

	stdoutOutput := strings.TrimSpace(stdout.String())
	assert.Equal(t, "no-timeout", stdoutOutput, "Expected echo output")
}
