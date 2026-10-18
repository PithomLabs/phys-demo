# PhysMath / Phys / phys-* Architecture
## Final Cross-Project Architecture Specification

**Document ID:** PM-ARCH-FOUNDATION-001  
**Version:** 1.5  
**Status:** **PROPOSED COMPLETE — BLOCKER-ONLY AUDIT PENDING**  
**Role:** Normative architecture specification; not an implementation plan.  
**Supersedes:** v1.4 and all earlier architecture decisions listed in §28.  
**Next artifact after audit:** project-specific implementation plan.  

> **Purpose.** Define a mechanically checkable boundary between pure mathematics, physical-domain semantics, trusted physical formalization, research state, and numerical execution. The document is designed so that a competent coding agent, given this document and the frozen Phys authority identified in §0, does not have to invent an architectural decision.

> **Non-goal.** This document does not define a Theory of Everything, adjudicate empirical truth, or turn Dancer/Dance/Stage into library ontology.

---

# 0. Authority, scope, and closure status

## 0.1 Frozen external authority

The current trusted Phys authority is:

```text
Module:             github.com/PithomLabs/phys
Frozen specification: specs_v2_3.md
Operational freeze: exact 42-file repository freeze record
Revision:           exact Git commit supplied to the implementation agent at run time
```

The implementation agent MUST treat the supplied repository revision plus `specs_v2_3.md` as authoritative for existing Phys semantics.

`README.md` is descriptive documentation. `AGENTS.md` is process/instruction authority. Neither may silently override `specs_v2_3.md` for semantic behavior.

Phys is **byte-frozen for this architecture**:

```text
NO new Phys imports
NO PhysMath import
NO new Phys API
NO architecture-driven edits to Phys
NO new Phys dependency solely to support this architecture
```

Any future change to the Phys kernel follows the post-MVP Growth Gate in §27; it is not part of this architecture's implementation.

## 0.2 Normative language

The terms **MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, **SHALL NOT**, **SHOULD**, **SHOULD NOT**, and **MAY** are normative.

When this document and an earlier architecture document conflict, this document wins.

When this document and the frozen Phys specification conflict about existing Phys behavior, the frozen Phys specification wins and this architecture must be amended rather than silently redefining Phys behavior.

## 0.3 Closure rule

This document is the **last architecture document in the current sequence**.

No v1.6 architecture document may be created for style, preference, optimization, or reviewer taste.

A new architecture revision is permitted only when a blocker meeting §32.2 is demonstrated with a concrete counterexample against this document and the frozen Phys authority.

After the blocker-only audit passes, the next artifact is an implementation plan, not another closure plan.

---

# 1. Constitutional invariants

## 1.1 Four distinct responsibilities

```text
PhysMath
    pure mathematics

phys-*
    physical-domain semantics

Phys
    trusted physical formalization and corpus authority

TOE-Agent
    research state, hypotheses, correspondences, formalization requests

Numerical
    numerical execution
```

These responsibilities are independent even when one project consumes artifacts produced by another.

## 1.2 Pure mathematics rule

PhysMath MUST NOT infer physical meaning from mathematical structure.

Examples that remain mathematical until a physical contract says otherwise:

```text
Manifold
Metric
Connection
Curvature
HilbertSpace
WaveFunction
Operator
Tensor
Function
```

A mathematical object does not become a physical object merely because it is commonly used in physics.

## 1.3 Trust rule

Phys owns the existing trusted physical metadata model:

```text
provenance
assumptions
conventions
corpus status
MRC/trust rules
trusted ledger/session semantics
```

External `phys-*` packages MUST NOT create a competing trust, corpus, or provenance ontology.

## 1.4 No automatic truth

No component in this architecture may infer:

```text
empirical truth
scientific confirmation
corpus establishment
physical correctness
```

from mathematical consistency, replay success, projection success, or contract validation alone.

---

# 2. Project boundaries and dependency topology

## 2.1 Projects

The target ecosystem contains:

```text
Project A — TOE-Agent / research program
Project B — PhysMath
Project C — phys-* physical-domain libraries
Existing authority — Phys
Numerical — numerical execution layer
Mechanical identity library — phys-artifact
```

`phys-artifact` is infrastructure, not a semantic authority.

## 2.2 Dependency graph

Production dependency edges are:

```text
PhysMath
    -> phys-artifact

Numerical / contract
    -> phys-artifact

Numerical / engine
    -> Numerical / contract
    -> phys-artifact

phys-*
    -> PhysMath
    -> Phys
    -> phys-artifact
    -> Numerical / contract

TOE-Agent
    -> PhysMath public API
    -> Phys public API
    -> phys-* public APIs
    -> phys-artifact
```

`Phys` has no new dependency edge.

There is no cycle:

```text
phys-* -> Numerical/contract -> phys-artifact
```

and Numerical does not import any `phys-*` implementation package.

## 2.3 Meaning of `PhysMath ⊥ Phys`

The constitutional dependency relation is:

```text
PhysMath  -X->  Phys
Phys      -X->  PhysMath
```

Both may independently depend on neutral infrastructure only when that infrastructure is non-semantic and the frozen authority permits the dependency.

For this architecture, `phys-artifact` is allowed for PhysMath and external projects; Phys itself remains frozen and does not import it.

## 2.4 `phys-artifact` ownership

`phys-artifact` owns only:

```text
ArtifactRef
artifact-family identity mechanics
canonical primitive framing helpers
reference-related validation errors
```

It MUST NOT own:

```text
math semantics
physics ontology
trust
corpus status
interpretation
research policy
registry state
```

It contains no universal artifact interface.

---

# 3. Artifact model

## 3.1 Artifact is a category, not a universal Go interface

The architecture recognizes durable artifacts as package-owned values with package-owned schemas.

There is **no** universal:

```go
type Artifact interface{}
```

There is **no** universal envelope carrying all math/physics/trust/research fields.

A concrete artifact is owned by its producer package.

## 3.2 In-memory composition

Within one process, components compose using ordinary typed Go values.

Examples:

```text
physmath.Expression
phys-gr.GRMetricalObject
phys.Object
numerical.Contract
```

An in-memory object MUST NOT require an `ArtifactRef` lookup merely to be usable inside the same typed API call.

## 3.3 Serialized composition

Across persistence or process boundaries, durable relationships use `ArtifactRef`.

Therefore:

```text
typed Go value != ArtifactRef
```

An `ArtifactRef` names an artifact. It does not resolve the artifact.

## 3.4 Reference resolution

`ArtifactRef` is a **content-addressed name**, not a resolver.

Resolution is an environment concern.

A storage environment MAY implement:

```text
Resolver.Resolve(ArtifactRef)
```

but:

```text
NO mandatory global registry
NO implicit network lookup
NO universal artifact service
NO semantic authority attached to resolution
```

A semantic-layer implementation MUST NOT assume that an `ArtifactRef` is network-resolvable.

## 3.5 Persistence rule

A durable artifact contains its canonical payload.

Its `ArtifactRef` is derived from that payload and the artifact-family identity.

The `ArtifactRef` is not recursively embedded in the payload that determines its own digest.

---

# 4. ArtifactRef and artifact identity

## 4.1 Canonical Go type

The field declaration order is intentionally identical to the canonical wire order:

```go
type ArtifactRef struct {
    Namespace     string
    SchemaVersion string
    ArtifactType  string
    Digest        [32]byte
}
```

This field order is normative.

Reflection-based serialization of `ArtifactRef` is forbidden because field declaration order alone does not define canonical artifact encoding; the explicit encoder in §5 does.

## 4.2 Identity tuple

Artifact identity is exactly:

```text
(Namespace, SchemaVersion, ArtifactType, Digest)
```

Two refs are equal iff all four fields are byte-for-byte equal after canonical UTF-8 encoding and exact digest comparison.

Digest equality alone is insufficient.

## 4.3 Content addressing

There is no separate `artifact_id` versus `content_digest` choice.

The architecture is content-addressed only.

The digest is SHA-256 over the exact framed bytes in §5.2.

## 4.4 Artifact family identifier

The non-instance family identity is the string:

```text
<Namespace>/<ArtifactType>/<SchemaVersion>
```

Example:

```text
phys/object/1
physmath/expression/1
physmath/function/1
```

This string is an **artifact-family identifier**, not an instance identity.

## 4.5 ArtifactRef is not itself an artifact

`ArtifactRef` MUST NOT be assigned another `ArtifactRef` merely because it is serialized.

It is a reference primitive.

---

# 5. Canonical byte encoding

## 5.1 Encoding principles

All canonical bytes are deterministic.

Canonical encodings MUST NOT depend on:

```text
Go map iteration order
pointer addresses
native memory layout
locale
pretty-print formatting
float decimal formatting
JSON field ordering
```

Transport encodings are separate from canonical identity encodings.

## 5.2 Artifact digest framing

For an artifact with canonical payload bytes `P`:

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

`Len32(S)` means the unsigned 32-bit big-endian integer equal to the number of UTF-8 **bytes** in `S`, not the number of Unicode runes.

This is the only normative artifact digest formula in v1.5.

## 5.3 Embedded ArtifactRef encoding

Whenever an `ArtifactRef` appears inside another canonical payload:

```text
EncodeArtifactRef(R) =
    Len32(R.Namespace UTF-8 bytes)
    || R.Namespace UTF-8 bytes
    || Len32(R.SchemaVersion UTF-8 bytes)
    || R.SchemaVersion UTF-8 bytes
    || Len32(R.ArtifactType UTF-8 bytes)
    || R.ArtifactType UTF-8 bytes
    || R.Digest[32]
```

Embedded refs MUST NOT be encoded in canonical payloads as:

```text
JSON objects
hex strings
bare 32-byte digests
pointer addresses
```

A transport document MAY display a ref as documented JSON fields; that transport representation is not the canonical payload.

## 5.4 Primitive canonical types

The canonical primitive vocabulary is:

```text
String        = UInt32BE(byte_length) || UTF-8 bytes
Bytes         = UInt32BE(byte_length) || raw bytes
UInt8         = 1 byte
UInt16        = 2-byte unsigned big-endian
UInt32        = 4-byte unsigned big-endian
UInt64        = 8-byte unsigned big-endian
Bool          = 0x00 or 0x01
Digest32      = exactly 32 bytes
Optional(T)   = presence byte 0x00/0x01 followed by T when present
Sequence(T)   = UInt32 count followed by elements in declared order
Set(T)        = UInt32 count followed by elements sorted by canonical element bytes
Enum          = canonical String of one closed value
TaggedUnion   = UInt8 discriminator followed by variant payload
```

A Set MUST reject duplicate canonical elements.

## 5.5 Arbitrary-precision integer

`BigInt` canonical encoding:

```text
Sign     = 0x00 for non-negative, 0x01 for negative
Magnitude = minimal big-endian unsigned magnitude
Encoding = Sign || UInt32BE(MagnitudeLength) || Magnitude
```

Zero is encoded as:

```text
0x00 || 0x00000000
```

Negative zero is invalid.

## 5.6 Rational

A rational is stored as two `BigInt` values:

```text
Numerator
Denominator
```

Canonical constraints:

```text
Denominator > 0
GCD(abs(Numerator), Denominator) = 1
```

A zero numerator MUST have denominator `1`.

## 5.7 Float64Bits

`Float64Bits` is the eight-byte big-endian IEEE-754 bit pattern of a `float64`.

Canonical numeric execution schemas MUST reject:

```text
NaN
+Inf
-Inf
```

`Float64Bits` is forbidden in PhysMath mathematical-expression artifacts and all trust/interpretation/research artifacts. It is permitted only in explicitly numerical schemas such as `ComputationalContract`.

## 5.8 Record encoding

A record's fields are encoded in the exact order declared by its schema.

No field names are included in canonical bytes unless the schema explicitly defines them as data.

Unknown fields are invalid for schema version `1`.

## 5.9 Collection ordering

For every collection the owning schema MUST declare `SET` or `SEQUENCE`.

For a `SET`, elements are sorted by lexicographic byte comparison of their exact canonical element encoding.

