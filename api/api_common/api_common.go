// Package api_common api/api_common/api_common.go
package api_common

import (
	"github.com/0xYeah/project_template_go/api/api_config"
	"github.com/0xYeah/project_template_go/api/api_json_rpc/json_rpc_permissions"
)

// CheckAllowedMethods 允许方法
func CheckAllowedMethods(method string, permission *json_rpc_permissions.Permission) bool {
	p := api_config.CurrentApiConfig.Permissions[permission.PermissionTag]

	for _, aMethod := range p.AllowedMethods {
		if aMethod == method {
			return true
		}
	}
	return false
}

// CheckAllowedUserAgent 检查UA是否在白名单
func CheckAllowedUserAgent(uaName string, permission *json_rpc_permissions.Permission) bool {
	p := api_config.CurrentApiConfig.Permissions[permission.PermissionTag]

	for _, aUA := range p.AllowedUAs {
		if aUA == uaName {
			return true
		}
	}
	return false
}
