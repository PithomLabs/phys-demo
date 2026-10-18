# Problem 8(a): Axion-to-Photon Conversion in a Dielectric Medium

Note: 

1. meta/muse-spark-1.3-contributor first pass (no assistance from any AI)
2. deepseek/deepseek-v4.1-flash second pass (after adv_review.md, AI-assisted by ChatGPT)
3. normalization-audit pass (plan.md Tasks 1–5: module context, assembly
   audit with explicit photon-measure convention, $\kappa=1$ limit,
   domain guards, test semantics)

Benchmark: TheoPhysEval, arXiv:2610.16733v1, App. A Prob. 8(a) (A.8.1),
eq. (A62). The benchmark is the **problem definition only**; its
solutions appendix was not consulted.

## Result

For an on-shell axion $(E_a,{\bf k}_a)$ in the stated $T=0$ dielectric
with plasmon-pole response
$\operatorname{Im}[-1/\epsilon_L] = \frac{\pi}{2}\omega_p\,
\delta(\omega-\omega_p)$:

$$
\boxed{
\Gamma(E_a,{\bf k}_a) =
\frac{g_{a\gamma\gamma}^2\,\omega_p\,q_\gamma^2\,k_a^2}
{16\pi\,E_a\,\omega_\gamma}\,
J(\kappa)\;\Theta(E_a-\omega_p),
}
\qquad
\begin{aligned}
\omega_\gamma &= E_a-\omega_p,\\
q_\gamma &= n\,\omega_\gamma,\\
\kappa &= k_a/q_\gamma,\\
J(\kappa) &= 2\pi\!\left[
\frac{1+\kappa^2}{2\kappa^2}
-\frac{(1-\kappa^2)^2}{4\kappa^3}
\ln\!\frac{1+\kappa}{|1-\kappa|}\right].
\end{aligned}
$$

Non-relativistic specialization ($J(0)=8\pi/3$ exactly,
$E_a=m(1+v^2/2)$, $k_a=mv$):

$$
\boxed{
\Gamma_{\rm NR} =
\frac{g_{a\gamma\gamma}^2\,\omega_p\,q_\gamma^2\,k_a^2}
{6\,E_a\,\omega_\gamma}
\quad\text{(exact rational skeleton)}\quad
\times\underbrace{\frac{J(\kappa)}{8\pi/3}}_{|C-1|<10^{-5}}
}
$$

Requested numbers ($g=10^{-10}\,{\rm GeV}^{-1}$,
$\omega_p=16.6\,{\rm eV}$, $n=1.5$, $v_a=10^{-3}$),
in order $[(1),(2)]$:

$$
\boxed{
[\,\Gamma_{\rm NR}(25\,{\rm eV}),\;
\Gamma_{\rm NR}(15\,{\rm eV})\,]
=
[\,1.31\times10^{-41}\,{\rm eV},\; 0\,].
}
$$

Under the zero-temperature plasmon-pole model the channel is
**kinematically closed and the rate is zero when
$E_a\le\omega_p$**; above threshold the rate is positive for the
stated model when the coupling and axion momentum are nonzero,
subject to the adopted photon-normalization convention (see
"Assembly audit" below). For the 15 eV case this is a kinematic
result, not a computed cancellation.

## Sources and fidelity

| ID | Source | Role | Fidelity |
|---|---|---|---|
| benchmark | arXiv:2610.16733v1, A.8.1, (A62) | problem, parameters, response model, output order | problem-definition |
| caputo | Caputo–Millar–Vitagliano, arXiv:2005.00078, Sec. II, (2)–(4), App. A | self-energy $\to$ rate template; narrow-width $\delta$ pattern; vertex-factor care | method |
| kkl | Knapen–Kozaczuk–Lin, arXiv:2101.08275, Sec. II, (1), (3)–(6), (13) | density response $\leftrightarrow$ loss-function identity; loss function as inclusive excitation rate | response |
| jackson | Jackson, Classical Electrodynamics, 3rd ed., energy-loss chapter | $D_{00}=1/(q^2\epsilon_L)$; loss-function reading | identity |
| liang | Liang et al., arXiv:2608.20715 | transverse medium response 4–6 orders down | background |

