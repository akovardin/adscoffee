package geoip

import (
	"github.com/ip2location/ip2location-go/v9"
)

type Geoip struct {
	DBv4 *ip2location.DB
	DBv6 *ip2location.DB
}

func New() (*Geoip, error) {
	filev4, err := NewEmbeddedFile(location, "IP2LOCATION-LITE-DB3.BIN")
	if err != nil {
		panic(err)
	}

	dbv4, err := ip2location.OpenDBWithReader(filev4)
	if err != nil {
		return nil, err
	}

	filev6, err := NewEmbeddedFile(location, "IP2LOCATION-LITE-DB3.IPV6.BIN")
	if err != nil {
		panic(err)
	}

	dbv6, err := ip2location.OpenDBWithReader(filev6)
	if err != nil {
		return nil, err
	}

	return &Geoip{
		DBv4: dbv4,
		DBv6: dbv6,
	}, nil
}
