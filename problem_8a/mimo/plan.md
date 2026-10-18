# Execution Plan: Type-B benchmark-solution package `phys-demo/problem_8a/mimo`

Tool: OpenCode on Mimo-V2.6-Flash Free

## 0. Scope and discipline

- **Writes only under** `phys-demo/problem_8a/mimo/` (plus nothing else — no edits to `PROTOCOL.md`, the ADR, sibling solutions, or any infrastructure repo).
- **Independence:** the derivation is written from `problem.md`'s (A62) only. I will not open `muse/`, `muse+deepseek/`, or `inkling/` during implementation; no values or code from siblings are used as oracles. README discloses: author has prior structure-level audit exposure; derivation reproduced independently; correctness anchored by internal checks (dimensions, thresholds, limits, three-path arithmetic agreement), never by sibling agreement.
- Baseline infra revisions recorded in README (`git rev-parse` for phys-math/phys-artifact).

## 1. Physics route (independent derivation → what gets represented)

**General (relativistic)** from the axion dipole in the dielectric with response (A62), plasmon pole at ω = ω_p, energy conservation ω_γ = E_a − ω_p:

Γ(E_a, k_a) = g²ω_p q_γ² k_a² / (16π E_a ω_γ) · J(κ) · Θ(E_a − ω_p),  q_γ = n ω_γ,  κ = k_a/q_γ,

J(κ) = 2π[ (1+κ²)/(2κ²) − (1−κ²)²/(4κ³) · ln((1+κ)/|1−κ|) ], with hand-checked limits **J(0) = 8π/3** (series: J = 8π/3 − (8π/15)κ² + O(κ⁴), derived in README) and **J(1) = 2π** (limit).

