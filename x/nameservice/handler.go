package nameservice

import (
	"fmt"

	"github.com/rune/baseapp/x/nameservice/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case types.MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		case types.MsgDeleteName:
			return handleMsgDeleteName(ctx, keeper, msg)
		case types.MsgAuction:
			return handleMsgAuction(ctx, keeper, msg)
		case types.MsgBid:
			return handleMsgBid(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type()))
		}
	}
}

// Handle a message to set name
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg types.MsgSetName) (*sdk.Result, error) {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) { // Checks if the the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}
	keeper.SetName(ctx, msg.Name, msg.Value) // If so, set the name to the value specified in the msg.
	return &sdk.Result{}, nil                // return
}

// Handle a message to buy name
func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg types.MsgBuyName) (*sdk.Result, error) {
	// Checks if the the bid price is greater than the price paid by the current owner
	if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid not high enough") // If not, throw an error
	}
	if keeper.HasOwner(ctx, msg.Name) {
		err := keeper.CoinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := keeper.CoinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return nil, err
		}
	}
	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return &sdk.Result{}, nil
}

// Handle a message to delete name
// Handle a message to delete name
func handleMsgDeleteName(ctx sdk.Context, keeper Keeper, msg types.MsgDeleteName) (*sdk.Result, error) {
	if !keeper.IsNamePresent(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExist, msg.Name)
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	keeper.DeleteWhois(ctx, msg.Name)
	return &sdk.Result{}, nil
}

func handleMsgAuction(ctx sdk.Context, keeper Keeper, msg types.MsgAuction) (*sdk.Result, error) {
	if keeper.HasAuction(ctx, msg.Lot) {
		return nil, sdkerrors.Wrap(types.ErrAuctionExist, fmt.Sprintf("Auction %s has existed", msg.Lot))
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Lot)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}
	if !msg.ReservePrice.IsAllPositive() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFee, "Reserve Price should be positive")
	}
	auction := types.Auction{Lot: msg.Lot, Owner: msg.Owner, ReservePrice: msg.ReservePrice}
	keeper.SetAuction(ctx, msg.Lot, auction, false)
	return &sdk.Result{}, nil
}

func handleMsgBid(ctx sdk.Context, keeper Keeper, msg types.MsgBid) (*sdk.Result, error) {
	if !keeper.IsNamePresent(ctx, msg.Lot) {
		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExist, fmt.Sprintf("Name %s is not presented", msg.Lot))
	}
	if !keeper.HasAuction(ctx, msg.Lot) {
		return nil, sdkerrors.Wrap(types.ErrAuctionDoesNotExist, fmt.Sprintf("Auction %s is not existed", msg.Lot))
	}
	auction := keeper.GetAuction(ctx, msg.Lot)
	if !msg.BidPrice.IsAllGT(auction.ReservePrice) {
		return nil, sdkerrors.Wrap(types.ErrBidPriceTooLow, msg.BidPrice.String())
	}
	if !msg.BidPrice.IsAllGT(auction.BidPrice) {
		return nil, sdkerrors.Wrap(types.ErrBidPriceTooLow, msg.BidPrice.String())
	}
	keeper.SetBid(ctx, msg.Lot, msg.BidPrice, msg.Bidder)
	return &sdk.Result{}, nil
}
