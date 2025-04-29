package test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/your-github-name/swift-codes-service/internal/db"
	"github.com/your-github-name/swift-codes-service/internal/models"
	"github.com/your-github-name/swift-codes-service/internal/repository"
	"github.com/your-github-name/swift-codes-service/internal/transport/http"
)

func TestEndToEnd(t *testing.T) {
	ctx := context.Background()
	pg, _ := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER": "swift", "POSTGRES_PASSWORD": "swift", "POSTGRES_DB": "swift",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	defer pg.Terminate(ctx)

	host, _ := pg.Host(ctx)
	port, _ := pg.MappedPort(ctx, "5432/tcp")

	t.Setenv("PGHOST", host)
	t.Setenv("PGPORT", port.Port())

	// prepare DB
	conn, _ := db.Connect()
	_ = conn.AutoMigrate(&models.SwiftCode{})
	repo := repository.New(conn)
	_ = repo.UpsertMany([]models.SwiftCode{
		{SwiftCode: "AAISALTRXXX", CountryISO2: "AL", CountryName: "ALBANIA", BankName: "UBA", IsHeadquarter: true},
	})

	srv := http.Router(http.New(repo))
	go srv.Run(":8080")
	time.Sleep(time.Second) // minimal wait for Gin to start

	res, err := http.Get("http://localhost:8080/v1/swift-codes/AAISALTRXXX")
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	var data map[string]any
	_ = json.Unmarshal(body, &data)
	require.Equal(t, "ALBANIA", data["countryName"])
}
