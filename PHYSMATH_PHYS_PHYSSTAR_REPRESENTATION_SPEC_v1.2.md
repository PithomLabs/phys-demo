# PhysMath + phys-* Representation Specification
## Phase-1 Implementation-Planning Baseline

**Document ID:** PM-REP-001  
**Version:** 1.2
**Status:** **IMPLEMENTATION-PLANNING BASELINE — ARCHITECTURALLY CLOSED FOR PHASE-1**  
**Role:** Normative representation specification for PhysMath and `phys-*` external physics libraries.  
**Supersedes:** PM-REP-001 v1.1 and all earlier PM-REP-001 versions. Where this document conflicts with an earlier version, v1.4 wins.  
**Primary scope:** PhysMath, `phys-artifact`, `phys-*`, and the executable reference protocol.  
**Secondary systems:** TOE-Agent, Numerical, and research/provenance/trust workflow are explicitly deferred and MUST NOT block Phase-1 implementation.

> **Primary purpose.** Define the smallest coherent representation architecture needed to implement PhysMath and `phys-*` libraries around a nontrivial physics-paper outcome: the representation of equations, assumptions, definitions, mathematical transformations, physical interpretation, and a complete derivation of an equivariance theorem. The first implementation is intentionally developed by refactoring the existing `phys-gr` implementation so that the common denominator and domain-specific differences are tested against real code rather than designed in isolation.

> **Central rule.** PhysMath represents mathematics. `phys-*` represents physical meaning, assumptions, framework-specific definitions, equations, and theorem/derivation context. Neither layer silently answers the other's question.

> **Formalization rule.** A theorem, hypothesis, derivation, or physics-paper formalization is not considered formally represented merely because prose, LaTeX, or an abstract diagram exists. The formal representation boundary is the executable project protocol established by the General Relativity example: `AGENTS.md`, `manifest.json`, and `protocol/derivation_test.go`. Future formalization packages MUST follow that protocol shape exactly. The existing GR example is the protocol precedent and MUST be reused or refactored rather than replaced by a newly invented package protocol. Otherwise the material remains prose or descriptive input to a future formalization process.

> **Complexity rule.** The architecture is governed throughout its lifecycle by three principles: (1) coupling means conflating different concerns and is fixed by separation of concerns; (2) cohesion means separating concerns that must be unified and is fixed by explicit interfaces; (3) core/specific mixing means conflating common structure with domain-specific differences and is fixed by identifying the true common denominator for the core and moving variation into specific frameworks.

---

# 0. Scope and authority

## 0.1 Frozen Phys authority

The existing Phys project remains authoritative for its existing semantics:

```text
Module:              github.com/PithomLabs/phys
Specification:       specs_v2_3.md
Process authority:   AGENTS.md
Operational freeze:  exact frozen repository state supplied to implementation
```

This specification MUST NOT modify the Phys kernel, dependency graph, or public API.

The frozen Phys authority wins whenever a statement about existing Phys behavior conflicts with this specification. Such a conflict is an architecture issue and MUST NOT be resolved by changing Phys during Phase 1.

## 0.2 Normative language

**MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, and **SHALL NOT** are normative.  
**SHOULD**, **SHOULD NOT**, and **MAY** are advisory.

## 0.3 Phase-1 scope

Phase 1 implements only:

```text
phys-artifact
PhysMath
phys-* external physical-domain libraries
```

The implementation target is representation plus executable reference protocol, not the complete research platform.

## 0.4 Explicitly deferred

The following are separate systems and MUST NOT block Phase-1 implementation:

```text
TOE-Agent research orchestration
Numerical execution
ComputationalContract execution
research lifecycle state machines
formalization search/orchestration algorithms
trust/corpus promotion workflow
scientific adjudication
paper-management workflow
human/AI agent orchestration
semantic embedding infrastructure
cross-domain autoformalization infrastructure
```

A Phase-1 implementation MUST NOT introduce these concerns merely because a future platform will eventually need them.

## 0.5 v1.4 closure basis

Version 1.4 incorporates the consolidated adversarial findings from the v2 and v3 review rounds and closes the remaining implementation-level ambiguities without reopening the macro-architecture. It also establishes the plan-first implementation handoff and the existing `phys-gr` refactoring path.

The closed architectural decisions are:

```text
Phase-1 / Phase-2 separation
PhysMath purity
phys-* physical-semantics ownership
many-frameworks principle
frozen Phys boundary
artifact family registry
canonical primitive encodings
ParameterID / BoundParameterRef rules
relation value-type encoding
closed Phase-1 rule registry
ProducerNamespace in durable rule records
closed rule arity and input families
deterministic transformation traversal
explicit EXPR-APPLY boundary
explicit base/specialized dependency model
single TransformResult type
closed ArtifactRef target map
closed Function binding/capture rules
closed equivariance reference graph
GR-style executable protocol requirement
```

No coding agent is permitted to reopen these decisions as implementation choices. A coding agent MAY report a demonstrated architectural blocker discovered while preparing the implementation plan, but it MUST NOT silently choose a new architecture.

## 0.6 Existing `phys-gr` as refactoring substrate

The existing `phys-gr` implementation predates the implementation of PhysMath and the external `phys-*` library family. It is therefore the preferred concrete substrate for the Phase-1 implementation plan.

The roles are deliberately distinct:

```text
PM-REP-001 v1.4
    = normative architecture and representation contract

General Relativity example
    = normative protocol packaging precedent

existing phys-gr implementation
    = refactoring substrate and implementation evidence

frozen Phys
    = authority only for existing frozen Phys semantics
```

Before coding, the implementation plan MUST inventory existing `phys-gr` responsibilities and map them to:

```text
phys-artifact
PhysMath core
phys-gr-specific semantics
protocol/package concerns
legacy or removable implementation detail
```

The plan MUST explicitly apply the core/specific rule: extract only the true common denominator into PhysMath or other shared infrastructure, while keeping GR-specific assumptions, definitions, interpretations, and algebra in `phys-gr`.

Existing `phys-gr` behavior MUST NOT silently override this specification. The specification MUST NOT be distorted merely to preserve an existing implementation detail. Any contradiction is handled by the §24.4 stop rule.


---

# 1. Research-to-representation boundary

## 1.1 Physics paper is the unit of research

A physics paper is the normal unit of research context.

A paper may contain:

```text
motivation
analysis
hypotheses
interpretation
experimental discussion
prose
mathematical definitions
assumptions
equations
derivations
theorems
```

The Phase-1 representation substrate does **not** model the whole paper as a PhysMath ontology object.

The representation target is the paper's formal mathematical/physical payload:

```text
definitions
assumptions
mathematical objects
equations
transformations
derivation steps
theorem statements
framework-specific physical interpretations
```

A human or AI research process may extract these objects from a paper. The extraction process itself is outside PhysMath semantics.

## 1.2 Analysis / Representation / Interpretation

The Phase-1 boundary is:

```text
ANALYSIS
    humans + AI agents determine what a paper appears to say
    and what hypotheses/theorems may be worth representing

REPRESENTATION
    PhysMath + phys-* encode the resulting mathematical and
    framework-specific structures deterministically

INTERPRETATION
    phys-* declares what physical meaning a framework assigns
    to represented mathematical structures
```

Analysis remains outside PhysMath semantics.

Interpretation is explicit in `phys-*`; it is never inferred from names or mathematical shape.

## 1.3 Equation-first rule

A paper's formal output is represented primarily through:

```text
Equation
EquationSet
Definition
MathCondition
Derivation
Theorem
```

Prose may exist in an enclosing research or paper-management system, but prose MUST NOT be required to reconstruct mathematical identity or replay a derivation.

## 1.4 Formalization is protocol execution

For any future theorem, hypothesis, derivation, or physics-paper formalization, the canonical project packaging is:

```text
<formalization-package>/
    AGENTS.md
    manifest.json
    protocol/
        derivation_test.go
```

The structure and role of these files MUST follow the same model established by the General Relativity example.

`AGENTS.md` is process/instruction authority for the formalization package.

`manifest.json` declares the package's formalization composition according to the established protocol model.

`protocol/derivation_test.go` is the executable derivation/protocol test and MUST exercise the formal content rather than merely checking static files.

Phase 1 does not invent a second protocol format for the equivariance theorem.

The package-level `manifest.json` MUST use the exact schema and loading convention of the existing General Relativity example. This specification intentionally does not invent a competing manifest schema. The coding agent MUST inspect the GR example and reuse its manifest structure, required keys, value forms, and protocol-loading pattern.

`manifest.json` is package/protocol metadata; it is not a PhysMath durable artifact and is not part of any ArtifactRef canonical payload. If the existing GR protocol cannot express the Phase-1 fixture without inventing a new semantic package-protocol field, the coding agent MUST report a `PROTOCOL GAP` rather than silently extending the protocol.

A later theorem, hypothesis, or paper formalization MAY introduce different mathematical artifacts and `phys-*` frameworks, but it MUST retain this protocol boundary.

## 1.5 Prose-only status

The following are descriptive rather than formal by themselves:

```text
natural-language theorem statement
LaTeX equation
diagram
LLM-generated explanation
plain Markdown derivation
JSON without the canonical artifact/protocol contract
```

They may be inputs to an analysis/formalization process.

They become a formal research representation only when the corresponding mathematical and physical structures are encoded according to this specification and exercised through the project protocol.

---

# 2. Constitutional architecture

## 2.1 PhysMath

PhysMath owns theory-neutral mathematics:

```text
mathematical types
mathematical objects
expressions
functions
maps
operators
relations/equations
conditions
transformations
derivations as mathematical structure
canonicalization
mathematical equality
canonical serialization
```

PhysMath MUST NOT own:

```text
physical interpretation
physical theory names
physical assumptions
experimental truth
research status
scientific preference
paper provenance
trust/corpus status
```

## 2.2 phys-* libraries

A `phys-*` library represents one specific physical framework/domain.

It owns:

```text
physical types
physical assumptions
framework vocabulary
framework-specific definitions
framework-specific equations
framework-specific theorem statements
framework-specific interpretation of mathematical objects
framework-specific mappings into PhysMath
framework-local validation logic
```

A `phys-*` library MAY define a local mathematical algebra when PhysMath does not yet contain the required abstraction, but that algebra remains package-local and MUST NOT dynamically extend the PhysMath AST or durable PhysMath family registry.

## 2.3 Many answers are legal

The architecture deliberately permits:

```text
same PhysMath object
    -> multiple phys-* interpretations
```

No global mechanism chooses one physical interpretation as universally correct.

A consumer that requires physical meaning MUST identify the framework explicitly.

## 2.4 Representation is not truth

The following are representation successes, not scientific truth claims:

```text
canonicalization succeeds
serialization succeeds
derivation replays
equation matches
theorem representation validates
physical framework record validates
```

A `phys-*` library represents a framework's claims and rules. It does not prove that the framework is empirically correct.

## 2.5 Complexity architecture

The following principles are permanent architectural heuristics:

### Coupling

Coupling occurs when separate concerns are conflated.

```text
Fix -> separation of concerns
```

### Cohesion

Cohesion problems occur when concerns that must be unified are separated without an interface.

```text
Fix -> identify the explicit interface that unifies them
```

### Core / specific split

Core/specific mixing occurs when common structure and framework-specific differences are placed in the same layer.

```text
Fix -> identify the common denominator for the core
       and move differences into the specific domain/framework
```

These principles apply recursively to:

```text
research concepts
mathematical representation
physical semantics
framework design
APIs
Go packages
schemas
tests
future domain extensions
```

---

# 3. Dependency boundary

The Phase-1 production dependency graph is:

```text
PhysMath
    -> phys-artifact

phys-*
    -> PhysMath
    -> phys-artifact
    -> Phys (existing public API only, when required)

Phys
    -> no Phase-1 dependency changes
```

There is no dependency from PhysMath to any `phys-*` package.

`phys-artifact` is mechanical infrastructure only.

It MUST NOT own:

```text
physics ontology
math semantics
interpretation
trust
research state
```

A `phys-*` framework-to-framework dependency is permitted only through the explicit framework dependency mechanism defined in §14.8 and §14.9. Framework dependency is never inferred merely because two packages happen to be imported in Go source.

The dependency graph MUST remain acyclic.

---

# 4. Artifact and identity model

## 4.1 Typed values versus durable references

In-memory composition uses concrete typed values:

```text
physmath.Expression
physmath.Function
physmath.Relation
phys-foo.<framework type>
```

A durable cross-object relationship uses `ArtifactRef`.

```text
typed value != ArtifactRef
```

## 4.2 ArtifactRef

The canonical Phase-1 type is:

```go
type ArtifactRef struct {
    Namespace     string
    SchemaVersion string
    ArtifactType  string
    Digest        [32]byte
}
```

Field order is normative.

Identity is exactly:

```text
(Namespace, SchemaVersion, ArtifactType, Digest)
```

The digest alone is not sufficient to identify an artifact family instance.

## 4.3 Canonical digest

For canonical payload `P`:

```text
Header =
    Len32(Namespace UTF-8 bytes)
    || Namespace UTF-8 bytes
    || Len32(SchemaVersion UTF-8 bytes)
    || SchemaVersion UTF-8 bytes
    || Len32(ArtifactType UTF-8 bytes)
    || ArtifactType UTF-8 bytes

Digest = SHA-256(Header || P)
```