For an `ArtifactRef` set, that element encoding is §5.3.

For a homogeneous artifact family, sorting by the 32-byte digest is equivalent, but §5.3 remains the normative representation.

Traversal order MUST NOT determine canonical Set order.

---

# 6. Namespace registry

## 6.1 Current namespaces

The following namespaces are reserved:

| Namespace | Owner | Status |
|---|---|---|
| `physmath` | PhysMath | current |
| `phys` | frozen Phys authority | current |
| `toe` | TOE-Agent / Project A | current target |
| `numerical` | Numerical Layer | current target |
| `phys-qm` | quantum physical domain | post-MVP |
| `phys-bohm` | Bohmian physical domain | post-MVP |
| `phys-gr` | general-relativity physical domain | post-MVP |
| `phys-cst` | causal-set physical domain | post-MVP |
| `phys-cdt` | causal-dynamical-triangulation physical domain | post-MVP |

A future `phys-*` package MUST declare a new namespace before publishing artifacts.

`phys-artifact` has no semantic artifact namespace.

## 6.2 Namespace grammar

A namespace is ASCII lowercase and consists of one to sixty-four characters from:

```text
[a-z0-9][a-z0-9._-]{0,63}
```

No spaces or uppercase letters are permitted.

## 6.3 Artifact-type grammar

Artifact types consist of one or more lowercase path segments, each matching the same ASCII-safe identifier grammar; `/` separates segments.

Maximum encoded UTF-8 length is 128 bytes.

## 6.4 Schema-version grammar

`SchemaVersion` is a String, not an integer.

MVP schemas use:

```text
"1"
```

Schema-version strings inherit the identifier grammar and MAY use dotted semantic subversions only when a future schema explicitly defines them.

## 6.5 Free-form identifier grammar

Unless a schema explicitly declares a different grammar, these identifiers inherit the same ASCII-safe identifier grammar and 128-byte maximum:

```text
InterpretationID
ProposalID
RegimeID
ExtensionID
PhysicalType
OperatorID
FrameworkID
FunctionID
```

Human-readable display names are not identity-bearing unless a schema explicitly says they are.

---

# 7. Frozen Phys artifact references

## 7.1 Phys namespace

The ecosystem namespace for the frozen Phys authority is exactly:

```text
phys
```

## 7.2 MVP Phys referenceable artifact family

Only this frozen Phys family is referenceable through the ecosystem ArtifactRef protocol in MVP:

```text
phys/object/1
```

No session, ledger, manifest, or other internal Phys structure becomes an ecosystem artifact family through this document.

## 7.3 Frozen Phys canonical bytes

For a frozen Phys object, a consumer MUST obtain the exact canonical object bytes defined by the frozen Phys specification.

The consumer then computes the ecosystem digest using §5.2 with:

```text
Namespace     = "phys"
ArtifactType  = "object"
SchemaVersion = "1"
CanonicalPayload = exact frozen Phys canonical object bytes
```

The frozen Phys internal hash and the ecosystem ArtifactRef digest are distinct values and MUST NOT be conflated.

## 7.4 Phys participation rule

Phys is a **reference target**, not a new importer of `phys-artifact`.

Consumers compute refs to frozen Phys artifacts. Phys does not need a new import to participate in the ecosystem reference graph.

---

# 8. PhysMath mathematical IR — exact 16-sort set

PhysMath is closed **per schema version**.

The exact v1 sort inventory is:

```text
1.  Expression
2.  Function
3.  FunctionApplication
4.  Functional
5.  FunctionalApplication
6.  Operator
7.  OperatorApplication
8.  Tensor
9.  Index
10. Map
11. Relation
12. Collection
13. Transformation
14. Limit
15. AsymptoticTruncation
16. DerivationObligation
```

`Symbol` is not a separate sort. It is an `Expression` node kind.

`MathCondition` is not an `Expression` node kind or a fifteenth/sixteenth mathematical sort; it is a support artifact defined in §14.

`TransformResult` is an ephemeral in-memory result type defined in §15.

`DerivationObligation` is **DEFERRED + CONTAINMENT** in schema version `1`; no v1 runtime operation may create or evaluate one. This keeps the sort reserved without forcing an unfinished proof-object subsystem into MVP.

## 8.1 Per-sort contract requirement

Each MVP sort MUST specify:

```text
domain
codomain where applicable
construction
canonicalization
equality
serialization
reference targets
projection behavior
operation behavior
resource bounds
```

A contract cell is exactly one of:

```text
N/A
RULE <ContractRuleRef>
DEFERRED + CONTAINMENT
```

A plain "yes" or "supported" is not a contract value.

---

# 9. Expression

## 9.1 Expression node kinds

The closed v1 `Expression` node-kind set is exactly:

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

## 9.2 NodeKind canonical names

For loss reports and diagnostics, the canonical string names are exactly:

```text
Symbol
Rational
Add
Mul
Neg
Pow
Sqrt
Call
Relation
BranchSet
```

These names are canonical identifiers, not UI labels.

## 9.3 Symbol schema

```text
Symbol {
    DisplayName string
    ParameterID Optional<ParameterID>
    ConstantID  Optional<ConstantID>
}
```

Validation:

```text
durable expression:
    exactly one of ParameterID / ConstantID is present

ephemeral interactive expression:
    neither may be present
    both may never be present
```

A Symbol with both fields present is malformed everywhere.

A durable expression containing an unbound Symbol is malformed.

## 9.4 Rational node

```text
Rational {
    Numerator   BigInt
    Denominator BigInt
}
```

The value MUST satisfy §5.6 canonical rational rules.

## 9.5 Add and Mul canonicalization

`Add` and `Mul` use n-ary child sequences internally.

Construction-time canonicalization MUST:

```text
flatten nested same-operator nodes
canonicalize children first
sort children by canonical child payload bytes
```

Arity rules:

```text
n = 0:
    Add([]) -> Rational(0/1)
    Mul([]) -> Rational(1/1)

n = 1:
    Add([x]) -> x
    Mul([x]) -> x

n >= 2:
    retain n-ary node
```

The empty forms never survive canonicalization.

Semantic combination such as rational addition or multiplication is a normalization rule, not merely construction-time canonicalization.

## 9.6 Neg

`Neg` has one child.

The structural canonicalizer retains it.

`EXPR-SIMPLIFY-001` may apply:

```text
Neg(Neg(x)) -> x
```

and records no condition.

## 9.7 Pow

```text
Pow {
    Base     Expression
    Exponent Expression
}
```

PhysMath v1 permits any expression as an exponent.

Projection to frozen Phys permits only a Rational exponent matching the frozen Phys `Pow` exponent contract.

Non-rational exponents therefore map to:

```text
ProjectionClass = UNREPRESENTABLE
ReasonCode      = NON_RATIONAL_EXPONENT_IN_POW
NodeKind        = Pow
```

No silent rational approximation is permitted.

## 9.8 Sqrt

`Sqrt` has one Expression child.

It is a distinct node kind even though a future normalization may map it to `Pow(x,1/2)`.

The v1 projection contract treats Sqrt as directly representable because frozen Phys has a distinct Sqrt node.

## 9.9 Call

```text
Call {
    FunctionRef ArtifactRef
    Args       []ArtifactRef
}
```

`FunctionRef` MUST target:

```text
physmath/function/1
```

`Args` are a `SEQUENCE`; positional order is semantic.

## 9.10 Relation — two-form rule

There are two distinct concepts:

```text
Expression node: Relation
Durable artifact: physmath/relation/1
```

The expression node is embedded within an Expression tree.

The durable relation artifact contains the relation body without an additional nested node tag.

```text
Relation {
    Operator ComparisonOperator
    Left     ArtifactRef
    Right    ArtifactRef
}
```

`ComparisonOperator` is the closed v1 set:

```text
EQ
NEQ
LT
LTE
GT
GTE
```

A `MathCondition.Predicate` points to a durable `physmath/relation/1` artifact.

## 9.11 BranchSet

```text
BranchSet {
    Branches []ArtifactRef
}
```

`Branches` is a `SEQUENCE`.

BranchSet is valid mathematical structure even when frozen Phys cannot represent it.

## 9.12 Expression equality

Expression equality is structural equality of canonical node payloads.

It is not physical equivalence.

No theorem prover is invoked.

---

# 10. Mathematical parameters and binding

## 10.1 ParameterID type

```go
type ParameterID [32]byte
```

`ParameterID` is a deterministic local mathematical binding identity.

It is **not globally unique** and MUST NOT be treated as provenance equivalence.

## 10.2 Scope path

A parameter is identified structurally by its declaration position.

`ScopePathBytes` is:

```text
UInt32BE(depth)
|| UInt32BE(scope_0_ordinal)
|| ...
|| UInt32BE(scope_(depth-1)_ordinal)
```

where each ordinal identifies the lexical binder position within its parent declaration.

Display names do not participate in ParameterID construction.

## 10.3 ParameterID construction

MVP ParameterID is:

```text
ParameterID = SHA-256(
    Len32("physmath" UTF-8 bytes)
    || "physmath" UTF-8 bytes
    || Len32(ScopePathBytes)
    || ScopePathBytes
    || UInt32BE(parameter_ordinal)
)
```

`parameter_ordinal` is zero-based within the declaring function/declaration.

The construction is deterministic across independent implementations.

## 10.4 Alpha-renaming invariant

Changing display text:

```text
t -> τ
```

MUST NOT change ParameterID.

The declaration position remains the identity source.

## 10.5 Cross-artifact parameter binding

Because ParameterID is intentionally local, a cross-artifact target uses:

```text
BoundParameterRef {
    DeclaringArtifactRef ArtifactRef
    ParameterID         ParameterID
}
```

There is no recursive dependency because ParameterID does not contain the declaring artifact's ref.

The declaring artifact does not embed its own ArtifactRef merely to derive its ParameterID.

## 10.6 Resolution rule

A `BoundParameterRef` resolves only against its explicitly named declaring artifact.

An operation receiving an ambiguous bare ParameterID across artifacts MUST reject the request.

The operation MUST NOT infer a declaration from matching display names.

## 10.7 t/k rule

`t` and `k` are not distinguished by names such as `"t"` and `"k"`.

They are distinct only when they are distinct ParameterIDs within their declaration contexts.

Physical roles such as:

```text
time
RG scale
```

are supplied by physical contracts outside PhysMath.

PhysMath does not define `TIME`, `RG_SCALE`, `COORDINATE`, or equivalent physics-role enums.

## 10.8 Differential operator invariance

A differentiation operation contains:

```text
Target BoundParameterRef
Order UInt32
```

Renaming a displayed target does not alter the target ParameterID.

A different parameter identity cannot be substituted merely by string comparison.

---

# 11. Mathematical constants

## 11.1 ConstantID

```go
type ConstantID string
```

`ConstantID` values use the reserved namespace:

```text
math.const
```

MVP closed registry:

```text
math.const.pi
math.const.e
math.const.i
```

No other ConstantID is valid in schema version `1`.

## 11.2 Constant semantics

The constant is retained symbolically in PhysMath.

Examples:

```text
Symbol{DisplayName:"π", ConstantID:"math.const.pi"}
Symbol{DisplayName:"e", ConstantID:"math.const.e"}
```

The display name is presentation. `ConstantID` carries the mathematical identity.

## 11.3 Physical constants

Physical constants such as `c`, `hbar`, or `G` are not mathematical ConstantIDs in PhysMath v1.

They are ordinary mathematical Symbols whose physical meaning is supplied by a `phys-*` contract.

## 11.4 Numeric lowering of constants

Exact mathematical or physical constants MUST be resolved by the owning mathematical/physical contract **before** a `ComputationalContract` is constructed.

Numerical MUST NOT:

```text
look up symbolic constants
infer symbolic constants
reinterpret symbolic constants
recompute symbolic constants
```

The phrase "numerical lowering" means the producer-side lowering phase only; it does not mean the Numerical engine.

---

# 12. Higher mathematical sorts

The following schemas are normative. Operations marked `DEFERRED + CONTAINMENT` are represented and validated but their higher-order transformation semantics are not implemented in PhysMath v1.

