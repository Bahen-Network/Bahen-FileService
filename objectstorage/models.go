// objectstorage/models.go

package objectstorage

type Folder struct {
	Objects []Object `json:"objects"`
}

type Object struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
}