There is exactly one Phase-1 artifact digest formula.

## 4.4 Embedded ArtifactRef

An embedded `ArtifactRef` is encoded exactly as:

```text
String(Namespace)
|| String(SchemaVersion)
|| String(ArtifactType)
|| Digest32(Digest)
```

Expanded:

```text
UInt32BE(byte_length(Namespace))
|| Namespace UTF-8 bytes
|| UInt32BE(byte_length(SchemaVersion))
|| SchemaVersion UTF-8 bytes
|| UInt32BE(byte_length(ArtifactType))
|| ArtifactType UTF-8 bytes
|| Digest[32]
```

No JSON, hexadecimal text, pointer value, or bare digest is canonical embedded-ref material.

## 4.5 Primitive encoding

```text
String(T)      = UInt32BE(byte_length(T)) || UTF-8(T)
Bytes(B)       = UInt32BE(len(B)) || B
UInt8          = 1 byte unsigned integer
UInt16         = 2-byte unsigned big-endian integer
UInt32         = 4-byte unsigned big-endian integer
UInt64         = 8-byte unsigned big-endian integer
Bool           = 0x00 or 0x01
Digest32       = exactly 32 bytes
Optional(T)    = 0x00 when absent; 0x01 || CanonicalEncoding(T) when present
Sequence(T)    = UInt32BE(count) || elements in declared order
Set(T)         = UInt32BE(count) || elements sorted by canonical element bytes
Enum           = String(canonical-token)
TaggedUnion    = UInt8(discriminator) || selected variant payload
```

A `Set` MUST reject duplicate canonical elements.

`Optional<T>` uses a one-byte presence discriminator and contains no additional framing around the contained canonical value beyond the contained type's own canonical encoding.

`Sequence<T>` and `Set<T>` use a four-byte element count.

A `Set<T>` is sorted by bytewise lexicographic comparison of the exact canonical encoding of each element. The resulting order is deterministic across implementations.

## 4.6 BigInt

`BigInt` is an exact signed integer used by mathematical rational values.

Canonical encoding:

```text
Sign:             UInt8
MagnitudeLength:  UInt32BE
Magnitude:        minimal unsigned big-endian bytes
```

Sign values are:

```text
0x00 = non-negative
0x01 = negative
```

The magnitude has no redundant leading zero bytes.

Zero is encoded as:

```text
0x00 || UInt32BE(0)
```

Negative zero is invalid.

The magnitude length is the byte length of the minimal magnitude, not the mathematical bit length.

## 4.7 Floating-point boundary

`Float64Bits` is not permitted in exact PhysMath mathematical identity payloads.

It MAY be used by future numerical schemas.

Phase 1 MUST NOT introduce floating-point approximation into mathematical canonical identity.

## 4.8 Canonical record rule

A record's fields are encoded in the exact order declared by its schema.

No field names are included in canonical bytes unless the schema explicitly defines them as data.

Unknown fields are invalid for schema version `1`.

## 4.9 Collection ordering rule

For every collection the owning schema MUST declare `SET` or `SEQUENCE`.

For a `SET`, elements are sorted by lexicographic byte comparison of their exact canonical element encoding.

Traversal order MUST NOT determine canonical Set order.

## 4.10 General BoundParameterRef encoding

Whenever a `BoundParameterRef` appears in a canonical payload, its canonical representation is:

```text
Encode(BoundParameterRef)
    = EncodeArtifactRef(DeclaringArtifactRef)
    || Digest32(ParameterID)
```

No `Bytes()` wrapper surrounds the 32-byte ParameterID.

This encoding is normative wherever a later schema references `BoundParameterRef` directly.

---

# 5. Namespace and artifact-family rules

## 5.1 Namespace grammar

A namespace is ASCII lowercase and consists of one to sixty-four characters from:

```text
[a-z0-9][a-z0-9._-]{0,63}
```

No spaces or uppercase letters are permitted.

Reserved Phase-1 namespace classes are:

```text
physmath
phys
phys-<domain>
```

Namespace validity is syntactic in Phase 1. Any namespace beginning with `phys-` and satisfying the namespace grammar is a valid `phys-*` namespace. Phase 1 does NOT define or require a global namespace registry.

A `phys-*` namespace becomes a conforming framework namespace when its package publishes a valid `<namespace>/framework/1` FrameworkArtifact and implements the six framework-owned families in §5.6. The package MUST use its own namespace consistently as `PackageNamespace` in its FrameworkArtifact and as the namespace portion of its six family identifiers.

The package itself establishes namespace ownership through its FrameworkArtifact and conformance tests. A central semantic registry is deliberately not part of Phase 1.

Examples:

```text
phys-qm
phys-bohm
phys-gr
phys-cst
phys-cdt
```

`phys-artifact` has no semantic artifact namespace.

## 5.2 Artifact-type grammar

Artifact types consist of one or more lowercase path segments, each matching the namespace identifier grammar; `/` separates segments.

Maximum encoded UTF-8 length is 128 bytes.

## 5.3 Schema-version grammar

`SchemaVersion` is a String and all Phase-1 artifact schemas use:

```text
"1"
```

Schema-version strings inherit the lowercase identifier grammar for Phase-1 version identifiers.

A changed canonical payload schema requires a new schema version and MUST NOT silently reinterpret an existing schema version.

## 5.4 Free-form identifier grammar

Unless a schema explicitly declares another grammar, the following identifiers use the ASCII-safe grammar:

```text
[a-z0-9][a-z0-9._-]{0,63}
```

with a maximum encoded UTF-8 length of 128 bytes:

```text
FrameworkID
AssumptionID
DefinitionID
EquationID
TheoremID
InterpretationID
FunctionID
OperatorID
PhysicalType
```

Human-readable display names are not identity-bearing unless the schema explicitly says they are. `FunctionID` is reserved vocabulary only; no Phase-1 `Function` payload contains it and implementations MUST NOT add it merely because it appears in the identifier grammar.

## 5.5 FrameworkVersion

`FrameworkVersion` is a non-empty opaque UTF-8 String with a maximum encoded length of 128 bytes. Its exact byte sequence is package-owned data and is not interpreted by PhysMath.

It is package-owned data, not a globally interpreted semantic version ontology.

Canonical encoding is exactly:

```text
String(FrameworkVersion)
```

The value MAY use forms such as `1`, `1.0.0`, or another package-defined release token, but the exact value is part of the framework artifact's canonical payload.

The reference fixture MUST pin its exact framework-version value in fixture data.

## 5.6 Durable artifact-family registry

The following durable artifact families are the complete Phase-1 registry.

### PhysMath families

```text
physmath/expression/1
physmath/function/1
physmath/function_application/1
physmath/map/1
physmath/tensor/1
physmath/index/1
physmath/collection/1
physmath/operator_application/1
physmath/transformation/1
physmath/relation/1
physmath/equation/1
physmath/equation_set/1
physmath/definition/1
physmath/condition/1
physmath/derivation/1
physmath/theorem/1
```

### phys-* families

For every conforming `phys-*` namespace `<ns>`:

```text
<ns>/framework/1
<ns>/assumption/1
<ns>/definition/1
<ns>/equation/1
<ns>/theorem/1
<ns>/interpretation/1
```

No other durable family may be introduced during Phase 1.

The following PhysMath sort concepts are **in-memory only in Phase 1** because they have no Phase-1 durable family:

```text
Functional
FunctionalApplication
Operator
Limit
AsymptoticTruncation
DerivationObligation
```

They MUST NOT be serialized into invented families.

## 5.7 Normative ArtifactRef target-family registry

The following table is the complete Phase-1 target-family contract for every durable `ArtifactRef`-typed field. A field not listed here does not acquire an open target by implication.

```text
Function.BodyRef
    -> physmath/expression/1

FunctionApplication.FunctionRef
    -> physmath/function/1

FunctionApplication.Args[*]
    -> physmath/expression/1

Map.DefinitionRef
    -> physmath/definition/1

Tensor.IndexRefs[*]
    -> physmath/index/1

Tensor.ComponentRefs[*]
    -> physmath/expression/1

Index.SpaceRef
    -> physmath/collection/1 when present

Collection.ElementRefs[*]
    -> physmath/expression/1

OperatorApplication.OperatorRef
    -> absent in Phase 1; non-absent value is invalid

OperatorApplication.Args[*]
    -> physmath/expression/1

Transformation.Inputs[*]
    -> exact family dictated by the selected Phase-1 rule in §17.2

Transformation.Output
    -> physmath/expression/1 for every Phase-1 executable rule

Equation.RelationRef
    -> physmath/relation/1

EquationSet.EquationRefs[*]
    -> physmath/equation/1

Definition.DefinedRef
    -> physmath/function/1 OR physmath/expression/1

Definition.DefiningRef
    -> physmath/expression/1

MathCondition.PredicateRef
    -> physmath/relation/1

DerivationStep.InputRefs[*]
    -> one of the six closed families allowed by the selected rule in §17.2:
       physmath/expression/1
       physmath/function/1
       physmath/relation/1
       physmath/equation/1
       physmath/definition/1
       physmath/condition/1
    The exact position-specific family MUST be validated against §17.2.

DerivationStep.OutputRef
    -> physmath/expression/1

DerivationStep.ConditionRefs[*]
    -> physmath/condition/1

Derivation.PremiseRefs[*]
    -> physmath/expression/1
       OR physmath/function/1
       OR physmath/relation/1
       OR physmath/equation/1
       OR physmath/definition/1
       OR physmath/condition/1

Derivation.ResultRefs[*]
    -> physmath/expression/1
       OR physmath/function/1
       OR physmath/relation/1
       OR physmath/equation/1
       OR physmath/definition/1

Theorem.PremiseRefs[*]
    -> physmath/expression/1
       OR physmath/function/1
       OR physmath/relation/1
       OR physmath/equation/1
       OR physmath/definition/1
       OR physmath/condition/1

Theorem.ConclusionRef
    -> physmath/relation/1

Theorem.AssumptionRefs[*]
    -> physmath/condition/1

Theorem.DerivationRef
    -> physmath/derivation/1

FrameworkArtifact.BaseFrameworkRefs[*]
    -> <valid phys-namespace publishing a FrameworkArtifact>/framework/1

PhysicalAssumption.ConditionRefs[*]
    -> physmath/condition/1

PhysicalDefinition.FrameworkRef
    -> <owning phys-namespace>/framework/1

PhysicalDefinition.MathematicalRefs[*]
    -> one of the 16 registered Phase-1 PhysMath durable families in §5.6

PhysicalDefinition.AssumptionRefs[*]
    -> <owning phys-namespace>/assumption/1
       OR an assumption family in the owning framework's dependency closure

PhysicalDefinition.DefinitionEquation
    -> physmath/equation/1 when present

PhysicalEquation.FrameworkRef
    -> <owning phys-namespace>/framework/1

PhysicalEquation.EquationRef
    -> physmath/equation/1

PhysicalEquation.AssumptionRefs[*]
    -> <owning phys-namespace>/assumption/1
       OR an assumption family in the owning framework's dependency closure

PhysicalEquation.ScopeRefs[*]
    -> <owning phys-namespace>/assumption/1
       OR an assumption family in the owning framework's dependency closure

PhysicalTheorem.FrameworkRef
    -> <owning phys-namespace>/framework/1

PhysicalTheorem.MathematicalTheoremRef
    -> physmath/theorem/1

PhysicalTheorem.AssumptionRefs[*]
    -> <owning phys-namespace>/assumption/1
       OR an assumption family in the owning framework's dependency closure

PhysicalTheorem.InterpretationRefs[*]
    -> <owning phys-namespace>/interpretation/1
       OR an interpretation family in the owning framework's dependency closure

PhysicalInterpretation.FrameworkRef
    -> <owning phys-namespace>/framework/1

PhysicalInterpretation.MathematicalRefs[*]
    -> one of the 16 registered Phase-1 PhysMath durable families in §5.6

PhysicalInterpretation.AssumptionRefs[*]
    -> <owning phys-namespace>/assumption/1
       OR an assumption family in the owning framework's dependency closure
```

The phrase `owning framework's dependency closure` has the exact meaning defined in §14.9. A `phys-*` namespace is valid when it satisfies §5.1 and publishes the required FrameworkArtifact; no global namespace registry is consulted. No other target family is permitted in Phase 1.

## 5.8 Family registration rule

An artifact family is identified by the complete triple:

```text
Namespace / ArtifactType / SchemaVersion
```

The family MUST appear in the registry above or in one of the six allowed framework families published by a conforming `phys-*` package.

No implementation may invent an unregistered durable family.

---

# 6. PhysMath mathematical IR

## 6.1 Closed core sort inventory

The Phase-1 PhysMath sort vocabulary is exactly:

```text
Expression
Function
FunctionApplication
Functional
FunctionalApplication
Operator
OperatorApplication
Tensor
Index
Map
Relation
Collection
Transformation
Limit
AsymptoticTruncation
DerivationObligation
```

This is a closed 16-sort vocabulary for PhysMath schema version `1`.

A new sort requires a new versioned specification. A `phys-*` package MUST NOT inject a runtime sort into PhysMath.

## 6.2 Implementation priority within the sort set

The reference implementation MUST first implement the subset required for the equivariance vertical slice:

```text
Expression
Function
FunctionApplication
Map
OperatorApplication
Relation
Transformation
Tensor
Index
Collection
```

