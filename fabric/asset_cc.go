package fabric

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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

func (cc *AssetCC) GetCCID() string {
	return "asset"
}

func (cc *AssetCC) AddAsset(asset Asset) ([]byte, error) {
	b, _ := json.Marshal(asset)
	return cc.contract.SubmitTransaction("AddAsset", string(b))
}

func (cc *AssetCC) BindAuction(args BindAuctionArgs) ([]byte, error) {
	b, _ := json.Marshal(args)
	return cc.contract.SubmitTransaction("BindAuction", string(b))
}

func (cc *AssetCC) EndAuction(args EndAuctionArgs) ([]byte, error) {
	b, _ := json.Marshal(args)
	return cc.contract.SubmitTransaction("EndAuction", string(b))
}

func (cc *AssetCC) GetAsset(assetID []byte) (Asset, error) {
	var asset Asset
	s := base64.StdEncoding.EncodeToString(assetID)
	res, err := cc.contract.EvaluateTransaction("GetAsset", s)
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
