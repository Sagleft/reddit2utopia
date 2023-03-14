package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"

	utopiago "github.com/Sagleft/utopialib-go/v2"
)

/*
       _              _
      | |            (_)
 _   _| |_ ___  _ __  _  __ _
| | | | __/ _ \| '_ \| |/ _` |
| |_| | || (_) | |_) | | (_| |
 \__,_|\__\___/| .__/|_|\__,_|
               | |
               |_|
*/

type utopiaService struct {
	Token        string
	Host         string
	Port         int
	HTTPSEnabled bool

	Client utopiago.Client
}

func newUtopiaService() *utopiaService {
	return &utopiaService{}
}

func (u *utopiaService) setToken(token string) *utopiaService {
	u.Token = token
	return u
}

func (u *utopiaService) setHost(host string) *utopiaService {
	u.Host = host
	return u
}

func (u *utopiaService) setPort(port int) *utopiaService {
	u.Port = port
	return u
}

func (u *utopiaService) setHTTPS(enabled bool) *utopiaService {
	u.HTTPSEnabled = enabled
	return u
}

func (u *utopiaService) connect() error {
	protocol := "http"
	if u.HTTPSEnabled {
		protocol += "s"
	}

	u.Client = utopiago.NewUtopiaClient(utopiago.Config{
		Protocol: protocol,
		Host:     u.Host,
		Token:    u.Token,
		Port:     u.Port,
	})

	if !u.Client.CheckClientConnection() {
		return errors.New("failed to connect to Utopia client")
	}
	return nil
}

func (u *utopiaService) postMedia(channelID string, media mediaPost) error {
	var imageBytes []byte
	var err error
	if media.IsLocalImage {
		imageBytes, err = ioutil.ReadFile(media.ImageURL)
	} else {
		imageBytes, err = getRemoteFileBytes(media.ImageURL)
	}
	if err != nil {
		return err
	}

	base64Image := base64.StdEncoding.EncodeToString(imageBytes)
	_, err = u.Client.SendChannelPicture(channelID, base64Image, media.Text, "photo.jpg")
	return err
}

func (u *utopiaService) updateAccountName() error {
	data, err := u.Client.GetOwnContact()
	if err != nil {
		return fmt.Errorf("get own contact: %w", err)
	}

	if data.Nick == defaultAccountName {
		if err := u.Client.SetProfileData(botNickname, "", ""); err != nil {
			return fmt.Errorf("set account nickname: %w", err)
		}
	}
	return nil
}