The remaining sort concepts remain part of the schema vocabulary but have no Phase-1 executable semantics unless explicitly activated by a future versioned specification:

```text
Functional
FunctionalApplication
Operator
Limit
AsymptoticTruncation
DerivationObligation
```

A deferred sort has no Phase-1 durable family and is therefore in-memory-only for Phase 1.

## 6.3 Phase-1 operation lock

A sort being present in the vocabulary does not create an executable operation for it.

The only executable PhysMath transformations in Phase 1 are the seven rule identities defined in §11.5 and §17.

---

# 7. Expression representation

## 7.1 Closed node set

The Phase-1 Expression tagged union is exactly:

```text
0x01 Symbol
0x02 Rational
0x03 Add
0x04 Mul
0x05 Neg
0x06 Pow
0x07 Sqrt
0x08 Call
0x09 Relation
0x0A BranchSet
```

No other node tag is legal in schema version `1`.

## 7.2 Expression artifact family

A durable `physmath/expression/1` artifact contains exactly one complete expression tree.

Its canonical payload is:

```text
UInt8(root-node-tag)
|| node-specific canonical payload recursively
```

There is no wrapper object around the root node.

An embedded expression tree and a standalone `physmath/expression/1` artifact use the same canonical expression encoding for the represented tree.

## 7.3 Recursive payload rule

Every child expression is encoded recursively as:

```text
UInt8(NodeTag) || node-specific fields
```

No pointer, runtime type name, JSON representation, or Go memory layout participates in canonical bytes.

## 7.4 Symbol

```text
Symbol {
    DisplayName string
    ParameterID Optional<ParameterID>
    ConstantID  Optional<ConstantID>
}
```

Validation:

```text
durable symbol:
    exactly one of ParameterID or ConstantID is present

ephemeral unbound symbol:
    both are absent

both present:
    always invalid
```

A standalone expression MAY contain a ParameterID-bearing Symbol as a schematic/free parameter symbol. That Symbol is not globally bound merely by its ParameterID.

Any operation that interprets a ParameterID as a bound target MUST receive the relevant declaring Function context explicitly according to the applicable rule contract.

`DisplayName` is presentation-only with respect to mathematical meaning, but it IS identity-bearing because it participates in the canonical `Symbol` payload. Changing it can therefore change the `physmath/expression/1` ArtifactRef.

PhysMath MUST NOT inspect display names to assign physical meaning.

## 7.5 Rational

```text
Rational {
    Numerator   BigInt
    Denominator BigInt
}
```

The denominator is positive, the fraction is reduced, and zero is represented as `0/1`.

## 7.6 Add / Mul

```text
Add {
    Children []Expression // SEQUENCE at construction; canonical child order is deterministic
}

Mul {
    Children []Expression // SEQUENCE at construction; canonical child order is deterministic
}
```

Canonicalization recursively canonicalizes children, flattens nested same-kind nodes, and sorts commutative children by canonical expression bytes.

Arity normalization is:

```text
Add([])  -> Rational(0/1)
Mul([])  -> Rational(1/1)
Add([x]) -> x
Mul([x]) -> x
```

## 7.7 Neg

```text
Neg {
    Operand Expression
}
```

Construction preserves the node.

`EXPR-SIMPLIFY-001` may apply:

```text
Neg(Neg(x)) -> x
```

## 7.8 Pow

```text
Pow {
    Base     Expression
    Exponent Expression
}
```

PhysMath v1 permits an arbitrary Expression exponent.

A `phys-*` framework MAY impose a stricter exponent rule in its own contract.

## 7.9 Sqrt

```text
Sqrt {
    Radicand Expression
}
```

`Sqrt` is a distinct mathematical node. Phase-1 canonical identity preserves the explicit node kind.

## 7.10 Call

```text
Call {
    FunctionRef ArtifactRef
    Args        []Expression // SEQUENCE
}
```

`FunctionRef` MUST target `physmath/function/1`.

Argument order is semantic.

`Call` is representational structure and is not itself a theorem.

## 7.11 Relation node

An Expression-level relation is:

```text
Relation {
    Operator ComparisonOperator
    Left     Expression
    Right    Expression
}
```

Closed comparison vocabulary:

```text
EQ
NEQ
LT
LTE
GT
GTE
```

Expression-node canonical encoding is:

```text
0x09
|| Enum(Operator)
|| CanonicalExpression(Left)
|| CanonicalExpression(Right)
```

A relation is a mathematical proposition. It is not a trust assertion.

## 7.12 BranchSet

```text
BranchSet {
    Branches []Expression // SEQUENCE
}
```

Branch order is semantic until an explicit selection transformation occurs.

No branch is selected by name or physical convention.

---

# 8. Parameters and mathematical identity

## 8.1 ParameterID

```go
type ParameterID [32]byte
```

`ParameterID` identifies a mathematical binding slot within a declaration context.

It is not a global provenance identifier.

## 8.2 Phase-1 binder scope

Phase 1 supports **Function parameter binding** as the canonical binding mechanism.

Nested lexical binders are not part of the Phase-1 binder vocabulary and are deferred.

For a `Function` declaration:

```text
ScopePath = empty sequence
depth = 0
ParameterOrdinal = zero-based position in Function.ParameterIDs
```

Therefore:

```text
ScopePathBytes = UInt32BE(0)
```

The general lexical-scope encoding remains defined as:

```text
ScopePathBytes =
    UInt32BE(depth)
    || UInt32BE(scope_0)
    || ...
    || UInt32BE(scope_(depth-1))
```

but Phase-1 Function parameters always use depth zero.

## 8.3 ParameterID formula

The canonical formula is:

```text
ParameterID = SHA-256(
    String("physmath")
    || Bytes(ScopePathBytes)
    || UInt32BE(ParameterOrdinal)
)
```

Because `Bytes()` is length-prefixed, the resulting hash input is unambiguous.

`DisplayName` is excluded.

An artifact's own `ArtifactRef` is excluded. This prevents recursive identity dependence.

`ParameterID` is intentionally local slot identity. Two different Function artifacts may therefore contain the same ordinal-derived ParameterID. Cross-artifact identity is disambiguated by `BoundParameterRef`, which includes the declaring artifact.

For every durable `Function`, validation MUST require:

```text
Function.ParameterIDs[i] == DeriveParameterID(i)
```

for every declared ordinal `i`. A Function with ordinal-derived parameter IDs that do not match the stored values is invalid.

The ParameterID hash is a distinct mathematical-binding digest domain. It is NOT an ArtifactRef digest and MUST NOT be interpreted or reused as one. The exact framing `String("physmath") || Bytes(ScopePathBytes) || UInt32BE(ParameterOrdinal)` is intentional and MUST NOT be harmonized with the ArtifactRef digest framing in §4.3.

## 8.4 Alpha-renaming

Changing:

```text
t -> tau
```

MUST NOT alter ParameterID when the declaration position is unchanged.

## 8.5 BoundParameterRef

When a parameter is referenced from another durable artifact:

```text
BoundParameterRef {
    DeclaringArtifactRef ArtifactRef
    ParameterID         ParameterID
}
```

Canonical encoding is exactly:

```text
EncodeArtifactRef(DeclaringArtifactRef)
|| Digest32(ParameterID)
```

A bare ParameterID is insufficient when multiple declaring artifacts are in scope.

## 8.6 Function body binding rule

For every durable `Function` artifact:

```text
Every Symbol in BodyRef MUST have:
    ParameterID ∈ Function.ParameterIDs
OR
    ConstantID present
```

A `Function` body containing a ParameterID not declared by the Function is invalid.

An unbound Symbol with both ParameterID and ConstantID absent is not permitted in a durable Function body.

Within one expression artifact, equal ParameterID values denote the same positional parameter slot for that expression's current declaration context.

## 8.7 Capture-safety rule

`EXPR-APPLY-001` is capture-free by construction because every bound Function-body ParameterID is among the Function's declared parameters and all declared parameters are substituted by the corresponding argument.

Implementations MUST NOT alpha-freshen, rename, or otherwise mutate ParameterIDs during application.

No fresh-variable allocator, closure object, or capture-avoidance engine is introduced in Phase 1.

## 8.8 Standalone schematic parameters

A standalone top-level expression, relation, or equation MAY contain ParameterID-bearing Symbols representing schematic variables such as `g` and `x`.

Such Symbols are not global variables and do not imply a new global binder system.

The strict Function-body rule in §8.6 applies to durable Function bodies and operation-referenced contexts that require a declaration. Standalone schematic expressions MAY contain ParameterID-bearing Symbols without a declaring artifact. These are deliberately different validation regimes.

For operations that need a declaration context, the applicable rule MUST supply or validate the declaring Function explicitly.

This is the Phase-1 boundary between schematic mathematical expressions and general quantified theorem binding.

## 8.9 t/k rule

PhysMath does not contain physical roles such as:

```text
TIME
RG_SCALE
```

The strings `t` and `k` have no physical significance inside PhysMath.

If a framework declares:

```text
parameter X = physical time
parameter Y = renormalization-group scale
```

that declaration belongs in `phys-*`.

---

# 9. Mathematical constants

## 9.1 ConstantID

```go
type ConstantID string
```

Canonical encoding is:

```text
String(ConstantID)
```

Phase-1 built-ins are:

```text
math.const.pi
math.const.e
math.const.i
```

These values are closed in Phase 1.

## 9.2 Physical constants

A physical constant is represented initially as a mathematical symbol and acquires physical meaning only through a `phys-*` contract.

The same mathematical symbol may therefore have different physical interpretations under different frameworks.

---

# 10. Higher mathematical representations

## 10.1 Function

```text
Function {
    ParameterIDs []ParameterID // SEQUENCE
    BodyRef      ArtifactRef
}
```

Family:

```text
physmath/function/1
```

`BodyRef` MUST target `physmath/expression/1`.

Parameter order is semantic.

The Function-body validation rules in §8.6 apply before the Function artifact is accepted.

## 10.2 FunctionApplication

```text
FunctionApplication {
    FunctionRef ArtifactRef
    Args        []ArtifactRef // SEQUENCE
}
```

Family:

```text
physmath/function_application/1
```

`FunctionRef` MUST target `physmath/function/1`.

`Args[*]` MUST target `physmath/expression/1`.

The reference equivariance fixture uses inline `Call` expression nodes rather than durable FunctionApplication artifacts. The durable FunctionApplication family remains valid general representation but is not another fixture representation choice.

## 10.3 Map

```text
Map {
    DomainSort    string
    CodomainSort  string
    DefinitionRef ArtifactRef
}
```

Family:

```text
physmath/map/1
```

`DefinitionRef` MUST target `physmath/definition/1`.

`DomainSort` and `CodomainSort` MUST be registered Phase-1 PhysMath sort identifiers from §6.1.

A package MUST NOT introduce a new sort merely by putting an unknown string in one of these fields.

## 10.4 Tensor

```text
Tensor {
    Rank          UInt16
    Dimension     UInt16
    IndexRefs     []ArtifactRef  // SEQUENCE, slot order is semantic
    ComponentRefs []ArtifactRef  // SEQUENCE, canonical row-major order
}
```

Family:

```text
physmath/tensor/1
```

Target families:

```text
IndexRefs[*]     -> physmath/index/1
ComponentRefs[*] -> physmath/expression/1
```

For a dense tensor:

```text
len(IndexRefs)     = Rank
len(ComponentRefs) = Dimension^Rank
```

General sparse tensor semantics are deferred.

## 10.5 Index

```text
Index {
    Label     string
    Variance  IndexVariance
    SpaceRef  Optional<ArtifactRef>
}
```

Family:

```text
physmath/index/1
```

Closed v1 variance:

```text
COVARIANT
CONTRAVARIANT
UNSPECIFIED
```

`SpaceRef`, when present, MUST target `physmath/collection/1`.

Variance remains mathematical structure. It does not itself choose a physical metric or sign convention.

## 10.6 Collection

```text
Collection {
    ElementRefs []ArtifactRef // SET
}
```

Family:

```text
physmath/collection/1
```

In Phase 1:

```text
ElementRefs[*] -> physmath/expression/1
```

Duplicate canonical elements are invalid.

## 10.7 OperatorApplication

```text
OperatorApplication {
    OperatorID  string
    OperatorRef Optional<ArtifactRef>
    Args        []ArtifactRef // SEQUENCE
    Parameters  Bytes
}
```

Family:

```text
physmath/operator_application/1
```

`Args[*]` MUST target `physmath/expression/1`.

`OperatorRef` -> absent in Phase 1; non-absent value is invalid. The free-form `OperatorID` grammar in §5.4 applies, but the executable Phase-1 operator vocabulary is exactly `{differentiate}`. Any other OperatorID with an absent OperatorRef is invalid.

The Phase-1 reserved mathematical operator identity is:

```text
differentiate
```

For `differentiate`:

```text
OperatorRef = absent
Parameters =
    Encode(BoundParameterRef)
    || UInt32BE(Order)
```

`Order >= 1`.

## 10.8 Transformation

```text
Transformation {
    ProducerNamespace string
    RuleID            TransformationRuleID
    Inputs            []ArtifactRef // SEQUENCE
    Output            ArtifactRef
    Parameters        Bytes
}
```

Family:

```text
physmath/transformation/1
```

