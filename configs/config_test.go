package configs

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	assert.Equal(t, nil, LoadConfig("./"))
	assert.Equal(t, "host", viper.GetString("DB.Host"))
	assert.Equal(t, 5432, viper.GetInt("DB.Port"))
	assert.Equal(t, "user", viper.GetString("DB.User"))
	assert.Equal(t, "password", viper.GetString("DB.Password"))
	assert.Equal(t, "dbname", viper.GetString("DB.DBName"))
}
