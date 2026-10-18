# PROTOCOL.md — Adopting Phys Infrastructure for Physics Problem Solving and Knowledge Representation

**Status:** Definitive workflow contract for `phys-demo` (benchmark-solution packages) and
for external users/agents adopting the Phys infrastructure libraries.
**Scope date:** 2026-10-18. Baseline: Architecture v1.5 (`PM-ARCH-FOUNDATION-001`), Phase-1
Representation Specification v1.2 (`PM-REP-001`).
**Authority position:** subordinate to every authority in §2. This document governs the
benchmark-solution workflow and package contract; it overrides no representation
specification, no frozen repository rule, and no package-level normative instruction.

Normative language: MUST / MUST NOT / REQUIRED / SHALL are normative; SHOULD / MAY are
advisory (Spec v1.2 §0.2).

> **Non-authoritative sibling:** `phys-docs/PROTOCOL.md` is a stale descriptive layout
> sketch for a different purpose. It is NOT this document and grants no authority
> (decision record Part 1 item 7). Do not modify it in this task series.

---

## 1. Purpose, scope, and non-goals

**Purpose.** Make Phys infrastructure the default, reusable, enforceable way to
represent, reuse, validate, and report theoretical physics work — the equivalent of the
mathematical and conceptual work done **with pen and paper** (decision record, Part 2
items 1–2).

**In scope.**

1. **Benchmark problems** (use case A): solving and documenting a specified physics
   problem (TheoPhysEval or equivalent) as a Type-B package (ADR §2.2).
2. **Paper-scale knowledge representation** (use case B): formalizing theorems,
   derivations, and physical frameworks as Type-A packages (Spec §1.4).
3. The interface between the two: reuse, capability discovery, evidence recording,
   conformance, reproducibility.

**Explicit non-goals** (decision record Part 1 item 10, Part 2 items 3–4, 8; ADR §5):

- **No numerical infrastructure.** Numerical computation, numerical integration, Monte
  Carlo, and solvers are out of scope. Do not expand PhysMath or PhysLib into numerical
  solvers to accommodate TheoPhysEval. `phys-calc` is a possible future layer, not an
  authorized project (§9.3).
- **No dimensional-semantics system.** §9.1 is metadata-only.
- **No autonomous discovery.** Scientific analysis, interpretation, intuition-guided
  research, and Jev or other probabilistic decision services stay outside this protocol.
- **No new executable rules, PhysArtifact families, PhysMath types, or registry
  changes** via document-writing tasks.
- **No frozen-repository changes** (`phys`, `phys-gr`, frozen schemas) without separate
  authorization (Spec §0.1).

**Deliverable standard:** a reproducible representation workflow for theoretical physics —
not universal numerical problem-solving capability and not an automated
scientific-discovery system.

---

## 2. Authority and document precedence

Read in this order. Where instructions conflict, follow the higher authority and record
the conflict; never silently discard either (`phys-system/AGENTS.md` §21).

| # | Authority | Governs | Canonical location |
|---|---|---|---|
| 1 | Frozen Phys: `specs_v2_3.md` + `AGENTS.md` | Existing Phys semantics only (Spec §0.1: frozen authority wins on Phys behavior; never modify Phys) | `phys/` (frozen) |
| 2 | **Architecture v1.5** `PHYSMATH_PHYS_EXTERNAL_ARCHITECTURE_v1.5.md` | Cross-project architecture: dependency topology (§2), identity (§4–§5), numerical boundary (§25), human admission of `phys-*` (§27.5), growth gate (§30), No-Invention Matrix (§32.2) | `phys-docs/` (byte-identical copies in `physmath/plans/`, `TOE-v4/plan2/`) |
| 3 | **Representation Spec v1.2** `PHYSMATH_PHYS_PHYSSTAR_REPRESENTATION_SPEC_v1.2.md` | API/representation semantics: artifact families, §1.4 formalization contract, §12 Derivation, §17 operation semantics, §20 replay, §21 conformance, §24.4 stop rule | `phys-docs/` (byte-identical copies elsewhere) |
| 4 | **ADR: `adr_package_types.md`** (Task 3 planning area) | The two package types and mixed-purpose composition — nothing else | `TOE-v4/plan2/adr_package_types.md` |
| 5 | **This document** | Benchmark-solution workflow, Type-B schema, adoption/units/evidence rules for `phys-demo` | `phys-demo/PROTOCOL.md` |
| 6 | Repository `AGENTS.md` files | Local process: `phys-system/AGENTS.md` (collective, §21 decision hierarchy), `phys-math/AGENTS.md` (test sharding, registry gates), `phys-trace/AGENTS.md` (engineering + evidence principles), `phys-lib/AGENTS.md`, fixture-local `AGENTS.md` (package contract) | each repo |
| 7 | `DERIVATION_SPEC_GUIDELINES.md` (DSG) | Derivation specifications for paper-scale work (§3 pin source, §5 quantity map, §8 derivation units, §9 capability audit, §10 independent expected results, §12 package contract, §14 readiness gate) | canonical copy: `phys-trace/` (phys-docs copy byte-identical) |
| 8 | `GUIDELINES.md` | Additive lessons/casebook — never stands alone against a normative document | canonical copy: **`phys-trace/`** (its `AGENTS.md` pointer is correct; the `phys-docs/` copy contains a stale `SKILLS.md` pointer — do not replicate) |
| 9 | READMEs, capability maps (`capability_map_*.md`), plans, reviews | Descriptive/history only | — |

