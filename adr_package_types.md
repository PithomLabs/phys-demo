# ADR — Two Package Types: Formalization Packages vs Benchmark-Solution Packages

**Status:** ACCEPTED (architect decision, Task 3 review items 3/9, Part 2 item 9)
**Date:** 2026-10-18
**Scope:** Narrowly scoped architecture decision establishing a second package-level
contract class. It does not amend Representation Specification v1.2, Architecture v1.5,
or any frozen repository.
**Preceded by:** `plan.md` (Task 3 audit + accepted decision record),
`phys-task3-decision-snapshot.md`.
**Followed by:** `phys-demo/PROTOCOL.md` (this ADR is its authority for the package-type
distinction).

Normative language (MUST/MUST NOT/REQUIRED/MAY) follows Representation Specification
v1.2 §0.2.

---

## 1. Context and demonstrated problem

The Task 3 audit (`plan.md`) established, with executed evidence:

- The governing representation contract for formalization packages is Representation
  Specification v1.2 §1.4: canonical packaging `AGENTS.md + manifest.json +
  protocol/derivation_test.go`, where `manifest.json` MUST use **the exact schema and
  loading convention of the General Relativity example**, and any unexpressible need
  MUST be reported as a `PROTOCOL GAP` rather than a schema extension.
- No document defines a package contract for **benchmark-problem solutions** — solving a
  specified problem (e.g. TheoPhysEval Problem 8(a)) and recording parameters,
  assumptions, numerical conventions, solver-local computation, and answer-level
  verification.
- Consequently `phys-demo` operates under no manifest authority at all:
  - `problem_8a/muse+deepseek` (scientifically the strongest solution) produced its own
    ad hoc manifest key set — neither the GR schema nor a registered deviation;
  - `problem_8a/inkling` declared entry points that panic when called and a protocol
    test that is an empty stub, yet every suite reports PASS;
  - zero infrastructure adoption in one variant and meaningless import in the other both
    pass all existing acceptance signals.
- The architect's accepted decision (decision record Part 1, item 3) requires resolving
  this with **two package types**, via an explicit narrowly scoped ADR, because the
  specification imposes an exact-schema requirement on formalization packages and no
  exception may be created silently.

## 2. Decision

Two distinct package types exist as **separate normative contracts**, each with its own
governing manifest schema and conformance rules.

### 2.1 Type A — Formalization package (unchanged)

- Trigger: the package's declared purpose is to formalize a theorem, hypothesis,
  derivation, or physics-paper content — i.e. it claims content inside the formalization
  boundary of Specification v1.2 §1.4/§1.5.
- Contract: Representation Specification v1.2 §1.4 verbatim. Canonical layout
  `AGENTS.md + manifest.json + protocol/derivation_test.go`, exact GR schema and
  loading convention, executable protocol test exercising the formal content,
  `PROTOCOL GAP` as the only escape valve.
- This ADR **does not amend, weaken, or exempt** any part of that contract. The GR
  example remains the protocol precedent; no second protocol format for formalization
  packages is created anywhere by this ADR.

### 2.2 Type B — Benchmark-solution package (new, package-level contract only)

- Trigger: the package's declared purpose is to solve and document a specified physics
  benchmark (TheoPhysEval or equivalent), and it asserts **no formalization-status claim
  for any contained artifact** (see §3).
- Contract: a separately specified, rigorously validated manifest schema, identified as:

  ```text
  schema_id: phys-demo/benchmark-solution-manifest/1
  ```

  specified field-by-field in `phys-demo/PROTOCOL.md` §11 (allowed fields, required
  fields, validation rules, package-type declaration). Its required capabilities:
  problem identity and source specification; parameters with units or stated unit
  conventions; assumptions, conditions, regime; theoretical equations and derivation
  dependencies; references to Phys infrastructure representations and actually reused
  APIs (`infrastructure_use`); solver-local calculations; verification methods, results,
  limitations, open obligations; registered deviations; declared protocol gaps.
- **Package-level contract only.** This schema is metadata validation for demo packages.
  It is **not** part of the canonical PhysArtifact family registry, creates no
  PhysArtifact family, no PhysMath type, no executable symbolic rule, no dimensional
  semantics, and no numerical infrastructure.

