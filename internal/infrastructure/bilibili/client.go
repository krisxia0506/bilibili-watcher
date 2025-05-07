package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Base URLs - Can be extended if needed
const (
	apiBaseURL = "https://api.bilibili.com"
)

// Client Bilibili API 客户端结构体。
// 负责维护 HTTP 客户端和通用请求逻辑。
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	sessData   string
}

// NewClient 创建一个新的 Bilibili API 客户端实例。
// sessData: 从配置中获取的 SESSDATA cookie 字符串。
func NewClient(sessData string) *Client {
	baseURL, _ := url.Parse(apiBaseURL) // Error ignored for constant URL
	return &Client{
		httpClient: &http.Client{},
		baseURL:    baseURL,
		sessData:   sessData,
	}
}

// Get 发送一个 GET 请求到指定的 API 路径，并将 JSON 响应解码到 target 中。
// path: 相对于 baseURL 的 API 路径 (例如 "/x/web-interface/view")。
// params: URL 查询参数。
// target: 用于解码 JSON 响应体的目标结构体指针。
func (c *Client) Get(path string, params url.Values, target interface{}) error {
	// 构建完整的请求 URL
	requestURL := c.baseURL.ResolveReference(&url.URL{Path: path})
	if params != nil {
		requestURL.RawQuery = params.Encode()
	}
	log.Printf("Request URL: %s", requestURL.String())
	// 创建 HTTP GET 请求
	req, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request for %s: %w", path, err)
	}

	// 设置通用请求头
	req.Header.Set("Accept", "application/json")
	if c.sessData != "" {
		req.Header.Set("Cookie", c.sessData)
	}
	// TODO: 可能需要添加 User-Agent 等其他通用 Header

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send GET request to %s: %w", path, err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body from %s: %w", path, err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		// 尝试提取 Bilibili 错误信息（如果响应是 JSON 格式的话）
		var baseResp struct { // 尝试解析通用错误结构
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		jsonErr := json.Unmarshal(body, &baseResp)
		if jsonErr == nil && baseResp.Message != "" {
			return fmt.Errorf("unexpected status code %d from %s. Bilibili error: code=%d, message=%s",
				resp.StatusCode, path, baseResp.Code, baseResp.Message)
		}
		// 如果无法解析为 Bilibili 错误，返回通用 HTTP 错误
		return fmt.Errorf("unexpected status code %d from %s, body: %s", resp.StatusCode, path, string(body))
	}

	// 解析 JSON 响应到目标结构体
	if target != nil {
		if err := json.Unmarshal(body, target); err != nil {
			// 提供更多上下文信息帮助调试
			bodyStr := string(body)
			if len(bodyStr) > 500 { // 避免打印过长的 body
				bodyStr = bodyStr[:500] + "..."
			}
			return fmt.Errorf("failed to unmarshal json response from %s into %T: %w, body snippet: %s",
				path, target, err, bodyStr)
		}
	}

	return nil
}
