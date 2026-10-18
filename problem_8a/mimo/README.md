# Problem 8(a) — Axion-to-photon conversion in a dielectric medium (mimo)

Type-B benchmark-solution package under `phys-demo`, built to
`phys-demo/PROTOCOL.md` (schema `phys-demo/benchmark-solution-manifest/1`).
Execution plan: [`plan.md`](plan.md).

**Question** (TheoPhysEval, arXiv:2610.16733v1, App. A Prob. 8(a); response
eq. (A62)): compute the rate for axion conversion to photons by scattering
off the electrons of a conventional (non-magnetic, zero-temperature)
dielectric with no external electric field, the axion coupling only via
`a F F̃`, in terms of the dielectric function ε; then, for the plasmon-pole
response

    Im[-1/ε_L(ω,q)] = (π/2) ω_p δ(ω - ω_p),      (A62)

evaluate the non-relativistic rate at g = 10⁻¹⁰ GeV⁻¹, ω_p = 16.6 eV,
n = 1.5, v_a = 10⁻³ for (1) m_a = 25 eV and (2) m_a = 15 eV.

**Result** (exact rational arithmetic; 3-significant-figure display):

| case | m_a | Γ_NR | basis |
|---|---|---|---|
| (1) | 25 eV | **1.31 × 10⁻⁴¹ eV** (exact 5229/(4·10⁴⁴) eV) | above threshold, κ ≈ 1.98×10⁻³, J → 8π/3 |
| (2) | 15 eV | **0 eV** | E_a = m(1+v_a²/2) = 15.0000075 eV < ω_p — kinematically closed |

Reproduce: `cd problem_8a/mimo && go test -count=1 ./...` (project-local
workspace; no env overrides). All 29 tests pass; `manifest_test.go` computes
the package status as **CONFORMANT** (schema valid, both entry points
exercised, 5 `infrastructure_use` entries adopted with exercised tests, 5
explicit deviations, 4 solver-local computations, 4 advisory gaps, 0
mandatory gaps).

---

## 1. Source and fidelity

- Source: `problem.md` (verbatim problem statement), SHA-256
  `7e79f2fc516cc1a5eca23efd8e1af9bb649c0b0d6805792f4dfa7b8f880a0263`.
- Fidelity class: problem-definition only. The solutions appendix and **all
  sibling solution folders were not consulted** — no value, code, or answer
  was taken from another agent or folder, per `problem.md`. The derivation
  below was reproduced independently from (A62) and anchored by internal
  consistency checks (dimensions, thresholds, limits, three-path arithmetic
  agreement), never by agreement with another solution.
- Disclosed exposure (honesty, not an oracle): the package author knew
  protocol-level structure of sibling solutions from the earlier workflow
  audit. No numbers or code were carried over; the manifest records this as
  a limitation.

## 2. Independent derivation (from eq. A62)

Pen-and-paper route; steps marked **[hand]** have no executable rule in the
closed Phys registry and are recorded `REPRESENTATIONAL ONLY`
(`manifest.derivations`), never claimed executable.

1. **[hand]** The `a F F̃` coupling gives the axion a dipole coupling to the
   photon field; in-medium propagation dresses the photon with the
   dielectric response, so the conversion rate follows from the dissipative
   part of the longitudinal response, Im[−1/ε_L], together with the on-shell
   photon kinematics. With the plasmon pole (A62) the medium absorbs at
   ω = ω_p, so energy conservation fixes the outgoing photon energy
   ω_γ = E_a − ω_p and the rate carries a threshold Θ(E_a − ω_p).
2. **[hand]** With the dispersion q_γ = n ω_γ and k_a the axion momentum, the
   two-body phase space reduces to one dimensionless angular integral,
   κ = k_a/q_γ:

   Γ(E_a, k_a) = g² ω_p q_γ² k_a² / (16π E_a ω_γ) · J(κ) · Θ(E_a − ω_p),

   J(κ) = 2π[ (1+κ²)/(2κ²) − (1−κ²)²/(4κ³) · ln((1+κ)/|1−κ|) ].

   The overall normalization follows from the photon final-state measure
   convention (free-photon cut 2πδ(q_γ²)); a residual one-power-of-n
   ambiguity relative to medium-modified phase-space conventions is recorded
   as a limitation (decision G-11; PROTOCOL §14.4), not resolved.
3. **[hand]** κ → 0 limit of J (checked numerically in
   `protocol.TestJSeriesAndLimits`): using
   ln((1+κ)/(1−κ)) = 2(κ + κ³/3 + κ⁵/5 + …),

       (1−κ²)²/(4κ³) · ln((1+κ)/(1−κ)) = 1/(2κ²) − 5/6 + (4/15)κ² + O(κ⁴),

   so J/(2π) = [1/(2κ²) + 1/2] − [1/(2κ²) − 5/6 + (4/15)κ² + …]
   = 4/3 − (4/15)κ² + O(κ⁴), i.e.

       J(0) = 8π/3,   J(κ) = 8π/3 − (8π/15)κ² + O(κ⁴),   J(1) = 2π.
