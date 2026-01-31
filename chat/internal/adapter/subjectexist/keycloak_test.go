package subjectexist_test

import (
	"testing"
	"time"

	"github.com/1ocknight/mess/chat/internal/adapter/subjectexist"
)

var Config = subjectexist.Config{
	KeycloakURL:  "http://localhost:7070",
	Realm:        "main",
	ClientID:     "user-checker-service",
	ClientSecret: "user-checker-secret",
	Timeout:      1 * time.Second,
}

func TestKeycloak_SubjectExists(t *testing.T) {
	tests := []struct {
		name      string
		cfg       subjectexist.Config
		subjectID string
		want      bool
		wantErr   bool
	}{
		{
			name:      "existing user",
			cfg:       Config,
			subjectID: "e1bafe11-fff5-4700-bc7a-fc6248509882",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "non-existing user",
			cfg:       Config,
			subjectID: "00000000-0000-0000-0000-000000000000",
			want:      false,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := subjectexist.New(tt.cfg)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := k.Exists(t.Context(), tt.subjectID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SubjectExists() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SubjectExists() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("SubjectExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
