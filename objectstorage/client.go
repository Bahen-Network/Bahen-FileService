// objectstorage/client.go

package objectstorage

import (
	"bytes"
	"context"
	"errors"
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
	cli client.Client
}

func (c *BNBClient) CreateObject(ctx context.Context, bucketName string, objectName string, buffer []byte) (string, error) {
	txnBucketHash, _ := c.CreateBucket(ctx, bucketName)
	log.Printf("Created/Checked bucket with txnHash: %v", txnBucketHash)

	// Upload the object
	txnHash, err := c.cli.CreateObject(ctx, bucketName, objectName, bytes.NewReader(buffer), types.CreateObjectOptions{})
	waitObjectSeal(c.cli, bucketName, objectName)
	if util.HandleErr(err, "CreateObject failed --------2---------.") {
		return "", err
	}

	log.Printf("CreateObject txnHash : %v", txnHash)

	log.Printf("object: %s has been uploaded to SP\n", objectName)
	return txnHash, nil
}

func (c *BNBClient) PutObject(ctx context.Context, bucketName string, objectName string, txnHash string, buffer []byte) error {
	// Put the object

	err := c.cli.PutObject(ctx, bucketName, objectName, int64(len(buffer)),
		bytes.NewReader(buffer), types.PutObjectOptions{TxnHash: txnHash})
	if util.HandleErr(err, "PutObject failed -------3-----------") {
		return err
	}

	log.Printf("object: %s has been update\n", objectName)
	return nil
}

func (c *BNBClient) GetObject(ctx context.Context, bucketName string, objectName string) ([]byte, error) {
	// Get the Object
	reader, info, err := c.cli.GetObject(ctx, bucketName, objectName, types.GetObjectOptions{})
	if util.HandleErr(err, "GetObject failed ------------4------------}") {
		log.Printf("object name: %v", objectName)
		return nil, err
	}

	log.Printf("get object %s successfully, size %d.. ---   %v \n", info.ObjectName, info.Size, info)

	objectBytes, err := io.ReadAll(reader)

	log.Printf("%v", string(objectBytes))
	return objectBytes, nil
}

func (c *BNBClient) CreateBucket(ctx context.Context, bucketName string) (string, error) {

	// get storage providers list
	spLists, err := c.cli.ListStorageProviders(ctx, true)
	util.HandleErr(err, "fail to list in service sps")

	//choose the first sp to be the primary SP
	primarySP := spLists[0].GetOperatorAddress()
	log.Printf("primarySP:  %v", primarySP)

	chargedQuota := uint64(100)
	visibility := storageTypes.VISIBILITY_TYPE_PUBLIC_READ
	opts := types.CreateBucketOptions{Visibility: visibility, ChargedQuota: chargedQuota}
	// bucketName : testbucket
	txnBucketHash, err := c.cli.CreateBucket(ctx, bucketName, primarySP, opts)
	util.HandleErr(err, "CreateBucket failed ------------1-------.")
	log.Printf("Create Bucket: txnHash: %v", txnBucketHash)

	return txnBucketHash, nil
}

func NewClient(privateKey string, chainId string, rpcAddr string) *BNBClient {
	account, err := types.NewAccountFromPrivateKey("test", privateKey)
	util.HandleErr(err, "New account from private key error")

	cli, err := client.New(chainId, rpcAddr, client.Option{DefaultAccount: account})
	util.HandleErr(err, "unable to new greenfield client")
	return &BNBClient{cli}
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
