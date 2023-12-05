// Package sample provides methods to initialize sample object of various types for test purposes
package sample

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cosmosed25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/require"
)

// ExtraRegistries is a function passed to sample.InterfaceRegistry() where any added interface registries can be
// registered.
type ExtraRegistries func(codectypes.InterfaceRegistry)

func InterfaceRegistry(register ExtraRegistries) codectypes.InterfaceRegistry {
	interfaceRegistry := codectypes.NewInterfaceRegistry()

	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	authz.RegisterInterfaces(interfaceRegistry)
	stakingtypes.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	consensusparamtypes.RegisterInterfaces(interfaceRegistry)
	slashingtypes.RegisterInterfaces(interfaceRegistry)
	upgradetypes.RegisterInterfaces(interfaceRegistry)
	distrtypes.RegisterInterfaces(interfaceRegistry)
	vestingtypes.RegisterInterfaces(interfaceRegistry)
	feegrant.RegisterInterfaces(interfaceRegistry)
	group.RegisterInterfaces(interfaceRegistry)
	govtypes.RegisterInterfaces(interfaceRegistry)
	evidencetypes.RegisterInterfaces(interfaceRegistry)
	crisistypes.RegisterInterfaces(interfaceRegistry)
	minttypes.RegisterInterfaces(interfaceRegistry)

	// call extra registry function
	register(interfaceRegistry)

	return interfaceRegistry
}

// Codec returns a codec with preregistered interfaces
func Codec(register ExtraRegistries) codec.Codec {
	return codec.NewProtoCodec(InterfaceRegistry(register))
}

// Bool returns randomly true or false
func Bool(r *rand.Rand) bool {
	b := r.Intn(100)
	return b < 50
}

// Bytes returns a random array of bytes
func Bytes(r *rand.Rand, n int) []byte {
	return []byte(String(r, n))
}

// Uint64 returns a random uint64
func Uint64(r *rand.Rand) uint64 {
	return uint64(r.Intn(10000))
}

// String returns a random string of length n
func String(r *rand.Rand, n int) string {
	letter := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	randomString := make([]rune, n)
	for i := range randomString {
		randomString[i] = letter[r.Intn(len(letter))]
	}
	return string(randomString)
}

// AlphaString returns a random string with lowercase alpha char of length n
func AlphaString(r *rand.Rand, n int) string {
	letter := []rune("abcdefghijklmnopqrstuvwxyz")

	randomString := make([]rune, n)
	for i := range randomString {
		randomString[i] = letter[r.Intn(len(letter))]
	}
	return string(randomString)
}

// NonAlphaString returns a random string with non alpha char of length n
func NonAlphaString(r *rand.Rand, n int) string {
	letter := []rune("0123456789!@#$%^&*()_+")

	randomString := make([]rune, n)
	for i := range randomString {
		randomString[i] = letter[r.Intn(len(letter))]
	}
	return string(randomString)
}

// PubKey returns a sample account PubKey
func PubKey(r *rand.Rand) crypto.PubKey {
	seed := []byte(strconv.Itoa(r.Int()))
	return ed25519.GenPrivKeyFromSecret(seed).PubKey()
}

// ConsAddress returns a sample consensus address
func ConsAddress(r *rand.Rand) sdk.ConsAddress {
	return sdk.ConsAddress(PubKey(r).Address())
}

// AccAddress returns a sample account address
func AccAddress(r *rand.Rand) sdk.AccAddress {
	addr := PubKey(r).Address()
	return sdk.AccAddress(addr)
}

// Address returns a sample string account address
func Address(r *rand.Rand) string {
	return AccAddress(r).String()
}

// ValAddress returns a sample validator operator address
func ValAddress(r *rand.Rand) sdk.ValAddress {
	return sdk.ValAddress(PubKey(r).Address())
}

// OperatorAddress returns a sample string validator operator address
func OperatorAddress(r *rand.Rand) string {
	return ValAddress(r).String()
}

// Validator returns a sample staking validator
func Validator(t testing.TB, r *rand.Rand) stakingtypes.Validator {
	seed := []byte(strconv.Itoa(r.Int()))
	val, err := stakingtypes.NewValidator(
		ValAddress(r),
		cosmosed25519.GenPrivKeyFromSecret(seed).PubKey(),
		stakingtypes.Description{})
	require.NoError(t, err)
	return val
}

// Delegation returns staking delegation with the given address
func Delegation(t testing.TB, r *rand.Rand, addr string) stakingtypes.Delegation {
	delAcc, err := sdk.AccAddressFromBech32(addr)
	require.NoError(t, err)

	return stakingtypes.NewDelegation(
		delAcc,
		ValAddress(r),
		sdk.NewDec(int64(r.Intn(10000))),
	)
}

// Coin returns a sample coin structure
func Coin(r *rand.Rand) sdk.Coin {
	return sdk.NewCoin(AlphaString(r, 5), sdkmath.NewInt(r.Int63n(10000)+1))
}

// CoinWithRange returns a sample coin structure where the amount is a random number between provided min and max values
// with a random denom
func CoinWithRange(r *rand.Rand, min, max int64) sdk.Coin {
	return sdk.NewCoin(AlphaString(r, 5), sdkmath.NewInt(r.Int63n(max-min)+min))
}

// CoinWithRangeAmount returns a sample coin structure where the amount is a random number between provided min and max values
// with a given denom
func CoinWithRangeAmount(r *rand.Rand, denom string, min, max int64) sdk.Coin {
	return sdk.NewCoin(denom, sdkmath.NewInt(r.Int63n(max-min)+min))
}

// Coins returns a sample coins structure
func Coins(r *rand.Rand) sdk.Coins {
	return sdk.NewCoins(Coin(r), Coin(r), Coin(r))
}

// CoinsWithRange returns a sample coins structure where the amount is a random number between provided min and max values
func CoinsWithRange(r *rand.Rand, min, max int64) sdk.Coins {
	return sdk.NewCoins(CoinWithRange(r, min, max), CoinWithRange(r, min, max), CoinWithRange(r, min, max))
}

// CoinsWithRangeAmount returns a sample coins structure where the amount is a random number between provided min and max values
// with a set of given denoms
func CoinsWithRangeAmount(r *rand.Rand, denom1, denom2, denom3 string, min, max int64) sdk.Coins {
	return sdk.NewCoins(CoinWithRangeAmount(r, denom1, min, max), CoinWithRangeAmount(r, denom2, min, max), CoinWithRangeAmount(r, denom3, min, max))
}

// Duration returns a sample time.Duration between a second and 21 days
func Duration(r *rand.Rand) time.Duration {
	return time.Duration(r.Int63n(int64(time.Hour*24*21-time.Second))) + time.Second
}

// DurationFromRange returns a sample time.Duration between the min and max values provided
func DurationFromRange(r *rand.Rand, min, max time.Duration) time.Duration {
	return time.Duration(r.Int63n(int64(max-min))) + min
}

// Int returns a sample sdkmath.Int
func Int(r *rand.Rand) sdkmath.Int {
	return sdkmath.NewInt(r.Int63())
}

// Time returns a sample time
func Time(r *rand.Rand) time.Time {
	return time.UnixMilli(r.Int63n(1000) + 1).UTC()
}

// ZeroTime returns time.Time that represents 0
func ZeroTime() time.Time {
	return time.UnixMilli(0).UTC()
}

// Rand returns a sample Rand object for randomness
func Rand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().Unix()))
}