**READMEs and capability maps are never executable authority.** Whether an operation is
executable is established solely by the live rule registry
(`phys-math/families.go` + `rules.go`) and its tests (ADR §5; decision record Part 1
item 9). A capability-map entry whose operation is absent or fail-closed (e.g.
`ExecuteOperator("covariant_derivative")` always returns `ErrOperatorNotExecutable`;
`integral/1` has no consuming rule) MUST be read as proposal-status, not capability.

**Document-control note (nonblocking):** Spec v1.2's header says "Version: 1.2" while its
body contains "v1.4 wins" references. Per decision item 6 this is recorded as a
document-control anomaly; cite **section numbers**, not version claims. Escalate only if
a material conflict actually affects a requirement (none found for this protocol).

**`phys-demo` is not a Git repository.** Record it; do not initialize it here
(decision item 12). Until resolved, exact revisions cannot be pinned for `phys-demo`
itself; pin the infrastructure repos (§13.2).

---

## 3. The two package types

Authority: ADR §2–§3 (`TOE-v4/plan2/adr_package_types.md`).

| | **Type A — Formalization package** | **Type B — Benchmark-solution package** |
|---|---|---|
| Purpose | Formally represent a theorem, hypothesis, derivation, or paper content | Solve and document a specified benchmark problem |
| Trigger | Declared purpose *or* any artifact claiming formalization status (Spec §1.4/§1.5) | Declared purpose *and* no formalization-status claim for any contained artifact |
| Layout | `AGENTS.md` + `manifest.json` + `protocol/derivation_test.go` (Spec §1.4, GR example) | `manifest.json` + `protocol/derivation_test.go` + `README.md` (+ optional package `AGENTS.md`); see §11.1 |
| Manifest schema | **Exact GR schema and loading convention** (Spec §1.4 verbatim) | `phys-demo/benchmark-solution-manifest/1` (§11.2) |
| Governing doc | Spec v1.2 §1.4 + DSG | This document + ADR |

**Composition and anti-evasion (ADR §3 — normative):**

1. **Category labels never control.** Declaring `package_type:
   benchmark-solution-manifest/1` exempts nothing from Spec §1.4. The contract follows
   declared purpose and contained-artifact status, not the label.
2. **Formalization claims pull in the full Type-A contract.** If any contained artifact
   claims to be a formally represented theorem/hypothesis/paper derivation (Spec §1.4/§1.5),
   either **(A)** relocate that content into its own conformant Type-A package and
   reference it (`formalization_refs`), or **(B)** adopt the full Type-A contract in this
   package and report any unexpressible benchmark metadata as a `PROTOCOL GAP`
   (Spec §24.4) — never as a schema extension, never in an unregistered sidecar.
3. **Ambiguity resolves to Type A**, the stricter contract.
4. **Evasion fails conformance:** a Type-B-declared package asserting formalization
   claims without (A) or (B) is NONCONFORMANT; its conformance test MUST fail.

This protocol never defines a competing protocol for Type A: for formalization, Spec §1.4
plus DSG remain the sole contract.

---

## 4. Component responsibilities and dependency boundaries

Dependency direction (Arch §2.2): `phys-lib/phys-* → phys-math → phys-artifact`.
No cycles. Frozen `phys` does not import any of them (Arch §2.3).

| Library | Owns (decision record Part 2 item 2) | Must NOT own |
|---|---|---|
| **PhysArtifact** (`github.com/PithomLabs/phys-artifact`) | `ArtifactRef` identity (namespace/schema/type/digest), canonical primitive framing, BigInt, reference validation. Enforced stdlib-only (`boundary_test.go`) | Math semantics, physics ontology, trust (Arch §2.4) |
| **PhysMath** (`github.com/PithomLabs/phys-math`) | Theory-neutral mathematical representation — expression IR, functions/parameters/constants, relations, derivations, conditions, obstructions — **and only the symbolic operations actually implemented** by the closed 9-rule registry | Physical meaning, units, numerics, analysis (Spec §2, §17) |
| **PhysLib / `phys-\*`** (`github.com/PithomLabs/phys-lib`) | Physical-framework semantics: assumptions, framework-specific equations, interpretations, theorem/derivation context (`gr`, `bohm`, future frameworks). New `phys-*` admission is a human governance decision (Arch §27.5) | Theory-neutral math, numerical execution |
| **PhysTrace** (`github.com/PithomLabs/phys-trace`) | Derivation protocols, replay validation, evidence of what has been mechanically verified; DSG/GUIDELINES process authority; MCP inspection tooling (currently scoped to phys-lib/bohm — `phys-demo` packages cannot use it, a recorded limitation) | Physics content itself |

