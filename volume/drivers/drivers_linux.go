package volumedrivers

asdf

import (
	"github.com/libopenstorage/openstorage/volume"
	"github.com/libopenstorage/openstorage/volume/drivers/aws"
	"github.com/libopenstorage/openstorage/volume/drivers/btrfs"
	"github.com/libopenstorage/openstorage/volume/drivers/buse"
	"github.com/libopenstorage/openstorage/volume/drivers/coprhd"
	"github.com/libopenstorage/openstorage/volume/drivers/fake"
	"github.com/libopenstorage/openstorage/volume/drivers/nfs"
	"github.com/libopenstorage/openstorage/volume/drivers/pwx"
	"github.com/libopenstorage/openstorage/volume/drivers/vfs"
)

var (
	// AllDrivers is a slice of all existing known Drivers.
	AllDrivers = []Driver{
		// AWS driver provisions storage from EBS.
		{DriverType: aws.Type, Name: aws.Name},
		// BTRFS driver provisions storage from local btrfs.
		{DriverType: btrfs.Type, Name: btrfs.Name},
		// BUSE driver provisions storage from local volumes and implements block in user space.
		{DriverType: buse.Type, Name: buse.Name},
		// COPRHD driver
		{DriverType: coprhd.Type, Name: coprhd.Name},
		// NFS driver provisions storage from an NFS server.
		{DriverType: nfs.Type, Name: nfs.Name},
		// PWX driver provisions storage from PWX cluster.
		{DriverType: pwx.Type, Name: pwx.Name},
		// VFS driver provisions storage from local filesystem
		{DriverType: vfs.Type, Name: vfs.Name},
		// Fake driver is used to develop and test the API
		{DriverType: fake.Type, Name: fake.Name},
	}

	volumeDriverRegistry = volume.NewVolumeDriverRegistry(
		map[string]func(map[string]string) (volume.VolumeDriver, error){
			aws.Name:    aws.Init,
			btrfs.Name:  btrfs.Init,
			buse.Name:   buse.Init,
			coprhd.Name: coprhd.Init,
			nfs.Name:    nfs.Init,
			pwx.Name:    pwx.Init,
			vfs.Name:    vfs.Init,
			fake.Name:   fake.Init,
		},
	)
)
