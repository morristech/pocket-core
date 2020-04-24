package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"

	appexported "github.com/pokt-network/pocket-core/x/apps/exported"
	nodeexported "github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

const DEFAULTHTTPMETHOD = "POST"

var (
	globalClientBlockAllowance int
	globalSortJSONResponses    bool
)

// "Relay" - A read / write API request from a hosted (non native) external blockchain
type Relay struct {
	Payload Payload    `json:"payload"` // the data payload of the request
	Meta    RelayMeta  `json:"meta"`    // metadata for the relay request
	Proof   RelayProof `json:"proof"`   // the authentication scheme needed for work
}

// "Validate" - Checks the validity of a relay request using store data
func (r *Relay) Validate(ctx sdk.Ctx, keeper PosKeeper, node nodeexported.ValidatorI, hb *HostedBlockchains, sessionBlockHeight int64,
	sessionNodeCount int, app appexported.ApplicationI) sdk.Error {
	// validate payload
	if err := r.Payload.Validate(); err != nil {
		return NewEmptyPayloadDataError(ModuleName)
	}
	// validate the metadata
	if err := r.Meta.Validate(ctx); err != nil {
		return err
	}
	// validate the relay hash = request hash
	if r.Proof.RequestHash != r.RequestHashString() {
		return NewRequestHashError(ModuleName)
	}
	// ensure the blockchain is supported locally
	if !hb.Contains(r.Proof.Blockchain) {
		return NewUnsupportedBlockchainNodeError(ModuleName)
	}
	evidenceHeader := SessionHeader{
		ApplicationPubKey:  r.Proof.Token.ApplicationPublicKey,
		Chain:              r.Proof.Blockchain,
		SessionBlockHeight: r.Proof.SessionBlockHeight,
	}
	// validate unique relay
	totalRelays := GetTotalProofs(evidenceHeader, RelayEvidence)
	// get evidence key by proof
	if !IsUniqueProof(evidenceHeader, r.Proof) {
		return NewDuplicateProofError(ModuleName)
	}
	// validate not over service
	if totalRelays >= int64(math.Ceil(float64(app.GetMaxRelays().Int64())/float64(len(app.GetChains())))/(float64(sessionNodeCount))) {
		return NewOverServiceError(ModuleName)
	}
	// validate the Proof
	if err := r.Proof.ValidateLocal(app.GetChains(), sessionNodeCount, sessionBlockHeight, node.GetPublicKey().RawString()); err != nil {
		return err
	}
	// get the sessionContext
	sessionContext, er := ctx.PrevCtx(sessionBlockHeight)
	if er != nil {
		return sdk.ErrInternal(er.Error())
	}
	// generate the header
	header := SessionHeader{
		ApplicationPubKey:  app.GetPublicKey().RawString(),
		Chain:              r.Proof.Blockchain,
		SessionBlockHeight: sessionBlockHeight,
	}
	// check cache
	session, found := GetSession(header)
	// if not found generate the session
	if !found {
		var err sdk.Error
		session, err = NewSession(sessionContext, ctx, keeper, header, BlockHash(sessionContext), sessionNodeCount)
		if err != nil {
			return err
		}
		// add to cache
		SetSession(session)
	}
	// validate the session
	err := session.Validate(node, app, sessionNodeCount)
	if err != nil {
		return err
	}
	// if the payload method is empty, set it to the default
	if r.Payload.Method == "" {
		r.Payload.Method = DEFAULTHTTPMETHOD
	}
	return nil
}

// "Execute" - Attempts to do a request on the non-native blockchain specified
func (r Relay) Execute(hostedBlockchains *HostedBlockchains) (string, sdk.Error) {
	// retrieve the hosted blockchain url requested
	url, err := hostedBlockchains.GetChainURL(r.Proof.Blockchain)
	if err != nil {
		return "", err
	}
	url = strings.Trim(url, `/`) + "/" + strings.Trim(r.Payload.Path, `/`)
	// do basic http request on the relay
	res, er := executeHTTPRequest(r.Payload.Data, url, r.Payload.Method, r.Payload.Headers)
	if er != nil {
		return res, NewHTTPExecutionError(ModuleName, er)
	}
	return res, nil
}

// "Requesthash" - The cryptographic hash representation of the request
func (r Relay) RequestHash() []byte {
	relay := struct {
		Payload Payload   `json:"payload"` // the data payload of the request
		Meta    RelayMeta `json:"meta"`    // metadata for the relay request
	}{r.Payload, r.Meta}
	res, err := json.Marshal(relay)
	if err != nil {
		panic(fmt.Sprintf("cannot marshal relay request hash: %s", err.Error()))
	}
	return res
}

// "RequestHashString" - The hex string representation of the request hash
func (r Relay) RequestHashString() string {
	return hex.EncodeToString(r.RequestHash())
}

// "Payload" - A data being sent to the non-native chain
type Payload struct {
	Data    string            `json:"data"`              // the actual data string for the external chain
	Method  string            `json:"method"`            // the http CRUD method
	Path    string            `json:"path"`              // the REST Path
	Headers map[string]string `json:"headers,omitempty"` // http headers
}

