package nfs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	losetup "gopkg.in/freddierice/go-losetup.v1"

	"math/rand"
	"strings"

	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/config"
	"github.com/libopenstorage/openstorage/pkg/mount"
	"github.com/libopenstorage/openstorage/pkg/seed"
	"github.com/libopenstorage/openstorage/pkg/util"
	"github.com/libopenstorage/openstorage/volume"
	"github.com/libopenstorage/openstorage/volume/drivers/common"
	"github.com/pborman/uuid"
	"github.com/portworx/kvdb"
)

const (
	Name         = "nfs"
	Type         = api.DriverType_DRIVER_TYPE_FILE
	NfsDBKey     = "OpenStorageNFSKey"
	nfsMountPath = "/var/lib/openstorage/nfs/"
	nfsBlockFile = ".blockdevice"
)

// Implements the open storage volume interface.
type driver struct {
	volume.IODriver
	volume.StoreEnumerator
	volume.StatsDriver
	volume.QuiesceDriver
	volume.CredsDriver
	volume.CloudBackupDriver
	volume.CloudMigrateDriver
	nfsServers []string
	nfsPath    string
	mounter    mount.Manager
	driverType api.DriverType
}

func Init(params map[string]string) (volume.VolumeDriver, error) {
	path, ok := params["path"]
	if !ok {
		return nil, errors.New("No NFS path provided")
	}
	server, ok := params["server"]
	if !ok {
		logrus.Printf("No NFS server provided, will attempt to bind mount %s", path)
	} else {
		logrus.Printf("NFS driver initializing with %s:%s ", server, path)
	}
	//support more than one server using CSV
	//TB-FIXME: modify driver params flow to support map[string]struct/array
	servers := strings.Split(server, ",")

	// Create a mount manager for this NFS server. Blank sever is OK.
	mounter, err := mount.New(mount.NFSMount, nil, servers, nil, []string{}, "")
	if err != nil {
		logrus.Warnf("Failed to create mount manager for server: %v (%v)", server, err)
		return nil, err
	}
	inst := &driver{
		IODriver:           volume.IONotSupported,
		StoreEnumerator:    common.NewDefaultStoreEnumerator(Name, kvdb.Instance()),
		StatsDriver:        volume.StatsNotSupported,
		QuiesceDriver:      volume.QuiesceNotSupported,
		nfsServers:         servers,
		CredsDriver:        volume.CredsNotSupported,
		nfsPath:            path,
		mounter:            mounter,
		CloudBackupDriver:  volume.CloudBackupNotSupported,
		CloudMigrateDriver: volume.CloudMigrateNotSupported,
	}
	blockEnabled, ok := params["block"]
	if ok && blockEnabled == "true" {
		logrus.Info("NFS driver now in block mode")
		inst.driverType = api.DriverType_DRIVER_TYPE_BLOCK
	} else {
		inst.driverType = api.DriverType_DRIVER_TYPE_FILE
	}

	//make directory for each nfs server
	for _, v := range servers {
		logrus.Infof("Calling mkdirAll: %s", nfsMountPath+v)
		if err := os.MkdirAll(nfsMountPath+v, 0744); err != nil {
			return nil, err
		}
	}
	src := inst.nfsPath
	if server != "" {
		src = ":" + inst.nfsPath
	}

	//mount each nfs server
	for _, v := range inst.nfsServers {
		// If src is already mounted at dest, leave it be.
		o, err := exec.Command("/bin/mount", "-t", "nfs").CombinedOutput()
		if !strings.Contains(string(o), src+" on "+nfsMountPath+v) {
			if server != "" {
				err = syscall.Mount(
					src,
					nfsMountPath+v,
					"nfs",
					0,
					"nolock,addr="+v,
				)
			} else {
				err = syscall.Mount(src, nfsMountPath+v, "", syscall.MS_BIND, "")
			}
			if err != nil {
				logrus.Errorf("Unable to mount %s:%s at %s (%+v)",
					v, inst.nfsPath, nfsMountPath+v, err)
				return nil, err
			} else {
				logrus.Infof("NFS: %s mounted", nfsMountPath+v)
			}
		} else {
			logrus.Infof("NFS: %s already mounted", nfsMountPath+v)
		}
	}

	volumeInfo, err := inst.StoreEnumerator.Enumerate(&api.VolumeLocator{}, nil)
	if err == nil {
		for _, info := range volumeInfo {
			if info.Status == api.VolumeStatus_VOLUME_STATUS_NONE {
				info.Status = api.VolumeStatus_VOLUME_STATUS_UP
				inst.UpdateVol(info)
			}
		}
	}

	logrus.Println("NFS initialized and driver mounted at: ", nfsMountPath)
	return inst, nil
}

