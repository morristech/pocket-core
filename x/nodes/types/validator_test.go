package types

import (
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestNewValidator(t *testing.T) {
	type args struct {
		addr          sdk.Address
		consPubKey    crypto.PublicKey
		tokensToStake sdk.Int
		chains        []string
		serviceURL    string
	}
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name string
		args args
		want Validator
	}{
		{"defaultValidator", args{sdk.Address(pub.Address()), pub, sdk.ZeroInt(), []string{"00"}, "https://www.google.com:443"},
			Validator{
				Address:                 sdk.Address(pub.Address()),
				PublicKey:               pub,
				Jailed:                  false,
				Status:                  sdk.Staked,
				StakedTokens:            sdk.ZeroInt(),
				Chains:                  []string{"00"},
				ServiceURL:              "https://www.google.com:443",
				UnstakingCompletionTime: time.Unix(0, 0).UTC(), // zero out because status: staked
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewValidator(tt.args.addr, tt.args.consPubKey, tt.args.chains, tt.args.serviceURL, tt.args.tokensToStake); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_ABCIValidatorUpdate(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   abci.ValidatorUpdate
	}{
		{"testingABCIValidatorUpdate Unstaked", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unstaked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, abci.ValidatorUpdate{
			PubKey: tmtypes.TM2PB.PubKey(pub.PubKey()),
			Power:  int64(0),
		}},
		{"testingABCIValidatorUpdate Staked", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, abci.ValidatorUpdate{
			PubKey: tmtypes.TM2PB.PubKey(pub.PubKey()),
			Power:  int64(0),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.ABCIValidatorUpdate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ABCIValidatorUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_AddStakedTokens(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	type args struct {
		tokens sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Validator
	}{
		{"Default Add Token Test", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{tokens: sdk.NewInt(100)},
			Validator{
				Address:                 sdk.Address(pub.Address()),
				PublicKey:               pub,
				Jailed:                  false,
				Status:                  sdk.Staked,
				StakedTokens:            sdk.NewInt(100),
				UnstakingCompletionTime: time.Time{},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.AddStakedTokens(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddStakedTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_ConsAddress(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.Address
	}{
		{"Default Test", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.Address(pub.Address())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetAddress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_ConsensusPower(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{"Default Test / 0 power", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, 0},
		{"Default Test / 1 power", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.NewInt(1000000),
			UnstakingCompletionTime: time.Time{},
		}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.ConsensusPower(); got != tt.want {
				t.Errorf("ConsensusPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_Equals(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	type args struct {
		v2 Validator
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"Default Test Equal", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{Validator{
			Address:                 sdk.Address(pub.Address()),
			PublicKey:               pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}}, true},
		{"Default Test Not Equal", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{Validator{
			Address:                 sdk.Address(pub.Address()),
			PublicKey:               pub,
			Jailed:                  false,
			Status:                  sdk.Unstaked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.Equals(tt.args.v2); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetAddress(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.Address
	}{
		{"Default Test", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.Address(pub.Address())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetAddress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetConsAddr(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.Address
	}{
		{"Default Test", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.Address(pub.Address())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetAddress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetConsPubKey(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}
	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   crypto.PublicKey
	}{
		{"Default Test", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, pub},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetPublicKey(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetStatus(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.StakeStatus
	}{
		{"Default Test", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.Staked},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetStatus(); got != tt.want {
				t.Errorf("GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetTokens(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   sdk.Int
	}{
		{"Default Test", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, sdk.ZeroInt()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetTokens(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_IsJailed(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Default Test", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.IsJailed(); got != tt.want {
				t.Errorf("IsJailed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_IsStaked(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Default Test / staked true", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, true},
		{"Default Test / Unstaking false", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unstaking,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
		{"Default Test / Unstaked false", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unstaked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.IsStaked(); got != tt.want {
				t.Errorf("IsStaked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_IsUnstaked(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Default Test / staked false", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
		{"Default Test / Unstaking false", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unstaking,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
		{"Default Test / Unstaked true", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unstaked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.IsUnstaked(); got != tt.want {
				t.Errorf("IsUnstaked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_IsUnstaking(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Default Test / staked false", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
		{"Default Test / Unstaking true", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unstaking,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, true},
		{"Default Test / Unstaked false", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unstaked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.IsUnstaking(); got != tt.want {
				t.Errorf("IsUnstaking() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_PotentialConsensusPower(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{"Default Test / potential power 0", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.ConsensusPower(); got != tt.want {
				t.Errorf("ConsensusPower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_RemoveStakedTokens(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	type args struct {
		tokens sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Validator
	}{
		{"Remove 0 tokens having 0 tokens ", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{tokens: sdk.ZeroInt()}, Validator{
			Address:                 sdk.Address(pub.Address()),
			PublicKey:               pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}},
		{"Remove 99 tokens having 100 tokens ", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.NewInt(100),
			UnstakingCompletionTime: time.Time{},
		}, args{tokens: sdk.NewInt(99)}, Validator{
			Address:                 sdk.Address(pub.Address()),
			PublicKey:               pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.OneInt(),
			UnstakingCompletionTime: time.Time{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.RemoveStakedTokens(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveStakedTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_UpdateStatus(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	type args struct {
		newStatus sdk.StakeStatus
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Validator
	}{
		{"Test staked -> Unstaking", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{newStatus: sdk.Unstaking}, Validator{
			Address:                 sdk.Address(pub.Address()),
			PublicKey:               pub,
			Jailed:                  false,
			Status:                  sdk.Unstaking,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}},
		{"Test Unstaking -> Unstaked", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unstaking,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{newStatus: sdk.Unstaked}, Validator{
			Address:                 sdk.Address(pub.Address()),
			PublicKey:               pub,
			Jailed:                  false,
			Status:                  sdk.Unstaked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}},
		{"Test Unstaked -> staked", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Unstaked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, args{newStatus: sdk.Staked}, Validator{
			Address:                 sdk.Address(pub.Address()),
			PublicKey:               pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.UpdateStatus(tt.args.newStatus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetServiceURL(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		Chains                  []string
		ServiceURL              string
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Test Service URL", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			Chains:                  []string{"00"},
			ServiceURL:              "www.pokt.network",
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, "www.pokt.network"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				Chains:                  tt.fields.Chains,
				ServiceURL:              tt.fields.ServiceURL,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetServiceURL(); got != tt.want {
				t.Errorf("GetServiceURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_GetChains(t *testing.T) {
	type fields struct {
		Address                 sdk.Address
		ConsPubKey              crypto.PublicKey
		Jailed                  bool
		Status                  sdk.StakeStatus
		Chains                  []string
		ServiceURL              string
		StakedTokens            sdk.Int
		UnstakingCompletionTime time.Time
	}

	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])

	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"Test Service URL", fields{
			Address:                 sdk.Address(pub.Address()),
			ConsPubKey:              pub,
			Jailed:                  false,
			Status:                  sdk.Staked,
			Chains:                  []string{"00"},
			ServiceURL:              "www.pokt.network",
			StakedTokens:            sdk.ZeroInt(),
			UnstakingCompletionTime: time.Time{},
		}, []string{"00"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Validator{
				Address:                 tt.fields.Address,
				PublicKey:               tt.fields.ConsPubKey,
				Jailed:                  tt.fields.Jailed,
				Status:                  tt.fields.Status,
				Chains:                  tt.fields.Chains,
				ServiceURL:              tt.fields.ServiceURL,
				StakedTokens:            tt.fields.StakedTokens,
				UnstakingCompletionTime: tt.fields.UnstakingCompletionTime,
			}
			if got := v.GetChains(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChains() = %v, want %v", got, tt.want)
			}
		})
	}
}
