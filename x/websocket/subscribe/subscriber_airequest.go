package subscribe

import (
	"context"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	aiRequest "github.com/oraichain/orai/x/airequest"
	"github.com/oraichain/orai/x/provider"

	"time"

	"github.com/oraichain/orai/x/websocket/types"
)

// SubmitReport creates a new MsgCreateReport and submits it to the Oraichain to create a new report
func (subscriber *Subscriber) submitReport(msgReport *types.MsgCreateReport) {

	if err := msgReport.ValidateBasic(); err != nil {
		subscriber.log.Error(":exploding_head: Failed to validate basic with error: %s", err.Error())
		return
	}

	txf := subscriber.newTxFactory("websocket")
	for try := uint64(1); try <= subscriber.config.MaxTry; try++ {
		subscriber.log.Info(":e-mail: Try to broadcast report transaction(%d/%d)", try, subscriber.config.MaxTry)
		err := tx.BroadcastTx(*subscriber.cliCtx, txf, msgReport)
		if err == nil {
			break
		}
		subscriber.log.Info(":warning: Failed to broadcast tx with error: %s", err.Error())
		time.Sleep(subscriber.config.RPCPollInterval)
	}
}

func (subscriber *Subscriber) handleAIRequestLog(queryClient types.QueryClient, attrMap map[string][]string) {
	subscriber.log.Info(":delivery_truck: Processing incoming request event before checking validators")

	// Skip if not related to this validator
	validators := attrMap[aiRequest.AttributeRequestValidator]
	hasMe := false
	for _, validator := range validators {
		subscriber.log.Debug(":delivery_truck: validator: %s", validator)
		if validator == subscriber.cliCtx.GetFromAddress().String() {
			hasMe = true
			break
		}
	}

	if !hasMe {
		subscriber.log.Info(":next_track_button: Skip request not related to this validator")
		return
	}

	subscriber.log.Info(":delivery_truck: Processing incoming request event")

	requestID := attrMap[aiRequest.AttributeRequestID][0]

	oscriptName := attrMap[aiRequest.AttributeOracleScriptName][0]
	creatorStr := attrMap[aiRequest.AttributeRequestCreator][0]
	// valCountStr := attrMap[aiRequest.AttributeRequestValidatorCount][0]
	inputStr := attrMap[aiRequest.AttributeRequestInput][0]
	expectedOutputStr := attrMap[aiRequest.AttributeRequestExpectedOutput][0]

	creator, _ := sdk.AccAddressFromBech32(creatorStr)

	// collect ai data sources and test cases from the ai request event.
	aiDataSources := attrMap[aiRequest.AttributeRequestDSources]
	testCases := attrMap[aiRequest.AttributeRequestTCases]

	var finalResultStr string
	// create data source results to store in the report
	var dataSourceResultsTest []*types.DataSourceResult
	var dataSourceResults []*types.DataSourceResult
	var testCaseResults []*types.TestCaseResult

	// we have different test cases, so we need to loop through them
	ctx := context.Background()
	for i := range testCases {
		//put the results from the data sources into the test case to verify if they are good enough
		for j := range aiDataSources {
			//// collect test case result from the script
			// outTestCase, err := ExecPythonFile("python", getTCasePath(testCases[i]), []string{provider.DataSourceStoreKeyString(aiDataSources[j]), req.Input, req.ExpectedOutput})
			query := &types.QueryTestCaseContract{
				Name:          testCases[i],
				DataSourcName: aiDataSources[j],
				Request: &types.RequestTestCase{
					Input:  inputStr,
					Output: expectedOutputStr,
				},
			}
			outTestCase, err := queryClient.TestCaseContract(ctx, query)

			if err != nil {
				subscriber.log.Info(":skull: failed to execute test case 1st loop: %s", err.Error())
			}
			result := string(outTestCase.Data)

			subscriber.log.Info("result after running test case: ", result)

			dataSourceResult := types.NewDataSourceResult(aiDataSources[j], []byte(result), types.ResultSuccess)

			// By default, we consider returning null as failure. If any datasource does not follow this rule then it should not be used by any oracle scripts.
			if result == types.FailedResult || len(result) == 0 {
				// change status to fail so the datasource cannot be rewarded afterwards
				dataSourceResult.Status = types.ResultFailure
				dataSourceResult.Result = []byte(types.FailedResponseTc)
			}
			// append an data source result into the list
			dataSourceResultsTest = append(dataSourceResultsTest, dataSourceResult)
		}

		// add test case result
		testCaseResults = append(testCaseResults, types.NewTestCaseResult(testCases[i], dataSourceResultsTest))
	}

	// after passing the test cases, we run the actual data sources to collect their results
	// create data source results to store in the report
	// we use dataSourceResultsTest since this list is the complete list of data sources that have passed the test cases
	for i := range dataSourceResultsTest {
		// run the data source script

		var dataSourceResult *types.DataSourceResult
		if dataSourceResultsTest[i].GetStatus() == types.ResultSuccess {

			// outDataSource, err = ExecPythonFile("python", getDSourcePath(dataSourceResultsTest[i].GetName()), []string{})
			query := &types.QueryDataSourceContract{
				Name: dataSourceResultsTest[i].GetName(),
			}
			outDataSource, err := queryClient.DataSourceContract(ctx, query)
			if err != nil {
				subscriber.log.Error(":skull: failed to execute data source script: %s", err.Error())
			}
			result := string(outDataSource.Data)

			subscriber.log.Info("star: result from data sources: ", result)
			// By default, we consider returning null as failure. If any datasource does not follow this rule then it should not be used by any oracle scripts.
			dataSourceResult = types.NewDataSourceResult(dataSourceResultsTest[i].GetName(), []byte(result), types.ResultSuccess)
			if result == types.FailedResult || len(result) == 0 {
				// change status to fail so the datasource cannot be rewarded afterwards
				dataSourceResult.Status = types.ResultFailure
				dataSourceResult.Result = []byte(types.FailedResponseDs)
			} else {
				//resultArr = append(resultArr, resultFloat)
				finalResultStr = finalResultStr + result
			}
		} else {
			dataSourceResult = types.NewDataSourceResult(dataSourceResultsTest[i].GetName(), []byte(dataSourceResultsTest[i].GetResult()), types.ResultFailure)
		}
		// append an data source result into the list
		dataSourceResults = append(dataSourceResults, dataSourceResult)
	}
	finalResultStr = strings.TrimSuffix(finalResultStr, "-")
	subscriber.log.Info("star: final result after trimming: ", finalResultStr)
	// Create a new MsgCreateReport with a new reporter to the Oraichain
	reporter := types.NewReporter(creator, subscriber.config.FromValidator, sdk.ValAddress(subscriber.cliCtx.FromAddress))
	msgReport := types.NewMsgCreateReport(requestID, dataSourceResults, testCaseResults, reporter, sdk.NewCoins(sdk.NewCoin(provider.Denom, sdk.NewInt(int64(5000)))), []byte(finalResultStr), types.ResultSuccess)
	if len(finalResultStr) == 0 {
		msgReport.AggregatedResult = []byte(types.FailedResponseOs)
		msgReport.ResultStatus = types.ResultFailure
		// Create a new MsgCreateReport to the Oraichain
	} else {
		// ExecPythonFile("python", oscriptPath, []string{"aggregation", finalResultStr})

		query := &types.QueryOracleScriptContract{
			Name: oscriptName,
			Request: &types.RequestOracleScript{
				Aggregation: finalResultStr,
			},
		}
		outOScript, err := queryClient.OracleScriptContract(ctx, query)
		if err != nil {
			subscriber.log.Error(":skull: failed to aggregate results: %s", err.Error())
		}
		subscriber.log.Info("final result from oScript: %s", string(outOScript.Data))
		msgReport.AggregatedResult = outOScript.Data
	}

	// TODO: check aggregated result value of the oracle script to verify if it fails or success
	subscriber.submitReport(msgReport)

}