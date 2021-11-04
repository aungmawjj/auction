package events

type OnBindAuction struct {
	AssetCC   []byte
	AssetID   []byte
	AuctionID []byte
}

type OnEndAuction struct {
	AssetCC       []byte
	AssetID       []byte
	AuctionID     []byte
	HighestBidder []byte
}