`ProducerNamespace` is part of canonical identity and MUST be encoded using `String()`.

`RuleID` is part of canonical identity and MUST be encoded using `String()`.

For Phase-1 PhysMath transformations:

```text
ProducerNamespace = "physmath"
Output family     = physmath/expression/1
```

`Inputs` are governed by the exact rule contract in §17. There is no generic open-ended input-family rule.

## 10.9 Transformation usage boundary

`Transformation` and `DerivationStep` intentionally share rule-execution fields but have different lifecycle roles.

`Transformation` is:

```text
standalone
durable
content-addressed
independently referenceable
```

It records one isolated mathematical operation.

`DerivationStep` is:

```text
inline
contained within a Derivation
not independently referenceable in Phase 1
```

It records one operation in an ordered proof/derivation trace.

A `Derivation` MUST NOT be transformed into a graph of tiny Transformation artifacts merely because the schemas share fields.

---

# 11. Mathematical support artifacts

These are durable artifacts, not additional Math sorts.

## 11.1 Equation

```text
Equation {
    RelationRef ArtifactRef
}
```

Family:

```text
physmath/equation/1
```

`RelationRef` MUST target `physmath/relation/1`.

An Equation packages a durable relation for reference.

## 11.2 EquationSet

```text
EquationSet {
    EquationRefs []ArtifactRef // SEQUENCE
}
```

Family:

```text
physmath/equation_set/1
```

`EquationRefs[*]` MUST target `physmath/equation/1`.

The sequence order is canonical and semantic.

Duplicate references are invalid in Phase 1 because an EquationSet represents a sequence of distinct equation artifacts. This is a deliberate uniqueness check on a SEQUENCE: order remains semantic and canonical, but duplicate canonical references are rejected. `EquationSet.EquationRefs` is never sorted as a SET.

## 11.3 Definition

```text
Definition {
    DefinedRef  ArtifactRef
    DefiningRef ArtifactRef
}
```

Family:

```text
physmath/definition/1
```

Phase-1 target families are:

```text
DefinedRef
    -> physmath/function/1 OR physmath/expression/1

DefiningRef
    -> physmath/expression/1
```

A Definition artifact does not itself become an executable transformation. Definitional application in expressions is performed only through `EXPR-APPLY-001` as specified in §17.

## 11.4 MathCondition

```text
MathCondition {
    PredicateRef ArtifactRef
}
```

Family:

```text
physmath/condition/1
```

`PredicateRef` MUST target `physmath/relation/1`.

MathCondition contains no physical interpretation or trust metadata.

## 11.5 TransformationRuleID

```go
type TransformationRuleID string
```

The exhaustive Phase-1 PhysMath rule vocabulary is exactly:

```text
EXPR-CANON-001
EXPR-SIMPLIFY-001
EXPR-SUBSTITUTE-001
RAT-NORMALIZE-001
DIFFERENTIATE-001
SELECT-BRANCH-001
EXPR-APPLY-001
```

`PROJECT-STRICT-001` is outside Phase 1.

No other `(ProducerNamespace, RuleID)` pair is executable by the Phase-1 PhysMath runtime.

---

# 12. Derivation representation

## 12.1 Purpose

A `Derivation` represents how a mathematical result was obtained from explicitly identified premises and operations.

It is not a universal proof assistant.

It records representation-level derivation structure sufficient for deterministic replay of the operations that the owning implementation supports.

## 12.2 DerivationStep

```text
DerivationStep {
    Ordinal           UInt32
    ProducerNamespace string
    RuleID            TransformationRuleID
    InputRefs         []ArtifactRef // SEQUENCE
    OutputRef         ArtifactRef
    ConditionRefs     []ArtifactRef // SET
    Parameters        Bytes
}
```

Family containment:

```text
DerivationStep is inline data inside physmath/derivation/1
```

It has no independent Phase-1 ArtifactRef.

Canonical field order is exactly the declaration order above.

Rules:

```text
Ordinal          = unique sequential position starting at 0
ProducerNamespace = rule producer namespace
RuleID            = rule identity within producer namespace
InputRefs         = SEQUENCE
ConditionRefs     = SET
Parameters        = exact canonical rule parameters
```

Each step may refer only to:

```text
an imported prerequisite artifact
or an output from an earlier step
```

unless the rule contract explicitly permits another referenced artifact.

## 12.3 Derivation

```text
Derivation {
    PremiseRefs []ArtifactRef // SET
    Steps       []DerivationStep // SEQUENCE
    ResultRefs  []ArtifactRef // SET
}
```

Family:

```text
physmath/derivation/1
```

Target families:

```text
PremiseRefs[*]
    -> one of
       physmath/expression/1
       physmath/function/1
       physmath/relation/1
       physmath/equation/1
       physmath/definition/1
       physmath/condition/1

ResultRefs[*]
    -> one of
       physmath/expression/1
       physmath/function/1
       physmath/relation/1
       physmath/equation/1
       physmath/definition/1
```

`Steps` preserve order.

`PremiseRefs` and `ResultRefs` are sets; duplicate canonical references are invalid.

## 12.4 Derivation validity

A derivation is structurally valid when:

```text
all references resolve
step ordinals are contiguous and unique
all step producer/rule pairs are registered for the applicable runtime
all step parameters satisfy their rule schema
all step InputRefs and OutputRef satisfy the selected rule's exact family contract in §17.2
all step ConditionRefs target physmath/condition/1
all same-run inputs are available before use
all generated/required conditions are explicit
all ResultRefs are produced by a step or are explicitly permitted carry-through artifacts under the derivation's declared premises
```

Structural validity MUST NOT be interpreted as scientific truth.

## 12.5 No hidden inference

A derivation implementation MUST NOT silently insert:

```text
unstated assumption
unstated equation
unstated theorem
physical interpretation
algebraic theorem not represented by a rule
```

If an operation requires a condition, the condition MUST be explicit in `ConditionRefs` or represented as a generated condition in the `TransformResult`; it MUST NOT be silently discharged.

## 12.6 Step inputs and same-run availability

The following artifacts are available to a step:

```text
all Derivation.PremiseRefs
all artifacts explicitly supplied as protocol fixtures
all OutputRefs from steps with lower Ordinal
```

An input reference to a future step output is invalid.

A step MUST NOT depend on another step's output merely because that output happens to be discoverable in a repository cache.

---

# 13. Theorem representation

## 13.1 Theorem

```text
Theorem {
    PremiseRefs    []ArtifactRef // SET
    ConclusionRef  ArtifactRef
    AssumptionRefs []ArtifactRef // SET
    DerivationRef  ArtifactRef
}
```

Family:

```text
physmath/theorem/1
```

Target families:

```text
PremiseRefs[*]
    -> one of
       physmath/expression/1
       physmath/function/1
       physmath/relation/1
       physmath/equation/1
       physmath/definition/1
       physmath/condition/1

ConclusionRef
    -> physmath/relation/1

AssumptionRefs[*]
    -> physmath/condition/1

DerivationRef
    -> physmath/derivation/1
```

## 13.2 Theorem meaning

`Theorem` means:

```text
this artifact represents a theorem claim
```

It does not mean:

```text
the theorem is true in nature
```

In Phase 1, `DerivationRef` is mandatory.

A theorem representation MUST reference a complete Derivation. Partial proofs and unfulfilled derivation obligations are not representable in Phase 1.

`DerivationObligation` remains deferred and has no Phase-1 durable family.

## 13.3 Theorem conclusion relation rule

Theorem conclusion is a durable `physmath/relation/1` artifact.

For an equality conclusion used by the reference protocol:

```text
ConclusionRef -> relation(EQ, Left, Right)
```

The relation's embedded left/right expression bytes MUST equal the canonical payload bytes of the corresponding final derivation result expression artifacts used as conclusion sides.

No hidden equality inference operation is introduced by theorem construction. The theorem-construction code is responsible for checking the fixture's explicit mapping from the first designated final result to the relation's left operand and from the second designated final result to the relation's right operand. This is construction-time validation, not a repository-wide inference pass.

## 13.4 Theorem versus physical theorem

PhysMath theorem:

```text
mathematical statement + mathematical derivation representation
```

`phys-*` physical theorem:

```text
physical assumptions + framework interpretation
+ mathematical theorem/equations
+ framework-specific context
```

The physical theorem remains owned by `phys-*`.

---

# 14. phys-* framework contract

## 14.1 Framework identity

Each `phys-*` package defines a framework artifact:

```text
FrameworkArtifact {
    FrameworkID        string
    PackageNamespace   string
    FrameworkVersion   string
    BaseFrameworkRefs  []ArtifactRef // SET
}
```

Family:

```text
<phys-namespace>/framework/1
```

Field order is exactly as declared above.

Validation:

```text
PackageNamespace = artifact namespace of the owning package
FrameworkID       = valid free-form identifier
FrameworkVersion  = non-empty opaque String
BaseFrameworkRefs = SET of framework refs
```

One framework version has exactly one framework identity artifact.

## 14.2 Physical assumptions

A framework defines package-owned assumptions:

```text
PhysicalAssumption {
    AssumptionID   string
    ConditionRefs  []ArtifactRef // SET
    Category       string
    Description    string
}
```

Family:

```text
<phys-namespace>/assumption/1
```

Target:

```text
ConditionRefs[*] -> physmath/condition/1
```

The architecture imposes no universal physical assumption vocabulary beyond the package contract.

`Category` is a framework-owned semantic label and intentionally remains an open string rather than a global ontology key.

`Description` is explanatory package-owned content and is identity-bearing because it is part of the canonical payload.

## 14.3 PhysicalDefinition

A framework may define physical meaning for a mathematical object:

```text
PhysicalDefinition {
    DefinitionID       string
    FrameworkRef       ArtifactRef
    PhysicalType       string
    MathematicalRefs   []ArtifactRef // SET
    AssumptionRefs     []ArtifactRef // SET
    DefinitionEquation Optional<ArtifactRef>
}
```

Family:

```text
<phys-namespace>/definition/1
```

Targets:

```text
FrameworkRef
    -> <same phys-namespace>/framework/1

MathematicalRefs[*]
    -> one of the following closed Phase-1 PhysMath families:
       physmath/expression/1
       physmath/function/1
       physmath/function_application/1
       physmath/map/1
       physmath/tensor/1
       physmath/index/1
       physmath/collection/1
       physmath/operator_application/1
       physmath/transformation/1
       physmath/relation/1
       physmath/equation/1
       physmath/equation_set/1
       physmath/definition/1
       physmath/condition/1
       physmath/derivation/1
       physmath/theorem/1

AssumptionRefs[*]
    -> <same framework namespace>/assumption/1
       OR an assumption family in the framework dependency closure

DefinitionEquation
    -> physmath/equation/1 when present
```

The framework decides the meaning of `PhysicalType`.

PhysMath does not interpret it.

### PhysicalDefinition versus PhysicalInterpretation

`PhysicalDefinition` states **what mathematical structure is taken to mean as a physical object in the framework**. `PhysicalInterpretation` states **the framework-specific interpretive context in which that meaning is assigned**.

They remain separate concerns and are not merged.

## 14.4 PhysicalEquation

```text
PhysicalEquation {
    EquationID     string
    FrameworkRef   ArtifactRef
    EquationRef    ArtifactRef
    AssumptionRefs []ArtifactRef // SET
    ScopeRefs      []ArtifactRef // SET
}
```

Family:

```text
<phys-namespace>/equation/1
```

Targets:

```text
FrameworkRef
    -> <same phys-namespace>/framework/1

EquationRef
    -> physmath/equation/1

AssumptionRefs[*]
    -> <same framework namespace>/assumption/1
       OR an assumption family in the framework dependency closure

ScopeRefs[*]
    -> <same framework namespace>/assumption/1
       OR an assumption family in the framework dependency closure
```

In Phase 1, `ScopeRefs` are intentionally limited to framework assumption artifacts because no separate generic scope ontology exists.

The mathematical equation remains a PhysMath artifact.

The physical equation adds framework meaning and assumptions.

## 14.5 PhysicalTheorem

```text
PhysicalTheorem {
    TheoremID               string
    FrameworkRef            ArtifactRef
    MathematicalTheoremRef  ArtifactRef
    AssumptionRefs          []ArtifactRef // SET
    InterpretationRefs      []ArtifactRef // SET
}
```

Family:

```text
<phys-namespace>/theorem/1
```

Targets:

```text
FrameworkRef
    -> <same phys-namespace>/framework/1

MathematicalTheoremRef
    -> physmath/theorem/1

AssumptionRefs[*]
    -> <same framework namespace>/assumption/1
       OR an assumption family in the framework dependency closure

InterpretationRefs[*]
    -> <same framework namespace>/interpretation/1
       OR an interpretation family in the framework dependency closure
```

The mathematical theorem and its mathematical derivation remain owned by PhysMath.

## 14.6 PhysicalInterpretation

A framework-specific interpretation is explicit:

```text
PhysicalInterpretation {
    InterpretationID string
    FrameworkRef     ArtifactRef
    PhysicalType     string
    MathematicalRefs []ArtifactRef // SET
    AssumptionRefs   []ArtifactRef // SET
}
```

Family:

```text
<phys-namespace>/interpretation/1
```

Targets:

