<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version"/>
  <img src="https://img.shields.io/badge/License-Apache_2.0-blue?style=for-the-badge" alt="License"/>
  <img src="https://img.shields.io/badge/Kubernetes-Native-326CE5?style=for-the-badge&logo=kubernetes&logoColor=white" alt="K8s Native"/>
</p>

<h1 align="center">🔍 Trace</h1>
<h3 align="center">Supply Chain Security & Vulnerability Analysis for Cloud-Native Infrastructure</h3>

<p align="center">
  <a href="#-quick-start">Quick Start</a> •
  <a href="#%EF%B8%8F-how-it-works">How It Works</a> •
  <a href="#-api-reference">API Reference</a> •
  <a href="#-contributing">Contributing</a>
</p>

---

## 🚨 The Problem

**Modern container deployments face a critical blind spot: software supply chain visibility.**

Every container image you deploy contains hundreds of dependencies—each one a potential attack vector. Without visibility into your software bill of materials (SBOM) and real-time vulnerability intelligence:

- 🔴 **Security teams** can't assess risk before deployment
- 🔴 **DevOps engineers** spend hours manually tracking dependencies
- 🔴 **Compliance audits** become a nightmare of spreadsheets
- 🔴 **Zero-day vulnerabilities** go undetected in production workloads

> *"The average container image contains 180+ packages. Are you tracking vulnerabilities in all of them?"*

---

## ✅ The Solution

**Trace is a Go-based supply chain security engine** that provides deep visibility into containerized deployments through automated SBOM generation and vulnerability correlation.

### Key Capabilities

| Feature | Description |
|---------|-------------|
| 📦 **Deep Container Inspection** | Layer-by-layer analysis of OCI-compliant images |
| 🔒 **SBOM Generation** | CycloneDX and SPDX-compliant manifests |
| ⚡ **Vulnerability Correlation** | Real-time CVE matching via Grype database |
| ☁️ **Cloud-Native Storage** | S3-backed artifact retention with presigned URLs |
| 🔌 **REST API** | High-performance API for CI/CD pipeline integration |

---

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- Docker
- AWS credentials (for S3 storage)

### Installation

```bash
# Clone the repository
git clone https://github.com/dakshhhhh16/trace.git
cd trace

# Build the binary
make build
```

### Running Trace

```bash
# Configure environment
export AWS_REGION="us-east-1"
export S3_BUCKET_NAME="trace-scans"

# Start the server
./bin/trace
```

The Trace engine will start on port `7789`.

### Your First Scan

```bash
# Submit a container image for scanning
curl -X POST http://localhost:7789/api/scan \
  -H "Content-Type: application/json" \
  -d '{
    "image_name": "ubuntu:latest",
    "org_id": "my-org",
    "image_id": "ubuntu-prod-001"
  }'
```

**Response:**
```json
{
  "message": "Manifest and vulnerabilities generated and uploaded successfully",
  "files": {
    "manifest": {
      "s3_key": "trace/my-org/ubuntu-prod-001/manifest.json",
      "download_url": "https://..."
    },
    "vulnerabilities": {
      "s3_key": "trace/my-org/ubuntu-prod-001/vulnerabilities.json",
      "download_url": "https://..."
    }
  }
}
```

---

## ⚙️ How It Works

Trace operates as a modular scanning engine with three core components:

```
┌─────────────────────────────────────────────────────────────┐
│                     Trace Architecture                      │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │
│   │   REST API  │───▶│ Scan Engine │───▶│  S3 Storage │     │
│   │    (Gin)    │    │             │    │             │     │
│   └─────────────┘    └──────┬──────┘    └─────────────┘     │
│                             │                               │
│                    ┌────────┴────────┐                      │
│                    ▼                 ▼                      │
│            ┌─────────────┐   ┌─────────────┐                │
│            │    Syft     │   │    Grype    │                │
│            │ (SBOM Gen)  │   │ (Vuln Scan) │                │
│            └─────────────┘   └─────────────┘                │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Scanning Pipeline

1. **Image Acquisition**: Trace fetches the container image layers from the specified registry
2. **SBOM Generation**: The Syft engine catalogs all packages, libraries, and dependencies
3. **Vulnerability Matching**: Grype correlates packages against the vulnerability database with multi-ecosystem support:
   - Go, Java, JavaScript, Python, Ruby, Rust, .NET
4. **Artifact Storage**: Results are uploaded to S3 with presigned URLs for secure access
5. **API Response**: Clients receive download links for both SBOM and vulnerability reports

### Multi-Matcher Vulnerability Engine

Trace uses Grype's sophisticated matching system to reduce false positives:

```go
// Each ecosystem has specialized matching logic
matcher.Config{
    Golang:     golang.MatcherConfig{UseCPEs: true},
    Java:       java.MatcherConfig{MavenBaseURL: "..."},
    Javascript: javascript.MatcherConfig{UseCPEs: false},
    Python:     python.MatcherConfig{UseCPEs: false},
    // ... additional ecosystems
}
```

---

## 📡 API Reference

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Health check |
| `POST` | `/api/scan` | Submit image for scanning |
| `GET` | `/api/scans/:org_id/:image_id` | List scans for an image |
| `GET` | `/api/scans/:org_id/:image_id/:filename` | Download scan file |
| `DELETE` | `/api/scans/:org_id/:image_id` | Delete scans for an image |

### Scan Request Schema

```json
{
  "image_name": "string",   // Required: Full image reference
  "arch": "string",         // Optional: Architecture (amd64, arm64)
  "registry": "string",     // Optional: Registry URL
  "org_id": "string",       // Optional: Organization identifier
  "image_id": "string"      // Optional: Custom image identifier
}
```

---

### Development Setup

```bash
# Clone your fork
git clone https://github.com/<your-username>/trace.git
cd trace

# Install dependencies
go mod download

# Run linting
make lint

# Build
make build
```

### Contribution Guidelines

1. **Fork & Clone**: Fork the repository and clone locally
2. **Branch**: Create a feature branch (`git checkout -b feature/amazing-feature`)
3. **Code**: Make your changes following Go best practices
4. **Test**: Ensure your code passes linting (`make lint`)
5. **Commit**: Write clear, semantic commit messages
6. **PR**: Open a Pull Request with a detailed description

---

<p align="center">
  <b>Built with ❤️ for the Cloud-Native Community</b>
</p>
