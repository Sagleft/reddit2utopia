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

type utopiaService struct {
	Conn                *uchatbot.ChatBot
	ConnEstablishedOnce bool
	Pubkey              string
	AccountNickname     string
}

func utopiaConnect(
	cfg utopiaConfig,
	nickname string,
	chats []uchatbot.Chat,
) (*utopiaService, error) {
	srv := &utopiaService{
		AccountNickname: nickname,
	}

	protocol := "http"
	if cfg.HTTPSEnabled {
		protocol += "s"
	}

	var err error
	srv.Conn, err = uchatbot.NewChatBot(uchatbot.ChatBotData{
		Config: utopiago.Config{
			Protocol: protocol,
			Host:     cfg.Host,
			Token:    cfg.Token,
			Port:     cfg.Port,
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
		ErrorCallback:    srv.onError,
	})
	if err != nil {
		return nil, fmt.Errorf("create chatbot: %w", err)
	}

	return srv, nil
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

	if u.AccountNickname != "" && data.Nick != u.AccountNickname {
		log.Println("update account name..")
		if err := u.Conn.SetAccountNickname(u.AccountNickname); err != nil {
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
