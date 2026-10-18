# Problem 8(a): Axion-to-Photon Conversion in a Dielectric — gemini

**Type-B benchmark-solution package** under `phys-demo/PROTOCOL.md` by Antigravity using Claude 4.6 Thinking and Gemini 3.8 Flash Medium

## Question

Compute the non-relativistic rate $\Gamma_{\rm NR}$ for axion conversion
to photons in a dielectric medium, for $m_a = 25\,{\rm eV}$ and
$m_a = 15\,{\rm eV}$, given the plasmon-pole response model (A62).

## Result

$$
\boxed{
\begin{aligned}
\Gamma_{\rm NR}(m_a=25\,{\rm eV}) &\approx 1.31 \times 10^{-41}\,{\rm eV}, \\
\Gamma_{\rm NR}(m_a=15\,{\rm eV}) &= 0 \quad\text{(threshold-closed)}.
\end{aligned}
}
$$

Ordered: $[\,1.31\times10^{-41}\,{\rm eV},\; 0\,]$.

## Derivation

### Interaction and vertex

The axion-photon Lagrangian $\mathcal{L} = -(g/4)\,a\,F_{\mu\nu}\tilde{F}^{\mu\nu}$
gives an $a\gamma\gamma$ vertex with amplitude $\propto g\,\varepsilon^{\mu\nu\rho\sigma}
k_{1\rho}k_{2\sigma}\,\epsilon_{1\mu}\epsilon_{2\nu}$. Squaring and
summing over polarisations, the matrix element produces
$|{\cal M}|^2 \propto g^2\,|{\bf q}\times{\bf p}|^2$ for spatial
momentum vectors ${\bf q}$ (outgoing photon) and ${\bf p}$ (medium
excitation). This is verified independently by Levi-Civita contraction
(numeric, 12 digits agreement plus collinear-zero check).

### Self-energy / cutting route

The conversion rate is extracted from the axion self-energy via the
optical theorem:

$$
\Gamma = -\frac{\operatorname{Im}\Pi(\omega)}{\omega}\bigg|_{\rm on\text{-}shell}.
$$

The cut of the self-energy loop produces one outgoing transverse photon
(free-photon spectral cut) and one medium excitation (longitudinal
response). The medium enters through the dressed longitudinal propagator
$D_{00} = 1/(q^2\,\epsilon_L)$, with the absorptive weight

$$
-2\operatorname{Im}D_{00}
= \frac{2}{|{\bf p}|^2}\operatorname{Im}\!\left[-\frac{1}{\epsilon_L}\right].
$$

### Rate template (calibrated on vacuum)

The rate template, calibrated against the vacuum $a\to\gamma\gamma$
width $g^2 m^3/(64\pi)$, is

$$
\Gamma = \frac{1}{2E_a}\,\frac{1}{2}\!\int\!\frac{d^4p}{(2\pi)^4}\,
[2\pi\delta(q^2)]\,[-2\operatorname{Im}D_{00}(p)]\,g^2|{\bf q}\times{\bf p}|^2.
$$

### Delta collapse

With (A62), the $\omega$ integral sets $\omega = \omega_p$,
$\omega_\gamma = E_a - \omega_p$, $\Theta(E_a - \omega_p)$. The model
normalisation satisfies the Bethe sum rule
$\int\omega\operatorname{Im}[-1/\epsilon_L]\,d\omega = (\pi/2)\omega_p^2$
exactly.

### Angular integral

With ${\bf p} = {\bf k}_a - {\bf q}_\gamma$,
$|{\bf q}_\gamma \times {\bf p}|^2 = q_\gamma^2 k_a^2 \sin^2\theta$,
and

$$
I_\Omega = \int\!d\Omega\,\frac{|{\bf q}_\gamma\times{\bf p}|^2}{p^2}
= k_a^2\,J(\kappa),
$$

with $\kappa = k_a/q_\gamma$ and

$$
J(\kappa) = 2\pi\!\left[\frac{1+\kappa^2}{2\kappa^2}
- \frac{(1-\kappa^2)^2}{4\kappa^3}\ln\frac{1+\kappa}{|1-\kappa|}\right].
$$

Special values: $J(0) = 8\pi/3$, $J(1) = 2\pi$ (removable singularity).
Checked against numerical quadrature to 8 decimals.

