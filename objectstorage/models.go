// objectstorage/models.go

package objectstorage

// ObjectListResponse represents the response structure for listing objects.
type ObjectListResponse struct {
	Objects []ObjectInfo `json:"objects"`
}

// ObjectInfo provides detailed information for a single object.
type ObjectInfo struct {
	ObjectName string `json:"objectName"`
	Data       []byte `json:"data"`
	Type       string
}
