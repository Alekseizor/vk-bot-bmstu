package state

import (
	"context"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	log "github.com/sirupsen/logrus"
	"main/internal/app/ds"
	"main/internal/app/model"
	"main/internal/app/redis"
	"main/internal/pkg/clients/bitop"
	"strconv"
	"strings"
	"time"
)

///////////////////////////////////////////////////////////
type ChatContext struct {
	User        *ds.User
	Vk          *api.VK
	RedisClient *redis.RedClient
	Ctx         *context.Context
	BitopClient *bitop.Client
	//получаем информацию о пользователе
	//используем для записи информации о выборе пользователя, на каком состоянии он находится
}

func (chc ChatContext) ChatID() int {
	return chc.User.VkID
}
func (chc ChatContext) Get(VkID int, Field string) string { //получаем информацию о пользователе(либо состояние, либо uuid)
	//в стрингу(входной параметр) будем писать нужный нам атрибут из БД, возвращаем
	var err error
	chc.User, err = chc.RedisClient.GetUser(*chc.Ctx, VkID)
	if err != nil {
		log.Println("Failed to get record")
		log.Fatal(err)
	}
	if Field == "State" {
		return chc.User.State
	}
	if Field == "BranchUUID" {
		return chc.User.BranchUUID
	}
	if Field == "FacultyUUID" {
		return chc.User.FacultyUUID
	}
	if Field == "DepartmentUUID" {
		return chc.User.DepartmentUUID
	}
	if Field == "GroupUUID" {
		return chc.User.GroupUUID
	}
	if Field == "IsNumerator" {
		return strconv.FormatBool(chc.User.IsNumerator)
	}

	return "not found"
}
func (chc ChatContext) Set() { //записываем информацию в бд
	err := chc.RedisClient.SetUser(*chc.Ctx, *chc.User)
	if err != nil {
		log.WithError(err).Error("cant set user")
		return
	}
}

type State interface {
	Name() string                      //получаем название состояния в виде строки, чтобы в дальнейшем куда-то записать(БД)
	Process(ChatContext, string) State //нужно взять контекст, посмотреть на каком состоянии сейчас пользователь, метод должен вернуть состояние
}

//////////////////////////////////////////////////////////
type StartState struct {
}

var RefStartState = &StartState{}

func (state StartState) Process(ctc ChatContext, messageText string) State {
	if messageText == "Узнать расписание" {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Укажи свою группу")
		b.PeerID(ctc.User.VkID)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		return RefGroupState
	} else {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.PeerID(ctc.User.VkID)
		b.Message("Рады приветствовать тебя у нас в сообществе, давай найдем твоё расписание!")
		k := &object.MessagesKeyboard{}
		k.AddRow()
		k.AddTextButton("Узнать расписание", "", "primary")
		b.Keyboard(k)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		return RefStartState
	}
}

func (state StartState) Name() string {
	return "StartState"
}

//////////////////////////////////////////////////////////
type BranchState struct {
}

var RefBranchState = &BranchState{}

func (state BranchState) Process(ctc ChatContext, messageText string) State {

	resp, _ := ctc.BitopClient.GetBranch(*ctc.Ctx, "messageText")
	if resp == nil {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.Message("Рады приветствовать тебя у нас в сообществе, давай найдем твоё расписание!")
		k := &object.MessagesKeyboard{}
		k.AddRow()
		k.AddTextButton("Узнать расписание", "", "primary")
		b.Keyboard(k)
		b.RandomID(0)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		return RefBranchState
	} else {

	}
	return RefFacultyState
}

func (state BranchState) Name() string {
	return "BranchState"
}

//////////////////////////////////////////////////////////
type FacultyState struct {
}

var RefFacultyState = &FacultyState{}

func (state FacultyState) Process(ctx ChatContext, messageText string) State {
	return RefDepartmentState
}

func (state FacultyState) Name() string {
	return "FacultyState"
}

//////////////////////////////////////////////////////////
type DepartmentState struct {
}

var RefDepartmentState = &DepartmentState{}

func (state DepartmentState) Process(ctx ChatContext, messageText string) State {
	return RefGroupState
}
func (department DepartmentState) Name() string {
	return "DepartmentState"
}

