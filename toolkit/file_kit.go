package toolkit

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func PathJoin(args ...string) string {
	buf := bytes.NewBuffer([]byte{})
	for i, p := range args {
		tmp := strings.TrimRight(p, "/")
		if i != 0 && tmp[0] != '/' {
			buf.WriteString("/")
		}
		buf.WriteString(tmp)
	}

	return buf.String()
}

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func EnsureDirExists(dirPaths ...string) {
	for _, dirPath := range dirPaths {
		if !IsExists(dirPath) {
			os.MkdirAll(dirPath, 0755)
		}
	}
}

func RemoveFileIfExists(path string) bool {
	if IsExists(path) {
		err := os.Remove(path)
		if err != nil {
			return false
		}
	}
	return true
}

type WalkFileFunc func(fileInfo os.FileInfo)

func WalkDirFiles(root string, handler WalkFileFunc) {
	filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		handler(fi)
		return nil
	})
}

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

func IsWindows() bool {
	return (runtime.GOOS == "windows")
}
