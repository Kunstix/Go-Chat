package config

import (
	"context"
)

type Configuration struct {
	PORT       string
	SECRET     string
	DBUSER     string
	DBPASSWORD string
	REDIS_HOST string
}

const GeneralChannel = "general"

var Ctx = context.Background()
