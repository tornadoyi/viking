package argparse

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	parsers []IParser
)

type IParser interface {
	Parse(app *Application) *CmdClause
	Execute(args interface{})
}


func AddParser(p IParser){
	parsers = append(parsers, p)
}


func ParseAndRun(app *Application, args interface{}){

	if len(parsers) == 0 {
		Parse()
		return
	}

	cmds := make([]*CmdClause, len(parsers))
	for i, p := range parsers{
		cmds[i] = p.Parse(app)
	}

	name := kingpin.MustParse(app.Parse(os.Args[1:]))
	index := 0
	for i := 0; i < len(parsers); i++ {
		index = i
		if cmds[i].FullCommand() != name { continue }
		break
	}

	if index >= len(parsers) {
		Fatalf("Invalid command %v", os.Args[1:])
	}

	parser := parsers[index]
	parser.Execute(args)
}






// export all struct
type Application = kingpin.Application

type CmdClause = kingpin.CmdClause

type FlagClause = kingpin.FlagClause

type ArgClause = kingpin.ArgClause

type ParseContext = kingpin.ParseContext

func New(name, help string) *Application { return kingpin.New(name, help)}

func Command(name, help string) *CmdClause { return kingpin.Command(name, help)}

func Flag(name, help string) *FlagClause { return kingpin.Flag(name, help)}

func Arg(name, help string) *ArgClause { return kingpin.Arg(name, help)}

func Parse() string { return kingpin.Parse()}

func Errorf(format string, args ...interface{}) { kingpin.Errorf(format, args...)}

func Fatalf(format string, args ...interface{}) { kingpin.Fatalf(format, args...)}

func FatalIfError(err error, format string, args ...interface{}) { FatalIfError(err, format, args...)}

func FatalUsage(format string, args ...interface{}) { FatalUsage(format, args...)}

func FatalUsageContext(context *ParseContext, format string, args ...interface{}) { FatalUsageContext(context, format, args...)}

func Usage() { kingpin.Usage()}

func UsageTemplate(template string) *Application { return UsageTemplate(template)}

func MustParse(command string, err error) string { return MustParse(command, err)}

func Version(version string) *Application { return kingpin.Version(version)}