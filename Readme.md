# Trace

**Next-Generation Supply Chain Security**

Trace is an advanced SBOM and vulnerability analysis engine designed for modern containerized infrastructure. It provides deep visibility into software supply chains, generating CycloneDX-compliant manifests and identifying security risks with precision.

## Core Capabilities

- **Deep Container Inspection**: Layer-by-layer analysis of OCI-compliant images.
- **Vulnerability Correlation**: Real-time CVE matching against comprehensive vulnerability databases.
- **SBOM Generation**: Automated production of high-fidelity software bills of materials.
- **Cloud Native Integration**: Seamless operation within AWS and Kubernetes environments.

## Architecture

Trace operates as a modular scanning engine.

- **Vulnerability Engine**: Powered by Grype.
- **SBOM Engine**: Powered by Syft.
- **Storage**: S3-backed artifact retention (SBOMs, Vulnerability Reports).
- **API**: High-performance REST API for integration with CI/CD pipelines.

## Getting Started

### Prerequisites

- Go 1.24+
- Docker
- AWS Credentials (for S3 storage)

### Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/dakshhhhh16/trace.git
cd trace
make build
```

### Running Trace

Ensure your environment is configured:

```bash
export AWS_REGION="us-east-1"
export S3_BUCKET_NAME="trace-scans"
./bin/trace
```

The server will initialize on port `7789`.

## API Documentation

Trace exposes a RESTful API for submitting scans and retrieving reports.

**Health Check**
```http
GET /
```

**Submit Scan**
```http
POST /api/scan
Content-Type: application/json

{
  "image": "ubuntu:latest",
  "org_id": "org-123"
}
```

## Contributing

We welcome contributions to the Trace engine. Please review our contribution guidelines before submitting pull requests.

## License

(c) Daksh Pathak. All rights reserved.
