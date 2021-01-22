package lib

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)


type Filesystem struct {
	Device string
	MountPoint string
	FsType string
	Options string
	Ro bool
}



func ListFileSystems() ([]Filesystem, error) {
	var filesystems []Filesystem
	file, err := os.Open(procFilePath("1/mounts"))
	if errors.Is(err, os.ErrNotExist) {
		// Fallback to `/proc/mounts` if `/proc/1/mounts` is missing due hidepid.
		log.Debugf("Reading root mounts failed, falling back to system mounts (err:%v)", err)
		file, err = os.Open(procFilePath("mounts"))
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())

		if len(parts) < 4 {
			return nil, fmt.Errorf("malformed mount point information: %q", scanner.Text())
		}

		// Ensure we handle the translation of \040 and \011
		// as per fstab(5).
		parts[1] = strings.Replace(parts[1], "\\040", " ", -1)
		parts[1] = strings.Replace(parts[1], "\\011", "\t", -1)

		fs := Filesystem{
			Device:     parts[0],
			MountPoint: rootfsStripPrefix(parts[1]),
			FsType:     parts[2],
			Options:    parts[3],
		}
		fs.Ro = false
		for _, option := range strings.Split(fs.Options, ",") {
			if option == "ro" {
				fs.Ro = true
				break
			}
		}
		filesystems = append(filesystems, fs)
	}

	return filesystems, scanner.Err()
}


func rootfsStripPrefix(path string) string {
	if rootfsPath == "/" {
		return path
	}
	stripped := strings.TrimPrefix(path, rootfsPath)
	if stripped == "" {
		return "/"
	}
	return stripped
}
