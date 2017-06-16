package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"io"
	"os/exec"
)

func execute(name string, arg ...string) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	cmd := exec.Command(name, arg...)

	stdin, err := cmd.StdinPipe()

	if err != nil {
		fmt.Println(err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Start(); err != nil {
		fmt.Println("An error occured: ", err)
	}

	stdin.Close()
	cmd.Wait()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	return string(out)
}

func execute_pipe(in []string, name string, arg ...string) string {
        rescueStdout := os.Stdout
        r, w, _ := os.Pipe()
        os.Stdout = w
        os.Stderr = w

        cmd := exec.Command(name, arg...)

        stdin, err := cmd.StdinPipe()

        if err != nil {
                fmt.Println(err)
        }

        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        if err = cmd.Start(); err != nil {
                fmt.Println("An error occured: ", err)
        }

        go func() {
                defer stdin.Close()
		if len(in) > 0 {
			io.WriteString(stdin, in[0])
		}
		for i := 1; i < len(in); i++ {
                	io.WriteString(stdin, "\n" + in[i])
		}
        }()

        cmd.Wait()

        w.Close()
        out, _ := ioutil.ReadAll(r)
        os.Stdout = rescueStdout

        return string(out)
}
