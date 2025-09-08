package mime

// MIME represent the valid MIME types that are supported
type MIME string

const (
	JSON       MIME = "application/json"
	XML        MIME = "application/xml"
	URLFORM    MIME = "application/x-www-form-urlencoded"
	MULTIFORM  MIME = "multipart/form-data"
	PLAINTEXT  MIME = "text/plain"
	HTML       MIME = "text/html"
	JAVASCRIPT MIME = "application/javascript"
	YAML       MIME = "application/yaml"
	PDF        MIME = "application/pdf"
	CSV        MIME = "text/csv"
	BINARY     MIME = "application/octet-stream"
)