### Assembly

Combining all factors with $\omega_\gamma = E_a - \omega_p$,
$q_\gamma = n\,\omega_\gamma$:

$$
\boxed{
\Gamma = \frac{g^2\,\omega_p\,q_\gamma^2\,k_a^2\,J(\kappa)}{16\pi\,E_a\,\omega_\gamma}.
}
$$

### NR specialisation

For $v_a = 10^{-3}$: $E_a = m(1 + v^2/2)$, $k_a = mv$,
$\kappa \approx 2\times10^{-3}$, $J(\kappa) \to 8\pi/3$:

$$
\Gamma_{\rm NR} = \frac{g^2\,\omega_p\,q_\gamma^2\,k_a^2}{6\,E_a\,\omega_\gamma}.
$$

This skeleton is evaluated in exact-rational arithmetic (`math/big.Rat`)
with a bounded float64 correction $C = J(\kappa)/(8\pi/3)$,
$|C - 1| < 10^{-5}$ at benchmark kinematics.

### Assembly audit (prefactor and in-medium photon normalisation)

The cutting route above is calibrated against the vacuum width in
the same notation, so every phase-space factor is accounted for
exactly once.

- *Template, calibrated on vacuum.* In vacuum,
  $[2\pi\delta(q^2)][2\pi\delta(p^2)]$ with the identical-line factor
  $1/2!$ gives $\Gamma(a\to\gamma\gamma) = g^2 m^3/64\pi$: the $1/2$
  and the $1/(2\pi)^4$ are fixed by the exact vacuum width.
- *Free-photon cut (the stated convention).* The outgoing transverse
  photon is cut as the free massless density with the given-$n$
  dispersion:
  $2\pi\delta(q^2)|_+ = (\pi/|{\bf q}|)\,\delta(E_a - p^0 - |{\bf q}|/n)$.
- *Medium weight.*
  $-2\operatorname{Im}D_{00} = (2/|{\bf p}|^2)(\pi/2)\omega_p\,
  \delta(p^0 - \omega_p)$.
- *Collapse.* $p^0$ collapses the medium delta; radial $q$ gives
  $n\,q_\gamma$; angular integral gives $2k_a^2 J(\kappa)$. Assembly:

$$
\Gamma = \frac{1}{2E_a}\frac{1}{32\pi^4}\,g^2\,\frac{\pi^2\omega_p}{2}\,
4\pi n\,q_\gamma\cdot 2k_a^2 J(\kappa)
= \frac{g^2\omega_p\,n\,q_\gamma\,k_a^2\,J(\kappa)}{16\pi E_a},
$$

which is exactly the boxed form ($n\,q_\gamma = q_\gamma^2/\omega_\gamma$).
`TestPrefactorAssembly` checks the end-to-end assembly numerically
(Simpson quadrature, relative agreement $< 10^{-9}$).

- *Residual convention dependence (one power of $n$).* Relative to
  the adopted convention (rate $\propto n^2$ here), standard
  alternatives give $0.67\times$ to $1.5\times$ the quoted value, i.e.
  $[\,8.7\times10^{-42},\;1.3\times10^{-41},\;2.0\times10^{-41}\,]$
  eV. The model does not by itself fix this power; it is the manifest
  limitation `prefactor_normalization`.

## Parameter roles

- $g_{a\gamma\gamma}$: overall $g^2$; GeV$\to$eV is $10^{-9}$
  ($g = 10^{-19}\,{\rm eV}^{-1}$); a GeV-as-eV slip moves the rate
  by exactly $10^{18}$ (tested).
- $\omega_p$: pole position (threshold) and linear prefactor.
- $n$: photon momentum $q_\gamma = n\omega_\gamma$ (phase space), and the
  convention-dependent spectral-density power (see Assembly audit).
- $v_a$: enters only via $k_a^2$ — quadratic, essential.
- $m_a$: sets $E_a, k_a$ and hence which side of threshold.

## Capability audit

