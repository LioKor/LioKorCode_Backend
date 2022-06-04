package constants

import "time"

const (
	Host                = "127.0.0.1:9091"
	RedisAddr           = "127.0.0.1:6379"
	IdKey               = "id"
	UserIdKey           = "uid"
	CodeKey             = "code"
	PageKey             = "page"
	CountKey            = "count"
	TaskId              = "taskId"
	SolutionId          = "solutionId"
	TasksPerPage        = 100
	WeekSec             = 604800
	DBConnect           = " dbname=liokoredu host=localhost port=5432 sslmode=disable pool_max_conns=10"
	CookieLength        = uint8(32)
	SessionCookieName   = "SID"
	SaltLength          = 8
	PythonAddress       = "http://178.62.57.180/check_solution?api_key=secret_key_here"
	SolutionsDir        = "/store/"
	AvatartDir          = "/media/avatars/"
	AvatartSalt         = 8
	PrivateLength       = 10
	Localhost           = "127.0.0.1"
	RedactorServicePort = ":3001"
	WSLength            = 16
	MaxSizeKB           = 8184
	SignKey             = "liokoredu"

	// Time allowed to read the next pong message from the peer.
	PongWait = 10 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * 9) / 10
	// Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second
)

var LetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
