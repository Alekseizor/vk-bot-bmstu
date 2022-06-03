package app

import (
	"context"
	"main/internal/app/config"
	"strings"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	log "github.com/sirupsen/logrus"
)

type App struct {
	// корневой контекст
	ctx context.Context
	vk  *api.VK
	lp  *longpoll.LongPoll

	redisClient *redis.Client
}

func New(ctx context.Context) (*App, error) {
	cfg := config.FromContext(ctx)
	vk := api.NewVK(cfg.Token)

	c, err := redis.New(ctx)
	if err != nil {
		return nil, err
	}

	// get information about the group
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.WithError(err).Error("cant get groups by id")

		return nil, err
	}

	log.WithField("group_id", group[0].ID).Info("init such group")

	log.Info("start init longpool")

	// Initializing Long Poll
	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		ctx:         ctx,
		vk:          vk,
		lp:          lp,
		redisClient: c,
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	// New message event
	a.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)

		messageText := obj.Message.Text

		if strings.HasPrefix(messageText, "запомни") {
			memoryString := messageText[13:]
			user := ds.User{
				VkID:   obj.Message.PeerID,
				Memory: memoryString,
			}

			err := a.redisClient.SetUser(ctx, user)
			if err != nil {
				log.WithError(err).Error("cant set user")

				return
			}

			b := params.NewMessagesSendBuilder()
			b.Message("запомнил")
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)

			_, err = a.vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
		}

		if messageText == "вспомни" {

			user, err := a.redisClient.GetUser(ctx, obj.Message.PeerID)
			if err != nil {
				log.WithError(err).Error("cant set user")

				return
			}

			b := params.NewMessagesSendBuilder()
			b.Message(user.Memory)
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)

			_, err = a.vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
		}

		if messageText == "ping" {
			b := params.NewMessagesSendBuilder()
			b.Message("pong")
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)

			_, err := a.vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	// Run Bots Long Poll
	log.Info("Start Long Poll")
	if err := a.lp.Run(); err != nil {
		return err
	}

	return nil
}
