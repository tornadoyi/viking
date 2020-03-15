package file

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	vfilepath "viking/goplus/path/filepath"
)


func Save(filename string, data []byte, perm os.FileMode) error {
	// create all dirs
	filename = vfilepath.Render(filename)
	dir := filepath.Dir(filename)
	if !vfilepath.IsDir(dir) {
		err := os.MkdirAll(dir, perm)
		if err != nil { return err }
	}

	// write to file
	return ioutil.WriteFile(filename, data, perm)
}

func Load(filename string) ([]byte, error){
	return ioutil.ReadFile(filename)
}

func SaveYaml(filename string, v interface{}) error {
	bs, err := yaml.Marshal(v)
	if err != nil { return err }
	return Save(filename, bs, os.ModePerm)
}

func LoadYaml(filename string, v interface{}) (error){
	bs, err := Load(filename)
	if err != nil { return err}
	return  yaml.Unmarshal(bs, v)
}

func SaveJson(filename string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil { return err }
	return Save(filename, bs, os.ModePerm)
}

func LoadJson(filename string, v interface{}) (error){
	bs, err := Load(filename)
	if err != nil { return err}
	return json.Unmarshal(bs, v)
}