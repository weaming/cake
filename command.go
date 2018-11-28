package main

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
	"syscall"
)

const defaultFailedCode = 1

// Run a shell command
// envs format is "a=1 b=c"
func RunCommand(envs string, dir string, command string, args ...string) (exitCode int, stdout string, stderr string) {
	cmd := exec.Command(command, args...)
	// parse envs
	if envs != "" {
		envList := strings.Fields(strings.TrimSpace(envs))
		for _, env := range envList {
			cmd.Env = append(cmd.Env, env)
		}
	}
	// parse work directory
	if dir != "" {
		cmd.Dir = dir
	}

	log.Println(cmd)

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			exitCode = defaultFailedCode
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	return exitCode, stdout, stderr
}

func RunSimpleCommand(command string) (int, string, string) {
	split := strings.Fields(command)
	l := len(split)
	if l > 1 {
		return RunCommand("", "", split[0], split[1:]...)
	} else {
		return RunCommand("", "", split[0])
	}
}
