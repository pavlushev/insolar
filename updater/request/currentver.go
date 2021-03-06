package request

import (
	"errors"
	"github.com/insolar/insolar/log"
	"regexp"
)

type version struct {
	Latest   string `json:"latest"`
	Major    int    `json:"major"`
	Minor    int    `json:"minor"`
	Revision int    `json:"revision"`
}

func NewVersion(ver string) *version {
	v := version{}
	v.Latest = ver
	re := regexp.MustCompile("[0-9]+")
	arr := re.FindAllString(ver, -1)
	v.Major = extractIntValue(arr, 0)
	v.Minor = extractIntValue(arr, 1)
	v.Revision = extractIntValue(arr, 2)
	return &v
}

func ReqCurrentVer(addresses []string) (string, *version, error) {
	log.Debug("Found update server addresses: ", addresses)

	for _, address := range addresses {

		log.Info("Found update server address: ", address)
		ver, err := ReqCurrentVerFromAddress(GetProtocol(address), address)

		if err == nil && ver != "" {
			currentVer := ExtractVersion(ver)
			return address, currentVer, err
		}
	}
	log.Warn("No Update Servers available")
	return "", nil, nil
}

func ReqCurrentVerFromAddress(request RequestUpdateNode, address string) (string, error) {
	log.Debug("Check latest version info from remote server: ", address)
	if request == nil {
		return "", errors.New("Unknown protocol")
	}
	return request.getCurrentVer(address)
}
