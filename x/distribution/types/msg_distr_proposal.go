//nolint
package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/libs/tendermint/types"
)

// Verify interface at compile time
var _ = &MsgWithdrawDelegatorReward{}

// msg struct for delegation withdraw from a single validator
type MsgWithdrawDelegatorReward struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
}

func NewMsgWithdrawDelegatorReward(delAddr sdk.AccAddress, valAddr sdk.ValAddress) MsgWithdrawDelegatorReward {
	return MsgWithdrawDelegatorReward{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
	}
}

func (msg MsgWithdrawDelegatorReward) Route() string { return ModuleName }
func (msg MsgWithdrawDelegatorReward) Type() string  { return "withdraw_delegator_reward" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawDelegatorReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddress)}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawDelegatorReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawDelegatorReward) ValidateBasic() error {
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr()
	}
	if msg.ValidatorAddress.Empty() {
		return ErrNilValidatorAddr()
	}

	//will delete it after upgrade venus2
	if !types.HigherThanVenus2(global.GetGlobalHeight()) {
		return ErrCodeNotSupportWithdrawDelegationRewards()
	}

	return nil
}
