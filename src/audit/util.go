package audit

import (
	"ats/src/config"
	"ats/src/pkg/common"
	"ats/src/slog"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

// DoHttpV1 向后端接口发送http请求
func DoHttpV1(c *app.RequestContext, res *httpCli) (*http.Response, error) {
	klog := slog.FromContext(c)
	client := &http.Client{}
	req, err := http.NewRequest(res.Method, res.Url, res.Body)
	if err != nil {
		klog.Error("Get error: ", err)
		return nil, err
	}
	// 设置请求头
	for key, header := range res.Headers {
		req.Header.Set(key, header)
	}
	resp, err := client.Do(req)
	if err != nil {
		klog.Error("Failed to send request:", err)
		return nil, err
	}
	return resp, nil
}

// CheckAction 向UIAS请求验证权限验证
func CheckAction(c *app.RequestContext, action string, token string) (bool, *map[string]map[string]interface{}) {
	klog := slog.FromContext(c)
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
	resp, err := DoHttpV1(c, &req)
	if err != nil {
		klog.Error("Error making GET request:", err)
		return false, nil
	}

	// 关闭请求
	defer func() {
		if resp != nil {
			if err := resp.Body.Close(); err != nil {
				klog.Errorf("Close request failed: %s", err)
			}
		}
	}()

	var result map[string]map[string]interface{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		klog.Error("read body error: ", err)
		return false, &result
	}

	klog.Debug("resp body: ", string(body))
	if err := json.Unmarshal(body, &result); err != nil {
		klog.Error("json Unmarshal err", err)
		return false, &result
	}

	if resp.StatusCode == 404 {
		klog.Error("api not found: ", string(body))
		return false, &result
	}

	if resp.StatusCode != 200 {
		klog.Error("check action Status: ", resp.Status)
		return false, &result
	}
	authentication := result["payload"]["authentication"]
	if authentication != float64(1) {
		klog.Warnf("Permission denied.")
		return false, &result
	}
	return true, &result
}

func strToInt64(str string) (int64, error) {
	intValue, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

func strToInt(str string) (int, error) {
	intValue, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(intValue), nil
}
