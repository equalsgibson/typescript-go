export namespace GoGenerated {
	export interface ExtendedType {
		ID: number
		Name: string
	}

	export type GroupMapA = { [key: string]: group } | null

	export type GroupMapB = { [key: string]: group } | null

	export interface GroupResponse {
		updated_at: string
		group_map: { [key: string]: group } | null
		data: group
	}

	export interface SystemUser {
		Reports: { [key: TestUserID]: boolean } | null
		userID: TestUserID
		primaryGroup: group
		UnknownType: unknown
		secondaryGroup?: group | null
		user_tags: string[] | null
	}

	export type TestUserID = number

	export interface UserResponse {
		updated_at: string
		group_map: { [key: string]: group } | null
		data: SystemUser[] | null
	}

	export interface group {
		groupName: string
		UpdatedAt: string
		DeletedAt: string | null
		Data: any
		MoreData: any
	}

	export const userCreate = (payload: SystemUser) => {
		return fetch("/api/user/create", {
			method: "POST",
			body: JSON.stringify(payload),
		}).then<UserResponse>((response => response.json()))
	}

	export const userGet = (userID: TestUserID) => {
		const params = {
			userID: userID,
		}

		const queryString = Object.keys(params).map((key) => {
			return encodeURIComponent(key) + "=" + encodeURIComponent(params[key])
		}).join("&")

		return fetch(`/api/user?${queryString}`, {
			method: "GET",
		}).then<UserResponse>((response => response.json()))
	}
}
