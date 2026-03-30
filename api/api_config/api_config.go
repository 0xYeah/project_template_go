// Package api_config api/api_config/api_config.go
package api_config

import (
	"github.com/0xYeah/project_template_go/api/api_json_rpc/json_rpc_permissions"
	"github.com/george012/gtbox"
)

type ApiProxy struct {
	ApiPath     string `json:"api_path"`
	Enabled     bool   `toml:"enabled"`
	Address     string `json:"address"`
	AuthEnabled bool   `toml:"auth_enabled"`
	User        string `json:"user"`
	Pwd         string `json:"pwd"`
}

type ApiConfig struct {
	Enabled       bool                                        `json:"enabled"`
	Port          int                                         `json:"port"`
	Apis          []*ApiProxy                                 `json:"apis"`
	ClientTimeout string                                      `json:"client_timeout"`
	Permissions   map[string]*json_rpc_permissions.Permission `json:"permissions"`
}

var (
	CurrentApiConfig *ApiConfig
	CurrentRunMode   gtbox.RunMode
	ApiCommonMethods = []string{"auth", "logout", "test"}
)
