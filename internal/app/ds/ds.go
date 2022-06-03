package ds

import "github.com/google/uuid"

type State string

type User struct {
	VkID int

	Memory string

	Schedule ScheduleInfo
}

type ScheduleInfo struct {
	BranchUUID     uuid.UUID
	FacultyUUID    uuid.UUID
	DepartmentUUID uuid.UUID
	GroupUUID      uuid.UUID
}