**Deferred boundaries preserved** (Spec §0.4/§23.2, Arch §25): Numerical execution,
ComputationalContract execution, trust/corpus promotion, scientific adjudication,
paper-management, autoformalization — separate systems; MUST NOT be simulated by this
protocol. Numerical MUST NOT import PhysMath or `phys-*` (Arch §25.13).

---

## 5. Capability discovery (mandatory, before writing solver code)

**REQUIRED sequence** (DSG §9; decision record Part 1 item 2):

1. **Read the live rule registry.** List every executable rule from `phys-math/families.go`
   / `rules.go` — currently exactly nine, each with direct tests:
   `EXPR-CANON-001`, `EXPR-SIMPLIFY-001`, `EXPR-SUBSTITUTE-001`, `RAT-NORMALIZE-001`,
   `DIFFERENTIATE-001`, `TOTAL-DERIVATIVE-001`, `EXPR-COMBINE-001`, `SELECT-BRANCH-001`,
   `EXPR-APPLY-001`. This list, not documentation, establishes executability.
2. **Read the public API of the three libraries** (`go doc` or source): families in
   `phys-math/families.go` (33 registered families incl. `expression/1`, `relation/1`,
   `derivation/1`, `theorem/1`, `integral/1`), `artifact.ArtifactRef`/`ComputeRef`,
   `phys-lib` framework constructors (`gr.NewFramework`, `PhysicalAssumption`,
   `PhysicalTheorem`), claim statuses (`phys-lib/gr/ltb/claim.go`).
3. **Inspect tests** of any API you intend to use — usage patterns and rejection cases.
4. **Inspect conformant fixtures**: `phys-lib/gr/equivariance` (the conformant reference,
   §14.1) and the registered-gap fixtures (§14.2).
5. **Classify each needed capability** with one of the five states and record it (e.g. in
   README or manifest `capability_audit` prose):

   | State | Meaning | Example |
   |---|---|---|
   | `SUPPORTED` | Implemented, tested, fits the need | exact rationals (`NewRational`), differentiation |
   | `PARTIALLY SUPPORTED` | Exists with material limits | `EXPR-SIMPLIFY-001` (closed rewrite list only) |
   | `REPRESENTATIONAL ONLY` | Type/family exists, **no rule executes it** | `integral/1` (no consuming rule) |
   | `DEFERRED` | Architecturally deferred | Numerical execution (Spec §0.4) |
   | `UNAVAILABLE` | No representation exists | δ-distributions, limits, units |

   (Vocabulary: `phys-trace/AGENTS.md` §8 capability states — do not invent synonyms.)

**Only after this audit may implementation start.** Implementing a local alternative to a
`SUPPORTED` capability without first recording the capability audit is a protocol
violation (§8.4). A represented operation MUST NEVER be described as executable when the
rule registry does not implement it (decision record Part 2 item 5).

---

## 6. Canonical workflow — Type B (benchmark solutions)

Twelve stages; each names its authority, required artifact, and test.

1. **Select and identify the problem.** Record benchmark id, version, and exact location
   (e.g. `arXiv:2610.16733v1`, App. A Prob. 8(a) eq. (A62)). → `manifest.problem`.
2. **Pin the source specification.** DSG §3: source id/version/digest, locations, fidelity
   class (the five evidence classes), and explicitly which parts were *not* consulted
   (e.g. "solutions appendix not consulted"). → `manifest.problem.fidelity`.
3. **Extract assumptions, conditions, conventions, symbols, requested outputs.** DSG §5/§7:
   quantity map; every parameter declared with unit state (§9.1); regime recorded.
   → `manifest.assumptions`, `manifest.parameters`, `manifest.regime`,
   `manifest.requested_outputs`.
4. **Capability audit.** Section 5 sequence, completed and recorded before coding.
   → README `Capability audit` section or `manifest.derivations[].support`.
5. **Reuse decision per operation.** For each operation: reuse a `SUPPORTED` capability,
   or record a deviation per §8.3. Imports alone never count (§8.1).