// "Bytes" - The bytes reprentation of a payload object
func (p Payload) Bytes() []byte {
	bz, err := json.Marshal(p)
	if err != nil {
		panic(fmt.Sprintf("an error occured converting the payload to bytes:\n%v", err))
	}
	return bz
}

// "Hash" - The cryptographic hash representation of the payload object
func (p Payload) Hash() []byte {
	return Hash(p.Bytes())
}

// "HashString" - The hex encoded string representation of the payload object
func (p Payload) HashString() string {
	return hex.EncodeToString(p.Hash())
}

// "Validate" - Validity check for the payload object
func (p Payload) Validate() sdk.Error {
	if p.Data == "" && p.Path == "" {
		return NewEmptyPayloadDataError(ModuleName)
	}
	return nil
}

// "payload" - A structure used for custom json marshalling/unmarshalling
type payload struct {
	Data    string            `json:"data"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
}

// "MarshalJSON" - Overrides json marshalling
func (p Payload) MarshalJSON() ([]byte, error) {
	pay := payload{
		Data:    p.Data,
		Method:  p.Method,
		Path:    p.Path,
		Headers: p.Headers,
	}
	return json.Marshal(pay)
}

// "RelayMeta" - Metadata that is included in the relay request
type RelayMeta struct {
	BlockHeight int64 `json:"block_height"` // the block height when the request is made
}

// "Validate" - Validates the relay meta object
func (m RelayMeta) Validate(ctx sdk.Ctx) sdk.Error {
	// ensures the block height is within the acceptable range
	if ctx.BlockHeight()+int64(globalClientBlockAllowance) < m.BlockHeight || ctx.BlockHeight()-int64(globalClientBlockAllowance) > m.BlockHeight {
		return NewOutOfSyncRequestError(ModuleName)
	}
	return nil
}

func InitClientBlockAllowance(allowance int) {
	globalClientBlockAllowance = allowance
}

// response structure for the relay
type RelayResponse struct {
	Signature string     `json:"signature"` // signature from the node in hex
	Response  string     `json:"payload"`   // response to relay
	Proof     RelayProof `json:"proof"`     // to be signed by the client
}

// "Validate" - The node validates the response after signing
func (rr RelayResponse) Validate() sdk.Error {
	// cannot contain empty response
	if rr.Response == "" {
		return NewEmptyResponseError(ModuleName)
	}
	// cannot contain empty signature (nodes must be accountable)
	if rr.Signature == "" || len(rr.Signature) == crypto.Ed25519SignatureSize {
		return NewResponseSignatureError(ModuleName)
	}
	return nil
}

// "Hash" - The cryptographic hash representation of the relay response
func (rr RelayResponse) Hash() []byte {
	seed, err := json.Marshal(relayResponse{
		Signature: "",
		Response:  rr.Response,
		Proof:     rr.Proof.HashString(),
	})
	if err != nil {
		panic(fmt.Sprintf("an error occured hashing the relay response:\n%v", err))
	}
	return Hash(seed)
}

// "HashString" - The hex string representation of the hash
func (rr RelayResponse) HashString() string {
	return hex.EncodeToString(rr.Hash())
}

// "relayResponse" - a structure used for custom json
type relayResponse struct {
	Signature string `json:"signature"`
	Response  string `json:"payload"`
	Proof     string `json:"Proof"`
}

// "ChallengeReponse" - The response object used in challenges
type ChallengeResponse struct {
	Response string `json:"response"`
}

// "DispatchResponse" - The response object used in dispatching
type DispatchResponse struct {
	Session     Session `json:"session"`
	BlockHeight int64   `json:"block_height"`
}

// "executeHTTPRequest" takes in the raw json string and forwards it to the RPC endpoint
func executeHTTPRequest(payload string, url string, method string, headers map[string]string) (string, error) {
	// generate an http request
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}
	// add headers if needed
	if len(headers) == 0 {
		req.Header.Set("Content-Type", "application/json")
	} else {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	// execute the request
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// ensure code is 200
	if resp.StatusCode != 200 {
		return "", NewHTTPStatusCodeError(ModuleName, resp.StatusCode)
	}
	// read all bz
	body, _ := ioutil.ReadAll(resp.Body)
	if globalSortJSONResponses {
		body = []byte(sortJSONResponse(string(body)))
	}
	// return
	return string(body), nil
}

func InitJSONSorting(doSorting bool) {
	globalSortJSONResponses = doSorting
}

// "sortJSONResponse" - sorts json from a relay response
func sortJSONResponse(response string) string {
	var rawJSON map[string]interface{}
	// unmarshal into json
	if err := json.Unmarshal([]byte(response), &rawJSON); err != nil {
		return response
	}
	// marshal into json
	bz, err := json.Marshal(rawJSON)
	if err != nil {
		return response
	}
	return string(bz)
}
