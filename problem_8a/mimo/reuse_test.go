package mimo

import (
	"math/big"
	"testing"

	"github.com/PithomLabs/phys-artifact"
	physmath "github.com/PithomLabs/phys-math"

	"github.com/PithomLabs/phys-demo/problem_8a/mimo/axion"
	"github.com/PithomLabs/phys-demo/problem_8a/mimo/representation"
)

// This file is the executable evidence behind manifest.infrastructure_use[]:
// each TestInfra* function exercises the exact API claimed in an
// infrastructure_use entry (PROTOCOL 8.1/8.2). The manifest test proves each
// recorded test name exists and reaches the claimed symbol; these tests prove
// the reuse behaves as the manifest states. Imports alone never count as
// adoption; every import below runs in an assertion.

// TestInfraRationalEncoding exercises physmath.NewRational: the five exact
// rational parameter artifacts are minted, are byte-identical on rebuild
// (content addressing), resolve as expression/1, and decode to values equal
// to the independent solver-local big.Rat parameters in the axion package.
func TestInfraRationalEncoding(t *testing.T) {
	store := physmath.NewMemStore()
	refs, err := representation.MintNRInputs(store)
	if err != nil {
		t.Fatalf("MintNRInputs: %v", err)
	}
	if len(refs) != 5 {
		t.Fatalf("want 5 parameter artifacts, got %d", len(refs))
	}
	// Determinism: a fresh store mints identical references.
	store2 := physmath.NewMemStore()
	refs2, err := representation.MintNRInputs(store2)
	if err != nil {
		t.Fatalf("MintNRInputs rebuild: %v", err)
	}
	for i := range refs {
		if refs[i] != refs2[i] {
			t.Fatalf("parameter ref %d not deterministic: %s vs %s", i, refs[i], refs2[i])
		}
	}
	// Every artifact resolves as expression/1 and decodes to the exact value.
	decoded := make([]*big.Rat, len(refs))
	for i, ref := range refs {
		fam, payload, err := store.Lookup(ref)
		if err != nil {
			t.Fatalf("lookup %d: %v", i, err)
		}
		if fam != physmath.FamilyExpression {
			t.Fatalf("family %d is %s/%s, want expression/1", i, fam.Namespace, fam.ArtifactType)
		}
		expr, err := physmath.DecodeExpr(artifact.NewDecoder(payload))
		if err != nil {
			t.Fatalf("decode %d: %v", i, err)
		}
		if expr.Kind != physmath.TagRational {
			t.Fatalf("kind %d is %v, want Rational", i, expr.Kind)
		}
		decoded[i] = new(big.Rat).SetFrac(expr.Numerator.ToBig(), expr.Denominator.ToBig())
	}
	// Cross-path: the phys-math-encoded values equal the local big.Rat params.
	p := axion.ExactParameters()
	want := []*big.Rat{p.GEVInv, p.OmegaP, p.N, big.NewRat(25, 1), p.VA}
	for i := range want {
		if decoded[i].Cmp(want[i]) != 0 {
			t.Fatalf("parameter %d mismatch: represented %s vs local %s", i, decoded[i], want[i])
		}
	}
}

// TestInfraArtifactIdentity exercises artifact.ComputeRef: the derivation/1
// artifact reference is re-derived from its canonical payload and equals the
// store-minted reference; a single-byte payload corruption changes the
// computed reference and is rejected. This ties the package to the single
// artifact digest formula (phys-artifact).
func TestInfraArtifactIdentity(t *testing.T) {
	store := physmath.NewMemStore()
	bodyRef, fnRef, err := representation.MintNRFormula(store)
	if err != nil {
		t.Fatalf("MintNRFormula: %v", err)
	}
	repRefs, err := representation.MintNRInputs(store)
	if err != nil {
		t.Fatalf("MintNRInputs: %v", err)
	}
	d, _, err := representation.BuildNRDerivation(store, bodyRef, fnRef, repRefs)
	if err != nil {
		t.Fatalf("BuildNRDerivation: %v", err)
	}
	payload, err := d.Encode()
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	ref, err := store.Put(physmath.FamilyDerivation, payload)
	if err != nil {
		t.Fatalf("Put derivation: %v", err)
	}
	if err := representation.VerifyContentAddress(store, ref, physmath.FamilyDerivation, payload); err != nil {
		t.Fatalf("VerifyContentAddress: %v", err)
	}
	// Corruption: one flipped byte must re-derive a different reference.
	corrupt := make([]byte, len(payload))
	copy(corrupt, payload)
	corrupt[len(corrupt)-1] ^= 0x01
	badRef, err := artifact.ComputeRef(physmath.FamilyDerivation.Namespace,
		physmath.FamilyDerivation.SchemaVersion, physmath.FamilyDerivation.ArtifactType, corrupt)
	if err != nil {
		t.Fatalf("ComputeRef corrupt: %v", err)
	}
	if badRef == ref {
		t.Fatal("corrupted payload produced the same reference")
	}
}

