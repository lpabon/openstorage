package graphdrivers

import (
	"github.com/libopenstorage/openstorage/graph/drivers/chainfs"
	"github.com/libopenstorage/openstorage/graph/drivers/layer0"
	"github.com/libopenstorage/openstorage/graph/drivers/proxy"
)

var (
	// AllDrivers is a slice of all existing known Drivers.
	AllDrivers = []Driver{
		// ChainFS driver implements a chained filesystem using FUSE.
		{DriverType: chainfs.Type, Name: chainfs.Name},
		// Layer0 driver provides persistent storage for the writable layer.
		{DriverType: layer0.Type, Name: layer0.Name},
		// Proxy driver simply uses the Docker overlay driver.
		{DriverType: proxy.Type, Name: proxy.Name},
	}
)
