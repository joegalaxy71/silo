package main

import (
	"fmt"
	"github.com/mitchellh/go-ps"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main_old() {
	switch os.Args[1] {
	case "run":
		run2()
	case "child":
		child()
	default:
		panic("invalid command!")
	}
}

func run2() {
	//cmd := exec.Command(os.Args[2], os.Args[3:]...)
	//cmd.Stdin = os.Stdin
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//
	//
	//err := cmd.Run()
	//if err!=nil {
	//	fmt.Println("error executing initial cmd")
	//	fmt.Println(err)
	//}
	//
	//os.Exit(0)

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	err := cmd.Run()
	if err != nil {
		fmt.Println("error executing first run")
		fmt.Println(err)
	}
}

func child() {
	fmt.Printf("running %v as PID %d\n", os.Args[2:], os.Getpid())

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	err = syscall.Chroot("/root/go/bin/cross/ubuntu_fs")
	if err != nil {
		fmt.Println("error executing chroot")
		fmt.Println(err)
	}

	err = os.Chdir("/")
	if err != nil {
		fmt.Println("error executing chdir")
		fmt.Println(err)
	}
	dir, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	err = syscall.Mount("proc", "proc", "proc", 0, "")
	if err != nil {
		fmt.Println("error mounting proc")
		fmt.Println(err)
	}

	processes, err := ps.Processes()
	if err != nil {
		fmt.Println("error getting processes")
		fmt.Println(err)
		os.Exit(-1)
	}

	for i, process := range processes {
		fmt.Printf("Entry: %v, PID:%v, Exec:%s\n", i, process.Pid(), process.Executable())
	}

	//var files []string
	//root := "/"
	//err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
	//	files = append(files, path)
	//	return nil
	//})
	//if err != nil {
	//	panic(err)
	//}
	//for _, file := range files {
	//	fmt.Println(file)
	//}

	for i, arg := range os.Args {
		fmt.Printf("arg#%v=%s\n", i, arg)
	}

	if stat, err := os.Stat(os.Args[2]); os.IsNotExist(err) {
		fmt.Printf("File %s doesn't exist\n", os.Args[2])
	} else {
		fmt.Printf("File %s exists\n", os.Args[2])
		fmt.Printf("stat:%+v\n", stat)
	}

	//dat, err := ioutil.ReadFile("/bin/bash")
	//if err!= nil {
	//	fmt.Println("error printing file")
	//	fmt.Println(err)
	//	os.Exit(-1)
	//} else {
	//	fmt.Print(string(dat))
	//}

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	//cmd := exec.Command("/bin/bash", os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("error executing command")
		fmt.Println(err)
	}

	//os.StartProcess()

	err = syscall.Unmount("proc", 0)
	if err != nil {
		fmt.Println("error unmounting proc")
		fmt.Println(err)
	}
}
