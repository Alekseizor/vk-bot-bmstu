package ds

import "github.com/google/uuid"

type User struct {
	VkID int

	State string

	Schedule ScheduleInfo
}

type ScheduleInfo struct {
	BranchUUID     uuid.UUID
	FacultyUUID    uuid.UUID
	DepartmentUUID uuid.UUID
	GroupUUID      uuid.UUID
}
