package json_rpc_request

import (
	"errors"
	"fmt"

	"github.com/0xYeah/project_template_go/api/api_json_rpc/json_rpc_protocol"
	"github.com/goccy/go-json"
)

func ParserRequest(body []byte, reqModel *json_rpc_protocol.RPCRequest) error {
	var tmpMap map[string]interface{}
	err := json.Unmarshal(body, &tmpMap)
	if err != nil {
		return err
	}

	if method, ok := tmpMap["method"].(string); ok {
		reqModel.Method = method
	} else {
		return errors.New("invalid or missing 'method' field")
	}

	reqModel.Params = tmpMap["params"]

	if jsonrpc, ok := tmpMap["jsonrpc"].(string); ok {
		reqModel.JsonRPC = jsonrpc
	} else {
		return errors.New("invalid or missing 'jsonrpc' field")
	}

	if id, ok := tmpMap["id"]; ok {
		reqModel.ID = fmt.Sprintf("%v", id)
	} else {
		reqModel.ID = "" // 设置默认ID为空字符串
	}

	return nil
}
