package file

const (
	// file status
	StatusDfsCreated     = 0
	StatusDfsUploading   = 1
	StatusDfsCompleted   = 2
	StatusDfsDelivering  = 3
	StatusDfsDelete      = 4
	StatusDfsConcatErr   = 5
	StatusDfsCorrupt     = 6
	StatusCdnAssignJob   = 7
	StatusCdnRetrieved   = 8
	StatusCdnCorrupt     = 9
	StatusCdnCompleted   = 10
	StatusCdnMarkDeleted = 11
	StatusWaitDelete     = 12

	// roles
	RoleUnknown      uint8 = 0
	RoleMajesty      uint8 = 1
	RoleFileUpload   uint8 = 2
	RoleFileChecker  uint8 = 3
	RoleCdnPublisher uint8 = 4
	RoleCloudOps     uint8 = 5
	RoleAbyss        uint8 = 6

	// OpCode
	OpCodePublish    uint8 = 0
	OpCodeDelete     uint8 = 1
	OpCodeMarkDelete uint8 = 2
	OpCodeUpdate     uint8 = 3

	// Majesty events
	EventInit           uint8 = 1
	EventDeleteFile     uint8 = 2
	EventMarkDeleteFile uint8 = 3
	EventUpdateFile     uint8 = 4
	EventSaved          uint8 = 16 // File Upload and Abyss
	// CDN Publisher
	EventCdnSentRequest uint8 = 20
	EventCdnReady       uint8 = 21
	EventCdnFailed      uint8 = 22
	EventCdnDeleted     uint8 = 23
	EventCdnMarkDeleted uint8 = 24
	EventCdnUpdated     uint8 = 25
	// file checker
	EventFileScanPass      uint8 = 40
	EventFileScanInvalid   uint8 = 41
	EventFileScanSpecious  uint8 = 42
	EventFileScanError     uint8 = 43
	EventImageScanPass     uint8 = 44
	EventImageScanInvalid  uint8 = 45
	EventImageScanSpecious uint8 = 46
	EventImageScanError    uint8 = 47
	EventTextScanPass      uint8 = 48
	EventTextScanInvalid   uint8 = 49
	EventTextScanSpecious  uint8 = 50
	EventTextScanError     uint8 = 51

	LogAppID     = "AppID"
	LogRequestID = "RequestID"
	LogFileID    = "FileID"
	LogVendor    = "Vendor"
	// Log fields: Majesty
	LogOptional    = "Op"
	LogRoleID      = "RoleID"
	LogEventID     = "EventID"
	LogEventMsg    = "EventMsg"
	LogCallbackMsg = "CallbackMsg"

	// Log fields: CDN Publisher
	LogCdnURL        = "CdnURL"
	LogSourceURL     = "SourceURL"
	LogScanURL       = "ScanURL"
	LogPublishAction = "Act"
	LogItemID        = "ItemID"
	LogElapsedTime   = "ElapsedTime" // seconds

	// CDN vendors
	VendorGosunName   = "gosun"
	VendorGosunDomain = "newmarket1.oo523.com"
	VendorQiniuName   = "qiniu"
	VendorQiniuDomain = "droibaascdn.com"

	// CDNRootPath - the unified root path of CDN URL
	CDNRootPath = "/droi"

	// Slack Integration
	SlackChannelGeneral = "file_service"
	SlackChannelAlert   = "file_service_alert"
	SlackChannelTest    = "file_service_test"
)

// See http://mkdocs/file-majesty
var (
	// Events - the mapping between eventID and name
	Events = map[uint8]string{
		// Global
		EventInit:              "EventInit",
		EventDeleteFile:        "EventDeleteFile",
		EventUpdateFile:        "EventUpdateFile",
		EventSaved:             "EventSaved",
		EventCdnSentRequest:    "EventCdnSentRequest",
		EventCdnReady:          "EventCdnReady",
		EventCdnFailed:         "EventCdnFailed",
		EventCdnDeleted:        "EventCdnDeleted",
		EventCdnMarkDeleted:    "EventCdnMarkDeleted",
		EventCdnUpdated:        "EventCdnUpdated",
		EventFileScanPass:      "EventFileScanPass",
		EventFileScanInvalid:   "EventFileScanInvalid",
		EventFileScanSpecious:  "EventFileScanSpecious",
		EventFileScanError:     "EventFileScanError",
		EventImageScanPass:     "EventImageScanPass",
		EventImageScanInvalid:  "EventImageScanInvalid",
		EventImageScanSpecious: "EventImageScanSpecious",
		EventImageScanError:    "EventImageScanError",
		EventTextScanPass:      "EventTextScanPass",
		EventTextScanInvalid:   "EventTextScanInvalid",
		EventTextScanSpecious:  "EventTextScanSpecious",
		EventTextScanError:     "EventTextScanError",
	}
	// Roles - the mapping between roleID and name
	Roles = map[uint8]string{
		RoleUnknown:      "Unknown",
		RoleMajesty:      "Majesty",
		RoleFileUpload:   "FileUpload",
		RoleFileChecker:  "FileChecker",
		RoleCdnPublisher: "CdnPublisher",
		RoleCloudOps:     "CloudOps",
		RoleAbyss:        "Abyss",
	}

	// OpCodes - the mapping between operation code and name
	OpCodes = map[uint8]string{
		OpCodePublish:    "publish",
		OpCodeDelete:     "delete",
		OpCodeMarkDelete: "mark-delete",
		OpCodeUpdate:     "update",
	}
)