func (d *driver) Name() string {
	return Name
}

func (d *driver) Type() api.DriverType {
	return d.driverType
}

func (d *driver) Version() (*api.StorageVersion, error) {
	return &api.StorageVersion{
		Driver:  d.Name(),
		Version: "1.0.0",
	}, nil
}

// Status diagnostic information
func (d *driver) Status() [][2]string {
	return [][2]string{}
}

//
//Utility functions
//
func (d *driver) getNewVolumeServer() (string, error) {
	//randomly select one
	if d.nfsServers != nil && len(d.nfsServers) > 0 {
		return d.nfsServers[rand.Intn(len(d.nfsServers))], nil
	}

	return "", errors.New("No NFS servers found")
}

//get nfsPath for specified volume
func (d *driver) getNFSPath(v *api.Volume) (string, error) {
	locator := v.GetLocator()
	server, ok := locator.VolumeLabels["server"]
	if !ok {
		logrus.Warnf("No server label found on volume")
		return "", fmt.Errorf("No server label found on volume: " + v.Id)
	}

	return path.Join(nfsMountPath, server), nil
}

//get nfsPath for specified volume
func (d *driver) getNFSPathById(volumeID string) (string, error) {
	v, err := d.GetVol(volumeID)
	if err != nil {
		return "", err
	}

	return d.getNFSPath(v)
}

//get nfsPath plus volume name for specified volume
func (d *driver) getNFSVolumePath(v *api.Volume) (string, error) {
	parentPath, err := d.getNFSPath(v)
	if err != nil {
		return "", err
	}

	return path.Join(parentPath, v.Id), nil
}

//get nfsPath plus volume name for specified volume
func (d *driver) getNFSVolumePathById(volumeID string) (string, error) {
	v, err := d.GetVol(volumeID)
	if err != nil {
		return "", err
	}

	return d.getNFSVolumePath(v)
}

//append unix time to volumeID
func (d *driver) getNewSnapVolName(volumeID string) string {
	return volumeID + "-" + strconv.FormatUint(uint64(time.Now().Unix()), 10)
}

//
// These functions below implement the volume driver interface.
//