4. **[hand]** Non-relativistic limit: E_a = m(1+v_a²/2) ≈ m, k_a = m v_a,
   ω_γ ≈ m − ω_p, q_γ = n(m − ω_p), and J(0) = 8π/3; the π's cancel exactly
   against the 16π prefactor, leaving the **pure rational**

       Γ_NR = g² ω_p n² (m − ω_p) m v_a² / 6.

   For (1): g = 10⁻¹⁹ eV⁻¹ exactly, ω_p = 83/5, n = 3/2, m = 25,
   v_a = 10⁻³ → Γ = 5229/(4·10⁴⁴) eV = 1.30725×10⁻⁴¹ eV → **1.31×10⁻⁴¹ eV**
   (κ ≈ 2×10⁻³, so the O(κ²) correction to J, ~4×10⁻⁶ relative, is below the
   displayed precision). For (2): E_a = 15.0000075 eV < ω_p = 16.6 eV → the
   threshold closes and Γ = 0 exactly.

## 3. Capability audit (PROTOCOL §5; recorded before coding, manifest mirror)

| Needed capability | State | Disposition |
|---|---|---|
| Exact rational construction | `SUPPORTED` | reused — `physmath.NewRational` (`TestInfraRationalEncoding`) |
| Content-addressed artifact identity | `SUPPORTED` | reused — `artifact.ComputeRef` (`TestInfraArtifactIdentity`) |
| Equation representation (expression/function/relation/derivation families) | `SUPPORTED` | reused — `EncodeRelationBody` + Store minting (`TestInfraRelationMinting`) |
| Parameter substitution + exact rational folding | `SUPPORTED` | reused — `physmath.Apply` with `EXPR-SUBSTITUTE-001`, `EXPR-SIMPLIFY-001` (`TestInfraDerivationReplay`) |
| Byte-exact derivation replay | `SUPPORTED` | reused — `physmath.Replay` (`TestInfraDerivationReplay`) |
| `EXPR-SIMPLIFY-001` generality | `PARTIALLY SUPPORTED` | closed rewrite list only; sufficient (substituted tree is all-rational) |
| Schematic general rate relation | `REPRESENTATIONAL ONLY` | minted as relation/1; no rule evaluates J(κ) or Θ |
| `integral/1` | `REPRESENTATIONAL ONLY` | no consuming rule in the registry |
| Numerical execution | `DEFERRED` | Spec §0.4/§23.2, Arch §25 — all float64 work solver-local |
| Dimensional semantics | `DEFERRED` | Arch §32.2 open — units are metadata (§9.1) |
| δ-distributions | `UNAVAILABLE` | no family, no rule — the (A62) delta-collapse is hand-derived |
| Limits / asymptotics | `UNAVAILABLE` | J(0), J(1) hand-derived, checked solver-locally |
| Heaviside/step rules | `UNAVAILABLE` | exact solver-local rational gate |
| axion/QED `phys-*` framework | `UNAVAILABLE` | none exists; admission is a human decision (Arch §27.5) — deviation `DEV-FRAMEWORK` |

## 4. Executed vs represented vs local (the boundary)

| Content | Where it lives | Status |
|---|---|---|
| NR formula as bound-symbol expression/1 + declaring function/1 | `representation` | **represented** (premise; physics not protocol-established) |
| Substituting the five exact rationals, folding to a single reduced Rational | `physmath.Apply` (`EXPR-SUBSTITUTE-001` → `EXPR-SIMPLIFY-001`) | **executed by Phys**, replayed byte-exactly (`physmath.Replay`) |
| General relativistic rate with J(κ), Θ | `representation.MintGeneralRateRelation` (relation/1, schematic symbols + `ConstPi`) | **represented only** — no rule evaluates it |
| δ-collapse, tensor contraction, phase-space integral, κ→0 limit | hand derivation (§2) | **represented only**; gaps `GAP-NO-DISTRIBUTION-RULE`, `GAP-NO-INTEGRAL-RULE`, `GAP-NO-LIMIT-RULE` |
| Exact Γ_NR arithmetic (second path), threshold gate, GeV→eV conversion, float64 J(κ)/general formula, formatting | `axion` (stdlib only) | **solver-local** (§9.2), never attributed to Phys |

**Three-path agreement** (`protocol.TestThreePathAgreement`): the replayed
rational (Phys bytes) ≡ the solver-local `big.Rat` closed form ≡ a float64
evaluation of the general formula with the full J(κ) and full kinematics
(relative < 10⁻⁴). Replay establishes reproducibility of the recorded
arithmetic under the closed rule registry — **not** the physics of the
premises (PROTOCOL §10).

## 5. Reuse and deviations (manifest `infrastructure_use[]`, `deviations[]`)

