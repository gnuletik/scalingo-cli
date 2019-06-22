package config

import (
	"encoding/json"
	"os"

	"gopkg.in/errgo.v1"
)

func SetRegion(regionName string) error {
	authenticator := &CliAuthenticator{}
	_, token, err := authenticator.LoadAuth()
	if err != nil {
		return errgo.Notef(err, "fail to get current user")
	}

	region, err := GetRegion(C, token.Token, regionName)
	if err != nil {
		return errgo.Notef(err, "fail to select region")
	}

	C.ConfigFile.Region = region.Name
	fd, err := os.OpenFile(C.ConfigFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return errgo.Notef(err, "fail to open config file")
	}
	defer fd.Close()

	err = json.NewEncoder(fd).Encode(C.ConfigFile)
	if err != nil {
		return errgo.Notef(err, "fail to persist config file %v", C.ConfigFilePath)
	}

	return nil
}
