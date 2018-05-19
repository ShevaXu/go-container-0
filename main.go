package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func help() {
	// docker run <image> cmd args
	fmt.Println("Usage: main run cmd [args]")
	os.Exit(0)
}

func run() {
	fmt.Printf("Running %v\n\n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// UTS namespace gives its processes their own view of the system’s hostname and domain name
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	must(syscall.Sethostname([]byte("docker-0")))

	must(cmd.Run())
}

func main() {
	if len(os.Args) < 3 {
		help()
	}
	switch os.Args[1] {
	case "run":
		run()
	default:
		fmt.Println(`Try "main run cmd [args]"`)
	}
}