Literature was located and fetched directly (arXiv abstract pages,
ar5iv/arXiv HTML renditions, and web search), and the quoted
identities were checked against the fetched text during the
normalization-audit pass. No answer key was consulted.

## Assumptions (see `axion.Assumptions()` for the closed set)

Late-time language aside: $T=0$ non-magnetic solid dielectric, no
external electric field, ${\cal L}=-(g/4)aF\tilde F$ only coupling,
$\hbar=c=1$, longitudinal pole response (A62) with no transverse
medium response in the model, outgoing transverse photon with
$|{\bf q}|=n\omega$ in the given-$n$ approximation, medium recoil
energy neglected ($\sim 10^{-9}$ eV), $T=0$ detailed balance
($\to\Theta(\omega)$), NR kinematics $E=m(1+v^2/2)$,
$k_a=mv$ (remainder vs the exact on-shell pair $O(v^2)\sim10^{-6}$
relative in the rate).

## Derivation

**Vertex.** From ${\cal L}=-(g/4)aF_{\mu\nu}\tilde F^{\mu\nu}$,
$V=-ig\,\epsilon_{\mu\nu\alpha\beta}k_1^\mu\varepsilon_1^\nu
k_2^\alpha\varepsilon_2^\beta$ (all-incoming). Validated:
$\sum|V|^2 = 2g^2(k_1\!\cdot\! k_2)^2$ reproduces the vacuum width
$\Gamma(a\to\gamma\gamma)=g^2m_a^3/64\pi$ (checked numerically).

**Cutting route (method after Caputo).** Disappearance rate from
the on-shell self-energy, $\Gamma=-\operatorname{Im}\Pi_a/E_a$.
Cut one photon leg as the outgoing real transverse photon, dress
and cut the other with the medium: with Coulomb-gauge
$D_{00}=1/(p^2\epsilon_L)$,

$$
\Gamma = \frac{1}{2E_a}\!\int\!
\frac{d^3q_\gamma}{(2\pi)^3 2\omega_\gamma}\,
(2\pi)\delta(E_a-\omega_\gamma-\omega)\,
\sum_T|V_0|^2\,
\frac{2}{p^2}\operatorname{Im}\!\left[\frac{-1}{\epsilon_L}\right],
$$

using $|D_{00}|^2\,2\operatorname{Im}\Pi_{00}=-2\operatorname{Im}D_{00}
=(2/p^2)\operatorname{Im}[-1/\epsilon_L]\ge0$,
$p=|{\bf k}_a-{\bf q}_\gamma|$. The variable $\omega$ here is the
medium's spectral variable; the exact bookkeeping (including the
$(2\pi)$ pairing and the photon-measure convention) is fixed in the
Assembly audit below.

**Vertex contraction.** $V_0=g\,\epsilon_{mna0}\,
q_\gamma^m e_T^n p^a$; transverse sum gives exactly

$$
\sum_T|V_0|^2 = g^2|{\bf q}_\gamma\times{\bf p}|^2,
$$

verified by explicit Levi–Civita summation at random kinematics
(ratio 1 to 12 digits) plus the collinear check ($=0$ exactly by
antisymmetry). Consequence: for ${\bf k}_a\to0$,
${\bf p}\parallel{\bf q}_\gamma$ and the rate vanishes — the
$v_a^2$ suppression is structural, so $v_a$ is load-bearing, and
$a\to\gamma+{}$plasmon with a rest axion is forbidden
(longitudinal–transverse collinear selection).

**Delta collapse.** With (A62), the $\omega$ integral sets
$\omega=\omega_p$, $\omega_\gamma=E_a-\omega_p$,
$\Theta(E_a-\omega_p)$. The model normalization satisfies the
Bethe sum rule $\int\omega\operatorname{Im}[-1/\epsilon_L]d\omega
=(\pi/2)\omega_p^2$ exactly, so no strength is missing.

**Angular integral.** With ${\bf p}={\bf k}_a-{\bf q}_\gamma$,
$|{\bf q}_\gamma\times{\bf p}|^2=q_\gamma^2k_a^2\sin^2\theta$ and

$$
I_\Omega = \int\!d\Omega\,
\frac{|{\bf q}_\gamma\times{\bf p}|^2}{p^2}
= k_a^2\,J(\kappa),
$$

