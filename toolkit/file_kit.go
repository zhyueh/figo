package toolkit

import (
	"bytes"
	"io/ioutil"
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

func ParentPath(path string) string {
	if index := strings.LastIndex(path, "/"); index != -1 {
		return path[0 : index+1]
	}
	return ""
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

func ReadAll(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

func WriteAll(body []byte, path string) error {
	return ioutil.WriteFile(path, body, 0666)
	//f, err := os.OpenFile(path, os.O_CREATE, 0666)
	//if err != nil {
	//	return err
	//}
	//defer f.Close()
	//_, err = f.Write(body)
	//if err != nil {
	//	return err
	//}

	//return nil

}

func RemoveFileIfExists(path string) error {
	if IsExists(path) {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
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
