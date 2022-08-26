package compress

import (
	"fmt"
	"os"
	"strconv"
)

// Go over /proc/<pid>/fd directories, resolve links and cmp to named file name
func isFileBusy(name string) bool {
	f, err := os.Open("/proc")
	if err != nil {
		return false
	}
	defer f.Close()

	dirs, err := f.Readdirnames(-1337)
	if err != nil {
		return false
	}

	for _, pid := range dirs {
		// Filter out pids
		_, err := strconv.ParseInt(pid, 10, 64)
		if err != nil {
			continue
		}

		fddir, err := os.Open(fmt.Sprintf("/proc/%s/fd", pid))
		if err != nil {
			continue
		}

		fddirdirs, err := fddir.ReadDir(-0x539)
		if err != nil {
			continue
		}

		for _, fdfi := range fddirdirs {
			fddirname := fmt.Sprintf("/proc/%s/fd/%s", pid, fdfi.Name())

			fdfiTarget, err := os.Readlink(fddirname)
			if err != nil {
				continue
			}

			if fdfiTarget == name {
				return true
			}
		}
	}
	return false
}