## 12.1 Function

```text
Function {
    ParameterIDs []ParameterID
    BodyRef      ArtifactRef
}
```

Rules:

```text
BodyRef -> physmath/expression/1
ParameterIDs -> SEQUENCE
parameter order is semantic
```

Construction and serialization are MVP.

Call projection relies on this family.

## 12.2 FunctionApplication

```text
FunctionApplication {
    FunctionRef ArtifactRef
    Args        []ArtifactRef
}
```

```text
FunctionRef -> physmath/function/1
Args -> SEQUENCE
```

Higher-order evaluation is deferred.

## 12.3 Functional

```text
Functional {
    DomainSort    MathSortID
    CodomainSort  MathSortID
    DefinitionRef ArtifactRef
}
```

`DefinitionRef` MUST target `physmath/transformation/1` in v1.

Operational semantics are deferred.

## 12.4 FunctionalApplication

```text
FunctionalApplication {
    FunctionalRef ArtifactRef
    Args         []ArtifactRef
}
```

`Args` is a `SEQUENCE`.

Operational semantics are deferred.

## 12.5 Operator

User-defined Operator artifacts are **DEFERRED + CONTAINMENT** in v1.

The v1 core differential operator does not require an Operator artifact; see §13.

If represented by future schemas, `DefinitionRef` MUST target `physmath/transformation/1`.

## 12.6 OperatorApplication

This sort is MVP because it provides the canonical differentiation operation.

```text
OperatorApplication {
    OperatorID string
    OperatorRef Optional<ArtifactRef>
    Args       []ArtifactRef
    Parameters []byte
}
```

Rules:

```text
OperatorID = "differentiate" in MVP
OperatorRef MUST be absent for "differentiate"
Args = SEQUENCE
Parameters = Rule-specific canonical bytes
```

The field `Parameters` is the explicit home for differential parameters. There is no hidden side channel.

## 12.7 Tensor

```text
Tensor {
    Rank        UInt16
    IndexRefs   []ArtifactRef
    ComponentRefs []ArtifactRef
}
```

Rules:

```text
IndexRefs -> SEQUENCE, slot order is semantic
ComponentRefs -> SEQUENCE, row-major slot order is semantic
```

General tensor algebra is deferred.

## 12.8 Index

```text
Index {
    Label    string
    Variance IndexVariance
    SpaceRef Optional<ArtifactRef>
}
```

`IndexVariance`:

```text
COVARIANT
CONTRAVARIANT
UNSPECIFIED
```

The labels and variance are mathematical metadata, not physical ontology.

General index contraction is deferred.

## 12.9 Map

```text
Map {
    DomainSort    MathSortID
    CodomainSort  MathSortID
    DefinitionRef ArtifactRef
}
```

`DefinitionRef` MUST target `physmath/expression/1` in v1.

Operational semantics are deferred.

## 12.10 Collection

```text
Collection {
    ElementRefs []ArtifactRef
}
```

`ElementRefs` is a `SET`.

Duplicates are invalid.

## 12.11 Limit

```text
Limit {
    ExpressionRef ArtifactRef
    Target        BoundParameterRef
    ApproachRef   ArtifactRef
}
```

`ApproachRef` MUST target `physmath/expression/1`.

Construction and serialization are MVP.

General limit solving beyond explicitly implemented rules is `OUT_OF_SCOPE`.

## 12.12 AsymptoticTruncation

```text
AsymptoticTruncation {
    ExpressionRef ArtifactRef
    Target        BoundParameterRef
    Order         UInt32
}
```

Representation is reserved but operational truncation is deferred.

## 12.13 DerivationObligation

Reserved only:

```text
Status = DEFERRED + CONTAINMENT
```

No v1 implementation may claim proof obligations are solved merely because a derivation trace exists.

---

# 13. Canonical differential operator

## 13.1 Reserved operator identity

The canonical v1 differential operator is a reserved `OperatorID` value:

```text
differentiate
```

It is **not** represented by a separate core `Operator` artifact.

Therefore no implementation-specific `OperatorRef` can cause two installations to assign different identities to differentiation.

This explicitly chooses the reserved-tag variant of the previous design.

## 13.2 Differential canonical parameters

`DIFFERENTIATE-001` parameter bytes are:

```text
Encode(BoundParameterRef.Target)
|| UInt32BE(Order)
```

`BoundParameterRef.Target` uses §5.3 encoding for its embedded ArtifactRef followed by the exact 32-byte ParameterID.

`Order = 1` and `Order = 2` are valid in MVP.

For repeated differentiation, the canonical form is one `DIFFERENTIATE-001` operation with the summed `Order`, not nested repeated operator applications merely to express repetition.

## 13.3 Differential semantics

For `Order >= 1`, `DIFFERENTIATE-001` is defined only for the bounded rules in §15.5.

An unsupported derivative request returns:

```text
OperationOutcome = OUT_OF_SCOPE
```

It MUST NOT invent derivative rules for arbitrary functions or operators.

---

# 14. MathCondition

## 14.1 Durable support artifact

```text
MathCondition
```

is a durable PhysMath support artifact with family:

```text
physmath/condition/1
```

Schema:

```text
MathCondition {
    PredicateRef ArtifactRef
}
```

`PredicateRef` MUST target:

```text
physmath/relation/1
```

## 14.2 Purity

MathCondition contains:

```text
mathematical predicate only
```

It contains no:

```text
trust
provenance
corpus status
physical interpretation
physical assumption authority
```

## 14.3 Generation

A mathematical transformation MAY generate mathematical conditions.

Conditions are carried explicitly as an operation result; they are never attached by mutable global state.

## 14.4 Condition finalization

Before a TransformResult is finalized:

```text
1. canonicalize each MathCondition
2. compute its canonical digest
3. sort ascending by digest bytes
4. deduplicate identical conditions
```

Duplicate conditions are removed.

The transformation's tree traversal order cannot determine the final condition order.

## 14.5 Condition discharge

MVP contains **no condition-discharge operation**.

Conditions accumulate.

Future discharge requires an explicit versioned operation and cannot be smuggled into `Simplify`.

---

# 15. Transformations and operations

## 15.1 TransformResult is ephemeral

```go
type TransformResult struct {
    Output     ArtifactRef
    Conditions []ArtifactRef
    Trace      MathTrace
}
```

`TransformResult` is strictly an in-memory return type.

It is **not** a durable artifact, has no ArtifactRef of its own, and has no canonical wire format.

The referenced output/condition/trace artifacts are durable when the caller chooses to persist them.

## 15.2 Durable Transformation schema

```text
Transformation {
    RuleID     TransformationRuleID
    Inputs     []ArtifactRef
    Output     ArtifactRef
    Parameters Bytes
}
```

`Inputs` is a `SEQUENCE` because operation argument position can be semantic.

`Parameters` is never an open-ended JSON/object bag.

## 15.3 TransformationRuleID closed vocabulary

The v1 serialized Transformation registry is exactly:

```text
EXPR-CANON-001
EXPR-SIMPLIFY-001
EXPR-SUBSTITUTE-001
RAT-NORMALIZE-001
DIFFERENTIATE-001
SELECT-BRANCH-001
PROJECT-STRICT-001
```

No wildcard and no unregistered RuleID is legal.

`ID-001` is deliberately removed from the transformation registry; serialization and identity are not mathematical transformation operations.

## 15.4 RuleID versus ContractRuleRef

The architecture uses two different concepts:

```text
TransformationRuleID
    serialized runtime operation identity

ContractRuleRef
    documentation-only identifier used by the sort contract matrix
```

A `ContractRuleRef` MUST NOT be accepted as a `Transformation.RuleID` unless it is separately and explicitly present in §15.3.

This prevents the earlier Appendix-B/Appendix-C semantic collision.

## 15.5 Exact TransformationRule semantics

### EXPR-CANON-001

Input:

```text
one Expression artifact
```

Parameters:

```text
empty byte sequence
```

Output:

```text
construction-canonical Expression
```

Semantics:

```text
canonicalize children recursively
flatten Add/Mul
apply Add/Mul arity normalization
sort associative children deterministically
preserve mathematical meaning
perform no semantic theorem proving
```

### EXPR-SIMPLIFY-001

Input:

```text
one Expression artifact
```

Parameters:

```text
empty byte sequence
```

Closed rewrite rules:

```text
Neg(Neg(x)) -> x
Add(..., 0, ...) -> Add(...)
Mul(..., 1, ...) -> Mul(...)
Mul(..., 0, ...) -> 0
Rational arithmetic -> canonical Rational
Pow(nonzero rational, 0) -> 1
```

For `Pow(x,0)` where nonzero cannot be established from the local exact form, the implementation MAY rewrite only by producing the explicit MathCondition `x != 0`. It MUST NOT assume the condition silently.

No general theorem prover is used.

### EXPR-SUBSTITUTE-001

Parameters are a `SET` of:

```text
SubstitutionEntry {
    Target      BoundParameterRef
    Replacement ArtifactRef
}
```

`Replacement` MUST target `physmath/expression/1`.

Targets are canonically sorted by their exact canonical bytes.

Duplicate targets are invalid.

Substitution is simultaneous, not sequential.

Ambiguous targets fail closed.

### RAT-NORMALIZE-001

Parameters:

```text
empty
```

Input and output are Rational artifacts whose numerator/denominator satisfy §5.6.

### DIFFERENTIATE-001

Parameters:

```text
Target BoundParameterRef
Order  UInt32
```

Closed derivative rules:

```text
d(Rational)/dp = 0

d(Symbol(target))/dp = 1

d(Symbol(other parameter))/dp = 0

d(Add(xs))/dp = Add(d(x1), ..., d(xn))

d(Mul(xs))/dp = sum_i(product_except_i(x) * d(x_i))

d(Neg(x))/dp = Neg(d(x))

d(Pow(b, n))/dp = n * Pow(b, n-1) * d(b)
    only when n is a nonnegative integer Rational
```

`Sqrt`, `Call`, `Relation`, and `BranchSet` differentiation are `OUT_OF_SCOPE` in v1 unless an explicitly versioned future rule adds them.

### SELECT-BRANCH-001

Parameters:

```text
BranchOrdinal UInt32
```

Input MUST be a `BranchSet` artifact.

The selected ordinal MUST be within range.

Output is the selected branch expression.

### PROJECT-STRICT-001

Parameters:

```text
RequestedTarget ArtifactFamilyID
ProjectionPolicy string
```

MVP requires:

```text
RequestedTarget = phys/object/1
ProjectionPolicy = STRICT_MIRROR_ONLY
```

It performs the projection algorithm in §20.

## 15.6 Operation outcomes

An individual operation invocation returns exactly one terminal outcome:

```text
SUCCESS
UNDEFINED
TRANSFORMATION_FAILED
OUT_OF_SCOPE
RESOURCE_LIMIT
```

These are result classifications, not lifecycle states.

No operation outcome can assert scientific truth.

---

# 16. Operation schema versioning and MathTrace

## 16.1 OperationSchemaVersion

`MathTrace.OperationSchemaVersion` is the version of the **operation parameter schema and replay contract**, not the whole PhysMath version and not a date.

MVP value:

```text
"1"
```

If the canonical parameter bytes or replay semantics of a RuleID change incompatibly, the operation schema version MUST change.

## 16.2 MathTrace schema

```text
MathTrace {
    TraceSchemaVersion       string
    ProducerNamespace        string
    OperationSchemaVersion   string
    RuleID                   TransformationRuleID
    Inputs                   []ArtifactRef
    Output                   ArtifactRef
    ParameterBytes           Bytes
    ConditionRefs            []ArtifactRef
}
```

Rules:

```text
TraceSchemaVersion = "1" in MVP
Inputs = SEQUENCE
ConditionRefs = SET
ParameterBytes exactly equal canonical parameters for RuleID
```

## 16.3 Replay authority

A trace is replayed only by the authority that owns the RuleID vocabulary.

Examples:

```text
PhysMath trace -> PhysMath replay
phys-gr trace  -> phys-gr replay
Phys ledger    -> Phys replay
```

An authority MUST NOT claim to replay another authority's internal operations merely because it has an ArtifactRef to the output.

## 16.4 Byte replay rule

Replay validation compares canonical stored bytes and exact recorded parameters.

A replay implementation MUST NOT regenerate an artifact and then trust only a matching semantic summary.

---

# 17. Theory-local mathematical algebras

A `phys-*` package MAY have a theory-local mathematical algebra.

Example:

```text
phys-gr.GRXExpr
```

The theory-local algebra remains package-local.

It MUST NOT become a dynamic extension node inside PhysMath.

A theory-local algebra MAY expose an explicit deterministic projection into PhysMath when a value lies in the frozen PhysMath subset.

No forced migration is required in MVP.

The promotion rule is:

```text
theory-local algebra
    -> deterministic mapping
    -> existing PhysMath sort
```

not:

```text
theory-local algebra
    -> runtime PhysMath plugin node
```

PhysMath is closed per schema version, not closed forever.

A future mathematical extension requires a new versioned PhysMath schema/architecture decision.

---

# 18. Mathematical projection to Phys

## 18.1 Exact projection subset

The v1 projection mirror is exactly these Expression node kinds:

```text
Symbol
Rational
Add
Mul
Neg
Pow
Sqrt
```

The projection subset is a semantic mirror of the frozen Phys expression subset.

It is not a promise that every PhysMath structure is Phys-representable.

## 18.2 Non-projectable expression kinds

These remain valid PhysMath mathematics but are not projectable under the v1 strict mirror:

```text
Call
Relation
BranchSet
```

when they cannot be reduced to a supported expression form under an explicit rule.

Higher mathematical sorts such as Tensor, Function, Functional, Limit, and AsymptoticTruncation are not projectable under the MVP mirror unless explicitly reduced to the supported expression subset by a deterministic upstream transformation.

## 18.3 Projection classes

```text
FULLY_REPRESENTABLE
PARTIALLY_REPRESENTABLE
UNREPRESENTABLE
```

MVP permits **zero permitted-loss contracts**.

Therefore `PARTIALLY_REPRESENTABLE` exists in the vocabulary but is unreachable in MVP.

A future version may add a versioned permitted-loss contract.

## 18.4 Composite projection rule

For a composite parent:

```text
FULL
    iff the parent itself is representable
    and every required child is FULL
    and no required information is lost

PARTIAL
    iff a versioned projection contract explicitly permits information loss
    and every discarded component is recorded

UNREPRESENTABLE
    otherwise
```

A projectable child never makes an unprojectable parent projectable.

## 18.5 Phys canonicalization handoff

Projection does not have to be injective.

A valid PhysMath source MAY contain mathematically equivalent structure that is canonicalized differently from frozen Phys construction.

The projection implementation MUST construct the Phys expression through the frozen Phys public construction path, allowing frozen Phys canonicalization to produce the trusted target form.

Example:

```text
PhysMath Neg(Neg(x))
    -> projection
    -> Phys constructor canonicalizes to x
```

The projection is still valid when the target is canonical.

## 18.6 Binding metadata rule

Phys does not represent PhysMath ParameterID/ConstantID metadata in the same way.

Such metadata is not automatically treated as information-loss failure.

It is safe to omit only when omission cannot merge distinct source semantics.

### Symbol projection collision algorithm

For every projected Expression:

```text
1. Traverse all Symbol nodes in the source expression.
2. Group them by DisplayName.
3. Within each group, compare their binding metadata:
       ParameterID, or ConstantID.
4. If two source symbols with the same DisplayName have different binding metadata,
   projection MUST fail.
```

Failure:

```text
Classification = UNREPRESENTABLE
ReasonCode    = SYMBOL_BINDING_COLLISION
NodeKind      = Symbol
```

If all same-spelling symbols share the same binding identity, the omitted source metadata does not constitute a merge and projection may continue.

## 18.7 Pow projection

A PhysMath Pow projects only if its exponent is a Rational expression satisfying frozen Phys's exact exponent contract.

Otherwise:

```text
ReasonCode = NON_RATIONAL_EXPONENT_IN_POW
```

The source is not rewritten merely to force projectability.

## 18.8 Associative projection

PhysMath Add/Mul is canonicalized to deterministic n-ary children.

Frozen Phys uses binary nodes.

Projection folds children left-associatively in canonical order:

```text
Add([a,b,c]) -> Phys.Add(Phys.Add(a,b),c)
Mul([a,b,c]) -> Phys.Mul(Phys.Mul(a,b),c)
```

Arity `0` and `1` are already eliminated by §9.5.

---

# 19. Projection loss reporting

## 19.1 ProjectionLossReport

```text
ProjectionLossReport {
    Reasons []ProjectionReason
}
```

`Reasons` is a `SET`.

## 19.2 ProjectionReason

```text
ProjectionReason {
    Code     ProjectionReasonCode
    NodeKind NodeKind
    Detail   string
}
```

## 19.3 Closed reason-code vocabulary

MVP:

```text
NON_RATIONAL_EXPONENT_IN_POW
SYMBOL_BINDING_COLLISION
UNSUPPORTED_NODE_KIND
UNSUPPORTED_TOP_LEVEL_SORT
UNRESOLVED_FUNCTION_REFERENCE
LOSS_NOT_PERMITTED
INVALID_SOURCE_ARTIFACT
```

No `ProjectionReasonCode` outside this set is valid in v1.

## 19.4 NodeKind

`NodeKind` is exactly the ten names in §9.2.

It MUST NOT be encoded as an arbitrary spelling or numeric label.

Canonical encoding is the canonical Enum/String value.

## 19.5 Reason ordering

When reasons are materialized as a deterministic sequence for transport or display, they are sorted by:

```text
Code
then NodeKind
then Detail
```

using canonical UTF-8 byte order.

---

# 20. Physical interpretation and physical contracts

## 20.1 PhysicalInterpretation

A `PhysicalInterpretation` is a `phys-*` owned artifact that gives mathematical structures physical meaning within a declared framework.

Conceptual category:

```text
PhysicalInterpretation
```

is not a universal shared Go interface.

Concrete artifact families remain package-owned.

## 20.2 Schema

The canonical interpretation contract is:

```text
PhysicalInterpretation {
    InterpretationID       string
    FrameworkRef            ArtifactRef
    PhysicalType            string
    MathRefs                []ArtifactRef
    RequiredConditions      []ArtifactRef
    RequiredConventions     []ConventionRequirement
    InterpretationClass     InterpretationClass
}
```

Collection rules:

```text
MathRefs             = SET
RequiredConditions   = SET
RequiredConventions  = SET
```

## 20.3 InterpretationClass

Closed v1 vocabulary:

```text
INTERPRETATION_NEUTRAL
FRAMEWORK_SPECIFIC
```

Canonical encoding is the exact enum string.

A physical framework SHOULD normally use `FRAMEWORK_SPECIFIC`.

`INTERPRETATION_NEUTRAL` is reserved for a physical contract that explicitly asserts framework-neutral semantics.

## 20.4 FrameworkRef

`FrameworkRef` MUST reference a durable framework artifact:

```text
<phys-namespace>/framework/1
```

Examples:

```text
phys-gr/framework/1
phys-bohm/framework/1
```

A framework artifact is not a trust assertion. It identifies the owning semantic framework.

## 20.5 Framework artifact

Each `phys-*` package that publishes physical interpretations MUST publish exactly one framework artifact per framework version:

```text
FrameworkArtifact {
    FrameworkID      string
    PackageNamespace string
    FrameworkVersion string
}
```

The artifact's namespace is the package namespace.

## 20.6 PhysicalContract

`PhysicalContract` is declarative metadata.

It is not executable code.

```text
PhysicalContract {
    ExtensionID          string
    SchemaVersion        string
    PhysicalType         string
    FrameworkRef         ArtifactRef
    MathRefs             []ArtifactRef
    RequiredConditions   []ArtifactRef
    RequiredConventions  []ConventionRequirement
    InterpretationClass  InterpretationClass
}
```

`CompatibilityRules` is deliberately absent from v1.

Validation logic lives in the owning `phys-*` package and is not serialized as code.

## 20.7 ConventionRequirement

```text
ConventionRequirement {
    OwnerNamespace string
    Key            string
    RequiredValue  string
}
```

This is a declarative requirement, not a second convention authority.

It MUST NOT claim ownership of Phys corpus/convention truth.

## 20.8 Physical artifact rule

There is no universal `PhysicalArtifact` Go envelope.

Each `phys-*` package defines concrete physical artifact schemas within its namespace.

A physical artifact MAY reference:

```text
MathArtifact refs
Framework refs
PhysicalContract refs
Phys refs
```

as explicitly defined by its package contract.

---

# 21. ProposedInterpretation and CorrespondenceProposal

## 21.1 ProposedInterpretation

Owned by Project A.

It is a research claim about possible physical meaning.

It is not a trusted physical contract.

Canonical form:

```text
ProposedInterpretation {
    ProposalID          string
    FrameworkRef        Optional<ArtifactRef>
    MathRefs            []ArtifactRef
    ProposedPhysicalType string
    Scope               string
    NotesRef            Optional<ArtifactRef>
}
```

`MathRefs` is a SET.

No promotion API exists.

## 21.2 Relationship to PhysicalInterpretation

There is no generic:

```text
Promote()
Accept()
Approve()
```

A `phys-*` framework MAY independently publish a `PhysicalInterpretation` that corresponds to a Project A proposal.

The two artifacts remain distinct.

## 21.3 CorrespondenceProposal

Owned by Project A.

```text
CorrespondenceProposal {
    ProposalID        string
    ProposalVersion   string
    SourceRefs        []ArtifactRef
    TargetRefs        []ArtifactRef
    MappingRefs       []ArtifactRef
    FrameworkRefs     []ArtifactRef
    ClaimClass        CorrespondenceClaimClass
    Scope             string
}
```

Collections:

```text
SourceRefs   = SET
TargetRefs   = SET
MappingRefs  = SET
FrameworkRefs = SET
```

`CorrespondenceProposal` is not a `PhysicalInterpretation` and not a trust assertion.

## 21.4 CorrespondenceClaimClass

Closed v1 set:

```text
CONJECTURAL
STRUCTURAL
FORMAL
EMPIRICAL_CLAIM
```

`EMPIRICAL_CLAIM` remains a research statement; the architecture does not validate its truth.

---

# 22. Research hypotheses and lifecycle state

## 22.1 ResearchHypothesis

Owned by Project A.

```text
ResearchHypothesis {
    HypothesisID string
    Version      string
    State        ResearchState
    SubjectRefs  []ArtifactRef
    EvidenceRefs []ArtifactRef
    ParentRef    Optional<ArtifactRef>
    Summary      string
}
```

`SubjectRefs` and `EvidenceRefs` are SETs.

## 22.2 ResearchState closed vocabulary

```text
PROPOSED
UNDER_ATTACK
REVISION_REQUIRED
COUNTEREXAMPLE_FOUND
INCONSISTENT
AWAITING_FORMALIZATION
FORMALIZED
CLOSED
```

## 22.3 ResearchState transitions

```text
PROPOSED
    -> UNDER_ATTACK
    -> REVISION_REQUIRED
    -> AWAITING_FORMALIZATION
    -> CLOSED

UNDER_ATTACK
    -> REVISION_REQUIRED
    -> COUNTEREXAMPLE_FOUND
    -> INCONSISTENT
    -> AWAITING_FORMALIZATION

REVISION_REQUIRED
    -> PROPOSED

COUNTEREXAMPLE_FOUND
    -> REVISION_REQUIRED
    -> CLOSED

INCONSISTENT
    -> REVISION_REQUIRED
    -> CLOSED

AWAITING_FORMALIZATION
    -> FORMALIZED
    -> REVISION_REQUIRED

FORMALIZED
    -> UNDER_ATTACK
    -> REVISION_REQUIRED
    -> CLOSED

CLOSED
    -> terminal
```

The system MUST reject transitions not listed above.

