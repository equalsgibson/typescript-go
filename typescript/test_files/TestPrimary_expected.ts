export type GroupMapA = Map<string, group> | null

export type GroupMapB = Map<string, group> | null

export interface GroupResponse {
	updated_at: string
	group_map: Map<string, group> | null
	data: group
}

export interface SystemUser {
	Reports: Map<TestUserID, boolean> | null
	userID: TestUserID
	primaryGroup: group
	UnknownType: unknown
	secondaryGroup?: group | null
	user_tags: string[] | null
}

export type TestUserID = number

export interface UserResponse {
	updated_at: string
	group_map: Map<string, group> | null
	data: SystemUser[] | null
}

export interface group {
	groupName: string
	UpdatedAt: string
	DeletedAt: string | null
	Data: any
	MoreData: any
}
