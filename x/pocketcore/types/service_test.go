package types

import (
	"encoding/hex"
	appsType "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"reflect"
	"testing"
	"time"
)

func TestRelay_Validate(t *testing.T) { // TODO add overservice, and not unique relay here
	clientPrivateKey := GetRandomPrivateKey()
	clientPubKey := clientPrivateKey.PublicKey().RawString()
	appPrivateKey := GetRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	npk := getRandomPubKey()
	nodePubKey := npk.RawString()
	InitCacheTest()
	ethereum := hex.EncodeToString([]byte{01})
	bitcoin := hex.EncodeToString([]byte{02})
	p := Payload{
		Data:    "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[],\"id\":67}",
		Method:  "",
		Path:    "",
		Headers: nil,
	}
	validRelay := Relay{
		Payload: p,
		Meta:    RelayMeta{BlockHeight: 1},
		Proof: RelayProof{
			Entropy:            1,
			SessionBlockHeight: 1,
			ServicerPubKey:     nodePubKey,
			Blockchain:         ethereum,
			Token: AAT{
				Version:              "0.0.1",
				ApplicationPublicKey: appPubKey,
				ClientPublicKey:      clientPubKey,
				ApplicationSignature: "",
			},
			Signature: "",
		},
	}
	validRelay.Proof.RequestHash = validRelay.RequestHashString()
	appSig, er := appPrivateKey.Sign(validRelay.Proof.Token.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay.Proof.Token.ApplicationSignature = hex.EncodeToString(appSig)
	clientSig, er := clientPrivateKey.Sign(validRelay.Proof.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay.Proof.Signature = hex.EncodeToString(clientSig)
	// invalid payload empty data and method
	invalidPayloadEmpty := validRelay
	invalidPayloadEmpty.Payload.Data = ""
	selfNode := nodesTypes.Validator{
		Address:                 sdk.Address(npk.Address()),
		PublicKey:               npk,
		Jailed:                  false,
		Status:                  sdk.Staked,
		Chains:                  []string{ethereum, bitcoin},
		ServiceURL:              "https://www.google.com:443",
		StakedTokens:            sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	}
	var allNodes []exported.ValidatorI
	for i := 0; i < 4; i++ {
		pubKey := getRandomPubKey()
		allNodes = append(allNodes, nodesTypes.Validator{
			Address:                 sdk.Address(pubKey.Address()),
			PublicKey:               pubKey,
			Jailed:                  false,
			Status:                  sdk.Staked,
			Chains:                  []string{ethereum, bitcoin},
			ServiceURL:              "https://www.google.com:443",
			StakedTokens:            sdk.NewInt(100000),
			UnstakingCompletionTime: time.Time{},
		})
	}
	var noEthereumNodes []exported.ValidatorI
	for i := 0; i < 4; i++ {
		pubKey := getRandomPubKey()
		noEthereumNodes = append(noEthereumNodes, nodesTypes.Validator{
			Address:                 sdk.Address(pubKey.Address()),
			PublicKey:               pubKey,
			Jailed:                  false,
			Status:                  sdk.Staked,
			Chains:                  []string{bitcoin},
			ServiceURL:              "https://www.google.com:443",
			StakedTokens:            sdk.NewInt(100000),
			UnstakingCompletionTime: time.Time{},
		})
	}
	allNodes = append(allNodes, selfNode)
	noEthereumNodes = append(noEthereumNodes, selfNode)
	hb := HostedBlockchains{
		M: map[string]HostedBlockchain{ethereum: {
			ID:  ethereum,
			URL: "https://www.google.com:443",
		}},
	}
	hbNotSupported := HostedBlockchains{
		M: map[string]HostedBlockchain{bitcoin: {
			ID:  bitcoin,
			URL: "https://www.google.com:443",
		}},
	}
	pubKey := getRandomPubKey()
	app := appsType.Application{
		Address:                 sdk.Address(pubKey.Address()),
		PublicKey:               pubKey,
		Jailed:                  false,
		Status:                  sdk.Staked,
		Chains:                  []string{ethereum},
		StakedTokens:            sdk.NewInt(1000),
		MaxRelays:               sdk.NewInt(1000),
		UnstakingCompletionTime: time.Time{},
	}
	tests := []struct {
		name     string
		relay    Relay
		node     nodesTypes.Validator
		app      appsType.Application
		allNodes []exported.ValidatorI
		hb       HostedBlockchains
		hasError bool
	}{
		{
			name:     "valid relay",
			relay:    validRelay,
			node:     selfNode,
			app:      app,
			allNodes: allNodes,
			hb:       hb,
			hasError: false,
		},
		{
			name:     "invalid relay: payload empty",
			relay:    invalidPayloadEmpty,
			node:     selfNode,
			app:      app,
			allNodes: allNodes,
			hb:       hb,
			hasError: true,
		},
		{
			name:     "invalid relay: unsupported blockchain",
			relay:    validRelay,
			node:     selfNode,
			app:      app,
			allNodes: allNodes,
			hb:       hbNotSupported,
			hasError: true,
		},
		{
			name:     "invalid relay: not enough service nodes",
			relay:    validRelay,
			node:     selfNode,
			app:      app,
			allNodes: noEthereumNodes,
			hb:       hb,
			hasError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			k := MockPosKeeper{Validators: tt.allNodes}
			assert.Equal(t, tt.relay.Validate(newContext(t, false).WithAppVersion("0.0.0"), k, tt.node,
				&tt.hb, 1, 5, tt.app) != nil, tt.hasError)
		})
		ClearSessionCache()
	}
}

func TestRelay_Execute(t *testing.T) {
	clientPrivateKey := GetRandomPrivateKey()
	clientPubKey := clientPrivateKey.PublicKey().RawString()
	appPrivateKey := GetRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	npk := getRandomPubKey()
	nodePubKey := npk.RawString()
	ethereum := hex.EncodeToString([]byte{01})
	p := Payload{
		Data:    "foo",
		Method:  "POST",
		Path:    "",
		Headers: nil,
	}
	validRelay := Relay{
		Payload: p,
		Proof: RelayProof{
			Entropy:            1,
			SessionBlockHeight: 1,
			ServicerPubKey:     nodePubKey,
			Blockchain:         ethereum,
			Token: AAT{
				Version:              "0.0.1",
				ApplicationPublicKey: appPubKey,
				ClientPublicKey:      clientPubKey,
				ApplicationSignature: "",
			},
			Signature: "",
		},
	}
	validRelay.Proof.RequestHash = validRelay.RequestHashString()
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://server.com").
		Post("/relay").
		Reply(200).
		BodyString("bar")

	hb := HostedBlockchains{
		M: map[string]HostedBlockchain{ethereum: {
			ID:  ethereum,
			URL: "https://server.com/relay/",
		}},
	}
	response, err := validRelay.Execute(&hb)
	assert.True(t, err == nil)
	assert.Equal(t, response, "bar")
}

func TestRelay_HandleProof(t *testing.T) {
	clientPrivateKey := GetRandomPrivateKey()
	clientPubKey := clientPrivateKey.PublicKey().RawString()
	appPrivateKey := GetRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	npk := getRandomPubKey()
	nodePubKey := npk.RawString()
	ethereum := hex.EncodeToString([]byte{01})
	p := Payload{
		Data:    "foo",
		Method:  "POST",
		Path:    "",
		Headers: nil,
	}
	validRelay := Relay{
		Payload: p,
		Proof: RelayProof{
			Entropy:            1,
			SessionBlockHeight: 1,
			ServicerPubKey:     nodePubKey,
			Blockchain:         ethereum,
			Token: AAT{
				Version:              "0.0.1",
				ApplicationPublicKey: appPubKey,
				ClientPublicKey:      clientPubKey,
				ApplicationSignature: "",
			},
			Signature: "",
		},
	}
	validRelay.Proof.RequestHash = validRelay.RequestHashString()
	validRelay.Proof.Store()
	res := GetProof(SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}, RelayEvidence, 0)
	assert.True(t, reflect.DeepEqual(validRelay.Proof, res))
}