### 2.3 Composition by construction

The two contracts are composed by **reference and separation**, never by merger: a
benchmark-solution package MAY reference a conformant formalization package (or any
artifact) by `ArtifactRef`/path; it MAY NOT carry that package's formalization
composition inside its own manifest.

## 3. Anti-evasion rule for mixed-purpose packages

**Question this section answers:** *what happens when a benchmark-solution package also
contains a paper-scale theorem formalization?*

1. **Category labels never control.** Declaring `package_type:
   benchmark-solution-manifest/1` does not exempt any contained artifact from
   Specification §1.4. The applicable contract is determined by the package's declared
   purpose **and by the claimed status of its contained artifacts**, not by the label on
   the manifest.
2. **Formalization claims pull in the full contract.** If any contained artifact claims
   formalization status under §1.4/§1.5 — a formally represented theorem, hypothesis, or
   paper derivation, a `Theorem`/paper `Derivation` asserted as formally represented and
   protocol-replayed — then that artifact is subject to the formalization contract, and
   the package MUST choose exactly one of:
   - **Option A (preferred):** relocate the formalization content into a separate,
     conformant Type-A package (sibling or nested directory following §1.4 exactly),
     referenced from the benchmark package. The two contracts compose by reference; the
     benchmark manifest MUST NOT restate the formalization composition.
   - **Option B:** the containing package itself adopts the full Type-A contract
     (GR-exact manifest schema and protocol). Benchmark metadata that does not fit the
     GR schema without a new semantic package-protocol field MUST be reported as a
     `PROTOCOL GAP` per §1.4 — never added as a schema extension, never moved into an
     unregistered sidecar.
3. **Splitting by artifact, not by directory naming.** Under Option A the formalization
   artifacts (manifest keys, expected results, claims of formalized content) remain
   governed by §1.4 in their own package; the benchmark package retains only benchmark
   metadata and references.
4. **Evasion is nonconformance.** A package that declares Type B while asserting
   formalization-status claims for contained artifacts, without satisfying §1.4 through
   Option A or B, is NONCONFORMANT under this ADR, and its conformance test MUST fail.
5. **Honest Type-B packages are not second-class.** A benchmark package that correctly
   records its derivation as solver-local/represented-only (making no formalization
   claim) is fully governed by Type B; it is not thereby required to become a Type-A
   package. This preserves §1.4's trigger as written while giving benchmark work a real,
   enforceable contract instead of none.

**Conflict check (recorded rather than assumed):** Specification §1.4's sentence
"This specification intentionally does not invent a competing manifest schema" and
§24.4's `PROTOCOL GAP -> do not invent a second protocol` govern *formalization
packages* — their purpose is preventing a substitute protocol for the GR precedent. This
ADR reads them as scoped to that purpose and does **not** override them: for Type A they
remain fully binding. Creating a package class for a different, non-formalization
purpose (answering a benchmark) is the narrowly scoped decision the architect authorized
in decision record Part 1 item 3; the interpretive step is recorded here so it is
explicit rather than silent. If a future reviewer reads §24.4's prohibition as extending
to *any* second package contract anywhere, that reading would nullify decision item 3
itself; the exact sections are named (§1.4 sentence above, §24.4) and the resolution
required is confirmation of this scope reading — none is assumed further.

## 4. Governing authorities

In precedence order (as fixed by decision record Part 1 items 3/6/7):

1. Frozen Phys authority (Specification §0.1) — for existing Phys semantics only.
2. Cross-project Architecture v1.5 (PM-ARCH-FOUNDATION-001) — dependency topology,
   identity, numerical boundary (§25), growth gate (§30), No-Invention Matrix (§32.2).
3. Representation Specification v1.2 (PM-REP-001) — representation semantics,
   §1.4 formalization contract, §24.4 stop rule. Document-control note: header says
   v1.2, body contains "v1.4 wins" references — recorded as a nonblocking anomaly
   (decision item 6); this ADR cites section numbers, not version claims.