**NR specialization** (E_a = m, J(0) = 8π/3; the π's cancel exactly):

**Γ_NR = g² ω_p n² (m − ω_p) m v_a² / 6** — pure rational arithmetic, no π.

**Results (computed at runtime, not pre-baked):**
- (1) m = 25 eV: all inputs rational (g = 10⁻¹⁹ eV⁻¹ exactly, ω_p = 83/5, n = 3/2, v = 1/1000) → exact rational oracle; expected ~1.31×10⁻⁴¹ eV.
- (2) m = 15 eV: E_a = m(1+v²/2) = 15.0000075 < ω_p → threshold gate → **0**, exact rational comparison.
- **Conventions stated as conventions:** photon measure/dispersion (given-n on-shell; residual one-power-n normalization ambiguity recorded as a `limitation`, not hidden); NR convention E_a = m with a test showing full-E_a insensitivity (relative ~1.5×10⁻⁶ ≪ 3 s.f.).

**Independent checks (all solver-local):** dimensional analysis; threshold behavior; J-series + J(1) special value; full-J(κ) float64 vs exact-rational J(0) value agreeing within 10⁻⁴ relative; E_a-convention insensitivity; unit-conversion exactness.

## 2. Package layout

```text
phys-demo/problem_8a/mimo/
  go.mod          module github.com/PithomLabs/phys-demo/problem_8a/mimo (go 1.25.7)
  go.work         use . + ../../../phys-artifact + ../../../phys-math   (phys-lib NOT wired: no genuine reuse)
  manifest.json   schema phys-demo/benchmark-solution-manifest/1 (PROTOCOL §11.2, all required keys)
  manifest_test.go     strict schema conformance + corruption negatives + computed status
  reuse_test.go        infrastructure_use exercise tests (root package, so the manifest test
                       can source-scan and prove each recorded test name exists)
  axion/          doc.go, axion.go — records, exact params (big.Rat), Γ_NR rational,
                  threshold gate, general float formula, J(κ) closed form + series, unit helpers
  representation/ representation.go — physmath.Store minting (formula/params/assumptions as
                  artifacts), Derivation build, ReplayCheck() error, Sources register
  protocol/       derivation_test.go — entry points, three-path results, thresholds,
                  unit slip, J checks, replay-corruption rejection
  README.md       capability audit, derivation, executed-vs-represented boundary,
                  limitations, reproduction commands
```

## 3. Infrastructure reuse (real, each with an exercising test)

Spike-verified chain (source-read, `rules.go`/`apply.go`): `EXPR-SUBSTITUTE-001` (target expression + sorted declaring-function refs + bound-parameter → rational-ref entries) then `EXPR-SIMPLIFY-001` (`simplifyStep` rule 5 folds an all-rational Mul/Add tree to a single Rational — confirmed in source), optionally a final `RAT-NORMALIZE-001` if it yields a distinct declared output (spike decides; drop if redundant, since Replay requires recomputed ≡ declared bytes).

| `infrastructure_use` concern | library | api_symbol | exercising test |
|---|---|---|---|
| Exact-rational parameters & artifact identity | phys-math / phys-artifact | `physmath.NewRational`, `MemStore.Put`, `ArtifactRef` determinism | `TestInfraRationalEncoding` |
| Replayable substitution→fold derivation (Γ_NR exact value) | phys-math | `physmath.Replay` (steps `EXPR-SUBSTITUTE-001`, `EXPR-SIMPLIFY-001`) | `TestInfraDerivationReplay` (+ corrupted-output rejection) |
| Equation/assumption records as artifacts | phys-math | `Relation`/expression minting through `Store` | covered by encoding test |

**Deviations (all four fields each):** `DEV-FRAMEWORK` (no axion/QED `phys-*`; `gr.NewFramework` hardcodes `phys-gr` → forbidden misuse, PROTOCOL §8.4), `DEV-ANALYTIC-STEPS` (δ-collapse/tensor/phase-space: δ-distributions and consuming rules absent → `derivations[].support: REPRESENTATIONAL ONLY`), `DEV-LIMIT` (no limit/asymptotic rule; J(0), J(1) hand-derived + solver-local checked), `DEV-NUMERICS` (Numerical deferred, Spec §0.4/§23.2 → float64 local), `DEV-UNITS` (no dimensional semantics, Arch §32.2 → metadata + local conversion test).

**Capability audit** states: `SUPPORTED` exact rationals, expression/relation/derivation artifacts, substitute/simplify/replay; `PARTIALLY SUPPORTED` simplify (closed rewrite list); `REPRESENTATIONAL ONLY` `integral/1`, analytic derivation steps; `DEFERRED` numerical execution, units; `UNAVAILABLE` δ-distributions, limits.

## 4. Manifest and conformance logic

- All §11.2 keys: `problem` (arXiv:2610.16733v1, App. A Prob. 8(a) eq. A62, fidelity "problem-definition only; solution not consulted"), `parameters` (unit states: g/n/v explicit or dimensionless per §9.1 table; convention ħ=c=1), `capability_audit`, `derivations[]` (four steps with honest support states), `infrastructure_use`, `deviations`, `solver_local_computations[]` (J(κ) float64, threshold gate, unit conversion, formatter — `executed_by: "solver-local"`), `results[]` (evidence levels from §10 only: (1) `Replayed/verified` for the arithmetic with limitation "premise physics is represented, not protocol-established"; (2) `Executed`, local gate), `limitations[]`, `open_obligations[]`, `declared_gaps[]` — **all advisory** (`GAP-NO-DISTRIBUTION-RULE`, `GAP-NO-INTEGRAL-RULE`, `GAP-NO-CENTRAL-VALIDATOR`), `formalization_refs: []` (zero formalization claims → firewall inert), `entry_points` = `{primary: axion.OrderedResults, secondary: representation.ReplayCheck}`.
- **Severity judgment stated in README:** *mandatory* = a §11.3 conformance element unfulfillable in this package; capability absences correctly recorded through required fields are *advisory* — otherwise every honest benchmark package would be permanently nonconformant. Auditability comes from `manifest_test.go` computing status, not the manifest declaring it.

## 5. Tests (PROTOCOL §12 ↔ prompt2 item 7)

1. `manifest_test.go`: exact allowed/required keys (unknown rejected), enums vs §5/§10 vocabularies, unit-state field rule, `infrastructure_use.test` names exist via source scan **and** live in the executed test set, deviation completeness, **corruption negatives** (wrong `schema_id`, missing `unit`, dangling test ref, injected formalization claim without `formalization_refs`, injected mandatory gap → status flips to NONCONFORMANT), computed `CONFORMANT` assertion.
2. `reuse_test.go`: the two infra tests above — real API calls, byte-determinism, replay success, corrupted-output failure.
3. `protocol/derivation_test.go`: entry points called and observed (order `[(1),(2)]`, (2) = 0); three-path agreement (replayed rational ≡ local `big.Rat` ≡ float64 within tol); `TestUnitSlip` (g conversion 10⁻¹⁰ GeV⁻¹ = 10⁻¹⁹ eV⁻¹ exact; deliberate slip differs by 10¹⁸ and is detected); convention-insensitivity; J-series/J(1) checks; `go vet` + `gofmt` clean.

## 6. Execution order

1. Module skeleton (`go.mod`, `go.work`).
2. **Spike:** substitute→simplify chain in a scratch test (fold to single rational? declaration/sort constraints?) → fixes `derivations[]` chain shape; delete scratch.
3. `axion/` records + exact params + formulas + J helpers.
4. `representation/` minting + derivation build + `ReplayCheck`.
5. `manifest.json` + `manifest_test.go` (+ `reuse_test.go`).
6. `protocol/derivation_test.go`.
7. `README.md` (independent derivation from A62, capability audit, boundary, limitations, conventions).
8. Validate: `gofmt -l`, `go vet ./...`, `go test -count=1 ./...` from `mimo/`; then `phys-demo/scripts/test-all.sh` (all 4 projects must pass). Report per PROTOCOL §13.2: commands, revisions, represented-vs-local split.

**Known risks:** (a) substitute-chain edge cases (sorted decls, whether a redundant-output step must be dropped) — resolved by the spike before manifest fields are frozen; (b) `ValidateFunctionBody` constraints on the declaring function — spike covers; (c) float J(κ) near κ=1 — irrelevant here (κ ≈ 2×10⁻³), guarded anyway.
