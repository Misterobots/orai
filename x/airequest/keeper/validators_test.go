package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/oraichain/orai/packages/rng"
	"github.com/oraichain/orai/x/airequest/keeper"
	"github.com/oraichain/orai/x/airequest/types"
)

func TestCalucateMol(t *testing.T) {
	k := &keeper.Keeper{}
	size := 10
	maxValidatorSize := 15
	totalPowers := int64(92)
	randomGenerator, _ := rng.NewRng(make([]byte, types.RngSeedSize), []byte("nonce"), []byte("Oraichain"))
	valOperators := make([]sdk.ValAddress, maxValidatorSize)
	for i := 0; i < maxValidatorSize; i++ {
		valOperators = append(valOperators, ed25519.GenPrivKey().PubKey().Address().Bytes())
	}

	validators := k.SampleIndexes(valOperators, size, randomGenerator, totalPowers)

	t.Logf("Validators :%v\n", validators)
}
