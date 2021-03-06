package app

import (
	"context"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	log "github.com/sirupsen/logrus"
	vk_client "main/internal/app/button"
	"main/internal/app/config"
	"main/internal/app/ds"
	"main/internal/app/redis"
	"main/internal/app/state"
	"main/internal/pkg/clients/bitop"
	"strings"
)

var start_message string = "Рады приветствовать тебя у нас в сообществе, выбирай пункт меню и полетели!"
var chatcontext state.ChatContext

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
	var ScheduleUser *ds.User
	var err error
	a.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {

		messageText := obj.Message.Text
		fmt.Println(messageText)
		fmt.Println(obj.Message.PeerID)
		ScheduleUser, err = a.redisClient.GetUser(ctx, obj.Message.PeerID)
		if err != nil {
			log.WithError(err).Error("cant set user")
			return
		}
		//if the user writes for the first time, add to the database
		if ScheduleUser == nil {
			ScheduleUser = &ds.User{}
			ScheduleUser.VkID = obj.Message.PeerID
			ScheduleUser.State = "StartState"
			err := a.redisClient.SetUser(ctx, *ScheduleUser)
			if err != nil {
				log.WithError(err).Error("cant set user")
				return
			}
		} else if ScheduleUser.State == "" {
			ScheduleUser.State = "StartState"
			err := a.redisClient.SetUser(ctx, *ScheduleUser)
			if err != nil {
				log.WithError(err).Error("cant set user")
				return
			}
		}
		fmt.Println(ScheduleUser.State) //норм

		if strings.EqualFold(messageText, "Сброс") {
			ScheduleUser.State = "StartState"
			err := a.redisClient.SetUser(ctx, *ScheduleUser)
			if err != nil {
				log.WithError(err).Error("cant set user")
				return
			}
		}

		strInState := map[string]state.State{
			state.RefStartState.Name():      state.RefStartState,
			state.RefBranchState.Name():     state.RefBranchState,
			state.RefFacultyState.Name():    state.RefFacultyState,
			state.RefDepartmentState.Name(): state.RefDepartmentState,
			state.RefGroupState.Name():      state.RefGroupState,
			state.RefWeekState.Name():       state.RefWeekState,
			state.RefNextWeekState.Name():   state.RefNextWeekState,
			state.RefDayState.Name():        state.RefDayState,
			state.RefErrorState.Name():      state.RefErrorState,
		}
		ctc := state.ChatContext{
			ScheduleUser,
			a.vk,
			a.redisClient,
			&ctx,
			a.bitopClient,
		}

		step := strInState[ScheduleUser.State]
		fmt.Println(step.Name())
		nextStep := step.Process(ctc, messageText)
		fmt.Println(nextStep.Name())
		ScheduleUser.State = nextStep.Name()
		err = a.redisClient.SetUser(ctx, *ScheduleUser)
		if err != nil {
			log.WithError(err).Error("cant set user")
			return
		}

		//ScheduleUser.State = nextStep.Name()
		//err = a.redisClient.SetUser(ctx, *ScheduleUser)
		//if err != nil {
		//	log.WithError(err).Error("cant set user")
		//	return
		//}
		/*messageText := obj.Message.Text
		ScheduleUser := &ds.User{}
		//check if we have such a user
		ScheduleUser, err := a.redisClient.GetUser(ctx, obj.Message.PeerID)
		if err != nil {
			log.WithError(err).Error("cant set user")

			return
		}
		//if the user writes for the first time, add to the database
		if ScheduleUser == nil {
			ScheduleUser.VkID = obj.Message.PeerID
			ScheduleUser.State = "StartState"
			err := a.redisClient.SetUser(ctx, *ScheduleUser)
			if err != nil {
				log.WithError(err).Error("cant set user")
				return
			}
		}
		//to get states
		strInState := map[string]state.State{
			state.RefStartState.Name():      state.RefStartState,
			state.RefBranchState.Name():     state.RefBranchState,
			state.RefFacultyState.Name():    state.RefFacultyState,
			state.RefDepartmentState.Name(): state.RefDepartmentState,
			state.RefGroupState.Name():      state.RefGroupState,
			state.RefWeekState.Name():       state.RefWeekState,
			state.RefNextWeekState.Name():   state.RefNextWeekState,
			state.RefDayState.Name():        state.RefDayState,
			state.RefErrorState.Name():      state.RefErrorState,
		}
		if strings.EqualFold(messageText, "Сброс") {
			ScheduleUser.State = "StartState"
			err := a.redisClient.SetUser(ctx, *ScheduleUser)
			if err != nil {
				log.WithError(err).Error("cant set user")
				return
			}
		}
		ctc := state.ChatContext{
			ScheduleUser,
			a.vk,
			a.redisClient,
			a.ctx,
			a.bitopClient,
		}

		step := strInState[ScheduleUser.State]
		fmt.Println(step.Name())
		nextStep := step.Process(ctc, messageText)
		ScheduleUser.State = nextStep.Name()
		err = a.redisClient.SetUser(ctx, *ScheduleUser)
		if err != nil {
			log.WithError(err).Error("cant set user")
			return
		}
		*/
	})
	log.Println("Start Long Poll")
	if err := a.lp.Run(); err != nil {
		log.Fatal(err)
		return nil
	}
	return nil
}