with $J(\kappa)$ above (elementary partial fractions; checked
against numerical quadrature to 8 decimals for
$\kappa\in[0,5]$, $J(0)=8\pi/3$). Assembling the prefactors gives
the boxed formula. Ultra-relativistic limit:
$\kappa\to1/n$, $J(2/3)\approx7.577$.

**Assembly audit (prefactor and in-medium photon normalization).**
The cutting route above is calibrated against the vacuum width in
the same notation, so every phase-space factor is accounted for
exactly once.

- *Template, calibrated on vacuum.* The rate template is
  $$
  \Gamma = \frac{1}{2E_a}\,\frac{1}{2}\!\int\!
  \frac{d^4p}{(2\pi)^4}\,[2\pi\delta(q^2)]\,
  [-2\operatorname{Im}D_{00}(p)]\,g^2|{\bf q}\times{\bf p}|^2 .
  $$
  In vacuum, $[2\pi\delta(q^2)][2\pi\delta(p^2)]$ with the
  identical-line factor $1/2!$ gives
  $\tfrac12\int\frac{d^4k}{(2\pi)^4}(2\pi\delta)^2 =
  \int\! d\Phi_2\,(2\pi)^4\delta^4$, hence
  $\Gamma(a\to\gamma\gamma) =
  \frac{1}{2m}\frac{1}{2!}\frac{1}{8\pi}\cdot 2g^2(k_1\!\cdot\!k_2)^2
  = g^2m^3/64\pi$: the $1/2$ and the $1/(2\pi)^4$ are fixed by the
  exact vacuum width.
- *Free-photon cut (the stated convention).* The outgoing transverse
  photon is cut as the free massless density with the given-$n$
  dispersion inserted only in the on-shell relation:
  $$
  2\pi\delta(q^2)\big|_+ =
  \frac{\pi}{|{\bf q}|}\,\delta\!\left(E_a-p^0-\frac{|{\bf q}|}{n}\right).
  $$
  The model fixes the dispersion $|{\bf q}|=n\omega$ but not the
  in-medium spectral-density normalization; this is the one
  convention input of the assembly, and the numerical result is
  quoted in it.
- *Medium weight.*
  $-2\operatorname{Im}D_{00} = (2/|{\bf p}|^2)
  \operatorname{Im}[-1/\epsilon_L] =
  (2/|{\bf p}|^2)(\pi/2)\omega_p\,\delta(p^0-\omega_p)$ (KKL Sec. II
  identity chain re-verified against the fetched text; Jackson
  $D_{00}$ relation; Bethe sum rule checked).
- *Collapse and angular integral.* The $p^0$ integral collapses the
  medium delta; the radial $q$ integral gives
  $\int q\,dq\,\delta(E_a-\omega_p-q/n) = n\,q_\gamma$; the angular
  integral is $2k_a^2J(\kappa)$. Assembly:
  $$
  \Gamma = \frac{1}{2E_a}\frac{1}{32\pi^4}\,g^2\,
  \frac{\pi^2\omega_p}{2}\,4\pi n\,q_\gamma\cdot 2k_a^2J(\kappa)
  = \frac{g^2\omega_p\,n\,q_\gamma\,k_a^2\,J(\kappa)}{16\pi E_a},
  $$
  which is exactly the boxed form. The end-to-end assembly integral
  is checked numerically against the closed form in
  `TestPrefactorAssembly` (independent Simpson quadrature, relative
  agreement better than $10^{-9}$).
- *Residual convention dependence (one power of $n$).* Relative to
  the adopted convention (rate $\propto n^2$ here), the standard
  dielectric-field mode normalization ($1/\epsilon=1/n^2$) gives
  $1/n\approx0.67\times$, and the unit-residue model-dispersion
  density $2\pi\delta(\omega^2-{\bf q}^2/n^2)$ gives
  $n=1.5\times$, i.e.
  $[\,8.7\times10^{-42},\;1.3\times10^{-41},\;2.0\times10^{-41}\,]\,
  {\rm eV}$. The required benchmark value
  $1.31\times10^{-41}\,{\rm eV}$ is the adopted convention. The
  stated model does not by itself fix this power; it is recorded as
  the manifest limitation `prefactor_normalization`.

