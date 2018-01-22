package majesty

// LaunchPayload - the payload to launch majesty
type LaunchPayload struct {
	EventID uint8  `json:"EventId"`
	RoleID  uint8  `json:"RoleId"`
	Files   []File `json:"Files"`
}

// File - basic file meta
type File struct {
	Appid     string  `json:"AppId"`
	FID       string  `json:"Fid"`
	UID       string  `json:"Uid"`
	Path      string  `json:"Path"`
	Type      string  `json:"Type"`
	Size      uint64  `json:"Size"`
	MD5       string  `json:"MD5"`
	SourceURL *string `json:"sourceURL"`
}

// CallBackPayload - payload sent to Majesty
type CallBackPayload struct {
	Code    int                    `json:"Code"`
	EventID uint8                  `json:"EventId"`
	RoleID  uint8                  `json:"RoleId"`
	Appid   string                 `json:"AppId"`
	Fid     string                 `json:"Fid"`
	Msg     string                 `json:"Msg"`
	Op      map[string]interface{} `json:"Op"`
}
