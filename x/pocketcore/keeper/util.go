package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

// "GetPKFromFile" - Returns the private key object from a file
func (k Keeper) GetPKFromFile(ctx sdk.Ctx) (crypto.PrivateKey, error) {
	// get the Private validator key from the file
	pvKey, err := types.GetPVKeyFile()
	if err != nil {
		return nil, err
	}
	// convert the privKey to a private key object (compatible interface)
	pk, er := crypto.PrivKeyToPrivateKey(pvKey.PrivKey)
	if er != nil {
		return nil, er
	}
	return pk, nil
}
