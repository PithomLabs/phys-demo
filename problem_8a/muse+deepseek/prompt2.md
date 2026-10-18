Use this prompt for the implementation agent. It combines the material findings from both reviews and adds specific guidance for deriving the physics correctly and implementing the numerical formula robustly in Go. It keeps the scope narrow and does not treat the disclosed limitations as reasons to redesign Phys.

# Targeted Correction: Axion-to-Photon Conversion, TheoPhysEval Problem 8(a)

## Mission

Address the remaining material issues in `phys-trace/demo/problem_8a/` following the consolidated adversarial review.

The review found that the reported solution is substantially correct within its stated dielectric model, but identified:

1. An unresolved normalization derivation for the general rate formula.
2. A singular numerical edge case in `JClosed` at (\kappa=1).
3. A misleading `TestResponseSignSensitivity` test whose comment and assertions describe different behavior.
4. A Go test command that was run from the wrong workspace context.

Your task is to resolve these issues with a minimal, scientifically defensible change. Do not blindly trust either the current implementation or either review. Independently verify the relevant physics and actual code.

The required numerical results remain:

[\
\Gamma\_{\rm NR}(25,\mathrm{eV})\
\approx1.31\times10^{-41},\mathrm{eV},\
\qquad\
\Gamma\_{\rm NR}(15,\mathrm{eV})=0,\
]

subject to confirming the derivation and normalization.

**Do not alter the expected answer merely to make tests pass.** If independent derivation reveals an actual normalization error, correct the formula, numerical results, goldens, and documentation consistently, explaining the discrepancy.

## 1. Reconnaissance and operating boundaries

Before editing, read the applicable `AGENTS.md` instructions, `phys-system/AGENTS.md`, the normative `DERIVATION_SPEC_GUIDELINES.md`, and the demo's README, manifest, source files, and tests.

Inspect the live module layout and `go.mod` before running commands. The previous invocation of `go test ./...` from an incompatible workspace directory failed because the current directory was not within the modules selected by `go.work`. That is an invocation-context failure, not evidence that the demo's tests failed.

### Scope

Prefer changes confined to `phys-trace/demo/problem_8a/`.

Do not modify frozen `phys` or `phys-gr`. Do not add PhysMath rules, artifact families, shared schemas, namespaces, or general-purpose integration/evaluation infrastructure. Do not commit, push, tag, or publish.

Keep the distinction between hand-derived physics, existing PhysMath execution, and numerical Go evaluation explicit.

## 2. Theoretical-physics task: independently audit the rate normalization

The current README derives the rate using the interaction

[\
\mathcal L\_{\rm int}\
\=-\frac{g\_{a\gamma\gamma}}4 aF\_{\mu\nu}\widetilde F^{\mu\nu},\
]

and the prescribed longitudinal dielectric response

[\
\operatorname{Im}\left[-\frac1{\epsilon_L(\omega,q)}\right]\
\=\frac{\pi}{2}\omega_p\delta(\omega-\omega_p).\
]

The implemented general rate is

[\
\Gamma =\
\frac{g\_{a\gamma\gamma}^{,2}\omega_p\
q\_\gamma^2 k_a^2}\
{16\pi E_a\omega\_\gamma},J(\kappa),\
]

where

[\
\omega\_\gamma=E_a-\omega_p,\qquad\
q\_\gamma=n\omega\_\gamma,\qquad\
\kappa=\frac{k_a}{q\_\gamma},\
]

and

[\
J(\kappa)=2\pi\left[\
\frac{1+\kappa^2}{2\kappa^2}\
-\frac{(1-\kappa^2)^2}{4\kappa^3}\
\ln\left(\frac{1+\kappa}{|1-\kappa|}\right)\
\right].\
]

The issue raised in review is that the final (1/(16\pi)) prefactor is not fully reconstructed in an auditable chain in the current README.

### Required work

Independently reconstruct the phase-space and self-energy calculation, using the original interaction convention and the cited literature as sources rather than treating the code's formula as authoritative.

Inspect the existing citations, including:

- Caputo–Millar–Vitagliano, arXiv:2005.00078, for the self-energy-to-rate method.
- Knapen–Kozaczuk–Lin, arXiv:2101.08275, for the longitudinal loss-function/dielectric-response relationship.
- The source of the axion-photon vertex convention and its relation to the vacuum width.

Record the exact source locations and distinguish source-backed identities from independently derived steps.

Pay particular attention to:

