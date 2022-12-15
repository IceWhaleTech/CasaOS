package command

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	println(cmd.String())

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

// 执行 lsblk 命令
func ExecLSBLK() []byte {
	output, err := exec.Command("lsblk", "-O", "-J", "-b").Output()
	if err != nil {
		fmt.Println("lsblk", err)
		return nil
	}
	return output
}

// 执行 lsblk 命令
func ExecLSBLKByPath(path string) []byte {
	output, err := exec.Command("lsblk", path, "-O", "-J", "-b").Output()
	if err != nil {
		fmt.Println("lsblk", err)
		return nil
	}
	return output
}

// exec smart
func ExecSmartCTLByPath(path string) []byte {
	timeout := 3
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, "smartctl", "-a", path, "-j").Output()
	if err != nil {
		fmt.Println("smartctl", err)
		return nil
	}
	return output
}

func ExecEnabledSMART(path string) {
	exec.Command("smartctl", "-s on", path).Output()
}

func ExecuteScripts(scriptDirectory string) {
	if _, err := os.Stat(scriptDirectory); os.IsNotExist(err) {
		fmt.Printf("No post-start scripts at %s\n", scriptDirectory)
		return
	}

	files, err := os.ReadDir(scriptDirectory)
	if err != nil {
		fmt.Printf("Failed to read from script directory %s: %s\n", scriptDirectory, err.Error())
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		scriptFilepath := filepath.Join(scriptDirectory, file.Name())

		f, err := os.Open(scriptFilepath)
		if err != nil {
			fmt.Printf("Failed to open script file %s: %s\n", scriptFilepath, err.Error())
			continue
		}
		f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Scan()
		shebang := scanner.Text()

		interpreter := "/bin/sh"
		if strings.HasPrefix(shebang, "#!") {
			interpreter = shebang[2:]
		}

		cmd := exec.Command(interpreter, scriptFilepath)

		fmt.Printf("Executing post-start script %s using %s\n", scriptFilepath, interpreter)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Printf("Failed to execute post-start script %s: %s\n", scriptFilepath, err.Error())
		}
	}
	fmt.Println("Finished executing post-start scripts.")
}
