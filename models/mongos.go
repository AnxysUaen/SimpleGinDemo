package models

// 接收的Log参数
type RequestData struct {
	Time       string `form:"time"`
	Name       string `form:"name"`
	Type       string `form:"type"`
	JSONData   string `form:"data"`
	RequestID  string `form:"id"`
	DataSource string `form:"source"`
}
