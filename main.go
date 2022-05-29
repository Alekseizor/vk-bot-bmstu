package main

type State interface {
	clickButton()
	clickUndo()
	changeSchedule(schedule Schedule)
}

type Schedule struct {
	state State
}

func (schedule *Schedule) changeState(state State) {

}

func main() {

}
