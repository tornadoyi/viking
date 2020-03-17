package filepath

import (
	"os"
	"os/user"
	_filepath "path/filepath"
	"strings"
)


func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil { return os.IsExist(err) }
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil { return false }
	return s.IsDir()
}

func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil { return false }
	return !s.IsDir()
}

func ExpandUser(path string) string{
	if len(path) == 0 { return path }
	if !strings.HasPrefix(path, "~/") { return path}
	usr, _ := user.Current()
	return _filepath.Join(usr.HomeDir, path[2:])
}


func Render(path string) string { return ExpandUser(_filepath.Clean(path)) }



// export

type WalkFunc = _filepath.WalkFunc

var Abs = _filepath.Abs
var Base = _filepath.Base
var Clean = _filepath.Clean
var Dir = _filepath.Dir
var EvalSymlinks = _filepath.EvalSymlinks
var Ext = _filepath.Ext
var FromSlash = _filepath.FromSlash
var Glob = _filepath.Glob
var HasPrefix = _filepath.HasPrefix
var IsAbs = _filepath.IsAbs
var Join = _filepath.Join
var Match = _filepath.Match
var Rel = _filepath.Rel
var Split = _filepath.Split
var SplitList = _filepath.SplitList
var ToSlash = _filepath.ToSlash
var VolumeName = _filepath.VolumeName
var Walk = _filepath.Walk
