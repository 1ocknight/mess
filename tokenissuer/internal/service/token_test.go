package service

import (
	"context"
	"reflect"
	"testing"
	"tokenissuer/internal/adapter/identifier"
)

func TestNewTokenImpl(t *testing.T) {
	type args struct {
		iden identifier.Service
	}
	tests := []struct {
		name string
		args args
		want *TokenImpl
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTokenImpl(tt.args.iden); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenImpl_RefreshTokenPair(t *testing.T) {
	type fields struct {
		iden identifier.Service
	}
	type args struct {
		ctx          context.Context
		refreshToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    identifier.TokenPair
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := &TokenImpl{
				iden: tt.fields.iden,
			}
			got, err := td.RefreshTokenPair(tt.args.ctx, tt.args.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Fatalf("TokenImpl.RefreshTokenPair() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenImpl.RefreshTokenPair() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenImpl_GetTokenPair(t *testing.T) {
	type fields struct {
		iden identifier.Service
	}
	type args struct {
		ctx         context.Context
		code        string
		redirectURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    identifier.TokenPair
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := &TokenImpl{
				iden: tt.fields.iden,
			}
			got, err := td.GetTokenPair(tt.args.ctx, tt.args.code, tt.args.redirectURL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("TokenImpl.GetTokenPair() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenImpl.GetTokenPair() = %v, want %v", got, tt.want)
			}
		})
	}
}
