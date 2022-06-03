package main

import "github.com/SevereCloud/vksdk/v2/api"

type S int

const (
	Main S = S(iota)
	SelectBranch
	SelectFaculty
	SelectDepartment
	SelectGroup
)

type ChatContext interface {
	Chat() *api.VK
	Get(string) interface{}
	Set(string, interface{})
}

type State interface {
	Process(ChatContext) State
}

type StartState struct {
}

func (s *StartState) Process(ctx ChatContext) State {
	//send buttons layout (all branches)

	return
}

var Chains = map[S][]S{
	Main: {SelectBranch},
}

type MainClass struct {
	chains []S
}
