package main

import (
	"errors"

	utopiago "github.com/Sagleft/utopialib-go"
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

	Client utopiago.UtopiaClient
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

	u.Client = utopiago.UtopiaClient{
		Protocol: protocol,
		Token:    u.Token,
		Host:     u.Host,
		Port:     u.Port,
	}
	if !u.Client.CheckClientConnection() {
		return errors.New("failed to connect to Utopia client")
	}
	return nil
}

func (u *utopiaService) postMedia(media mediaPost) error {
	// TODO
	return nil
}
