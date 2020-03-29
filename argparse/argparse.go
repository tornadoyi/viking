package argparse

import (
	"github.com/alecthomas/kong"
)

type ICLI interface {
	Context()	interface{}
}


func Parse(cli ICLI) {
	ctx := kong.Parse(cli)
	err := ctx.Run(cli.Context())
	ctx.FatalIfErrorf(err)
}

