package constants

// ServiceName The name of this module/service
const ServiceName = "users"

// GRPC Service Names
const (
	UserServiceName = "USERS"
)

// Dependency Injection Keys
const (
	DatabaseTransactionKey = "tx"
	UsersRepoKey           = "usersRepo"
	UsersCtrlKey           = "usersCtrl"
)

// Metric Names
const (
	UsersRegisteredCount = "users_registered_count"
)
