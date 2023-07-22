export type GroupMapA = Map<string, group> | null

export type GroupMapB = Map<string, group> | null

export interface GroupResponse {
	UpdatedAt: string
	Map: Map<string, group> | null
	Data: group[] | null
}

export interface SystemUser {
	Reports: Map<foobar, boolean> | null
	userID: foobar
	primaryGroup: group
	X: unknown
	secondaryGroup?: group | null
	tags: string[] | null
}

export interface Thing {
	UpdatedAt: string
	Map: Map<string, group> | null
	Data: SystemUser[] | null
}

export type foobar = number

export interface group {
	groupName: string
	UpdatedAt: string
	DeletedAt: string | null
	Data: any
	MoreData: any
}
