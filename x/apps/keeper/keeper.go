package keeper

import (
	"container/list"
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/libs/log"
)

const aminoCacheSize = 500

// Implements ApplicationSet interface
var _ types.ApplicationSet = Keeper{}

// keeper of the staking store
type Keeper struct {
	storeKey             sdk.StoreKey
	cdc                  *codec.Codec
	AccountsKeeper       types.AuthKeeper
	POSKeeper            types.PosKeeper
	Paramstore           sdk.Subspace
	applicationCache     map[string]cachedApplication
	applicationCacheList *list.List

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, posKeeper types.PosKeeper, supplyKeeper types.AuthKeeper,
	paramstore sdk.Subspace, codespace sdk.CodespaceType) Keeper {

	// ensure staked module accounts are set
	if addr := supplyKeeper.GetModuleAddress(types.StakedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.StakedPoolName))
	}

	return Keeper{
		storeKey:             key,
		cdc:                  cdc,
		AccountsKeeper:       supplyKeeper,
		POSKeeper:            posKeeper,
		Paramstore:           paramstore.WithKeyTable(ParamKeyTable()),
		applicationCache:     make(map[string]cachedApplication, aminoCacheSize),
		applicationCacheList: list.New(),
		codespace:            codespace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Ctx) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}