```text
FrameworkRef
    -> <same phys-namespace>/framework/1

MathematicalRefs[*]
    -> one of the following closed Phase-1 PhysMath families:
       physmath/expression/1
       physmath/function/1
       physmath/function_application/1
       physmath/map/1
       physmath/tensor/1
       physmath/index/1
       physmath/collection/1
       physmath/operator_application/1
       physmath/transformation/1
       physmath/relation/1
       physmath/equation/1
       physmath/equation_set/1
       physmath/definition/1
       physmath/condition/1
       physmath/derivation/1
       physmath/theorem/1

AssumptionRefs[*]
    -> <same framework namespace>/assumption/1
       OR an assumption family in the framework dependency closure
```

There is no global interpretation registry.

## 14.7 Framework disagreement

Two packages MAY define:

```text
same mathematical expression
same mathematical equation
same theorem structure
```

with different physical meaning or assumptions.

This is conforming behavior, not an identity conflict.

The framework identity is part of the physical artifact's semantic context.

## 14.8 Base/specialized framework dependency rule

A base framework contains only the minimum assumptions required to define that framework.

A specialized framework MAY depend on the base framework:

```text
phys-qm
   ↑
phys-bohm
```

but not the reverse.

The same recursive principle applies to other framework families:

```text
base framework
    ↑
specialized framework
```

Examples include `phys-cdt`, `phys-cst`, `phys-gr`, and future domain packages.

Interpretation-specific assumptions MUST NOT be placed in the base framework merely because one specialization currently uses them.

## 14.9 Framework dependency closure

`FrameworkArtifact.BaseFrameworkRefs` expresses explicit direct framework dependencies.

Rules:

```text
1. A framework MUST NOT reference itself.
2. BaseFrameworkRefs MUST form an acyclic directed graph.
3. A specialized framework MAY reference physical assumptions or interpretations
   owned by itself or by any framework in the transitive BaseFrameworkRefs closure.
4. It MUST NOT reference unrelated framework namespaces.
5. It need not duplicate inherited base assumptions.
6. "Depends on" means both:
       - package-level dependency as required for implementation, and
       - explicit artifact-level dependency recorded by BaseFrameworkRefs.
```

For the artifact-reference rules in §§14.3–14.6, “dependency closure” means the transitive closure of `BaseFrameworkRefs` beginning at the artifact's owning framework.

This is the canonical solution to the base/specialized assumption contradiction: inherited assumptions remain owned by the base package and may be referenced by a dependent specialized package without duplication.

## 14.10 Interpretation scope

There is no Phase-1 `InterpretationClass` field.

Phase-1 ruling:

> **All Phase-1 physical artifacts are framework-scoped by construction through `FrameworkRef`; Phase 1 contains no interpretation-neutral physical claims.**

## 14.11 phys-* canonical serialization

All six `phys-*` durable families use the same canonical primitive encoding defined in §4.

Each artifact's field order is exactly its declaration order.

Each collection field uses the declared `SET` or `SEQUENCE` semantics.

No ad-hoc JSON serialization is canonical identity material.

For every registered concrete framework implementation, the package MUST provide:

```text
canonical serialization vector
expected ArtifactRef
negative vector
```

for each durable family it instantiates.

---

# 15. Theory-local algebra rule

## 15.1 When local algebra is permitted

A `phys-*` library MAY introduce a local mathematical type when:

```text
required by the physical framework
AND
not yet expressible using existing PhysMath types
```

## 15.2 Locality

The local type MUST remain inside the package:

```text
phys-foo.LocalExpr
```

It MUST NOT be added dynamically to PhysMath's tagged union or durable family registry.

## 15.3 Mapping to PhysMath

A package MAY define:

```text
LocalExpr
    -> explicit deterministic mapping
    -> PhysMath representation
```

The mapping MUST document:

```text
source type
target type
conditions for mapping
information preserved
information rejected
canonical output
```

No implicit fallback mapping is permitted.

## 15.4 Framework-local transformations

A future `phys-*` package MAY define framework-local executable transformations only under a future versioned execution contract.

Their rule identity MUST use the pair:

```text
(ProducerNamespace, RuleID)
```

and their producer namespace MUST be their conforming package namespace.

No framework-local transformation is executable by the Phase-1 runtime merely because a package defines a local Go function.

No dynamic PhysMath AST extension follows from framework-local rules.

---

# 16. Equivariance theorem reference slice

## 16.1 Why this is the Phase-1 target

The reference implementation target is the representation of an equivariance theorem derivation because it requires the architecture to exercise:

```text
mathematical objects
parameters
functions
function application
relations/equations
conditions
assumptions
transformations
derivation sequencing
theorem representation
physical framework interpretation
executable protocol packaging
```

It is therefore a vertical-slice test of the architecture rather than a special feature of any one physical theory.

## 16.2 Canonical protocol packaging

The equivariance reference fixture MUST be packaged exactly in the project-standard form:

```text
<equivariance-fixture>/
    AGENTS.md
    manifest.json
    protocol/
        derivation_test.go
```

The fixture MUST follow the same structural and operational pattern used by the existing General Relativity example.

No alternative “equivariance fixture format” is introduced.

## 16.3 Generic equivariance structure

The representation MUST be capable of expressing:

```text
Given:
    a group/action parameter g
    a source action rho_X
    a target action rho_Y
    a map f : X -> Y

Equivariance:

    f(rho_X(g, x)) = rho_Y(g, f(x))
```

The abstract concepts `G`, `X`, and `Y` are semantic labels of the fixture rather than new PhysMath sorts.

No `GroupAction` sort is created.

## 16.4 Canonical representation choices

The reference fixture MUST use:

```text
f, rho_X, rho_Y
    -> physmath/function/1

LHS / RHS
    -> physmath/expression/1

Equivariance relation
    -> physmath/relation/1

Equation wrapper
    -> physmath/equation/1

Mathematical theorem
    -> physmath/theorem/1

Physical framework artifacts
    -> concrete phys-* families
```

The reference fixture MUST NOT use `Map` as an alternate representation of `f` and MUST NOT use durable `FunctionApplication` artifacts for the LHS/RHS when inline `Call` nodes are sufficient.

This removes the structural ambiguity identified by the v2 review.

## 16.5 Fixture parameters `g` and `x`

Phase 1 does not introduce a general top-level quantified variable/binder sort.

The reference fixture represents `g` and `x` as schematic ParameterID-bearing Symbols using the two declaration slots of the action Functions:

```text
rho_X.ParameterIDs = [PID_0, PID_1]
rho_Y.ParameterIDs = [PID_0, PID_1]

g = Symbol(ParameterID = PID_0, DisplayName = "g")
x = Symbol(ParameterID = PID_1, DisplayName = "x")
```

Here:

```text
PID_0 = ParameterID for Function parameter ordinal 0
PID_1 = ParameterID for Function parameter ordinal 1
```

The same positional values occur in both two-parameter action Functions by design.

The map `f` uses a separate one-parameter declaration slot for its local parameter:

```text
f.ParameterIDs = [PID_0]
f(v) uses Symbol(ParameterID = PID_0, DisplayName = "v")
```

That local PID is capture-safe because `EXPR-APPLY-001` substitutes the entire function body parameter slot with the supplied argument expression before the caller's symbols are retained.

The reference fixture does not introduce a general global variable identity mechanism. `g` and `x` are schematic theorem variables whose ParameterIDs are interpreted in the fixture's explicit action-function declaration context.

## 16.6 X and Y

`X` and `Y` are represented by framework-owned definitions associated with ordinary registered PhysMath structures rather than by new PhysMath domain sorts.

The fixture does not require a `Group`, `GroupAction`, or `Domain` sort.

The mathematical action/functions are sufficient for the Phase-1 representation target.

## 16.7 Concrete mathematical realization

The reference fixture MUST use the following concrete functions because they close under the Phase-1 rule set without requiring distributivity or any unregistered theorem:

```text
f(v)        = 2 * v
rho_X(g,v)  = g * v
rho_Y(g,v)  = g * v
```

Therefore:

```text
f(rho_X(g,x))
    = 2 * (g * x)
    -> 2 * g * x

rho_Y(g,f(x))
    = g * (2 * x)
    -> 2 * g * x
```

The final equality is achieved by canonicalization/simplification only.

The reference fixture MUST NOT require distributivity, factorization, `x+x -> 2x`, trigonometric identities, calculus identities, or any other unregistered theorem.

## 16.8 Exact expression structure

The canonical LHS is:

```text
Call(
    f,
    [
        Call(
            rho_X,
            [g, x]
        )
    ]
)
```

The canonical RHS is:

```text
Call(
    rho_Y,
    [
        g,
        Call(
            f,
            [x]
        )
    ]
)
```

The durable relation is:

```text
EQ(LHS, RHS)
```

The Equation wraps that durable relation.

The Theorem references the Equation's relation as its ConclusionRef and references the complete Derivation.

## 16.9 Normative derivation sequence

The reference derivation MUST use the following exact transformation sequence.

Fixture construction of Calls and function artifacts is setup data and is not itself a DerivationStep. Because Phase 1 contains no generic enclosing-AST rebuild rule, the fixture MUST pre-construct these staged Calls:

```text
LHS_0 = Call(f, [Call(rho_X, [g, x])])
LHS_1 = Call(f, [g * x])
RHS_0 = Call(rho_Y, [g, Call(f, [x])])
RHS_1 = Call(rho_Y, [g, 2 * x])
```

Before replay, the protocol MUST verify:

```text
canonical bytes of LHS_1.Call.Args[0] == expected Step 0 output bytes
canonical bytes of RHS_1.Call.Args[1] == expected Step 3 output bytes
```

The staged Calls are explicit fixture data, not hidden transformations. The linkage checks are construction-time fixture validation. No generic AST subtree-replacement rule is implied.

The six derivation transformations are:

```text
Step 0
    Rule = (physmath, EXPR-APPLY-001)
    Input = inner rho_X(g,x) Call expression
    Function = rho_X
    Output = g * x

Step 1
    Rule = (physmath, EXPR-APPLY-001)
    Input = staged LHS_1 = f(g * x) Call expression
    Function = f
    Fixture linkage = LHS_1 argument bytes equal Step 0 output bytes
    Output = 2 * (g * x)

Step 2
    Rule = (physmath, EXPR-CANON-001)
    Input = Step 1 output
    Output = 2 * g * x

Step 3
    Rule = (physmath, EXPR-APPLY-001)
    Input = inner f(x) Call expression
    Function = f
    Output = 2 * x

Step 4
    Rule = (physmath, EXPR-APPLY-001)
    Input = staged RHS_1 = rho_Y(g, 2 * x) Call expression
    Function = rho_Y
    Fixture linkage = RHS_1 second-argument bytes equal Step 3 output bytes
    Output = g * (2 * x)

Step 5
    Rule = (physmath, EXPR-CANON-001)
    Input = Step 4 output
    Output = 2 * g * x
```

The fixture MUST pin exact ArtifactRefs and expected output ArtifactRefs for every step. Those golden values are fixture data produced by the canonical implementation after the logical fixture is fixed by this specification; they are not additional architectural choices. A conforming independent implementation MUST reproduce them.

The final two expression outputs MUST have identical canonical bytes.

The fixture MUST pin the following Symbol display names exactly:

```text
f parameter/body slot:        "v"
rho_X parameters:             "g", "v"
rho_Y parameters:             "g", "v"
standalone schematic symbols: "g", "x"
```

These display names are part of the canonical Symbol payload, so changing them changes fixture ArtifactRefs and requires a fixture-version change.

## 16.10 Relation-side identity

The durable `physmath/relation/1` payload is constructed exactly as:

```text
Enum(EQ)
|| CanonicalExpression(Step 2 output)
|| CanonicalExpression(Step 5 output)
```

A standalone `physmath/expression/1` artifact representing either relation side MUST have a canonical payload byte-for-byte identical to the embedded expression copy.

No alternate reference-based relation encoding is permitted.

The expression bytes embedded in the relation are authoritative at construction time. The fixture/tooling MUST obtain them from the designated final expression artifacts and verify byte identity before constructing the relation artifact; the relation validator does not perform a repository-wide search to infer correspondence.

## 16.11 No hidden equality theorem

The reference fixture establishes equivariance through deterministic expression identity after explicit transformations.

It MUST NOT rely on an unregistered algebraic equality rule.

Relation packaging is deterministic artifact construction, not a theorem-proving transformation.

The theorem conclusion is accepted only when:

```text
left expression bytes == right expression bytes
AND
ConclusionRef relation embeds those exact bytes
AND
those side artifacts are produced by the declared derivation
```

## 16.12 Physical fixture requirements

Every concrete `phys-*` reference fixture MUST add at least:

```text
one FrameworkArtifact
at least one PhysicalAssumption
at least one PhysicalDefinition or PhysicalInterpretation
one PhysicalTheorem referencing the mathematical theorem
```

The fixture MUST NOT modify PhysMath to encode the physical interpretation.

## 16.13 Many-framework ambiguity test

The implementation suite MUST include two physical framework fixtures that interpret the same or overlapping mathematical representation differently.

Expected result:

```text
PhysMath identity remains unchanged
framework artifacts remain distinct
physical assumptions remain distinct
physical interpretations remain distinct
no global winner is inferred
```

This is the concrete implementation test for the many-frameworks principle.

---

# 17. Operation semantics

## 17.1 Common operation contract

Every Phase-1 operation has the abstract execution shape:

