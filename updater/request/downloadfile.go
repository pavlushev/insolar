package request

import (
	"errors"
	"github.com/insolar/insolar/log"
	"path"
)

func DownloadFiles(version string, binariesList []string, url string) (success bool) {
	success = false
	errors := 0
	total := 0

	pathToSave := createCurrentPath(version)
	request := GetProtocol(url)
	log.Info("Download updates from remote server: ", url)
	for _, binary := range binariesList {
		log.Info("Download file : ", binary)
		err := downloadFromAddress(request, path.Join(pathToSave, binary), url+"/"+version+"/"+binary)
		total++
		if err != nil {
			log.Error(err)
			errors++
		} else {
			log.Info("SUCCESS")
		}
	}
	log.Info("Download complete, TOTAL:", total, ", ERRORS: ", errors)
	if errors == 0 && total == len(binariesList) {
		success = true
	}
	return
}

func downloadFromAddress(request RequestUpdateNode, filePath string, url string) error {
	if request == nil {
		return errors.New("Unknown protocol")
	}
	return request.downloadFile(filePath, url)
}
