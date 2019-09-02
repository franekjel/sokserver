package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/franekjel/sokserver/logger"
)

//Fs handles operations on files at given path. All Fs method paths all relative to this path
type Fs struct {
	Path string
}

//Init initialize Fs at given path. Path must exist and be accesible
func Init(path string, dir string) *Fs {
	fs := Fs{filepath.Join(path, dir)}
	if fi, err := os.Stat(fs.Path); err != nil || !fi.IsDir() {
		log.Fatal("%s is inaccesible", fs)
	}
	return &fs
}

//Join is wrapper on filepath.join
func Join(path string, dir string) string {
	return filepath.Join(path, dir)
}

//FileExist check if file exist
func (fs *Fs) FileExist(file string) bool {
	_, err := os.Stat(filepath.Join(fs.Path, file))
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

//ReadFile copy content of file to table. Doesn't check errors (throw fatal)
func (fs *Fs) ReadFile(file string) []byte {
	buff, err := ioutil.ReadFile(filepath.Join(fs.Path, file))
	if err != nil {
		log.Fatal("Cannot read file %s, %s", filepath.Join(fs.Path, file), err.Error())
	}
	return buff
}

//WriteFile weites buff to file and return true on success
func (fs *Fs) WriteFile(name string, buff string) bool {
	err := ioutil.WriteFile(filepath.Join(fs.Path, name), []byte(buff), 0644)
	if err != nil {
		log.Warn("Cannot write to file %s, %s", filepath.Join(fs.Path, name), err.Error())
		return false
	}
	return true
}

//ListFiles return list of regular files (skip directories)
func (fs *Fs) ListFiles(dir string) []string {
	var s []string
	files, err := ioutil.ReadDir(filepath.Join(fs.Path, dir))
	if err != nil {
		log.Fatal("Cannot list files in %s, %s", filepath.Join(fs.Path, dir), err.Error())
	}
	for _, file := range files {
		if !file.IsDir() {
			s = append(s, file.Name())
		}
	}
	return s
}

//ListDirs return list of dirs
func (fs *Fs) ListDirs(dir string) []string {
	var s []string
	files, err := ioutil.ReadDir(filepath.Join(fs.Path, dir))
	if err != nil {
		log.Fatal("Cannot list dirs in %s, %s", filepath.Join(fs.Path, dir), err.Error())
	}
	for _, file := range files {
		if file.IsDir() {
			s = append(s, file.Name())
		}
	}
	return s
}

//CreateDirectory return true on success and false othervise
func (fs *Fs) CreateDirectory(name string) bool {
	err := os.Mkdir(filepath.Join(fs.Path, name), 0755)
	if err != nil {
		log.Warn("Cannot create directory %s, %s", filepath.Join(fs.Path, name), err.Error())
		return false
	}
	return true
}

//RemoveFile removes given file
func (fs *Fs) RemoveFile(name string) bool {
	err := os.Remove(name)
	if err != nil {
		return false
	}
	return true
}
