package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MinNamePrice is Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// NewWhois returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

// implement fmt.Stringer
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s
Price: %s`, w.Owner, w.Value, w.Price))
}

type Auction struct {
	Lot          string         `json:"lot"`
	Owner        sdk.AccAddress `json:"owner"`
	ReservePrice sdk.Coins      `json:"reserve_price"`
	Bidder       sdk.AccAddress `json:"bidder"`
	BidPrice     sdk.Coins      `json:"bid_price"`
	Deadline     int64          `json:"deadline"`
}

func NewAuction() Auction {
	return Auction{}
}

func (a Auction) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Lot: %s Owner: %s Reserve Price: %s Bidder:%s BidPrice: %s Deadline: %d`,
		a.Lot, a.Owner, a.ReservePrice, a.Bidder, a.BidPrice, a.Deadline))
}
