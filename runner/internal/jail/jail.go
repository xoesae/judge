package jail

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xoesae/judge/runner/internal/config"
	"github.com/xoesae/judge/runner/internal/filesystem"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"syscall"
)

type ExecutionStatus string

const (
	NoError      ExecutionStatus = "NO_ERROR"
	MemoryLimit  ExecutionStatus = "MEMORY_LIMIT"
	RuntimeError ExecutionStatus = "RUNTIME_ERROR"
	Timeout      ExecutionStatus = "TIMEOUT"
)

type ProcessResult struct {
	Output   []byte
	Error    ExecutionStatus
	ExitCode int
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
	nsjailPath := "/bin/nsjail"
	cmd := exec.Command(nsjailPath,
		"--mode", "o",
		"--chroot", "/",
		"--proc_rw",
		"--quiet",
		"--disable_clone_newnet",
		"--disable_clone_newuser",
		"--disable_clone_newns",
		"--disable_clone_newpid",
		"--disable_clone_newipc",
		"--disable_clone_newuts",
		"--disable_clone_newcgroup",
		// TODO: add this value to a variable param
		"--time_limit", "10", // (seconds)
		// TODO: add this value to a variable param
		"--rlimit_as", "15", // (MB)
		"--", "/usr/bin/python3", scriptFilePath,
	)

	// capture the stdout and the stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	code := "0"

	if err := cmd.Run(); err != nil {
		var exitError *exec.ExitError

		if errors.As(err, &exitError) {
			c := exitError.ExitCode()
			code = strconv.Itoa(c)

			// writing the exit code in the stderr
			cmd.Stderr.Write([]byte("[code]" + code + "[code]"))
			return
		}

		fmt.Println("error:", err)
		return
	}

	cmd.Stderr.Write([]byte("[code]" + code + "[code]"))
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

	result := ProcessResult{
		Output:   stdout.Bytes(),
		Error:    NoError,
		ExitCode: 0,
	}

	rgx := regexp.MustCompile(`\[code](.*?)\[code]`)
	matches := rgx.FindStringSubmatch(stderr.String())

	// debug status code
	fmt.Println(stderr.String())

	exitCode, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(err)
	}

	// generic
	if exitCode != 0 {
		result.Error = RuntimeError
	}

	// memory limit
	if exitCode == 139 {
		result.Error = MemoryLimit
	}

	// timeout
	if exitCode == 137 {
		result.Error = Timeout
	}

	result.ExitCode = exitCode

	return result, nil
}
