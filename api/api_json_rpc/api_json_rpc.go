package api_json_rpc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/0xYeah/project_template_go/api/api_config"
	"github.com/0xYeah/project_template_go/api/api_json_rpc/json_rpc_handler"
	"github.com/0xYeah/project_template_go/api/api_json_rpc/json_rpc_permissions"
	"github.com/0xYeah/project_template_go/config"
	"github.com/george012/gtbox/gtbox_log"
	"github.com/gorilla/mux"
)

var httpServer *http.Server

func extractBetweenSlashes(input string) string {
	// 找到第一个 `/`
	firstSlashIndex := strings.Index(input, "/")
	if firstSlashIndex == -1 {
		// 如果连第一个 `/` 都不存在，返回空字符串
		return ""
	}

	// 找到第一个 `/` 之后的子串
	remainder := input[firstSlashIndex+1:]

	// 找到第二个 `/`
	secondSlashIndex := strings.Index(remainder, "/")
	if secondSlashIndex == -1 {
		// 如果没有第二个 `/`，返回剩余部分
		return remainder
	}

	// 如果有第二个 `/`，取第一个 `/` 和第二个 `/` 之间的部分
	return remainder[:secondSlashIndex]
}

func StartAPIServiceWithJsonRPC(apiCfg *api_config.ApiConfig) {

	if apiCfg.Port < 1 || apiCfg.Port > 65535 {
		gtbox_log.LogErrorf("api port must be between 1 and 65535")
		return
	}

	api_config.CurrentApiConfig = apiCfg
	if apiCfg.Permissions != nil {
		for _, permissions := range apiCfg.Permissions {
			permissions.AllowedMethods = append(permissions.AllowedMethods, api_config.ApiCommonMethods...)
			permissions.AllowedUAs = append(permissions.AllowedUAs, config.ProjectName)
		}
	} else {
		apiCfg.Permissions = map[string]*json_rpc_permissions.Permission{
			json_rpc_permissions.APIPermissionTagDefault: &json_rpc_permissions.Permission{
				EncryptionEnabled: false,
				PermissionTag:     json_rpc_permissions.APIPermissionTagDefault,
				AllowedUAs:        []string{config.ProjectName},
				AllowedMethods:    api_config.ApiCommonMethods,
			},
		}
	}

	muxRouter := mux.NewRouter()
	//muxRouter.Use(json_rpc_handler.Middleware) // 使用中间件
	muxRouter.HandleFunc("/", json_rpc_handler.HomeHandler).Methods("GET")
	muxRouter.HandleFunc("/robots.txt", json_rpc_handler.RobotsHandler)

	for _, aPermission := range apiCfg.Permissions {
		for _, aapi := range apiCfg.Apis {
			aPath := ""
			if aapi.ApiPath != "" {
				aPath = fmt.Sprintf("/%s", aapi.ApiPath)
			}

			if aPermission.PermissionTag != "" && aPermission.PermissionTag != json_rpc_permissions.APIPermissionTagDefault {
				aPath = fmt.Sprintf("%s/%s", aPath, aPermission.PermissionTag)
			}

			muxRouter.HandleFunc(aPath, func(writer http.ResponseWriter, request *http.Request) {
				json_rpc_handler.ApiHandler(writer, request, aPermission)
			}).Methods("POST")
		}

	}

	addr := fmt.Sprintf("%s:%d", "0.0.0.0", apiCfg.Port)
	httpServer = &http.Server{
		Addr:    addr,
		Handler: muxRouter,
	}

	go func() {
		gtbox_log.LogInfof("API server Run On  [%s]", fmt.Sprintf("http://127.0.0.1:%d", apiCfg.Port))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			gtbox_log.LogErrorf("Failed to start HTTP server: %v\n", err)
		}
	}()

}

func StopApiServiceWithJsonRPC() {
	if httpServer == nil {
		gtbox_log.LogInfof("API server is not running")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gtbox_log.LogInfof("Shutting down API server...")
	if err := httpServer.Shutdown(ctx); err != nil {
		gtbox_log.LogErrorf("Error shutting down API server: %v\n", err)
	} else {
		gtbox_log.LogInfof("API server stopped successfully")
	}
}
