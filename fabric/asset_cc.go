package fabric

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type AssetCC struct {
	contract *gateway.Contract
}

func NewAssetCC() *AssetCC {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	ccpPath := filepath.Join(
		"/",
		"home",
		"ubuntu",
		"fabric2",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)

	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	return &AssetCC{
		contract: network.GetContract("asset"),
	}
}

type Asset struct {
	ID               string
	Owner            string
	PendingAuctionID int
}

type Auction struct {
	ID              int
	AssetID         string
	Platforms       []string
	CrossAuctionIDs []string
	Status          string

	HighestBid         int
	HighestBidder      string
	HighestBidPlatform string
}

type StartAuctionArgs struct {
	AssetID   string
	Platforms []string
}

type BindAuctionArgs struct {
	AuctionID       int
	CrossAuctionIDs []string
}

type EndAuctionArgs struct {
	AuctionID      int
	HighestBids    []int
	HighestBidders []string
}

func (cc *AssetCC) GetCCID() string {
	return "asset"
}

func (cc *AssetCC) AddAsset(id, owner string) ([]byte, error) {
	return cc.contract.SubmitTransaction("AddAsset", id, owner)
}

func (cc *AssetCC) StartAuction(args StartAuctionArgs) ([]byte, error) {
	b, _ := json.Marshal(args)
	return cc.contract.SubmitTransaction("StartAuction", string(b))
}

func (cc *AssetCC) BindAuction(args BindAuctionArgs) ([]byte, error) {
	b, _ := json.Marshal(args)
	return cc.contract.SubmitTransaction("BindAuction", string(b))
}

func (cc *AssetCC) EndAuction(args EndAuctionArgs) ([]byte, error) {
	b, _ := json.Marshal(args)
	return cc.contract.SubmitTransaction("EndAuction", string(b))
}

func (cc *AssetCC) SetAuctionEnding(assetID string) ([]byte, error) {
	return cc.contract.SubmitTransaction("SetAuctionEnding", assetID)
}

func (cc *AssetCC) GetAsset(assetID string) (*Asset, error) {
	var asset Asset
	res, err := cc.contract.EvaluateTransaction("GetAsset", assetID)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &asset)
	return &asset, err
}

func (cc *AssetCC) GetAuction(auctionID int) (*Auction, error) {
	var auction Auction
	res, err := cc.contract.EvaluateTransaction("GetAuction", strconv.Itoa(auctionID))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &auction)
	return &auction, err
}

func (cc *AssetCC) GetLastAuctionID() (int, error) {
	res, err := cc.contract.EvaluateTransaction("GetLastAuctionID")
	if err != nil {
		return 0, err
	}
	var id int
	err = json.Unmarshal(res, &id)
	return id, err
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"/",
		"home",
		"ubuntu",
		"fabric2",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "User1@org1.example.com-cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}
