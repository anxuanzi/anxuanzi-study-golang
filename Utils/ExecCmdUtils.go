package Utils

import (
	"az-ops/Global"
	"fmt"
	"github.com/go-cmd/cmd"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/pterm/pterm"
	"io/ioutil"
	"os"
	"time"
)

func ExecSimpleCmdWithResult(name string, args ...string) *cmd.Status {
	command := cmd.NewCmd(name, args...)
	execStatus := <-command.Start()
	return &execStatus
}

func ExecSimpleCmdStreaming(name string, args ...string) bool {
	// Disable output buffering, enable streaming
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	// Create Cmd with options
	envCmd := cmd.NewCmdOptions(cmdOptions, name, args...)

	// Print STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		for envCmd.Stdout != nil || envCmd.Stderr != nil {
			select {
			case line, open := <-envCmd.Stdout:
				if !open {
					envCmd.Stdout = nil
					continue
				}
				fmt.Println(line)
			case line, open := <-envCmd.Stderr:
				if !open {
					envCmd.Stderr = nil
					continue
				}
				_, _ = fmt.Fprintln(os.Stderr, line)
			}
		}
	}()

	// Run and wait for Cmd to return, discard Status
	<-envCmd.Start()

	// Wait for goroutine to print everything
	<-doneChan

	return true
}

func ExecShellScript(statikFileName string) {
	id, err := gonanoid.ID(32)
	shellFileName := id + ".sh"
	if err != nil {
		pterm.Error.Println("Nano ID error")
		pterm.Error.Println("error details: ", err)
		os.Exit(0)
	}
	shellFile, err := Global.GetStatikFS().Open("/scripts/" + statikFileName)
	if err != nil {
		pterm.Error.Println("Error when creating shell script")
		pterm.Error.Println("error details: ", err)
		os.Exit(0)
	}
	shellContent, err := ioutil.ReadAll(shellFile)
	if err != nil {
		pterm.Error.Println("Error when creating shell script")
		pterm.Error.Println("error details: ", err)
		os.Exit(0)
	}
	defer shellFile.Close()
	createdScript, err := Global.GetTmpFolder().Create(shellFileName)
	if err != nil {
		pterm.Error.Println("Error when creating shell script")
		pterm.Error.Println("error details: ", err)
		os.Exit(0)
	}
	_, err = createdScript.Write(shellContent)
	if err != nil {
		pterm.Error.Println("Error when creating shell script")
		pterm.Error.Println("error details: ", err)
		os.Exit(0)
	}
	ExecSimpleCmdWithResult("chmod", "+x", "/tmp/"+shellFileName)
	time.Sleep(2 * time.Second)
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	envCmd := cmd.NewCmdOptions(cmdOptions, "./tmp/"+shellFileName)
	// Print STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		// Done when both channels have been closed
		// https://dave.cheney.net/2013/04/30/curious-channels
		for envCmd.Stdout != nil || envCmd.Stderr != nil {
			select {
			case line, open := <-envCmd.Stdout:
				if !open {
					envCmd.Stdout = nil
					continue
				}
				fmt.Println(line)
			case line, open := <-envCmd.Stderr:
				if !open {
					envCmd.Stderr = nil
					continue
				}
				_, _ = fmt.Fprintln(os.Stderr, line)
			}
		}
	}()
	// Run and wait for Cmd to return, discard Status
	<-envCmd.Start()
	// Wait for goroutine to print everything
	<-doneChan
}
