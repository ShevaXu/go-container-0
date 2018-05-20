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
	fmt.Println("Usage: container0 run cmd [args]")
	os.Exit(0)
}

func run() {
	fmt.Printf("Running %v as parent PID %d\n\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// UTS namespace gives its processes their own view of the system’s hostname and domain name
		Cloneflags: syscall.CLONE_NEWUTS |
			// MNT mount namespace gives the process’s contained within it their own mount table.
			syscall.CLONE_NEWNS |
			// PID namespace gives a process and its children their own view of a subset of the processes in the system.
			// Think of it as a mapping table.
			syscall.CLONE_NEWPID,
		Unshareflags: syscall.CLONE_NEWNS, // to create a private mount namespace
	}

	must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v as child PID %d\n\n", os.Args[2:], os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	must(syscall.Sethostname([]byte("docker-0")))

	// change the root filesystem
	// (NOTE that /home/rootfs must be a valid dir, e.g., downloaded from
	// https://github.com/wsilva/container-from-scratch-demo/blob/master/ubuntu-rootfs.tar.gz
	must(syscall.Chroot("/home/rootfs"))
	must(os.Chdir("/"))
	// mount /proc so that `ps` can work
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	must(cmd.Run())

	must(syscall.Unmount("proc", 0))
}

func main() {
	if len(os.Args) < 3 {
		help()
	}
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		fmt.Println(`Try "main run cmd [args]"`)
	}
}
