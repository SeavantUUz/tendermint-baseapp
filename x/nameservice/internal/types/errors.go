package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
// var (
//	ErrInvalid = sdkerrors.Register(ModuleName, 1, "custom error message")
// )

var (
	ErrNameDoesNotExist    = sdkerrors.Register(ModuleName, 1, "name does not exist")
	ErrAuctionExist        = sdkerrors.Register(ModuleName, 2, "auction exists")
	ErrAuctionDoesNotExist = sdkerrors.Register(ModuleName, 3, "auction does not exist")
	ErrBidPriceTooLow      = sdkerrors.Register(ModuleName, 4, "bid price is too low")
)
