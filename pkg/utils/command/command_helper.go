package command

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os/exec"
)

func OnlyExec(cmdStr string) {
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		return
	}
	cmd.Wait()
	return
}

func ExecResultStrArray(cmdStr string) []string {
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer stdout.Close()
	if err = cmd.Start(); err != nil {
		fmt.Println(err)
		return nil
	}
	// str, err := ioutil.ReadAll(stdout)
	networklist := []string{}
	outputBuf := bufio.NewReader(stdout)
	for {
		output, _, err := outputBuf.ReadLine()
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Printf("Error :%s\n", err)
			}
			break
		}
		networklist = append(networklist, string(output))
	}
	cmd.Wait()
	return networklist
}

func ExecResultStr(cmdStr string) string {
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return ""
	}
	str, err := ioutil.ReadAll(stdout)
	cmd.Wait()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(str)
}