func TestRelayResponse_BytesAndHash(t *testing.T) {
	nodePrivKey := GetRandomPrivateKey()
	nodePubKey := nodePrivKey.PublicKey().RawString()
	appPrivKey := GetRandomPrivateKey()
	appPublicKey := appPrivKey.PublicKey().RawString()
	cliPrivKey := GetRandomPrivateKey()
	cliPublicKey := cliPrivKey.PublicKey().RawString()
	relayResp := RelayResponse{
		Signature: "",
		Response:  "foo",
		Proof: RelayProof{
			Entropy:            230942034,
			SessionBlockHeight: 1,
			RequestHash:        nodePubKey,
			ServicerPubKey:     nodePubKey,
			Blockchain:         hex.EncodeToString(hash([]byte("foo"))),
			Token: AAT{
				Version:              "0.0.1",
				ApplicationPublicKey: appPublicKey,
				ClientPublicKey:      cliPublicKey,
				ApplicationSignature: "",
			},
			Signature: "",
		},
	}
	appSig, err := appPrivKey.Sign(relayResp.Proof.Token.Hash())
	if err != nil {
		t.Fatalf(err.Error())
	}
	relayResp.Proof.Token.ApplicationSignature = hex.EncodeToString(appSig)
	assert.NotNil(t, relayResp.Hash())
	assert.Equal(t, hex.EncodeToString(relayResp.Hash()), relayResp.HashString())
	storedHashString := relayResp.HashString()
	nodeSig, err := nodePrivKey.Sign(relayResp.Hash())
	if err != nil {
		t.Fatalf(err.Error())
	}
	relayResp.Signature = hex.EncodeToString(nodeSig)
	assert.Equal(t, storedHashString, relayResp.HashString())
}

