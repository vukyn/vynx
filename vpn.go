package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/vukyn/vynx/pkg/cryp"
	"github.com/vukyn/vynx/pkg/otp"
)

func generateVPNConfig(username string, secret string) error {
	key := os.Getenv("AES_KEY")
	encSecret, err := cryp.EncryptAES(secret, key)
	if err != nil {
		return err
	}
	cfg := VPNConfig{
		Username:   username,
		Secret:     encSecret,
		OvpnConfig: "myvpn.ovpn",
	}
	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}
	if err := os.WriteFile("vpn_config.json", data, 0600); err != nil {
		return err
	}
	fmt.Println("vpn_config.json generated with encrypted secret.")
	return nil
}

func connectVPN() error {
	// Check if openvpn CLI is installed
	if _, err := exec.LookPath("openvpn"); err != nil {
		fmt.Println("Error: openvpn CLI is not installed or not found in PATH.")
		switch runtime.GOOS {
		case "linux":
			fmt.Println("To install on Linux: sudo apt-get install openvpn")
		case "windows":
			fmt.Println("To install on Windows: choco install openvpn")
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

	// Prepare auth file
	authFile := ".vpn_auth.tmp"
	authContent := fmt.Sprintf("%s\n%s\n", cfg.Username, otp)
	if err := os.WriteFile(authFile, []byte(authContent), 0600); err != nil {
		fmt.Println("Failed to write auth file:", err)
		return err
	}
	defer os.Remove(authFile)

	// Call OpenVPN
	cmd := exec.Command("sudo", "openvpn", "--config", cfg.OvpnConfig, "--auth-user-pass", authFile)

	// Optional: pipe stdout/stderr to your console
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
