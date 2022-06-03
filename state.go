package main

///////////////////////////////////////////////////////////
type ChatContext struct {
	//получаем информацию о пользователе
	//используем для записи информации о выборе пользователя, на каком состоянии он находится
}

func (chc ChatContext) ChatID() string {
	//получаем VK ID, возвращаем его
	return id
}
func (chc ChatContext) Get(string) string { //получаем информацию о пользователе(либо состояние, либо uuid)
	//в стрингу(входной параметр) будем писать нужный нам атрибут из БД, возвращаем
	return uuid
}
func (chc ChatContext) Set(string, string) { //записываем информацию в бд

}
func (chc ChatContext) PayLoad() string { //записываем информацию в бд
	return msg
}

type State interface {
	Name() string              //получаем название состояния в виде строки, чтобы в дальнейшем куда-то записать(БД)
	Process(ChatContext) State //нужно взять контекст, посмотреть на каком состоянии сейчас пользователь, метод должен вернуть состояние
}

//////////////////////////////////////////////////////////
type StartState struct {
}

var _StartState = &StartState{}

func (state StartState) Process(ctx ChatContext) State {
	return _BranchState
}
func (state StartState) Name() string {
	return "StartState"
}

//////////////////////////////////////////////////////////
type BranchState struct {
}

var _BranchState = &BranchState{}

func (state BranchState) Process(ctx ChatContext) State {
	return _FacultyState
}

func (state BranchState) Name() string {
	return "BranchState"
}

//////////////////////////////////////////////////////////
type FacultyState struct {
}

var _FacultyState = &FacultyState{}

func (state FacultyState) Process(ctx ChatContext) State {
	return _DepartmentState
}

func (state FacultyState) Name() string {
	return "FacultyState"
}

//////////////////////////////////////////////////////////
type DepartmentState struct {
}

var _DepartmentState = &DepartmentState{}

func (state DepartmentState) Process(ctx ChatContext) State {
	return _GroupState
}
func (department DepartmentState) Name() string {
	return "DepartmentState"
}

/////////////////////////////////////////////////////////
type GroupState struct {
}

var _GroupState = &GroupState{}

func (state GroupState) Process(ctx ChatContext) State {
	return _FacultyState
}

func (state GroupState) Name() string {
	return "GroupState"
}

//////////////////////////////////////////////////////////
type WeekState struct {
}

var _WeekState = &WeekState{}

func (state WeekState) Process(ctx ChatContext) State {
	return _DayState
}

func (state WeekState) Name() string {
	return "WeekState"
}

///////////////////////////////////////////////////////////

type NextWeekState struct {
}

var _NextWeekState = &WeekState{}

func (state NextWeekState) Process(ctx ChatContext) State {
	return _DayState
}

func (state NextWeekState) Name() string {
	return "NextWeekState"
}

///////////////////////////////////////////////////////////

type DayState struct {
}

var _DayState = &DayState{}

func (state DayState) Process(ctx ChatContext) State {
	return _DayState
}

func (state DayState) Name() string {
	return "DayState"
}

///////////////////////////////////////////////////////////

type ErrorState struct {
}

var _ErrorState = &DayState{}

func (state ErrorState) Process(ctx ChatContext) State {
	return _DayState
}

func (state ErrorState) Name() string {
	return "ErrorState"
}

///////////////////////////////////////////////////////////
//var AuthUsersList = make(map[string]State)
