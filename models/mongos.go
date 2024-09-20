package models

// 接收的Log参数
type RequestData struct {
	RequestID  string `form:"id"`
	Name       string `form:"name" binding:"required"`
	Time       string `form:"time" binding:"required"`
	JSONData   string `form:"data" binding:"required"`
	DataSource string `form:"source" binding:"required"`
}
