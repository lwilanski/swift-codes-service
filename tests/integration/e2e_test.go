//go:build integration
// +build integration

package integration_test

import (
	"context"
	"encoding/json"
	"io"
	stdhttp "net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/lwilanski/swift-codes-service/internal/db"
	"github.com/lwilanski/swift-codes-service/internal/models"
	"github.com/lwilanski/swift-codes-service/internal/repository"
	th "github.com/lwilanski/swift-codes-service/internal/transport/http"
)

func TestEndToEnd(t *testing.T) {
	ctx := context.Background()
	pg, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
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
	require.NoError(t, err)
	defer pg.Terminate(ctx)

	host, _ := pg.Host(ctx)
	port, _ := pg.MappedPort(ctx, "5432/tcp")

	t.Setenv("PGHOST", host)
	t.Setenv("PGPORT", port.Port())

	// DB + seed
	conn, _ := db.Connect()
	_ = conn.AutoMigrate(&models.SwiftCode{})
	r := repository.New(conn)
	_ = r.UpsertMany([]models.SwiftCode{
		{SwiftCode: "AAISALTRXXX", CountryISO2: "AL", CountryName: "ALBANIA", BankName: "UBA", IsHeadquarter: true},
	})

	// start API
	srv := th.Router(th.New(r))
	go srv.Run(":8080")
	time.Sleep(time.Second) // krótki wait

	res, err := stdhttp.Get("http://localhost:8080/v1/swift-codes/AAISALTRXXX")
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	var data map[string]any
	_ = json.Unmarshal(body, &data)
	require.Equal(t, "ALBANIA", data["countryName"])
}
