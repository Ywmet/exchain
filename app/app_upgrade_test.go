package app

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/store"
	cosmost "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	upgradetypes "github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/okex/exchain/x/params"
)

var (
	_ upgradetypes.UpgradeModule = (*SimpleBaseUpgradeModule)(nil)

	defaultHandleStore upgradetypes.HandleStore = func(st cosmost.CommitKVStore, h int64) {
		st.SetUpgradeVersion(h)
	}
)

type SimpleBaseUpgradeModule struct {
	h                  int64
	taskExecutedNotify func()
	appModule          module.AppModuleBasic
	handler            upgradetypes.HandleStore
}

func NewSimpleBaseUpgradeModule(h int64, appModule module.AppModuleBasic, handler upgradetypes.HandleStore, taskExecutedNotify func()) *SimpleBaseUpgradeModule {
	return &SimpleBaseUpgradeModule{h: h, appModule: appModule, handler: handler, taskExecutedNotify: taskExecutedNotify}
}

func (b *SimpleBaseUpgradeModule) ModuleName() string {
	return b.appModule.Name()
}

func (b *SimpleBaseUpgradeModule) RegisterTask() upgradetypes.HeightTask {
	return upgradetypes.NewHeightTask(0, func(ctx sdk.Context) error {
		b.taskExecutedNotify()
		return nil
	})
}

func (b *SimpleBaseUpgradeModule) UpgradeHeight() int64 {
	return b.h
}

func (b *SimpleBaseUpgradeModule) BlockStoreModules() map[string]upgradetypes.HandleStore {
	return map[string]upgradetypes.HandleStore{
		b.ModuleName(): b.handler,
	}
}

func (b *SimpleBaseUpgradeModule) RegisterParam() params.ParamSet {
	return nil
}

func (b *SimpleBaseUpgradeModule) HandleStoreWhenMeetUpgradeHeight() upgradetypes.HandleStore {
	return func(st store.CommitKVStore, h int64) {
		st.SetUpgradeVersion(h)
	}
}

var (
	_ module.AppModuleBasic = (*simpleDefaultAppModuleBasic)(nil)
)

type simpleDefaultAppModuleBasic struct {
	name string
}

func (s *simpleDefaultAppModuleBasic) Name() string {
	return s.name
}

func (s *simpleDefaultAppModuleBasic) RegisterCodec(c *codec.Codec) {}

func (s *simpleDefaultAppModuleBasic) DefaultGenesis() json.RawMessage { return nil }

func (s *simpleDefaultAppModuleBasic) ValidateGenesis(message json.RawMessage) error { return nil }

func (s *simpleDefaultAppModuleBasic) RegisterRESTRoutes(context context.CLIContext, router *mux.Router) {
	return
}

func (s *simpleDefaultAppModuleBasic) GetTxCmd(c *codec.Codec) *cobra.Command { return nil }

func (s *simpleDefaultAppModuleBasic) GetQueryCmd(c *codec.Codec) *cobra.Command { return nil }

var (
	_ module.AppModule = (*simpleAppModule)(nil)
)

type simpleAppModule struct {
	*SimpleBaseUpgradeModule
	*simpleDefaultAppModuleBasic
}

func newSimpleAppModule(h int64, name string, notify func()) *simpleAppModule {
	ret := &simpleAppModule{}
	ret.simpleDefaultAppModuleBasic = &simpleDefaultAppModuleBasic{name: name}
	ret.SimpleBaseUpgradeModule = NewSimpleBaseUpgradeModule(h, ret, func(st cosmost.CommitKVStore, h int64) {
		st.SetUpgradeVersion(h)
	}, notify)
	return ret
}

func (s2 *simpleAppModule) InitGenesis(s sdk.Context, message json.RawMessage) []abci.ValidatorUpdate {
	return nil
}

func (s2 *simpleAppModule) ExportGenesis(s sdk.Context) json.RawMessage {
	return nil
}

func (s2 *simpleAppModule) RegisterInvariants(registry sdk.InvariantRegistry) { return }

func (s2 *simpleAppModule) Route() string {
	return ""
}

func (s2 *simpleAppModule) NewHandler() sdk.Handler { return nil }

func (s2 *simpleAppModule) QuerierRoute() string {
	return ""
}

func (s2 *simpleAppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

func (s2 *simpleAppModule) BeginBlock(s sdk.Context, block abci.RequestBeginBlock) {}

func (s2 *simpleAppModule) EndBlock(s sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}

func setupModuleBasics(bs ...module.AppModule) *module.Manager {
	basis := []module.AppModule{}
	for _, v := range bs {
		basis = append(basis, v)
	}
	return module.NewManager(
		basis...,
	)
}

type testSimApp struct {
	*OKExChainApp
	// the module manager
}

func newTestSimApp(name string, logger log.Logger, db dbm.DB, txDecoder sdk.TxDecoder, mm *module.Manager) *testSimApp {
	ret := &testSimApp{}
	//ret.BaseApp = bam.NewBaseApp(name, logger, db, txDecoder)
	ret.OKExChainApp = NewOKExChainApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, 0)
	ret.mm = mm
	ret.setupUpgradeModules()
	return ret
}

type UpgradeCase struct {
	name     string
	upgradeH int64
}

func createCases(moduleCount int, beginHeight int) []UpgradeCase {
	ret := make([]UpgradeCase, moduleCount)
	for i := 0; i < moduleCount; i++ {
		ret[i] = UpgradeCase{
			name:     "m_" + strconv.Itoa(i),
			upgradeH: int64(beginHeight + i),
		}
	}
	return ret
}
func TestUpgradeWithConcreteHeight(t *testing.T) {
	db := dbm.NewMemDB()

	cases := createCases(5, 10)
	m := make(map[string]int)
	modules := make([]module.AppModule, 0)
	count := 0
	maxHeight := int64(0)
	for _, ca := range cases {
		c := ca
		m[c.name] = 0
		if maxHeight < c.upgradeH {
			maxHeight = c.upgradeH
		}
		modules = append(modules, newSimpleAppModule(c.upgradeH, c.name, func() {
			m[c.name]++
			count++
		}))
	}

	mm := setupModuleBasics(modules...)

	app := newTestSimApp("demo", log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, func(txBytes []byte, height ...int64) (sdk.Tx, error) {
		return nil, nil
	}, mm)
	genesisState := ModuleBasics.DefaultGenesis()
	stateBytes, err := codec.MarshalJSONIndent(app.Codec(), genesisState)
	require.NoError(t, err)
	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit(abci.RequestCommit{})

	for i := int64(2); i < maxHeight+5; i++ {
		header := abci.Header{Height: i}
		app.BeginBlock(abci.RequestBeginBlock{Header: header})
		app.Commit(abci.RequestCommit{})
	}
	for _, v := range m {
		require.Equal(t, 1, v)
	}
	require.Equal(t, count, len(cases))
}