```go
type OperationInvocation struct {
    ProducerNamespace string
    RuleID            TransformationRuleID
    InputRefs         []ArtifactRef // SEQUENCE
    ConditionRefs     []ArtifactRef // SET
    Parameters        []byte
}

Apply(invocation OperationInvocation) (TransformResult, error)
```

`OperationInvocation` is in-memory only. It is the explicit execution context for one rule invocation. `ConditionRefs` are context inputs and do NOT count toward rule input arity. The supplied condition references MUST be duplicate-free and normalized according to the canonical SET ordering rule before evaluation; their ordering therefore has no semantic effect.

```go
type TransformResult struct {
    OutputRef     ArtifactRef
    ConditionRefs []ArtifactRef // SET semantics
}
```

`TransformResult` is in-memory only and immutable after construction. Its `ConditionRefs` are the canonical union of supplied conditions and any conditions explicitly generated by the rule. Phase-1 rules generate no new mathematical conditions; therefore, for every Phase-1 rule, `TransformResult.ConditionRefs` is exactly the supplied `OperationInvocation.ConditionRefs` after canonical SET normalization. No Phase-1 rule may silently discharge a condition.

There is no `OperationResult` type and no `TraceRef` field in `TransformResult`.

On operation failure, the implementation MUST return a non-nil error and MUST NOT return a usable output artifact reference.

For a `DerivationStep`, the executor MUST construct `OperationInvocation.ConditionRefs` directly from that step's `ConditionRefs`.

## 17.2 Rule contract table

Every executable Phase-1 rule has a fixed input arity, positional input family, root constraint where applicable, output family, parameter schema, and invalid-input behavior.

| Rule | Exact arity | Positional input contract | Output | Invalid input |
|---|---:|---|---|---|
| `EXPR-CANON-001` | 1 | Input 0 = `physmath/expression/1` | `physmath/expression/1` | fail; no output |
| `EXPR-SIMPLIFY-001` | 1 | Input 0 = `physmath/expression/1` | `physmath/expression/1` | fail; no output |
| `EXPR-SUBSTITUTE-001` | `1 + D` | Input 0 = expression; Inputs 1..D = distinct declaring `physmath/function/1` refs required by substitution targets | `physmath/expression/1` | fail; no output |
| `RAT-NORMALIZE-001` | 1 | Input 0 = `physmath/expression/1`, root `Rational` | `physmath/expression/1` | fail; no output |
| `DIFFERENTIATE-001` | 2 | Input 0 = expression; Input 1 = declaring `physmath/function/1` | `physmath/expression/1` | fail; no output |
| `SELECT-BRANCH-001` | 1 | Input 0 = `physmath/expression/1`, root `BranchSet` | `physmath/expression/1` | fail; no output |
| `EXPR-APPLY-001` | 2 | Input 0 = `physmath/expression/1`, root `Call`; Input 1 = referenced `physmath/function/1` | `physmath/expression/1` | fail; no output |

For `EXPR-SUBSTITUTE-001`, `D >= 1`.

For all seven Phase-1 rules, an arity mismatch is invalid even when the supplied artifacts could otherwise be interpreted meaningfully. `ConditionRefs` are execution-context inputs and do not change rule arity. None of the seven Phase-1 rules generates a new mathematical condition; condition-dependent success is therefore determined entirely by the explicit supplied condition context.

## 17.3 EXPR-CANON-001

Parameters:

```text
empty byte sequence
```

Traversal and rewrite policy:

```text
1. Traverse expression nodes recursively in post-order.
2. Canonicalize every child before canonicalizing the parent.
3. Flatten nested Add and Mul nodes.
4. Sort commutative children by lexicographic canonical-expression bytes.
5. Normalize Add/Mul unary and empty arity according to §7.6.
6. Preserve explicit node kinds other than the closed canonicalization rules.
7. No global second pass is required after completion.
8. The operation is idempotent.
```

`EXPR-CANON-001` does not:

```text
combine x+x into 2x
apply distributivity
factor expressions
prove arbitrary equalities
infer physical meaning
```

## 17.4 EXPR-SIMPLIFY-001

Parameters:

```text
empty byte sequence
```

Traversal policy:

```text
1. Traverse recursively in post-order.
2. Simplify children before the parent.
3. At each node, repeatedly apply the closed rewrite list until that node is locally stable.
4. Rewrite priority is the order in the closed list below.
5. After each rewrite, restart rule selection at the top of the priority list for that node.
6. Any newly formed Add or Mul is kept in canonical child order.
7. The operation is deterministic and idempotent.
```

Closed rewrite list:

```text
Neg(Neg(x)) -> x

Add(..., Rational(0/1), ...) -> Add(...)
    when removing the zero preserves valid Add arity; unary/empty cases then use §7.6 normalization

Mul(..., Rational(1/1), ...) -> Mul(...)
    when removing the one preserves valid Mul arity; unary/empty cases then use §7.6 normalization

Mul(..., Rational(0/1), ...) -> Rational(0/1)

Rational arithmetic is permitted only when the current node is `Add` or `Mul` and every child at that node is a Rational expression node. The entire current node is evaluated as exact rational arithmetic and replaced by one reduced Rational node. Partial folding, adjacent-pair folding, or reduction of a subset of Rational children in a mixed symbolic Add/Mul node is NOT permitted.

Rational arithmetic on exact Rational children
    -> exact reduced Rational result

Pow(x, 1) -> x

Pow(x, 0) -> Rational(1/1)
    only when an explicit MathCondition x != 0 is present in the derivation context
```

The condition rule is fail-closed:

```text
without explicit x != 0 condition -> transformation fails
```

No condition is silently assumed or discharged.

The explicit condition MUST be supplied through `OperationInvocation.ConditionRefs` and, when used in a derivation, therefore appears in `DerivationStep.ConditionRefs`. The condition is satisfied only by an exact canonical `MathCondition` predicate of the form:

```text
NEQ(canonical-x, Rational(0/1))
```

where `canonical-x` is byte-identical to the operand of the `Pow(x,0)` node. No logical equivalence checking, theorem proving, or condition synthesis occurs.

A missing condition causes ordinary transformation failure with no output. The specification does not require a particular Go error type or error string.

`EXPR-SIMPLIFY-001` does not add:

```text
distributivity
factorization
x+x -> 2x
trigonometric identities
calculus identities
theorem proving
```

## 17.5 RAT-NORMALIZE-001

Parameters:

```text
empty byte sequence
```

Input:

```text
exactly one physmath/expression/1 artifact
root node MUST be Rational
```

Output is one `physmath/expression/1` artifact containing the normalized Rational node.

Normalization guarantees:

```text
denominator > 0
gcd(abs(numerator), denominator) = 1
zero -> 0/1
```

The operation is idempotent.

## 17.6 SELECT-BRANCH-001

Exactly one input is required:

```text
physmath/expression/1
```

The root node MUST be `BranchSet`.

Parameters are:

```text
UInt32BE(BranchOrdinal)
```

The ordinal is zero-based and MUST be in range.

Output is exactly the selected branch as a new `physmath/expression/1` artifact.

No condition is generated by this rule.

Branch selection is an explicit structural choice, not a mathematical predicate.

## 17.7 EXPR-APPLY-001

`EXPR-APPLY-001` is the only Phase-1 definitional/function-body expansion rule.

Inputs are exactly:

```text
Input[0] = physmath/expression/1
Input[1] = physmath/function/1
```

Applicability predicate, evaluated before rewriting:

```text
Inputs length = 2
Input[0] root node is Call
Input[1] is a Function artifact
Input[0].Call.FunctionRef == ArtifactRef(Input[1])
len(Input[1].ParameterIDs) == len(Input[0].Call.Args)
Function body passes the validation rules of §8.6
```

If any predicate fails, the transformation MUST fail with no output.

There is no NodePath parameter.

The root node MUST be the Call. A Call nested inside `Add`, `Mul`, `Pow`, `Relation`, or another node is not selected by the same invocation.

For a nested Call, the caller MUST isolate that Call subtree as a top-level `physmath/expression/1` artifact, apply `EXPR-APPLY-001`, and rebuild the enclosing expression explicitly. Phase 1 provides no generic NodePath selector and no generic AST subtree-replacement rule. The equivariance fixture therefore uses pre-constructed staged outer Calls and explicit byte-linkage validation as defined in §16.9.

The Function body is expanded by simultaneous substitution of its parameter slots with the Call arguments.

The same invocation does not recursively expand Calls that appear inside the selected Function body's result after substitution.

Function body symbols that represent constants remain constants.

Function body parameter symbols are substituted by ParameterID according to the Function's declared parameter order.

No alpha-freshening or renaming is permitted.

Parameters:

```text
empty byte sequence
```

## 17.8 EXPR-SUBSTITUTE-001

Substitution is simultaneous and recursive over the complete expression tree.

Input arrangement is exactly:

```text
Input[0]      = target physmath/expression/1
Input[1..D]   = distinct declaring physmath/function/1 artifacts
```

`D` is the number of distinct declaring Function refs used by the substitution targets.

The declaring Function refs in Inputs 1..D MUST be sorted by canonical ArtifactRef bytes.

The parameters byte payload is:

```text
UInt32BE(N)
|| Entry_0
|| ...
|| Entry_(N-1)
```

where entries are sorted lexicographically by their exact canonical entry bytes.

Each entry is:

```text
Encode(BoundParameterRef)
|| EncodeArtifactRef(ReplacementExpression)
```

Every ReplacementExpression MUST target `physmath/expression/1`.

Validation:

```text
N >= 1
all Target BoundParameterRefs have declaring artifacts in Inputs[1..D]
each target ParameterID is declared by its referenced Function
duplicate targets are invalid. Specifically, no two substitution entries may share the same `ParameterID`, regardless of `DeclaringArtifactRef`, because expression matching is by `ParameterID`. A ParameterID collision is therefore sufficient to make the substitution invalid.
matching inside Input[0] is by ParameterID
matching recurses through every expression child node
replacement expressions are not recursively substituted into again during the same invocation
```

The transformation MUST fail if a target ParameterID is not declared by its referenced Function.

The transformation MUST fail if the declaring Function context is missing from the inputs. Because matching is by `ParameterID`, the full set of substitution target ParameterIDs for one invocation MUST be pairwise distinct; no symbol can be assigned two replacements. The input expression is interpreted under its explicit schematic parameter context; the rule MUST NOT infer an additional declaring artifact from `DisplayName`, repository lookup, or the replacement expression.

## 17.9 DIFFERENTIATE-001

Exactly two inputs are required:

```text
Input[0] = physmath/expression/1
Input[1] = physmath/function/1
```

The parameters are:

```text
Encode(BoundParameterRef)
|| UInt32BE(Order)
```

Validation:

```text
BoundParameterRef.DeclaringArtifactRef == ArtifactRef(Input[1])
Target ParameterID is declared by Input[1]
Order >= 1
```

Matching within the target expression is by ParameterID.

Scope:

```text
single target ParameterID
scalar expressions
ordinary partial differentiation with respect to that target
```

Deferred:

```text
multi-target derivative operators
covariant differentiation
tensor differentiation
framework-specific physical derivative semantics
```

Closed derivative rules:

```text
d(Rational)/dp = 0

d(Symbol(target))/dp = 1

d(Symbol(other parameter))/dp = 0

d(Add(xs))/dp = Add(d(x1), ..., d(xn))
d(Mul(xs))/dp = sum_i(product_except_i(x) * d(x_i))
d(Neg(x))/dp = Neg(d(x))
d(Pow(b,n))/dp = n * Pow(b,n-1) * d(b)
```

For the `Pow` rule, the exponent `n` MUST be a nonnegative integer Rational and all required side conditions MUST be explicit.

`Call`, `Relation`, and `BranchSet` differentiation are out of scope for the Phase-1 core runtime.

Repeated differentiation is represented by `Order > 1` in one canonical operator application rather than by hidden chaining.

## 17.10 Condition handling

If an operation requires a mathematical condition:

```text
operation execution
    -> explicit MathCondition reference(s)
    -> ConditionRefs in TransformResult and/or DerivationStep context
```

Conditions MUST be explicit.

The implementation MUST NOT silently discharge them.

A generated condition is an obligation attached to the transformation result; it is not proof of the condition.

## 17.11 Rule identity and namespace

Durable rule identity is the pair:

```text
(ProducerNamespace, RuleID)
```

For all Phase-1 core PhysMath rules:

```text
ProducerNamespace = "physmath"
```

The durable `Transformation` and `DerivationStep` schemas contain the ProducerNamespace field explicitly.

Unknown rule pairs MUST be rejected during validation/deserialization.

For the Phase-1 PhysMath runtime, only `(physmath, <one of the seven closed Phase-1 RuleIDs>)` is executable. Framework-local rule pairs may exist as data only for future versioned contracts; the Phase-1 PhysMath validator MUST reject them for execution rather than guessing a delegation target.

No framework-local rule becomes executable in Phase 1 without a future versioned extension.

---

# 18. Canonicalization and deterministic behavior

## 18.1 Canonical child order

For commutative mathematical collections such as `Add` and `Mul`:

```text
sort by lexicographic canonical encoded child bytes
```

No runtime map iteration order may affect identity.

## 18.2 Sequence semantics

A `SEQUENCE` is used when order itself is meaningful.

Examples:

