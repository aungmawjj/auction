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
	PendingAuction []byte
}

type Auction struct {
	ID            []byte
	Ended         bool
	HighestBid    int64
	HighestBidder []byte
}

type SetAuctionArgs struct {
	AssetID   []byte
	AuctionID []byte
}

func (cc *AssetCC) Deploy() (string, error) {
	return cc.fabric.SendChaincodeRequest("deploy", "Init", nil)
}

func (cc *AssetCC) AddAsset(asset Asset) (string, error) {
	b, _ := json.Marshal(asset)
	return cc.fabric.SendChaincodeRequest("invoke", "addAsset", []string{string(b)})
}

func (cc *AssetCC) BindAuction(args SetAuctionArgs) (string, error) {
	b, _ := json.Marshal(args)
	return cc.fabric.SendChaincodeRequest("invoke", "bindAuction", []string{string(b)})
}

func (cc *AssetCC) UpdateAuction(auction Auction) (string, error) {
	b, _ := json.Marshal(auction)
	return cc.fabric.SendChaincodeRequest("invoke", "updateAuction", []string{string(b)})
}

func (cc *AssetCC) EndAuction(assetID []byte) (string, error) {
	s := base64.StdEncoding.EncodeToString(assetID)
	return cc.fabric.SendChaincodeRequest("invoke", "endAuction", []string{s})
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
