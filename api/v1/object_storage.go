package v1

import (
	"archive/zip"
	"bytes"
	"file-service/config"
	"file-service/objectstorage"
	"file-service/util"
	"fmt"
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

func (c *Controller) PostObject(ctx *gin.Context) {
	log.Printf("1-1. UploadFolder start!")

	objectName := ctx.PostForm("objectName")
	bucketName := ctx.PostForm("bucketName")
	if objectName == "" {
		ctx.JSON(400, gin.H{"error": "objectName is required in form data"})
		return
	}

	file, err := ctx.FormFile("folder")
	if err != nil {
		util.HandleErr(err, "ctx.FormFile")
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

func (c *Controller) GetObject(ctx *gin.Context) {
	objectName := ctx.PostForm("objectName")
	bucketName := ctx.PostForm("bucketName")
	if objectName == "" {
		ctx.JSON(400, gin.H{"error": "objectName is required in form data"})
		return
	}

	if bucketName == "" {
		ctx.JSON(400, gin.H{"error": "bucketName is required in form data"})
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

func (c *Controller) ListObjects(ctx *gin.Context) {
	bucketName := ctx.PostForm("bucketName")
	if bucketName == "" {
		ctx.JSON(400, gin.H{"error": "bucketName is required in form data"})
		return
	}

	folder, err := c.client.ListObjects(ctx.Request.Context(), bucketName)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to getObject from BNB to BNBClient."})
		return
	}

	// Create a new zip containing the decrypted objects.
	var buffer bytes.Buffer
	zipWriter := zip.NewWriter(&buffer)
	for _, obj := range folder.Objects {
		decryptedByte, err := util.Decrypt(obj.Data, config.PrivateAESKey)
		if err != nil {
			util.HandleErr(err, "")
			ctx.JSON(500, gin.H{"error": fmt.Sprintf("Decryption failed for the object: %s.", obj.ObjectName)})
			return
		}
		zipFile, err := zipWriter.Create(obj.ObjectName + ".zip")
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to create new zip entry."})
			return
		}
		_, err = zipFile.Write(decryptedByte)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to write decrypted file to new zip."})
			return
		}
	}
	err = zipWriter.Close()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to finalize new zip file."})
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=aggregated.zip")
	ctx.Data(200, "application/zip", buffer.Bytes())
}

func (c *Controller) HelloWorld(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"msg": "Hello world!",
	})
}
