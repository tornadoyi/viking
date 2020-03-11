package filepath

import (
	"os"
	"os/user"
	"path/filepath"
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
	return !IsDir(path)
}

func ExpandUser(path string) string{
	if len(path) == 0 { return path }
	if !strings.HasPrefix(path, "~/") { return path}
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, path[2:])
}


func Render(path string) string { return ExpandUser(filepath.Clean(path)) }