1. **Vertex normalization.** Derive the vertex from the stated interaction convention. Check its polarization sum and verify the vacuum (a\to\gamma\gamma) width (g^2m_a^3/(64\pi)) under that convention.
2. **Self-energy normalization.** Establish the precise relationship between the disappearance rate and the imaginary part of the on-shell axion self-energy, including any factors of (2E_a) and sign conventions.
3. **Phase-space measure.** Write all integration measures explicitly. The energy-transfer variable (\omega) and the outgoing photon momentum are distinct integration variables before the delta function is applied. Do not omit a (d\omega) integral or its normalization if the derivation requires it.
4. **Dielectric spectral weight.** Trace the relation between the longitudinal propagator (D\_{00}=1/(q^2\epsilon_L)) and the absorptive loss function. Establish the sign and normalization of the factor used in the cut.
5. **Photon normalization in the medium.** The solution assumes an outgoing transverse photon with (|\mathbf q\_\gamma|=n\omega\_\gamma). Derive or explicitly source the photon phase-space measure and normalization under the stated constant-(n) approximation. Do not assume that the refractive-index dependence or the density-of-states Jacobian cancels without showing why.
6. **Delta-function collapse.** Show precisely how the energy-conserving delta function and the plasmon-pole response fix the energy transfer, establish the threshold factor, and leave the surviving factors.
7. **Angular integral.** Starting from the vertex contraction, show how the angular measure produces (J(\kappa)), including the factors of (2\pi).
8. **Assembled prefactor.** Combine the verified pieces step by step to establish or correct the coefficient (1/(16\pi)). Account for every phase-space factor and normalization convention exactly once.

Do not silently import a coefficient from a related process merely because the method is similar.

### Possible outcomes

- If the coefficient and formula are confirmed, retain them and make the derivation auditable.
- If the coefficient is wrong, correct the implementation and all affected independent goldens, tests, and reported numbers.
- If some normalization convention cannot be established from the available sources, explicitly classify that portion as source-backed but unresolved, state its effect on the result, and do not claim the rate is fully independently derived.

This does not require adding symbolic integration or tensor-contraction operations to PhysMath. Complete the physics analytically and formalize the result honestly.

## 3. Go implementation task: fix `JClosed` at (\kappa=1)

The current closed form has a removable singularity at (\kappa=1). The expression contains a factor that approaches zero multiplied by a logarithm that diverges. Direct floating-point evaluation can therefore produce `NaN`.

The correct limiting value is

[\
\boxed{J(1)=2\pi}.\
]

### Required fix

Add an explicit, well-documented limit branch in `axion.JClosed` for (\kappa=1), ensuring that the mathematical limit is returned rather than evaluating the indeterminate floating-point expression.

Review the numerical function's valid input domain at the same time:

- (\kappa=0) should continue to return (8\pi/3).
- (\kappa=1) should return (2\pi).
- Valid positive inputs away from the removable singularity should retain the existing formula.
- Non-finite or otherwise invalid inputs must not silently produce a plausible physical rate. If input validation is appropriate under the existing API contract, handle it narrowly and document the behavior.

### Numerical-stability guidance

A robust numerical implementation should consider more than the exact singular point.

The closed form contains subtraction of large terms for small positive (\kappa), and the logarithm can also lose accuracy when its argument is close to one. Assess whether this materially affects the supported input domain, including the benchmark value (\kappa\approx0.001984).

Use stable elementary-function evaluations or a small-(\kappa) series only if needed to prevent a meaningful loss of accuracy. Do not add a general symbolic engine or an elaborate numerical framework for this fix.

### Required tests

Add tests that verify:

1. `JClosed(0) == 8π/3` within the existing justified tolerance.
2. `JClosed(1) == 2π` within floating-point tolerance.
3. Values on both sides of (\kappa=1) remain finite and approach the correct limit.
4. The benchmark non-relativistic correction and headline number remain within the documented accuracy, unless the normalization audit justifies a correction.
5. If additional input guards are added, invalid values fail according to the explicitly documented contract.

Retain the existing quadrature tests. They are useful independent checks, but they do not replace an explicit test of the singular point.

## 4. Correct the response-sign test semantics

Review `TestResponseSignSensitivity`.

Its current setup negates `omega_p`, expects a negative rate, but describes the test as exercising a positivity gate. Those are conflicting intentions. Negative (\omega_p) is outside the stated physical benchmark domain.

Make the smallest honest correction:

