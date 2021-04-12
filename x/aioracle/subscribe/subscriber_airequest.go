package subscribe

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/oraichain/orai/x/aioracle/types"
)

const (
	QuoteString = `"'`
	JoinString  = `-`
)

func (subscriber *Subscriber) handleAIOracleLog(queryClient types.QueryClient, ev *sdk.StringEvent) (*types.MsgSetAIOracleRes, error) {
	subscriber.log.Info(":delivery_truck: Processing incoming request event before checking validators")

	attrMap := GetAttributeMap(ev.GetAttributes())

	// Skip if not related to this validator
	validators := attrMap[types.AttributeRequestValidator]
	subscriber.log.Info(":delivery_truck: Validator lists: ", validators)
	currentValidator := sdk.ValAddress(subscriber.cliCtx.GetFromAddress()).String()
	hasMe := false
	for _, validator := range validators {
		subscriber.log.Debug(":delivery_truck: validator: %s", validator)
		if validator == currentValidator {
			hasMe = true
			break
		}
	}

	if !hasMe {
		subscriber.log.Info(":next_track_button: Skip request not related to this validator")
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Skip request not related to this validator")
	}

	subscriber.log.Info(":delivery_truck: Processing incoming request event")

	var requestID, inputStr, expectedOutputStr string
	fmt.Println(requestID)
	var oscriptContractAddr sdk.AccAddress = nil
	var err error = nil

	if val, ok := attrMap[types.AttributeRequestID]; ok {
		requestID = val[0]
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, types.AttributeRequestID)
	}
	if val, ok := attrMap[types.AttributeContract]; ok {
		oscriptContractAddr, err = sdk.AccAddressFromBech32(val[0])
		if err != nil {
			return nil, err
		}
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, types.AttributeContract)
	}
	if val, ok := attrMap[types.AttributeRequestInput]; ok {
		inputStr = val[0]
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, types.AttributeRequestInput)
	}
	if val, ok := attrMap[types.AttributeRequestExpectedOutput]; ok {
		expectedOutputStr = val[0]
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, types.AttributeRequestExpectedOutput)
	}

	// collect ai data sources and test cases from the ai request event.
	ctx := context.Background()

	aiDataSources, err := queryClient.DataSourceEntries(ctx, &types.QueryDataSourceEntriesContract{
		Contract: oscriptContractAddr,
		Request:  &types.EmptyParams{},
	})

	if err != nil {
		return nil, err
	}

	testCases, err := queryClient.DataSourceEntries(ctx, &types.QueryDataSourceEntriesContract{
		Contract: oscriptContractAddr,
		Request:  &types.EmptyParams{},
	})

	if err != nil {
		return nil, err
	}

	// create data source results to store in the report
	var dataSourceResultsTest []*types.DataSourceResult
	var dataSourceResults []*types.DataSourceResult
	var testCaseResults []*types.TestCaseResult
	var resultArr = []string{}

	// we have different test cases, so we need to loop through them
	for _, testCase := range testCases.Data {
		//put the results from the data sources into the test case to verify if they are good enough
		for _, aiDataSource := range aiDataSources.Data {
			// collect test case result from the script
			outTestCase, err := queryClient.TestCaseContract(ctx, &types.QueryTestCaseContract{
				Contract: oscriptContractAddr,
				Request: &types.RequestTestCase{
					Tcase:  testCase,
					Input:  inputStr,
					Output: expectedOutputStr,
				},
			})

			dataSourceResult := types.NewDataSourceResult(aiDataSource, []byte{}, types.ResultSuccess)
			// By default, we consider returning null as failure. If any datasource does not follow this rule then it should not be used by any oracle scripts. We dont return error here since the error relates to the external scripts, not the node.
			if err != nil {
				subscriber.log.Info(":skull: failed to execute test case, due to error: %s", err.Error())
				// change status to fail so the datasource cannot be rewarded afterwards
				dataSourceResult.Status = types.ResultFailure
				dataSourceResult.Result = []byte(types.FailedResponseTc)
			} else {
				// remove all quotes at start and begin
				result := strings.Trim(string(outTestCase.Data), QuoteString)
				fmt.Println("result after running test case: ", result)
				dataSourceResult.Result = []byte(result)
			}
			// append an data source result into the list
			dataSourceResultsTest = append(dataSourceResultsTest, dataSourceResult)
		}

		// add test case result
		testCaseResults = append(testCaseResults, types.NewTestCaseResult(testCase, dataSourceResultsTest, types.ResultSuccess))
	}

	// after passing the test cases, we run the actual data sources to collect their results
	// create data source results to store in the report
	// we use dataSourceResultsTest since this list is the complete list of data sources that have passed the test cases
	for _, dataSourceResultTest := range dataSourceResultsTest {
		// run the data source script

		var dataSourceResult *types.DataSourceResult
		if dataSourceResultTest.GetStatus() == types.ResultSuccess {
			outDataSource, err := queryClient.DataSourceContract(ctx, &types.QueryDataSourceContract{
				Contract: oscriptContractAddr,
				Request: &types.RequestDataSource{
					Dsource: dataSourceResultTest.GetEntryPoint(),
					Input:   inputStr,
				},
			})
			// By default, we consider returning null as failure. If any datasource does not follow this rule then it should not be used by any oracle scripts.
			dataSourceResult = types.NewDataSourceResult(dataSourceResultTest.GetEntryPoint(), []byte{}, types.ResultSuccess)
			if err != nil {
				subscriber.log.Error(":skull: failed to execute data source script: %s", err.Error())
				// change status to fail so the datasource cannot be rewarded afterwards
				dataSourceResult.Status = types.ResultFailure
				dataSourceResult.Result = []byte(types.FailedResponseDs)
			} else {
				// remove all quote at start and begin
				result := strings.Trim(string(outDataSource.Data), QuoteString)
				if len(outDataSource.Data) == 0 || result == "null" {
					// change status to fail so the datasource cannot be rewarded afterwards
					dataSourceResult.Status = types.ResultFailure
					dataSourceResult.Result = []byte(types.FailedResponseDs)
				} else {
					dataSourceResult.Result = []byte(result)
					resultArr = append(resultArr, result)
				}
			}
		} else {
			dataSourceResult = types.NewDataSourceResult(dataSourceResultTest.GetEntryPoint(), []byte(dataSourceResultTest.GetResult()), types.ResultFailure)
		}
		// append an data source result into the list
		dataSourceResults = append(dataSourceResults, dataSourceResult)
	}
	subscriber.log.Info("star: final result: ", resultArr)
	// Create a new MsgCreateReport with a new reporter to the Oraichain
	_ = types.NewReporter(
		subscriber.cliCtx.GetFromAddress(), subscriber.cliCtx.GetFromName(),
		sdk.ValAddress(subscriber.cliCtx.GetFromAddress()),
	)
	final := []byte(strings.Join(resultArr, JoinString))
	fmt.Println("final result: ", final)
	// msgReport := types.NewMsgCreateReport(
	// 	requestID, dataSourceResults,
	// 	testCaseResults, reporter,
	// 	subscriber.config.Fees, finalResult,
	// 	types.ResultSuccess,
	// )

	// if len(resultArr) == 0 {
	// 	msgReport.AggregatedResult = []byte(types.FailedResponseOs)
	// 	msgReport.ResultStatus = types.ResultFailure
	// } else {
	// 	query := &types.QueryOracleScriptContract{
	// 		Contract: oscriptContractAddr,
	// 		Request: &types.RequestOracleScript{
	// 			Results: resultArr,
	// 		},
	// 	}
	// 	fmt.Printf("query :%v\n", query)
	// 	outOScript, err := queryClient.OracleScriptContract(ctx, query)

	// 	if err != nil {
	// 		subscriber.log.Error(":skull: failed to aggregate results: %s", err.Error())
	// 		// do not return error since this is the error from the script side
	// 		msgReport.AggregatedResult = []byte(types.FailedResponseOs + ": " + err.Error())
	// 		msgReport.ResultStatus = types.ResultFailure
	// 		return msgReport, nil
	// 	}
	// 	// remove all quote at start and begin
	// 	result := strings.Trim(string(outOScript.Data), QuoteString)
	// 	subscriber.log.Info("final result from oScript: %s", result)
	// 	msgReport.AggregatedResult = []byte(result)
	// }

	// TODO: check aggregated result value of the oracle script to verify if it fails or success
	return nil, nil

}
