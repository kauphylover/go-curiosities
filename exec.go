package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/coreos/go-semver/semver"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	cmd  string
	args []string
}

//func main() {
//	// Collect directories from the command-line
//	var dirs []string
//	if len(os.Args) > 1 {
//		dirs = os.Args[1:]
//	} else {
//		dirs = []string{"."}
//	}
//
//	// Run the command on each directory
//	for _, dir := range dirs {
//		// find $DIR -type f # Find all files
//		ls := exec.Command("find", dir, "-type", "f")
//
//		// | grep -v '/[._]' # Ignore hidden/temporary files
//		visible := exec.Command("egrep", "-v", `/[._]`)
//
//		// | sort -t. -k2 # Sort by file extension
//		sort := exec.Command("sort", "-t.", "-k2")
//
//		// Run the pipeline
//		output, stderr, err := Pipeline(ls, visible, sort)
//		if err != nil {
//			fmt.Printf("dir %q: %s", dir, err)
//		}
//
//		// Print the stdout, if any
//		if len(output) > 0 {
//			fmt.Printf("%q:\n%s", dir, output)
//		}
//
//		// Print the stderr, if any
//		if len(stderr) > 0 {
//			fmt.Printf("%q: (stderr)\n%s", dir, stderr)
//		}
//	}
//}

func main() {
	git := exec.Command("git", "ls-remote", "-t", "--sort", "v:refname", "git@gitlab.eng.vmware.com:nsx-allspark_users/common-apis.git")
	//git := exec.Command("ls", "-lrt")
	tail := exec.Command("tail", "-n1")
	cut := exec.Command("cut", "-d/", "-f3")

	fmt.Println(git.String())
	output, stderr, err := Pipeline(git, tail, cut)
	//output, stderr, err := Pipeline(git)

	outputStr := string(output)
	if strings.HasPrefix(outputStr, "v") {
		outputStr = strings.TrimPrefix(outputStr, "v")
	}
	if strings.HasSuffix(outputStr, "\n") {
		outputStr = strings.TrimSuffix(outputStr, "\n")
	}

	if err != nil {
		fmt.Println(err)
	}

	if len(outputStr) > 0 {
		fmt.Printf("output:\n%s\n", outputStr)
	}
	v1 := semver.New(string(outputStr))

	fmt.Println(v1)

	// Print the stderr, if any
	if len(stderr) > 0 {
		fmt.Printf("err: (stderr)\n%s", stderr)
	}

	//stdout := os.Stdout
	//c1 := exec.Command("ls")
	//run(c1, nil, nil, stdout)

	//run2(c1, nil, stdout)

	//c2 := exec.Command("tail", "-n1")
	//run2(c2, stdout, stdout)
	//c3 := exec.Command("wc -c")
	//run(c1, nil, c2, os.Stdout)
	//run(c3, os.Stdout, nil, nil)

	//pr, pw := io.Pipe()
	//c1.Stdout = pw
	//c2.Stdin = pr
	//c2.Stdout = os.Stdout
	//
	//c1.Start()
	//c2.Start()
	//
	//go func() {
	//	defer pw.Close()
	//
	//	c1.Wait()
	//}()
	//c2.Wait()
}
func run2(c1 *exec.Cmd, stdin *os.File, stdout *os.File) {
	if stdin != nil {
		c1.Stdin = stdin
	}
	if stdout != nil {
		c1.Stdout = stdout
	}
	err := c1.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
	}

	err = c1.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
	}
}

func run(c1 *exec.Cmd, c1Stdin *os.File, c2 *exec.Cmd, c2Stdout *os.File) {
	if c2 != nil {
		pr, pw := io.Pipe()
		if c1Stdin != nil {
			c1.Stdin = c1Stdin
		}
		c1.Stdout = pw
		c2.Stdin = pr

		if c2Stdout != nil {
			c2.Stdout = c2Stdout
		} else {
			c2.Stdout = os.Stdout
		}

		c1.Start()
		c2.Start()
		go func() {
			defer pw.Close()

			c1.Wait()
		}()
		c2.Wait()
	} else {
		if c1Stdin != nil {
			c1.Stdin = c1Stdin
		}
		stdout, err := c1.StdoutPipe()
		scanner := bufio.NewScanner(stdout)
		go func() {
			for scanner.Scan() {
				fmt.Printf("\t > %s\n", scanner.Text())
			}
		}()

		err = c1.Start()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		}

		err = c1.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		}
	}
}

//func SystemCommand(envList []string, silent bool, cmds []Cmd) error {
//	if !silent {
//		if len(envList) != 0 {
//			fmt.Printf("envList: %v\n", envList)
//		}
//		fmt.Printf("command: %v\n", name)
//		fmt.Printf("args: %v\n", args)
//	}
//	command := exec.Command(name, args...)
//	command.Env = os.Environ()
//
//	if len(envList) > 0 {
//		command.Env = append(command.Env, envList...)
//	}
//
//	stdout, err := command.StdoutPipe()
//	if err != nil {
//		return fmt.Errorf(err.Error())
//
//	}
//	stderr, err := command.StderrPipe()
//	if err != nil {
//		return fmt.Errorf(err.Error())
//	}
//	scanner := bufio.NewScanner(stdout)
//	go func() {
//		for scanner.Scan() {
//			if !silent {
//				fmt.Printf("\t > %s\n", scanner.Text())
//			}
//		}
//	}()
//
//	errScanner := bufio.NewScanner(stderr)
//	go func() {
//		for errScanner.Scan() {
//			fmt.Printf("\t > %s\n", errScanner.Text())
//		}
//	}()
//	err = command.Start()
//	if err != nil {
//		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
//		return err
//	}
//
//	err = command.Wait()
//	if err != nil {
//		if !silent {
//			fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
//		}
//		return err
//	}
//
//	return nil
//}

func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	// Require at least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	// Collect the output from the command(s)
	var output bytes.Buffer
	var stderr bytes.Buffer

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		var err error
		// Connect each command's stdin to the previous command's stdout
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return nil, nil, err
		}
		// Connect each command's stderr to a buffer
		cmd.Stderr = &stderr
	}

	// Connect the output and error for the last command
	cmds[last].Stdout, cmds[last].Stderr = &output, &stderr

	// Start each command
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Wait for each command to complete
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Return the pipeline output and the collected standard error
	return output.Bytes(), stderr.Bytes(), nil
}
