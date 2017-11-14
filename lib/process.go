package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func ScanLFLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\n', '\n'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, data, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func CreateProcess(label string, command string) (*Process, error) {
	commandParts := strings.Split(command, " ")
	name := commandParts[0]
	args := commandParts[1:]
	cmd := exec.Command(name, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("Error when opening stdin: %+v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Error when opening stdout: %+v", err)
	}
	out := bufio.NewScanner(stdout)
	out.Split(ScanLFLF)
	process := &Process{
		Label:  label,
		logger: log.New(os.Stderr, label+": ", log.LstdFlags),
		name:   name,
		args:   args,
		cmd:    cmd,
		stdin:  stdin,
		out:    out,
	}
	return process, nil
}

func (p *Process) StartProcess() error {
	log.Printf(
		"Starting process %s with args %v and label %s",
		p.name,
		p.args,
		p.Label,
	)
	return p.cmd.Start()
}

func (p *Process) send(message string, args ...interface{}) (string, error) {
	str := fmt.Sprintf("%s\n", fmt.Sprintf(message, args...))
	p.logger.Printf("-> %s", str)

	_, err := io.WriteString(p.stdin, str)
	if err != nil {
		return "", err
	}
	if !p.out.Scan() {
		return "", p.out.Err()
	}
	response := p.out.Text()
	if err != nil {
		return "", err
	}

	if response[0] != '=' {
		return "", fmt.Errorf("Expected '=...', got '%s'", response)
	}

	return response[2 : len(response)-2], nil
}

// Commands
func (p *Process) Name() (string, error) {
	return p.send("name")
}

func (p *Process) Version() (string, error) {
	return p.send("version")
}

func (p *Process) Boardsize(n int) error {
	_, err := p.send("boardsize %d", n)
	return err
}

func (p *Process) Komi(komi string) error {
	_, err := p.send("komi %s", komi)
	return err
}

func (p *Process) ClearBoard() error {
	p.logger.Printf("Clearing board")
	_, err := p.send("clear_board")
	return err
}

func (p *Process) ShowBoard() (string, error) {
	return p.send("showboard")
}

func (p *Process) GenMove(color int) (string, error) {
	var err error
	var str string
	if color == Black {
		str, err = p.send("genmove B")
	} else if color == White {
		str, err = p.send("genmove W")
	}
	if err != nil {
		return "", err
	}
	return str, nil
}

func (p *Process) Play(color int, move string) error {
	var err error
	if color == Black {
		_, err = p.send("play B %s", move)
	} else if color == White {
		_, err = p.send("play W %s", move)
	}
	if err != nil {
		return err
	}
	return nil
}

func (p *Process) FinalScore() (string, error) {
	return p.send("final_score")
}

func (p *Process) Close() error {
	if _, err := p.send("quit"); err != nil {
		return err
	}
	if err := p.stdin.Close(); err != nil {
		return err
	}
	if err := p.cmd.Wait(); err != nil {
		return err
	}
	if !p.cmd.ProcessState.Success() {
		return fmt.Errorf("Process did not quit successfully")
	}
	return nil
}
