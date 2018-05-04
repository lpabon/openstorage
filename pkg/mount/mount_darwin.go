package mount

import (
	"fmt"
	"strings"

	"github.com/docker/docker/pkg/mount"
	"github.com/libopenstorage/openstorage/pkg/chattr"
)

// Mount default mount implementation is syscall.
func (m *DefaultMounter) Mount(
	source string,
	target string,
	fstype string,
	flags uintptr,
	data string,
	timeout int,
) error {
	return nil
}

// Unmount default unmount implementation is syscall.
func (m *DefaultMounter) Unmount(target string, flags int, timeout int) error {
	return nil
}

// String representation of Mounter
func (m *Mounter) String() string {
	s := struct {
		mounts        DeviceMap
		paths         PathMap
		allowedDirs   []string
		trashLocation string
	}{
		mounts:        m.mounts,
		paths:         m.paths,
		allowedDirs:   m.allowedDirs,
		trashLocation: m.trashLocation,
	}

	return fmt.Sprintf("%#v", s)
}

// Inspect mount table for device
func (m *Mounter) Inspect(sourcePath string) []*PathInfo {
	m.Lock()
	defer m.Unlock()

	v, ok := m.mounts[sourcePath]
	if !ok {
		return nil
	}
	return v.Mountpoint
}

// Mounts returns  mount table for device
func (m *Mounter) Mounts(sourcePath string) []string {
	m.Lock()
	defer m.Unlock()

	v, ok := m.mounts[sourcePath]
	if !ok {
		return nil
	}

	mounts := make([]string, len(v.Mountpoint))
	for i, v := range v.Mountpoint {
		mounts[i] = v.Path
	}

	return mounts
}

// GetSourcePaths returns all source paths from the mount table
func (m *Mounter) GetSourcePaths() []string {
	m.Lock()
	defer m.Unlock()

	sourcePaths := make([]string, len(m.mounts))
	i := 0
	for path := range m.mounts {
		sourcePaths[i] = path
		i++
	}
	return sourcePaths
}

// HasMounts determines returns the number of mounts for the device.
func (m *Mounter) HasMounts(sourcePath string) int {
	m.Lock()
	defer m.Unlock()

	v, ok := m.mounts[sourcePath]
	if !ok {
		return 0
	}
	return len(v.Mountpoint)
}

// HasTarget returns true/false based on the target provided
func (m *Mounter) HasTarget(targetPath string) (string, bool) {
	m.Lock()
	defer m.Unlock()

	for k, v := range m.mounts {
		for _, p := range v.Mountpoint {
			if p.Path == targetPath {
				return k, true
			}
		}
	}
	return "", false
}

// Exists scans mountpaths for specified device and returns true if path is one of the
// mountpaths. ErrEnoent may be retuned if the device is not found
func (m *Mounter) Exists(sourcePath string, path string) (bool, error) {
	m.Lock()
	defer m.Unlock()

	v, ok := m.mounts[sourcePath]
	if !ok {
		return false, ErrEnoent
	}
	for _, p := range v.Mountpoint {
		if p.Path == path {
			return true, nil
		}
	}
	return false, nil
}

// GetSourcePath scans mount for a specified mountPath and returns the sourcePath
// if found or returnes an ErrEnoent
func (m *Mounter) GetSourcePath(mountPath string) (string, error) {
	m.Lock()
	defer m.Unlock()

	for k, v := range m.mounts {
		for _, p := range v.Mountpoint {
			if p.Path == mountPath {
				return k, nil
			}
		}
	}
	return "", ErrEnoent
}

func normalizeMountPath(mountPath string) string {
	if len(mountPath) > 1 && strings.HasSuffix(mountPath, "/") {
		return mountPath[:len(mountPath)-1]
	}
	return mountPath
}

func (m *Mounter) maybeRemoveDevice(device string) {
	m.Lock()
	defer m.Unlock()
	if info, ok := m.mounts[device]; ok {
		// If the device has no more mountpoints, remove it from the map
		if len(info.Mountpoint) == 0 {
			delete(m.mounts, device)
		}
	}
}

// reload from newM
func (m *Mounter) reload(device string, newM *Info) error {
	m.Lock()
	defer m.Unlock()

	// New mountable has no mounts, delete old mounts.
	if newM == nil {
		delete(m.mounts, device)
		return nil
	}

	// Old mountable had no mounts, copy over new mounts.
	oldM, ok := m.mounts[device]
	if !ok {
		m.mounts[device] = newM
		return nil
	}

	// Overwrite old mount entries into new mount table, preserving refcnt.
	for _, oldP := range oldM.Mountpoint {
		for j, newP := range newM.Mountpoint {
			if newP.Path == oldP.Path {
				newM.Mountpoint[j] = oldP
				break
			}
		}
	}

	// Purge old mounts.
	m.mounts[device] = newM
	return nil
}

func (m *Mounter) load(prefixes []string, fmp findMountPoint) error {
	info, err := mount.GetMounts()
	if err != nil {
		return err
	}
DeviceLoop:
	for _, v := range info {
		var (
			sourcePath, devicePath string
			foundPrefix            bool
		)
		for _, devPrefix := range prefixes {
			foundPrefix, sourcePath, devicePath = fmp(v, devPrefix, info)
			if foundPrefix {
				break
			}
		}
		if !foundPrefix {
			continue
		}
		mount, ok := m.mounts[sourcePath]
		if !ok {
			mount = &Info{
				Device:     devicePath,
				Fs:         v.Fstype,
				Minor:      v.Minor,
				Mountpoint: make([]*PathInfo, 0),
			}
			m.mounts[sourcePath] = mount
		}
		// Allow Load to be called multiple times.
		for _, p := range mount.Mountpoint {
			if p.Path == v.Mountpoint {
				continue DeviceLoop
			}
		}
		mount.Mountpoint = append(
			mount.Mountpoint,
			&PathInfo{
				Root: normalizeMountPath(v.Root),
				Path: normalizeMountPath(v.Mountpoint),
			},
		)
		m.paths[v.Mountpoint] = sourcePath
	}
	return nil
}

// Mount new mountpoint for specified device.
func (m *Mounter) Mount(
	minor int,
	devPath, path, fs string,
	flags uintptr,
	data string,
	timeout int,
	opts map[string]string,
) error {
	return nil
}

// Unmount device at mountpoint and from the matrix.
// ErrEnoent is returned if the device or mountpoint for the device is not found.
func (m *Mounter) Unmount(
	devPath string,
	path string,
	flags int,
	timeout int,
	opts map[string]string,
) error {
	return nil
}

// RemoveMountPath makes the path writeable and removes it after a fixed delay
func (m *Mounter) RemoveMountPath(mountPath string, opts map[string]string) error {
	return nil
}

func (m *Mounter) EmptyTrashDir() error {
	return nil
}

// isPathSetImmutable returns true on error in getting path info or if path
// is immutable .
func (m *Mounter) isPathSetImmutable(mountpath string) bool {
	return chattr.IsImmutable(mountpath)
}

// makeMountpathReadOnly makes given mountpath read-only
func (m *Mounter) makeMountpathReadOnly(mountpath string) error {
	return chattr.AddImmutable(mountpath)
}

// makeMountpathWriteable makes given mountpath writeable
func (m *Mounter) makeMountpathWriteable(mountpath string) error {
	return chattr.RemoveImmutable(mountpath)
}

// New returns a new Mount Manager
func New(
	mounterType MountType,
	mountImpl MountImpl,
	identifiers []string,
	customMounter CustomMounter,
	allowedDirs []string,
	trashLocation string,
) (Manager, error) {
	return nil, ErrUnsupported
}
