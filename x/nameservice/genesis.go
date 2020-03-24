package nameservice

import (
	"fmt"
	"github.com/rune/baseapp/x/nameservice/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	WhoisRecords   []Whois   `json:"whois_records"`
	AuctionRecords []Auction `json:"auction_records"`
}

func NewGenesisState(whoIsRecords []Whois, auctionRecords []Auction) GenesisState {
	return GenesisState{WhoisRecords: nil, AuctionRecords: auctionRecords}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.WhoisRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Owner", record.Value)
		}
		if record.Value == "" {
			return fmt.Errorf("invalid WhoisRecord: Owner: %s. Error: Missing Value", record.Owner)
		}
		if record.Price == nil {
			return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Price", record.Value)
		}
	}
	for _, record := range data.AuctionRecords {
		if record.Lot == "" {
			return fmt.Errorf("invalid Auction: Owner: %s. Error: Missing Lot", record.Owner)
		}
		if record.Owner == nil {
			return fmt.Errorf("invalid Auction: Lot: %s. Error: Missing Owner", record.Lot)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		WhoisRecords:   []Whois{},
		AuctionRecords: []Auction{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.WhoisRecords {
		keeper.SetWhois(ctx, record.Value, record)
	}
	for _, record := range data.AuctionRecords {
		keeper.SetAuction(ctx, util.AuctionName(record.Lot), record, true)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []Whois
	var auctionRecords []Auction
	iterator := k.GetNamesIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {

		name := string(iterator.Key())
		whois := k.GetWhois(ctx, name)
		records = append(records, whois)

	}
	for ; iterator.Valid(); iterator.Next() {
		key := string(iterator.Key())
		auction := k.GetRawAuction(ctx, key)
		auctionRecords = append(auctionRecords, auction)
	}
	return GenesisState{WhoisRecords: records, AuctionRecords: auctionRecords}
}
