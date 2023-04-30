package events

const (
	CREATE_GUILD = "CREATE_GUILD"
	DELETE_GUILD = "DELETE_GUILD"
	UPDATE_GUILD = "UPDATE_GUILD"

	NOT_OWNER = "NOT_OWNER" //changes owner (removes old)
	NEW_OWNER = "NEW_OWNER" //changes owner (gives new)

	CREATE_INVITE = "CREATE_INVITE"
	DELETE_INVITE = "DELETE_INVITE"

	CREATE_GUILD_MESSAGE = "CREATE_GUILD_MESSAGE"
	DELETE_GUILD_MESSAGE = "DELETE_GUILD_MESSAGE"
	UPDATE_GUILD_MESSAGE = "UPDATE_GUILD_MESSAGE"

	CREATE_DM = "CREATE_DM"
	DELETE_DM = "DELETE_DM"

	CREATE_DM_MESSAGE      = "CREATE_DM_MESSAGE"
	DELETE_DM_MESSAGE      = "DELETE_DM_MESSAGE"
	UPDATE_DM_MESSAGE      = "UPDATE_DM_MESSAGE"
	CLEAR_USER_DM_MESSAGES = "CLEAR_USER_DM_MESSAGES"

	USER_DM_TYPING = "USER_DM_TYPING"

	ADD_FRIEND_REQUEST    = "ADD_FRIEND_REQUEST"
	REMOVE_FRIEND_REQUEST = "REMOVE_FRIEND_REQUEST"

	ADD_FRIEND_INCOMING_REQUEST    = "ADD_FRIEND_INCOMING_REQUEST"
	REMOVE_FRIEND_INCOMING_REQUEST = "REMOVE_FRIEND_INCOMING_REQUEST"

	ADD_USER_FRIENDLIST    = "ADD_USER_FRIENDLIST"
	REMOVE_USER_FRIENDLIST = "DELETE_USER_FRIENDLIST"

	ADD_USER_BLOCKEDLIST    = "ADD_USER_BLOCKEDLIST"
	REMOVE_USER_BLOCKEDLIST = "REMOVE_USER_BLOCKEDLIST"

	CLEAR_USER_MESSAGES  = "CLEAR_USER_MESSAGES"
	CLEAR_GUILD_MESSAGES = "CLEAR_GUILD_MESSAGES"
	USER_TYPING          = "USER_TYPING"

	ADD_USER_GUILDLIST    = "ADD_USER_GUILDLIST"
	REMOVE_USER_GUILDLIST = "REMOVE_USER_GUILDLIST"

	ADD_USER_BANLIST    = "ADD_USER_BANLIST"
	REMOVE_USER_BANLIST = "REMOVE_USER_BANLIST"

	ADD_USER_GUILDADMIN    = "ADD_USER_GUILDADMIN"
	REMOVE_USER_GUILDADMIN = "REMOVE_USER_GUILDADMIN"

	LOG_OUT = "LOG_OUT"

	UPDATE_USER_INFO      = "UPDATE_USER_INFO"
	UPDATE_SELF_USER_INFO = "UPDATE_SELF_USER_INFO"
)
