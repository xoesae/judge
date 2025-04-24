package jail

import (
	"bytes"
	"fmt"
	"github.com/xoesae/judge/runner/internal/config"
	"github.com/xoesae/judge/runner/internal/filesystem"
	"os"
	"os/exec"
	"syscall"
)

type ProcessResult struct {
	Output []byte
	Error  []byte
}

func RunIsolated(scriptFilePath string) {
	// mount /proc on rootfs
	filesystem.Mount()

	// isolate the hostname and change the root
	syscall.Sethostname([]byte("sandbox"))

	// set the chroot
	if err := syscall.Chroot(config.GetConfig().RootFs); err != nil {
		fmt.Println("Erro no chroot:", err)
		return
	}

	// set the chdir
	err := syscall.Chdir("/")
	if err != nil {
		panic(err)
	}

	// execute the python file

	// TODO: exec the nsjail command
	cmd := exec.Command("/usr/bin/python3", scriptFilePath)

	// capture the stdout and the stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("error on run script:", err)
	}
}

func InitChildProcess(filename string) (ProcessResult, error) {
	var stdout, stderr bytes.Buffer

	// create a child process to run the script i another process
	cmd := exec.Command("/proc/self/exe", "child", "/tmp/scripts/"+filename)
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWNET,
		UidMappings: []syscall.SysProcIDMap{{ContainerID: 0, HostID: os.Getuid(), Size: 1}},
		GidMappings: []syscall.SysProcIDMap{{ContainerID: 0, HostID: os.Getgid(), Size: 1}},
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return ProcessResult{}, err
	}

	return ProcessResult{
		Output: stdout.Bytes(),
		Error:  stderr.Bytes(),
	}, nil
}