## Parameter roles

- $g_{a\gamma\gamma}$: overall $g^2$; GeV$\to$eV is $10^{-9}$
  ($g=10^{-19}\,{\rm eV}^{-1}$); a GeV-as-eV slip moves the rate
  by exactly $10^{18}$ (tested).
- $\omega_p$: pole position (threshold) and linear prefactor.
- $n$: photon momentum $q_\gamma=n\omega_\gamma$ (phase space), and the
  convention-dependent spectral-density power (see Assembly audit).
- $v_a$: enters only via $k_a^2$ — quadratic, essential.
- $m_a$: sets $E_a,k_a$ and hence which side of threshold.

## Independent checks (analysis record, `/tmp` scripts discarded)

1. Vertex normalization vs vacuum width: exact.
2. Tensor structure $|{\bf q}\times{\bf p}|^2$: numeric, 12 digits + collinear zero.
3. $J(\kappa)$ closed vs quadrature: 8 decimals, 9 points (now including
   $\kappa=1$ and the benchmark $\kappa\approx0.001984$).
4. Dimensions: $\lambda$-scaling gives $\lambda^3$ at fixed $g$ (rate is mass-dimension 1).
5. Sum rule: model saturates Bethe exactly.
6. Threshold continuity: rate $\to0$ smoothly ($\sim q_\gamma^4$); pole-mass case open only via kinetic energy.
7. Assembly audit: cut template calibrated on the vacuum width; end-to-end
   assembly integral vs closed form, relative agreement $<10^{-9}$
   (`TestPrefactorAssembly`).
8. $\kappa=1$ removable singularity: exact branch $2\pi$; both-sided
   approach finite and first-order close (`TestJClosedKappaOneLimit`).
9. Domain guards: negative or non-finite $\kappa$ returns NaN; out-of-domain
   rate inputs return `(NaN, false)`; the kinematic gate still returns
   exactly `(0, false)` (`TestDomainGuards`).

## Executed vs represented-only

- **Executed (exact rational Go):** parameter substitution, NR
  skeleton, threshold gate, ordering.
- **Evaluated numerically (float64, documented):** $J(\kappa)$,
  correction factor, final decimal, end-to-end assembly quadrature
  (`TestPrefactorAssembly`), domain-guard behavior.
- **Analytical hand derivation (represented as records + this
  file):** vertex algebra, cutting/spectral identity, delta-collapse,
  phase-space assembly and its normalization audit (including the
  convention statement and the residual one-power-of-$n$ dependence),
  angular integral.
- **Unsupported by the PhysMath rule registry:**
  $\delta$-distribution calculus, tensor contraction with a
  medium projector, phase-space integration, spectral-function
  identities. No PhysMath derivation is claimed; per the
  growth gate no new rule is added for one demo.

## Growth-gate note

No new PhysMath operation, family, schema, or shared abstraction
is introduced. Physical assumptions live in demo-local records
because no `phys-lib` axion-electrodynamics domain package
exists; creating one would change shared contracts and requires
authorization (proposed, not implemented). This demo is a
standalone Go module
(`github.com/PithomLabs/phys-demo/problem_8a/muse-deepseek`,
stdlib only; the directory keeps its name `muse+deepseek`, since
`+` is not a legal character in a Go module path). It is not listed
in the parent workspace file, so run demo commands with
`GOWORK=off` (workspace file unmodified).

## Limitations

- The $q$-independent pole response is given; real materials
  have continua, gaps, dispersion, local-field effects.
- In-medium photon normalization: the prefactor is audited and
  confirmed under the stated convention (README "Assembly audit"),
  but the model does not fix the spectral-density normalization of
  the in-medium photon; the residual dependence is one power of $n$
  (standard alternatives give $0.67\times$ to $1.5\times$ the quoted
  value). The quoted result is conditional on the stated convention
  (manifest limitation `prefactor_normalization`).
- The 15 eV zero is kinematic within the model.
- Transverse medium response absent by model (background
  literature supports the suppression).
- $n=1.5$ constant; recoil neglected; NR truncation remainder
  $O(v^2)\sim10^{-6}$ relative.
- No empirical claims; benchmark solutions appendix not read.
