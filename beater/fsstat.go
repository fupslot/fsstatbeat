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
	Name  string      `json:"name"`
	Path  string      `json:"path"`
	Umask string      `json:"umask"`
	Owner string      `json:"owner"`
	Group string      `json:"group"`
	Perm  os.FileMode `json:"perm"`
	Octal string      `json:"octal"`
}

func Fsstat(r config.Resource) (state *FileState, err error) {
	if info, err := os.Stat(r.File.Path); err == nil {
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
			Name:  filepath.Base(r.File.Path),
			Path:  r.File.Path,
			Umask: umask,
			Octal: octal,
			Perm:  info.Mode().Perm(),
			Owner: u.Username,
			Group: g.Name,
		}

		return state, nil
	}

	return nil, err
}
