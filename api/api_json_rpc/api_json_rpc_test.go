package api_json_rpc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/0xYeah/project_template_go/api/api_config"
	"github.com/0xYeah/project_template_go/config"
	"github.com/george012/gtbox/gtbox_encryption"
	"github.com/goccy/go-json"
)

type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	UA         string
}

// NewAPIClient 创建一个新的 API 客户端
func NewAPIClient(baseURL string, timeout time.Duration) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		UA:      fmt.Sprintf("%s/%s", config.ProjectName, config.ProjectVersion),
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// SendRequest 通用请求方法
func (c *APIClient) SendRequest(method, endpoint string, headers map[string]string, body interface{}) ([]byte, error) {
	// 构建完整的 URL
	url := c.BaseURL + "/" + endpoint

	// 将 body 转换为 JSON
	var requestBody []byte
	if body != nil {
		var err error
		requestBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest(method, url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置 Headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 默认设置 JSON Content-Type
	if _, ok := headers["Content-Type"]; !ok {
		req.Header.Set("Content-Type", "application/json")
	}

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {

		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}

	// 读取响应
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return responseData, nil
}

func TestMethods(t *testing.T) {
	testApiCfg := &api_config.ApiConfig{
		Enabled: true,
		Port:    12345,
		Apis: []*api_config.ApiProxy{
			&api_config.ApiProxy{
				ApiPath:     "test_rpc_01",
				Enabled:     false,
				Address:     "",
				AuthEnabled: false,
				User:        "",
				Pwd:         "",
			},
			&api_config.ApiProxy{
				ApiPath:     "test_rpc_02",
				Enabled:     true,
				Address:     "",
				AuthEnabled: false,
				User:        "",
				Pwd:         "",
			},
		},
		ClientTimeout: "3s",
	}

	StartAPIServiceWithJsonRPC(testApiCfg)
	time.Sleep(1 * time.Second)

	apiBaseUrl := "http://127.0.0.1:12345"

	client := NewAPIClient(apiBaseUrl, 10*time.Second)

	headers := map[string]string{"Authorization": "Bearer example-token"}
	headers["User-Agent"] = fmt.Sprintf("%s/%s", config.ProjectName, config.ProjectVersion)

	// 测试 GET 请求
	response, err := client.SendRequest("GET", "", headers, nil)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	t.Logf("GET Response:[%s]", string(response))

	// 测试 POST 请求

	postBody := map[string]interface{}{
		"method":  "test",
		"params":  []interface{}{"test_params_01", "test_params_02"},
		"jsonrpc": "2.0",
		"id":      "1",
	}
	for _, tApi := range testApiCfg.Apis {
		response, err = client.SendRequest("POST", tApi.ApiPath, headers, postBody)
		if err != nil {
			t.Fatalf("[%s] POST request failed: %v", tApi.ApiPath, err)
		}

		t.Logf("[%s] POST Response:[%s]", tApi.ApiPath, string(response))

		t.Logf("starting encryption post request\n")
		tpb, _ := json.Marshal(postBody)
		tpbStr := string(tpb)
		t_ua := fmt.Sprintf("%s/%s", config.ProjectName, config.ProjectVersion)
		t_key := fmt.Sprintf("%s/%s", t_ua, "1")

		encStr := gtbox_encryption.GTEnc(tpbStr, t_key)
		t.Logf("enc str [%s]\n", encStr)

		response, err = client.SendRequest("POST", fmt.Sprintf("%s/tx", tApi.ApiPath), headers, postBody)
		if err != nil {
			t.Fatalf("POST encryption request failed: %v", err)
		}
		t.Logf("POST encryption Response:\n%s", response)

		var wait_decData map[string]interface{}
		json.Unmarshal(response, &wait_decData)

		decStr := gtbox_encryption.GTDec(wait_decData["result"].(string), t_key)
		wait_decData["result"] = decStr
		deByte, _ := json.Marshal(&wait_decData)

		t.Logf("POST encryption Response with decryption\n%s\n", deByte)

		time.Sleep(1 * time.Second)
	}

	StopApiServiceWithJsonRPC()
}

func TestGenerateReqParamsString(t *testing.T) {
	testParams := []interface{}{
		"test_params_01",
		"test_params_02",
	}
	tpb, _ := json.Marshal(testParams)
	tpbStr := string(tpb)

	t_reqCode := strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
	fmt.Printf("reqCode [%s]\n", t_reqCode)

	t_ua := fmt.Sprintf("%s/%s", config.ProjectName, config.ProjectVersion)
	t_key := fmt.Sprintf("%s/%s", t_ua, t_reqCode)

	encStr := gtbox_encryption.GTEnc(tpbStr, t_key)
	fmt.Printf("enc str [%s]\n", encStr)

	decStr := gtbox_encryption.GTDec(encStr, t_key)
	fmt.Printf("dec str [%s]\n", decStr)
}
