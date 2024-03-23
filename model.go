package logger

// LogTdrModel or Transaction Data Record
type LogTdrModel struct {
	AppName    string `json:"app"`
	AppVersion string `json:"ver"`
	ThreadID   string `json:"xid"`
	JourneyID  string `json:"jid"`
	ChainID    string `json:"cid"`

	Path         string `json:"path"`
	Method       string `json:"method"`
	IP           string `json:"ip"`
	Port         int    `json:"port"`
	SrcIP        string `json:"srcIP"`
	RespTime     int64  `json:"rt"`
	ResponseCode string `json:"rc"`

	Header   interface{} `json:"header"` // better to pass data here as is, don't cast it to string. use map or array
	Request  interface{} `json:"req"`
	Response interface{} `json:"resp"`
	Error    string      `json:"error"`

	AdditionalData interface{} `json:"addData"`
}
