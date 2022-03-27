package constants

const (
	IdKey             = "id"
	CodeKey           = "code"
	PageKey           = "page"
	TaskId            = "taskId"
	SolutionId        = "solutionId"
	TasksPerPage      = 10
	WeekSec           = 604800
	DBConnect         = " dbname=liokoredu host=localhost port=5432 sslmode=disable pool_max_conns=10"
	CookieLength      = uint8(32)
	SessionCookieName = "SID"
	SaltLength        = 8
	PythonAddress     = "http://10.106.0.2/check_task/long"
	SolutionsDir      = "/store/"
	PrivateLength     = 10
)

var LetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
