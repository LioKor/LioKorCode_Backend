package constants

const (
	IdKey               = "id"
	UserIdKey           = "uid"
	CodeKey             = "code"
	PageKey             = "page"
	TaskId              = "taskId"
	SolutionId          = "solutionId"
	TasksPerPage        = 10
	WeekSec             = 604800
	DBConnect           = " dbname=liokoredu host=localhost port=5432 sslmode=disable pool_max_conns=10"
	CookieLength        = uint8(32)
	SessionCookieName   = "SID"
	SaltLength          = 8
	PythonAddress       = "http://10.106.0.2/check_task/multiple_files"
	SolutionsDir        = "/store/"
	AvatartDir          = "/public/"
	AvatartSalt         = 8
	PrivateLength       = 10
	Localhost           = "127.0.0.1"
	RedactorServicePort = ":3001"
	WSLength            = 16
	MaxSizeKB           = 8184
)

var LetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
