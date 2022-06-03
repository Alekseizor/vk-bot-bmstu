package ds

type User struct {
	VkID int

	State string

	Schedule ScheduleInfo
}

type ScheduleInfo struct {
	BranchUUID     string
	FacultyUUID    string
	DepartmentUUID string
	GroupUUID      string
}
