package toolkit

import (
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

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

func DiskUsage(path string) DiskStatus {
	disk := DiskStatus{}
	if IsWindows() {
		return disk
	}

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err == nil {
		disk.All = fs.Blocks * uint64(fs.Bsize)
		disk.Free = fs.Bfree * uint64(fs.Bsize)
		disk.Used = disk.All - disk.Free
	}
	return disk
}

func IsWindows() bool {
	return (runtime.GOOS == "windows")
}
