package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// 请求结构体
type Request struct {
	URL    string
	Params url.Values
	Data   interface{}
}

// 递归处理嵌套对象
func m(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, val := range v {
			if val != nil {
				result[k] = m(val)
			}
		}
		return result
	case []interface{}:
		return v
	default:
		return v
	}
}

// URL路径处理
func processURL(rawURL string) string {
	u, _ := url.Parse(rawURL)
	path := u.Path
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		base := strings.Join(parts[:2], "/")
		re := regexp.MustCompile("^" + regexp.QuoteMeta(base))
		return re.ReplaceAllString(path, "")
	}
	return path
}

// 生成签名
func GenerateSign(e *Request, t string) (string) {
	processedURL := processURL(e.URL)

	// 解析URL参数
	urlParams := url.Values{}
	if idx := strings.Index(e.URL, "?"); idx != -1 {
		queryPart := e.URL[idx+1:]
		urlParams, _ = url.ParseQuery(queryPart)
	}

	// 合并参数
	mergedParams := urlParams
	if e.Params != nil {
		mergedParams = mergeParams(urlParams, e.Params)
	}

	// 生成查询字符串
	queryStr := mergedParams.Encode()

	// 处理data
	processedData := m(e.Data)
	dataStr := toJSON(processedData)
	if dataStr == "{}" || dataStr == "null" {
		dataStr = ""
	}

	// 拼接签名字符串
	str := fmt.Sprintf("SwadSign%s%s%s%s", 
		processedURL, 
		queryStr, 
		dataStr, 
		t)

	// fmt.Println("String to be hashed:", str)
	// 计算MD5
	hash := md5.Sum([]byte(str))
	sign := hex.EncodeToString(hash[:])

	return sign
}

// 合并查询参数
func mergeParams(a, b url.Values) url.Values {
	result := url.Values{}
	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		result[k] = v
	}
	return result
}

// JSON序列化
func toJSON(data interface{}) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}