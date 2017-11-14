package lib

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

// Process stores GTP compatible process
type Process struct {
	Label  string
	logger *log.Logger
	cmd    *exec.Cmd
	name   string
	args   []string
	stdin  io.WriteCloser
	out    *bufio.Scanner
}

const (
	Black = iota
	White = iota
)

const (
	Pass = iota
	Move = iota
)