6. **Represent the theoretical content.** Express equations and derivation dependencies
   with PhysMath/PhysLib structures where supported: `Relation`/`Equation`/`Derivation`
   artifacts minted through `physmath.Store` (pattern: `phys-lib/gr/ltb`); framework
   assumptions via `phys-*` **only inside their own framework** (never force GR semantics
   onto non-GR physics — §8.4). Unsupported steps get `support: UNSUPPORTED` or
   `REPRESENTATIONAL ONLY` with reasons — never fabricated rule IDs
   (inkling's `delta_substitution_analytic` was impossible in the closed registry;
   an invented rule ID MUST cause rejection).
7. **Solver-local computation, if required.** Section 9.2. Keep it explicitly outside the
   representation libraries; record algorithm, conventions, inputs, outputs, independent
   checks; attribute results to local code only.
8. **Record evidence at the correct epistemic level.** Section 10 status vocabulary — no
   level may be asserted upward because "tests pass".
9. **Write the tests.** `manifest_test.go` (strict schema validation, §11.3) +
   `protocol/derivation_test.go` (§12): every declared entry point exercised, reuse
   claims tested, deliberate-corruption negatives, unit-error tests where §9.1 applies.
10. **Conformance computation.** The package test computes CONFORMANT/NONCONFORMANT from
    schema conformance + unresolved mandatory gaps (§11.3). Status is computed, never
    self-declared in the manifest.
11. **Deviation/gap registration.** §8.3, §11.4 — all deviations and gaps explicit.
12. **Reproducibility.** §13: standard `go test ./...` from package root in the
    project-local workspace context; `scripts/test-all.sh` at `phys-demo` root.

---

## 7. Canonical workflow — Type A (paper-scale representations)

**This protocol does not define it; it exists and is binding.** Follow:

1. `phys-system/AGENTS.md` (collective decision hierarchy) and the target repo
   `AGENTS.md`.
2. **`phys-trace/DERIVATION_SPEC_GUIDELINES.md`** — the normative derivation
   specification standard (read §1–§14 before drafting; §14 readiness gate).
3. **Spec v1.2 §1.4** — exact package layout and exact GR manifest schema; any
   unexpressible need is a `PROTOCOL GAP` (§24.4), never a schema extension.
4. Established fixture precedent: `phys-lib/gr/equivariance` as conformant reference;
   `gr/ltb`, `gr/maxwellbox`, `bohm/minisuperspace` as registered-gap exemplars
   (§14.2), each governed by its fixture-local `AGENTS.md`.

