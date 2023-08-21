package v1

import (
	"archive/zip"
	"bytes"
	"file-service/config"
	"file-service/storageclient"
	"file-service/util"
	"io"
	"log"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	client *storageclient.BNBClient
}

func NewController(client *storageclient.BNBClient) *Controller {
	return &Controller{client: client}
}
func (c *Controller) PostObject(ctx *gin.Context) {
	log.Printf("1-1. UploadFolder start!")

	objectName := ctx.PostForm("objectName")
	bucketName := ctx.PostForm("bucketName")
	if objectName == "" {
		util.ReportError(ctx, nil, util.ObjectNameArgumentError)
		return
	}

	if bucketName == "" {
		util.ReportError(ctx, nil, util.BucketNameArgumentError)
		return
	}

	file, err := ctx.FormFile("folder")
	if err != nil {
		util.ReportError(ctx, err, util.FormDataFileRetrieveError)
		return
	}

	fileBytes, err := file.Open()
	if err != nil {
		util.ReportError(ctx, err, util.FileOpenError)
		return
	}

	defer func(fileBytes multipart.File) {
		err := fileBytes.Close()
		if err != nil {
			util.ReportError(ctx, err, util.FileOpenError)
		}
	}(fileBytes)

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, fileBytes); err != nil {
		util.ReportError(ctx, err, util.FileReadError)
		return
	}

	encryptedBytes, err := util.Encrypt(buffer.Bytes(), config.PrivateAESKey)
	if err != nil {
		util.ReportError(ctx, err, util.EncryptionError)
		return
	}

	_, err = c.client.CreateObject(ctx.Request.Context(), bucketName, objectName, encryptedBytes)
	if err != nil {
		util.ReportError(ctx, err, util.BNBClientUploadError)
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
		util.ReportError(ctx, nil, util.ObjectNameArgumentError)
		return
	}

	if bucketName == "" {
		util.ReportError(ctx, nil, util.BucketNameArgumentError)
		return
	}

	zipBytes, err := c.client.GetObject(ctx.Request.Context(), bucketName, objectName)
	if err != nil {
		util.ReportError(ctx, err, util.BNBClientUploadError)
		return
	}

	decryptedBytes, err := util.Decrypt(zipBytes, config.PrivateAESKey)
	if err != nil {
		util.ReportError(ctx, err, util.DecryptionError)
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=folder.zip")
	ctx.Data(200, "application/zip", decryptedBytes)
}

func (c *Controller) ListObjects(ctx *gin.Context) {
	bucketName := ctx.PostForm("bucketName")
	if bucketName == "" {
		util.ReportError(ctx, nil, util.BucketNameArgumentError)
		return
	}

	folder, err := c.client.ListObjects(ctx.Request.Context(), bucketName)
	if err != nil {
		util.ReportError(ctx, err, util.BNBClientUploadError)
		return
	}

	var buffer bytes.Buffer
	zipWriter := zip.NewWriter(&buffer)
	for _, obj := range folder.Objects {
		decryptedByte, err := util.Decrypt(obj.Data, config.PrivateAESKey)
		if err != nil {
			util.ReportError(ctx, err, util.DecryptionError)
			return
		}
		zipFile, err := zipWriter.Create(obj.ObjectName)
		if err != nil {
			util.ReportError(ctx, err, util.ZipEntryCreationError)
			return
		}
		_, err = zipFile.Write(decryptedByte)
		if err != nil {
			util.ReportError(ctx, err, util.ZipWriteError)
			return
		}
	}
	err = zipWriter.Close()
	if err != nil {
		util.ReportError(ctx, err, util.ZipFinalizeError)
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