func (d *driver) Create(
	locator *api.VolumeLocator,
	source *api.Source,
	spec *api.VolumeSpec) (string, error) {

	if len(locator.Name) == 0 {
		return "", fmt.Errorf("volume name cannot be empty")
	}

	if hasSpaces := strings.Contains(locator.Name, " "); hasSpaces {
		return "", fmt.Errorf("volume name cannot contain space characters")
	}

	volumeID := strings.TrimSuffix(uuid.New(), "\n")

	if _, err := d.GetVol(volumeID); err == nil {
		return "", fmt.Errorf("volume with that id already exists")
	}

	//snapshot passes nil volumelabels
	if locator.VolumeLabels == nil {
		locator.VolumeLabels = make(map[string]string)
	}

	//check if user passed server as option
	labels := locator.GetVolumeLabels()
	_, ok := labels["server"]
	if !ok {
		server, err := d.getNewVolumeServer()
		if err != nil {
			logrus.Infof("no nfs servers found...")
			return "", err
		} else {
			logrus.Infof("Assigning random nfs server: %s to volume: %s", server, volumeID)
		}

		labels["server"] = server
	}

	// Create a directory on the NFS server with this UUID.
	volPathParent := path.Join(nfsMountPath, labels["server"])
	volPath := path.Join(volPathParent, volumeID)
	err := os.MkdirAll(volPath, 0744)
	if err != nil {
		logrus.Println(err)
		return "", err
	}

	// Setup volume object
	if source != nil {
		if len(source.Seed) != 0 {
			seed, err := seed.New(source.Seed, locator.VolumeLabels)
			if err != nil {
				logrus.Warnf("Failed to initailize seed from %q : %v",
					source.Seed, err)
				return "", err
			}
			err = seed.Load(path.Join(volPath, config.DataDir))
			if err != nil {
				logrus.Warnf("Failed to  seed from %q to %q: %v",
					source.Seed, volPathParent, err)
				return "", err
			}
		}
	}

	// Create volume
	var v *api.Volume
	if d.driverType == api.DriverType_DRIVER_TYPE_BLOCK {
		v = common.NewVolume(
			volumeID,
			spec.GetFormat(),
			locator,
			source,
			spec,
		)
		if err := d.CreateVol(v); err != nil {
			return "", err
		}

		if source != nil && len(source.GetParent()) != 0 {
			// Need to clone
			if err := d.clone(volumeID, source.GetParent()); err != nil {
				return "", err
			}
		} else {
			// This is not a snapshot, but a new volume
			// so let's format the file
			blockFile := path.Join(volPathParent, volumeID+nfsBlockFile)
			f, err := os.Create(blockFile)
			if err != nil {
				logrus.Println(err)
				return "", err
			}
			defer f.Close()

			// Create sparse file
			if err := f.Truncate(int64(spec.Size)); err != nil {
				logrus.Println(err)
				return "", err
			}

			// Format
			if spec.GetFormat() != api.FSType_FS_TYPE_NONE {
				dev, err := losetup.Attach(blockFile, 0, false)
				if err != nil {
					return "", err
				}
				logrus.Infof("Formatting %s with %v", dev, spec.Format)
				cmd := "/sbin/mkfs." + spec.Format.SimpleString()
				o, err := exec.Command(cmd, blockFile).Output()
				if err != nil {
					logrus.Warnf("Failed to run command %v %v: %v", cmd, dev, o)
					return "", err
				}
				dev.Detach()
			}
		}
	} else {
		// File based
		v = common.NewVolume(
			volumeID,
			api.FSType_FS_TYPE_NFS,
			locator,
			source,
			spec,
		)
		if err := d.CreateVol(v); err != nil {
			return "", err
		}
	}

	return v.Id, err
}

func (d *driver) Delete(volumeID string) error {
	v, err := d.GetVol(volumeID)
	if err != nil {
		logrus.Println(err)
		return err
	}

	nfsVolPath, err := d.getNFSVolumePath(v)
	if err != nil {
		return err
	}

	// Delete the simulated block volume
	os.Remove(nfsVolPath + nfsBlockFile)

	// Delete the directory on the nfs server.
	os.RemoveAll(nfsVolPath)

	err = d.DeleteVol(volumeID)
	if err != nil {
		logrus.Println(err)
		return err
	}

	return nil
}

func (d *driver) MountedAt(mountpath string) string {
	return ""
}

