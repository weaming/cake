package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Use supervisor like config
type Program struct {
	Name            string   `yaml:"name"`      // required
	Dir             string   `yaml:"directory"` // required
	Envs            string   `yaml:"environment"`
	Command         string   `yaml:"command"` // required
	Args            []string `yaml:"arguments"`
	RedirectStderr  bool     `yaml:"redirect_stderr"`
	StdoutFile      string   `yaml:"stdout_logfile"`
	StderrFile      string   `yaml:"stderr_logfile"`
	ProcessesNumber int      `yaml:"numprocs"`
	Cron            string   `yaml:"cron"` // required
}

var (
	taskNameReplacer  = "%(program_name)s"
	defaultStdoutFile = fmt.Sprintf("~/.cake/%v.stdout.log", taskNameReplacer)
	defaultStderrFile = fmt.Sprintf("~/.cake/%v.stderr.log", taskNameReplacer)
)

func ExpandUser(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if path == "~" {
		// In case of "~", which won't be caught by the "else if"
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		path = filepath.Join(dir, path[2:])
	}
	return path
}

func ParseProgramConfig(path string) (programs []*Program) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("read file err #%v\n", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, &programs)
	if err != nil {
		log.Printf("unmarshal: %v\n", err)
		os.Exit(2)
	}

	for _, p := range programs {
		// check for constraints
		if p.Name == "" {
			log.Println("name is required")
			os.Exit(3)
		}
		if p.Dir == "" {
			log.Println("directory is required")
			os.Exit(3)
		}
		if p.Command == "" {
			log.Println("command is required")
			os.Exit(3)
		}
		if p.Cron == "" {
			log.Println("cron is required")
			os.Exit(3)
		}
		if p.StdoutFile == "" {
			p.StdoutFile = defaultStdoutFile
		}
		if p.StderrFile == "" {
			p.StderrFile = defaultStderrFile
		}

		// parse string
		if strings.Contains(p.Dir, taskNameReplacer) {
			p.Dir = strings.Replace(p.Dir, taskNameReplacer, p.Name, -1)
		}
		if strings.Contains(p.Envs, taskNameReplacer) {
			p.Envs = strings.Replace(p.Envs, taskNameReplacer, p.Name, -1)
		}
		if strings.Contains(p.Command, taskNameReplacer) {
			p.Command = strings.Replace(p.Command, taskNameReplacer, p.Name, -1)
		}
		for i, arg := range p.Args {
			if strings.Contains(arg, taskNameReplacer) {
				p.Args[i] = strings.Replace(arg, taskNameReplacer, p.Name, -1)
			}
		}
		if strings.Contains(p.StdoutFile, taskNameReplacer) {
			p.StdoutFile = strings.Replace(p.StdoutFile, taskNameReplacer, p.Name, -1)
		}
		if strings.Contains(p.StderrFile, taskNameReplacer) {
			p.StderrFile = strings.Replace(p.StderrFile, taskNameReplacer, p.Name, -1)
		}

		// expand user
		p.StdoutFile = ExpandUser(p.StdoutFile)
		p.StderrFile = ExpandUser(p.StderrFile)
		p.Dir = ExpandUser(p.Dir)
	}
	return
}
