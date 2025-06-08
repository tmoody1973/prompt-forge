# 🚀 GitHub Actions CI/CD

This directory contains GitHub Actions workflows for the PromptForge project.

## 📋 Workflows

### 1. 🧪 Test Workflow (`test.yml`)

**Triggers:**
- Pull requests to `main` and `develop` branches
- Pushes to `main` branch
- Manual dispatch

**Jobs:**
- **Unit Tests**: Runs Go unit tests with coverage reporting
- **Integration Tests**: Builds Docker image and tests health endpoint
- **Security Scan**: Runs Gosec security scanner on PRs

**Features:**
- ✅ Go 1.21 with module caching
- ✅ Race condition detection
- ✅ Coverage reporting with HTML output
- ✅ golangci-lint code quality checks
- ✅ Docker integration testing
- ✅ Security vulnerability scanning
- ✅ Artifact uploads for coverage reports

### 2. 🐳 Docker Workflow (`docker.yml`)

**Triggers:**
- Pushes to `main` branch
- Git tags (releases)
- Pull requests to `main` branch
- Manual dispatch

**Jobs:**
- **Build & Push**: Multi-platform Docker image build and push to GHCR
- **Security Scan**: Trivy vulnerability scanning
- **Notification**: Deployment summary and quick deploy instructions

**Features:**
- ✅ Multi-platform builds (AMD64, ARM64)
- ✅ GitHub Container Registry (GHCR) publishing
- ✅ Smart tagging strategy
- ✅ Docker layer caching
- ✅ Security scanning with Trivy
- ✅ Deployment manifests generation
- ✅ Health check validation

## 🏷️ Image Tagging Strategy

| Event | Tags Generated |
|-------|----------------|
| Push to main | `latest`, `main-<sha>` |
| Git tag `v1.2.3` | `v1.2.3`, `v1.2`, `v1`, `latest` |
| Pull request | `pr-<number>` |
| Other branches | `<branch>-<sha>` |

## 🔐 Required Secrets

The workflows use the following GitHub secrets:

- `GITHUB_TOKEN` (automatically provided) - for GHCR authentication and API access

## 🛠️ Local Development

### Running Tests Locally

```bash
cd api
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Building Docker Image Locally

```bash
docker build -t promptforge:local .
docker run -p 8080:8080 promptforge:local
```

### Code Quality Checks

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
cd api
golangci-lint run
```

## 📦 Docker Image Usage

### Quick Start

```bash
# Pull and run latest image
docker run -p 8080:8080 ghcr.io/your-username/prompt-workbench:latest
```

### Production Deployment

```bash
# Using docker-compose (generated in workflow artifacts)
docker-compose up -d
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DATABASE_PATH` | `/data/promptforge.db` | SQLite database path |

## 🔍 Monitoring & Health Checks

The Docker image includes built-in health checks:

```bash
# Check container health
docker inspect --format='{{.State.Health.Status}}' <container-name>

# Manual health check
curl http://localhost:8080/api/health
```

## 🚀 Deployment Options

### 1. Docker Run

```bash
docker run -d \
  --name promptforge \
  -p 8080:8080 \
  -v promptforge_data:/data \
  --restart unless-stopped \
  ghcr.io/your-username/prompt-workbench:latest
```

### 2. Docker Compose

```yaml
version: '3.8'
services:
  promptforge:
    image: ghcr.io/your-username/prompt-workbench:latest
    ports:
      - "8080:8080"
    volumes:
      - promptforge_data:/data
    restart: unless-stopped
    environment:
      - PORT=8080
      - DATABASE_PATH=/data/promptforge.db

volumes:
  promptforge_data:
```

### 3. Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: promptforge
spec:
  replicas: 1
  selector:
    matchLabels:
      app: promptforge
  template:
    metadata:
      labels:
        app: promptforge
    spec:
      containers:
      - name: promptforge
        image: ghcr.io/your-username/prompt-workbench:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: DATABASE_PATH
          value: "/data/promptforge.db"
        volumeMounts:
        - name: data
          mountPath: /data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: promptforge-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: promptforge-service
spec:
  selector:
    app: promptforge
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

## 🔧 Troubleshooting

### Common Issues

1. **Build Failures**: Check Go version compatibility and dependency issues
2. **Test Failures**: Ensure all environment variables are set correctly
3. **Docker Push Failures**: Verify GITHUB_TOKEN permissions for packages
4. **Health Check Failures**: Check if the application starts properly in container

### Debug Commands

```bash
# View workflow logs
gh run list
gh run view <run-id>

# Test Docker image locally
docker run --rm -it promptforge:local sh

# Check container logs
docker logs <container-name>
```

## 📊 Coverage Reports

Coverage reports are automatically generated and uploaded as artifacts in the test workflow. You can:

1. Download coverage artifacts from the GitHub Actions run
2. View the HTML coverage report in your browser
3. Check coverage percentages in the workflow output

## 🎯 Next Steps

- [ ] Add end-to-end tests
- [ ] Set up staging environment deployment
- [ ] Add performance benchmarks
- [ ] Implement blue-green deployment strategy
- [ ] Add monitoring and alerting 