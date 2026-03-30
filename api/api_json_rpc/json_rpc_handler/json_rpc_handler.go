package json_rpc_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/0xYeah/project_template_go/api/api_json_rpc/json_rpc_config"
	"github.com/0xYeah/project_template_go/api/api_json_rpc/json_rpc_permissions"
	"github.com/0xYeah/project_template_go/api/api_json_rpc/json_rpc_protocol"
	"github.com/0xYeah/project_template_go/api/api_json_rpc/json_rpc_request"
	"github.com/0xYeah/project_template_go/api/api_json_rpc/json_rpc_response"
	"github.com/george012/gtbox/gtbox_encryption"
)

func encResult(respModel *json_rpc_protocol.RPCResponse, encKey string) error {
	if respModel.Result == nil {
		return nil
	}

	// TODO 如果是开启加密 则对result 进行全字段加密成字符串
	if resultStr, ok := respModel.Result.(string); ok {
		// 尝试直接 加密
		encStr := gtbox_encryption.GTEnc(resultStr, encKey)
		respModel.Result = encStr
	} else {
		needEncrypt := false
		// 检查是否是空字典或空数组
		switch params := respModel.Result.(type) {
		case map[string]interface{}:
			if len(params) > 0 {
				needEncrypt = true
			}
		case []interface{}:
			if len(params) > 0 {
				needEncrypt = true
			}
		}
		if needEncrypt == true {
			jBytes, _ := json.Marshal(respModel.Result)
			jStr := string(jBytes)
			encStr := gtbox_encryption.GTEnc(jStr, encKey)
			respModel.Result = encStr
		}
	}
	return nil
}

func decParams(reqModel *json_rpc_protocol.RPCRequest, decKey string) error {
	// TODO 如果是开启加密 则对params 进行全字段字符串解密操作
	if reqModel.Params == nil {
		return nil
	}
	if paramsStr, ok := reqModel.Params.(string); ok {
		// 尝试直接解密
		decStr := gtbox_encryption.GTDec(paramsStr, decKey)
		var jObj interface{}
		err := json.Unmarshal([]byte(decStr), &jObj)
		if err != nil {
			return err
		}
		reqModel.Params = jObj
	}
	return nil
}

// ApiHandler 处理 HTTP 请求并转发给 TCP 服务器
func ApiHandler(w http.ResponseWriter, r *http.Request, permission *json_rpc_permissions.Permission) {

	reqModel := &json_rpc_protocol.RPCRequest{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		bodyErr := errors.New("body read error")
		json_rpc_response.HandleResponse(w, bodyErr, nil, reqModel.ID)
		return
	}

	err = json_rpc_request.ParserRequest(body, reqModel)
	if err != nil {
		json_rpc_response.HandleResponse(w, err, nil, reqModel.ID)
		return
	}

	uaOK := checkUA(r, permission)
	if !uaOK {
		uaOKErr := errors.New("UA not allowed")
		json_rpc_response.HandleResponse(w, uaOKErr, nil, reqModel.ID)
		return
	}

	err = checkMethods(r, reqModel, permission)
	if err != nil {
		json_rpc_response.HandleResponse(w, err, nil, reqModel.ID)
		return
	}

	decKey := fmt.Sprintf("%s/%s", r.UserAgent(), reqModel.ID)
	if permission.EncryptionEnabled == true {
		err = decParams(reqModel, decKey)
		if err != nil {
			json_rpc_response.HandleResponse(w, err, nil, reqModel.ID)
			return
		}
	}

	var respError error
	respModel := &json_rpc_protocol.RPCResponse{
		Result:  nil,
		Error:   nil,
		JsonRPC: json_rpc_config.JSONRPCVersion,
		ID:      reqModel.ID,
	}

	switch reqModel.Method {
	case "test":
		respModel.Result = "this is test method, request is success"
	default:
		// TODO api action handler

		// Examples ::

		//switch permission.PermissionTag {
		//case "test_a":
		//	// TODO TestA 相关方法处理
		//	respError = handlePowApis(reqModel, permission, apiProxy, respModel)
		//case "test_b":
		//	// TODO TestB 相关方法处理
		//	respError = handleTXApis(reqModel, permission, apiProxy, respModel)
		//default:
		//	respError = errors.New(fmt.Sprintf("%s method support is error", reqModel.Method))
		//}
	}

	if permission.EncryptionEnabled == true {
		err = encResult(respModel, decKey)
		if err != nil {
			json_rpc_response.HandleResponse(w, err, nil, reqModel.ID)
			return
		}
	}
	json_rpc_response.HandleResponse(w, respError, respModel, reqModel.ID)

}