## 22.4 Status versus truth

`FORMALIZED` means a formalization path succeeded.

It does not mean:

```text
scientifically true
empirically confirmed
Phys corpus established
```

---

# 23. Formalization protocol

## 23.1 FormalizationProposal

```text
FormalizationProposal {
    ProposalRef        ArtifactRef
    ProposalRevision    UInt32
    ParentProposalRef   Optional<ArtifactRef>
    ProducerNamespace   string
    SourceArtifactRefs  []ArtifactRef
    RequestedTarget     ArtifactFamilyID
    ProjectionPolicy    ProjectionPolicy
    RequiredFrameworkRefs []ArtifactRef
}
```

Rules:

```text
ProposalRevision initial value = 1
ParentProposalRef absent on first revision
ParentProposalRef points to the immediately preceding resolved proposal
SourceArtifactRefs = SET
RequiredFrameworkRefs = SET
RequestedTarget = "phys/object/1" in MVP
ProjectionPolicy = STRICT_MIRROR_ONLY in MVP
```

## 23.2 ProposalState

```text
DRAFT
SUBMITTED
RESOLVED
```

## 23.3 ProposalOutcome

Closed v1 vocabulary:

```text
PROJECTABLE
PROJECTED_AS_HYPOTHESIS
PARTIALLY_PROJECTABLE
UNREPRESENTABLE
REJECTED
REQUEST_REVISION
```

Meanings:

```text
PROJECTABLE
    full projection is mechanically possible, but the proposal requested validation-only resolution

PROJECTED_AS_HYPOTHESIS
    full projection was executed into the frozen Phys hypothesis path

PARTIALLY_PROJECTABLE
    projection-class result is PARTIALLY_REPRESENTABLE under a future permitted-loss contract;
    unreachable in MVP

UNREPRESENTABLE
    projection cannot satisfy the strict mirror contract

REJECTED
    proposal fails protocol validation for a reason unrelated to representation alone

REQUEST_REVISION
    proposal is syntactically valid enough to identify a concrete correction needed
```

## 23.4 Proposal state/output relationship

Persisted proposal records always contain the compound tuple:

```text
(State, Outcome)
```

Allowed combinations:

```text
DRAFT      -> Outcome absent
SUBMITTED  -> Outcome absent
RESOLVED   -> exactly one Outcome
```

`RESOLVED` means protocol processing is complete, not scientific success.

## 23.5 Revision lineage

When `REQUEST_REVISION` is returned:

```text
parent -> (RESOLVED, REQUEST_REVISION)
child  -> new proposal
child.ParentProposalRef = parent.ProposalRef
child.ProposalRevision = parent.ProposalRevision + 1
```

A resolved proposal is immutable.

A proposal in `RESOLVED` state MUST NOT be resubmitted or mutated.

## 23.6 FormalizationResponse

```text
FormalizationResponse {
    ProposalRef         ArtifactRef
    ProposalRevision    UInt32
    Outcome             ProposalOutcome
    ProjectionClass     Optional<ProjectionClass>
    ProjectedArtifactRefs []ArtifactRef
    LossReportRef       Optional<ArtifactRef>
    ValidationReportRef Optional<ArtifactRef>
}
```

Collection:

```text
ProjectedArtifactRefs = SET
```

`ProjectionClass` presence rule:

```text
PROJECTABLE              -> FULLY_REPRESENTABLE
PROJECTED_AS_HYPOTHESIS  -> FULLY_REPRESENTABLE
PARTIALLY_PROJECTABLE    -> PARTIALLY_REPRESENTABLE
UNREPRESENTABLE          -> UNREPRESENTABLE
REJECTED                 -> absent
REQUEST_REVISION         -> absent
```

No other mapping is valid.

## 23.7 Formalization implementation ownership

The protocol is implemented by Project A integration tooling and/or a future mechanical interop layer.

The tooling may orchestrate existing Phys public APIs.

It MUST NOT add a new Phys API merely to satisfy this architecture.

It MUST NOT mutate the Phys corpus through an invented acceptance path.

## 23.8 Known Phys hypothesis mint path

The current architecture recognizes the frozen Phys hypothesis-grade raw mint path:

```text
hypothesis.NewCandidateConcept
```

but the implementation plan MUST verify the exact signature and availability against the frozen Phys revision before coding.

If the supplied frozen revision differs, implementation MUST stop and report a contract mismatch rather than invent a replacement API.

## 23.9 Provenance law

When external content is projected through the frozen Phys hypothesis path, the resulting Phys object is hypothesis-grade.

A session correspondence using hypothesis-grade operands remains hypothesis-grade under the frozen Phys contamination law.

The formalization protocol records this outcome and never treats it as an error or an implicit corpus promotion.

The exact inherited contamination rule is identified in §0.1's frozen Phys authority; it is not redefined here.

---

# 24. Trust and corpus boundary

## 24.1 Current Phys trust model

Existing Phys provenance and corpus status remain entirely within Phys.

This architecture does not introduce a standalone runtime `TrustAssertion` type into Phys.

## 24.2 TrustAssertion status

```text
TrustAssertion = FUTURE DESIGN TARGET
```

No coding agent may implement:

```text
AcceptTrust
AcceptFormalization
PromoteCorpus
```

as a new generic authority path in Phys or this foundation.

## 24.3 External physical objects

A `phys-*` object remains external unless and until an explicitly supported formalization path produces a Phys artifact through existing frozen APIs.

Even then, successful projection is hypothesis-grade unless an existing frozen Phys curation path independently changes corpus status under human authority.

---

# 25. Numerical boundary

## 25.1 Principle

Numerical consumes an executable numerical contract.

It does not reconstruct symbolic mathematics or physical semantics.

## 25.2 ComputationalContract owner

The `ComputationalContract` schema is owned by the Numerical Layer's **contract subpackage**:

```text
numerical/contract
```

The execution engine consumes that contract.

`phys-*` packages produce contract instances.

This is dual ownership by role:

```text
schema ownership -> Numerical/contract
instance production -> phys-*
execution          -> Numerical
```

## 25.3 ComputationalContract schema

```text
ComputationalContract {
    SchemaVersion       string
    PhysicalModelRefs   []ArtifactRef
    MathematicalModelRefs []ArtifactRef
    Inputs              []NumericInput
    Parameters          []NumericParameter
    RequestedOutputs    []RequestedOutput
    NumericalMethod     NumericalMethodID
    Precision           PrecisionRequirement
    ExecutionLimits     ExecutionLimits
    ProvenanceRefs      []ArtifactRef
}
```

Collection rules:

```text
PhysicalModelRefs     = SET
MathematicalModelRefs = SET
Inputs                = SET
Parameters            = SET
RequestedOutputs      = SET
ProvenanceRefs        = SET
```

Each Set element carries its own name/key where positional distinction is required.

## 25.4 NumericInput

```text
NumericInput {
    Name      string
    Shape     []UInt32
    ValuesRef ArtifactRef
}
```

`Shape` is a `SEQUENCE`.

`ValuesRef` resolves to a numerical data artifact owned by the execution environment.

## 25.5 NumericParameter

```text
NumericParameter {
    Name  string
    Value NumericValue
}
```

## 25.6 NumericValue

Closed v1 variants:

```text
RATIONAL
FLOAT64_BITS
INTEGER
TEXT
```

`FLOAT64_BITS` uses §5.7.

`RATIONAL` and `INTEGER` use exact §5 encodings.

`TEXT` is permitted only for parameters explicitly declared as non-numeric configuration data.

## 25.7 RequestedOutput

```text
RequestedOutput {
    Name string
}
```

`Name` is identity within the contract's RequestedOutputs Set.

## 25.8 NumericalMethodID

A closed v1 architectural schema registry contains:

```text
DIRECT
ITERATIVE
LINEAR_SYSTEM
USER_PROVIDED_EXECUTION_PLAN
```

These values classify execution strategy; they do not describe physical validity.

## 25.9 PrecisionRequirement

```text
PrecisionRequirement {
    AbsoluteTolerance Optional<Float64Bits>
    RelativeTolerance Optional<Float64Bits>
    Digits             Optional<UInt32>
}
```

All fields absent is valid and means:

```text
no additional precision constraint declared by the physical contract
```

Numerical may choose an implementation default; it MUST NOT reinterpret that default as a physical claim.

## 25.10 ExecutionLimits

```text
ExecutionLimits {
    MaxSteps Optional<UInt64>
    MaxBytes Optional<UInt64>
    MaxSeconds Optional<UInt64>
}
```

All fields absent is valid and means no additional contract-level limit is declared.

## 25.11 Constant-resolution boundary

Before a `ComputationalContract` is constructed:

```text
phys-* / contract owner
    resolves exact mathematical and physical constants
    into NumericValue
```

Then:

```text
ComputationalContract
    -> Numerical
```

Numerical never resolves symbolic constants.

## 25.12 In-process handoff

When the producer and Numerical execute within one process, the handoff MUST contain fully materialized in-memory numeric arrays/execution structures.

Numerical MUST NOT perform ArtifactRef resolution as part of the numerical inner loop.

Persisted ArtifactRefs may be retained only for provenance and external reload workflows.

## 25.13 Numerical dependency rule

Numerical MUST NOT import:

```text
PhysMath
phys-* implementations
```

The Numerical engine receives formalized/executable data, not symbolic ASTs.

---

# 26. Status vocabularies

Status vocabularies are owner-specific.

## 26.1 PhysMath operation outcomes

```text
SUCCESS
UNDEFINED
TRANSFORMATION_FAILED
OUT_OF_SCOPE
RESOURCE_LIMIT
```

## 26.2 Project A research states

```text
PROPOSED
UNDER_ATTACK
REVISION_REQUIRED
COUNTEREXAMPLE_FOUND
INCONSISTENT
AWAITING_FORMALIZATION
FORMALIZED
CLOSED
```

## 26.3 Formalization outcomes

```text
PROJECTABLE
PROJECTED_AS_HYPOTHESIS
PARTIALLY_PROJECTABLE
UNREPRESENTABLE
REJECTED
REQUEST_REVISION
```

## 26.4 Projection classes

```text
FULLY_REPRESENTABLE
PARTIALLY_REPRESENTABLE
UNREPRESENTABLE
```

## 26.5 Phys corpus status

Defined only by the frozen Phys specification.

This architecture does not re-declare it.

## 26.6 Meaning separation

For example:

```text
FormalizationOutcome.UNREPRESENTABLE
```

and:

```text
ProjectionClass.UNREPRESENTABLE
```

are different typed values with different owners, even though their canonical strings are identical.

`UNREPRESENTABLE` MUST NOT be used as a generic free-form status outside its declared type.

---

# 27. phys-* extension and conformance contract

## 27.1 Required conformance

Every `phys-*` package MUST demonstrate:

```text
1. namespace ownership
2. concrete artifact schemas
3. ArtifactRef compatibility
4. canonical serialization determinism
5. explicit PhysicalInterpretation ownership
6. no duplicate trust/corpus authority
7. explicit MathRefs
8. explicit framework identity
9. explicit projection/containment behavior
```

## 27.2 No duplicate semantic metadata

A `phys-*` package MUST NOT create a second system claiming to be:

```text
Phys provenance
Phys corpus status
Phys trust authority
Phys assumption authority
```

Domain-specific physical metadata is allowed when it is clearly package-owned physical semantics rather than a competing trust authority.

## 27.3 Physical interpretation conflicts

Multiple frameworks MAY publish different interpretations for the same MathArtifact.

There is no global winner.

Consumers requiring physical meaning MUST name the required `FrameworkRef` or reject an ambiguous resolution.

## 27.4 Theory-local algebra conformance

Theory-local mathematical types remain local.

A mapping into PhysMath is explicit and deterministic.

A dynamic runtime AST extension registry is forbidden.

## 27.5 Human admission

Admission of a new `phys-*` package is a human governance decision.

The conformance suite checks structure; it does not adjudicate scientific truth.

---

# 28. phys-interop decision

## 28.1 MVP

There is no mandatory shared runtime `phys-interop` module.