```text
function parameters
function arguments
tensor index slots
tensor components
derivation steps
branch alternatives
operation inputs
EquationSet equations
```

`EquationSet.EquationRefs` is exclusively `SEQUENCE` and MUST NOT appear in any SET-oriented validation or sorting path.

## 18.3 Set semantics

A `SET` is used when order carries no mathematical or framework meaning.

Examples:

```text
assumption references
premises
physical interpretation references
framework dependency references
```

The schema MUST explicitly state `SET` or `SEQUENCE` for every collection field.

## 18.4 Canonicalization is not theorem proving

Canonicalization may make structurally equivalent mathematical representations identical.

It MUST NOT:

```text
solve arbitrary equations
prove theorems
infer physical meaning
invent assumptions
apply unregistered identities
```

## 18.5 Canonical operation idempotence

The following are idempotent:

```text
EXPR-CANON-001
EXPR-SIMPLIFY-001
RAT-NORMALIZE-001
```

Idempotence means applying the same operation with the same applicable context to its own output produces the same canonical artifact reference.

---

# 19. Equality rules

## 19.1 Structural equality

Canonical mathematical equality means:

```text
canonical payload bytes are identical
```

for the same mathematical artifact family and schema.

## 19.2 Relation equality

For `physmath/relation/1`, equality is canonical payload equality of:

```text
Enum(Operator)
|| CanonicalExpression(Left)
|| CanonicalExpression(Right)
```

The relation's left/right operands are embedded value data, not reference indirection.

## 19.3 Physical equivalence

Physical equivalence is framework-specific.

PhysMath MUST NOT define:

```text
GR equivalence
QM interpretive equivalence
Bohmian equivalence
```

## 19.4 Theorem equality

Two theorem artifacts are equal only when their canonical theorem payloads are equal.

Two mathematically equivalent theorems may nevertheless be different artifacts if they use different canonical theorem structures; equivalence reasoning is not automatically identity normalization.

## 19.5 Cross-framework identity

A PhysMath artifact may be referenced by multiple physical artifacts in different frameworks.

This does not alter the PhysMath ArtifactRef.

Physical artifact identity remains framework-scoped by its own namespace and payload.

---

# 20. Serialization and replay

## 20.1 Canonical serialization

Every durable Phase-1 artifact MUST have:

```text
fixed field order
explicit lengths
explicit optional presence bytes
explicit sequence/set semantics
exact integer widths
UTF-8 for strings
closed enum values
```

## 20.2 Transport separation

JSON MAY be used for debugging or transport.

Transport JSON MUST NOT be used as canonical identity material unless a future schema explicitly declares it canonical.

## 20.3 Replay authority

A mathematical transformation is replayed by the authority identified by `ProducerNamespace`.

For Phase-1 core rules:

```text
ProducerNamespace = physmath
replay authority = PhysMath
```

A future framework-local rule would be replayed only by its owning namespace under a future versioned contract.

A consumer MUST NOT reinterpret a RuleID under a different ProducerNamespace.

No Phase-1 PhysMath or `phys-*` artifact may reference the frozen `phys` namespace. Frozen-Phys projection begins only under a future projection specification.

## 20.4 Replay prerequisites

To replay a `DerivationStep`, an implementation MUST have:

```text
rule producer namespace
rule ID
rule parameters
input artifacts
applicable conditions
```

in canonical form.

Replay MUST NOT depend on:

```text
runtime pointer addresses
Go map iteration order
repository discovery order
human-readable display names
implicit physical conventions
```

## 20.5 Golden vectors

Every durable family and every executable rule contract MUST have at least one positive canonical vector and one negative vector in the reference implementation.

Golden vectors include:

```text
ArtifactRef
ParameterID
BoundParameterRef
physmath/expression/1
physmath/relation/1
physmath/transformation/1
physmath/derivation/1
physmath/theorem/1
all six phys-* families
all seven Phase-1 rule parameter encodings
```

The equivariance fixture additionally pins the exact step sequence and expected ArtifactRefs.

---

# 21. phys-* conformance contract

A `phys-*` package implementing this specification MUST demonstrate:

```text
namespace ownership
framework identity
explicit BaseFrameworkRefs
concrete physical artifact schemas
ArtifactRef compatibility
canonical serialization
explicit assumptions
explicit mathematical references
explicit framework interpretation
no dynamic PhysMath AST extension
no duplicate Phys trust/corpus authority
```

## 21.1 Required negative tests

At minimum:

```text
physical meaning inferred from a PhysMath display name -> reject
unregistered physical artifact type -> reject
unregistered RuleID pair -> reject
missing required mathematical reference -> reject
duplicate SET element -> reject
invalid ParameterID binding -> reject
implicit condition discharge -> reject
dynamic PhysMath AST extension attempt -> reject
self-referential framework dependency -> reject
cyclic framework dependency -> reject
unrelated framework assumption reference -> reject
```

## 21.2 Framework disagreement test

Two framework packages MUST be able to attach different physical definitions or interpretations to the same PhysMath representation without mutating the PhysMath artifact.

## 21.3 Base/specialized assumption test

At minimum:

```text
phys-qm assumption referenced by phys-qm artifact      -> valid
phys-bohm base phys-qm assumption referenced           -> valid
phys-bohm private assumption referenced                -> valid
phys-qm referencing phys-bohm assumption               -> invalid
unrelated phys-cst assumption referenced by phys-bohm  -> invalid
cyclic framework dependency                            -> invalid
```

## 21.4 Protocol conformance

A formalization package claiming a theorem, hypothesis, derivation, or physics-paper formalization MUST contain:

```text
AGENTS.md
manifest.json
protocol/derivation_test.go
```

The executable protocol MUST pass its declared derivation test.

A package with only static prose is not a Phase-1 formalization artifact set.

---

# 22. Explicit Phase-1 no-go list

The Phase-1 implementation MUST NOT introduce:

```text
TOE-specific ontology into PhysMath
Dancer/Dance/Stage as PhysMath types
physical-role enums such as TIME or RG_SCALE in PhysMath
Numerical execution into PhysMath
Numerical dependencies into the core PhysMath package
research state machines into PhysMath
paper/citation databases into PhysMath
trust/corpus/provenance ontology into PhysMath
universal Artifact interfaces
global semantic artifact registry
dynamic PhysMath AST extension
physical interpretation inferred from symbol names
scientific truth flags attached to mathematical artifacts
automatic theorem proving hidden in canonicalization
automatic assumption generation without explicit rule contracts
frozen-Phys Kind/Dimension/AssumptionSet projection in Phase 1
TOE-Agent artifact schemas in Phase 1
Generic PhysicalContract ontology in Phase 1
generic node-path transformation selectors in Phase 1
```

---

# 23. Secondary-system boundary

## 23.1 TOE-Agent

TOE-Agent may later consume:

```text
PhysMath artifacts
phys-* artifacts
paper references
hypotheses
correspondences
research states
```

None of those research workflow concerns belongs in this Phase-1 implementation contract.

The later research platform may use autoformalization, semantic retrieval, candidate generation, verifier feedback, multi-interpretation search, and other orchestration mechanisms without changing the Phase-1 representation ontology.

## 23.2 Numerical

Numerical later consumes explicit executable numerical contracts produced by an upstream lowering process.

Numerical is not part of PhysMath's representation core and MUST NOT parse arbitrary PhysMath ASTs as an implicit execution layer in Phase 1.

## 23.3 Provenance and trust

Phase 1 preserves only the minimum identity/reference information necessary to locate the represented object.

It does not define:

```text
scientific trust
corpus promotion
established/hypothesis status
human approval
provenance policy
```

Such metadata can be layered later without changing the mathematical identity of an artifact.

## 23.4 Future autoformalization layer

A future higher layer may use the Phase-1 representation substrate as its target for:

```text
natural-language analysis
semantic retrieval
candidate generation
search-guided formalization
verifier-in-the-loop repair
adversarial counterexample generation
semantic-fidelity evaluation
paper-scale formalization
```

Those functions remain deliberately outside Phase 1.

The core boundary is:

```text
AI / human analysis
        -> candidate formal meaning
        -> PhysMath / phys-* representation
        -> executable protocol
        -> deterministic replay / validation
```

The future layer may be generalised beyond physics, but it MUST NOT require PhysMath itself to understand chemistry, biology, economics, or another scientific domain.

---

# 24. Implementation contract for the coding agent

## 24.0 Plan-first implementation handoff

The first coding-agent deliverable after this specification is a **project-specific implementation plan**, not source code. The plan is prepared against:

```text
PM-REP-001 v1.4
frozen Phys authority where applicable
existing General Relativity example as protocol precedent
existing phys-gr implementation as refactoring substrate
```

The implementation plan MUST:

```text
1. inventory the current phys-gr code and responsibilities;
2. map existing code to phys-artifact, PhysMath, phys-gr-specific semantics, and protocol roles;
3. identify the common denominator that belongs in the core;
4. identify differences that must remain in phys-gr or other phys-* packages;
5. identify refactors/migrations needed to satisfy this specification;
6. identify the tests needed to preserve valid behavior while changing boundaries;
7. identify the exact path to the executable equivariance protocol fixture;
8. identify any architectural blocker or missing decision discovered during planning.
```

The coding agent MAY ask the human architect for input only when an issue satisfies §24.4's blocker test: contradiction, missing representation decision, ownership conflict, or materially different architectural choices that remain legitimate.

When requesting architectural input, the agent MUST provide:

```text
problem statement
exact specification section(s) involved
exact existing-code or protocol evidence
competing architectural interpretations
concrete consequence of each interpretation
recommended resolution when defensible
minimal specification amendment required
```

The agent MUST NOT ask the human architect to choose permitted implementation details such as Go package layout, helper naming, local storage, caching, or CI mechanics.

Planning may continue through non-blocking details. Once no architectural blocker remains, the plan becomes the basis for implementation. The coding agent MUST NOT silently edit this specification to resolve a newly discovered architectural issue; it must report the issue under §24.4 and wait for an explicitly resolved architectural decision before implementation of the affected behavior.

## 24.1 Primary objective

Implement:

```text
existing phys-gr implementation
    -> refactor/extract
phys-artifact
    +
PhysMath
    +
phys-gr as the first concrete phys-* framework
    -> complete executable equivariance-theorem representation
```

`phys-gr` is the first real-world implementation used to validate the core/specific split. Its GR-specific semantics MUST remain outside PhysMath unless they are demonstrated to be the common denominator required by multiple domains/frameworks.

## 24.2 Architectural decisions already fixed

The coding agent MUST NOT decide any of the following:

```text
whether PhysMath knows physics
whether phys-* may define physical meaning
whether multiple frameworks may coexist
what ArtifactRef means
what canonical digest framing is
what primitive encoding means
whether symbol names define parameter identity
whether physical roles belong in PhysMath
whether physical assumptions belong in phys-*
whether local framework algebra may dynamically extend PhysMath
whether TOE-Agent is required
whether Numerical is required
whether a PhysicalContract ontology is required
whether a new GroupAction sort is needed
whether the equivariance fixture may use an alternate artifact graph
whether the protocol may omit AGENTS.md, manifest.json, or derivation_test.go
whether the existing GR example protocol may be replaced by a new package protocol
how rule conditions are supplied to the runtime
whether a global phys-* namespace registry is required
whether implementation begins with a coding-agent plan
whether existing phys-gr code is the refactoring substrate
```

These are fixed by this specification.

## 24.3 Permitted implementation decisions

The implementation may decide only details that do not alter semantic contracts:

```text
Go package/file structure
repository layout
constructor/helper naming where schema semantics are unchanged
error type naming where semantics are unchanged
test fixture organization
CI configuration
local cache/storage implementation
codec implementation strategy
benchmark harnesses
```

## 24.4 Stop rule

When implementation encounters behavior not defined here, the coding agent MUST NOT invent a semantic rule.

It must classify the gap as:

```text
MATH GAP        -> PhysMath specification issue
PHYS GAP        -> phys-* contract issue
FROZEN-PHYS GAP -> inspect frozen Phys authority; do not modify Phys
INTEGRATION GAP -> defer to secondary-system specification
PROTOCOL GAP    -> compare against the GR example protocol; do not invent a second protocol
```

Only a demonstrated contradiction or missing representation decision may reopen this specification.

## 24.5 Phase-1 execution lock

The only runtime transformation pairs in Phase 1 are:

```text
(physmath, EXPR-CANON-001)
(physmath, EXPR-SIMPLIFY-001)
(physmath, EXPR-SUBSTITUTE-001)
(physmath, RAT-NORMALIZE-001)
(physmath, DIFFERENTIATE-001)
(physmath, SELECT-BRANCH-001)
(physmath, EXPR-APPLY-001)
```

Any unknown rule pair MUST be rejected during deserialization or runtime validation.

---

# 25. Reference implementation acceptance suite

## 25.1 Identity

**Given** identical canonical payloads and identical namespace/type/version.  
**When** ArtifactRef is computed twice.  
**Then** digests are identical.

## 25.2 Primitive framing

**Given** two distinct payload sequences that would be ambiguous under raw concatenation.  
**When** encoded using `String`, `Bytes`, `Optional`, `Sequence`, and `Set` framing.  
**Then** their canonical byte streams remain distinct.

## 25.3 BigInt vectors

**Given** positive, negative, zero, and non-minimal magnitude candidates.  
**When** BigInt canonical encoding validates.  
**Then** minimal valid forms are accepted and negative zero/redundant leading zeros are rejected.

