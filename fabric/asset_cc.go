package fabric

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type AssetCC struct {
	fabric *FabricClient
}

func NewAssetCC(fabric *FabricClient) *AssetCC {
	return &AssetCC{
		fabric: fabric,
	}
}

type Asset struct {
	ID             []byte
	Owner          []byte
	PendingAuction *Auction
}

type Auction struct {
	ID       []byte
	Platform string
}

type AuctionResult struct {
	Auction
	HighestBidder []byte
}

type BindAuctionArgs struct {
	AssetID []byte
	Auction Auction
}

type EndAuctionArgs struct {
	AssetID       []byte
	AuctionResult AuctionResult
}

func (cc *AssetCC) Deploy() (string, error) {
	return cc.fabric.SendChaincodeRequest("deploy", "Init", nil)
}

func (cc *AssetCC) SetChaincodeID(ccid string) {
	cc.fabric.ChaincodeID = ccid
}

func (cc *AssetCC) GetCCID() string {
	return cc.fabric.ChaincodeID
}

func (cc *AssetCC) AddAsset(asset Asset) (string, error) {
	b, _ := json.Marshal(asset)
	return cc.fabric.SendChaincodeRequest("invoke", "addAsset", []string{string(b)})
}

func (cc *AssetCC) BindAuction(args BindAuctionArgs) (string, error) {
	b, _ := json.Marshal(args)
	return cc.fabric.SendChaincodeRequest("invoke", "bindAuction", []string{string(b)})
}

func (cc *AssetCC) EndAuction(args EndAuctionArgs) (string, error) {
	b, _ := json.Marshal(args)
	return cc.fabric.SendChaincodeRequest("invoke", "endAuction", []string{string(b)})
}

func (cc *AssetCC) GetAsset(assetID []byte) (Asset, error) {
	var asset Asset
	s := base64.StdEncoding.EncodeToString(assetID)
	res, err := cc.fabric.SendChaincodeRequest("query", "getAsset", []string{s})
	if err != nil {
		return asset, err
	}
	b := []byte(res)
	if b == nil {
		return asset, fmt.Errorf("asset not found")
	}
	err = json.Unmarshal(b, &asset)
	return asset, err
}