Five adoption records, each with an exercised test (imports alone never
count): `physmath.NewRational`, `artifact.ComputeRef`, `physmath.Apply`,
`physmath.EncodeRelationBody`, `physmath.Replay` — code paths in
`representation/`, tests `TestInfra*` in the package root; `manifest_test.go`
proves each recorded test exists *and* reaches the claimed symbol through
this module's call graph (§8.2.2).

Five justified deviations (§8.3), each with capability sought / API examined
/ reason / smallest local alternative: no axion framework (`DEV-FRAMEWORK`,
never misusing `gr`'s namespace), no executable analytic steps
(`DEV-ANALYTIC-STEPS`), no limit rules (`DEV-LIMIT`), no numerical layer
(`DEV-NUMERICS`), no unit engine (`DEV-UNITS`).

## 6. Conventions and honest limitations

- Conventions (manifest `conventions`): natural units ħ=c=1; exact coupling
  conversion 10⁻¹⁰ GeV⁻¹ = 10⁻¹⁹ eV⁻¹ (unit-tested); given-n dispersion
  q_γ = n ω_γ; free-photon measure cut; plasmon pole (A62); J(κ) closed form
  with endpoint limits J(0) = 8π/3, J(1) = 2π (the closed form is a
  0·∞ indeterminate exactly at κ=1 — handled by the documented limit).
- **Threshold gate**: the physical condition is E_a = m(1+v_a²/2) > ω_p,
  while the leading-order prefactor requires m − ω_p > 0. The exact gate is
  therefore m > ω_p; the razor-thin window ω_p/(1+v_a²/2) < m ≤ ω_p (width
  O(v_a² ω_p) ≈ 8 meV), where the exact kinematics is open but the NR
  prefactor is invalid, is excluded by the same gate. No requested case lies
  inside it. (This replaced an earlier gate formulation that could return a
  negative "rate" in that window — caught by `protocol.TestThresholdBehavior`.)
- Limitations (manifest `limitations[]`): replay ≠ premise truth;
  `prefactor_normalization` one-power-of-n ambiguity (G-11) stated, not
  resolved; NR E_a convention tested insensitive to < 10⁻⁵ relative; exact
  rationals use J(0) with the κ² correction (~4×10⁻⁶) below display
  precision; units are metadata only; prior protocol-level exposure to
  sibling structure disclosed; `phys-demo` is not a git repository (no
  revision pinning); no central Type-B validator exists — this package's
  `manifest_test.go` is the entire enforcement story (§11.3).
- **Gap-severity judgment** (auditable, stated because it decides status):
  a `mandatory` gap is a §11.3 conformance element unfulfillable by *this
  package* (schema field, entry point, evidence level, firewall). Capability
  absences that are correctly recorded through the schema's own required
  fields (`support` states, `deviations`, `limitations`) are `advisory` —
  §9.2.5/§11.4 direct that out-of-scope operations be recorded as boundary,
  not treated as missing features. All four declared gaps are advisory;
  `TestConformanceFlipsOnMandatoryGap` proves the computation still flips to
  NONCONFORMANT the moment any gap is mandatory, so the classification is
  enforced, not decorative.

## 7. Testing (PROTOCOL §12) — what is and is not established

Established: strict manifest schema validity (15 corruption negatives:
wrong schema_id, unknown/missing keys, missing unit, dimensionless-with-unit,
dangling test ref, forgotten reuse marker, invented rule ID, bad evidence
level/support state/`executed_by`, incomplete deviation, empty
`infrastructure_use`+`deviations`, firewall claim, mandatory-gap flip,
unexercised-entry-point flip); entry points called and observed; byte-exact
replay incl. corrupted-output and invented-rule-ID rejection; content-address
integrity incl. flipped-byte rejection; exact dual-path arithmetic; the
exact threshold gate on both sides of ω_p and inside the NR-invalid window;
unit-slip detection (slips ≥ 10⁸ relative); J(κ) series/endpoint/monotonic
checks.

Never established: physical truth, empirical support, correctness of the
hand-derived rate formula, correctness of source mapping.

## 8. Reproduction and baseline revisions

```text
cd phys-demo/problem_8a/mimo
go test -count=1 ./...        # 29 tests, all packages
go vet ./...                  # clean
gofmt -l .                    # empty
```

The root orchestrator `phys-demo/scripts/test-all.sh` keeps a hardcoded
project list; `mimo` is not yet in it (deferred out-of-package one-line
addition per the script's documented process, recorded in the manifest's
`open_obligations`). The package passes the script's exact per-project
invocation (`unset GOWORK && go test -count=1 ./...`); the three pre-existing
projects still pass under the script unchanged.

Infrastructure baselines at completion (`git rev-parse HEAD`):
`phys-math 7f3e0a7e6960c5a918e558d4b0d9715b13191a45`,
`phys-artifact 31c21fe93215b60238205ce7f5fb32ba29c195f2` (both clean work
trees). `phys-lib`/`phys-trace` were not imported (see `DEV-FRAMEWORK`;
`phys-lib` is not in this module's `go.work`). `phys-demo` itself is not a
git repository, so this package cannot be revision-pinned (recorded
limitation, decision item 12).
