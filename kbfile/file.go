package kbfile

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	ProjectPath = "."
)

// RealPath get an absolute path
func RealPath(path string, addSlash ...bool) (realPath string) {
	if len(path) == 0 || (path[0] != '/' && !filepath.IsAbs(path)) {
		path = ProjectPath + "/" + path
	}
	realPath, _ = filepath.Abs(path)
	realPath = strings.Replace(realPath, "\\", "/", -1)
	realPath = pathAddSlash(realPath, addSlash...)
	return
}

func pathAddSlash(path string, addSlash ...bool) string {
	if len(addSlash) > 0 && addSlash[0] && !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return path
}

// RealPathMkdir get an absolute path, create it if it doesn't exist
func RealPathMkdir(path string, addSlash ...bool) string {
	realPath := RealPath(path, addSlash...)
	if DirExist(realPath) {
		return realPath
	}
	_ = os.MkdirAll(realPath, os.ModePerm)
	return realPath
}

// DirExist Is it an existing directory
func DirExist(path string) bool {
	state, _ := PathExist(path)
	return state == 1
}

// FileExist Is it an existing file?
func FileExist(path string) bool {
	state, _ := PathExist(path)
	return state == 2
}

// PathExist PathExist
// 1 exists and is a directory path, 2 exists and is a file path, 0 does not exist
func PathExist(path string) (int, error) {
	path = RealPath(path)
	f, err := os.Stat(path)
	if err == nil {
		isFile := 2
		if f.IsDir() {
			isFile = 1
		}
		return isFile, nil
	}

	return 0, err
}
