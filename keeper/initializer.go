package keeper

import (
	"maps"

	"cosmossdk.io/log"
	"cosmossdk.io/store/metrics"
	"github.com/cosmos/cosmos-sdk/runtime"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"

	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/skip-mev/chaintestutil/sample"
)

var moduleAccountPerms = map[string][]string{
	authtypes.FeeCollectorName:     nil,
	distrtypes.ModuleName:          nil,
	minttypes.ModuleName:           {authtypes.Minter},
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
}

// Initializer allows initializing of each module keeper.
type Initializer struct {
	Codec      codec.Codec
	Amino      *codec.LegacyAmino
	DB         *dbm.MemDB
	StateStore store.CommitMultiStore
	Logger     log.Logger
}

func newInitializer() Initializer {
	logger := log.NewNopLogger()
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db, logger, metrics.NewNoOpMetrics())

	return Initializer{
		DB:         db,
		Codec:      sample.Codec(),
		StateStore: cms,
		Logger:     logger,
	}
}

// ModuleAccountAddrs returns all the app's module account addresses.
func ModuleAccountAddrs(maccPerms map[string][]string) map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (i *Initializer) Auth(maccPerms map[string][]string) authkeeper.AccountKeeper {
	storeKey := storetypes.NewKVStoreKey(authtypes.StoreKey)
	i.StateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, i.DB)
	kvStoreService := runtime.NewKVStoreService(storeKey)

	maps.Copy(moduleAccountPerms, maccPerms)

	return authkeeper.NewAccountKeeper(
		i.Codec,
		kvStoreService,
		authtypes.ProtoBaseAccount,
		moduleAccountPerms,
		authcodec.NewBech32Codec(sdk.Bech32MainPrefix),
		sdk.Bech32PrefixAccAddr,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
}

func (i *Initializer) Bank(authKeeper authkeeper.AccountKeeper, maccPerms map[string][]string) bankkeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(banktypes.StoreKey)
	i.StateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, i.DB)
	kvStoreService := runtime.NewKVStoreService(storeKey)

	maps.Copy(moduleAccountPerms, maccPerms)
	modAccAddrs := ModuleAccountAddrs(moduleAccountPerms)

	return bankkeeper.NewBaseKeeper(
		i.Codec,
		kvStoreService,
		authKeeper,
		modAccAddrs,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		i.Logger,
	)
}

// create mock ProtocolVersionSetter for UpgradeKeeper

type ProtocolVersionSetter struct{}

func (vs ProtocolVersionSetter) SetProtocolVersion(uint64) {}

func (i *Initializer) Upgrade() *upgradekeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(upgradetypes.StoreKey)
	i.StateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, i.DB)
	kvStoreService := runtime.NewKVStoreService(storeKey)

	skipUpgradeHeights := make(map[int64]bool)
	vs := ProtocolVersionSetter{}

	return upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		kvStoreService,
		i.Codec,
		"",
		vs,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
}

func (i *Initializer) Staking(
	authKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
) *stakingkeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(stakingtypes.StoreKey)
	i.StateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, i.DB)
	kvStoreService := runtime.NewKVStoreService(storeKey)

	return stakingkeeper.NewKeeper(
		i.Codec,
		kvStoreService,
		authKeeper,
		bankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		authcodec.NewBech32Codec(sdk.Bech32PrefixValAddr),
		authcodec.NewBech32Codec(sdk.Bech32PrefixConsAddr),
	)
}

func (i *Initializer) Distribution(
	authKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
) distrkeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(distrtypes.StoreKey)
	i.StateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, i.DB)
	kvStoreService := runtime.NewKVStoreService(storeKey)

	return distrkeeper.NewKeeper(
		i.Codec,
		kvStoreService,
		authKeeper,
		bankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
}

func (i *Initializer) FeeGrant(
	authKeeper authkeeper.AccountKeeper,
) feegrantkeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(feegrant.StoreKey)
	i.StateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, i.DB)
	kvStoreService := runtime.NewKVStoreService(storeKey)

	return feegrantkeeper.NewKeeper(
		i.Codec,
		kvStoreService,
		authKeeper,
	)
}

func (i *Initializer) LoadLatest() error {
	return i.StateStore.LoadLatestVersion()
}