4. This ADR — the package-type distinction only.
5. `phys-demo/PROTOCOL.md` — benchmark workflow and the Type-B schema, subordinate to
   all of the above.
6. Repository `AGENTS.md` files and `DERIVATION_SPEC_GUIDELINES.md` (canonical
   `phys-trace` copy) — process authority for derivation work.
7. READMEs, capability maps, plans — descriptive only; never executable authority.

## 5. Scope and explicit non-goals

In scope: defining the two package types, the Type-B schema identity, the mixed-purpose
composition rule, conformance expectations, and the boundary statements in §6/§7.

Explicit non-goals (decision record Part 1 items 9/10, Part 2 items 3/4/8):

- No new numerical engine; no `phys-calc` package created; numerical computation
  remains outside the representation layer.
- No dimensional-semantics system; unit handling stays metadata-only.
- No new PhysMath types, PhysArtifact families, executable symbolic rules, or registry
  changes.
- No export or cross-repository sharing of `phys-lib/internal/manifestschema`; no new
  module dependencies; no `go.mod`/`go.work` changes.
- No frozen-repository changes; no git initialization of `phys-demo`.
- No Jev, autonomous discovery, interpretation services, or research-layer expansion.
- No new claim-status vocabulary; the protocol reuses existing vocabularies.

## 6. Consequences

**Validation.** Type A validation is unchanged (GR-schema conformance + executable
protocol, as in `phys-lib/internal/manifestschema` for phys-lib fixtures). Type B
validation MUST be performed by a per-package strict validator in `phys-demo` (a
`manifest_test.go` checking the exact schema of `PROTOCOL.md` §11), because
`phys-lib/internal/manifestschema` is repository-internal and cannot be imported from
`phys-demo` without an unauthorized cross-repository dependency. **Known limitation,
recorded here:** until a separately authorized shared checker exists, Type-B conformance
is enforced per package rather than by one central registry; `PROTOCOL.md` MUST state
this limitation rather than imply a shared validator. Registered gaps MUST be declared
in the manifest and machine-identifiable; any unresolved mandatory gap forces a
NONCONFORMANT status that the conformance test MUST compute, not the manifest
self-declare.

**External users.** An agent can now determine applicability by declared purpose: solving
a specified problem → Type B (`phys-demo/PROTOCOL.md`); formalizing a theorem/paper →
Type A (Specification §1.4 + `DERIVATION_SPEC_GUIDELINES.md`). Mixed content follows
§3. Ambiguity defaults to Type A (the stricter contract), never to evasion.

**Negative consequence accepted:** two schemas exist. This is the accepted cost of giving
benchmark work an enforceable contract (decision item 3); the anti-evasion rule and
reference-only composition keep the contracts from blurring.

## 7. Acceptance criteria for this ADR

1. `phys-demo/PROTOCOL.md` declares the Type-B schema id
   `phys-demo/benchmark-solution-manifest/1` with exact allowed/required fields and
   validation rules, and states the per-package validation limitation.
2. The protocol contains the §3 anti-evasion rule verbatim in substance: category label
   never controls; formalization claims trigger §1.4; ambiguous packages resolve to
   Type A; evasion fails the conformance test.
3. A check shows no requirement in the protocol forces numerical computation from Phys
   infrastructure, and every solver-local numerical step is attributable outside
   PhysMath/PhysLib.
4. Specification §1.4 remains byte-unchanged and is quoted, not paraphrased, wherever
   the protocol relies on it.
5. No file outside the authorized planning/documentation paths changes (verified by
   repository status before/after).

## 8. Explicit statements (normative)

- **The existing formalization contract remains unchanged.** Representation
  Specification v1.2 §1.4 for Type-A packages is untouched by this ADR: exact GR schema,
  canonical layout, executable protocol test, and the mandatory `PROTOCOL GAP` procedure
  all remain fully in force.
- **Numerical computation remains outside the representation layer.** Architecture v1.5
  §25 and Specification §23.2 are unchanged: Phys represents theoretical mathematics,
  declared physical semantics, and derivation structure; a future `phys-calc` may sit
  on top but is not authorized, specified, or created here.
