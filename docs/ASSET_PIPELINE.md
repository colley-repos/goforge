# Asset Pipeline

## Overview

GoForge's asset system uses the `AssetProvider` interface to abstract asset sources.
The engine never knows where assets come from — it requests them by ID and format,
and the provider delivers bytes.

## Providers

### Meshy (AI-generated 3D)
- REST API at `https://api.meshy.ai`
- Text-to-3D: prompt → preview mesh → refine with textures
- Image-to-3D: reference image → textured model
- Output: GLB (binary glTF) format
- Async workflow: submit task → poll status → download result
- Test mode API key available for integration tests (zero cost)

### Local Filesystem
- Reads from a configured directory
- Supports PNG (2D sprites) and GLB (3D models)
- File watching for hot-reload during development

### Self-Hosted (Future)
- Same REST interface as Meshy but pointed at a local server
- Enables running open-source 3D generation models
- Interface swap: change provider config, zero code changes

## Cache Layer

Content-addressed storage wraps any provider:
- SHA256 hash of asset ID → filename
- Download once, serve from disk thereafter
- Configurable cache directory and max size

## Swapping Providers

```go
// Meshy in development
assetMgr := assets.NewManager(meshy.NewProvider(apiKey))

// Local files in production
assetMgr := assets.NewManager(local.NewProvider("./assets"))

// Self-hosted in CI
assetMgr := assets.NewManager(selfhosted.NewProvider("http://localhost:8080"))
```

The game code never changes — only the provider configuration.
