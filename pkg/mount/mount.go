package mount

import (
	"errors"
	"sync"
	"time"

	"github.com/docker/docker/pkg/mount"
	"github.com/libopenstorage/openstorage/pkg/keylock"
)

// Manager defines the interface for keep track of volume driver mounts.
type Manager interface {
	// String representation of the mount table
	String() string
	// Reload mount table for specified device.
	Reload(source string) error
	// Load mount table for all devices that match the list of identifiers
	Load(source []string) error
	// Inspect mount table for specified source. ErrEnoent may be returned.
	Inspect(source string) []*PathInfo
	// Mounts returns paths for specified source.
	Mounts(source string) []string
	// HasMounts determines returns the number of mounts for the source.
	HasMounts(source string) int
	// HasTarget determines returns the number of mounts for the target.
	HasTarget(target string) (string, bool)
	// Exists returns true if the device is mounted at specified path.
	// returned if the device does not exists.
	Exists(source, path string) (bool, error)
	// GetSourcePath scans mount for a specified mountPath and returns the
	// sourcePath if found or returnes an ErrEnoent
	GetSourcePath(mountPath string) (string, error)
	// GetSourcePaths returns all source paths from the mount table
	GetSourcePaths() []string
	// Mount device at mountpoint
	Mount(
		minor int,
		device string,
		path string,
		fs string,
		flags uintptr,
		data string,
		timeout int,
		opts map[string]string) error
	// Unmount device at mountpoint and remove from the matrix.
	// ErrEnoent is returned if the device or mountpoint for the device
	// is not found.
	Unmount(source, path string, flags int, timeout int, opts map[string]string) error
	// RemoveMountPath removes the given path
	RemoveMountPath(path string, opts map[string]string) error
	// EmptyTrashDir removes all directories from the mounter trash directory
	EmptyTrashDir() error
}

// MountImpl backend implementation for Mount/Unmount calls
type MountImpl interface {
	Mount(source, target, fstype string, flags uintptr, data string, timeout int) error
	Unmount(target string, flags int, timeout int) error
}

// MountType indicates different mount types supported
type MountType int

const (
	// DeviceMount indicates a device mount type
	DeviceMount MountType = 1 << iota
	// NFSMount indicates a NFS mount point
	NFSMount
	// CustomMount indicates a custom mount type with its
	// own defined way of handling mount table
	CustomMount
	// BindMount indicates a bind mount point
	BindMount
)

const mountPathRemoveDelay = 30 * time.Second

var (
	// ErrExist is returned if path is already mounted to a different device.
	ErrExist = errors.New("Mountpath already exists")
	// ErrEnoent is returned for a non existent mount point
	ErrEnoent = errors.New("Mountpath is not mounted")
	// ErrEinval is returned is fields for an entry do no match
	// existing fields
	ErrEinval = errors.New("Invalid arguments for mount entry")
	// ErrUnsupported is returned for an unsupported operation or a mount type.
	ErrUnsupported = errors.New("Not supported")
	// ErrMountpathNotAllowed is returned when the requested mountpath is not
	// a part of the provided allowed mount paths
	ErrMountpathNotAllowed = errors.New("Mountpath is not allowed")
)

// DeviceMap map device name to Info
type DeviceMap map[string]*Info

// PathMap map path name to device
type PathMap map[string]string

// PathInfo is a reference counted path
type PathInfo struct {
	Root string
	Path string
}

// Info per device
type Info struct {
	sync.Mutex
	Device     string
	Minor      int
	Mountpoint []*PathInfo
	Fs         string
}

// Mounter implements Ops and keeps track of active mounts for volume drivers.
type Mounter struct {
	sync.Mutex
	mountImpl     MountImpl
	mounts        DeviceMap
	paths         PathMap
	allowedDirs   []string
	kl            keylock.KeyLock
	trashLocation string
}

type findMountPoint func(source *mount.Info, destination string, mountInfo []*mount.Info) (bool, string, string)

// DefaultMounter defaults to syscall implementation.
type DefaultMounter struct {
}
