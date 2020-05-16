package pty

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/creack/pty"
	"github.com/thearchitect/thearchitect.github.io/internal/environment"
	"golang.org/x/crypto/ssh/terminal"
)

type PTY struct {
	Master *os.File
	Slave  *os.File
}

func New() (*PTY, error) {
	master, slave, err := pty.Open()
	if err != nil {
		return nil, err
	}

	if _, err := terminal.MakeRaw(int(master.Fd())); err != nil {
		return nil, err
	}

	return &PTY{
		Master: master,
		Slave:  slave,
	}, nil
}

func (p *PTY) Close() {
	if err := p.Master.Close(); err != nil {
		panic(err)
	}
	if err := p.Slave.Close(); err != nil {
		panic(err)
	}
}

func (p *PTY) Resize(cols, rows int) error {
	if err := pty.Setsize(p.Master, &pty.Winsize{
		Cols: uint16(cols),
		Rows: uint16(rows),
		//X:    1920,
		//Y:    1080,
	}); err != nil {
		return err
	}

	return nil
}

func (p *PTY) Command(ctx context.Context, name string, args ...string) *exec.Cmd {
	wd, err := os.Executable()
	if err != nil {
		panic(err)
	}

	wd = filepath.Dir(wd)

	cmd := exec.CommandContext(ctx, name, args...)

	cmd.Dir = wd

	cmd.Stdin, cmd.Stdout, cmd.Stderr = p.Slave, p.Slave, p.Slave

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid:  true,
		Setctty: true,
		Ctty:    int(p.Slave.Fd()),
	}

	cmd.Env = environment.Environ().
		//Set("DOCKER_BUILDKIT", "1").
		Set("TERM", "xterm-256color").
		Slice()

	return cmd
}
