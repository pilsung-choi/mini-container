# mini-container

Go로 직접 만들어보는 미니 컨테이너 런타임입니다.
`docker run`이 내부에서 리눅스 커널 기능(네임스페이스, cgroups, 파일시스템 격리)을 어떻게 조합하는지,
한 단계씩 직접 구현하며 이해하는 것이 목표입니다.

## 실행 환경

컨테이너의 핵심 기능은 리눅스 커널 기능이라 macOS에서는 동작하지 않습니다.
[Lima](https://lima-vm.io/)로 띄운 Ubuntu VM 안에서 개발·실행합니다.

```bash
limactl shell container-lab
sudo go run . run /bin/sh
```

## 단계별 진행

| 단계 | 컨테이너 기능 | 배운 Go | 상태 |
|---|---|---|---|
| 1 | `run` 명령으로 프로세스 실행 | 모듈, `os.Args`, `os/exec`, 에러 처리 | ✅ 완료 |
| 2 | UTS·PID 네임스페이스 | `syscall`, 구조체 필드 설정 | ✅ 완료 |
| 3 | 루트 파일시스템 격리 | 파일·디렉터리, `defer` | |
| 4 | cgroups로 메모리·CPU 제한 | 파일 쓰기, 구조체 + 메서드 | |
| 5 | 네트워크 격리 | 인터페이스, 고루틴 | |
| 6 | 이미지 레이어 (overlayfs) | 패키지 분리, 테스트 | |

## 단계별 기록

각 단계에서 만든 것, 배운 것, 막혔던 것을 기록합니다.

### 1단계 — 프로세스 실행

- **만든 것**: `run <명령>`으로 받은 명령을 자식 프로세스로 실행하고, 표준 입출력을 터미널에 연결
- **배운 것**
  - `os.Args`는 명령줄 인자 목록이고 `[0]`은 프로그램 자신
  - `exec.Command`는 명령서만 만들고, 실제 실행은 `cmd.Run()`에서 일어남
  - 자식 프로세스의 Stdin·Stdout·Stderr는 기본적으로 막혀 있어서 직접 연결해야 함
  - Go는 예외 대신 에러를 반환값으로 돌려주고 `if err != nil`로 처리
- **관찰**: 격리가 없어서 자식 bash의 hostname과 PID가 VM과 똑같이 보임

### 2단계 — UTS·PID 네임스페이스

- **만든 것**: `SysProcAttr.Cloneflags`에 `CLONE_NEWUTS | CLONE_NEWPID`를 넣어, 자식 프로세스가 새 네임스페이스에서 태어나게 함
- **확인한 것**

  | 실험 | 결과 |
  |---|---|
  | 컨테이너 안에서 `hostname mini` | 안쪽만 `mini`, VM은 `lima-container-lab` 그대로 (UTS 격리) |
  | 컨테이너 안에서 `echo $$` | `1` (PID 격리). VM에서 보면 같은 bash가 32000번대 PID |
  | `sudo` 없이 실행 | `operation not permitted` (네임스페이스 생성은 root 권한 필요) |

- **배운 것**
  - 네임스페이스 플래그는 `syscall`의 상수이고, `go doc -all syscall | grep CLONE_NEW`로 전체 목록을 찾을 수 있음
  - 플래그는 비트 하나씩 켜진 숫자라 `|`로 합쳐 여러 네임스페이스를 동시에 만듦
  - `SysProcAttr` 칸은 포인터 타입이라 `&syscall.SysProcAttr{...}`로 넘김
- **남은 문제**: 안에서 `echo $$`는 1인데 `ps aux`에는 VM 프로세스 146개가 전부 보임.
  `ps`는 `/proc` 파일시스템을 읽는데, 파일시스템은 아직 VM과 공유하고 있기 때문 → 3단계에서 마운트 네임스페이스와 새 `/proc`으로 해결

### 3단계 — 루트 파일시스템 격리

_(진행 후 작성)_
