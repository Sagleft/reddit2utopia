package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/Sagleft/uchatbot-engine"
	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
	"github.com/fatih/color"
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
	Token           string
	Host            string
	Port            int
	HTTPSEnabled    bool
	ChannelID       string
	ChannelPassword string

	Conn                *uchatbot.ChatBot
	ConnEstablishedOnce bool
	Pubkey              string
}

func newUtopiaService() *utopiaService {
	return &utopiaService{}
}

func (u *utopiaService) setChannelID(ID, password string) *utopiaService {
	u.ChannelID = ID
	u.ChannelPassword = password
	return u
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

	chats := []uchatbot.Chat{}
	if u.ChannelID != "" {
		chats = append(chats, uchatbot.Chat{
			ID:       u.ChannelID,
			Password: u.ChannelPassword,
		})
	}

	var err error
	u.Conn, err = uchatbot.NewChatBot(uchatbot.ChatBotData{
		Config: utopiago.Config{
			Protocol: protocol,
			Host:     u.Host,
			Token:    u.Token,
			Port:     u.Port,
		},
		Chats: chats,
		Callbacks: uchatbot.ChatBotCallbacks{
			OnContactMessage:        func(im structs.InstantMessage) {},
			OnChannelMessage:        func(wcm structs.WsChannelMessage) {},
			OnPrivateChannelMessage: func(wcm structs.WsChannelMessage) {},
			WelcomeMessage: func(userPubkey string) string {
				return welcomeMessage
			},
		},
		UseErrorCallback: true,
		DisableEvents:    true,
		ErrorCallback:    u.onError,
	})
	return err
}

func (u *utopiaService) onError(err error) {
	if err == nil {
		return
	}

	if strings.Contains(err.Error(), errConnectionMessage) {
		if !u.ConnEstablishedOnce {
			log.Println("wait for reconnect to Utopia client..")

			u.ConnEstablishedOnce = true
			return
		}
	}

	color.Red(err.Error())
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
	_, err = u.Conn.GetClient().SendChannelPicture(channelID, base64Image, media.Text, "photo.jpg")
	return err
}

func (u *utopiaService) updateAccountName() error {
	data, err := u.Conn.GetClient().GetOwnContact()
	if err != nil {
		return fmt.Errorf("get own contact: %w", err)
	}

	if data.Nick == defaultAccountName {
		log.Println("update account name..")
		if err := u.Conn.SetAccountNickname(botNickname); err != nil {
			return fmt.Errorf("set account nickname: %w", err)
		}
	}
	return nil
}

func (u *utopiaService) loadBotPubkey() error {
	var err error
	u.Pubkey, err = u.Conn.GetOwnPubkey()
	if err != nil {
		return fmt.Errorf("get own pubkey: %w", err)
	}

	log.Printf("bot pubkey: %s\n", u.Pubkey)
	return nil
}