Mixed-purpose content: ADR §3 (reference, don't merge; ambiguity → Type A).

---

## 8. Infrastructure adoption requirements

### 8.1 `infrastructure_use` record (required in Type-B manifests)

Every Type-B manifest MUST include an `infrastructure_use` array. Each entry MUST
contain:

| Field | Content |
|---|---|
| `concern` | The operation or representation concern (e.g. "exact-rational parameter encoding") |
| `library` | `phys-artifact` \| `phys-math` \| `phys-lib` (exact module name) |
| `api_symbol` | Exact public symbol, e.g. `physmath.NewRational`, `artifact.ComputeRef`, `physmath.Replay` |
| `code_path` | File:symbol where the API is called in this package |
| `test` | Name of the test that **exercises** the claimed code path |
| `result` | What the reuse produced (artifact ref, replay outcome, validation result) |
| `limitations` | Honest limits of this reuse; `""` only if genuinely none |

Rules (decision record Part 1 item 4):

- **Imports alone never constitute adoption.** An import without an exercised code path
  in `test` is unverified and MUST NOT be listed.
- `infrastructure_use` MAY be empty **only if** `deviations` is non-empty (§8.3) and
  covers why nothing was reusable. Empty-and-undeclared fails conformance (§11.3).
- **No import-only compliance:** a package whose tests never execute the imported API
  (inkling's pattern) is NONCONFORMANT under §12.

### 8.2 Acceptance tests detect both failure modes

1. *Bypass despite supported capability* → test: for each `SUPPORTED` capability claimed
   in the capability audit as reused, a test must exercise it; audit entries marked
   reused with no corresponding `infrastructure_use` entry fail validation.
2. *Import without meaningful use* → test: every `infrastructure_use.code_path` symbol
   referenced by `test` must actually run (see §12 entry-point rule); a test that never
   touches the API cannot satisfy `infrastructure_use.test`.

### 8.3 Deviation record (when no suitable API exists)

`deviations[]` entries MUST contain: `capability_sought` (what was needed),
`api_examined` (exact API/contract inspected and why it was examined),
`reason` (why it cannot faithfully represent the problem), `local_alternative` (the
smallest justified local replacement). **An honest, justified deviation is conformant;
silent bypass is not.** Growing a local workaround for a `SUPPORTED` capability without
this record is a violation.

### 8.4 Prohibitions

- MUST NOT misuse a GR (or any) framework for non-GR physics to inflate usage
  (Spec §14 framework semantics; `gr.NewFramework` hardcodes `phys-gr` namespace —
  the honest option for new domains is demo-local records + a deviation entry, never a
  namespace violation).
- MUST NOT invent APIs, rule IDs, or types "as the framework should work".
- MUST NOT present `phys-lib`'s internal conformance as importable: it is
  repository-internal (§11.3).

---

## 9. Units and numerical boundaries

### 9.1 Interim unit metadata policy (metadata only)

Every parameter entry in `manifest.parameters[]` MUST declare exactly one unit state:

| `unit_state` | Required extra field |
|---|---|
| `explicit` | `unit` (string, e.g. `"eV"`) |
| `dimensionless` | — |
| `natural-units` | `convention` (string, e.g. `"hbar=c=1; all quantities in eV"`) |
| `unresolved` | `note` (why unresolved) |

- **This is metadata, not dimensional analysis.** It does NOT establish a
  dimensional-semantics system in Phys; that decision remains separate and open
  (decision record Part 1 item 5; Arch §32.2 open row).
- **Unit tests are solver-local or package-local**: a package MAY (and for
  unit-bearing numerics SHOULD) include unit-conversion and *deliberate unit-error* tests
  — its own workflow tests for unit errors and records the method (pattern:
  `muse+deepseek/protocol/derivation_test.go` `TestUnitSlipSensitivity`). Such a test is
  **not** a PhysMath or PhysLib capability and MUST NOT be attributed to them.
- MUST NOT add physical dimension types, quantity kinds, or a shared conversion engine
  to PhysMath/PhysLib.

### 9.2 Solver-local numerical computation

Where a benchmark requires numerical calculation, integration, Monte Carlo, ODEs, or any
computation beyond the nine rules (decision record Part 2 items 3, 7; ADR §5):

1. **Keep it outside the Phys representation libraries** — plain solver-local Go
   (or a separate module), never a PhysMath/PhysLib extension.
2. **Record it in the package** as `solver_local_computations[]`: `id`, `algorithm`,
   `conventions` (e.g. quadrature rule, precision), `inputs`, `outputs`,
   `independent_checks` (a check the library cannot provide, e.g. Simpson-vs-closed-form),
   and attribution (`"executed_by": "solver-local"`).
3. **State which parts are represented by Phys and which were evaluated outside it**
   (README boundary section + `derivations[].support`).
4. **Never attribute local numbers to Phys**: a numerical result MUST NOT be described as
   computed or verified by PhysMath/PhysLib. The only PhysMath-verified claims are
   replayed derivation results under the closed rule registry.
5. **Never require Phys infrastructure to execute the numerical problem.** If a required
   operation is out of scope for representation, record the boundary (§11.4) — do not
   treat it as a missing feature this protocol must supply.

### 9.3 `phys-calc`

A possible future layer built on top of Phys (Arch §25 shapes this direction).
**Not authorized, specified, or created by this protocol.** Any mention of future
numerical capability MUST be labeled "future possible layer, not current capability".

---

## 10. Representation and epistemic status

Use existing vocabularies only — no competing status system.

**Evidence levels** (`phys-trace/AGENTS.md` §8) — a claim MUST state its level:

1. `Represented` — faithful structural representation exists.
2. `Structurally validated` — schema/types/references/shapes passed.
3. `Executed` — a **supported** operation actually computed the result.
4. `Replayed/verified` — recorded steps reproduce the result from bound premises under
   the operation semantics.
5. `Physically interpreted` — connected to a physical model under explicit assumptions.
6. `Empirically supported` — observations/experiments support the claim.

Rules:

- **Success at one level never implies a stronger level.** A passing test is not
  scientific truth; replay proves only what the operations and premises justify; a
  structural certificate is not a proof; a derivation is not empirical confirmation
  (phys-trace AGENTS §8; `phys-lib/README.md`).
- **Claim statuses** (existing vocabularies): `literature-claim-requiring-verification`,
  `candidate-construction-incomplete`, `open-conjecture`, `deferred`, `proved`
  (`phys-lib/gr/ltb/claim.go`); formalization outcomes `PROJECTABLE`/`UNREPRESENTABLE`/
  …(Arch §26.3) for research-layer context.
- **Solver-local numerical results** (§9.2) record their level from these same six
  levels where applicable, and MUST carry the `executed_by: "solver-local"` attribution —
  local means honest, never weaker, and never implies Phys executed the computation.
  (`phys-demo/PROTOCOL.md` creates no status labels outside this vocabulary and the
  existing claim/formalization vocabularies above.)
- **Never** a bare `verified: true` flag or caller assertion in place of verifiable
  evidence (phys-trace AGENTS §8).
- Distinguish explicitly: represented mathematics / replayed derivations / verified
  claims / unverified source mappings (bohm's `source_mapping_structural_only`) /
  open conjectures (ltb `ZEH-Q-03`) / local numerical results / empirical claims.

---

## 11. Package structure, manifest validation, and protocol gaps

### 11.1 Type-B package layout (required files)

```text
<phys-demo>/problem_<id>/<solution>/
    manifest.json                      # schema phys-demo/benchmark-solution-manifest/1
    manifest_test.go                   # strict schema validation (§11.2/§11.3)
    protocol/derivation_test.go        # executable protocol test (§12)
    README.md                          # human orientation: question, result, capability
                                       # audit, executed-vs-represented boundary
```

Optional: package `AGENTS.md` (process authority for the package). Layout is
package-level convention — it creates no artifact family and is not Spec §1.4's
formalization layout.

### 11.2 Type-B manifest schema — `phys-demo/benchmark-solution-manifest/1`

**Required top-level keys** (exactly these allowed; unknown keys are rejected):

```text
schema_id            "phys-demo/benchmark-solution-manifest/1"   (exact)
package_type         "benchmark-solution"                        (exact; = ADR §2.2)
package_id           stable kebab id
problem              {benchmark, location, fidelity, requested_outputs[]}
parameters[]         {symbol, value, unit_state, unit|convention|note}
regime               string
conventions          object (interaction/dispersion/measure conventions, etc.)
assumptions[]        {id, statement}
capability_audit     {supported[], partial[], representational_only[], deferred[], unavailable[]}
derivations[]        {id, support: SUPPORTED|PARTIALLY SUPPORTED|REPRESENTATIONAL ONLY|UNSUPPORTED,
                      executed_rules[] (rule IDs actually replayed), reason}
infrastructure_use[] §8.1 records (may be empty iff deviations non-empty)
deviations[]         §8.3 records: {capability_sought, api_examined, reason, local_alternative}
solver_local_computations[]   §9.2 records (algorithm, conventions, inputs, outputs,
                      independent_checks, executed_by: "solver-local")
results[]            {id, value, evidence_level (§10 vocabulary), evidence}
limitations[]        strings
open_obligations[]   strings
declared_gaps[]      {gap_id, requirement, reason, severity: mandatory|advisory}
formalization_refs[] optional: {path, artifact_ref} to conformant Type-A packages
entry_points         {primary, secondary} — each MUST be exercised (§12)
```

Validation rules: strict JSON parsing (unknown keys rejected, bidirectional —
manifest keys vs schema keys must match exactly); enum fields validated against §10/§5
vocabularies; `schema_id`/`package_type` exact; if `unit_state == "explicit"` then
`unit` REQUIRED (analogously §9.1); `infrastructure_use[].test` MUST name a test that
exists in the package; `deviations[]` entries complete per §8.3;
`solver_local_computations[].executed_by` MUST be `"solver-local"`.

**Category firewall (ADR §3):** if any `results[]`/`derivations[]` claim formalization
status, the package MUST either carry non-empty `formalization_refs` pointing to
conformant Type-A packages (Option A) or be declared under the Type-A contract itself
(Option B); otherwise validation MUST fail.

### 11.3 Validation mechanism and its recorded limitation

- **Validation is performed by each package's own `manifest_test.go`**, using strict
  parsing and the schema above. This is a package-level contract (ADR §2.2).
- **Recorded limitation (decision item 4; ADR §6):** `phys-lib/internal/manifestschema`
  is **repository-internal** and MUST NOT be imported from `phys-demo`; it is NOT
  exported, and its implementation MUST NOT be duplicated across repositories. A shared
  Type-B checker would be a separate, unauthorized implementation decision — until one
  exists, per-package validation is the whole enforcement story, and this document does
  not imply a central registry.
- **Conformance is computed, not declared.** The package test computes:
  `CONFORMANT` = schema valid ∧ every entry point exercised ∧ every claim has evidence at
  a stated level ∧ **zero unresolved `severity: mandatory` gaps** ∧ no category-firewall
  violation. Any unresolved mandatory gap ⇒ `NONCONFORMANT`, logged explicitly
  (mirrors phys-lib's standing rule: a registered gap is not conformance).

### 11.4 Protocol-gap registration

- `declared_gaps[]` is the Type-B gap registry (Type A uses Spec §1.4's `PROTOCOL GAP`
  procedure — unchanged).
- Registered ≠ resolved: mandatory unresolved gaps keep the package explicitly
  nonconformant (§11.3), machine-identifiable through the manifest test.
- Gaps record missing capabilities honestly (e.g. "no δ-distribution rule exists —
  capability UNAVAILABLE, derivation step representational-only"). This protocol
  **documents** missing infrastructure; it never pretends to supply it.

---

## 12. Testing and independent verification

Normative safeguards (decision record Part 1 item 9; Task 4 §5):

1. **Every declared entry point is exercised by a meaningful test** that observes actual
   behavior: `manifest.entry_points.primary/secondary` MUST resolve to functions called
   by `protocol/derivation_test.go`. A test suite that never calls them is vacuous.
   *(Anti-example: inkling declared `problem_8a.Build` while `Build()` panics — §14.3.)*
2. **Empty or vacuous protocol tests can never count as successful verification.** A
   `derivation_test.go` containing no assertion over the fixture's content (an empty
   `TestEvidenceNonEmpty`-style stub) MUST be treated as no verification; the conformance
   computation treats "protocol test exists but exercises nothing" as failure.
3. **Imports alone never count as adoption** (§8.1); applicable reuse MUST be traceable
   to `api_symbol` + `code_path` + a test that executes it.
4. **Unsupported/unavailable operations are recorded honestly** (`support` states, gaps) —
   claiming executability for a non-registry operation fails review; an invented rule ID
   in any `Derivation` step MUST fail `physmath.Replay` (closed registry,
   `ValidPhase1RulePair`).
5. **Manifest conformance validated against the package type's governing schema** (§11).
6. **Every allowed deviation is explicit; unresolved mandatory gaps prevent conformant
   status** (§11.3/§11.4).
7. **Deliberately corrupted metadata MUST fail:** negative tests mutate the manifest
   (wrong `schema_id`, missing `unit`, broken `infrastructure_use.test` reference,
   formalization claim without `formalization_refs`) and assert rejection — mirroring
   phys-lib's mutation-drill pattern (`TestMutationNegatives`, bohm's named-corruption
   tests).
8. **Replay success ≠ scientific truth** (§10); tests state what they establish and what
   they do not (README boundary section; test comments).
9. **Solver-local numerics never attributed to PhysMath/PhysLib** (§9.2).
10. **Frozen repositories, canonical encodings, artifact-family registries, and the
    executable rule registry are unchanged** by any protocol work (decision item 9).

**What tests establish (summary):** schema validity, exercised entry points, byte-exact
replay of *represented* steps, honest status bookkeeping, rejection of corruptions.
**What they never establish:** physical truth, empirical support, correctness of
solver-local numbers, correctness of source mappings (bohm: structural only).

---

## 13. Reproducibility and external consumption

### 13.1 Module/workspace setup (as verified; do not modify `go.mod`/`go.work`)

- Every `phys-demo` solution directory is its own Go module under
  `github.com/PithomLabs/phys-demo/...`, with a project-local `go.work` (pre-existing —
  do not relocate or attribute to documentation tasks).
- Type-B **self-contained** solutions (`muse`, `muse+deepseek`): `go.work` = `use .`
  only; stdlib-only `go.mod`. Verified command:
  `cd <solution> && GOWORK=off go test -count=1 ./...`
- **Infrastructure-consuming solutions** (`inkling`): `go.work` adds `use` entries for
  `../../../phys-artifact`, `../../../phys-math`, `../../../phys-lib` so the sibling
  multi-repo checkout resolves. Verified command: `go test ./...` from the solution root.
- There is **no root `go.work`** anywhere above; do not assume one exists. Do not run
  `go mod tidy`; do not add `require`/`replace` lines (that is an implementation
  decision — phys-math currently lacks a `require` on phys-artifact; recorded as
  integration friction, not fixed here).
- Root orchestration: `phys-demo/scripts/test-all.sh` runs `go test ./...` per project,
  continues after failures, exits nonzero on any failure.
- **Limitation:** `phys-demo` is not a git repository — no revision pinning for demo
  packages (decision item 12; §2). Infrastructure repos: record `git rev-parse HEAD`
  before citing results.

### 13.2 Reporting conventions

A completion report MUST state: exact commands run and their outcomes; baseline
revisions of every repository touched by a claim; executed vs represented vs local
boundaries; open gaps. **A documentation-only task MUST NOT claim executable test
suites were validated beyond what was actually run.**

---

## 14. Worked examples and anti-examples

### 14.1 Conformant formalization reference — `phys-lib/gr/equivariance`

The only manifest conforming to the GR baseline today (machine-verified by
`phys-lib/internal/manifestschema`: `gr/equivariance/manifest.json CONFORMANT`). Copy its
pattern for Type A: canonical nested layout, exact 6-step pinned rule order
(APPLY/APPLY/CANON/APPLY/APPLY/CANON), byte-linkage between proposition sides and
derivation steps, hash-pinned `expected_result`, `TestManifestProseIsNotMath` (prose can
never substitute for the hash), fresh-verifier anti-goldens. **Reused here as:** the
evidence discipline every package should emulate — claim-to-conclusion trace, not
assertion.

### 14.2 Registered-gap exemplars — `gr/ltb`, `gr/maxwellbox`, `bohm/minisuperspace`

All NONCONFORMANT with explicit registered gaps (`phys-lib/PROTOCOL_GAPS.md` +
`internal/manifestschema`, machine-checked): LTB (claims/research metadata relocated to
sidecars; main theorem `open-conjecture`; executed: `dM/dr` + EdS theta only), Maxwellbox
(PARTIAL execution by design — `TestOnlyCanonStepsExecuted` proves only CANON ran;
K1 not end-to-end), bohm minisuperspace (heaviest rule use — 14 replayed steps — yet a
toy model; `source_mapping` structural, explicitly *not* source-verified). **Pattern for
Type B:** gaps declared and machine-identifiable; status honestly nonconformant; sidecar
relocation is a *registered gap*, not quiet practice (Type-B schema has no sidecar
equivalent — use `declared_gaps[]`).

### 14.3 Anti-example — `phys-demo/problem_8a/inkling`

Imports all three libraries yet: `Build()` **panics** (`gr` namespace validation rejects
assumption ID `medium_nonmagnetic_zeroT`) and no test ever calls it; the protocol test is
an empty stub; `numerical_check_test.go` asserts hardcoded placeholders with a wrong unit
conversion (g stated as `0.1 eV^-1` vs true `10^-19 eV^-1`); a derivation step invents
rule ID `delta_substitution_analytic`, impossible in the closed registry; axion physics
forced under the GR framework namespace. All green anyway. **Every §12 safeguard exists
because of this case.** Status: pending disposition (repair so entry points run, or
reclassify) — until then, treat as NONCONFORMANT exemplar, never as a model.

### 14.4 Honest solver-local pattern — `phys-demo/problem_8a/muse+deepseek`

Correct physics with zero library imports — the legitimate-boundary exemplar, not an
adoption model: exact-rational arithmetic via stdlib `big.Rat` (would be `physmath.NewRational`
under §8.1), float64 evaluation kept explicitly local, derivation steps honestly declared
non-executed, growth-gate note and `prefactor_normalization` limitation exemplary
(decision record G-11). **What to copy:** the executed/represented/local boundary
discipline, `TestUnitSlipSensitivity`-style unit-error tests, declared limitations.
**What §8.1 now adds:** it would additionally require either `infrastructure_use`
entries or explicit `deviations[]` explaining why reuse did not apply — the honest
deviation path, which was previously unmeasured.

### 14.5 Paper-scale process precedent

For Type A workflows, follow `phys-trace` protocol material (`M1_checklist_v2.md`,
`protocol/prompts/`) and DSG — referenced, never copied into this document.

---

## 15. Glossary and status vocabulary

| Term | Definition | Source |
|---|---|---|
| **Type-A / formalization package** | Theorem/hypothesis/derivation/paper formalization under the exact GR contract | Spec §1.4; ADR §2.1 |
| **Type-B / benchmark-solution package** | Solution of a specified benchmark under `phys-demo/benchmark-solution-manifest/1` | ADR §2.2; §11 |
| **`PROTOCOL GAP`** | Mandatory report when the GR protocol cannot express a fixture without a new semantic package-protocol field (Type A only) | Spec §1.4/§24.4 |
| **Evidence level** | `Represented → Structurally validated → Executed → Replayed/verified → Physically interpreted → Empirically supported` | phys-trace AGENTS §8; §10 |
| **Capability state** | `SUPPORTED / PARTIALLY SUPPORTED / REPRESENTATIONAL ONLY / DEFERRED / UNAVAILABLE` | phys-trace AGENTS §8; §5 |
| **Claim status** | `literature-claim-requiring-verification / candidate-construction-incomplete / open-conjecture / deferred / proved` | phys-lib `gr/ltb/claim.go` |
| **Formalization outcome** | `PROJECTABLE / PROJECTED_AS_HYPOTHESIS / PARTIALLY PROJECTABLE / UNREPRESENTABLE / REJECTED / REQUEST_REVISION` | Arch §26.3 |
| **Conformance status** | `CONFORMANT` / `NONCONFORMANT` — computed by the package test; a registered gap is not conformance | §11.3 |
| **`infrastructure_use`** | Adoption record proving reuse with exercised code path (§8.1) | this protocol |
| **Deviation record** | Justified non-reuse: capability sought, API examined, reason, smallest local alternative (§8.3) | this protocol |
| **Solver-local computation** | Numerical/other computation executed outside Phys libraries, never attributed to them (§9.2) | this protocol |
| **Growth Gate** | Post-MVP evidence-driven route for missing capability ownership (pure-math → PhysMath; physical → phys-*; …) | Arch §30.2–30.3 |

---

*End of `phys-demo/PROTOCOL.md`. Authority: architect's accepted decision record
(`TOE-v4/plan2/plan.md`) → ADR (`TOE-v4/plan2/adr_package_types.md`) → this document,
all subordinate to Architecture v1.5, Representation Spec v1.2, and the frozen Phys
authority.*