func (d *driver) Mount(volumeID string, mountpath string, options map[string]string) error {
	v, err := d.GetVol(volumeID)
	if err != nil {
		logrus.Println(err)
		return err
	}

	nfsPath, err := d.getNFSPath(v)
	if err != nil {
		logrus.Printf("Could not find server for volume: %s", volumeID)
		return err
	}

	if v.GetSpec().GetSize() == 0 {
		// File access
		srcPath := path.Join(":", nfsPath, volumeID)
		mountExists, err := d.mounter.Exists(srcPath, mountpath)
		if err != nil {
			return err
		}
		if !mountExists {
			d.mounter.Unmount(path.Join(nfsPath, volumeID), mountpath,
				syscall.MNT_DETACH, 0, nil)
			if err := d.mounter.Mount(
				0, path.Join(nfsPath, volumeID),
				mountpath,
				string(v.Spec.Format),
				syscall.MS_BIND,
				"",
				0,
				nil,
			); err != nil {
				logrus.Printf("Cannot mount %s at %s because %+v",
					path.Join(nfsPath, volumeID), mountpath, err)
				return err
			}
		}
	} else {
		// Block access
		if v.GetState() != api.VolumeState_VOLUME_STATE_ATTACHED {
			return fmt.Errorf("Voume %s is not attached", volumeID)
		}
		mountExists, _ := d.mounter.Exists(v.DevicePath, mountpath)
		if !mountExists {
			if err := syscall.Mount(v.DevicePath, mountpath, v.Spec.Format.SimpleString(), 0, ""); err != nil {
				return fmt.Errorf("Failed to mount %v at %v: %v", v.DevicePath, mountpath, err)
			}
		}
	}
	if v.AttachPath == nil {
		v.AttachPath = make([]string, 0)
	}
	v.AttachPath = append(v.AttachPath, mountpath)
	return d.UpdateVol(v)

}

func (d *driver) Unmount(volumeID string, mountpath string, options map[string]string) error {
	v, err := d.GetVol(volumeID)
	if err != nil {
		return err
	}
	if len(v.AttachPath) == 0 {
		return fmt.Errorf("Device %v not mounted", volumeID)
	}

	nfsVolPath, err := d.getNFSVolumePath(v)
	if err != nil {
		return err
	}

	if v.GetSpec().GetSize() == 0 {
		err = d.mounter.Unmount(nfsVolPath, mountpath, syscall.MNT_DETACH, 0, nil)
		if err != nil {
			return err
		}
		v.AttachPath = d.mounter.Mounts(nfsVolPath)
	} else {
		if err := syscall.Unmount(mountpath, 0); err != nil {
			return err
		}
		v.AttachPath = nil
	}
	return d.UpdateVol(v)
}

func (d *driver) clone(newVolumeID, volumeID string) error {
	nfsVolPath, err := d.getNFSVolumePathById(volumeID)
	if err != nil {
		return err
	}

	newNfsVolPath, err := d.getNFSVolumePathById(newVolumeID)
	if err != nil {
		return err
	}

	// NFS does not support snapshots, so just copy the files.
	if err := copyDir(nfsVolPath, newNfsVolPath); err != nil {
		d.Delete(newVolumeID)
		return nil
	}

	// First try reflinks
	_, err = exec.Command(
		"/bin/cp",
		"--reflink=always",
		nfsVolPath+nfsBlockFile,
		newNfsVolPath+nfsBlockFile,
	).Output()
	if err == nil {
		logrus.Infof("Cloned %s to %s using reflink copy",
			nfsVolPath+nfsBlockFile,
			newNfsVolPath+nfsBlockFile)
	} else {
		// Second try sparse copy
		_, err = exec.Command(
			"/bin/cp",
			"--sparse=always",
			nfsVolPath+nfsBlockFile,
			newNfsVolPath+nfsBlockFile,
		).Output()
		if err == nil {
			logrus.Infof("Cloned %s to %s using sparse copy",
				nfsVolPath+nfsBlockFile,
				newNfsVolPath+nfsBlockFile)
		} else {
			// if block, copy the block file also
			if err := copyFile(nfsVolPath+nfsBlockFile, newNfsVolPath+nfsBlockFile); err == nil {
				logrus.Infof("Cloned %s to %s using slow copy",
					nfsVolPath+nfsBlockFile,
					newNfsVolPath+nfsBlockFile)
			} else {
				logrus.Errorf("Failed to clone %s to %s: %v",
					nfsVolPath+nfsBlockFile,
					newNfsVolPath+nfsBlockFile,
					err)
				d.Delete(newVolumeID)
				return nil
			}
		}
	}
	return nil
}

