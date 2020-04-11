package strcase

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/tornadoyi/viking/goplus/core"
	"strings"
)


/*
s := "AnyKind of_string"
ToSnake(s)										any_kind_of_string
ToSnakeWithIgnore(s, '.')						any_kind.of_string
ToScreamingSnake(s)								ANY_KIND_OF_STRING
ToKebab(s)										any-kind-of-string
ToScreamingKebab(s)								ANY-KIND-OF-STRING
ToDelimited(s, '.')								any.kind.of.string
ToScreamingDelimited(s, '.', '', true)			ANY.KIND.OF.STRING
ToScreamingDelimited(s, '.', ' ', true)			ANY.KIND OF.STRING
ToCamel(s)										AnyKindOfString
ToLowerCamel(s)									anyKindOfString

*/
type Case uint
func (k Case) String() string {
	if int(k) < len(caseNames) {
		return caseNames[k]
	}
	return "Unknown case"
}

const (
	BeginTitle				 					Case = iota
	Title
	Snake
	SnakeWithIgnore
	ScreamingSnake
	Kebab
	ScreamingKebab
	Delimited
	ScreamingDelimited
	Camel
	LowerCamel
)

var(
	ToBeginTitle								= strings.Title
	ToTitle										= strings.ToTitle
	ToSnake 									= strcase.ToSnake
	ToSnakeWithIgnore 							= strcase.ToSnakeWithIgnore
	ToScreamingSnake 							= strcase.ToScreamingSnake
	ToKebab 									= strcase.ToKebab
	ToScreamingKebab 							= strcase.ToScreamingKebab
	ToDelimited 								= strcase.ToDelimited
	ToScreamingDelimited 						= strcase.ToScreamingDelimited
	ToCamel 									= strcase.ToCamel
	ToLowerCamel 								= strcase.ToLowerCamel


	caseNames = []string{
		BeginTitle:								"BeginTitle",
		Title:									"Title",
		Snake:									"Snake",
		SnakeWithIgnore:						"SnakeWithIgnore",
		ScreamingSnake:							"ScreamingSnake",
		Kebab:									"Kebab",
		ScreamingKebab:							"ScreamingKebab",
		Delimited:								"Delimited",
		ScreamingDelimited:						"ScreamingDelimited",
		Camel:									"Camel",
		LowerCamel:								"LowerCamel",
	}

	lowerName2Cases								= func() map[string]Case {
		m := make(map[string]Case)
		for c, n := range caseNames { m[strings.ToLower(n)] = Case(c) }
		return m
	}()
)


func ToCase(s string, c Case, opts... interface{}) (string, error) {
	switch c {
	case BeginTitle: 					return ToBeginTitle(s), nil
	case Title: 						return ToTitle(s), nil
	case Snake: 						return ToSnake(s), nil
	case SnakeWithIgnore:
		if len(opts) != 1 { return "", fmt.Errorf("invalid arguments, expect 1")}
		p1, err := core.UInt8(opts[0])
		if err != nil { return "", err}
		return ToSnakeWithIgnore(s, p1), nil
	case ScreamingSnake: 				return ToScreamingSnake(s), nil
	case Kebab: 						return ToKebab(s), nil
	case ScreamingKebab: 				return ToScreamingKebab(s), nil
	case Delimited:
		if len(opts) != 1 { return "", fmt.Errorf("invalid arguments, expect 1")}
		p1, err := core.UInt8(opts[0])
		if err != nil { return "", err}
		return ToDelimited(s, p1), nil
	case ScreamingDelimited:
		if len(opts) != 3 { return "", fmt.Errorf("invalid arguments, expect 3")}
		p1, err := core.UInt8(opts[0])
		if err != nil { return "", err}
		p2, err := core.UInt8(opts[0])
		if err != nil { return "", err}
		p3, err := core.Bool(opts[0])
		if err != nil { return "", err}
		return ToScreamingDelimited(s, p1, p2, p3), nil
	case Camel: 						return ToCamel(s), nil
	case LowerCamel: 					return ToLowerCamel(s), nil
	}
	return "", fmt.Errorf("Invalid case %v", c)
}


func ToNameCase(s string, c string, opt... interface{}) (string, error) {
	c = strings.ToLower(c)
	t, ok := lowerName2Cases[c];
	if !ok { return "", fmt.Errorf("Invalid case name %v", c)}
	return ToCase(s, t, opt...)
}

