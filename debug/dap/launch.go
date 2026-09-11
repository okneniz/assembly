package dap

import (
	"context"
	"errors"
	"io"

	"github.com/okneniz/assembly/debug"
	"github.com/okneniz/assembly/debug/load"
	"github.com/okneniz/assembly/debug/qemu"
	"github.com/okneniz/assembly/debug/session"
)

// Boot is the standard Launcher: load the input per the arch table
// (the load package), boot the executor frozen with its console
// drained into the writer, bind the session to the halted machine.
func Boot(args launchArgs, console io.Writer) (*Runtime, error) {
	setup, err := load.SetupFor(args.Arch)
	if err != nil {
		return nil, err
	}

	res, err := load.Build(setup, load.Input{
		SrcPath: args.Source,
		BinPath: args.Bin,
		SymPath: args.Sym,
		Base:    args.Base,
	})
	if err != nil {
		return nil, err
	}

	return NewImageLauncher(setup.Tgt, res.Img, res.Syms, res.Lines)(args, console)
}

// NewImageLauncher boots ready-made parts: the executor image, the
// symbols, and the line map (a prog program assembles them in-process
// and serves its own debug session through this). The machine boots
// frozen at the dap launch request, its console drained into the
// writer the server hands over.
func NewImageLauncher(
	tgt debug.Target,
	img []byte,
	syms map[string]uint64,
	lines []session.Line,
) Launcher {
	return func(args launchArgs, console io.Writer) (*Runtime, error) {
		opts := qemu.NewOptions()
		opts.Serial = console
		opts.Icount = args.Icount
		machine, err := qemu.Start(context.Background(), tgt, img, opts)
		if err != nil {
			return nil, err
		}

		s, err := session.New(machine.Conn(), tgt, syms, lines)
		if err != nil {
			return nil, errors.Join(err, machine.Close())
		}

		return &Runtime{Tgt: tgt, Sess: s, Lines: lines, Closer: machine}, nil
	}
}