The neutral identity module is:

```text
phys-artifact
```

and has a strictly narrower purpose.

## 28.2 Future interop module

A future `phys-interop` MAY provide only:

```text
reference codecs
schema utilities
cross-project conformance tests
import-boundary checks
mechanical projection helpers
```

It MUST NOT provide:

```text
physics ontology
math ontology
trust authority
corpus authority
interpretation policy
research state
scientific adjudication
```

## 28.3 Duplicate implementation trade-off

Until a shared interop runtime is justified by a second genuine consumer, each producer MAY implement mechanically equivalent codecs locally.

The conformance contract is normative, so equivalence is testable.

---

# 29. Version compatibility and evolution

## 29.1 Schema compatibility

A consumer MUST:

```text
accept exact supported schema versions
reject incompatible versions
never silently reinterpret unknown versions
never silently ignore unknown required fields
```

## 29.2 Operation compatibility

`OperationSchemaVersion` determines compatibility of transformation parameter bytes and replay behavior.

## 29.3 Phys compatibility

Frozen Phys compatibility is determined by the exact supplied frozen revision.

The architecture does not invent a compatibility shim for incompatible Phys revisions.

## 29.4 Future PhysMath sort extension

A new sort is introduced only through an explicit versioned schema change.

It is never injected dynamically into a running v1 evaluator.

---

# 30. MVP / post-MVP boundary and Growth Gate

## 30.1 MVP

The frozen Phys MVP remains unchanged.

The new ecosystem projects are target architecture, not an instruction to modify the Phys repository.

## 30.2 Growth Gate

The Growth Gate is **post-MVP**.

During MVP:

```text
Phys kernel frozen
GR stress-test evidence recorded externally
no kernel modification
no automatic Growth Gate execution
```

Post-MVP:

```text
evidence from phys-* projects
    -> Growth Gate evaluation
    -> possible future Phys kernel evolution
```

GR stress-test success or failure does not itself authorize kernel changes.

## 30.3 Growth routing

When a capability is missing:

```text
pure mathematical capability
    -> PhysMath

physical-domain semantic capability
    -> phys-*

trusted generic physical primitive truly required by multiple domains
    -> Phys Growth Gate

numerical execution capability
    -> Numerical

research hypothesis/correspondence
    -> Project A
```

A TOE difficulty alone is not evidence for kernel growth.

---

# 31. TOE foundation boundary

## 31.1 Dancer / Dance / Stage status

The following are research-level conceptual labels, not foundation ontology:

```text
Dancer
Dance
Stage
```

Their current conceptual use is:

```text
Dancer -> matter/state sector
Dance  -> guidance/dynamical/RG flow
Stage  -> causal/discrete spacetime sector
```

CST and CDT remain distinct:

```text
CST -> causal-set substrate
CDT -> dynamical triangulated-spacetime realization
```

These mappings are not encoded as PhysMath types.

## 31.2 Research goal

The TOE-Agent may use this architecture to test candidate mappings of:

```text
physical concept
    -> mathematical object
    -> valid relation
    -> regime/conditions
    -> epistemic status
```

The architecture rejects category errors such as:

```text
spectral dimension == proof of asymptotic safety
index theorem == universal matter/geometry coupling
k == number of causal-set elements
CDT == purely kinematic substrate
```

unless an explicit domain contract actually defines such a relation.

## 31.3 Mathematical Reality Check

For every candidate cross-domain claim:

```text
physical object has a declared physical type
mathematical object has a declared mathematical type
mapping direction is explicit
regime is explicit
conditions are explicit
epistemic status is explicit
```

No silent interpretation.

---

# 32. Artifact contract template and no-invention matrix

## 32.1 Artifact contract template

Every durable artifact family MUST specify:

```text
OWNER
SCHEMA
IDENTITY
CANONICAL FORM
SERIALIZATION
REFERENCES
RESOLUTION
VALIDATION
LEGAL STATES
LEGAL TRANSITIONS
POSITIVE TEST
NEGATIVE TEST
```

### Resolution rule

Resolution is environment infrastructure.

A semantic artifact MUST NOT require a network/database lookup merely to validate its own structure.

### Identity rule

Identity is the §4 ArtifactRef tuple.

### Legal-state rule

If the artifact has lifecycle state, it references the owning state machine in §22 or §23.

### Test rule

Every durable artifact family has at least one executable positive and one executable negative test in Given/When/Then form.

## 32.2 Complete No-Invention Matrix

The matrix is complete only when every row below contains:

```text
Normative decision
Owner
Concrete schema/type
Validation rule
Acceptance test
Source ADR
```

Required rows:

```text
dependency direction
Phys freeze
phys-artifact ownership
ArtifactRef ownership
ArtifactRef identity
canonical digest framing
embedded ArtifactRef encoding
namespace ownership
in-memory composition
serialized composition
Phys referenceable artifact family
PhysMath sort set
Expression node set
ParameterID
BoundParameterRef
ConstantID
MathCondition
TransformResult
TransformationRuleID
ContractRuleRef
OperatorApplication parameters
differentiate identity
OperationSchemaVersion
projection subset
projection composition
projection loss policy
ProjectionClass
ProjectionReasonCode
NodeKind
PhysicalInterpretation
InterpretationClass
FrameworkRef
FrameworkArtifact
PhysicalContract
ProposedInterpretation
CorrespondenceProposal
ResearchHypothesis
ProposalState
ProposalOutcome
revision lineage
provenance inheritance
TrustAssertion future status
MathTrace
replay authority
ComputationalContract owner
ComputationalContract schema
Numerical dependency topology
constant-resolution boundary
phys-* conformance
phys-interop scope
version compatibility
MVP boundary
Growth Gate phase
TOE boundary
status vocabularies
supersession
review arbitration
```

A missing row means the architecture is open.

---

# 33. No-Go architectures

The following are explicitly forbidden:

```text
1. PhysMath -> Phys dependency
2. Phys -> PhysMath dependency
3. Phys modification merely to support this architecture
4. universal Artifact interface/envelope
5. central mandatory runtime artifact registry
6. ArtifactRef interpreted as a resolver
7. mutable AST nodes carrying hidden conditions
8. dynamic PhysMath AST extension registry
9. physical-role enums embedded in PhysMath core
10. generic theorem prover inside PhysMath
11. Numerical parsing PhysMath symbolic ASTs
12. Numerical resolving physical semantics
13. Numerical resolving symbolic constants
14. generic CompatibilityRules engine
15. generic trust/corpus authority outside Phys
16. runtime AcceptFormalization/PromoteCorpus API
17. automatic scientific truth adjudication
18. implicit physical interpretation of mathematical structure
19. theory-local algebra injected into PhysMath at runtime
20. global replay engine pretending to replay all authorities
21. artifact identity based on UUID/assigned IDs in addition to digest
22. float formatting used as canonical identity encoding
23. arbitrary open-ended Transformation parameter bags
24. RuleIDs used as unexplained prose labels
25. one status enum shared across research/projection/formalization/corpus domains
26. hidden merge semantics for ParameterID resolution
27. silent projection repair
28. treating a successful projection as scientific truth
```

---

# 34. Executable acceptance properties

The acceptance suite is part of the architecture.

## 34.1 Dependency

**Given** the production package graph  
**When** the dependency audit runs  
**Then** PhysMath has no Phys import, Phys has no PhysMath import, Phys has no new architecture dependency, and Numerical has no phys-* implementation imports.

Failure condition: any prohibited production import exists.

## 34.2 Artifact identity

**Given** identical canonical payload bytes and family identity  
**When** two implementations compute an ArtifactRef  
**Then** the full four-field ArtifactRef is identical.

Failure condition: differing digest or differing family tuple.

## 34.3 Encoding order

**Given** an ArtifactRef with fields populated by distinct values  
**When** canonical bytes are generated  
**Then** field order is Namespace, SchemaVersion, ArtifactType, Digest.

Failure condition: ArtifactType precedes SchemaVersion or any alternative encoding is accepted.

## 34.4 Embedded reference

**Given** a parent artifact containing ArtifactRef fields  
**When** its canonical payload is encoded  
**Then** each embedded ref uses §5.3 exactly.

Failure condition: JSON/hex/raw digest representation occurs in canonical bytes.

## 34.5 UTF-8 length

**Given** a non-ASCII string  
**When** its canonical string encoding is generated  
**Then** Len32 equals UTF-8 byte length.

Failure condition: rune count or locale-specific count is used.

## 34.6 Parameter identity

**Given** two declarations with the same structural scope path and parameter ordinal  
**When** ParameterID is computed independently  
**Then** the resulting 32-byte IDs match.

Failure condition: random or display-name-based IDs differ.

## 34.7 Alpha-renaming

**Given** a parameter display name changes from `t` to `tau`  
**When** the expression is re-canonicalized  
**Then** ParameterID is unchanged.

Failure condition: display rename changes parameter identity.

## 34.8 Cross-artifact parameter resolution

**Given** the same local ParameterID occurs in two declaration artifacts  
**When** a BoundParameterRef resolves the target  
**Then** the DeclaringArtifactRef determines the declaration.

Failure condition: display name or bare ParameterID silently selects a target.

## 34.9 Symbol validity

**Given** a durable Symbol  
**When** validation runs  
**Then** exactly one of ParameterID or ConstantID is present.

Failure condition: neither or both is accepted.

## 34.10 ConstantID

**Given** `math.const.pi`  
**When** encoded  
**Then** it uses ConstantID string encoding and the reserved registry.

Failure condition: arbitrary digest/UUID representation is accepted.

## 34.11 Add/Mul empty arity

**Given** Add([]) or Mul([])  
**When** canonicalization runs  
**Then** they become 0 or 1 respectively.

Failure condition: empty associative node survives canonicalization.

## 34.12 Add/Mul unary arity

**Given** Add([x]) or Mul([x])  
**When** canonicalization runs  
**Then** output is x.

Failure condition: binary identity wrapper is invented.

## 34.13 Condition determinism

**Given** duplicate and differently ordered generated conditions  
**When** TransformResult is finalized  
**Then** conditions are sorted by digest and deduplicated.

Failure condition: traversal order changes the serialized condition sequence.

## 34.14 Differentiate schema

**Given** an OperatorApplication with OperatorID `differentiate`  
**When** validation runs  
**Then** Parameters contain BoundParameterRef + UInt32 Order and OperatorRef is absent.

Failure condition: target/order have no canonical field or are represented only as free-form args.

## 34.15 Differentiate determinism

**Given** the same declaration context, target, expression, and order  
**When** two implementations execute DIFFERENTIATE-001  
**Then** their canonical parameters, trace, and output are identical.

Failure condition: operator identity or parameter bytes differ.

## 34.16 RuleID closure

**Given** a Transformation  
**When** its RuleID is validated  
**Then** it belongs to §15.3 exactly.

Failure condition: an unregistered RuleID is accepted.

## 34.17 Rule parameter closure

**Given** each RuleID in §15.3  
**When** its Parameters are validated  
**Then** its exact schema from §15.5 is enforced.

Failure condition: open-ended parameter bags or ambiguous field order are accepted.

## 34.18 Projection classification

**Given** a composite expression with one unrepresentable required child  
**When** projection runs  
**Then** the parent is UNREPRESENTABLE in MVP.

Failure condition: parent becomes PARTIAL merely because another child projected.

## 34.19 Symbol collision

**Given** two Symbol nodes with equal DisplayName but different ParameterID/ConstantID  
**When** projection runs  
**Then** it fails with SYMBOL_BINDING_COLLISION.

Failure condition: distinct source semantics merge silently.

## 34.20 Pow projection

**Given** Pow(base, symbolic_exponent)  
**When** strict projection runs  
**Then** classification is UNREPRESENTABLE with NON_RATIONAL_EXPONENT_IN_POW.

Failure condition: approximation or silent exponent coercion.

## 34.21 Associative projection

**Given** Add([a,b,c])  
**When** projected  
**Then** Phys receives the canonical left-associated binary tree.

Failure condition: right association or traversal-order association.

## 34.22 Multiple interpretations

