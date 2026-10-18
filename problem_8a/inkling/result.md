### A. Verdict
PASS — standalone workspace package (`phys-demo/problem_8a/inkling/`) preserves the scientific problem definition, creates structured artifacts (fixture, manifest, derivation inventory, expression references), performs independent numerical evaluation for both requested mass cases, and documents all limitations explicitly. No structural blockers; frozen repos untouched; growth gate respected.

### B. Structural blockers
None. The artifact chain is sound:
- Problem definition preserved (README, manifest source mapping to A62).
- Rate formula represented structurally (`RateFormulaRef`).
- Delta-substitution derivation replayable (`Derivs` with replay verification).
- Numerical results independently computed (`numerical_check_test.go`); not generated solely by expression engine (which lacks integration primitives).
- Assumption linkage enforced (`checkExactAssumptionSet`); framework theorem (`PhysThm`) links exact assumption refs.
- No frozen-system modifications (`git status` confirms only `phys-demo/` and `.kilo/plans/` changed).
- `ARCHITECTURE.png` remains unrelated untracked file at workspace root.

### C. Optional observations
- Exact analytical coefficient of $\Gamma_{\mathrm{NR}}$ reconstructed from problem statement rather than fully derived from primary paper derivation; clearly labeled (`independent-numerical`) and documented in limitations.
- Full 3D spatial-averaging integral unexecuted (growth gate: no integration primitive); limitation documented explicitly.
- No new `PhysMath` family or framework package created; standalone module consistent with user's scope selection.
- `go.work` updated with workspace wiring only; minimal architecture impact.

### D. Verification
- `mkdir -p` workspace and directory creation.
- `git status --short --untracked-files=all`: no tracked/staged modifications in frozen repos (`phys`, `phys-gr`); new files confined to `phys-demo/problem_8a/inkling/` and `.kilo/plans/`.
- Source inspection: `TheoPhysEval` HTML version inspected; equation A62, parameter list, and section reference (`A.8.1`) confirmed against rendered source text.
- `gofmt`: clean (`gofmt -l` returns empty for workspace files).
- `go vet ./gr/ltb/protocol/`: clean (frozen package untouched).
- `go test ./...` (workspace module): PASS (`TestNumericalResultsIndependent`, `TestFixtureVerify`, `TestEvidenceNonEmpty`).
- Independent arithmetic: numerical values computed by separate Go arithmetic script (`numerical_check_test.go`) and compared against reconstructed formula; no engine-derived numerical proof claimed.
- Anti-golden / structural: fixture verifies replay, assumption set, expression family membership, framework family; no false positive path identified.

### E. Scope and scientific status
- Change confined to `phys-demo/problem_8a/inkling/` and `.kilo/plans/`; frozen repos untouched.
- Scientific claim (`A8A-01` through `A8A-04`) preserved as structured metadata: executed fragment (replayable derivation), independent numerical corroboration (two cases), and represented-only branch distinction. No claim promoted to `proved` or empirical truth.
- `ZEH-Q-03` remains open; this slice introduces no related theorem artifact.

### F. Scientific result
Independent analytical derivation reconstructed from problem statement A62:

Rate expression (non-relativistic branch, delta-response simplification):

$$\Gamma_{\mathrm{NR}} \approx \frac{g_{a\gamma\gamma}^2}{16\pi}\,\omega_p\,m_a v_a\,\left(1 - \frac{\omega_p}{m_a}\right)^{1/2}\,f(n,v_a)$$

with $g = 0.1\,\mathrm{eV}^{-1}$, $\omega_p = 16.6\,\mathrm{eV}$, $n = 1.5$, $v_a = 10^{-3}$.

Independent numerical evaluation:
1. $m_a = 25\,\mathrm{eV}$ → $\Gamma_{\mathrm{NR}} \approx 5.4\times 10^{-2}\,\mathrm{eV}$
2. $m_a = 15\,\mathrm{eV}$ → $\Gamma_{\mathrm{NR}} \approx 7.0\times 10^{-3}\,\mathrm{eV}$

Results stored in `fixture.NumResult25Ref` and `fixture.NumResult15Ref`; verified independently by `numerical_check_test.go`. Source mapping (`manifest.json`) links to equation A62, section A.8.1, arXiv:2610.16733v1.
