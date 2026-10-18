# Final plan — Problem 8(a) targeted correction

## Progress (updated as work lands; `[x]` = done and verified)

- [x] Task 1 — `go.mod` added (sanitized module path), stale import fixed,
      `gofmt`/`vet`/`test` green from the demo root with `GOWORK=off`.
- [x] Task 2 — normalization audit: derivation reconstructed; prefactor
      confirmed under the documented convention; residual `n`-power
      convention dependence quantified and documented; assembly checked
      numerically (`TestPrefactorAssembly`).
- [x] Task 3 — `JClosed(1) = 2pi` branch; domain guards (negative /
      non-finite kappa cannot yield a plausible rate).
- [x] Task 4 — tests updated (kappa = 1 limit, both-sides approach,
      benchmark kappa in quadrature, response-sign semantics, guard tests,
      assembly check).
- [x] Task 5 — README / `doc.go` / manifest consistency
      (path, threshold wording, convention statement, limitation entry).
- [x] Task 6 — completion report (final message) with exact commands and
      results.

Verification command (reproduces the suite from the demo root):

    GOWORK=off go test ./...

Target: `/home/chaschel/Documents/go/phys-demo/problem_8a/muse+deepseek`
(the physical directory keeps this name).

Canonical repository: `github.com/PithomLabs/phys-demo`. PhysTrace continues
to own the normative derivation specification
(`https://github.com/PithomLabs/phys-trace/blob/main/DERIVATION_SPEC_GUIDELINES.md`).
The earlier `phys-trace/demo/problem_8a/` phrasing in the mission prompt is
superseded by the `phys-demo` ownership decision; historical input records
(`prompt.md`, `prompt2.md`, `protocol/adv_review.md`) are left unchanged.

## Operating constraints

- Scope: only files under this demo directory. No changes to frozen `phys` or
  `phys-gr`; no PhysMath rules, families, schemas, namespaces, or shared
  abstractions; no commits, pushes, tags, or merges.
- Keep explicit the distinction between hand-derived physics, Go numerical
  evaluation, and PhysMath execution (none is replayed).
- Do not change the headline number on an unverified hand derivation; do not
  hide a genuinely unresolved convention behind a passing test.

## Task 1 — Module context (enabling step)

- Add `go.mod` in the demo root (directory name stays `muse+deepseek`; the
  `+` character is invalid in Go module paths, so the module path uses the
  sanitized segment `muse-deepseek`):

  ```
  module github.com/PithomLabs/phys-demo/problem_8a/muse-deepseek

  go 1.25.7
  ```

  Stdlib only; no edits to the parent `go.work`; no external dependencies; no
  `go.sum` needed.
- Fix the stale internal import in `protocol/derivation_test.go` to
  `axion "github.com/PithomLabs/phys-demo/problem_8a/muse-deepseek/axion"`.
- Run all checks from the demo root with `GOWORK=off`:
  `gofmt -l .`, `go vet ./...`, `go test ./...`. Report actual output.
- Do not report the earlier workspace-root invocation as a passing suite.

## Task 2 — Normalization audit (scientific core)

Produce an auditable derivation of the full rate prefactor in README.md (and
code comments where they carry claims). Work through, in order:

1. **Anchor (re-confirm in the write-up).** From `L = -(g/4) a F Fdual`:
   `M(a->gg) = i g eps^{mu nu alpha beta} k1_mu eps1_nu k2_alpha eps2_beta`,
   `sum_T |V|^2 = 2 g^2 (k1.k2)^2`, vacuum width `g^2 m^3 / 64pi`;
   `sum_T |V0|^2 = g^2 |q x p|^2`; `J(0) = 8pi/3`, `J(1) = 2pi`.
2. **Cut -> rate bookkeeping.** Fix one convention for the self-energy cut and
   the `(2pi)` factors, and calibrate it against the vacuum width *in the same
   notation*. This is where the assembled `1/(16pi)` can slip.
3. **Medium spectral weight.** Re-derive `-2 Im D00 = (2/p^2) Im[-1/eps_L]`
   from the screened Coulomb propagator + FDT; label which part is
   source-backed (KKL/Jackson) and which is independently derived.
4. **Photon phase space and mode normalization with `|q| = n w` — the crux.**
   Investigate first, as instructed: attempt the cited literature
   (Caputo–Millar–Vitagliano, arXiv:2005.00078; Knapen–Kozaczuk–Lin,
   arXiv:2101.08275) via arXiv if network permits. If not, derive with an
   explicitly stated convention and cross-check against an exactly solvable
   toy medium (single-pole quantized oscillator, computed microscopically and
   via the eps_L formula) to pin the medium-side and phase-space factors.
   Determine the `n`-power explicitly: whether the `d^3q/((2pi)^3 2w)` measure
   picks up factors of eps/n.
5. **Delta collapse and angular integral.** Show the `dw` integration
   normalization explicitly (currently omitted), the pole collapse
   `w = w_p` and threshold `Theta(Ea - w_p)`, and `I_Omega = ka^2 J(kappa)`.
