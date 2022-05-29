package main

///////////////////////////////////////////////////////////
type State interface {
	clickButton()
	clickUndo()
	changeSchedule(schedule Schedule)
}

//////////////////////////////////////////////////////////
type StartState struct {
	schedule Schedule
}

func (state *StartState) clickButton() {

}

func (state *StartState) clickUndo() {

}

func (state *StartState) changeSchedule(schedule Schedule) {

}

//////////////////////////////////////////////////////////
type BranchState struct {
	schedule Schedule
}

func (state *BranchState) clickButton() {

}

func (state *BranchState) clickUndo() {

}

func (state *BranchState) changeSchedule(schedule Schedule) {

}

//////////////////////////////////////////////////////////
type GroupState struct {
	schedule Schedule
}

func (state *GroupState) clickButton() {

}

func (state *GroupState) clickUndo() {

}

func (state *GroupState) changeSchedule(schedule Schedule) {

}

//////////////////////////////////////////////////////////
type FacultyState struct {
	schedule Schedule
}

func (state *FacultyState) clickButton() {

}

func (state *FacultyState) clickUndo() {

}

func (state *FacultyState) changeSchedule(schedule Schedule) {

}

///////////////////////////////////////////////////////////
type Schedule struct {
	state State
}

func (schedule *Schedule) changeState(state State) {
}

///////////////////////////////////////////////////////////
func main() {

}
