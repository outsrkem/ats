package audit

// ReqCreateAudLogRaw 上传日志请求体
type ReqCreateAudLogRaw struct {
	Service string `json:"service"` // 服务名称
	// 事件主要内容
	Events []struct {
		UserID     string   `json:"user_id"`     // 账号ID
		Account    string   `json:"account"`     // 操作账号
		SourceIP   string   `json:"source_ip"`   // 请求的客户端IP
		ResourceId []string `json:"resource_id"` // 资源ID
		Name       string   `json:"name"`        // 事件名称
		Rating     string   `json:"rating"`      // 日志类型/级别
		Etime      int64    `json:"etime"`       // 事件发生时间
		Message    string   `json:"message"`     // 日志消息
		Reqdata    string   `json:"reqdata"`     // 请求体，Get请求一般没有请求体
		Uagent     string   `json:"uagent"`      // user-agent
		Method     string   `json:"method"`      // 请求方法，GET/POST/DELETE/...
		ReqUrl     string   `json:"requrl"`      // 请求的URl路径
	} `json:"events"`
}

// ResTracesAuditLog 返回的日志列表
type ResTracesAuditLog struct {
	Eid        string `json:"eid"`
	UserID     string `json:"user_id"`
	Account    string `json:"account"`
	Service    string `json:"service"`
	ResourceId string `json:"resource_id"`
	Name       string `json:"name"`
	Rating     string `json:"rating"`
	ETime      int64  `json:"etime"`
	Message    string `json:"message"`
	Extras     string `json:"extras"`
}
