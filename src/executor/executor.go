package executor

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunCmd(cmd string) (string, error) {
	// log.Println("RUN: ", cmd)
	stdout, stderr, err := Shellout(cmd)
	if err != nil {
		res := fmt.Sprintf("cmd error: %s\n %s\n %s", stdout, stderr, err)
		log.Println(res)
		return res, err
	}
	return stdout + stderr, err
}

func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func Run(command, args string) (stdout bytes.Buffer, stderr bytes.Buffer, err error) {
	log.Println("Run: ", command, args)
	cmd := exec.Command(command, strings.Split(args, " ")...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return
}

func RunAndProcess(comand, args string) (err error) {
	cmd := exec.Command(comand, strings.Split(args, " ")...)
	stdout, _ := cmd.StdoutPipe() // StderrPipe
	err = cmd.Start()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		log.Println("XXX:", m)
	}
	log.Println("scanner done1")
	err = cmd.Wait()
	return
}

func RunAndProcessStderr(comand, args string) (err error) {
	cmd := exec.Command(comand, strings.Split(args, " ")...)
	stderr, _ := cmd.StderrPipe() // StderrPipe
	err = cmd.Start()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		log.Println("XXX:", m)
	}
	log.Println("scanner done1")
	err = cmd.Wait()
	return
}

func FileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
