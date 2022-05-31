package main

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	longpoll "github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
)

func main() {
	vk := api.NewVK("e30a0a7d544a45c5b0483924b964f1848644f5ae3e27112ccd4870a4cf5d8e59631beaac913619cdc14fa")

	//Устанавливаем longpoll
	lp, err := longpoll.NewLongPoll(vk, 213613932)
	if err != nil {
		log.Println(err)
	}

	//Первое сообщение + блок на ввод с руки
	/*lp.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		b := params.NewMessagesSendBuilder()

		keyboard := object.NewMessagesKeyboard(true)
		keyboard.AddRow()
		keyboard.AddCallbackButton("МГТУ", "", "primary")
		keyboard.AddCallbackButton("КФ МГТУ", "", "primary")

		log.Println(obj.Message.Text)
		b.RandomID(0)
		b.Keyboard(keyboard)
		b.Message("Выбери пункт меню")
		b.PeerID(obj.Message.PeerID)
		_, err := vk.MessagesSend(b.Params)
		if err != nil {
			log.Println(err)
		}
	})*/

	//Инициализация контекста (паттерн Состояние)q
	//schedule := Schedule{}
	//st := &StartState{schedule: schedule}
	//schedule.state = st

	//schedule.state.start(lp, vk)

	//Обработка CallBack кнопок
	/*go lp.MessageEvent(func(ctx context.Context, obj events.MessageEventObject) {
		m := params.NewMessagesSendMessageEventAnswerBuilder()
		m.PeerID(obj.PeerID)
		m.EventID(obj.EventID)
		m.UserID(obj.UserID)
		m.EventData(ShowSnackbar("Отличный выбор"))
		_, err := vk.MessagesSendMessageEventAnswer(m.Params)
		if err != nil {
			log.Println(err)
		}
	})*/

	//Обработка события для подписки пользователя
	go lp.GroupJoin(func(ctx context.Context, obj events.GroupJoinObject) {
		b := params.NewMessagesSendBuilder()

		b.RandomID(0)
		b.Message("Добро пожаловать в нашу группу. Здесь ты можешь узнать свое расписание в самом лучшем вузе Казахстана!")
		b.PeerID(obj.UserID)
		_, err := vk.MessagesSend(b.Params)
		if err != nil {
			log.Println(err)
		}
	})

	lp.Run()
}
