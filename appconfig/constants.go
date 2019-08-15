package appconfig

import "fmt"

// RabbitMQURL ...
const RabbitMQURL = "159.65.220.217:5672"

// RabbitMQUsername ...
const RabbitMQUsername = "admin"

// RabbitMQPassword ...
const RabbitMQPassword = "Password123"

// RabbitMQRpcQueue ...
const RabbitMQRpcQueue = "rpc-die-queue"

// PublicAPIAddress ...
const PublicAPIAddress = ":9000"

// PublicAPIFullAddress ...
const PublicAPIFullAddress = "127.0.0.1:9000"

// PublicAPIUrl ...
const PublicAPIUrl = "http://localhost"

// PublicAPIWriteTimeout ...
const PublicAPIWriteTimeout = 15

// PublicAPIReadTimeout ...
const PublicAPIReadTimeout = 15

// PublicAPICountriesPattern ...
const PublicAPICountriesPattern = "/countries"

// PublicAPITeamsPattern ...
const PublicAPITeamsPattern = "/teams"

// PublicAPILeaguesPattern ...
const PublicAPILeaguesPattern = "/leagues"

// PublicAPIPlayerPattern ...
const PublicAPIPlayerPattern = "/player"

// DieAPIAddress ...
const DieAPIAddress = ":8000"

// DieAPIUrl ...
const DieAPIUrl = "http://localhost"

// DieAPICountriesPattern ...
const DieAPICountriesPattern = "/api//mysql/countries"

// DieAPITeamsPattern ...
const DieAPITeamsPattern = "/api/mysql/teams"

// DieAPILeaguesPattern ...
const DieAPILeaguesPattern = "/api/mysql/leagues"

// DieAPIPlayerPattern ...
const DieAPIPlayerPattern = "/api/mysql/player"

// GetPublicAPIBaseURL ...
func GetPublicAPIBaseURL() string {
	return fmt.Sprintf(
		"%s%s",
		DieAPIUrl,
		DieAPIAddress,
	)
}

// GetDieAPIBaseURL ...
func GetDieAPIBaseURL() string {
	return fmt.Sprintf(
		"%s%s",
		PublicAPIUrl,
		PublicAPIAddress,
	)
}

// GetAPIPatternMap ...
func GetAPIPatternMap() map[string]string {
	// MAP<PublicAPI><DieAPI>
	return map[string]string{
		PublicAPICountriesPattern: DieAPICountriesPattern,
		PublicAPITeamsPattern:     DieAPITeamsPattern,
		PublicAPILeaguesPattern:   DieAPILeaguesPattern,
		PublicAPIPlayerPattern:    DieAPIPlayerPattern,
	}
}

// GetPublicGetAPIPatterns ...
func GetPublicGetAPIPatterns() []string {
	m := GetAPIPatternMap()

	keys := []string{}
	for key := range m {
		keys = append(keys, key)
	}

	return keys
}
