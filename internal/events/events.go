package events

const (
	//Swap out the names from like CREATE_GUILD TO GUILD_CREATE
	GUILD_CREATE = "GUILD_CREATE"
	GUILD_DELETE = "GUILD_DELETE"
	GUILD_UPDATE = "GUILD_UPDATE"

	INVITE_CREATE = "INVITE_CREATE"
	INVITE_DELETE = "INVITE_DELETE"

	MESSAGE_CREATE = "MESSAGE_CREATE"
	MESSAGE_DELETE = "MESSAGE_DELETE"
	MESSAGE_UPDATE = "MESSAGE_UPDATE"

	MESSAGES_USER_CLEAR  = "MESSAGES_CLEAR"
	MESSAGES_GUILD_CLEAR = "MESSAGES_GUILD_CLEAR"

	DM_CREATE    = "DM_CREATE"
	DM_DELETE    = "DM_DELETE"
	TYPING_START = "TYPING_START"

	USER_FRIEND_REQUEST_ADD    = "USER_FRIEND_REQUEST_ADD"
	USER_FRIEND_REQUEST_REMOVE = "USER_FRIEND_REQUEST_REMOVE"

	USER_FRIEND_INCOMING_REQUEST_ADD    = "USER_FRIEND_INCOMING_REQUEST_ADD"
	USER_FRIEND_INCOMING_REQUEST_REMOVE = "USER_FRIEND_INCOMING_REQUEST_REMOVE"

	USER_BLOCKED_ADD    = "USER_BLOCKED_ADD"
	USER_BLOCKED_REMOVE = "USER_BLOCKED_REMOVE"

	USER_FRIEND_ADD    = "USER_FRIEND_ADD"
	USER_FRIEND_REMOVE = "USER_FRIEND_REMOVE"

	MEMBER_ADD    = "MEMBER_ADD"
	MEMBER_REMOVE = "MEMBER_REMOVE"

	MEMBER_BAN_ADD    = "MEMBER_BAN_ADD"
	MEMBER_BAN_REMOVE = "MEMBER_BAN_REMOVE"

	MEMBER_ADMIN_ADD    = "MEMBER_ADMIN_ADD"
	MEMBER_ADMIN_REMOVE = "MEMBER_ADMIN_REMOVE"

	LOG_OUT = "LOG_OUT"

	USER_INFO_UPDATE = "USER_INFO_UPDATE"
)
