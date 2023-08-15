package v1

import (
	"file-service/config"
	"file-service/objectstorage"
	"file-service/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"mime/multipart"
)

type Controller struct {
	client *objectstorage.BNBClient
}

func NewController(client *objectstorage.BNBClient) *Controller {

	c := &Controller{}
	c.client = client

	return c
}

type UploadFolderRequestBody struct {
	ObjectName string `json:"object_name"`
}

// TODO: optimize code to support large size file up/down load!
func (c *Controller) UploadFolder(ctx *gin.Context) {
	log.Printf("1-1. UploadFolder start!")

	objectName := ctx.PostForm("objectName")
	if objectName == "" {
		ctx.JSON(400, gin.H{"error": "objectName is required in form data"})
		return
	}

	file, err := ctx.FormFile("folder")
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error in retrieving the folder form file."})
		return
	}

	fileBytes, err := file.Open()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error in opening the uploaded file."})
		return
	}
	defer func(fileBytes multipart.File) {
		err := fileBytes.Close()
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Error in closing the uploaded file."})
		}
	}(fileBytes)

	allBytes, err := ioutil.ReadAll(fileBytes)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Reading uploaded file failed."})
		return
	}

	// Enciphered Data
	encryptedBytes, err := util.Encrypt(allBytes, config.PrivateAESKey)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Encryption failed."})
		return
	}

	_, err = c.client.CreateObject(ctx.Request.Context(), "testbucket", objectName, encryptedBytes)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to upload to BNBClient."})
		return
	}

	log.Printf("1-2. UploadFolder success!")
	ctx.JSON(200, gin.H{
		"msg": "folder uploaded",
	})
}

// Download Folder
func (c *Controller) DownloadFolder(ctx *gin.Context) {
	objectName := ctx.PostForm("objectName")

	zipBytes, err := c.client.GetObject(ctx.Request.Context(), "testbucket", objectName)
	if err != nil {
		util.HandleErr(err, "")
		ctx.JSON(500, gin.H{"error": "Failed to getObject from BNB to BNBClient."})
		return
	}

	decryptedBytes, err := util.Decrypt(zipBytes, config.PrivateAESKey)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Decryption failed."})
		return
	}
	ctx.Header("Content-Disposition", "attachment; filename=folder.zip")
	ctx.Data(200, "application/zip", decryptedBytes)
}
