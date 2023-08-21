// storageclient/bnb_client.go

package storageclient

import (
	"bytes"
	"context"
	"errors"
	"file-service/models"
	"file-service/util"
	"fmt"
	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	storageTypes "github.com/bnb-chain/greenfield/x/storage/types"
	"io"
	"log"
	"time"
)

type BNBClient struct {
	cli          client.Client
	primarySP    string
	chargedQuota uint64
	visibility   storageTypes.VisibilityType
	opts         types.CreateBucketOptions
}

func NewClient(privateKey string, chainId string, rpcAddr string) *BNBClient {
	account, err := types.NewAccountFromPrivateKey("test", privateKey)
	util.HandleErr(err, "New account from private key error")

	cli, err := client.New(chainId, rpcAddr, client.Option{DefaultAccount: account})
	util.HandleErr(err, "unable to new greenfield client")

	c := &BNBClient{cli: cli}

	// get storage providers list
	ctx := context.Background() // Create a background context, can be replaced with a more relevant context if necessary
	spLists, err := c.cli.ListStorageProviders(ctx, true)
	util.HandleErr(err, "fail to list in service sps")

	// choose the first sp to be the primary SP
	c.primarySP = spLists[0].GetOperatorAddress()
	log.Printf("primarySP:  %v", c.primarySP)

	c.chargedQuota = uint64(100)
	c.visibility = storageTypes.VISIBILITY_TYPE_PUBLIC_READ

	c.opts = types.CreateBucketOptions{Visibility: c.visibility, ChargedQuota: c.chargedQuota}

	return c
}

func (c *BNBClient) CreateObject(ctx context.Context, bucketName string, objectName string, buffer []byte) (string, error) {
	txnBucketHash, _ := c.CreateBucket(ctx, bucketName)
	log.Printf("Created/Checked bucket with txnHash: %v", txnBucketHash)

	log.Printf("Start upload object")
	// Upload the object
	txnHash, err := c.cli.CreateObject(ctx, bucketName, objectName, bytes.NewReader(buffer), types.CreateObjectOptions{})
	// waitObjectSeal(c.cli, bucketName, objectName)
	if util.HandleErr(err, "CreateObject failed --------2---------.") {
		return "", err
	}

	err = c.cli.PutObject(ctx, bucketName, objectName, int64(len(buffer)),
		bytes.NewReader(buffer), types.PutObjectOptions{TxnHash: txnHash})
	util.HandleErr(err, "PutObject")

	log.Printf("CreateObject txnHash : %v", txnHash)

	log.Printf("object: %s has been uploaded to SP\n", objectName)
	return txnHash, nil
}

func (c *BNBClient) GetObject(ctx context.Context, bucketName string, objectName string) ([]byte, error) {
	// Get the Object
	reader, info, err := c.cli.GetObject(ctx, bucketName, objectName, types.GetObjectOptions{})
	if util.HandleErr(err, "GetObject failed ------------4------------}") {
		log.Printf("object name: %v", objectName)
		return nil, err
	}

	log.Printf("get object %s successfully, size %d.. \n", info.ObjectName, info.Size)
	objectBytes, err := io.ReadAll(reader)
	return objectBytes, nil
}

func (c *BNBClient) ListObjects(ctx context.Context, bucketName string) (models.ObjectListResponse, error) {
	objects, err := c.cli.ListObjects(ctx, bucketName, types.ListObjectsOptions{
		ShowRemovedObject: false, Delimiter: "", MaxKeys: 100, EndPointOptions: &types.EndPointOptions{
			Endpoint:  "",
			SPAddress: "",
		}})
	if err != nil {
		return models.ObjectListResponse{}, err
	}

	var response models.ObjectListResponse
	for _, obj := range objects.Objects {
		objectBytes, err := c.GetObject(ctx, bucketName, obj.ObjectInfo.ObjectName)
		if err != nil {
			util.HandleErr(err, "")
			continue
		}
		objectInfo := models.ObjectInfo{
			ObjectName: obj.ObjectInfo.ObjectName,
			Data:       objectBytes,
			Type:       obj.ObjectInfo.ContentType,
		}
		response.Objects = append(response.Objects, objectInfo)
	}
	return response, nil
}

func (c *BNBClient) CreateBucket(ctx context.Context, bucketName string) (string, error) {
	// bucketName : testbucket
	txnBucketHash, err := c.cli.CreateBucket(ctx, bucketName, c.primarySP, c.opts)
	util.HandleErr(err, "CreateBucket failed ------------1-------.")
	log.Printf("Create Bucket: txnHash: %v", txnBucketHash)

	return txnBucketHash, nil
}

func waitObjectSeal(cli client.Client, bucketName, objectName string) {
	ctx := context.Background()
	// wait for the object to be sealed
	timeout := time.After(15 * time.Second)
	ticker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-timeout:
			err := errors.New("object not sealed after 15 seconds")
			util.HandleErr(err, "")
		case <-ticker.C:
			objectDetail, err := cli.HeadObject(ctx, bucketName, objectName)
			util.HandleErr(err, "HeadObject")
			if objectDetail.ObjectInfo.GetObjectStatus().String() == "OBJECT_STATUS_SEALED" {
				ticker.Stop()
				fmt.Printf("put object %s successfully \n", objectName)
				return
			}
		}
	}
}
