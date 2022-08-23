//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package allocdir

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/shoenig/netlog"
	"golang.org/x/sys/unix"
)

var (
	LOG = netlog.New("alloc_dir")
)

var (
	// SharedAllocContainerPath is the path inside container for mounted
	// directory shared across tasks in a task group.
	SharedAllocContainerPath = filepath.Join("/", SharedAllocName)

	// TaskLocalContainer is the path inside a container for mounted directory
	// for local storage.
	TaskLocalContainerPath = filepath.Join("/", TaskLocal)

	// TaskSecretsContainerPath is the path inside a container for mounted
	// secrets directory
	TaskSecretsContainerPath = filepath.Join("/", TaskSecrets)
)

// dropDirPermissions gives full access to a directory to all users and sets
// the owner to nobody.
func dropDirPermissions(path string, desired os.FileMode) error {
	LOG.Info("dropDirPermissions", "path", path, "desired", string(desired))

	if err := os.Chmod(path, desired|0777); err != nil {
		LOG.Error("Chmod failed", "error", err)
		return fmt.Errorf("Chmod(%v) failed: %v", path, err)
	}
	LOG.Info("Chmod ok")

	// Can't change owner if not root.
	if unix.Geteuid() != 0 {
		LOG.Error("Geteuid failed")
		return nil
	}
	LOG.Info("Geteuid ok")

	u, err := user.Lookup("nobody")
	if err != nil {
		LOG.Error("Lookup nobody failed", "error", err)
		return err
	}
	LOG.Info("Lookup ok")

	uid, err := getUid(u)
	if err != nil {
		LOG.Error("getUid failed", "error", err)
		return err
	}
	LOG.Info("getUid ok")

	gid, err := getGid(u)
	if err != nil {
		LOG.Error("getGid failed", "error", err)
		return err
	}
	LOG.Info("getGid ok")

	if err := os.Chown(path, uid, gid); err != nil {
		LOG.Error("Chown failed", "error", err)
		return fmt.Errorf("Couldn't change owner/group of %v to (uid: %v, gid: %v): %v", path, uid, gid, err)
	}
	LOG.Info("Chown ok")

	return nil
}

// getUid for a user
func getUid(u *user.User) (int, error) {
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return 0, fmt.Errorf("Unable to convert Uid to an int: %v", err)
	}

	return uid, nil
}

// getGid for a user
func getGid(u *user.User) (int, error) {
	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return 0, fmt.Errorf("Unable to convert Gid to an int: %v", err)
	}

	return gid, nil
}

// linkOrCopy attempts to hardlink dst to src and fallsback to copying if the
// hardlink fails.
func linkOrCopy(src, dst string, uid, gid int, perm os.FileMode) error {
	// Avoid link/copy if the file already exists in the chroot
	// TODO 0.6 clean this up. This was needed because chroot creation fails
	// when a process restarts.
	if fileInfo, _ := os.Stat(dst); fileInfo != nil {
		return nil
	}
	// Attempt to hardlink.
	if err := os.Link(src, dst); err == nil {
		return nil
	}

	return fileCopy(src, dst, uid, gid, perm)
}

func getOwner(fi os.FileInfo) (int, int) {
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return -1, -1
	}
	return int(stat.Uid), int(stat.Gid)
}
