package main

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"os"

	"github.com/joho/godotenv"

	"github.com/urfave/cli/v3"
)

//go:embed .env
var envTemplate embed.FS

const version = "1.2.0"
const AES_KEY_1_2_0 = "WBAUGSQ6FEAZMQL7MCMZOR53IX3MOQGC"

type VPNConfig struct {
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"` // Encrypted, base64 (optional)
	Secret     string `json:"secret,omitempty"`   // Encrypted, base64 (optional)
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
			case 2:
				// vynx <username> <password>
				username := cmd.Args().Get(0)
				password := cmd.Args().Get(1)
				return generateVPNConfig(username, password, "")
			case 3:
				// vynx <username> <password> <secret>
				username := cmd.Args().Get(0)
				password := cmd.Args().Get(1)
				secret := cmd.Args().Get(2)
				return generateVPNConfig(username, password, secret)
			default:
				return errors.New("usage: vynx [<username> <password>] or vynx [<username> <password> <secret>]")
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
		// load default aes key
		os.Setenv("AES_KEY", AES_KEY_1_2_0)
		return
		// panic("Failed to read embedded .env file")
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