**Given** one mathematical artifact  
**When** two `phys-*` frameworks publish interpretations  
**Then** both may coexist.

Failure condition: one is silently chosen as a global winner.

## 34.23 Formalization no-promotion

**Given** a valid projection from an external hypothesis  
**When** formalization succeeds  
**Then** the Phys result is hypothesis-grade and no corpus mutation occurs.

Failure condition: ESTABLISHED/IDENTIFIED appears solely because projection succeeded.

## 34.24 Proposal revision

**Given** a `REQUEST_REVISION` response  
**When** lineage is persisted  
**Then** parent is RESOLVED and child points to it via ParentProposalRef with revision +1.

Failure condition: parent remains SUBMITTED or is mutated in place.

## 34.25 ProjectionClass mapping

**Given** a FormalizationResponse outcome  
**When** it is validated  
**Then** ProjectionClass follows §23.6 exactly.

Failure condition: REJECTED or REQUEST_REVISION carries a fabricated projection class.

## 34.26 Numerical boundary

**Given** a ComputationalContract handed to Numerical in-process  
**When** numerical execution begins  
**Then** all required numeric data is already materialized in memory.

Failure condition: Numerical performs implicit ArtifactRef/network lookup.

## 34.27 Constant resolution boundary

**Given** a symbolic constant in a physical model  
**When** ComputationalContract is built  
**Then** the producer has already converted it to a NumericValue.

Failure condition: Numerical resolves the symbol itself.

## 34.28 Numerical dependency

**Given** the Numerical engine package graph  
**When** the import audit runs  
**Then** Numerical imports contract + mechanical identity dependencies but no phys-* implementation.

Failure condition: semantic reverse dependency appears.

## 34.29 No-Go physical inference

**Given** a pure mathematical artifact with no physical interpretation  
**When** an API consumes it  
**Then** no physical type is inferred from the mathematical sort alone.

Failure condition: mathematical structure silently acquires physical meaning.

## 34.30 Phys freeze

**Given** the implementation workspace  
**When** the architecture implementation gate runs  
**Then** the Phys repository remains at the supplied frozen revision with no architecture-driven changes.

Failure condition: new Phys imports, APIs, or files appear because of this architecture.

---

# 35. Supersession register

The following v1.5 decisions supersede the corresponding earlier architecture choices:

| v1.5 decision | Supersedes | Reason |
|---|---|---|
| Neutral `phys-artifact` exists as mechanical identity library | v1.4 provisional identity ownership | Gives one shared ArtifactRef definition without making Phys a consumer |
| Phys does not import `phys-artifact` | v1.4 dependency diagram | Preserves frozen Phys byte/dependency freeze |
| ArtifactRef is a name, not resolver | earlier registry-oriented interpretations | Separates naming from storage infrastructure |
| Canonical digest uses one exact framed formula | earlier abstract formula variants | Removes byte-level ambiguity |
| ArtifactRef struct order matches wire order | v1.4 declaration-order warning | Prevents encoder trap |
| ParameterID uses structural/positional identity | earlier global/cross-artifact identity proposals | Avoids circular artifact-reference identity |
| BoundParameterRef carries declaration context | earlier bare ParameterID cross-artifact use | Removes ambiguity |
| ConstantID = string with closed registry | v1.4 undefined ConstantID | Makes Symbol identity portable |
| `differentiate` is a reserved OperatorID | v1.4 undefined canonical Operator artifact | Removes bootstrap and artifact-identity ambiguity |
| OperatorApplication gains Parameters | v1.4 Args-only schema | Gives differential parameters a declared home |
| TransformationRuleID differs from ContractRuleRef | v1.4 Appendix-B/C dual use | Removes semantic registry confusion |
| `ID-001` removed as TransformationRuleID | v1.4 serialization-rule collision | Serialization is not a mathematical transform |
| OperationSchemaVersion = operation parameter/replay schema | v1.4 undefined field semantics | Makes trace portability explicit |
| ProjectionClass optional by outcome | v1.4 required-on-rejection ambiguity | Non-projecting outcomes carry no invented class |
| `InterpretationClass` closed to two values | v1.4 open string | Makes physical interpretation filtering deterministic |
| FrameworkRef -> `<namespace>/framework/1` | v1.4 target ambiguity | Gives framework identity a concrete target family |
| Numerical contract schema owned by Numerical/contract | v1.4 owner ambiguity | Keeps schema ownership at execution boundary without circular import |
| Numerical engine does not import phys-* | earlier consumer-direction model | Prevents producer/consumer cycle |
| Constant resolution upstream of ComputationalContract | v1.4 §13.3/§24.3 ambiguity | One unambiguous phase boundary |
| MVP permitted-loss set = empty | v1.3/v1.4 partial-projection ambiguity | Makes PARTIAL unreachable until explicitly enabled |
| MathCondition finalization = sort + dedup | earlier sort-only rule | Prevents duplicate canonical forms |
| TransformResult ephemeral only | earlier artifact-like ambiguity | Separates operation results from durable artifacts |
| Growth Gate post-MVP | earlier stress-test growth ambiguity | Preserves frozen Phys |
| Physical meaning never inferred | earlier philosophy-only wording | Promotes core boundary to hard prohibition |
```

Lineage tags for future supersession records MUST distinguish:

```text
INHERITED_FROM_FROZEN_PHYS
INHERITED_FROM_PRIOR_ARCHITECTURE
DECIDED_IN_V1_5
```

---

# 36. Frozen vocabulary registry

## 36.1 Artifact families — PhysMath core

```text
physmath/expression/1
physmath/function/1
physmath/function_application/1
physmath/functional/1
physmath/functional_application/1
physmath/operator/1          # reserved, DEFERRED operationally
physmath/operator_application/1
physmath/tensor/1
physmath/index/1
physmath/map/1
physmath/relation/1
physmath/collection/1
physmath/transformation/1
physmath/limit/1
physmath/asymptotic_truncation/1
physmath/condition/1         # support artifact
```

`DerivationObligation` is reserved but has no v1 durable artifact family because it is DEFERRED + CONTAINMENT.

## 36.2 NodeKind registry

```text
Symbol
Rational
Add
Mul
Neg
Pow
Sqrt
Call
Relation
BranchSet
```

## 36.3 TransformationRuleID registry

```text
EXPR-CANON-001
EXPR-SIMPLIFY-001
EXPR-SUBSTITUTE-001
RAT-NORMALIZE-001
DIFFERENTIATE-001
SELECT-BRANCH-001
PROJECT-STRICT-001
```

## 36.4 ProjectionReasonCode registry

```text
NON_RATIONAL_EXPONENT_IN_POW
SYMBOL_BINDING_COLLISION
UNSUPPORTED_NODE_KIND
UNSUPPORTED_TOP_LEVEL_SORT
UNRESOLVED_FUNCTION_REFERENCE
LOSS_NOT_PERMITTED
INVALID_SOURCE_ARTIFACT
```

## 36.5 InterpretationClass registry

```text
INTERPRETATION_NEUTRAL
FRAMEWORK_SPECIFIC
```

## 36.6 ProposalState registry

```text
DRAFT
SUBMITTED
RESOLVED
```

## 36.7 ProposalOutcome registry

```text
PROJECTABLE
PROJECTED_AS_HYPOTHESIS
PARTIALLY_PROJECTABLE
UNREPRESENTABLE
REJECTED
REQUEST_REVISION
```

## 36.8 OperationOutcome registry

```text
SUCCESS
UNDEFINED
TRANSFORMATION_FAILED
OUT_OF_SCOPE
RESOURCE_LIMIT
```

---

# 37. Artifact-family contract matrix

Legend:

```text
RULE <id> = exact rule named in the artifact-family section
N/A = behavior does not apply
DEFERRED + CONTAINMENT = schema reserved/contained; operation not implemented in v1
```

| Sort | Construct | Canonicalize | Equality | Substitute | Differentiate | Serialize | Project |
|---|---|---|---|---|---|---|---|
| Expression | RULE EXPR-CONSTRUCT-001 | RULE EXPR-CANON-001 | RULE EXPR-EQUAL-001 | RULE EXPR-SUBSTITUTE-001 | RULE EXPR-DIFFERENTIATE-001 | RULE EXPR-SERIALIZE-001 | RULE PROJECT-STRICT-001 |
| Function | RULE FUNCTION-CONSTRUCT-001 | RULE FUNCTION-CANON-001 | RULE FUNCTION-EQUAL-001 | RULE FUNCTION-SUBSTITUTE-001 | DEFERRED + CONTAINMENT | RULE FUNCTION-SERIALIZE-001 | DEFERRED + CONTAINMENT |
| FunctionApplication | RULE CALL-CONSTRUCT-001 | RULE CALL-CANON-001 | RULE CALL-EQUAL-001 | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | RULE CALL-SERIALIZE-001 | DEFERRED + CONTAINMENT |
| Functional | RULE FUNCTIONAL-CONSTRUCT-001 | RULE FUNCTIONAL-CANON-001 | RULE FUNCTIONAL-EQUAL-001 | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | RULE FUNCTIONAL-SERIALIZE-001 | DEFERRED + CONTAINMENT |
| FunctionalApplication | RULE FUNCTIONALAPP-CONSTRUCT-001 | RULE FUNCTIONALAPP-CANON-001 | RULE FUNCTIONALAPP-EQUAL-001 | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | RULE FUNCTIONALAPP-SERIALIZE-001 | DEFERRED + CONTAINMENT |
| Operator | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT |
| OperatorApplication | RULE OPAPP-CONSTRUCT-001 | RULE OPAPP-CANON-001 | RULE OPAPP-EQUAL-001 | DEFERRED + CONTAINMENT | RULE DIFFERENTIATE-001 | RULE OPAPP-SERIALIZE-001 | DEFERRED + CONTAINMENT |
| Tensor | RULE TENSOR-CONSTRUCT-001 | RULE TENSOR-CANON-001 | RULE TENSOR-EQUAL-001 | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | RULE TENSOR-SERIALIZE-001 | DEFERRED + CONTAINMENT |
| Index | RULE INDEX-CONSTRUCT-001 | RULE INDEX-CANON-001 | RULE INDEX-EQUAL-001 | DEFERRED + CONTAINMENT | N/A | RULE INDEX-SERIALIZE-001 | DEFERRED + CONTAINMENT |
| Map | RULE MAP-CONSTRUCT-001 | RULE MAP-CANON-001 | RULE MAP-EQUAL-001 | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | RULE MAP-SERIALIZE-001 | DEFERRED + CONTAINMENT |
| Relation | RULE RELATION-CONSTRUCT-001 | RULE RELATION-CANON-001 | RULE RELATION-EQUAL-001 | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | RULE RELATION-SERIALIZE-001 | DEFERRED + CONTAINMENT |
| Collection | RULE COLLECTION-CONSTRUCT-001 | RULE COLLECTION-CANON-001 | RULE COLLECTION-EQUAL-001 | DEFERRED + CONTAINMENT | N/A | RULE COLLECTION-SERIALIZE-001 | DEFERRED + CONTAINMENT |
| Transformation | RULE TRANSFORMATION-CONSTRUCT-001 | RULE TRANSFORMATION-CANON-001 | RULE TRANSFORMATION-EQUAL-001 | N/A | N/A | RULE TRANSFORMATION-SERIALIZE-001 | N/A |
| Limit | RULE LIMIT-CONSTRUCT-001 | RULE LIMIT-CANON-001 | RULE LIMIT-EQUAL-001 | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | RULE LIMIT-SERIALIZE-001 | DEFERRED + CONTAINMENT |
| AsymptoticTruncation | RULE ASYMPTOTIC-CONSTRUCT-001 | RULE ASYMPTOTIC-CANON-001 | RULE ASYMPTOTIC-EQUAL-001 | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | RULE ASYMPTOTIC-SERIALIZE-001 | DEFERRED + CONTAINMENT |
| DerivationObligation | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT | DEFERRED + CONTAINMENT |

### 37.1 ContractRuleRef rule

The `RULE ...` identifiers in this table are not serialized TransformationRuleID values unless separately listed in §36.3.

They identify the normative contract paragraphs for the corresponding sort.

---

# 38. Deterministic canonicalization examples

## 38.1 Add ordering

For source children:

```text
b, a, c
```

the canonical child order is the lexicographic order of each child's canonical payload bytes.

No source traversal order survives.

## 38.2 Condition sorting

For conditions `C1`, `C2`, `C3`, canonical order is:

```text
Digest(C1) < Digest(C2) < Digest(C3)
```

and duplicate digest entries are reduced to one.

## 38.3 ArtifactRef Set ordering

A set containing refs is ordered by the exact §5.3 canonical bytes.

This resolves the prior field-order ambiguity between embedded encoding and Set ordering.

---

# 39. Resource bounds and failure behavior

## 39.1 Architectural requirement

Operations MUST have deterministic resource behavior.

The architecture does not prescribe machine-specific default numbers because those are deployment parameters, not semantic decisions.

A runtime MUST expose explicit resource limits to the operation evaluator.

## 39.2 Failure rule

A resource exhaustion produces:

```text
OperationOutcome = RESOURCE_LIMIT
```

It does not produce a partial mathematical artifact and must not be interpreted as mathematical falsehood.

## 39.3 No hidden approximation

Resource exhaustion MUST NOT trigger an undocumented approximation or numerical fallback.

---

# 40. Research meta-ledger

## 40.1 Purpose

Project A MAY maintain a cross-authority research ledger.

It is a graph of artifact transitions, not a global authority replay log.

## 40.2 ResearchStep

```text
ResearchStep {
    Operation    string
    InputRefs    []ArtifactRef
    OutputRef    ArtifactRef
    ModuleOrigin string
}
```

`InputRefs` is a SEQUENCE.

The ledger records:

```text
what artifact was used
what artifact was produced
which authority produced it
```

It does not claim to reproduce internal authority operations.

## 40.3 Research replay

Replaying the research ledger means validating that referenced artifacts and declared transitions form a coherent research graph.

It does not replace authority-local replay from §16.

---

# 41. Numerical and research separation

The Numerical Layer consumes formalized execution contracts.

It does not consume:

```text
ResearchHypothesis directly
ProposedInterpretation directly
CorrespondenceProposal directly
unprojected PhysMath ASTs directly
```

A producer MUST first derive a valid ComputationalContract.

This prevents “research state” from accidentally becoming executable physics.

---

# 42. Implementation boundary discipline

The architecture deliberately leaves the following as implementation details:

```text
private Go helper names
private package layout beneath the declared public boundary
memory allocation strategy
caching
parallelism
storage backend
network transport
UI
logging format
performance optimization
```

These choices MUST NOT alter any normative identity, schema, dependency, or semantic rule in this document.

---

# 43. Golden-vector requirements

The implementation MUST provide machine-readable golden vectors for at least:

```text
ArtifactRef encoding
Artifact digest framing
BigInt
Rational
ParameterID
BoundParameterRef
ConstantID
Symbol
Add canonicalization
Condition sorting/dedup
DIFFERENTIATE-001 parameters
FormalizationResponse outcome/class mapping
ProjectionLossReport reason ordering
ComputationalContract numeric value encoding
```

A golden vector contains:

```text
input logical value
canonical bytes hex
expected digest where applicable
```

The golden vectors are normative test fixtures once frozen with the implementation plan.

---

# 44. Acceptance inventory by boundary

## 44.1 PhysMath

```text
closed node kinds
closed v1 sort set
exact canonicalization
ParameterID determinism
ConstantID determinism
MathCondition accumulation
closed TransformationRuleID registry
bounded differentiation
no physical-role semantics
```

## 44.2 phys-*

```text
framework identity
interpretation ownership
explicit math references
no trust duplication
explicit projection behavior
local theory algebra containment
```

## 44.3 Phys

```text
no new dependency
no new API
existing hypothesis mint only
frozen provenance contamination law
frozen corpus authority
```

## 44.4 Project A

```text
research states
formalization protocol
proposal lineage
cross-authority meta-ledger
no trust promotion
```

## 44.5 Numerical

```text
contract schema
materialized in-process inputs
constant resolution upstream
no symbolic AST parsing
no physical semantics inference
```

---

# 45. Reviewer arbitration and closure protocol

## 45.1 Standing reviewers

The blocker-only audit is performed by:

```text
DeepSeek reviewer
Gemini reviewer
Qwen reviewer
Z reviewer
Project owner as final human arbiter
```

## 45.2 What qualifies as a blocker

A finding is a **BLOCKER** only when it demonstrates at least one:

```text
1. contradiction between two normative rules;
2. two architecturally legitimate implementations remain;
3. a required schema/type/owner/identity/state/validation rule is absent;
4. the coding agent would have to invent an architectural decision.
```

The reviewer MUST provide:

```text
exact section
concrete counterexample
failure consequence
```

A preference, optimization, naming suggestion, or un demonstrated concern is not a blocker.

## 45.3 Evidence beats majority

A demonstrated blocker blocks freeze even if three other reviewers disagree.

A claimed blocker without the required evidence does not block freeze.

The project owner resolves classification disputes using the evidence test above.

## 45.4 Audit packet

The final audit packet consists only of:

```text
this v1.5 document
frozen Phys authority identified in §0.1
```

Prior review documents remain archival lineage, not additional normative inputs.

## 45.5 Stop condition

Architecture may be frozen when all are true:

```text
DeepSeek: 0 BLOCKER
Gemini:   0 BLOCKER
Qwen:     0 BLOCKER
Z:        0 BLOCKER

