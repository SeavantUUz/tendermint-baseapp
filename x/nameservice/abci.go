package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rune/baseapp/x/nameservice/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	// 	TODO: fill out if your application requires beginblock, if not you can delete this function
	iterator := k.GetAuctionIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		bz := iterator.Value()
		var auction types.Auction
		ModuleCdc.MustUnmarshalBinaryBare(bz, &auction)
		if req.Header.Height > auction.Deadline {
			if !auction.Bidder.Empty() && !auction.Bidder.Equals(auction.Owner) {
				_, err := k.CoinKeeper.SubtractCoins(ctx, auction.Bidder, auction.BidPrice)
				if err == nil {
					k.SetOwner(ctx, auction.Lot, auction.Bidder)
				}
			}
			k.DeleteAuction(ctx, auction.Lot)
		}
	}
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k Keeper) {
	// 	TODO: fill out if your application requires endblock, if not you can delete this function
}