6. **Assemble** the prefactor step by step and decide with evidence:
   - **Confirmed** -> keep `1.31e-41 eV`, make the full derivation auditable.
   - **Definitively wrong** -> correct formula, numerical result, independent
     goldens, tests, manifest, and documentation consistently, explaining the
     discrepancy.
   - **Model under-specifies the in-medium photon normalization** -> retain
     the result explicitly as *conditional on the stated photon normalization
     convention* (option (a)), document the unresolved dependence prominently.
     Present an alternative numerical value only if its convention and
     derivation are independently justified.

Do not silently import a coefficient from a related process; account for every
phase-space factor exactly once.

## Task 3 — `JClosed` at kappa = 1 (`axion/axion.go`)

- Add an explicit, well-documented limit branch: `kappa == 1 -> 2pi`
  (removable singularity; never evaluate the indeterminate form).
- Keep `kappa == 0 -> 8pi/3`; keep the closed form elsewhere.
- Physical domain: `kappa = ka / q_gamma >= 0` (a ratio of momentum
  magnitudes). Document it and enforce it narrowly; do not add generality for
  negative kappa. Negative or non-finite input must never yield a plausible
  rate.
- Narrow guards: `GammaGeneral` must not return `ok == true` with a NaN rate;
  non-finite inputs return `(NaN, false)` under an explicitly documented
  contract.
- Assess small-kappa cancellation honestly for the benchmark
  `kappa ~ 0.001984` (each term ~1/(2 kappa^2) ~ 1.3e5, difference ~1.33;
  loss ~5 digits, well inside existing tolerances). Use `log1p`-based
  evaluation only if it materially helps; otherwise document why not.

## Task 4 — Tests (`protocol/derivation_test.go`)

- Keep the existing quadrature tests.
- Add: `JClosed(0) == 8pi/3` (exists; keep), **`JClosed(1) == 2pi`**, and
  one-sided checks `JClosed(1 +/- delta)` finite and approaching `2pi` (the
  approach is first order in `delta`: `J = 2pi(1 -+ delta) +
  O(delta^2 ln(1/delta))`, so bound the difference by `4 pi delta`); include
  the benchmark `kappa` and `kappa = 1` in the quadrature check (with the
  removable point handled by its limit).
- Fix `TestResponseSignSensitivity` semantics: the physical part verifies
  `ResponseStrength == (pi/2) w_p > 0` and an open-channel positive rate; the
  `w_p < 0` case is relabeled explicitly as an **unphysical signed-input
  algebraic sensitivity check**, with a comment stating that no positivity
  validator exists or is claimed. Comment and assertions must agree.
- Guard tests: negative / NaN kappa yields NaN, never a plausible rate;
  `GammaGeneral` with non-finite input returns `ok == false` and no plausible
  numeric rate.
- Optional (non-blocking): assembled ultrarelativistic consistency check at
  `kappa = 2/3` against Simpson `J(2/3)`.

## Task 5 — Documentation and manifest consistency

- README: auditable derivation section; explicit statement of the in-medium
  photon normalization convention; threshold statement — "kinematically
  closed and the rate is zero when `Ea <= w_p`; above threshold the rate is
  positive for the stated model when the coupling and axion momentum are
  nonzero, subject to the adopted photon-normalization convention" (no
  "iff"); limitations unchanged or strengthened (delta-pole response,
  constant `n`, no transverse medium response, recoil neglected, no PhysMath
  replay).
- Path references: update `axion/doc.go` ("living under
  phys-demo/problem_8a/muse+deepseek") and the README `GOWORK=off` note (this
  standalone demo module is not listed in the parent workspace file).
- `manifest.json`: add a normalization limitation entry (e.g.
  `prefactor_normalization`) describing the audited or conditional status;
  keep claim IDs/statuses compatible with `manifest_test.go` unless the audit
  forces a documented reclassification (then update the test consistently).
- Leaves `muse/`, `inkling/`, and the parent `go.work` untouched.

## Task 6 — Completion report

Report:
1. Exact files changed, including the repository/move correction.
2. What the independent physics audit establishes about the full prefactor
   (confirmed / corrected / conditional with the specific unresolved
   dependency).
3. Key derivation steps linking vertex, spectral response, phase space,
   delta collapse, angular integral, prefactor.
4. kappa = 1 treatment and any other materially necessary numerical-stability
   fixes.
5. How `TestResponseSignSensitivity` was corrected or reclassified.
6. Final numerical results (unchanged or corrected).
7. Exact commands and actual results from the demo root (`GOWORK=off`).
8. Which stages remain analytical or represented-only, and why.
9. Confirmation of no frozen-system or architecture changes and no
   unauthorized Git operations.

## Acceptance criteria

- Normalization derived, correctly sourced, or explicitly conditional with
  the unresolved dependency identified.
- `kappa = 1` no longer produces NaN; function correct across the tested
  domain; invalid inputs cannot produce a plausible rate.
- Response-sign test matches what it verifies.
- Benchmark values follow from the independently checked physical formula.
- Module declaration, internal imports, and live directory structure agree;
  tests run from the correct module context; outcomes reported truthfully.
- Model assumptions, analytical work, numerical results, and (absent) Phys
  execution remain distinct.

Governing rule: resolve the physics where possible; preserve a conditional
result where necessary; never turn an unresolved convention into false
certainty.
