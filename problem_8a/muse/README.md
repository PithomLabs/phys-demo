# Problem 8(a): Axion-to-Photon Conversion in a Dielectric Medium

meta/muse-spark-1.3-contributor

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

The 15 eV channel is **exactly closed**: at $T=0$ the medium can
only absorb energy, and $E_a<\omega_p$ leaves no phase space once
the response collapses to the pole. This is a kinematic result,
not a computed cancellation.

## Sources and fidelity

| ID | Source | Role | Fidelity |
|---|---|---|---|
| benchmark | arXiv:2610.16733v1, A.8.1, (A62) | problem, parameters, response model, output order | problem-definition |
| caputo | Caputo–Millar–Vitagliano, arXiv:2005.00078, Sec. II, (2)–(4), App. A | self-energy $\to$ rate template; narrow-width $\delta$ pattern; vertex-factor care | method |
| kkl | Knapen–Kozaczuk–Lin, arXiv:2101.08275, Sec. II, (5)–(6), (13) | density response $\leftrightarrow$ loss-function identity; loss function as inclusive excitation rate | response |
| jackson | Jackson, Classical Electrodynamics, 3rd ed., energy-loss chapter | $D_{00}=1/(q^2\epsilon_L)$; loss-function reading | identity |
| liang | Liang et al., arXiv:2608.20715 | transverse medium response 4–6 orders down | background |

Web search was non-functional in this environment; literature was
located through the arXiv API and fetched in full (HTML) before
use. No answer key was consulted.

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
$p=|{\bf k}_a-{\bf q}_\gamma|$.

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

## Parameter roles

- $g_{a\gamma\gamma}$: overall $g^2$; GeV$\to$eV is $10^{-9}$
  ($g=10^{-19}\,{\rm eV}^{-1}$); a GeV-as-eV slip moves the rate
  by exactly $10^{18}$ (tested).
- $\omega_p$: pole position (threshold) and linear prefactor.
- $n$: photon momentum $q_\gamma=n\omega_\gamma$ (phase space).
- $v_a$: enters only via $k_a^2$ — quadratic, essential.
- $m_a$: sets $E_a,k_a$ and hence which side of threshold.

## Independent checks (analysis record, `/tmp` scripts discarded)

1. Vertex normalization vs vacuum width: exact.
2. Tensor structure $|{\bf q}\times{\bf p}|^2$: numeric, 12 digits + collinear zero.
3. $J(\kappa)$ closed vs quadrature: 8 decimals, 8 points.
4. Dimensions: $\lambda$-scaling gives $\lambda^3$ at fixed $g$ (rate is mass-dimension 1).
5. Sum rule: model saturates Bethe exactly.
6. Threshold continuity: rate $\to0$ smoothly ($\sim q_\gamma^4$); pole-mass case open only via kinetic energy.

## Executed vs represented-only

- **Executed (exact rational Go):** parameter substitution, NR
  skeleton, threshold gate, ordering.
- **Evaluated numerically (float64, documented):** $J(\kappa)$,
  correction factor, final decimal.
- **Analytical hand derivation (represented as records + this
  file):** vertex algebra, cutting/spectral identity,
  delta-collapse, phase-space assembly, angular integral.
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
authorization (proposed, not implemented). The demo is its own Go
module (`github.com/PithomLabs/phys-demo/problem_8a/muse`) with a
project-local `go.work`, so demo commands run bare, e.g.
`go test ./...` from the project root; no shared workspace file is
involved.

## Limitations

- The $q$-independent pole response is given; real materials
  have continua, gaps, dispersion, local-field effects.
- The 15 eV zero is kinematic within the model.
- Transverse medium response absent by model (background
  literature supports the suppression).
- $n=1.5$ constant; recoil neglected; NR truncation remainder
  $O(v^2)\sim10^{-6}$ relative.
- No empirical claims; benchmark solutions appendix not read.
