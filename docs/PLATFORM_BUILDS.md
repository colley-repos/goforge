# Platform Build Guide

## Supported Targets

| Platform | Renderer    | Build Method                        | Status  |
|----------|-------------|-------------------------------------|---------|
| Windows  | Ebitengine  | `go build`                          | Planned |
| Linux    | Ebitengine  | `go build` (native or container)    | Planned |
| macOS    | Ebitengine  | `go build` (native, requires Cgo)   | Planned |
| Android  | Ebitengine  | `ebitenmobile bind -target android`  | Planned |
| iOS      | Ebitengine  | `ebitenmobile bind -target ios`      | Planned |
| Web      | Ebitengine  | `GOOS=js GOARCH=wasm go build`      | Planned |
| Terminal | Terminal    | `go build` (any platform)           | Planned |

## Build Notes

### Desktop (Windows/Linux/macOS)
- Windows: no Cgo required
- Linux/macOS: Cgo required (link against OpenGL/Metal)
- Cross-compilation to Windows/WASM works; other targets need native builds

### Mobile
- Requires `ebitenmobile` tool: `go install github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile@latest`
- Android: outputs `.aar` file, needs Android Studio wrapper project
- iOS: outputs `.xcframework`, needs Xcode wrapper project
- Both need platform-specific project shells (provided in scaffold templates)

### Web (WASM)
- `GOOS=js GOARCH=wasm go build -o game.wasm`
- Serve with provided `index.html` + `wasm_exec.js`
- Note: no filesystem access, assets must be embedded or fetched via HTTP

## CI Strategy (Deferred)

GitHub Actions matrix build per platform. Each platform gets its own job
with the native toolchain. `GOWORK=off` in all CI jobs.

The CI setup will be implemented in a future phase. This doc captures
the intended configuration for when that work begins.
