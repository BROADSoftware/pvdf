package lib

import (
	"fmt"
	"golang.org/x/sys/unix"
	"github.com/BROADSoftware/pvdf/shared/common"
	v1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
	"sync"
	"time"
)


type Stats struct {
	Size uint64
	Free uint64
	Avail uint64
	Files uint64
	FilesFree uint64
	Error error
}


type Volume struct {
	FileSystem	Filesystem
	PodId      	string
	Type 		string
	Name 		string
	Stats		Stats
}

/*
 Filter filesystems to extract the ones which are pod's volume.
 We assume the pattern of the mount point is:
 /var/lib/kubelet/pods/<podId>/volumes/<type>/<pvname>[/mount]
*/
func GetVolumeByName(fsList []Filesystem) map[string]Volume {
	volumeByName := make(map[string]Volume)
	for _, fs := range fsList {
		if strings.HasPrefix(fs.MountPoint, "/var/lib/kubelet/pods/") {
			parts := strings.Split(fs.MountPoint, "/")
			if len(parts) < 9 {
				log.Warnf("Invalid volume name:%s (Less than 8 parts)", fs.MountPoint)
			} else {
				volumeByName[parts[8]] = Volume {
					FileSystem: fs,
					PodId:      parts[5],
					Type: 		parts[7],
					Name: 		parts[8],
				}
			}
		}
	}
	return volumeByName
}

var stuckVolumeSet = make(map[string]struct{})
var stuckVolumeSetMtx = &sync.Mutex{}

func (this *Volume) markStuck() {
	stuckVolumeSetMtx.Lock()
	stuckVolumeSet[this.Name] = struct{}{}
	stuckVolumeSetMtx.Unlock()
}

func (this *Volume) unmarkStuck() {
	stuckVolumeSetMtx.Lock()
	if _, ok := stuckVolumeSet[this.Name]; ok {
		log.Infof("PV '%s' has recovered", this.Name)
	}
	delete(stuckVolumeSet, this.Name)
	stuckVolumeSetMtx.Unlock()
}

func (this *Volume) isStuck() bool {
	stuckVolumeSetMtx.Lock()
	_, ok := stuckVolumeSet[this.Name]
	stuckVolumeSetMtx.Unlock()
	return ok
}


func (this *Volume) GetStats() {

	if this.isStuck() {
		this.Stats.Error = fmt.Errorf("PV '%s' is marked as stuck", this.Name)
		return
	}

	success := make(chan struct{})

	go func(r chan struct{}) {
		buf := new(unix.Statfs_t)
		//log.Tracef("Before Statfs on %s", this.FileSystem.MountPoint)
		err := unix.Statfs(rootfsFilePath(this.FileSystem.MountPoint), buf)
		//log.Tracef("After Statfs on %s", this.FileSystem.MountPoint)
		if err != nil {
			this.Stats.Error = fmt.Errorf("Error on statfs() system call on %s: %v", rootfsFilePath(this.FileSystem.MountPoint), err)
		} else {
			this.Stats.Size = buf.Blocks * uint64(buf.Bsize)
			this.Stats.Free = buf.Bfree * uint64(buf.Bsize)
			this.Stats.Avail = buf.Bavail * uint64(buf.Bsize)
			this.Stats.Files = uint64(buf.Files)
			this.Stats.FilesFree = uint64(buf.Ffree)
			this.Stats.Error = nil
		}
		this.unmarkStuck()
		close(success)
	}(success)

	select {
	case <-success:
		// Success. Nothing more to do
	case <-time.After(StatfsTimeout):
		select {
		case <-success:
			// Success came in just after the timeout was reached. Nothing to do anymore
		default:
			//log.Warnf("Mount point '%s' timeout!", this.FileSystem.MountPoint)
			this.markStuck()
			this.Stats.Error = fmt.Errorf("Access timeout!")
		}
	}
}


func b2mib(x uint64) int {
	return int(x/(1024*1025))
}
func (this *Volume) AdjustAnnotationsOn(pv v1.PersistentVolume) (dirty bool) {
	free := b2mib(this.Stats.Free)
	oldFreeStr, ok := pv.Annotations[common.FreeAnnotation]
	if ok {
		oldFree, _ := strconv.Atoi(oldFreeStr)
		if oldFree != free {
			dirty = true
			pv.Annotations[common.FreeAnnotation] = strconv.Itoa(free)
		}
	} else {
		dirty = true
		pv.Annotations[common.FreeAnnotation] = strconv.Itoa(free)
	}
	size := b2mib(this.Stats.Size)
	oldSizeStr, ok := pv.Annotations[common.SizeAnnotation]
	if ok {
		oldSize, _ := strconv.Atoi(oldSizeStr)
		if oldSize != size {
			dirty = true
			pv.Annotations[common.SizeAnnotation] = strconv.Itoa(size)
		}
	} else {
		dirty = true
		pv.Annotations[common.SizeAnnotation] = strconv.Itoa(size)
	}
	return dirty
}
