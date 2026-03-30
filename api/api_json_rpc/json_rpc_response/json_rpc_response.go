package json_rpc_response

import (
	"net/http"

	"github.com/0xYeah/project_template_go/api/api_json_rpc/json_rpc_config"
	"github.com/0xYeah/project_template_go/api/api_json_rpc/json_rpc_protocol"
	"github.com/george012/gtbox/gtbox_log"
	"github.com/goccy/go-json"
)

func HandleResponse(w http.ResponseWriter, err error, resp *json_rpc_protocol.RPCResponse, reqCode string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	aID := "1"
	if reqCode != "" {
		aID = reqCode
	}

	if resp == nil {
		resp = &json_rpc_protocol.RPCResponse{
			JsonRPC: json_rpc_config.JSONRPCVersion,
			ID:      aID,
		}
	}

	if err != nil {
		resp.Error = &json_rpc_protocol.RPCError{
			Code:    -1,
			Message: err.Error(),
		}
	}

	// 直接使用 Encoder 编码并写入响应体
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		// 如果编码失败，可以记录日志或进一步处理
		gtbox_log.LogErrorf("Failed to encode JSON response[%v]", http.StatusInternalServerError)
	}
}
