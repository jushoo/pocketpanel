package handlers

import (
	"testing"

	"pocketpanel/api/internal/models"
	"pocketpanel/api/internal/validator"
)

func TestCreateServerRequestValidation(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name    string
		req     CreateServerRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateServerRequest{
				Name:    "myserver",
				Type:    models.ServerTypeVanilla,
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  4,
				Port:    25565,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			req: CreateServerRequest{
				Name:    "",
				Type:    models.ServerTypeVanilla,
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  4,
			},
			wantErr: true,
		},
		{
			name: "name too short",
			req: CreateServerRequest{
				Name:    "ab",
				Type:    models.ServerTypeVanilla,
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  4,
			},
			wantErr: true,
		},
		{
			name: "name too long",
			req: CreateServerRequest{
				Name:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Type:    models.ServerTypeVanilla,
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  4,
			},
			wantErr: true,
		},
		{
			name: "missing type",
			req: CreateServerRequest{
				Name:    "myserver",
				Type:    "",
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  4,
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			req: CreateServerRequest{
				Name:    "myserver",
				Type:    "forge",
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  4,
			},
			wantErr: true,
		},
		{
			name: "fabric type is valid",
			req: CreateServerRequest{
				Name:    "myserver",
				Type:    models.ServerTypeFabric,
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  4,
			},
			wantErr: false,
		},
		{
			name: "missing version",
			req: CreateServerRequest{
				Name:    "myserver",
				Type:    models.ServerTypeVanilla,
				Version: "",
				MinMem:  2,
				MaxMem:  4,
			},
			wantErr: true,
		},
		{
			name: "min_mem is zero",
			req: CreateServerRequest{
				Name:    "myserver",
				Type:    models.ServerTypeVanilla,
				Version: "1.20.4",
				MinMem:  0,
				MaxMem:  4,
			},
			wantErr: true,
		},
		{
			name: "min_mem exceeds limit",
			req: CreateServerRequest{
				Name:    "myserver",
				Type:    models.ServerTypeVanilla,
				Version: "1.20.4",
				MinMem:  200,
				MaxMem:  4,
			},
			wantErr: true,
		},
		{
			name: "max_mem is zero",
			req: CreateServerRequest{
				Name:    "myserver",
				Type:    models.ServerTypeVanilla,
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  0,
			},
			wantErr: true,
		},
		{
			name: "max_mem exceeds limit",
			req: CreateServerRequest{
				Name:    "myserver",
				Type:    models.ServerTypeVanilla,
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  200,
			},
			wantErr: true,
		},
		{
			name: "port is optional (zero)",
			req: CreateServerRequest{
				Name:    "myserver",
				Type:    models.ServerTypeVanilla,
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  4,
				Port:    0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.req)

			if tt.wantErr && err == nil {
				t.Errorf("expected validation error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
