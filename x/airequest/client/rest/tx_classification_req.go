package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/oraichain/orai/x/airequest/types"
	"github.com/segmentio/ksuid"
)

type setClassificationRequestReq struct {
	FormAIRequest setFormAIRequestReq `json:"base_form_request"`
	Name          string              `json:"Name"`
	Hash          string              `json:"Hash"`
}

// newSetClassificationRequestReq is the constructor for the setClassificationRequestReq
func newSetClassificationRequestReq(formAIReq setFormAIRequestReq, name, hash string) setClassificationRequestReq {
	return setClassificationRequestReq{
		FormAIRequest: formAIReq,
		Name:          name,
		Hash:          hash,
	}
}

func (kyc setClassificationRequestReq) getFormAIRequest() setFormAIRequestReq {
	return kyc.FormAIRequest
}

func (kyc setClassificationRequestReq) getName() string {
	return kyc.Name
}

func (kyc setClassificationRequestReq) getHash() string {
	return kyc.Hash
}

// setClassificationRequestReqFn is the function that collects all the necessary info of image classification and return a new object out of it
func setClassificationRequestReqFn(cliCtx context.CLIContext, w http.ResponseWriter, r *http.Request) setClassificationRequestReq {
	req := setAIRequestHandlerFn(cliCtx, w, r)
	imageHash := r.FormValue("image_hash")
	imageName := r.FormValue("image_name")
	return newSetClassificationRequestReq(req, imageName, imageHash)
}

func setClassificationRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		req := setClassificationRequestReqFn(cliCtx, w, r)

		// Need to create a baseReq to write tx response. We cannot use baseReq in the AIRequest struct because AIRequest needs to be in form data to be able to send images
		baseReq := rest.BaseReq{
			From:          req.getFormAIRequest().From,
			Memo:          req.getFormAIRequest().Memo,
			ChainID:       req.getFormAIRequest().ChainID,
			AccountNumber: req.getFormAIRequest().AccountNumber,
			Sequence:      req.getFormAIRequest().Sequence,
			Fees:          req.getFormAIRequest().Fees,
			GasPrices:     req.getFormAIRequest().GasPrices,
			Gas:           req.getFormAIRequest().Gas,
			GasAdjustment: req.getFormAIRequest().GasAdjustment,
			Simulate:      req.getFormAIRequest().Simulate,
		}

		if !baseReq.ValidateBasic(w) {
			return
		}

		// collect valid address from the request address string
		addr, err := sdk.AccAddressFromBech32(req.getFormAIRequest().From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "AVXSD")
			return
		}

		// create the message
		msg := types.NewMsgSetClassificationRequest(req.getHash(), req.getName(), types.NewMsgSetAIRequest(ksuid.New().String(), req.getFormAIRequest().OracleScriptName, addr, req.getFormAIRequest().Fees.String(), req.getFormAIRequest().ValidatorCount, req.getFormAIRequest().Input, req.getFormAIRequest().ExpectedOutput))
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "GHYK")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
