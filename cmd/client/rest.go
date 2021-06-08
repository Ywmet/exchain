package client

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bankrest "github.com/cosmos/cosmos-sdk/x/bank/client/rest"
	supplyrest "github.com/cosmos/cosmos-sdk/x/supply/client/rest"
	"github.com/okex/exchain/app/rpc"
	"github.com/okex/exchain/app/types"
	ammswaprest "github.com/okex/exchain/x/ammswap/client/rest"
	backendrest "github.com/okex/exchain/x/backend/client/rest"
	dexclient "github.com/okex/exchain/x/dex/client"
	dexrest "github.com/okex/exchain/x/dex/client/rest"
	dist "github.com/okex/exchain/x/distribution"
	distr "github.com/okex/exchain/x/distribution"
	distrest "github.com/okex/exchain/x/distribution/client/rest"
	evmclient "github.com/okex/exchain/x/evm/client"
	evmrest "github.com/okex/exchain/x/evm/client/rest"
	farmclient "github.com/okex/exchain/x/farm/client"
	farmrest "github.com/okex/exchain/x/farm/client/rest"
	govrest "github.com/okex/exchain/x/gov/client/rest"
	orderrest "github.com/okex/exchain/x/order/client/rest"
	paramsclient "github.com/okex/exchain/x/params/client"
	stakingrest "github.com/okex/exchain/x/staking/client/rest"
	"github.com/okex/exchain/x/stream/nacos"
	"github.com/okex/exchain/x/token"
	tokensrest "github.com/okex/exchain/x/token/client/rest"
	"github.com/spf13/viper"
	"net/http"
	_ "net/http/pprof"
)

// RegisterRoutes registers the routes from the different modules for the LCD.
// NOTE: details on the routes added for each module are in the module documentation
// NOTE: If making updates here you also need to update the test helper in client/lcd/test_helper.go
func RegisterRoutes(rs *lcd.RestServer) {
	if !rs.CliCtx.Client.IsRunning() {
		err := rs.CliCtx.Client.Start()
		if err != nil {
			panic(err)
		}
	}
	rpc.RegisterRoutes(rs)
	pathPrefix := viper.GetString(server.FlagRestPathPrefix)
	if pathPrefix == "" {
		pathPrefix = types.EthBech32Prefix
	}
	registerRoutesV1(rs, pathPrefix)
	registerRoutesV2(rs, pathPrefix)

	// start nacos client for registering rest-server service
	if viper.GetBool(rpc.FlagRestServerServiceEnable) && viper.GetString(rpc.FlagRestServerNacosUrls) != "" {
		nacos.StartNacosClient(rs.Logger(), viper.GetString(rpc.FlagRestServerNacosUrls), viper.GetString(rpc.FlagRestServerNacosNamespaceId), viper.GetString(rpc.FlagRestServerNacosServiceName), viper.GetString(server.FlagExternalListenAddr))
	}
	if viper.GetString(rpc.FlagRestServerProfLaddr) != "" {
		go func() {
			err := http.ListenAndServe(viper.GetString(rpc.FlagRestServerProfLaddr), nil)
			if err != nil {
				panic(err)
			}
		}()
	}
}

func registerRoutesV1(rs *lcd.RestServer, pathPrefix string) {
	v1Router := rs.Mux.PathPrefix(fmt.Sprintf("/%s/v1", pathPrefix)).Name("v1").Subrouter()
	client.RegisterRoutes(rs.CliCtx, v1Router)
	authrest.RegisterRoutes(rs.CliCtx, v1Router, auth.StoreKey)
	bankrest.RegisterRoutes(rs.CliCtx, v1Router)
	stakingrest.RegisterRoutes(rs.CliCtx, v1Router)
	distrest.RegisterRoutes(rs.CliCtx, v1Router, dist.StoreKey)

	orderrest.RegisterRoutes(rs.CliCtx, v1Router)
	tokensrest.RegisterRoutes(rs.CliCtx, v1Router, token.StoreKey)
	backendrest.RegisterRoutes(rs.CliCtx, v1Router)
	dexrest.RegisterRoutes(rs.CliCtx, v1Router)
	ammswaprest.RegisterRoutes(rs.CliCtx, v1Router)
	supplyrest.RegisterRoutes(rs.CliCtx, v1Router)
	farmrest.RegisterRoutes(rs.CliCtx, v1Router)
	evmrest.RegisterRoutes(rs.CliCtx, v1Router)
	govrest.RegisterRoutes(rs.CliCtx, v1Router,
		[]govrest.ProposalRESTHandler{
			paramsclient.ProposalHandler.RESTHandler(rs.CliCtx),
			distr.ProposalHandler.RESTHandler(rs.CliCtx),
			dexclient.DelistProposalHandler.RESTHandler(rs.CliCtx),
			farmclient.ManageWhiteListProposalHandler.RESTHandler(rs.CliCtx),
			evmclient.ManageContractDeploymentWhitelistProposalHandler.RESTHandler(rs.CliCtx),
		},
	)
}

func registerRoutesV2(rs *lcd.RestServer, pathPrefix string) {
	v2Router := rs.Mux.PathPrefix(fmt.Sprintf("/%s/v2", pathPrefix)).Name("v1").Subrouter()
	client.RegisterRoutes(rs.CliCtx, v2Router)
	authrest.RegisterRoutes(rs.CliCtx, v2Router, auth.StoreKey)
	bankrest.RegisterRoutes(rs.CliCtx, v2Router)
	stakingrest.RegisterRoutes(rs.CliCtx, v2Router)
	distrest.RegisterRoutes(rs.CliCtx, v2Router, dist.StoreKey)

	orderrest.RegisterRoutesV2(rs.CliCtx, v2Router)
	tokensrest.RegisterRoutesV2(rs.CliCtx, v2Router, token.StoreKey)
	backendrest.RegisterRoutesV2(rs.CliCtx, v2Router)
}