| Capability | State | PhysMath API | Verdict |
|---|---|---|---|
| Exact rational arithmetic | SUPPORTED | `physmath.NewRational` (expression IR node) | DEVIATION: IR nodes, not evaluable scalars; stdlib `big.Rat` used instead |
| Expression construction | SUPPORTED | `physmath.PutExpression`, Expr AST | Not needed: solver does arithmetic, not symbolic representation |
| Differentiation | SUPPORTED | `DIFFERENTIATE-001` | Not needed for this problem |
| Simplification | PARTIALLY SUPPORTED | `EXPR-SIMPLIFY-001` (closed rewrite list) | Not needed |
| Delta-distribution calculus | UNAVAILABLE | No rule; `integral/1` REPRESENTATIONAL ONLY | DEVIATION: hand derivation |
| Tensor contraction | UNAVAILABLE | `tensor/1` exists, no contraction rule | DEVIATION: hand derivation |
| Phase-space integration | UNAVAILABLE | `integral/1` REPRESENTATIONAL ONLY | DEVIATION: closed form + quadrature |
| Axion-electrodynamics framework | UNAVAILABLE | `gr.NewFramework` rejects non-GR IDs | DEVIATION: demo-local records |
| Numerical execution | DEFERRED | Arch §25, Spec §23.2 | Solver-local float64 |

## Executed vs represented-only

- **Executed (exact rational Go):** parameter substitution, NR
  skeleton, threshold gate, ordering.
- **Evaluated numerically (float64, documented):** $J(\kappa)$,
  correction factor, final decimal, end-to-end assembly quadrature
  (`TestPrefactorAssembly`), domain-guard behaviour.
- **Analytical hand derivation (represented as records + this
  file):** vertex algebra, cutting/spectral identity, delta-collapse,
  phase-space assembly and its normalisation audit (including the
  convention statement and the residual one-power-of-$n$ dependence),
  angular integral.
- **Unsupported by the PhysMath rule registry:**
  $\delta$-distribution calculus, tensor contraction with a
  medium projector, phase-space integration, spectral-function
  identities. No PhysMath derivation is claimed; per the
  growth gate no new rule is added for one demo.

## Independent checks (analysis record)

1. Vertex normalisation vs vacuum width: exact.
2. Tensor structure $|{\bf q}\times{\bf p}|^2$: numeric, 12 digits + collinear zero.
3. $J(\kappa)$ closed vs quadrature: 8 decimals, 9 points (including
   $\kappa = 1$ and the benchmark $\kappa \approx 0.001984$).
4. Dimensions: $\lambda$-scaling gives $\lambda^3$ at fixed $g$ (rate is mass-dimension 1).
5. Sum rule: model saturates Bethe exactly.
6. Threshold continuity: rate $\to 0$ smoothly ($\sim q_\gamma^4$); pole-mass case open only via kinetic energy.
7. Assembly audit: cut template calibrated on the vacuum width; end-to-end
   assembly integral vs closed form, relative agreement $< 10^{-9}$
   (`TestPrefactorAssembly`).
8. $\kappa = 1$ removable singularity: exact branch $2\pi$; both-sided
   approach finite and first-order close (`TestJClosedKappaOneLimit`).
9. Domain guards: negative or non-finite $\kappa$ returns NaN; out-of-domain
   rate inputs return `(NaN, false)`; the kinematic gate still returns
   exactly `(0, false)` (`TestDomainGuards`).

## Growth-gate note

No new PhysMath operation, family, schema, or shared abstraction
is introduced. Physical assumptions live in demo-local records
because no `phys-lib` axion-electrodynamics domain package
exists; creating one would change shared contracts and requires
authorisation (proposed, not implemented). This demo is a
standalone Go module
(`github.com/PithomLabs/phys-demo/problem_8a/gemini`,
stdlib only).

## Limitations

- The $q$-independent pole response is given; real materials
  have continua, gaps, dispersion, local-field effects.
- In-medium photon normalisation: the prefactor is audited and
  confirmed under the stated convention (README "Assembly audit"),
  but the model does not fix the spectral-density normalisation of
  the in-medium photon; the residual dependence is one power of $n$
  (standard alternatives give $0.67\times$ to $1.5\times$ the quoted
  value). The quoted result is conditional on the stated convention
  (manifest limitation `prefactor_normalization`).
- The 15 eV zero is kinematic within the model.
- Transverse medium response absent by model (background
  literature supports the suppression).
- $n = 1.5$ constant; recoil neglected; NR truncation remainder
  $O(v^2) \sim 10^{-6}$ relative.
- No empirical claims; benchmark solutions appendix not read.
