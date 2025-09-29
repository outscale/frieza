# Frieza

[![Project Graduated](https://docs.outscale.com/fr/userguide/_images/Project-Graduated-green.svg)](https://docs.outscale.com/en/userguide/Open-Source-Projects.html) [![](https://dcbadge.limes.pink/api/server/HUVtY5gT6s?style=flat&theme=default-inverted)](https://discord.gg/HUVtY5gT6s)

<p align="center">
  <img alt="Terminal Icon" src="https://img.icons8.com/ios-filled/100/console.png" width="100px">
  <br/>
  <strong>Cleanup your cloud resources!</strong>
</p>

---

## 🌐 Links

* 📘 Documentation: [Supported Providers](./providers.md)
* 📦 Releases: [GitHub Releases](https://github.com/outscale-dev/frieza/releases)
* 🤝 Contribution Guide: [CONTRIBUTING.md](./CONTRIBUTING.md)
* 💬 Join us on [Discord](https://discord.gg/HUVtY5gT6s)

---

## 📄 Table of Contents

* [Overview](#-overview)
* [Use Cases](#-use-cases)
* [Features](#-features)
* [Requirements](#-requirements)
* [Installation](#-installation)
* [Usage](#-usage)

  * [Manage Profiles](#manage-profiles)
  * [Manage Snapshots](#manage-snapshots)
  * [Cleanup Resources](#cleanup-resources)
  * [Configuration](#configuration)
* [Building from Source](#-building-from-source)
* [License](#-license)

---

## 🧭 Overview

**Frieza** is a CLI tool to clean up cloud resources across multiple providers.

It can either wipe all resources in an account or compare a known “clean state” snapshot and delete only what was created since.

---

## 🚀 Use Cases

* **Full Cleanup:**
  Delete all cloud resources from a specific account or region.

  ```bash
  frieza nuke regionEu2
  ```

* **Snapshot-based Cleanup:**

  1. Keep your essential resources
  2. Create a snapshot

     ```bash
     frieza snap new cleanAccountState regionEu2
     ```
  3. Run temporary experiments that create resources
  4. Clean everything not present in the original snapshot

     ```bash
     frieza clean cleanAccountState
     ```

---

## ✨ Features

* Multi-provider support ([see list](./providers.md))
* Clean resources based on current state or snapshot delta
* Store multiple profiles and configurations
* Track and review deleted resources before execution
* Optional `--auto-approve` for automated cleanup

---

## ✅ Requirements

* Golang (to build from source)
* Make
* A valid cloud account (e.g., Outscale API credentials)

---

## ⚙ Installation

### macOS (via Homebrew)

```bash
brew tap outscale/tap
brew install frieza
```

### Other OS (manual)

Download the latest release from the [Releases page](https://github.com/outscale-dev/frieza/releases).

Or build from source:

```bash
git clone https://github.com/outscale/frieza.git
cd frieza
make install
```

---

## 🧪 Usage

Run the CLI to discover subcommands:

```bash
frieza --help
```

Use `--help` with subcommands for detailed usage.

---

### 🔐 Manage Profiles

Configure access to your cloud providers:

```bash
frieza profile new outscale_oapi --help
frieza profile new outscale_oapi myDevAccount --region=eu-west-2 --ak=XXX --sk=YYY
frieza profile test myDevAccount
frieza profile list
frieza profile describe myDevAccount
frieza profile rm myDevAccount
```

Profiles are stored in: `~/.frieza/config.json`

---

### 📸 Manage Snapshots

Create and manage snapshots of cloud state:

```bash
frieza snapshot new myFirstSnap myDevAccount myOtherAccount
frieza snapshot list
frieza snapshot describe myFirstSnap
frieza snapshot update myFirstSnap
frieza snapshot rm myFirstSnap
```

Snapshots are stored in: `~/.frieza/snapshots/`

---

### 💥 Cleanup Resources

Delete resources created **since a snapshot**:

```bash
frieza clean myFirstSnap
```

Delete **all resources** for a profile:

```bash
frieza nuke myDevAccount
```

You will see a preview of the deletions before execution.
Use `--auto-approve` to skip confirmation prompts.

---

### ⚙ Configuration

Use the `frieza config` subcommands to view and modify CLI options.

---

## 🏗 Building from Source

To build Frieza from source:

```bash
make build
# Output binary is located at: cmd/frieza/frieza
```

---

## 📜 License

**Frieza** is licensed under the BSD 3-Clause License.
© Outscale SAS

License files are available in the [`LICENSES`](./LICENSES) directory.
This project follows the [REUSE Specification](https://reuse.software/).
