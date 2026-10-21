package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("사용법: mini-container run <명령> [인자...]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		run(os.Args[2], os.Args[3:])
	default:
		fmt.Println("알 수 없는 명령:", os.Args[1])
		os.Exit(1)
	}
}

func run(command string, args []string) {
	fmt.Printf("실행: %s %v (PID %d)\n", command, args, os.Getpid())

	// TODO 1: exec.Command로 command와 args를 실행할 cmd를 만든다
	cmd := exec.Command(command, args...)
	// var cmd *exec.Cmd = exec.Command(command, args...)
	// TODO 2: cmd의 Stdin·Stdout·Stderr를 os.Stdin·os.Stdout·os.Stderr와 연결한다
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 2단계: 새 UTS(호스트명)·PID(프로세스 번호) 네임스페이스에서 자식 프로세스를 실행
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	// TODO 4: cmd.Run()으로 실행하고, 에러가 나면 출력한 뒤 os.Exit(1)
	if err := cmd.Run(); err != nil {
		fmt.Println("오류:", err)
		os.Exit(1)
	}
}