/////////////////////////////////////////////////////////
type GroupState struct {
}

var RefGroupState = &GroupState{}

func (state GroupState) Process(ctc ChatContext, messageText string) State {
	ctc.BitopClient = bitop.New(*ctc.Ctx)
	resp, _ := ctc.BitopClient.GetGroup(*ctc.Ctx, messageText)
	if resp.Total > 1 {
		for _, group := range resp.Items {
			if group.Caption == strings.ToUpper(messageText) {
				ctc.User.GroupUUID = resp.Items[0].Uuid
				b := params.NewMessagesSendBuilder()
				b.PeerID(ctc.User.VkID)
				b.RandomID(0)
				b.Message("Выберите интересующую Вас неделю")
				k := &object.MessagesKeyboard{}
				k.AddRow()
				k.AddTextButton("Числитель", "", "primary")
				k.AddRow()
				k.AddTextButton("Знаменатель", "", "primary")
				b.Keyboard(k)
				_, err := ctc.Vk.MessagesSend(b.Params)
				if err != nil {
					log.Fatal(err)
				}
				return RefWeekState
			}
		}
		b := params.NewMessagesSendBuilder()
		b.PeerID(ctc.User.VkID)
		b.RandomID(0)
		b.Message("Введите полное название группы")
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		return RefGroupState
	}
	if resp.Total == 1 {
		ctc.User.GroupUUID = resp.Items[0].Uuid
		b := params.NewMessagesSendBuilder()
		b.PeerID(ctc.User.VkID)
		b.RandomID(0)
		b.Message("Выберите интересующую Вас неделю")
		k := &object.MessagesKeyboard{}
		k.AddRow()
		k.AddTextButton("Числитель", "", "primary")
		k.AddRow()
		k.AddTextButton("Знаменатель", "", "primary")
		b.Keyboard(k)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		return RefWeekState
	}
	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.PeerID(ctc.User.VkID)
	b.Message("Введите нужную группу повторно, не удалось найти")
	_, err := ctc.Vk.MessagesSend(b.Params)
	if err != nil {
		log.Fatal(err)
	}
	return RefGroupState
}

func (state GroupState) Name() string {
	return "GroupState"
}

//////////////////////////////////////////////////////////
type WeekState struct {
}

var RefWeekState = &WeekState{}

func (state WeekState) Process(ctc ChatContext, messageText string) State {
	if messageText == "Числитель" {
		ctc.User.IsNumerator = true
		b := params.NewMessagesSendBuilder()
		b.PeerID(ctc.User.VkID)
		b.RandomID(0)
		b.Message("Выберите нужный день недели")
		k := &object.MessagesKeyboard{}
		k.AddRow()
		k.AddTextButton("Понедельник", "", "primary")
		k.AddRow()
		k.AddTextButton("Вторник", "", "primary")
		k.AddRow()
		k.AddTextButton("Среда", "", "primary")
		k.AddRow()
		k.AddTextButton("Четверг", "", "primary")
		k.AddRow()
		k.AddTextButton("Пятница", "", "primary")
		k.AddRow()
		k.AddTextButton("Суббота", "", "primary")
		b.Keyboard(k)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		return RefDayState
	} else if messageText == "Знаменатель" {
		ctc.User.IsNumerator = false
		b := params.NewMessagesSendBuilder()
		b.PeerID(ctc.User.VkID)
		b.RandomID(0)
		b.Message("Выберите нужный день недели")
		k := &object.MessagesKeyboard{}
		k.AddRow()
		k.AddTextButton("Понедельник", "", "primary")
		k.AddRow()
		k.AddTextButton("Вторник", "", "primary")
		k.AddRow()
		k.AddTextButton("Среда", "", "primary")
		k.AddRow()
		k.AddTextButton("Четверг", "", "primary")
		k.AddRow()
		k.AddTextButton("Пятница", "", "primary")
		k.AddRow()
		k.AddTextButton("Суббота", "", "primary")
		b.Keyboard(k)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		return RefDayState
	} else {
		b := params.NewMessagesSendBuilder()
		b.PeerID(ctc.User.VkID)
		b.RandomID(0)
		b.Message("Выберите интересующую Вас неделю")
		k := &object.MessagesKeyboard{}
		k.AddRow()
		k.AddTextButton("Числитель", "", "primary")
		k.AddRow()
		k.AddTextButton("Знаменатель", "", "primary")
		b.Keyboard(k)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		return RefWeekState
	}
}

