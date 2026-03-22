package handlers

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

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
				Port:    25565,
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
			name: "port is zero",
			req: CreateServerRequest{
				Name:    "myserver",
				Type:    models.ServerTypeVanilla,
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  4,
				Port:    0,
			},
			wantErr: true,
		},
		{
			name: "port below minimum",
			req: CreateServerRequest{
				Name:    "myserver",
				Type:    models.ServerTypeVanilla,
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  4,
				Port:    25564,
			},
			wantErr: true,
		},
		{
			name: "port at minimum",
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
			name: "port above maximum",
			req: CreateServerRequest{
				Name:    "myserver",
				Type:    models.ServerTypeVanilla,
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  4,
				Port:    65536,
			},
			wantErr: true,
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

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&models.Server{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})

	return db
}

func TestListServers(t *testing.T) {
	db := setupTestDB(t)
	handler := NewServerHandler(db)

	tests := []struct {
		name           string
		setupDB        func()
		expectedStatus int
	}{
		{
			name: "returns servers successfully",
			setupDB: func() {
				db.Create(&models.Server{Name: "server1", Type: "vanilla", Version: "1.20.4", MinMem: 2, MaxMem: 4, Port: 25565})
				db.Create(&models.Server{Name: "server2", Type: "fabric", Version: "1.20.4", MinMem: 2, MaxMem: 4, Port: 25566})
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "returns empty list when no servers",
			setupDB:        func() {},
			expectedStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.Exec("DELETE FROM servers")
			tt.setupDB()

			app := fiber.New()
			app.Get("/servers", handler.List)

			req := httptest.NewRequest("GET", "/servers", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("failed to test request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestCreateServer(t *testing.T) {
	db := setupTestDB(t)
	handler := NewServerHandler(db)

	tests := []struct {
		name           string
		server         models.Server
		expectedStatus int
	}{
		{
			name: "creates server successfully",
			server: models.Server{
				Name:    "testserver",
				Type:    models.ServerTypeVanilla,
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  4,
				Port:    25565,
			},
			expectedStatus: fiber.StatusCreated,
		},
		{
			name: "returns conflict for duplicate port",
			server: models.Server{
				Name:    "testserver2",
				Type:    models.ServerTypeVanilla,
				Version: "1.20.4",
				MinMem:  2,
				MaxMem:  4,
				Port:    25565,
			},
			expectedStatus: fiber.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := db.Begin()
			t.Cleanup(func() { tx.Rollback() })

			app := fiber.New()
			app.Post("/servers", handler.Create)

			body := `{"name":"` + tt.server.Name + `","type":"` + string(tt.server.Type) + `","version":"` + tt.server.Version + `","min_mem":` + fmt.Sprintf("%d", tt.server.MinMem) + `,"max_mem":` + fmt.Sprintf("%d", tt.server.MaxMem) + `,"port":` + fmt.Sprintf("%d", tt.server.Port) + `}`

			req := httptest.NewRequest("POST", "/servers", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("failed to test request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}