// TestInfraDerivationReplay exercises physmath.Apply (EXPR-SUBSTITUTE-001,
// EXPR-SIMPLIFY-001) and physmath.Replay: the recorded two-step derivation
// replays byte-exactly from its premises, its result equals the independent
// solver-local exact computation, and a deliberately corrupted declared
// output is rejected. This test establishes arithmetic reproducibility only,
// never the physics of the premise (PROTOCOL 10/12.8).
func TestInfraDerivationReplay(t *testing.T) {
	store := physmath.NewMemStore()
	bodyRef, fnRef, err := representation.MintNRFormula(store)
	if err != nil {
		t.Fatalf("MintNRFormula: %v", err)
	}
	repRefs, err := representation.MintNRInputs(store)
	if err != nil {
		t.Fatalf("MintNRInputs: %v", err)
	}
	d, resultRef, err := representation.BuildNRDerivation(store, bodyRef, fnRef, repRefs)
	if err != nil {
		t.Fatalf("BuildNRDerivation: %v", err)
	}
	if err := d.ValidateStructure(); err != nil {
		t.Fatalf("ValidateStructure: %v", err)
	}
	if err := physmath.Replay(store, d, nil); err != nil {
		t.Fatalf("Replay: %v", err)
	}
	if len(d.Steps) != 2 ||
		d.Steps[0].RuleID != "EXPR-SUBSTITUTE-001" || d.Steps[1].RuleID != "EXPR-SIMPLIFY-001" {
		t.Fatalf("unexpected step chain: %+v", d.Steps)
	}
	// The replayed rational equals the local exact computation for m = 25 eV.
	outcome, err := representation.ReplayCheck()
	if err != nil {
		t.Fatalf("ReplayCheck: %v", err)
	}
	if outcome.ResultRef != resultRef {
		t.Fatalf("replay result ref %s != built result ref %s", outcome.ResultRef, resultRef)
	}
	local := axion.GammaNRExact(big.NewRat(25, 1))
	if outcome.Value.Cmp(local) != 0 {
		t.Fatalf("replayed value %s != local exact %s", outcome.Value, local)
	}
	// Corruption: a declared output that the rules do not reproduce must fail.
	// Claim the substituted tree (step 0 output) as the final folded result.
	bad := d
	bad.Steps = append([]physmath.DerivationStep{}, d.Steps...)
	bad.Steps[1].OutputRef = d.Steps[0].OutputRef
	if bad.Steps[1].OutputRef == d.Steps[1].OutputRef {
		t.Fatal("corruption is a no-op; test cannot establish rejection")
	}
	if err := physmath.Replay(store, bad, nil); err == nil {
		t.Fatal("corrupted derivation replayed clean")
	}
	// Invented rule IDs are impossible in the closed registry (PROTOCOL 12.4).
	invented := d
	invented.Steps = append([]physmath.DerivationStep{}, d.Steps...)
	invented.Steps[1].RuleID = "DELTA-SUBSTITUTION-ANALYTIC"
	if err := physmath.Replay(store, invented, nil); err == nil {
		t.Fatal("derivation with invented rule ID replayed clean")
	}
}

// TestInfraRelationMinting exercises physmath.EncodeRelationBody: the general
// relativistic rate is minted as a durable relation/1 artifact, re-mints to
// the identical reference, and decodes back to a relation node. No rule
// evaluates this relation (REPRESENTATIONAL ONLY); this test establishes
// representation and identity only.
func TestInfraRelationMinting(t *testing.T) {
	store := physmath.NewMemStore()
	ref, err := representation.MintGeneralRateRelation(store)
	if err != nil {
		t.Fatalf("MintGeneralRateRelation: %v", err)
	}
	fam, payload, err := store.Lookup(ref)
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}
	if fam != physmath.FamilyRelation {
		t.Fatalf("family is %s/%s, want relation/1", fam.Namespace, fam.ArtifactType)
	}
	// relation/1 payload = Enum(comparison op) || left expr || right expr.
	dec := artifact.NewDecoder(payload)
	op, err := dec.String()
	if err != nil {
		t.Fatalf("decode relation op: %v", err)
	}
	if op != string(physmath.OpEQ) {
		t.Fatalf("relation op %q, want EQ", op)
	}
	left, err := physmath.DecodeExpr(dec)
	if err != nil {
		t.Fatalf("decode relation left: %v", err)
	}
	right, err := physmath.DecodeExpr(dec)
	if err != nil {
		t.Fatalf("decode relation right: %v", err)
	}
	if !dec.Exhausted() {
		t.Fatal("relation payload has trailing bytes")
	}
	if left.Kind != physmath.TagSymbol || left.DisplayName != "Gamma" {
		t.Fatalf("relation left is %v %q, want symbol Gamma", left.Kind, left.DisplayName)
	}
	if right.Kind != physmath.TagMul {
		t.Fatalf("relation right kind %v, want Mul", right.Kind)
	}
	// Determinism in a fresh store.
	ref2, err := representation.MintGeneralRateRelation(physmath.NewMemStore())
	if err != nil {
		t.Fatalf("MintGeneralRateRelation rebuild: %v", err)
	}
	if ref2 != ref {
		t.Fatalf("relation ref not deterministic: %s vs %s", ref2, ref)
	}
}