func TestSortJSON(t *testing.T) {
	// out of order json arrays
	j1 := `{"foo":0,"bar":1}`
	j2 := `{"bar":1,"foo":0}`
	// sort
	objs := sortJSONResponse(j1)
	objs2 := sortJSONResponse(j2)
	// compare
	assert.Equal(t, objs, objs2)
}

type MockPosKeeper struct {
	Validators []exported.ValidatorI
}

func (m MockPosKeeper) RewardForRelays(ctx sdk.Ctx, relays sdk.Int, address sdk.Address) {
	panic("implement me")
}

func (m MockPosKeeper) GetStakedTokens(ctx sdk.Ctx) sdk.Int {
	panic("implement me")
}

func (m MockPosKeeper) Validator(ctx sdk.Ctx, addr sdk.Address) exported.ValidatorI {
	for _, v := range m.Validators {
		if addr.Equals(v.GetAddress()) {
			return v
		}
	}
	return nil
}

func (m MockPosKeeper) TotalTokens(ctx sdk.Ctx) sdk.Int {
	panic("implement me")
}

func (m MockPosKeeper) BurnForChallenge(ctx sdk.Ctx, challenges sdk.Int, address sdk.Address) {
	panic("implement me")
}

func (m MockPosKeeper) JailValidator(ctx sdk.Ctx, addr sdk.Address) {
	panic("implement me")
}

func (m MockPosKeeper) AllValidators(ctx sdk.Ctx) (validators []exported.ValidatorI) {
	return m.Validators
}

func (m MockPosKeeper) GetStakedValidators(ctx sdk.Ctx) (validators []exported.ValidatorI) {
	return m.Validators
}

func (m MockPosKeeper) BlocksPerSession(ctx sdk.Ctx) (res int64) {
	panic("implement me")
}

func (m MockPosKeeper) StakeDenom(ctx sdk.Ctx) (res string) {
	panic("implement me")
}