func (state WeekState) Name() string {
	return "WeekState"
}

///////////////////////////////////////////////////////////

type NextWeekState struct {
}

var RefNextWeekState = &WeekState{}

func (state NextWeekState) Process(ctc ChatContext, messageText string) State {
	return NextWeekState{}
}

func (state NextWeekState) Name() string {
	return "NextWeekState"
}

///////////////////////////////////////////////////////////

type DayState struct {
}

var RefDayState = &DayState{}

func (state DayState) Process(ctc ChatContext, messageText string) State {
	var v string
	if (messageText == "Понедельник") || (messageText == "Вторник") || (messageText == "Среда") || (messageText == "Четверг") || (messageText == "Пятница") || (messageText == "Суббота") {
		fmt.Print("Я здесь")
		ctc.BitopClient = bitop.New(*ctc.Ctx)
		Schedule, err := ctc.BitopClient.GetSchedule(*ctc.Ctx, ctc.User.GroupUUID, ctc.User.IsNumerator, messageText)
		if err != nil {
			log.WithError(err).Error("failed to get schedule")
		}
		if Schedule == nil {
			v := "В этот день Вы отдыхаете!"
			b := params.NewMessagesSendBuilder()
			b.PeerID(ctc.User.VkID)
			b.RandomID(0)
			k := &object.MessagesKeyboard{}
			k.AddRow()
			k.AddTextButton("Сброс", "", "primary")
			k.AddRow()
			k.AddTextButton("Вернуться к выбору дня", "", "primary")
			k.AddRow()
			k.AddTextButton("Вернуться к выбору недели", "", "primary")
			b.Message(v)
			b.Keyboard(k)
			_, err := ctc.Vk.MessagesSend(b.Params)
			if err != nil {
				log.Fatal(err)
			}
			return RefDayState
		}
		var lessons []model.Lesson
		var less model.Lesson
		var teach model.Teacher
		var teachs model.Teachers
		for _, lesson := range Schedule.Lessons {
			less.Name = lesson.Name
			less.Cabinet = lesson.Cabinet
			less.Type = lesson.Type
			for _, teacher := range lesson.Teachers {
				teach.Name = teacher.Name
				teachs = append(teachs, teach)
			}
			less.Teachers = teachs
			teachs = nil
			less.StartAt = lesson.StartAt
			less.EndAt = lesson.EndAt
			less.Day = lesson.Day
			less.IsNumerator = lesson.IsNumerator
			lessons = append(lessons, less)
		}
		lessons = quickSort(&lessons)
		switch messageText {
		case "Понедельник":
			{
				v = "Ваше расписание на понедельник:\n\n"
			}
		case "Вторник":
			{
				v = "Ваше расписание на вторник:\n\n"
			}
		case "Среда":
			{
				v = "Ваше расписание на среду:\n\n"
			}
		case "Четверг":
			{
				v = "Ваше расписание на четверг:\n\n"
			}
		case "Пятница":
			{
				v = "Ваше расписание на пятницу:\n\n"
			}
		case "Суббота":
			{
				v = "Ваше расписание на субботу:\n\n"
			}
		}
		for _, lesson := range lessons {
			v += ("\t\tПредмет: " + lesson.Name + "\n\n")
			v += ("\t\t\t\tВремя начала пары: " + lesson.StartAt[0:5] + "\n")
			v += ("\t\t\t\tВремя окончания пары: " + lesson.EndAt[0:5] + "\n")
			if (lesson.Type) != "" {
				v += ("\t\t\t\tТип занятия: " + lesson.Type + "\n")
			}
			for _, teacher := range lesson.Teachers {
				v += ("\t\t\t\tПреподаватель: " + teacher.Name + "\n")
			}
			if (lesson.Cabinet) != "" {
				v += ("\t\t\t\tКабинет: " + lesson.Cabinet + "\n")
			}
			v += "\n\n"
		}
		b := params.NewMessagesSendBuilder()
		b.PeerID(ctc.User.VkID)
		b.RandomID(0)
		b.Message(v)
		k := &object.MessagesKeyboard{}
		k.AddRow()
		k.AddTextButton("Сброс", "", "primary")
		k.AddRow()
		k.AddTextButton("Вернуться к выбору дня", "", "primary")
		k.AddRow()
		k.AddTextButton("Вернуться к выбору недели", "", "primary")
		b.Keyboard(k)
		_, err = ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		return RefDayState
	} else if messageText == "Сброс" {
		b := params.NewMessagesSendBuilder()
		b.PeerID(ctc.User.VkID)
		b.RandomID(0)
		k := &object.MessagesKeyboard{}
		k.AddRow()
		k.AddTextButton("Узнать расписание", "", "primary")
		b.Keyboard(k)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		return RefStartState
	} else if messageText == "Вернуться к выбору дня" {
		b := params.NewMessagesSendBuilder()
		b.PeerID(ctc.User.VkID)
		b.RandomID(0)
		b.Message("Выберите нужный день недели")
		k := &object.MessagesKeyboard{}
		k.AddRow()
		k.AddTextButton("Понедельник", "", "primary")
		k.AddRow()
		k.AddTextButton("Вторник", "", "primary")
		k.AddRow()
		k.AddTextButton("Среда", "", "primary")
		k.AddRow()
		k.AddTextButton("Четверг", "", "primary")
		k.AddRow()
		k.AddTextButton("Пятница", "", "primary")
		k.AddRow()
		k.AddTextButton("Суббота", "", "primary")
		b.Keyboard(k)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		return RefDayState
	} else if messageText == "Вернуться к выбору недели" {
		b := params.NewMessagesSendBuilder()
		b.PeerID(ctc.User.VkID)
		b.RandomID(0)
		b.Message("Выберите интересующую Вас неделю")
		k := &object.MessagesKeyboard{}
		k.AddRow()
		k.AddTextButton("Числитель", "", "primary")
		k.AddRow()
		k.AddTextButton("Знаменатель", "", "primary")
		b.Keyboard(k)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		return RefWeekState
	} else {
		b := params.NewMessagesSendBuilder()
		v := "Проверьте правильность ввода введенного учебного дня"
		b.PeerID(ctc.User.VkID)
		b.RandomID(0)
		b.Message(v)
		k := &object.MessagesKeyboard{}
		k.AddRow()
		k.AddTextButton("Сброс", "", "primary")
		k.AddRow()
		k.AddTextButton("Вернуться к выбору дня", "", "primary")
		b.Keyboard(k)
		_, err := ctc.Vk.MessagesSend(b.Params)
		if err != nil {
			log.Fatal(err)
		}
		return RefDayState
	}
}

