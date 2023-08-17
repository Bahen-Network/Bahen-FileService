package v1

import (
	"bytes"
	"file-service/config"
	"file-service/objectstorage"
	"file-service/util"
	"github.com/gin-gonic/gin"
	"io"
	"log"
)

type Controller struct {
	client *objectstorage.BNBClient
}

func NewController(client *objectstorage.BNBClient) *Controller {
	return &Controller{client: client}
}

func (c *Controller) UploadFolder(ctx *gin.Context) {
	log.Printf("1-1. UploadFolder start!")

	objectName := ctx.PostForm("objectName")
	bucketName := ctx.PostForm("bucketName")
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
	defer fileBytes.Close()

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, fileBytes); err != nil {
		ctx.JSON(500, gin.H{"error": "Reading uploaded file failed."})
		return
	}

	encryptedBytes, err := util.Encrypt(buffer.Bytes(), config.PrivateAESKey)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Encryption failed."})
		return
	}

	_, err = c.client.CreateObject(ctx.Request.Context(), bucketName, objectName, encryptedBytes)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to upload to BNBClient."})
		return
	}

	log.Printf("1-2. UploadFolder success!")
	ctx.JSON(200, gin.H{
		"msg": "folder uploaded",
	})
}

func (c *Controller) DownloadFolder(ctx *gin.Context) {
	objectName := ctx.PostForm("objectName")
	bucketName := ctx.PostForm("bucketName")
	if objectName == "" {
		ctx.JSON(400, gin.H{"error": "objectName is required in form data"})
		return
	}

	zipBytes, err := c.client.GetObject(ctx.Request.Context(), bucketName, objectName)
	if err != nil {
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

func (c *Controller) HelloWorld(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"msg": "Hello world!",
	})
}
