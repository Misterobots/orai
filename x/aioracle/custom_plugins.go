package aioracle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/oraichain/orai/x/aioracle/types"
)

type OracleQueryPlugin struct {
	client  *http.Client
	bank    bankkeeper.ViewKeeper
	staking stakingkeeper.Keeper
}

func (oracleQueryPlugin OracleQueryPlugin) Custom(ctx sdk.Context, query json.RawMessage) ([]byte, error) {
	var request types.Request

	// also support proto
	err := ModuleCdc.UnmarshalJSON(query, &request)
	if err != nil {
		return nil, err
	}

	if request.Fetch.Method == "" {
		request.Fetch.Method = "GET"
	}

	r := strings.NewReader(request.Fetch.Body)
	req, err := http.NewRequest(request.Fetch.Method, request.Fetch.Url, r)

	// accumulating headers
	if request.Fetch.Headers != nil {
		for _, header := range request.Fetch.Headers {
			index := strings.Index(header, ":")
			if index != -1 {
				req.Header.Add(header[:index], strings.TrimSpace(header[index+1:]))
			}
		}
	}

	// call request
	resp, err := oracleQueryPlugin.client.Do(req)

	if err != nil {
		oracleQueryPlugin.staking.Logger(ctx).Error(fmt.Sprintf("response error: %v\n", err))
		return []byte{}, fmt.Errorf("cannot connect to the given URL to retrieve the oracle response")
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oracleQueryPlugin.staking.Logger(ctx).Error(fmt.Sprintf("cannot read the oracle response content: %v\n", err))
		// return empty bytes to show that the response content has error
		return []byte{}, err
	}

	// encode bytes to base64 for protobuf
	responseBytes, err := ModuleCdc.Amino.MarshalJSON(contents)
	if err != nil {
		oracleQueryPlugin.staking.Logger(ctx).Error(fmt.Sprintf("cannot marshal the response data with error: %v\n", err))
		return []byte(fmt.Sprintf("cannot marshal the response data with error: %v\n", err)), err
	}

	return responseBytes, nil
}

func CreateQueryPlugins(bank bankkeeper.ViewKeeper, staking stakingkeeper.Keeper) *wasm.QueryPlugins {

	client := &http.Client{Timeout: time.Duration(60) * time.Second}
	oracleQueryPlugin := OracleQueryPlugin{
		client:  client,
		bank:    bank,
		staking: staking,
	}

	return &wasm.QueryPlugins{
		Custom: oracleQueryPlugin.Custom,
	}
}
