package config

import "os"

var (
	CheckerUrl = os.Getenv("CHECKER_URL")
)
