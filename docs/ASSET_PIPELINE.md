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

## Asset Assessment & Cataloguing

### Never Trust Filenames

Asset filenames (from third-party packs, AI generators, or even internal teams)
do NOT reliably describe what the mesh looks like or how it should be used.
A file named `SM_Bld_Cover_01` might be a shop awning, not a tactical cover wall.

**All assets MUST be visually assessed before use.** The assessment pipeline:

1. **Inspect geometry**: read collision shapes, bounding boxes, vertex bounds to
   determine actual dimensions (height, width, depth)
2. **Classify by role**: based on actual shape, assign roles:
   - `half_cover` (waist-height, ~0.8-1.5m): dumpsters, traffic barriers, crates
   - `full_cover` (tall, ~2.0m+): walls, shipping containers, vehicles
   - `obstacle` (impassable): rubble, locked doors, pillars
   - `floor_tile`: road, sidewalk, ground surfaces
   - `prop` (decorative): street lights, signs, cones
   - `character`: humanoid meshes with rigs
3. **Record the assessment**: store classification, actual dimensions, recommended
   scale, rotation, and usage notes in the asset manifest

### Assessment Cache

```json
{
  "SM_Prop_Skip_02": {
    "assessed": true,
    "actual_description": "Industrial dumpster, rectangular",
    "dimensions": { "width": 2.3, "height": 1.87, "depth": 1.2 },
    "role": "half_cover",
    "recommended_scale": 0.85,
    "rotation_notes": "Long side should face attacker for max coverage",
    "source_pack": "polygon-city-01"
  }
}
```

Once assessed, an asset never needs re-evaluation unless the source file changes.
The manifest serves as the authoritative record of what each asset actually IS,
not what its filename suggests.

### Skeletal vs Static Mesh Detection

When loading 3D assets programmatically, distinguish between:
- **Static meshes** (`.res`, `.glb` props/environment): safe to load as plain
  MeshInstance3D. Just set `mesh_instance.mesh = loaded_mesh`.
- **Skeletal meshes** (character `.mesh`, `.glb` with armature): MUST be loaded
  with a Skeleton3D parent node, otherwise they render in bind pose (T-pose).
  Either load the full scene/prefab (which includes the skeleton) or construct
  a Skeleton3D programmatically.

The assessment cache should record `"is_skeletal": true/false` for each asset
so the loader knows which path to take.

### Animation Pipeline — Lessons Learned

**Do NOT:**
- Load raw `.fbx` animation files at runtime — Godot/engine importers only process
  these at import time, not runtime. They must be pre-compiled to native formats.
- Use programmatic `set_bone_pose_rotation()` to pose characters — this fights
  the skeleton's rest pose processing and produces broken/subtle results.
  Always use proper animation playback through AnimationPlayer.
- Assume bone names match between animation sources and target skeletons.
  Always verify track paths match skeleton bone names before playback.

**Do:**
- Use pre-compiled AnimationLibrary resources (`.res`) that can be loaded at runtime
  via `anim_player.add_animation_library("name", loaded_lib)`.
- Load full character prefab scenes (with Skeleton3D + skin bindings) rather than
  raw mesh files. Skeletal meshes without a skeleton render in T-pose.
- Inspect animation track paths before use: track paths contain
  `SkeletonNodeName:BoneName` and must match the actual scene tree.
- Build a runtime retargeting function that remaps bone names in track paths
  if the animation source uses different naming (e.g. Mixamo `mixamorig_Hips`
  vs Godot humanoid `Hips`). This is string replacement on `NodePath`, not
  complex math.
- Use `%NodeName` (unique name) syntax in track paths for robust resolution
  regardless of scene tree depth.

**Recommended sources for free animation libraries:**
- `catprisbrey/Godot4-OpenAnimationLibraries` — 150+ combat/locomotion animations
  as native Godot AnimationLibrary `.res` files, already using Godot humanoid
  bone names. Includes ShooterLib (rifle, pistol, melee, movement) and MeleeLib.
- `ZenXChaos/ThirdPersonShooter-AnimationSets` — 40 FBX files (Unlicense),
  need editor import + retargeting but cover idle/walk/shoot/death/crouch.
- Mixamo.com — 2300+ animations, requires Adobe account, download as FBX,
  need editor import with BoneMap retargeting.

### Assessment for AI-Generated Assets (Meshy)

Meshy outputs are even less predictable than named pack assets. After generation:
1. Load the GLB and compute bounding box automatically
2. Run a geometry classifier (convex hull analysis) to detect role
3. Flag for manual review if classification confidence is low
4. Never deploy an unassessed Meshy asset into gameplay
