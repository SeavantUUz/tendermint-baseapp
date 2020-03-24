package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/rune/baseapp/x/nameservice/internal/types"
	"github.com/rune/baseapp/x/nameservice/util"
)

// Keeper of the nameservice store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	CoinKeeper types.BankKeeper
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		CoinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// Sets the entire Whois metadata struct for a name
func (k Keeper) SetWhois(ctx sdk.Context, name string, whois types.Whois) {
	if whois.Owner.Empty() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(whois))
}

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetWhois(ctx sdk.Context, name string) types.Whois {
	store := ctx.KVStore(k.storeKey)
	if !k.IsNamePresent(ctx, name) {
		return types.NewWhois()
	}
	bz := store.Get([]byte(name))
	var whois types.Whois
	k.cdc.MustUnmarshalBinaryBare(bz, &whois)
	return whois
}

// Deletes the entire Whois metadata struct for a name
func (k Keeper) DeleteWhois(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(name))
}

// ResolveName - returns the string that the name resolves to
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	return k.GetWhois(ctx, name).Value
}

// SetName - sets the value string that a name resolves to
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	whois := k.GetWhois(ctx, name)
	whois.Value = value
	k.SetWhois(ctx, name, whois)
}

// HasOwner - returns whether or not the name already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	return !k.GetWhois(ctx, name).Owner.Empty()
}

// GetOwner - get the current owner of a name
func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetWhois(ctx, name).Owner
}

// SetOwner - sets the current owner of a name
func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	whois := k.GetWhois(ctx, name)
	whois.Owner = owner
	k.SetWhois(ctx, name, whois)
}

// GetPrice - gets the current price of a name
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	return k.GetWhois(ctx, name).Price
}

// SetPrice - sets the current price of a name
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois := k.GetWhois(ctx, name)
	whois.Price = price
	k.SetWhois(ctx, name, whois)
}

// Check if the name is present in the store or not
func (k Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

// it is not the best practice, different module should be stored in different keeper.
func (k Keeper) SetAuction(ctx sdk.Context, lot string, auction types.Auction, isGenesis bool) {
	if auction.Owner.Empty() {
		return
	}

	if !k.IsNamePresent(ctx, lot) {
		return
	}
	store := ctx.KVStore(k.storeKey)
	if !isGenesis {
		auction.Deadline = ctx.BlockHeight() + 100
	}
	store.Set([]byte(util.AuctionName(lot)), k.cdc.MustMarshalBinaryBare(auction))
}

func (k Keeper) GetAuction(ctx sdk.Context, lot string) types.Auction {
	store := ctx.KVStore(k.storeKey)
	if !k.IsNamePresent(ctx, lot) {
		return types.NewAuction()
	}
	bz := store.Get([]byte(util.AuctionName(lot)))
	var auction types.Auction
	k.cdc.MustUnmarshalBinaryBare(bz, &auction)
	return auction
}

func (k Keeper) GetRawAuction(ctx sdk.Context, lotKey string) types.Auction {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(lotKey))
	var auction types.Auction
	k.cdc.MustUnmarshalBinaryBare(bz, &auction)
	return auction
}

func (k Keeper) HasAuction(ctx sdk.Context, lot string) bool {
	return !k.GetAuction(ctx, lot).Owner.Empty()
}

func (k Keeper) DeleteAuction(ctx sdk.Context, lot string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(util.AuctionName(lot)))
}

// do not check in keeper
func (k Keeper) SetBid(ctx sdk.Context, lot string, price sdk.Coins, bidder sdk.AccAddress) {
	auction := k.GetAuction(ctx, lot)
	auction.BidPrice = price
	auction.Bidder = bidder
	k.SetAuction(ctx, lot, auction, false)
}

func (k Keeper) GetAuctionOwner(ctx sdk.Context, lot string) sdk.AccAddress {
	return k.GetAuction(ctx, lot).Owner
}

func (k Keeper) GetBid(ctx sdk.Context, lot string) sdk.Coins {
	return k.GetAuction(ctx, lot).BidPrice
}

func (k Keeper) ResolveBidder(ctx sdk.Context, lot string) string {
	return k.GetAuction(ctx, lot).Bidder.String()
}

func (k Keeper) GetDeadline(ctx sdk.Context, lot string) int64 {
	return k.GetAuction(ctx, lot).Deadline
}

func (k Keeper) GetAuctionIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(util.AuctionPrefix))
}
