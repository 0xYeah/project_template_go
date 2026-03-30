package json_rpc_protocol

// RPCRequest JSON-RPC 请求和响应结构
type RPCRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
	//ParamsDec interface{} `json:"-"` // 该字段不会参与 JSON 序列化
	JsonRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
}

type RPCResponse struct {
	Result  interface{} `json:"result"`
	Error   *RPCError   `json:"error"`
	JsonRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
}

type RPCError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
