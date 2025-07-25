# 🔐 Vynx – TOTP-powered VPN Connector

**Vynx** is a lightweight CLI tool that helps you connect to OpenVPN servers using TOTP-based dynamic passwords (2FA). It eliminates the hassle of manually opening your authenticator app and entering the OTP every time you connect.

---

## 🚀 Features

-   🔒 Secure TOTP (RFC 6238) password generation
-   ⚡ One-command VPN connection
-   📁 Encrypted configuration stored locally
-   🧰 Easily scriptable for DevOps or remote access workflows

---

## 🛠️ Requirements

-   A valid `.ovpn` configuration file
-   Golang (only required to build from source)
-   [OpenVPN](https://openvpn.net/) installed and accessible from the command line

For macOS:

```bash
brew install openvpn
```

For Linux:

```bash
sudo apt-get install openvpn
```

---

## 📦 Installation

### Build from source

Clone and build the tool:

```bash
git clone https://github.com/vukyn/vynx.git
cd vynx
touch .env && echo "AES_KEY=YOUR_BASE_32_SECRET_KEY" > .env
make build
```

Make it available system-wide:

```bash
make install
```

### Install from binary

<!-- Download the binary from the [releases](https://github.com/vukyn/vynx/releases) page. -->

Please build from source.

## ⚙️ Usage

### 🔁 1. Connect to VPN

Run the following command to automatically generate the OTP and connect to the VPN:

```bash
$ vynx
```

> **NOTE:** Make sure you have the openvpn profile in the same directory as the .ovpn file

### 🔁 2. Generate a new VPN config

To initialize or update your configuration:

```bash
$ vynx myusername mytotpsecret
```

-   myusername: Your VPN login username
-   mytotpsecret: The Base32-encoded TOTP secret (e.g., from your authenticator app)

> This will create or update vpn_config.json with your credentials.

## 🔐 How It Works

Vynx uses the following workflow:

1. Your secret and username are stored (with encryption).
2. When invoked, Vynx:

    - Reads the local vpn_config.json
    - Generates a current TOTP value (30-second window)
    - Creates a temporary file containing username\npassword
    - Launches OpenVPN using your .ovpn file and the temp credentials

3. The credentials file is deleted automatically after use.

## 🔐 Security Notes

-   vpn_config.json is encrypted and should be protected as sensitive information.
-   No TOTP secret or password is ever printed or stored in plain text.
-   The tool deletes temporary authentication files after execution.

## 📄 License

This project is licensed under the MIT License.

## 📝 Author

-   [@vukyn](https://github.com/vukyn)
