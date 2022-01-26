package command

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"
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
	//str, err := ioutil.ReadAll(stdout)
	var networklist = []string{}
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

//执行 lsblk 命令
func ExecLSBLK() []byte {
	output, err := exec.Command("lsblk", "-O", "-J", "-b").Output()
	if err != nil {
		fmt.Println("lsblk", err)
		return nil
	}
	return output
}

//执行 lsblk 命令
func ExecLSBLKByPath(path string) []byte {
	output, err := exec.Command("lsblk", path, "-O", "-J", "-b").Output()
	if err != nil {
		fmt.Println("lsblk", err)
		return nil
	}
	return output
}

//exec smart
func ExecSmartCTLByPath(path string) []byte {
	timeout := 3
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, "smartctl", "-a", path, "-j").Output()
	if err != nil {
		fmt.Println("smartctl", err)
		fmt.Println("output", string(output))
		return nil
	}
	return output
}

func ExecEnabledSMART(path string) {

	exec.Command("smartctl", "-s on", path).Output()
}