func (d *driver) Snapshot(volumeID string, readonly bool, locator *api.VolumeLocator, noRetry bool) (string, error) {
	volIDs := []string{volumeID}
	vols, err := d.Inspect(volIDs)
	if err != nil {
		return "", nil
	}
	source := &api.Source{Parent: volumeID}
	logrus.Infof("Creating snap vol name: %s", locator.Name)
	return d.Create(locator, source, vols[0].Spec)
}

func (d *driver) Restore(volumeID string, snapID string) error {
	if _, err := d.Inspect([]string{volumeID, snapID}); err != nil {
		return err
	}

	nfsVolPath, err := d.getNFSVolumePathById(volumeID)
	if err != nil {
		return err
	}

	snapNfsVolPath, err := d.getNFSVolumePathById(snapID)
	if err != nil {
		return err
	}

	// NFS does not support restore, so just copy the files.
	if err := copyDir(snapNfsVolPath, nfsVolPath); err != nil {
		return err
	}
	return nil
}

func (d *driver) SnapshotGroup(groupID string, labels map[string]string, volumeIDs []string) (*api.GroupSnapCreateResponse, error) {

	return nil, volume.ErrNotSupported
}

func (d *driver) Attach(volumeID string, attachOptions map[string]string) (string, error) {

	nfsPath, err := d.getNFSPathById(volumeID)
	if err != nil {
		return "", err
	}
	blockFile := path.Join(nfsPath, volumeID+nfsBlockFile)

	// Check if it is block
	v, err := util.VolumeFromName(d, volumeID)
	if err != nil {
		return "", err
	}

	// If it has no size, no need to attach
	if v.GetSpec().GetSize() == 0 {
		return blockFile, nil
	}

	// If it is a block device, create a loop device
	dev, err := losetup.Attach(blockFile, 0, false /* not read only: TODO change this */)
	if err != nil {
		return "", err
	}

	// Update volume info
	v.DevicePath = dev.Path()
	v.State = api.VolumeState_VOLUME_STATE_ATTACHED
	if err := d.UpdateVol(v); err != nil {
		dev.Detach()
		return "", err
	}

	return dev.Path(), nil
}

func (d *driver) Detach(volumeID string, options map[string]string) error {

	// Get volume info
	v, err := util.VolumeFromName(d, volumeID)
	if err != nil {
		return err
	}

	// If it has no size, no need to detach
	if v.GetSpec().GetSize() == 0 {
		return nil
	} else if v.GetState() != api.VolumeState_VOLUME_STATE_ATTACHED {
		// if it is not attached, just return
		return nil
	}

	// Detach -- code from https://github.com/freddierice/go-losetup
	loopFile, err := os.OpenFile(v.GetDevicePath(), os.O_RDONLY, 0660)
	if err != nil {
		return fmt.Errorf("could not open loop device")
	}
	defer loopFile.Close()

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, loopFile.Fd(), losetup.ClrFd, 0)
	if errno != 0 {
		return fmt.Errorf("error clearing loopfile: %v", errno)
	}

	// Update volume info
	v.DevicePath = ""
	v.State = api.VolumeState_VOLUME_STATE_NONE
	if err := d.UpdateVol(v); err != nil {
		return err
	}

	return nil
}

func (d *driver) Set(volumeID string, locator *api.VolumeLocator, spec *api.VolumeSpec) error {
	if spec != nil {
		return volume.ErrNotSupported
	}
	v, err := d.GetVol(volumeID)
	if err != nil {
		return err
	}
	if locator != nil {
		v.Locator = locator
	}
	return d.UpdateVol(v)
}

func (d *driver) Shutdown() {
	logrus.Printf("%s Shutting down", Name)

	for _, v := range d.nfsServers {
		logrus.Infof("Umounting: %s", nfsMountPath+v)
		syscall.Unmount(path.Join(nfsMountPath, v), 0)
	}
}

func copyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func copyDir(source string, dest string) (err error) {
	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = copyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = copyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func (d *driver) Catalog(volumeID, path, depth string) (api.CatalogResponse, error) {
	return api.CatalogResponse{}, volume.ErrNotSupported
}
