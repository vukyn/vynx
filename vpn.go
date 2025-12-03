package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/atotto/clipboard"
	"github.com/vukyn/vynx/pkg/cryp"
	"github.com/vukyn/vynx/pkg/otp"
)

func generateVPNConfig(username string, password string, secret string) error {
	cfg := VPNConfig{
		Username:   username,
		OvpnConfig: "myvpn.ovpn",
	}
	key := os.Getenv("AES_KEY")
	// Encrypt password if provided
	if password != "" {
		encPassword, err := cryp.EncryptAES(password, key)
		if err != nil {
			return err
		}
		cfg.Password = encPassword
	}
	// Encrypt secret if provided
	if secret != "" {
		encSecret, err := cryp.EncryptAES(secret, key)
		if err != nil {
			return err
		}
		cfg.Secret = encSecret
	}
	if password == "" && secret == "" {
		return errors.New("either password or secret must be provided")
	}
	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}
	if err := os.WriteFile("vpn_config.json", data, 0600); err != nil {
		return err
	}

	fmt.Println("vpn_config.json is generated.")
	return nil
}

func connectVPN() error {
	// Check if openvpn CLI is installed
	if _, err := exec.LookPath("openvpn"); err != nil {
		fmt.Println("Error: openvpn CLI is not installed or not found in PATH.")
		switch runtime.GOOS {
		case "linux":
			fmt.Println("To install on Linux: sudo apt-get install openvpn")
		case "darwin":
			fmt.Println("To install on macOS: brew install openvpn")
		default:
			fmt.Println("Please install openvpn CLI: https://openvpn.net/connect-docs/command-line.html.")
		}
		return nil
	}

	cfg, err := loadConfig("vpn_config.json")
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return err
	}

	if cfg.Password == "" && cfg.Secret == "" {
		return errors.New("no password or secret found in config")
	}

	var password string

	// Generate OTP from secret
	if cfg.Secret != "" {
		secret, err := cryp.DecryptAES(cfg.Secret, os.Getenv("AES_KEY"))
		if err != nil {
			fmt.Println("Failed to decrypt secret:", err)
			return err
		}
		otp, err := otp.GenerateTOTP(secret)
		if err != nil {
			fmt.Println("Failed to generate OTP:", err)
			return err
		}
		password = otp
		// Write text to clipboard
		if err := clipboard.WriteAll(otp); err != nil {
			fmt.Println("Failed to copy to clipboard:", err)
			return err
		}
		fmt.Println("OTP is copied to clipboard.")
	}

	// Decrypt password (overwrites OTP if password is also present)
	if cfg.Password != "" {
		decPassword, err := cryp.DecryptAES(cfg.Password, os.Getenv("AES_KEY"))
		if err != nil {
			fmt.Println("Failed to decrypt password:", err)
			return err
		}
		password = decPassword
	}

	// Prepare auth file with absolute path
	// Use /tmp to ensure it's accessible when running with sudo
	authFile := "/tmp/.vynx_auth.tmp"
	authContent := fmt.Sprintf("%s\n%s\n", cfg.Username, password)
	if err := os.WriteFile(authFile, []byte(authContent), 0600); err != nil {
		fmt.Println("Failed to write auth file:", err)
		return err
	}
	defer os.Remove(authFile)

	// Get absolute path for ovpn config
	ovpnConfigPath := cfg.OvpnConfig
	if !filepath.IsAbs(ovpnConfigPath) {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		ovpnConfigPath = filepath.Join(cwd, ovpnConfigPath)
	}

	// Call OpenVPN
	cmd := exec.Command("sudo", "openvpn", "--config", ovpnConfigPath, "--auth-user-pass", authFile)

	// Pipe stdin/stdout/stderr to allow interactive input (for sudo password prompt)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Connecting to VPN...")
	if err := cmd.Run(); err != nil {
		fmt.Println("VPN connection failed:", err)
	}
	return nil
}

func loadConfig(filename string) (*VPNConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg VPNConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
