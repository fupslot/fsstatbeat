package beater

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/fupslot/fsstatbeat/config"
)

type FileState struct {
	name  string
	path  string
	umask string
	owner string
	group string
	perm  os.FileMode
	octal string
}

func Fsstat(r config.Resource) (state *FileState, err error) {
	if info, err := os.Stat(r.Path); err == nil {
		stat := info.Sys().(*syscall.Stat_t)
		uid := stat.Uid
		gid := stat.Gid

		// Retriving file's User and Group
		u, _ := user.LookupId(strconv.FormatUint(uint64(uid), 10))
		g, _ := user.LookupGroupId(strconv.FormatUint(uint64(gid), 10))

		umask := info.Mode().Perm().String()
		octal := strconv.FormatUint(uint64(info.Mode().Perm()), 8)

		// fmt.Printf("%s %s %s %s %s\n", filePath, umask, octal, u.Username, g.Name)

		state := &FileState{
			name:  filepath.Base(r.Path),
			path:  r.Path,
			umask: umask,
			octal: octal,
			perm:  info.Mode().Perm(),
			owner: u.Username,
			group: g.Name,
		}

		// fmt.Printf("%+v", fileResource)
		return state, nil
	}

	return nil, err
}
