package rpc

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
)

func Dispatch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !cors(&w, r) {
		return
	}
	var res *types.DispatchResponse
	var err error
	d := types.SessionHeader{}

	if err := PopModel(w, r, ps, &d); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	if res, err = app.QueryDispatch(d); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, err := json.Marshal(res)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

func Relay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var relay = types.Relay{}
	if !cors(&w, r) {
		return
	}
	if err := PopModel(w, r, ps, &relay); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryRelay(relay)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, er := json.Marshal(res)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

func Challenge(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var challenge = types.ChallengeProofInvalidData{}
	if !cors(&w, r) {
		return
	}
	if err := PopModel(w, r, ps, &challenge); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.QueryChallenge(challenge)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, er := json.Marshal(res)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

type sendRawTxParams struct {
	Addr        string `json:"address"`
	RawHexBytes string `json:"raw_hex_bytes"`
}

func SendRawTx(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = sendRawTxParams{}
	if !cors(&w, r) {
		return
	}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	bz, err := hex.DecodeString(params.RawHexBytes)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.SendRawTx(params.Addr, bz)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, er := app.Codec().MarshalJSON(res)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteResponse(w, string(j), r.URL.Path, r.Host)
}
