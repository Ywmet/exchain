package simulation_test

import (
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	//	"github.com/cosmos/cosmos-sdk/codec"
	//	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	//	"github.com/cosmos/cosmos-sdk/types/module"
	//	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	simtypes "github.com/okex/exchain/libs/cosmos-sdk/x/simulation"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/simulation"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
)

// TestRandomizedGenState tests the normal scenario of applying RandomizedGenState.
// Abonormal scenarios are not tested here.
func TestRandomizedGenState(t *testing.T) {
	//interfaceRegistry := codectypes.NewInterfaceRegistry()
	// cdc := codec.NewProtoCodec(interfaceRegistry)
	cdc := codec.New()

	s := rand.NewSource(1)
	r := rand.New(s)

	simState := module.SimulationState{
		AppParams:    make(simtypes.AppParams),
		Cdc:          cdc,
		Rand:         r,
		NumBonded:    3,
		Accounts:     simtypes.RandomAccounts(r, 3),
		InitialStake: 1000,
		GenState:     make(map[string]json.RawMessage),
	}

	simulation.RandomizedGenState(&simState)

	var ibcTransferGenesis types.GenesisState
	simState.Cdc.MustUnmarshalJSON(simState.GenState[types.ModuleName], &ibcTransferGenesis)

	require.Equal(t, "euzxpfgkqegqiqwixnku", ibcTransferGenesis.PortId)
	require.True(t, ibcTransferGenesis.Params.SendEnabled)
	require.True(t, ibcTransferGenesis.Params.ReceiveEnabled)
	require.Len(t, ibcTransferGenesis.DenomTraces, 0)

}

// TestRandomizedGenState tests abnormal scenarios of applying RandomizedGenState.
func TestRandomizedGenState1(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	s := rand.NewSource(1)
	r := rand.New(s)
	// all these tests will panic
	tests := []struct {
		simState module.SimulationState
		panicMsg string
	}{
		{ // panic => reason: incomplete initialization of the simState
			module.SimulationState{}, "invalid memory address or nil pointer dereference"},
		{ // panic => reason: incomplete initialization of the simState
			module.SimulationState{
				AppParams: make(simtypes.AppParams),
				Cdc:       cdc,
				Rand:      r,
			}, "assignment to entry in nil map"},
	}

	for _, tt := range tests {
		require.Panicsf(t, func() { simulation.RandomizedGenState(&tt.simState) }, tt.panicMsg)
	}
}
