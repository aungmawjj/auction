//go:generate abigen --sol auction.sol --pkg contract --type Auction --out auction_gen.go
//go:generate abigen --sol asset.sol --pkg contract --type Asset --out asset_gen.go

package contract