- Do not claim that the existing test demonstrates a physical positivity gate.
- Verify the positive benchmark response normalization and positive rate for an open physical channel.
- If you retain a negative-(\omega_p) algebraic test, label it explicitly as an unphysical signed-input sensitivity test, not evidence that the physical model permits negative rates or that a positivity validator exists.
- Do not introduce a general validation framework solely to repair the test. Add a narrow invalid-parameter rejection only if it is already required by the existing API contract.

The intended scientific meaning is that the physical response model has positive plasmon frequency and the appropriate absorptive sign. A deliberately unphysical mutation should not obscure that meaning.

## 5. Preserve the numerical and scientific boundaries

Keep the threshold conclusion explicit:

[\
E_a=m_a\left(1+\frac{v_a^2}{2}\right),\
\qquad\
\Gamma\_{\rm NR}=0 \quad\text{when }E_a\le\omega_p\
]

under the stated (T=0) response model.

For the specified parameters, the 15 eV channel is below the 16.6 eV plasmon pole. It is zero by energy conservation, not by a computed cancellation.

For the 25 eV case, retain the exact-rational skeleton and the numerical angular correction as separate evidence types. Recompute the result independently after changing the general formula if the normalization audit requires it.

Preserve the model limitations: the delta-function plasmon pole, constant refractive index, neglected recoil, and omitted transverse medium response do not define a universal rate for every real dielectric.

Do not claim that the implementation has replayed a complete PhysMath derivation. The analytical vertex, spectral identity, delta-function reduction, and phase-space integration remain outside the current closed rule registry unless you demonstrate otherwise from the actual implementation.

## 6. Test execution: use the correct Go module context

The previous command

`go test ./...`

failed with:

`pattern ./...: directory prefix . does not contain modules listed in go.work or their selected dependencies`

The demo README says this module should be tested with `GOWORK=off`. Verify the live module configuration and run the tests from the demo module root, for example:

```
cd phys-trace/demo/problem_8a
GOWORK=off go test ./...
```

Run `gofmt` and `go vet` in the correct context as well.

Do not report the earlier workspace-root invocation as a passing test. If the demo uses any local dependencies or requires another test command, follow the actual module setup and document it.

## 7. Scope and change control

Keep changes limited to the axion demo and its relevant documentation/tests unless a genuine dependency constraint prevents that.

Do not:

- Modify frozen `phys` or `phys-gr`.
- Add PhysMath rules, families, schemas, namespaces, or shared abstractions.
- Treat the benchmark statement as a source for the solution.
- Change a golden solely to make a test green.
- Hide unresolved formula normalization behind a test pass.
- Commit, push, tag, merge, or publish without explicit authorization.

Optional strengthening: an assembled ultrarelativistic-rate regression test may be useful if it adds meaningful coverage beyond the existing (J(2/3)) quadrature test. It is not a mandatory blocker-resolution requirement unless the review uncovers a specific material inconsistency in the general branch.

## 8. Completion report

Report:

1. Exact files changed.
2. What the independent physics audit establishes about the full prefactor.
3. The key steps linking vertex normalization, spectral response, phase-space integration, delta collapse, angular integral, and final rate.
4. The treatment of the (\kappa=1) limit and any other materially necessary numerical stability fixes.
5. How `TestResponseSignSensitivity` was corrected or reclassified.
6. The requested numerical results after any independently justified correction.
7. Exact test commands and actual results from the correct module context.
8. Which stages remain analytical or represented-only, and why.
9. Confirmation that no frozen-system or architecture changes occurred and no unauthorized Git operations were performed.

## Acceptance criteria

The task passes when:

- The overall rate normalization is derived, correctly sourced, or explicitly left conditional with its unresolved dependency identified.
- The exact (\kappa=1) limit no longer produces `NaN`, and the function behaves correctly over the tested valid domain.
- The response-sign test matches what it actually verifies.
- The reported benchmark values follow from the independently checked physical formula.
- Tests are run from the correct Go module context and their actual outcomes are reported.
- Model assumptions, analytical work, numerical results, and Phys execution remain distinct.

Use expert physics judgment to resolve the physical factors and expert Go judgment to implement the numerical behavior robustly. Prefer a small correct fix over a broad redesign, and prefer an explicitly unresolved scientific issue over an unsupported declaration of success.

My recommendation: treat the normalization audit as the most scientifically important part, and the \\(\kappa=1\\) fix as the clearest code defect. The sign-test correction is small but necessary for the test suite to mean what it claims. The agent should not declare the whole problem fully verified merely because the updated Go tests pass.