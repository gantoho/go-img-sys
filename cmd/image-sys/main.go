package main

import (
	"flag"
	"os"

	"github.com/gantoho/go-img-sys/internal/app"
	"github.com/gantoho/go-img-sys/internal/config"
)

// @title           Go Image System API
// @version         1.0.0
// @description     Go Image System API
// @description    ## Authentication
// @description    - **API Key**: Pass `X-API-Key` header or `api_key` query parameter
// @description    - **JWT**: Login at `/api/auth/login` to get token, then use `Authorization: Bearer <token>`
// @description    ## Development API Keys
// @description    - `demo-key-12345` (30 days validity)
// @description    - `test-key-67890` (7 days validity)
// @contact.name   GitHub Repository
// @contact.url    https://github.com/gantoho/go-img-sys
// @license.name   MIT
// @license.url    https://opensource.org/licenses/MIT
// @host           localhost:3128
// @BasePath       /
// @schemes        http
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        X-API-Key
// @securityDefinitions.apikey  JWTAuth
// @in                          header
// @name                        Authorization
// @description                Bearer <token>

//go:generate swag init -g cmd/image-sys/main.go -o docs --parseDependency --parseInternal
func main() {
	port := flag.String("port", "", "server port (e.g. :3128)")
	env := flag.String("env", "", "runtime environment (development/release)")
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	cfg := config.Init()

	if *port != "" {
		cfg.Server.Port = *port
	}
	if *env != "" {
		cfg.Server.Env = *env
	}
	_ = configPath

	// Override from environment variables (highest priority)
	if v := os.Getenv("SERVER_PORT"); v != "" {
		cfg.Server.Port = v
	}
	if v := os.Getenv("SERVER_ENV"); v != "" {
		cfg.Server.Env = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.Auth.JWTSecret = v
	}
	if v := os.Getenv("UPLOAD_DIR"); v != "" {
		cfg.File.UploadDir = v
	}
	if v := os.Getenv("MAX_FILE_SIZE"); v != "" {
		if size, err := flagParseInt64(v); err == nil {
			cfg.File.MaxSize = size
		}
	}

	srv := app.New(cfg)
	defer srv.Close()
	srv.Start()
}

func flagParseInt64(s string) (int64, error) {
	var v int64
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, nil
		}
		v = v*10 + int64(c-'0')
	}
	return v, nil
}
