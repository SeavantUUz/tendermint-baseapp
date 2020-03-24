package util

const AuctionPrefix = "Auction:"

func AuctionName(name string) string {
	return AuctionPrefix + name
}
