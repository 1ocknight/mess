package workers_test

import(
	"context"
	"github.com/TATAROmangol/mess/profile/internal/storage"
	"github.com/TATAROmangol/mess/profile/internal/wokers"
	"github.com/TATAROmangol/mess/shared/messagequeue"
	"testing"
)

func TestProfileDelete(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cons    messagequeue.Consumer
		store   storage.Profile
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := workers.ProfileDelete[](context.Background(), tt.cons, tt.store)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ProfileDelete() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ProfileDelete() succeeded unexpectedly")
			}
		})
	}
}
