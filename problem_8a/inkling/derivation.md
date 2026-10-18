# Derivation Trace — Axion-to-Photon Conversion (Problem 8a)

## Source claim (from TheoPhysEval arXiv:2610.16733v1, Appendix A, Eq. A62)

The rate $\Gamma_{\mathrm{NR}}$ for axion-to-photon conversion in a non-magnetic, zero-temperature dielectric with delta-function response

$$\mathrm{Im}\left[-\frac{1}{\epsilon_L(\omega,q)}\right] = \frac{\pi}{2}\omega_p \delta(\omega - \omega_p),$$

coupling $g_{a\gamma\gamma} a F_{\mu\nu}\tilde{F}^{\mu\nu}$, parameters:
$g = 10^{-10}\,\mathrm{GeV}^{-1} = 10^{-1}\,\mathrm{eV}^{-1}$, $\omega_p = 16.6\,\mathrm{eV}$, $n = 1.5$, $v_a = 10^{-3}$.

## Assumptions and conventions (recorded in manifest)

- `axion_interaction`: only $aF\tilde{F}$ coupling.
- `medium_nonmagnetic_zeroT`: $\mu = 1$, $T = 0$, $E_{\mathrm{ext}} = 0$.
- `delta_response_exact`: delta substitution exact.
- `refractive_index_constant`: $n = 1.5$.
- `natural_units`: $\hbar = c = 1$.

## Derivation steps (structured inventory)

1. **Interaction term**: $\mathcal{L}_{\mathrm{int}} = g_{a\gamma\gamma}\, a\, F_{\mu\nu}\tilde{F}^{\mu\nu}$.
2. **Rate integral**: general form involves phase-space integral over $\mathbf{q}$ and $\omega$ with the response function $\mathrm{Im}[-1/\epsilon_L]$.
3. **Delta substitution**: using the supplied delta response collapses the $\omega$ integral to $\omega = \omega_p$.
4. **Non-relativistic approximation**: for $v_a \ll 1$ and $m_a \gg \omega_p$ or comparable, the rate simplifies to an analytical form $\Gamma_{\mathrm{NR}}(m_a, v_a, n, g, \omega_p)$.
5. **Parameter substitution**: insert $g = 0.1\,\mathrm{eV}^{-1}$, $\omega_p = 16.6\,\mathrm{eV}$, $n = 1.5$, $v_a = 10^{-3}$, evaluate for $m_a = 25\,\mathrm{eV}$ and $15\,\mathrm{eV}$.

## Intermediate results (independent arithmetic)

The numerical values computed independently (not by the engine's expression evaluator, but by independent arithmetic in the protocol test) are:

- Case (1): $m_a = 25\,\mathrm{eV}$ → rate in scientific eV form.
- Case (2): $m_a = 15\,\mathrm{eV}$ → rate in scientific eV form.

These are stored in `fixture.NumResult25Ref` and `fixture.NumResult15Ref`.

## Evidence status

- `RateFormulaRef`: structured expression (representational, analytical form preserved).
- `DeltaSubstDerivRef`: replayed derivation linking delta response to rate expression.
- `NumResult25Ref`, `NumResult15Ref`: independent numerical artifacts.
- `PhysThm`: framework-level theorem linking assumptions to derivation.
- Full spatial integration: **represented only**, not executed numerically.
- Relativistic branch: **represented only**, structural representation in derivation inventory; numerical evaluation focused on non-relativistic case per problem request.

## Limitations (explicit)

- The full 3D spatial-averaging integral is not executed numerically (growth-gate: no integration primitive in current PhysMath registry).
- The delta-substitution is analytical, not a numerical integration over $\omega$.
- Numerical results are independent arithmetic corroboration, not engine-derived proofs.
- The exact normalization convention for the rate formula follows the problem statement; source-paper ambiguity (if any) is not independently resolved.
