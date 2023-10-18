package localcommand

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/creack/pty"
	"github.com/labbs/webtty/utils"
	"github.com/pkg/errors"
)

const (
	DefaultCloseSignal  = syscall.SIGINT
	DefaultCloseTimeout = 10 * time.Second
)

type LocalCommand struct {
	command string
	argv    []string

	closeSignal  syscall.Signal
	closeTimeout time.Duration
	logFile      *os.File
	cmdBuffer    string

	cmd       *exec.Cmd
	pty       *os.File
	ptyClosed chan struct{}
}

func New(command string, argv []string, options ...Option) (*LocalCommand, error) {
	cmd := exec.Command(command, argv...)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to start command `%s`", command)
	}
	ptyClosed := make(chan struct{})
	logFile, err := os.OpenFile("terminal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open log file")
	}

	lcmd := &LocalCommand{
		command: command,
		argv:    argv,

		closeSignal:  DefaultCloseSignal,
		closeTimeout: DefaultCloseTimeout,

		cmd:       cmd,
		pty:       ptmx,
		ptyClosed: ptyClosed,
		logFile:   logFile,
	}

	for _, option := range options {
		option(lcmd)
	}

	// When the process is closed by the user,
	// close pty so that Read() on the pty breaks with an EOF.
	go func() {
		defer func() {
			lcmd.pty.Close()
			lcmd.cmd.Cancel()
			close(lcmd.ptyClosed)
		}()

		lcmd.cmd.Wait()
	}()

	return lcmd, nil
}

func (lcmd *LocalCommand) Read(p []byte) (n int, err error) {
	n, err = lcmd.pty.Read(p)
	if err != nil {
		return n, err
	}

	// Convertir le byte slice en string pour le traitement
	output := string(p[:n])

	// Tronquer ou masquer les informations sensibles
	output = CatchAndTruncate(output)

	// Convertir la string de sortie tronquée en byte slice et copier dans p
	copy(p, []byte(output))

	return n, err
}

func (lcmd *LocalCommand) Write(p []byte) (n int, err error) {
	n, err = lcmd.pty.Write(p)
	if err != nil {
		return n, err
	}

	output := GetTerminalState()

	diff, _ := strings.CutPrefix(output, lcmd.cmdBuffer)

	truncate := CatchAndTruncate(diff)

	// _, errLog := lcmd.logFile.WriteString(truncate)
	// if errLog != nil {
	// 	return n, errLog
	// }

	if utils.RecordingEnabled {
		go utils.PushRecording(truncate)
	}

	lcmd.cmdBuffer = output

	return n, err
}

func (lcmd *LocalCommand) Close() error {
	if lcmd.cmd != nil && lcmd.cmd.Process != nil {
		lcmd.cmd.Process.Signal(lcmd.closeSignal)
	}
	for {
		select {
		case <-lcmd.ptyClosed:
			return nil
		case <-lcmd.closeTimeoutC():
			lcmd.cmd.Process.Signal(syscall.SIGKILL)
		}
	}
}

func (lcmd *LocalCommand) WindowTitleVariables() map[string]interface{} {
	return map[string]interface{}{
		"command": lcmd.command,
		"argv":    lcmd.argv,
		"pid":     lcmd.cmd.Process.Pid,
	}
}

func (lcmd *LocalCommand) ResizeTerminal(width int, height int) error {
	window := struct {
		row uint16
		col uint16
		x   uint16
		y   uint16
	}{
		uint16(height),
		uint16(width),
		0,
		0,
	}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		lcmd.pty.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&window)),
	)
	if errno != 0 {
		return errno
	} else {
		return nil
	}
}

func (lcmd *LocalCommand) closeTimeoutC() <-chan time.Time {
	if lcmd.closeTimeout >= 0 {
		return time.After(lcmd.closeTimeout)
	}

	return make(chan time.Time)
}

func GetTerminalState() string {
	// Exemple : exécutez une commande tmux pour obtenir l'historique du terminal.
	// Vous devrez ajuster cela en fonction de votre mise en œuvre spécifique.
	cmd := exec.Command("tmux", "capture-pane", "-p")
	output, err := cmd.Output()
	if err != nil {
		log.Println("Error capturing tmux pane:", err)
		return ""
	}

	cleanOutput := strings.TrimRight(string(output), "\n")
	lines := strings.Split(cleanOutput, "\n")

	// On ne veut pas afficher la dernière ligne, car elle est vide.
	return strings.Join(lines[:len(lines)-1], "\n")
}

func CatchAndTruncate(s string) string {
	var blackList []string = []string{
		"password",
	}

	lines := strings.Split(s, "\n")

	for _, blackListed := range blackList {
		for i, line := range lines {
			if strings.Contains(line, blackListed) {
				lines[i] = string("[TRUNCATED]")
			}
			if isBase64(line) {
				lines[i] = string("[TRUNCATED]")
			}
		}
	}

	s = strings.Join(lines, "\n")

	return s
}

func isBase64(s string) bool {
	// Regex for verifying base64 strings
	base64RegExp := "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
	base64Re := regexp.MustCompile(base64RegExp)

	// Verify if the string is base64
	return base64Re.MatchString(s)
}
