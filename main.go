package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/urfave/cli/v3"
)

//go:embed .env
var envTemplate embed.FS

const version = "1.0.0"

type VPNConfig struct {
	Username   string `json:"username"`
	Secret     string `json:"secret"` // Encrypted, base64
	OvpnConfig string `json:"ovpn_config"`
}

func init() {
	loadEnv()
}

func main() {
	app := &cli.Command{
		Name:     "vynx",
		Version:  version,
		Usage:    "VPN connector with TOTP and config encryption",
		Commands: []*cli.Command{},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			switch cmd.NArg() {
			case 0:
				return connectVPN()
			default:
				if cmd.NArg() != 2 {
					return fmt.Errorf("Usage: vynx <username> <secret>")
				}
				username := cmd.Args().Get(0)
				secret := cmd.Args().Get(1)
				return generateVPNConfig(username, secret)
			}
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}

func loadEnv() {
	// Load the embedded .env file content
	envData, err := envTemplate.ReadFile(".env")
	if err != nil {
		panic("Failed to read embedded .env file")
	}

	// Parse the .env content into a map
	envMap, err := godotenv.Parse(bytes.NewReader(envData))
	if err != nil {
		panic("Failed to parse embedded .env")
	}

	// Set the parsed variables into the environment
	for k, v := range envMap {
		if err := os.Setenv(k, v); err != nil {
			panic("Failed to set env var")
		}
	}
}
