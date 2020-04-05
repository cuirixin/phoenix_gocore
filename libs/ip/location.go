package ip

import (
	"errors"
	"github.com/gogather/com"
	"github.com/gogather/iplocation"
)

// IpLocation struct
type IpLocation struct {
	Key string
	Il *iplocation.IpLocation
}

// NewIpLocation init
func NewIpLocation(key string) *IpLocation {
	ip := new(IpLocation)
	ip.Key = key
	ip.Il = iplocation.NewIpLocation(key)
	return ip
}

func (il *IpLocation) GetIpLocation(ip string) (string, error) {
	json, err := il.Location(ip)
	if json == nil {
		return "", errors.New("json is nil")
	}
	countryName := json["countryName"].(string)
	regionName := json["regionName"].(string)
	cityName := json["cityName"].(string)
	data := map[string]interface{}{
		"countryName": countryName,
		"regionName":  regionName,
		"cityName":    cityName,
	}
	str, err := com.JsonEncode(data)
	return str, err
}
