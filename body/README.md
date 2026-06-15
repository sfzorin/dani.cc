# body — exercise figure generator

Standalone workspace for the morning app's exercise illustrations.
Figures are **generated** from a logical 3D skeleton, not hand-drawn, and exported as
flat schematic **SVG** (light fill + thin dark contour, themeable, crisp at any size).

## Single source of truth: `core.mjs`
One dependency-free ES module holds ALL the figure math, used by **both** the browser
editor and the Node generator (no duplicated logic):
- `joints(pose, sex)` — builds a 3D skeleton: pelvis root, pelvis bar, spine, shoulder
  girdle (chest + shoulder cross-bars), neck, head, elbows/wrists/hands, knees/ankles/toes.
  A pose = joint **direction** vectors (bone lengths fixed → consistent proportions).
  Joint limits: spine cone 80°, shoulder 160°, hip 120°, elbow/knee hinge 3–145°,
  neck 45°, wrist 55°, foot 75°.
- `cameraBasis(yaw, pitch)` / `VIEWS.{side,front,top}` + `project(p, cam)` — orthographic
  projection from an **arbitrary orbit camera** or the 3 fixed views. `[screenX, screenY, depth]`.
- `figureSVG(J, cam, {sex, level, ghost, xf})` — the meat. **xf={sc,tx,ty}** bakes the
  placement transform straight into coordinates (no `<g transform>` wrapper). The flesh is
  built so the silhouette does NOT just snap to the skeleton joints (joints drive the body,
  they aren't its corners):
    - **shoulders** — a sloping trapezius *yoke* from the neck base down to a rounded
      **deltoid cap** that sits below+outside the joint (FLESH `shDrop`/`shOut`); the torso
      reaches the ribcage, the deltoid caps the arm root. No flat horizontal shoulder bar.
    - **torso** — lofted from three soft volumes (pelvis → waist → ribcage) with extra
      profile samples and **curved sides** (quadratic smoothing) → waist pinch, rounded
      corners, no trapezoid facets. Cross-section ellipse projected per-camera.
    - **hands** — rounded **mitten** (narrow wrist → wider palm → softly rounded tip,
      length > width so no "flipper") + a small soft **thumb bump** on the forward-facing
      side (auto-hidden when the hand is edge-on).
    - **feet** — rounded **wedge** with a flat sole: ankle on top/back, heel behind, toe
      ahead, instep curving up (built in screen space along the toe dir, up-normal flipped
      toward the leg). End-on (toe-to-camera) → small oval. (`#FT` path is now unused.)
    - limbs = tapered capsules; head = ellipse; knee/elbow knobs (×1.18) read as joints.
- `skeletonSVG`, `figureExtent` (auto-fit bbox), `neutralPose`.

**Three render levels** (`level`): L0 skeleton, L1 capsules (depth-sorted, visible joints),
**L2 = single smoothed silhouette** (two-pass: every part DARK with a fat dark stroke =
expanded outline, then every part GRAY on top → uniform contour, no internal seams).
**Production = L2.**

> **Render gotcha:** L2 stacks many same-colour overlapping fills. Headless Chrome's GPU
> rasterizer **conflates** these past a certain overdraw count and punches white holes
> (esp. small contained shapes like joint knobs). Mitigations baked in: single calf cap
> (fewer overlaps), baked coords (no scaled group), and dedup of identical shapes. When
> screenshot-verifying, **embed the SVG in an HTML page** — never open a standalone
> `file://x.svg` (that path is the most flake-prone).

## Editor — `web/body.html` (served at `/body`)
Imports `/body/core.mjs`. Big **3D orbit window** (drag to rotate, L0/L1/L2 + camera
presets) + three small **editable** projections (side/front/top) where you drag joint
handles. Served by the Go app: `//go:embed body/core.mjs` + `//go:embed body/web/body.html`
in `main.go`, routes `/body` and `/body/core.mjs` on the http mux.

## Generator — `gen.mjs`
`node body/gen.mjs` → writes SVGs to `svg3d/`. Imports `core.mjs`. Holds the `EX` exercise
poses (one 3D pose each + optional arrows / safety lines / ghost reference pose), auto-fits
(centre + scale + floor line), emits to `svg3d/{id}.svg`. Production level = L2.

## Dirs / files
- `core.mjs` — shared figure math (single source of truth).
- `web/body.html` — interactive editor.
- `gen.mjs` — SVG generator (run with node).
- `svg3d/` — generated figures.
- `svg/`, `gen.py` — deprecated hand-drawn 2D attempts / old 2D generator (kept for reference).
