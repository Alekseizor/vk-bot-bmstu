package main

type State interface {
	clickButton()
	clickUndo()
	changeSchedule(schedule *Schedule)
}

type Schedule struct {
	state State
}

func (schedule *Schedule) changeState(state State) {

}

type DayState struct {
	schedule Schedule
}

func (day DayState) clickButton() {

}
func (day DayState) clickUndo() {

}
func (day DayState) changeSchedule(schedule *Schedule) {

}

type WeekState struct {
	schedule Schedule
}

func (week WeekState) clickButton() {

}
func (week WeekState) clickUndo() {

}
func (week WeekState) changeSchedule(schedule *Schedule) {

}

type DepartmentState struct {
	schedule Schedule
}

func (department DepartmentState) clickButton() {

}
func (department DepartmentState) clickUndo() {

}
func (department DepartmentState) changeSchedule(schedule *Schedule) {

}
func main() {

}