func (state DayState) Name() string {
	return "DayState"
}

///////////////////////////////////////////////////////////

type ErrorState struct {
}

var RefErrorState = &ErrorState{}

func (state ErrorState) Process(ctx ChatContext, messageText string) State {
	return RefStartState
}

func (state ErrorState) Name() string {
	return "ErrorState"
}

///////////////////////////////////////////////////////////

func quickSort(lessons *[]model.Lesson) []model.Lesson {
	var lessonl, lessone, lessonm []model.Lesson
	if (len(*lessons) == 1) || (len(*lessons) == 0) {
		return *lessons
	}
	randomTime := (*lessons)[0].StartAt
	randomTimetime, _ := time.Parse("15:04:05", randomTime)
	fmt.Println(randomTimetime)
	for _, lesson := range *lessons {
		Timetime, _ := time.Parse("15:04:05", lesson.StartAt)
		fmt.Println(Timetime)
		if Timetime.Before(randomTimetime) { //если ли Timetime раньше randomTimetime
			lessonl = append(lessonl, lesson)
		} else if Timetime.After(randomTimetime) {
			lessonm = append(lessonm, lesson)
		} else {
			lessone = append(lessone, lesson)
		}
	}
	finalLessonsl := quickSort(&lessonl)
	for _, lesson := range lessone {
		finalLessonsl = append(finalLessonsl, lesson)
	}
	finalLessonsm := quickSort(&lessonm)
	for _, lesson := range finalLessonsm {
		finalLessonsl = append(finalLessonsl, lesson)
	}
	return finalLessonsl
}
