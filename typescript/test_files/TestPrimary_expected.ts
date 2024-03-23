export namespace GoGenerated {
	export type ExtendedType = {
		ID: number
		Name: string
	}

	export type GroupMapA = { [key: string]: group } | null

	export type GroupMapB = { [key: string]: group } | null

	export type GroupResponse = {
		updated_at: string
		group_map: { [key: string]: group } | null
		data: group
	}

	export type SystemUser = {
		Reports: { [key: TestUserID]: boolean } | null
		userID: TestUserID
		primaryGroup: group
		UnknownType: unknown
		secondaryGroup?: group | null
		user_tags: string[] | null
	}

	export type TestUserID = number

	export type UserResponse = {
		updated_at: string
		group_map: { [key: string]: group } | null
		data: SystemUser[] | null
	}

	export type group = {
		groupName: string
		UpdatedAt: string
		DeletedAt: string | null
		Timeout: number
		CreateAt: string
		Data: any
		MoreData: any
	}

	export const userCreate = (payload: SystemUser) => {
		return fetch("/api/user/create", {
			method: "POST",
			body: JSON.stringify(payload),
		}).then<UserResponse>((response) => response.json())
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
		}).then<UserResponse>((response) => response.json())
	}

	export const foobar: group = {
		"groupName": "hello there",
		"UpdatedAt": "0001-01-01T01:01:01.000000001Z",
		"DeletedAt": null,
		"Timeout": 0,
		"CreateAt": "0001-01-01T01:01:01.000000001Z",
		"Data": null,
		"MoreData": null
	}
}
