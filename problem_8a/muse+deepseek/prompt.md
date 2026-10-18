# Phys Stress Test: Axion-to-Photon Conversion in a Dielectric Medium

## 1. Mission

You are new to the Phys infrastructure. Your task is to solve and formalize a concrete theoretical-physics problem from:

**TheoPhysEval: Assessing AI agents at autonomous theoretical physics problem solving**, arXiv:2610.16733v1, Appendix A, Problem 8(a), equation (A62).

Use this problem to stress-test whether the existing Phys infrastructure can faithfully preserve a scientific calculation from its original problem definition through assumptions, source equations, analytical derivation, intermediate results, numerical evaluation, and independent verification.

This is a combined physics and formalization task. You must actually derive the answer; do not merely create a manifest, write explanatory prose, or make an unsupported assertion that the problem has been solved.

However, the scientific derivation and its formalization are separate responsibilities. A mathematically correct result does not imply that Phys has represented and executed every step, and a passing test does not establish physical truth.

The intended approach is **analytical derivation first, targeted numerical evaluation second**. Do not default to numerical simulation.

Do not modify frozen predecessor systems, redesign the architecture, or add unsupported mathematical capabilities simply to make this particular example pass.

## 2. Mandatory workspace entry point

Include and follow this standard before doing any work:

**# AGENTS.md — workspace entry pointer**

For work involving `phys-artifact`, `phys-math`, `phys-lib`, `phys`, or cross-repository architecture, read `phys-system/AGENTS.md` first.

Package-local `AGENTS.md` files remain authoritative for local implementation details and constraints. The frozen `phys/AGENTS.md` governs the legacy system and is never modified by new-architecture work.

## 3. Reading order: understand Phys before coding

### First: workspace and scientific-derivation rules

Read the following in order:

1. `phys-system/AGENTS.md` — mandatory workspace entry point and cross-repository architectural guidance.
2. Relevant package-local `AGENTS.md` files in the repositories you intend to use. These instructions govern local implementation details.
3. [DERIVATION_SPEC_GUIDELINES.md](https://github.com/PithomLabs/phys-trace/blob/main/DERIVATION_SPEC_GUIDELINES.md) — the normative specification for formalizing scientific derivations.
4. `phys-trace/GUIDELINES.md` — the additive incident casebook, if present in the active checkout.
5. The root READMEs and relevant `doc.go` files for `phys-artifact`, `phys-math`, and `phys-lib`.
6. `phys-system/system.md` and other applicable architecture documents.

Do not assume a document exists at a particular path without inspecting the live workspace. Do not silently reconcile contradictory instructions; follow the specified authority hierarchy and report any material conflict.

### Second: understand each layer's responsibility

**PhysArtifact — canonical artifact substrate**

Inspect its actual implementation and tests, including `artifact.go`, `primitives.go`, `bigint.go`, `doc.go`, and `boundary_test.go` where present.

It owns artifact identity, references, digesting, deterministic encoding, canonical primitive framing, and grammar validation. It does not own physical meaning or scientific derivation semantics.

**PhysMath — theory-neutral mathematical representation**

Inspect `ast.go`, `identity.go`, `families.go`, `rules.go`, `replay.go`, and their tests, following the current repository's package instructions.

PhysMath represents mathematical structures and provides a bounded set of explicitly implemented operations. The current documented rules include expression canonicalization, simplification, substitution, rational normalization, differentiation, total differentiation, expression combination, branch selection, and expression application. Verify the live rule registry rather than assuming this list is current.

It is not an unlimited symbolic algebra system. Do not assume it supports arbitrary integration, unrestricted algebraic equivalence, or every expression transformation required by this problem.

**PhysLib — physical-framework semantics**

Inspect `phys-lib/core`, `phys-lib/gr`, `phys-lib/bohm`, and their package instructions and contracts. Relevant examples include:

- `phys-lib/bohm/minisuperspace` — a reference for representing derivation steps, assumptions, provenance, replay, independent goldens, and limitations.
- `phys-lib/gr/ltb` — a reference for formalizing an ongoing research question while preserving its open conjectures and unresolved proof obligations.

The target problem concerns axion-photon conversion in a dielectric medium. Do not place it under general relativity or Bohmian mechanics merely because those packages exist. Inspect the available packages and determine an appropriate physical-semantic location. If the existing architecture lacks a suitable domain package, apply the growth gate and report any material package or namespace decision that requires authorization before changing shared contracts.

**PhysTrace — derivation workflow**

Understand how its normative guidelines govern source fidelity, assumptions, derivation units, intermediate results, evidence, replay, and limitations.

### Third: inspect the existing reference pattern

Before writing new code, trace one complete existing derivation fixture from its source expressions through stored artifacts, derivation records, verification entry point, manifest, and tests.

Determine which existing representations can be reused directly. Do not invent APIs or types based on how you think the framework should work.

## 4. Scientific problem definition

Preserve the following problem faithfully.

### Axion-to-photon conversion in a dielectric medium

Compute the rate for axion conversion to photons via scattering off the electrons in a dielectric medium. Assume that there is no external electric field and the medium is a conventional solid state material (non-magnetic, zero temperature), that the axion only interacts via an \(aF\widetilde F\) coupling. Express the result in terms of the dielectric function \(\epsilon\) of the medium. Consider both relativistic and non-relativistic axions.

Consider the case when the dielectric function of the medium satisfies

\[
\operatorname{Im}\left[-\frac{1}{\epsilon_L(\omega,q)}\right]
=
\frac{\pi}{2}\omega_p\,\delta(\omega-\omega_p).
\tag{A62}
\]

Working in natural units \((\hbar=c=1)\), take the axion-photon coupling

\[
g_{a\gamma\gamma}=10^{-10}\,\mathrm{GeV}^{-1},
\]

frequency

\[
\omega_p=16.6\,\mathrm{eV},
\]

refractive index

\[
n=1.5,
\]

and axion velocity

\[
v_a=10^{-3}.
\]

Evaluate the rate \(\Gamma_{\mathrm{NR}}\) for:

1. \(m_a=25\,\mathrm{eV}\)
2. \(m_a=15\,\mathrm{eV}\)

Express the final numerical answer in scientific form in eV units, in the order \([(1),(2)]\).

**The problem is supplied, but its solution is not supplied in this task. Derive it independently. Do not copy an answer from another agent or treat the numerical output as known.**

## 5. Research and analytical derivation

Before formalizing the rate, establish what physical formula is actually justified by the stated assumptions.

1. Locate and inspect the relevant primary literature or authoritative source deriving the axion-to-photon conversion rate in a dielectric medium. Use the benchmark paper as the source of the problem definition, not as a substitute for the physical derivation.
2. Record source identifiers and exact equation or section locations. Distinguish source-derived equations from your own deductions.
3. Derive the rate expression in terms of the longitudinal dielectric response, using the stated interaction and material assumptions.
4. Derive or identify the expressions for both relativistic and non-relativistic axions, as required by the original problem.
5. Apply the supplied delta-function response analytically wherever justified. Show which integral or response factor it simplifies and how.
6. Substitute the parameters for both specified axion masses and obtain the final rates in eV units.
7. Independently check dimensions, conversion of \(g_{a\gamma\gamma}\) from \(\mathrm{GeV}^{-1}\) to the required energy units, normalization factors, and use of the refractive index and axion velocity.

Do not assume that every supplied parameter must survive into the final expression. Trace the role of each parameter from the physical formula and explain whether it contributes, cancels, or is irrelevant under the stated regime.

Do not silently choose among incompatible conventions for the interaction normalization, rate definition, dielectric response, or phase-space factors. Resolve conventions from sources where possible; otherwise record the ambiguity and its impact.

### Calculation strategy

Use pen-and-paper derivation and symbolic reasoning as the default.

The delta-function dielectric response suggests that part of the integration may be performed analytically, but establish this from the rate formula rather than assuming it. Use numerical computation for the final parameter substitution and targeted arithmetic verification, not as a replacement for deriving the rate.

If a remaining integral cannot be evaluated analytically, establish why numerical quadrature is necessary and isolate that computation. Do not build an extensive simulation or numerical framework for this problem.

If an equation or operation cannot be executed by existing PhysMath rules, keep the mathematical derivation explicit and independently checked; document the unsupported execution step. Do not introduce a fictitious symbolic-integration capability.

## 6. Formalize the research artifact in Phys

Once the physical derivation is established, create a small, auditable reference artifact using the live repository architecture and package contracts.

The artifact should preserve the following chain:

**Problem definition → source equations → assumptions → structured quantities and expressions → ordered derivation steps → intermediate results → final rate expression → parameter substitution → numerical results → independent checks → limitations.**

Where supported, represent the following as structured artifacts rather than opaque prose:

- The source-defined physical quantities and their units.
- The axion-photon coupling and dielectric-response relation.
- The assumptions on the medium, external field, temperature, interaction, units, and kinematic regime.
- The analytical relationship between the response function and the rate.
- The use of the delta function and any resulting simplification.
- The relativistic and non-relativistic branches as distinct derivation conditions or paths, as appropriate.
- The two parameter sets and the final rates.
- Derivation steps, input and output artifact references, applicable conditions, provenance, and verification obligations.

Preserve the identity of physical quantities and their units, operator order, argument order, and domain or regime restrictions. Do not represent quantities such as energy, mass, frequency, and coupling as interchangeable plain numbers without recording the necessary units and conventions.

Use existing PhysArtifact and PhysMath structures where they faithfully express the needed mathematics. Put physical assumptions and interpretation in an appropriate PhysLib package rather than embedding them into the theory-neutral mathematical layer.

The representation must reflect what was actually derived. Do not create a theorem artifact claiming more than the derivation establishes.

### Evidence and verification

Provide independent checks for the derivation and calculation where possible:

- Compare critical intermediate equations against the source and independently reconstructed expectations.
- Check units and dimensional consistency at every important transition.
- Test both mass cases and the requested output ordering.
- Add targeted negative tests for meaningful errors, such as the wrong dielectric-response sign, omitted factors, unit-conversion errors, a swapped mass assignment, or the wrong kinematic branch, where those mutations are applicable to the implemented representation.
- Ensure expected numerical values are computed independently of the code under test.
- Ensure stored derivation outputs and provenance actually correspond to the operations being claimed.

Distinguish structural validation, canonical decoding, executed mathematical operations, replayed derivations, independent analytical checks, numerical evaluation, and physical interpretation.

A reproducible numerical answer is useful evidence, but matching a reference number alone does not prove that the formula, assumptions, or physical interpretation are correct.

## 7. Apply the Phys growth gate

Before adding any mathematical operation, artifact type, schema, shared physical abstraction, namespace, or evaluator, establish:

1. The problem genuinely requires it.
2. Existing representations and operations cannot express or validate the requirement faithfully.
3. The capability belongs in the proposed architecture layer.
4. Its semantics, assumptions, valid inputs, and failure behavior can be specified and tested.
5. It is reusable or has a clear justification for being domain-specific.

Prefer the smallest coherent implementation with existing capabilities.

An exact integral or a domain-specific numerical evaluation that is necessary for this problem does not automatically justify adding a general-purpose integration API to PhysMath. If a real capability gap exists, document it, demonstrate its effect on the derivation, and propose the smallest appropriate extension. Do not modify shared contracts without authorization.

### Frozen-system rule

The legacy repositories `phys` and `phys-gr` are frozen predecessors/reference systems.

Do not change their code, instructions, schemas, or history. New work belongs in the current architecture and must follow its live package-level rules.

Do not commit, push, merge, tag, or publish without explicit authorization.

## 8. Acceptance criteria

The task succeeds when all of the following are true:

1. The original problem definition and its assumptions are faithfully preserved.
2. The relevant physical rate formula is sourced, and the analytical derivation is independently explained and checked.
3. Both the relativistic and non-relativistic cases are addressed as requested, with the specified non-relativistic numerical evaluations in the required order.
4. Numerical computation supports the derivation rather than substituting for it.
5. The structured artifacts and derivation evidence are consistent with the mathematical work actually performed.
6. Tests validate the implemented claims and fail on material malformed inputs.
7. Unsupported mathematical or physical operations are documented, not falsely presented as executed.
8. The result does not overstate what the tests prove.
9. No frozen-system or unauthorized architecture changes occur.
10. A reviewer can trace the final values back to sources, assumptions, equations, intermediate results, and independently checked arithmetic.

If the physical derivation cannot be completed or a critical source convention cannot be resolved, record the precise blocker and preserve the work as incomplete. Do not manufacture a final result merely to make the artifact appear complete.

## 9. Required completion report

Return a concise, evidence-based report containing:

1. **Scientific result:** the derived rate expression, its regime, and the two requested numerical results, with sources and actual assumptions.
2. **Derivation trace:** the important equations and intermediate steps, including how the dielectric response enters.
3. **Formalized artifacts:** exact paths and the role of each file or artifact.
4. **Executed versus represented-only:** which operations were actually executed and replayed, which checks were independent, and which steps remained analytical or unsupported.
5. **Verification:** exact commands executed and their actual outcomes, plus the independent checks and negative tests.
6. **Growth-gate decisions:** any requested capability that existing Phys could not express faithfully and why.
7. **Limitations:** unresolved source conventions, mathematical gaps, and anything not independently verified.
8. **Repository safety:** confirmation that frozen repositories were not modified and no unauthorized publication or Git history operation occurred.

## Final instruction

Treat this as a scientific derivation with a formalization exercise attached—not a coding exercise that happens to contain physics.

First establish the physics. Then encode what is justified. Execute only what the current infrastructure actually supports. Verify independently. Preserve the boundary between mathematical validity, physical interpretation, and empirical truth.

The goal is to stress-test whether Phys can serve as an auditable substrate for real theoretical-physics calculations, without turning unsupported assumptions or green tests into scientific conclusions.