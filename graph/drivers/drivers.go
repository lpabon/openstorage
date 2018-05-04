package graphdrivers

import (
	"github.com/libopenstorage/openstorage/api"
)

// Driver is the description of a supported OST driver. New Drivers are added to
// the drivers array
type Driver struct {
	DriverType api.DriverType
	Name       string
}
