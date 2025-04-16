package audit

import "bytes"

type httpCli struct {
	Headers map[string]string
	Method  string
	Url     string
	Body    *bytes.Buffer
}

//	{
//	   "service": "xxxxServic",
//	   "events": [
//	       {
//	           "user_id": "xxx",
//	           "account": "admin",
//	           "source_ip": "222.10.10.111",
//	           "resource_id": [
//	               "3f507f27f2c54f609c60ab9dd6ae2a58"
//	           ],
//	           "name": "deleuser",
//	           "rating": "normal",
//	           "etime": 1732182636954,
//	           "message": "{[1,2,3][]}",
//	           "reqdata": "aaassssaa",
//	           "uagent": "ccc",
//	           "method": "method",
//	           "requrl": "requrl"
//	       }
//	   ]
//	}
//
// ReqCreateAudLogRaw 上传日志请求体
type ReqCreateAudLogRaw struct {
	Service string `json:"service"`
	Events  []struct {
		UserID     string   `json:"user_id"`
		Account    string   `json:"account"`
		SourceIP   string   `json:"source_ip"`
		ResourceId []string `json:"resource_id"`
		Name       string   `json:"name"`
		Rating     string   `json:"rating"`
		Etime      int64    `json:"etime"`
		Message    string   `json:"message"`
		Reqdata    string   `json:"reqdata"`
		Uagent     string   `json:"uagent"` // user-agent
		Method     string   `json:"method"`
		ReqUrl     string   `json:"requrl"`
	} `json:"events"`
}

// ResTracesAuditLog 返回的日志列表
type ResTracesAuditLog struct {
	Eid        *string `json:"eid"`
	UserID     *string `json:"user_id"`
	Account    *string `json:"account"`
	Service    *string `json:"service"`
	ResourceId *string `json:"resource_id"`
	Name       *string `json:"name"`
	Rating     *string `json:"rating"`
	ETime      *int64  `json:"etime"`
	Message    *string `json:"message"`
	Extras     *string `json:"extras"`
	CreateTime *int64  `json:"create_time"`
}