## 25.4 Parameter alpha-renaming

**Given** the same Function declaration structure with display name `t` versus `tau`.  
**When** ParameterID is computed.  
**Then** ParameterID is identical.

## 25.5 Parameter separation

**Given** two declaration positions in one Function.  
**When** ParameterIDs are computed.  
**Then** their ordinal-based identities differ.

## 25.6 Function body closure

**Given** a Function whose body contains a foreign ParameterID not present in `Function.ParameterIDs`.  
**When** the Function artifact is validated.  
**Then** validation rejects it.

## 25.7 Symbol physical neutrality

**Given** symbols named `t`, `k`, `c`, `G`, and `psi`.  
**When** PhysMath canonicalization runs.  
**Then** no physical role is attached from spelling alone.

## 25.8 Add/Mul normalization

**Given** nested Add/Mul structures with different insertion order.  
**When** canonicalized.  
**Then** equivalent commutative structures have one canonical child order.

## 25.9 Simplification recursion

**Given** nested simplifiable expressions.  
**When** `EXPR-SIMPLIFY-001` executes.  
**Then** child simplification is performed recursively in post-order and the result is deterministic and idempotent.

## 25.10 Condition preservation

**Given** an operation requiring `x != 0`.  
**When** the operation executes.  
**Then** the condition is explicit and is not silently discharged.

## 25.11 Rule registry and arity

**Given** an unknown RuleID pair, wrong input arity, or wrong target family.  
**When** runtime validation occurs.  
**Then** the operation is rejected.

## 25.12 EXPR-APPLY boundary

**Given** a root Call and its exact referenced Function artifact.  
**When** `EXPR-APPLY-001` executes.  
**Then** the Function body is expanded deterministically by simultaneous substitution.

**Given** a non-root nested Call.  
**When** the same operation is invoked without isolating the Call.  
**Then** validation rejects the operation.

**Given** a mismatched FunctionRef or arity.  
**When** validation occurs.  
**Then** validation rejects the operation.

## 25.13 Substitution context

**Given** an `EXPR-SUBSTITUTE-001` target with a BoundParameterRef.  
**When** the declaring Function is absent from the rule inputs.  
**Then** validation rejects the operation.

## 25.14 Differentiation context

**Given** `DIFFERENTIATE-001`.  
**When** the BoundParameterRef does not match the declaring Function input.  
**Then** validation rejects the operation.

## 25.15 Branch selection

**Given** a BranchSet and an in-range ordinal.  
**When** `SELECT-BRANCH-001` executes.  
**Then** the selected branch is returned exactly and no condition is generated.

**Given** an out-of-range ordinal or non-BranchSet root.  
**When** validation occurs.  
**Then** the operation fails.

## 25.16 Relation serialization

**Given** an equality Relation node.  
**When** serialized as an embedded expression and as a durable `physmath/relation/1`.  
**Then** the durable relation payload is exactly the relation body without the outer node tag, and left/right operands are embedded canonical expression bytes rather than ArtifactRefs.

## 25.17 Derivation determinism

**Given** the same premises, ordered derivation steps, producer namespaces, rule IDs, parameters, and conditions.  
**When** serialized and replayed.  
**Then** the canonical result references are reproducible.

## 25.18 Theorem completeness

**Given** a theorem with a claimed derivation.  
**When** validated.  
**Then** all premise, conclusion, assumption, and derivation references are present and type-correct, and the derivation is complete.

## 25.19 Equivariance protocol structure

**Given** the reference fixture.  
**When** the package is inspected.  
**Then** all of the following are present:

```text
AGENTS.md
manifest.json
protocol/derivation_test.go
```

The executable derivation test MUST be the mechanism that demonstrates the formal theorem representation.

## 25.20 Equivariance exact graph

**Given** the normative functions:

```text
f(v)       = 2*v
rho_X(g,v) = g*v
rho_Y(g,v) = g*v
```

**When** the derivation executes.  
**Then** the exact six transformation steps of §16.9 produce identical canonical final expressions `2*g*x`.

## 25.21 Equivariance theorem conclusion

**Given** the final LHS and RHS expression artifacts from the reference derivation.  
**When** the theorem conclusion is constructed.  
**Then** the durable relation embeds those exact canonical expressions, the Equation wraps that relation, and the Theorem references a complete Derivation.

## 25.22 Many-framework interpretation

**Given** the same mathematical theorem artifact.  
**When** two independent `phys-*` frameworks attach different physical assumptions or meanings.  
**Then** the PhysMath theorem identity remains unchanged while the framework artifacts remain distinct.

## 25.23 Framework dependency

**Given** `phys-bohm` declares `phys-qm` as a base framework.  
**When** a `phys-bohm` physical artifact references a `phys-qm` assumption.  
**Then** the reference is valid.

**Given** an unrelated `phys-cst` assumption.  
**When** a `phys-bohm` artifact references it without a dependency edge.  
**Then** validation rejects it.

## 25.24 Dynamic extension prohibition

**Given** a framework-local mathematical node not registered in PhysMath v1.  
**When** a package attempts to inject it into the PhysMath runtime AST.  
**Then** the operation is rejected; the package must keep the type local or define an explicit future PhysMath mapping.

## 25.25 Protocol manifest inheritance

**Given** the existing General Relativity formalization example and a new Phase-1 package.  
**When** the package is loaded.  
**Then** its `manifest.json` uses the same schema and loading convention as the GR example, with no alternative package-manifest protocol.

**Given** the GR protocol cannot express a required Phase-1 fixture without a new semantic manifest field.  
**When** the coding agent discovers the case during planning.  
**Then** it reports `PROTOCOL GAP` rather than silently inventing a second schema.

## 25.26 Condition-context determinism

**Given** `EXPR-SIMPLIFY-001` applied to `Pow(x,0)`.  
**When** an exact `MathCondition` for `x != 0` is supplied in `OperationInvocation.ConditionRefs`.  
**Then** the transformation succeeds.

**Given** the same input with no exact required condition.  
**When** the rule executes.  
**Then** the transformation fails with no output.

## 25.27 Rational-arithmetic trigger

**Given** `Add([Rational(1/2), Symbol(x), Rational(1/3)])`.  
**When** `EXPR-SIMPLIFY-001` executes.  
**Then** no partial rational folding occurs.

**Given** `Add([Rational(1/2), Rational(1/3)])`.  
**When** `EXPR-SIMPLIFY-001` executes.  
**Then** the entire node reduces to `Rational(5/6)`.

## 25.28 Duplicate ParameterID substitution

**Given** two substitution entries with different declaring artifacts but the same `ParameterID`.  
**When** `EXPR-SUBSTITUTE-001` validates.  
**Then** the operation is rejected.

## 25.29 DerivationStep rule-family enforcement

**Given** a DerivationStep whose input or output family violates the selected rule's §17.2 contract.  
**When** derivation validation occurs.  
**Then** the step is rejected.

## 25.30 ParameterID ordinal consistency

**Given** a Function whose stored `ParameterIDs[i]` differs from `DeriveParameterID(i)`.  
**When** the Function is validated.  
**Then** validation rejects the Function.

## 25.31 phys-* namespace validity

**Given** a namespace matching the `phys-*` grammar and publishing a valid FrameworkArtifact.  
**When** it is used as a framework namespace.  
**Then** it is valid without consulting a global registry.

**Given** a namespace that violates the namespace grammar.  
**When** it is used in an ArtifactRef.  
**Then** validation rejects it.

## 25.32 Equivariance staged-call linkage

**Given** `LHS_1` and `RHS_1` setup Calls.  
**When** the fixture validates staging.  
**Then** the designated argument bytes match the outputs of Steps 0 and 3 exactly.

## 25.33 Protocol refactoring safety

**Given** the existing `phys-gr` implementation.  
**When** the coding agent prepares the implementation plan.  
**Then** the plan distinguishes shared core extraction from GR-specific semantics and records every conflict with this specification before source changes begin.

## 25.34 Protocol replay determinism

**Given** two independent implementations following this specification and the same reference fixture.  
**When** both execute the same `protocol/derivation_test.go` derivation.  
**Then** the expected canonical outputs and golden ArtifactRefs are identical.

---

# 26. Implementation order

The implementation sequence is intentionally short and explicitly **plan-first** and **phys-gr-refactor-first**.

```text
0. Inspect existing General Relativity protocol example + existing phys-gr
      |
      v
1. Write implementation plan and blocker report
      |
      v
2. phys-artifact
      |
      v
3. PhysMath core identity + canonical codecs
      |
      v
4. Expression + Symbol + Rational + Add/Mul + Neg + Pow + Sqrt + Call + Relation
      |
      v
5. ParameterID + BoundParameterRef + constants
      |
      v
6. Function + FunctionApplication + Map + Tensor + Index + Collection
      |
      v
7. Transformation registry + TransformResult + explicit condition context
      |
      v
8. Derivation + Theorem artifacts
      |
      v
9. Refactor/extract phys-gr into the first concrete `phys-*` implementation
      |
      v
10. Refactor the GR protocol fixture into the canonical package pattern
      |
      v
11. Implement the equivariance reference fixture
      |
      v
12. second phys-* framework/interpretation fixture
      |
      v
13. canonical vectors + negative tests + conformance suite
```

TOE-Agent, Numerical, and research/trust layers do not appear in this sequence.

The existing `phys-gr` implementation is the primary concrete refactoring case for testing the architecture's core/specific split.

---

# 27. Architecture invariants to preserve permanently

```text
PhysMath = mathematical representation

phys-* = framework-specific physical representation

Paper = research context, not PhysMath ontology

Equation = primary formal representation unit

Derivation = explicit ordered transformation structure

Theorem = represented mathematical claim + derivation

Physical assumption = phys-* owned

Physical interpretation = phys-* owned

Mathematical identity = explicit structural identity

Symbol spelling != physical meaning

Multiple physical interpretations may coexist

No dynamic PhysMath AST extension

Phys remains frozen

Numerical is downstream and deferred

TOE-Agent is downstream and deferred

Trust/provenance/corpus is downstream and deferred

Formalization requires executable GR-style protocol packaging

Protocol packaging is inherited from the GR example, not reinvented

Implementation begins with a coding-agent plan, not source code

Existing phys-gr is refactoring substrate, not semantic authority

Representation success != scientific truth
```

The three complexity invariants are also permanent:

```text
Coupling
    -> separate conflated concerns

Cohesion
    -> unify necessary commonality through explicit interfaces

Core / specific
    -> common denominator in core, differences in specific layers
```

---

# 28. Final handoff rule

This specification is architecturally closed for Phase-1 implementation planning when the coding agent can satisfy all of the following:

```text
1. Every Phase-1 durable artifact family is registered.
2. Every ArtifactRef field has a closed target family or explicitly closed family set.
3. Every canonical primitive encoding is deterministic.
4. ParameterID and BoundParameterRef are deterministic and ordinal-consistent.
5. Every Phase-1 rule has exact arity, input families, parameters, condition context, traversal behavior, and failure behavior.
6. Durable Transformation and DerivationStep records carry ProducerNamespace.
7. Function body binding and capture behavior are deterministic.
8. phys-* namespace validity and dependency closure are explicit and acyclic.
9. DerivationStep input/output families are validated against the selected rule contract.
10. The equivariance fixture uses exactly one normative mathematical graph plus explicitly declared staged setup artifacts.
11. The equivariance fixture follows the GR-style AGENTS.md + manifest.json + protocol/derivation_test.go protocol.
12. The protocol package manifest uses the exact GR example schema and loading convention; no alternate Phase-1 manifest protocol exists.
13. The equivariance derivation closes using only the seven registered Phase-1 rules.
14. Staged Call linkage is explicit and requires no hidden AST subtree-replacement rule.
15. The theorem conclusion is constructed deterministically from the designated final derivation outputs.
16. Golden vectors and negative vectors cover the closed contracts.
17. The implementation plan identifies what must be extracted from phys-gr into the common denominator and what must remain GR-specific.
18. Any remaining implementation-time uncertainty is either a permitted implementation detail or a reported architectural blocker under §24.4.
```

The governing handoff statement is:

> **Given this specification, the frozen Phys authority where explicitly applicable, the existing General Relativity example as protocol precedent, and the existing `phys-gr` implementation as refactoring substrate, a competent Go developer can first produce an implementation plan for PhysMath and `phys-*` and can proceed to coding without inventing any mathematical identity rule, physical interpretation rule, serialization rule, derivation representation rule, protocol format, namespace-registration mechanism, condition-provision mechanism, or cross-layer ownership rule.**

The **next artifact is the coding agent's project-specific implementation plan**. That plan MUST use the existing `phys-gr` code as its primary concrete refactoring case.

During planning, the coding agent is expected to surface any genuine architectural blocker or missing decision to the human architect with concrete evidence and competing interpretations rather than silently choosing one. Once those issues are resolved—or once the plan demonstrates that none exist—implementation may begin.

No secondary architecture work is required merely to start planning or Phase-1 coding.

A future architecture revision is justified only by a demonstrated blocker or contradiction against this specification, the frozen Phys authority where applicable, or the inherited GR protocol precedent—not by stylistic preference, optimization, reviewer taste, or the desire to prematurely implement Phase-2 systems.

**END OF PM-REP-001 v1.4**
