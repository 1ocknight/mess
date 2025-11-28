package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"reflect"
	"sync"
	"testing"
	"time"
	"tokenissuer/internal/adapter/identifier"
	identifiermocks "tokenissuer/internal/adapter/identifier/mocks"
	"tokenissuer/internal/model"
	"tokenissuer/pkg/jwks"
	jwksmocks "tokenissuer/pkg/jwks/mocks"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
)

const (
	KID = "test_kid"
)

func initTestKeyAndToken(t *testing.T, kid string, claims jwt.MapClaims) (*rsa.PrivateKey, *rsa.PublicKey, string) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("cannot generate rsa: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid

	signedToken, err := token.SignedString(priv)
	if err != nil {
		t.Fatalf("cannot sign token: %v", err)
	}

	return priv, &priv.PublicKey, signedToken
}

func TestVerifyImpl_findKeyByKid(t *testing.T) {
	type fields struct {
		iden        identifier.JWKSLoader
		jwks        map[string]jwks.JWKS
		jwksUpdated time.Time
		jwksTTL     time.Duration
		parser      *jwt.Parser
		mu          *sync.RWMutex
	}
	type args struct {
		ctx context.Context
		kid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *rsa.PublicKey
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VerifyImpl{
				iden:        tt.fields.iden,
				jwks:        tt.fields.jwks,
				jwksUpdated: tt.fields.jwksUpdated,
				jwksTTL:     tt.fields.jwksTTL,
				parser:      tt.fields.parser,
				mu:          tt.fields.mu,
			}
			got, err := v.findKeyByKid(tt.args.ctx, tt.args.kid)
			if (err != nil) != tt.wantErr {
				t.Fatalf("VerifyImpl.findKeyByKid() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyImpl.findKeyByKid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyImpl_VerifyToken(t *testing.T) {
	claims := jwt.MapClaims{
		"sub":                "123",
		"email":              "test@example.com",
		"preferred_username": "tester",
	}

	_, key, token := initTestKeyAndToken(t, KID, claims)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwksMock := jwksmocks.NewMockJWKS(ctrl)
	jwksMock.EXPECT().GetPublicKey().Return(key, nil).AnyTimes()

	jwksLoader := identifiermocks.NewMockJWKSLoader(ctrl)
	jwksLoader.EXPECT().LoadJWKS(gomock.Any()).Return(
		map[string]jwks.JWKS{
			KID: jwksMock,
		},
		nil,
	).AnyTimes()

	verify := NewVerifyImpl(jwksLoader, time.Hour)

	type args struct {
		ctx         context.Context
		typeToken   string
		accessToken string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "valid token",
			args: args{
				ctx:         context.Background(),
				typeToken:   BearerType,
				accessToken: token,
			},
			want: &model.User{
				ID:    "123",
				Name:  "tester",
				Email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid token",
			args: args{
				ctx:         context.Background(),
				typeToken:   BearerType,
				accessToken: "234",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid token type",
			args: args{
				ctx:         context.Background(),
				typeToken:   "InvalidType",
				accessToken: token,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong kid",
			args: args{
				ctx:       context.Background(),
				typeToken: BearerType,
				accessToken: func() string {
					_, _, tok := initTestKeyAndToken(t, "wrong_kid", claims)
					return tok
				}(),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := verify.VerifyToken(tt.args.ctx, tt.args.typeToken, tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Fatalf("VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
