module github.com/your-github-name/swift-codes-service

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/xuri/excelize/v2 v2.8.1
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.5

	github.com/stretchr/testify v1.8.4
	gorm.io/driver/sqlite v1.5.1
	github.com/testcontainers/testcontainers-go v0.30.0
)

# cofamy bibliotekę wymagającą Go 1.23 → kompatybilna z 1.22
replace github.com/rogpeppe/go-internal v1.14.1 => github.com/rogpeppe/go-internal v1.13.1
