package main

import (
	"bytes"
	"context"
	"errors"
	"file-service/util"
	"fmt"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	storageTypes "github.com/bnb-chain/greenfield/x/storage/types"
	"log"
	"time"

	"github.com/bnb-chain/greenfield-go-sdk/client"
)

const (
	rpcAddr        = "https://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443"
	chainId        = "greenfield_5600-1"
	privateKey     = "86c6252d772b7a85fd566e19d1dab0a7f6b246348bc133689633db4c0322cb14"
	privateKeyTest = "0x49993248a7c8d748aa68ff249a1bceac91f58eb88b77e9a67daf767c981ea1fd"
)

func handleErr(err error, msg string) {
	if err != nil {
		log.Printf("%s, %v", msg, err)
	}
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
			// log.Printf("%v", ticker)
		}
	}
}

func main() {
	account, err := types.NewAccountFromPrivateKey("test", privateKey)
	handleErr(err, "New account from private key error")

	cli, err := client.New(chainId, rpcAddr, client.Option{DefaultAccount: account})
	handleErr(err, "unable to new greenfield client")
	ctx := context.Background()

	// get storage providers list
	spLists, err := cli.ListStorageProviders(ctx, true)
	handleErr(err, "fail to list in service sps")

	//choose the first sp to be the primary SP
	primarySP := spLists[7].GetOperatorAddress()
	log.Printf("primarySP:  %v", primarySP)

	bucketName := "testbucket07"

	buffer := []byte(bucketName)
	// Create Bucket (Created testbucket)

	chargedQuota := uint64(100)
	visibility := storageTypes.VISIBILITY_TYPE_PUBLIC_READ
	opts := types.CreateBucketOptions{Visibility: visibility, ChargedQuota: chargedQuota}
	txnBucketHash, err := cli.CreateBucket(ctx, bucketName, primarySP, opts)
	handleErr(err, "CreateBucket failed ------------1-------.")
	log.Printf("Create Bucket: txnHash: %v", txnBucketHash)

	// Upload the object
	objectName := "testString2"

	txnHash, err := cli.CreateObject(ctx, bucketName, objectName, bytes.NewReader(buffer), types.CreateObjectOptions{})

	log.Printf("CreateObject txnHash : %v", txnHash)
	handleErr(err, "CreateObject failed --------2---------.")
	// waitObjectSeal(cli, bucketName, objectName)

	/*err = cli.PutObject(ctx, bucketName, objectName, int64(len(buffer)),
		bytes.NewReader(buffer), types.PutObjectOptions{TxnHash: txnHash})
	cli.WaitForTx(ctx, txnHash)
	handleErr(err, "PutObject failed -------3-----------")
	*/
	log.Printf("object: %s has been uploaded to SP\n", objectName)

	// Get the Object
	// get object
	_, info, err := cli.GetObject(ctx, bucketName, objectName, types.GetObjectOptions{
		Range:            "",
		SupportRecovery:  true,
		SupportResumable: true,
		PartSize:         1024 * 1024, // For example, set to 1MB
	})
	handleErr(err, "GetObject failed ------------4------------")
	log.Printf("get object %s successfully, size %d.. ---   %v \n", info.ObjectName, info.Size, info)

	log.Printf("listObject strat!")
	// list objects
	objects, err := cli.ListObjects(ctx, bucketName, types.ListObjectsOptions{
		ShowRemovedObject: false, Delimiter: "", MaxKeys: 100, EndPointOptions: &types.EndPointOptions{
			Endpoint:  "",
			SPAddress: "",
		}})
	log.Println("list objects result:")
	for _, obj := range objects.Objects {
		i := obj.ObjectInfo
		log.Printf("object: %s, status: %s  %v\n", i.ObjectName, i.ObjectStatus, i.ContentType)
	}

	//log.Printf("%v", string(objectBytes))
}
