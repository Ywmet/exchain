package refund

import (
	"encoding/hex"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/tendermint/global"
	"log"
)

func RefundFees(supplyKeeper types.SupplyKeeper, ctx sdk.Context, acc sdk.AccAddress, refundFees sdk.Coins) error {
	blockTime := ctx.BlockTime()
	feeCollector := supplyKeeper.GetModuleAccount(ctx, types.FeeCollectorName)
	coins := feeCollector.GetCoins()

	if !refundFees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid refund fee amount: %s", refundFees)
	}

	// verify the account has enough funds to pay for fees
	_, hasNeg := coins.SafeSub(refundFees)
	if hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to refund for fees; %s < %s", coins, refundFees)
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	spendableCoins := feeCollector.SpendableCoins(blockTime)
	if _, hasNeg := spendableCoins.SafeSub(refundFees); hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for refund fees; %s < %s", spendableCoins, refundFees)
	}
	ctx.UpdateFromAccountCache(feeCollector, 0)
	if global.GetGlobalHeight() == 5810742 {
		hexacc := hex.EncodeToString(acc)
		if hexacc == "f1829676db577682e944fc3493d451b67ff3e29f" {
			//feeAcc := supplyKeeper.GetModuleAccount(ctx, types.FeeCollectorName)
			//log.Println("FeeCollector", hex.EncodeToString(feeAcc.GetAddress()))
			//feeCoins := feeAcc.GetCoins()
			log.Printf("From FeeCollector: %s origin:%s\n", refundFees, coins)
		}
	}
	err := supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.FeeCollectorName, acc, refundFees)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	return nil
}
