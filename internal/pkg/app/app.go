package app

import (
	"context"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
	vk_client "main/internal/app/button"
	"main/internal/app/config"
	"main/internal/app/ds"
	"main/internal/app/redis"
	"main/internal/pkg/clients/bitop"
	"strings"
)

var start_message string = "Рады приветствовать тебя у нас в сообществе, выбирай пункт меню и полетели!"

type App struct {
	// корневой контекст
	ctx context.Context
	vk  *api.VK
	lp  *longpoll.LongPoll

	vkClient    *vk_client.VkClient
	redisClient *redis.RedClient
	bitopClient *bitop.Client
}

func New(ctx context.Context) (*App, error) {
	cfg := config.FromContext(ctx)
	vk := api.NewVK(cfg.VKToken)
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.WithError(err).Error("cant get groups by id")

		return nil, err
	}

	log.WithField("group_id", group[0].ID).Info("init such group")

	c, err := redis.New(ctx)
	if err != nil {
		return nil, err
	}

	vkClient, err := vk_client.New(ctx)
	if err != nil {
		return nil, err
	}

	//starting long poll
	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		ctx:         ctx,
		vk:          vk,
		lp:          lp,
		vkClient:    vkClient,
		redisClient: c,
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	// New message event
	a.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)

		user := ds.User{
			VkID: obj.Message.PeerID,
			//Ставим начальный стейт State: startState,
		}

		messageText := obj.Message.Text

		if strings.EqualFold(messageText, "Начать") {

			err := a.redisClient.SetUser(ctx, user)
			if err != nil {
				log.WithError(err).Error("cant set user")

				return
			}

			b := params.NewMessagesSendBuilder()
			b.Message(start_message)
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			k := &object.MessagesKeyboard{}
			k.AddRow()
			k.AddTextButton("Гайд", "", "primary")
			k.AddRow()
			k.AddTextButton("Узнать расписание", "", "primary")
			k.AddRow()
			k.AddTextButton("Информация обо мне", "", "primary")
			b.Keyboard(k)

			_, err = a.vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
		} else if /* u,*/ _, err := a.redisClient.GetUser(ctx, user.VkID); obj.Message.Text != "Начать" && err == nil /*&& u.State==startState*/ {
			b := params.NewMessagesSendBuilder()
			b.Message("Если хотите начать работу с ботом, нажмите на кнопку 'Начать'")
			b.RandomID(0)
			b.PeerID(obj.Message.PeerID)
			k := &object.MessagesKeyboard{}
			k.AddRow()
			k.AddTextButton("Начать", "", "primary")
			b.Keyboard(k)
			a.vk.MessagesSend(b.Params)
		}
	})

	a.bitopClient = bitop.New(ctx)
	resp, _ := a.bitopClient.GetBranch(ctx, "")

	resp, _ = a.bitopClient.GetFaculty(ctx, resp.Items[0].Uuid)
	fmt.Println(resp.Total)

	fmt.Println()

	log.Println("Start Long Poll")
	if err := a.lp.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}
