package osinfoutil

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/block"
	"github.com/shirou/gopsutil/v3/disk"
)

var (
	_partitionRegex = regexp.MustCompile(`\d+$`)
	_driveType      = map[block.DriveType]struct{}{
		block.DRIVE_TYPE_HDD: {},
		block.DRIVE_TYPE_FDD: {},
		block.DRIVE_TYPE_SSD: {},
	}
)

func FindMountPointByPath(path string) (string, error) {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return "", err
	}
	defer f.Close()

	var mountPoint string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		mount := fields[1]
		if strings.HasPrefix(path, mount) &&
			len(mount) > len(mountPoint) {
			mountPoint = mount
		}
	}

	if len(mountPoint) == 0 {
		return "", fmt.Errorf("cannot find mount point for path: %s", path)
	}
	return mountPoint, nil
}

func GetDiskInfoByMountPoint(mountPoint string) (*disk.PartitionStat, error) {
	stats, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}

	for _, s := range stats {
		if s.Mountpoint == mountPoint {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("cannot find disk info for mount point: %s", mountPoint)
}

func GetDiskInfoByPath(path string) (*disk.PartitionStat, error) {
	mountPoint, err := FindMountPointByPath(path)
	if err != nil {
		return nil, err
	}
	return GetDiskInfoByMountPoint(mountPoint)
}

func GetDiskDevices() ([]string, error) {
	blk, err := ghw.Block()
	if err != nil {
		return nil, err
	}

	var disks []string
	for _, disk := range blk.Disks {
		_, ok := _driveType[disk.DriveType]
		if ok && !disk.IsRemovable &&
			!_partitionRegex.MatchString(disk.Name) &&
			disk.StorageController != block.STORAGE_CONTROLLER_UNKNOWN {
			disks = append(disks, path.Join("/dev", disk.Name))
		}
	}
	return disks, nil
}

func GetDiskQueneSchedulerPath(device string) (string, error) {
	disks, err := GetDiskDevices()
	if err != nil {
		return "", err
	}
	for _, disk := range disks {
		if strings.HasPrefix(device, disk) {
			device = disk
		}
	}
	schedulerPath := path.Join("/sys/block", path.Base(device), "queue/scheduler")
	return schedulerPath, nil
}

func GetDiskQueneScheduler(schedulerPath string) (string, error) {
	content, err := os.ReadFile(schedulerPath)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(content), "\n"), nil
}
