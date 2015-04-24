package toolkit

import (
	"syscall"
)

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
