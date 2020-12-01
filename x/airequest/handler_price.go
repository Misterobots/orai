package airequest

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/oraichain/orai/x/airequest/keeper"
	"github.com/oraichain/orai/x/airequest/types"
)

func handleMsgSetPriceRequest(ctx sdk.Context, k keeper.Keeper, msg types.MsgSetPriceRequest) (*sdk.Result, error) {
	validators, err := k.RandomValidators(ctx, msg.MsgAIRequest.ValidatorCount, []byte(msg.MsgAIRequest.RequestID))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrCannotRandomValidators, err.Error())
	}

	// we can safely parse fees to coins since we have validated it in the Msg already
	fees, _ := sdk.ParseCoins(msg.MsgAIRequest.Fees)
	// Compute the fee allocated for oracle module to distribute to active validators.
	rewardRatio := sdk.NewDecWithPrec(int64(k.GetParam(ctx, types.KeyOracleScriptRewardPercentage)), 2)
	// We need to calculate the final 70% fee given by the user because the remaining 30% must be reserved for the proposer and validators.
	providedCoins, _ := sdk.NewDecCoinsFromCoins(fees...).MulDecTruncate(rewardRatio).TruncateDecimal()

	// get data source and test case names from the oracle script
	aiDataSources, testCases, err := k.ProviderKeeper.GetDNamesTcNames(ctx, msg.MsgAIRequest.OracleScriptName)
	if err != nil {
		return nil, err
	}

	fmt.Println("ai data sources: ", aiDataSources)

	// collect data source and test case objects to store into the request
	dataSourceObjs, testcaseObjs, err := getDSourcesTCases(ctx, k, aiDataSources, testCases)
	if err != nil {
		return nil, err
	}

	fmt.Println("data source obj: ", dataSourceObjs)

	finalFees, err := k.ProviderKeeper.GetMinimumFees(ctx, aiDataSources, testCases, len(validators))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Error getting minimum fees from oracle script")
	}
	fmt.Println("final fees needed: ", finalFees.String())

	// If the total fee is larger than the fee provided by the user then we return error
	if finalFees.IsAnyGT(providedCoins) {
		return nil, sdkerrors.Wrap(types.ErrNeedMoreFees, "Fees given by the users are less than the total fees needed")
	}
	// set a new request with the aggregated result into blockchain
	request := types.NewAIRequest(msg.MsgAIRequest.RequestID, msg.MsgAIRequest.OracleScriptName, msg.MsgAIRequest.Creator, validators, ctx.BlockHeight(), dataSourceObjs, testcaseObjs, fees, msg.MsgAIRequest.Input, msg.MsgAIRequest.ExpectedOutput)

	//fmt.Printf("request result: %s\n", outOracleScript.String())

	fmt.Println("request information: ", request)

	k.SetAIRequest(ctx, request.RequestID, request)

	// TODO: Define your msg events
	// Emit an event describing a data request and asked validators.
	event := sdk.NewEvent(types.EventTypeSetPriceRequest)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeRequestID, string(request.RequestID[:])),
	)
	for _, validator := range validators {
		event = event.AppendAttributes(
			sdk.NewAttribute(types.AttributeRequestValidator, validator.String()),
		)
	}
	ctx.EventManager().EmitEvent(event)

	// TODO: Define your msg events
	// Emit an event describing a data request and asked validators.
	event = sdk.NewEvent(types.EventTypeRequestWithData)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeRequestID, string(request.RequestID[:])),
		sdk.NewAttribute(types.AttributeOracleScriptName, request.OracleScriptName),
		sdk.NewAttribute(types.AttributeRequestCreator, msg.MsgAIRequest.Creator.String()),
		sdk.NewAttribute(types.AttributeRequestValidatorCount, fmt.Sprint(msg.MsgAIRequest.ValidatorCount)),
		sdk.NewAttribute(types.AttributeRequestInput, string(msg.MsgAIRequest.Input)),
		sdk.NewAttribute(types.AttributeRequestExpectedOutput, string(msg.MsgAIRequest.ExpectedOutput)),
		sdk.NewAttribute(types.AttributeRequestDSources, strings.Join(aiDataSources, "-")),
		sdk.NewAttribute(types.AttributeRequestTCases, strings.Join(testCases, "-")),
	)
	ctx.EventManager().EmitEvent(event)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
