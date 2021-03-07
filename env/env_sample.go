// +build sample

// 1. Copy this file
//   for local      : cp env_sample.go env_local.go
//   for production : cp env_sample.go env_production.go
//
// 2. Edit build tag
//   for local      : +build local
//   for production : +build production

package env

const (
	DiscordBotToken     = "Discord Bot Token"
	DiscordBotClientId  = "Client ID"
	DbDriver            = "mysql"
	DbUser              = "user"
	DbPassword          = "password"
	DbHost              = "127.0.0.1"
	DbPort              = "3306"
	DbName              = "database_name"
	DbLogLevel          = "info" // silent, error, warn, info
	SlackToken          = "TOKEN"
	SlackChannelWarning = "#channel-name-warning"
	SlackChannelAlert   = "#channel-name-alert"
	SlackTitleWarning   = "robotyosu-go warning notification"
	SlackTitleAlert     = "robotyosu-go alert notification"
	ScheduledDuration   = 60
)
