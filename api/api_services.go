// Package api api/api_services.go
package api

import (
	"github.com/0xYeah/project_template_go/api/api_config"
	"github.com/0xYeah/project_template_go/api/api_json_rpc"
)

func StartAPIServices(apiCfg *api_config.ApiConfig) {
	api_json_rpc.StartAPIServiceWithJsonRPC(apiCfg)
}

func StopApiServices() {
	api_json_rpc.StopApiServiceWithJsonRPC()
}
