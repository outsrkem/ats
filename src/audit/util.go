package audit

import (
	"ats/src/config"
	"ats/src/pkg/common"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// DoHttpV1 向后端接口发送http请求
func DoHttpV1(res *httpCli) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(res.Method, res.Url, res.Body)
	if err != nil {
		hlog.Error("Get error: ", err)
		return nil, err
	}
	// 设置请求头
	for key, header := range res.Headers {
		req.Header.Set(key, header)
	}
	resp, err := client.Do(req)
	if err != nil {
		hlog.Error("Failed to send request:", err)
		return nil, err
	}
	return resp, nil
}

// CheckAction 向UIAS请求验证权限验证
func CheckAction(action string, token string) (bool, *map[string]map[string]interface{}) {
	var req httpCli
	req.Headers = map[string]string{"X-Auth-Token": token}

	req.Method = "POST"
	req.Url = config.Cfg.Ats.Uias.Endpoint + common.CheckAction
	raw, _ := json.Marshal(map[string]map[string]string{
		"uias": {
			"action": action,
		},
	})
	req.Body = bytes.NewBuffer(raw)
	resp, err := DoHttpV1(&req)
	if err != nil {
		hlog.Error("Error making GET request:", err)
		return false, nil
	}

	// 关闭请求
	defer func() {
		if resp != nil {
			if err := resp.Body.Close(); err != nil {
				hlog.Error("Close request failed: %v", err)
			}
		}
	}()

	var result map[string]map[string]interface{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		hlog.Error("read body error: ", err)
		return false, &result
	}

	hlog.Debug("resp body: ", string(body))
	if err := json.Unmarshal(body, &result); err != nil {
		hlog.Warn("json Unmarshal err", err)
		return false, &result
	}

	if resp.StatusCode == 404 {
		hlog.Warn("Interface not found: ", string(body))
		return false, &result
	}

	if resp.StatusCode != 200 {
		hlog.Warn("check action Status: ", resp.Status)
		return false, &result
	}
	authentication := result["payload"]["authentication"]
	if authentication != float64(1) {
		return false, &result
	}
	return true, &result
}

func strToInt64(str string) (int64, error) {
	hlog.Info("strToInt64, str: ", str)
	intValue, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		hlog.Error("Error converting string to int64: ", err)
		return 0, err
	}
	return intValue, nil
}

func strToInt(str string) (int, error) {
	hlog.Info("strToInt, str: ", str)
	intValue, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		hlog.Error("Error converting string to int64: ", err)
		return 0, err
	}
	return int(intValue), nil
}
