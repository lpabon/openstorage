package volumedrivers

import (
	"github.com/libopenstorage/openstorage/volume"
	"github.com/libopenstorage/openstorage/volume/drivers/fake"
)

var (
	// AllDrivers is a slice of all existing known Drivers.
	AllDrivers = []Driver{
		// Fake driver is used to develop and test the API
		{DriverType: fake.Type, Name: fake.Name},
	}

	volumeDriverRegistry = volume.NewVolumeDriverRegistry(
		map[string]func(map[string]string) (volume.VolumeDriver, error){
			fake.Name: fake.Init,
		},
	)
)
