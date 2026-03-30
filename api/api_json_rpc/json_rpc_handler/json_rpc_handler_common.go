package json_rpc_handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/0xYeah/project_template_go/api/api_common"
	"github.com/0xYeah/project_template_go/api/api_json_rpc/json_rpc_permissions"
	"github.com/0xYeah/project_template_go/api/api_json_rpc/json_rpc_protocol"
)

// HomeHandler 处理根路径请求
func HomeHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(""))
}

func RobotsHandler(w http.ResponseWriter, r *http.Request) {
	robotsTxt := `
User-agent: *
Disallow: /
`
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(robotsTxt))
}

// checkUA 检查UA
func checkUA(req *http.Request, permission *json_rpc_permissions.Permission) bool {

	// 检查 User-Agent
	ua := req.Header.Get("User-Agent")
	uas := strings.Split(ua, "/")
	if len(uas) > 1 {
		uaName := uas[0]
		// 如果 UA 校验通过，直接调用下一个处理程序
		if api_common.CheckAllowedUserAgent(uaName, permission) {
			return true
		}
	}
	return false
}

func checkMethods(r *http.Request, reqModel *json_rpc_protocol.RPCRequest, permission *json_rpc_permissions.Permission) error {

	if r.Method != http.MethodPost {
		return errors.New("Only POST requests are allowed")
	}

	//	检测 请求rpc 方法权限
	methodOK := api_common.CheckAllowedMethods(reqModel.Method, permission)
	if methodOK == false {
		return errors.New("request method is not allowed")
	}
	return nil
}
