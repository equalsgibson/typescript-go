export interface GroupResponse {
	UpdatedAt: string
	Data: group[] | null
}

export interface SystemUser {
	userID: foobar
	primaryGroup: group
	secondaryGroup?: group | null
	tags: string[] | null
}

export type foobar = number

export interface group {
	groupName: string
	UpdatedAt: string
	DeletedAt: string | null
	Data: any
	MoreData: any
}
