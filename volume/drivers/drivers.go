//go:generate mockgen -package=mock -destination=mock/driver.mock.go github.com/libopenstorage/openstorage/volume VolumeDriver

package volumedrivers

import (
	"github.com/libopenstorage/openstorage/api"
	"github.com/libopenstorage/openstorage/volume"
)

// Driver is the description of a supported OST driver. New Drivers are added to
// the drivers array
type Driver struct {
	DriverType api.DriverType
	Name       string
}

// Get returns a VolumeDriver based on input name.
func Get(name string) (volume.VolumeDriver, error) {
	return volumeDriverRegistry.Get(name)
}

// Register registers a new driver.
func Register(name string, params map[string]string) error {
	return volumeDriverRegistry.Register(name, params)
}

// Add adds a new driver.
func Add(name string, init func(map[string]string) (volume.VolumeDriver, error)) error {
	return volumeDriverRegistry.Add(name, init)
}

// Remove removes driver from registry. Does nothing if driver does not exist
func Remove(name string) {
	volumeDriverRegistry.Remove(name)
}

// Shutdown stops the volume driver registry
func Shutdown() error {
	return volumeDriverRegistry.Shutdown()
}
