package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Route should return the name of the module
func (msg *MsgSetAIOracleReq) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSetAIOracleReq) Type() string { return EventTypeSetAIOracle }

// ValidateBasic runs stateless checks on the message
func (msg *MsgSetAIOracleReq) ValidateBasic() error {
	// if msg.Owner.Empty() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	// }
	if len(msg.Contract) == 0 || msg.ValidatorCount <= 0 {
		return sdkerrors.Wrap(ErrRequestInvalid, "Name or / and validator count cannot be empty")
	}
	_, err := sdk.ParseCoinsNormalized(msg.Fees)
	if err != nil {
		return sdkerrors.Wrap(ErrRequestFeesInvalid, err.Error())
	}
	if len(msg.Fees) == 0 {
		return sdkerrors.Wrap(ErrRequestFeesInvalid, "The fee format is not correct")
	}
	return nil
}

// // ValidateBasic runs stateless checks on the message
// func (msg *MsgCreateReport) ValidateBasic() error {
// 	reporter := msg.GetReporter()
// 	if reporter.GetAddress().Empty() || len(reporter.GetName()) == 0 || !provider.IsStringAlphabetic(reporter.GetName()) || len(reporter.GetName()) >= ReporterNameLen {
// 		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, reporter.String())
// 	} else if len(msg.GetRequestID()) == 0 || reporter.Validator.Empty() {
// 		return sdkerrors.Wrap(ErrMsgReportInvalid, "Request ID / validator address cannot be empty")
// 	} else if len(msg.GetDataSourceResults()) == 0 || len(msg.GetTestCaseResults()) == 0 || len(msg.GetAggregatedResult()) == 0 {
// 		return sdkerrors.Wrap(ErrMsgReportInvalid, "lengths of the data source and test case must be greater than zero, and there must be an aggregated result")
// 	} else if msg.GetResultStatus() != ResultSuccess && msg.GetResultStatus() != ResultFailure {
// 		return sdkerrors.Wrap(ErrMsgReportInvalid, "result status of the report is not valid")
// 	} else {
// 		var dsResultSize int
// 		for _, dsResult := range msg.DataSourceResults {
// 			dsResultSize += len(dsResult.Result)
// 		}
// 		var tcResultSize int
// 		for _, tcResult := range msg.TestCaseResults {
// 			for _, dsResult := range tcResult.DataSourceResults {
// 				tcResultSize += len(dsResult.Result)
// 			}
// 		}
// 		aggregatedResultSize := len(msg.AggregatedResult)
// 		requestIdSize := len(msg.RequestID)
// 		finalLen := dsResultSize + tcResultSize + aggregatedResultSize + requestIdSize
// 		if finalLen >= MsgLen {
// 			return sdkerrors.Wrap(ErrMsgReportInvalid, "Size of the report should not be larger than 200KB")
// 		}

// 		_, err := sdk.ParseCoinsNormalized(msg.Fees.String())
// 		if err != nil {
// 			return sdkerrors.Wrap(ErrReportFeeTypeInvalid, err.Error())
// 		}
// 		return nil
// 	}
// }

// GetSignBytes encodes the message for signing
func (msg *MsgSetAIOracleReq) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgSetAIOracleReq) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}
