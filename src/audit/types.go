package audit

import "bytes"

type httpCli struct {
	Headers map[string]string
	Method  string
	Url     string
	Body    *bytes.Buffer
}

type ReqCreateAudLogRaw struct {
	UserID   string `json:"user_id"`
	Account  string `json:"account"`
	SourceIP string `json:"source_ip"`
	Service  string `json:"service"`
	Name     string `json:"name"`
	Rating   string `json:"rating"`
	Etime    int64  `json:"etime"`
	Message  string `json:"message"`
}

type ResTracesAuditLog struct {
	Eid        *string `json:"eid"`
	UserID     *string `json:"user_id"`
	Account    *string `json:"account"`
	SourceIP   *string `json:"source_ip"`
	Service    *string `json:"service"`
	Name       *string `json:"name"`
	Rating     *string `json:"rating"`
	ETime      *int64  `json:"etime"`
	Message    *string `json:"message"`
	CreateTime *int64  `json:"create_time"`
}
