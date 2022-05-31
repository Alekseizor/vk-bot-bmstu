package main

type S int

const (
	Main S = S(iota)
	SelectBranch
	SelectFaculty
	SelectDepartment
	SelectGroup
)

///////////////////////////////////////////////////////////
type State interface {
	start()
	handler()
	clickUndo()
	changeSchedule()
}

var Chains = map[S][]S{
	Main: {SelectBranch},
}

type MainClass struct {
	chains []S
}

//////////////////////////////////////////////////////////
type StartState struct {
	schedule Schedule
}

func (state StartState) start() {
}
func (state StartState) clickUndo() {

}

func (state StartState) handler() {

}

func (state StartState) changeSchedule() {

}

//////////////////////////////////////////////////////////
type BranchState struct {
	schedule Schedule
}

func (state BranchState) start() {

}

func (state BranchState) clickUndo() {

}

func (state BranchState) changeSchedule() {

}

//////////////////////////////////////////////////////////
type GroupState struct {
	schedule Schedule
}

func (state GroupState) start() {

}

func (state GroupState) clickUndo() {

}

func (state GroupState) changeSchedule() {

}

//////////////////////////////////////////////////////////
type FacultyState struct {
	schedule Schedule
}

func (state FacultyState) start() {

}

func (state FacultyState) clickUndo() {

}

func (state FacultyState) changeSchedule() {

}

///////////////////////////////////////////////////////////
type Schedule struct {
	state State
}

func (schedule *Schedule) changeState(state State) {
}

type DayState struct {
	schedule Schedule
}

func (day DayState) start() {

}
func (day DayState) clickUndo() {

}
func (day DayState) changeSchedule() {

}

type WeekState struct {
	schedule Schedule
}

func (week WeekState) start() {

}
func (week WeekState) clickUndo() {

}
func (week WeekState) changeSchedule() {

}

type DepartmentState struct {
	schedule Schedule
}

func (department DepartmentState) start() {

}
func (department DepartmentState) clickUndo() {

}
func (department DepartmentState) changeSchedule() {

}

///////////////////////////////////////////////////////////
