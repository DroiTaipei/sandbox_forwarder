package cdnpublisher

// LaunchPayload - Define the JSON template for Majesty to trigger CDN Publusher
type LaunchPayload struct {
	RequestTime string     `json:"RequestTime"`
	Type        uint8      `json:"Type"`
	Files       []FileMeta `json:"Files"`
}

// FileMeta - the payload sent from Majesty
type FileMeta struct {
	Appid     string  `json:"AppID"`
	FID       string  `json:"FID"`
	Size      uint64  `json:"Size"`
	MD5       string  `json:"MD5"`
	RID       string  `json:"RID"`
	Path      string  `json:"Path"`
	SourceURL *string `json:"sourceURL"`
}
