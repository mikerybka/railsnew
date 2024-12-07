package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var dockerfile = `FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update
RUN apt upgrade -y
RUN apt install -y build-essential ruby-dev sqlite3 git libyaml-dev tzdata
RUN gem install rails
`

func main() {
	defer func() {
		os.RemoveAll("/tmp/railsnew")
	}()
	name := os.Args[1]
	err := os.MkdirAll("/tmp/railsnew", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	write("/tmp/railsnew/Dockerfile", dockerfile)
	run("/tmp/railsnew", "docker build . -t railsnew")
	run("/tmp/railsnew", fmt.Sprintf("docker run -v .:/src -w /src railsnew rails new %s", name))
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	run(".", fmt.Sprintf("cp -r /tmp/railsnew/%s %s", name, pwd))
}
func write(path, data string) {
	err := os.WriteFile(path, []byte(data), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
func run(dir, c string) {
	parts := strings.Split(c, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