No-Invention Matrix complete
Artifact contract templates complete
All closed vocabularies enumerated
All state machines explicit
All identity/framing rules exact
No demonstrated internal contradiction
```

No discretionary fifth architecture round is required after this condition.

---

# 46. Final architectural decision register

| ID | Decision | Status |
|---|---|---|
| ADR-001 | Projects have separate semantic ownership. | CLOSED |
| ADR-002 | PhysMath has no Phys dependency. | CLOSED |
| ADR-003 | Phys has no PhysMath dependency. | CLOSED |
| ADR-004 | phys-* may consume PhysMath and Phys. | CLOSED |
| ADR-005 | Project A owns research state. | CLOSED |
| ADR-006 | PhysMath is pure mathematics. | CLOSED |
| ADR-007 | phys-* owns physical-domain semantics. | CLOSED |
| ADR-008 | Phys owns trusted physical metadata. | CLOSED |
| ADR-009 | Numerical owns numerical execution. | CLOSED |
| ADR-010 | In-memory composition is typed Go composition. | CLOSED |
| ADR-011 | Serialized cross-boundary references use ArtifactRef. | CLOSED |
| ADR-012 | Artifact identity is content-addressed. | CLOSED |
| ADR-013 | ArtifactRef is a name, not a resolver. | CLOSED |
| ADR-014 | Canonical encoding is independent of transport encoding. | CLOSED |
| ADR-015 | `phys-artifact` is neutral identity infrastructure. | CLOSED |
| ADR-016 | Phys remains frozen and does not import `phys-artifact`. | CLOSED |
| ADR-017 | Only `phys/object/1` is ecosystem-referenceable from frozen Phys in MVP. | CLOSED |
| ADR-018 | PhysMath v1 has exactly 16 sorts. | CLOSED |
| ADR-019 | Symbol is an Expression node kind. | CLOSED |
| ADR-020 | PhysMath Expression has exactly 10 node kinds. | CLOSED |
| ADR-021 | ParameterID is deterministic structural/positional identity. | CLOSED |
| ADR-022 | BoundParameterRef supplies cross-artifact declaration context. | CLOSED |
| ADR-023 | Physical parameter roles are external to PhysMath. | CLOSED |
| ADR-024 | ConstantID is a closed mathematical-constant string registry. | CLOSED |
| ADR-025 | Physical constants remain ordinary PhysMath Symbols. | CLOSED |
| ADR-026 | MathCondition is a pure mathematical support artifact. | CLOSED |
| ADR-027 | TransformResult is ephemeral. | CLOSED |
| ADR-028 | TransformationRuleID is a closed serialized operation registry. | CLOSED |
| ADR-029 | ContractRuleRef is distinct from TransformationRuleID. | CLOSED |
| ADR-030 | Serialization is not a TransformationRuleID. | CLOSED |
| ADR-031 | `differentiate` is a reserved OperatorID, not a core Operator artifact. | CLOSED |
| ADR-032 | OperatorApplication has explicit Parameters. | CLOSED |
| ADR-033 | OperationSchemaVersion identifies parameter/replay schema version. | CLOSED |
| ADR-034 | Projection subset is the seven-node frozen-Phys mirror. | CLOSED |
| ADR-035 | MVP permitted-loss set is empty. | CLOSED |
| ADR-036 | Projection is compositional and fails closed. | CLOSED |
| ADR-037 | Symbol binding collisions fail projection. | CLOSED |
| ADR-038 | PhysicalInterpretation is phys-* owned. | CLOSED |
| ADR-039 | InterpretationClass has two v1 values. | CLOSED |
| ADR-040 | FrameworkRef targets `<namespace>/framework/1`. | CLOSED |
| ADR-041 | PhysicalContract is declarative metadata only. | CLOSED |
| ADR-042 | ProposedInterpretation remains Project A owned. | CLOSED |
| ADR-043 | CorrespondenceProposal remains research state. | CLOSED |
| ADR-044 | Multiple framework interpretations may coexist. | CLOSED |
| ADR-045 | Proposal State and Outcome are separate typed concepts. | CLOSED |
| ADR-046 | Formalization revision lineage is immutable and parent-linked. | CLOSED |
| ADR-047 | ProjectionClass is optional when no projection occurs. | CLOSED |
| ADR-048 | Hypothesis-grade formalization does not promote Phys corpus status. | CLOSED |
| ADR-049 | TrustAssertion remains a future design target. | CLOSED |
| ADR-050 | Replay authority is local to the producing authority. | CLOSED |
| ADR-051 | Research replay is a meta-ledger graph, not global authority replay. | CLOSED |
| ADR-052 | ComputationalContract schema is owned by Numerical/contract. | CLOSED |
| ADR-053 | phys-* produces ComputationalContract instances. | CLOSED |
| ADR-054 | Numerical engine has no phys-* implementation dependency. | CLOSED |
| ADR-055 | Exact constant resolution occurs upstream of ComputationalContract construction. | CLOSED |
| ADR-056 | `Float64Bits` is confined to numerical schemas. | CLOSED |
| ADR-057 | phys-* conformance is mandatory. | CLOSED |
| ADR-058 | No shared runtime phys-interop module is required in MVP. | CLOSED |
| ADR-059 | Physical meaning may never be inferred from mathematical structure. | CLOSED |
| ADR-060 | Growth Gate is post-MVP. | CLOSED |
| ADR-061 | TOE Dancer/Dance/Stage are research concepts, not foundation ontology. | CLOSED |
| ADR-062 | v1.5 is the final architecture specification before implementation planning. | CLOSED |

---

# 47. Final closure statement

The architecture is considered **complete for implementation planning** when the blocker-only audit in §45 passes and no demonstrated blocker remains.

The implementation plan that follows this document MUST treat the decisions above as fixed inputs.

The implementation plan MUST NOT reopen:

```text
dependency direction
artifact identity
canonical framing
namespace ownership
Phys freeze
PhysMath sort set
ParameterID
ConstantID
differentiate representation
RuleID semantics
projection semantics
interpretation ownership
formalization states
provenance law
replay authority
Numerical contract ownership
constant-resolution boundary
Growth Gate timing
```

Those are architecture, not implementation choices.

The remaining implementation-plan work is deliberately mechanical:

```text
repository/package sequence
file changes
public API placement consistent with this contract
unit tests
golden vectors
import checks
serialization tests
projection tests
formalization tests
CI gates
acceptance execution
```

**This document is not the TOE theory. It is the architectural substrate that prevents mathematics, physics, trust, research state, and numerical execution from being silently conflated while the TOE research program evolves.**
