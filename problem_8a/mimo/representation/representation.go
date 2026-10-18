// Package representation mints the Phys representation of Problem 8(a) and
// the replayable substitution derivation. This package owns every Phys
// library import in the mimo module; the sibling axion package is
// deliberately Phys-free (PROTOCOL §9.2 solver-local boundary).
//
// What is represented and at what support state (PROTOCOL §5/§6 stage 6):
//
//   - The parameterized non-relativistic rate
//     Gamma = g^2 omega_p n^2 (m - omega_p) m v_a^2 / 6 is minted as a
//     physmath expression/1 (bound symbols) plus its declaring function/1.
//     Its substitution of the exact rational parameter values and exact
//     rational folding are EXECUTED by the closed rule registry
//     (EXPR-SUBSTITUTE-001 then EXPR-SIMPLIFY-001) and replayed byte-exactly
//     by physmath.Replay.
//
//   - The general relativistic rate with the J(kappa) angular factor and the
//     Heaviside threshold is minted as a schematic relation/1 (ephemeral
//     symbols plus the registered pi constant). No rule in the closed
//     registry evaluates it; it is REPRESENTATIONAL ONLY, never claimed as
//     executable.
//
// The analytic physics steps (delta-collapse of the plasmon-pole response,
// tensor contraction, phase-space integration, the kappa -> 0 limit) have no
// executable rule and are recorded honestly in the manifest, never simulated.
package representation

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/PithomLabs/phys-artifact"
	physmath "github.com/PithomLabs/phys-math"
)

// nrParamNames are the display names of the five ordinal-derived parameter
// slots of the NR rate function, in ParameterID ordinal order.
var nrParamNames = []string{"g", "omega_p", "n", "m", "v_a"}

// nrFormulaBody builds the bound-symbol expression
//
//	g^2 * omega_p * n^2 * (m - omega_p) * m * v_a^2 * (1/6).
//
// It uses v_a * v_a (not Pow) so that EXPR-SIMPLIFY-001's all-rational fold
// applies to the whole tree after substitution.
func nrFormulaBody() physmath.Expr {
	pids := make([]physmath.ParameterID, len(nrParamNames))
	for i := range pids {
		pids[i] = physmath.DeriveParameterID(uint32(i))
	}
	g := physmath.NewBoundSymbol("g", pids[0])
	wp := physmath.NewBoundSymbol("omega_p", pids[1])
	n := physmath.NewBoundSymbol("n", pids[2])
	m := physmath.NewBoundSymbol("m", pids[3])
	v := physmath.NewBoundSymbol("v_a", pids[4])
	return physmath.Mul(
		g, g, wp, n, n,
		physmath.Add(m, physmath.Neg(wp)),
		m, v, v,
		physmath.MustRational(1, 6),
	)
}

// nrParameterIDs returns the ordinal-derived parameter slots.
func nrParameterIDs() []physmath.ParameterID {
	ids := make([]physmath.ParameterID, len(nrParamNames))
	for i := range ids {
		ids[i] = physmath.DeriveParameterID(uint32(i))
	}
	return ids
}

// MintNRFormula mints the parameterized NR rate expression (expression/1)
// and its declaring function/1 (whose ParameterIDs are ordinal-derived and
// whose ParamNames record the display authority for each slot).
// It returns (bodyRef, functionRef).
func MintNRFormula(store physmath.Store) (artifact.ArtifactRef, artifact.ArtifactRef, error) {
	body := nrFormulaBody()
	bodyRef, err := physmath.PutExpression(store, body)
	if err != nil {
		return artifact.ArtifactRef{}, artifact.ArtifactRef{},
			fmt.Errorf("representation: mint NR body: %w", err)
	}
	fn := physmath.Function{
		ParameterIDs: nrParameterIDs(),
		BodyRef:      bodyRef,
		ParamNames:   append([]string{}, nrParamNames...),
	}
	fnPayload, err := fn.Encode()
	if err != nil {
		return artifact.ArtifactRef{}, artifact.ArtifactRef{},
			fmt.Errorf("representation: encode NR function: %w", err)
	}
	fnRef, err := store.Put(physmath.FamilyFunction, fnPayload)
	if err != nil {
		return artifact.ArtifactRef{}, artifact.ArtifactRef{},
			fmt.Errorf("representation: mint NR function: %w", err)
	}
	return bodyRef, fnRef, nil
}

// MintNRInputs mints the five exact-rational parameter replacement
// artifacts (expression/1) through physmath.NewRational: g = 10^-19 eV^-1,
// omega_p = 83/5 eV, n = 3/2, m = 25 eV, v_a = 1/1000, in slot order.
func MintNRInputs(store physmath.Store) ([]artifact.ArtifactRef, error) {
	ten19 := new(big.Int).Exp(big.NewInt(10), big.NewInt(19), nil)
	vals := [][2]*big.Int{
		{big.NewInt(1), ten19},            // g = 1e-19 eV^-1
		{big.NewInt(83), big.NewInt(5)},   // omega_p = 16.6 eV
		{big.NewInt(3), big.NewInt(2)},    // n = 1.5
		{big.NewInt(25), big.NewInt(1)},   // m = 25 eV (case 1)
		{big.NewInt(1), big.NewInt(1000)}, // v_a = 1e-3
	}
	refs := make([]artifact.ArtifactRef, 0, len(vals))
	for i, v := range vals {
		e, err := physmath.NewRational(v[0], v[1])
		if err != nil {
			return nil, fmt.Errorf("representation: rational %d: %w", i, err)
		}
		ref, err := physmath.PutExpression(store, e)
		if err != nil {
			return nil, fmt.Errorf("representation: mint rational %d: %w", i, err)
		}
		refs = append(refs, ref)
	}
	return refs, nil
}

// substitutionParameters encodes the EXPR-SUBSTITUTE-001 parameter payload:
// UInt32BE(N) || N entries (BoundParameterRef || replacement expression ref),
// strictly sorted by exact canonical entry bytes (rule contract §17.8).
func substitutionParameters(fnRef artifact.ArtifactRef, repRefs []artifact.ArtifactRef) ([]byte, error) {
	ids := nrParameterIDs()
	if len(ids) != len(repRefs) {
		return nil, fmt.Errorf("representation: %d slots vs %d replacements", len(ids), len(repRefs))
	}
	raws := make([][]byte, len(ids))
	for i := range ids {
		bpr := physmath.BoundParameterRef{
			DeclaringArtifactRef: fnRef,
			ParameterID:          ids[i],
		}
		raws[i] = append(bpr.Encode(), repRefs[i].Encode()...)
	}
	for i := 1; i < len(raws); i++ {
		for j := i; j > 0 && bytes.Compare(raws[j-1], raws[j]) > 0; j-- {
			raws[j-1], raws[j] = raws[j], raws[j-1]
		}
	}
	return artifact.EncodeSequence(raws), nil
}

// BuildNRDerivation executes the two closed-registry steps (construction
// phase) and assembles the replayable physmath/derivation/1 record:
//
//	0: EXPR-SUBSTITUTE-001 (formula expression, declaring function) -> substituted tree
//	1: EXPR-SIMPLIFY-001  (exact all-rational fold)                   -> single Rational
//
// It returns the derivation and its result reference.
func BuildNRDerivation(store physmath.Store, bodyRef, fnRef artifact.ArtifactRef, repRefs []artifact.ArtifactRef) (physmath.Derivation, artifact.ArtifactRef, error) {
	params, err := substitutionParameters(fnRef, repRefs)
	if err != nil {
		return physmath.Derivation{}, artifact.ArtifactRef{}, err
	}
	sub, err := physmath.Apply(store, physmath.Invocation{
		ProducerNamespace: "physmath",
		RuleID:            "EXPR-SUBSTITUTE-001",
		InputRefs:         []artifact.ArtifactRef{bodyRef, fnRef},
		Parameters:        params,
	})
	if err != nil {
		return physmath.Derivation{}, artifact.ArtifactRef{},
			fmt.Errorf("representation: EXPR-SUBSTITUTE-001: %w", err)
	}
	simp, err := physmath.Apply(store, physmath.Invocation{
		ProducerNamespace: "physmath",
		RuleID:            "EXPR-SIMPLIFY-001",
		InputRefs:         []artifact.ArtifactRef{sub.OutputRef},
	})
	if err != nil {
		return physmath.Derivation{}, artifact.ArtifactRef{},
			fmt.Errorf("representation: EXPR-SIMPLIFY-001: %w", err)
	}
	d := physmath.Derivation{
		PremiseRefs: []artifact.ArtifactRef{bodyRef, fnRef},
		Steps: []physmath.DerivationStep{
			{
				Ordinal:           0,
				ProducerNamespace: "physmath",
				RuleID:            "EXPR-SUBSTITUTE-001",
				InputRefs:         []artifact.ArtifactRef{bodyRef, fnRef},
				OutputRef:         sub.OutputRef,
				Parameters:        params,
			},
			{
				Ordinal:           1,
				ProducerNamespace: "physmath",
				RuleID:            "EXPR-SIMPLIFY-001",
				InputRefs:         []artifact.ArtifactRef{sub.OutputRef},
				OutputRef:         simp.OutputRef,
			},
		},
		ResultRefs: []artifact.ArtifactRef{simp.OutputRef},
	}
	return d, simp.OutputRef, nil
}

// MintGeneralRateRelation mints the schematic general relativistic rate as a
// relation/1 artifact:
//
//	Gamma = g^2 omega_p q_gamma^2 k_a^2 J_kappa Theta (16 pi E_a omega_gamma)^-1.
//
// Symbols are ephemeral (schematic expression, Spec §8.8) except pi, which
// uses the registered constant physmath.ConstPi. No closed-registry rule
// evaluates this relation: its support state is REPRESENTATIONAL ONLY
// (the J(kappa) closed form and the Heaviside gate are recorded in the
// manifest conventions and executed solver-locally where needed).
func MintGeneralRateRelation(store physmath.Store) (artifact.ArtifactRef, error) {
	sym := physmath.NewEphemeralSymbol
	inv := func(e physmath.Expr) physmath.Expr { return physmath.Pow(e, physmath.MustRational(-1, 1)) }
	right := physmath.Mul(
		sym("g"), sym("g"),
		sym("omega_p"),
		physmath.Pow(sym("q_gamma"), physmath.MustRational(2, 1)),
		physmath.Pow(sym("k_a"), physmath.MustRational(2, 1)),
		sym("J_kappa"),
		sym("Theta_step"),
		inv(sym("E_a")),
		inv(sym("omega_gamma")),
		physmath.NewConstantSymbol("pi", physmath.ConstPi),
		physmath.MustRational(1, 16),
	)
	rel, err := physmath.Relation(physmath.OpEQ, sym("Gamma"), right)
	if err != nil {
		return artifact.ArtifactRef{}, fmt.Errorf("representation: general-rate relation: %w", err)
	}
	payload, err := physmath.EncodeRelationBody(physmath.OpEQ, *rel.Left, *rel.Right)
	if err != nil {
		return artifact.ArtifactRef{}, fmt.Errorf("representation: encode general-rate relation: %w", err)
	}
	ref, err := store.Put(physmath.FamilyRelation, payload)
	if err != nil {
		return artifact.ArtifactRef{}, fmt.Errorf("representation: mint general-rate relation: %w", err)
	}
	return ref, nil
}

// VerifyContentAddress re-derives the content-addressed reference of
// family/payload through artifact.ComputeRef (the single artifact digest
// formula) and requires it to equal ref. It is the package's direct
// phys-artifact identity check.
func VerifyContentAddress(store physmath.Store, ref artifact.ArtifactRef, fam physmath.FamilyID, payload []byte) error {
	want, err := artifact.ComputeRef(fam.Namespace, fam.SchemaVersion, fam.ArtifactType, payload)
	if err != nil {
		return fmt.Errorf("representation: ComputeRef: %w", err)
	}
	if want != ref {
		return fmt.Errorf("representation: content-address mismatch: computed %s, store says %s", want, ref)
	}
	if _, _, err := store.Lookup(ref); err != nil {
		return fmt.Errorf("representation: minted ref does not resolve: %w", err)
	}
	return nil
}

// ReplayOutcome is what ReplayCheck establishes: the durable derivation
// artifact, its replay-verified result reference, and the exact rational
// value decoded from the replayed result bytes.
type ReplayOutcome struct {
	DerivationRef artifact.ArtifactRef
	ResultRef     artifact.ArtifactRef
	Value         *big.Rat
}

// ReplayCheck runs the full package pipeline in a fresh store and replays
// the recorded derivation under the closed rule registry (physmath.Replay,
// Spec §20.4): the recomputed output of every step must equal the declared
// output byte-exactly. It returns the decoded exact rational result.
//
// What this establishes (PROTOCOL §10/§12): reproducibility of the recorded
// substitution/fold steps from the bound premises under the operation
// semantics. It does NOT establish the physics of the premises — the rate
// formula is a hand derivation recorded at evidence level Represented.
func ReplayCheck() (ReplayOutcome, error) {
	store := physmath.NewMemStore()
	bodyRef, fnRef, err := MintNRFormula(store)
	if err != nil {
		return ReplayOutcome{}, err
	}
	repRefs, err := MintNRInputs(store)
	if err != nil {
		return ReplayOutcome{}, err
	}
	d, resultRef, err := BuildNRDerivation(store, bodyRef, fnRef, repRefs)
	if err != nil {
		return ReplayOutcome{}, err
	}
	if err := d.ValidateStructure(); err != nil {
		return ReplayOutcome{}, fmt.Errorf("representation: derivation structure: %w", err)
	}
	// Mint the durable derivation/1 artifact and verify its identity.
	payload, err := d.Encode()
	if err != nil {
		return ReplayOutcome{}, fmt.Errorf("representation: encode derivation: %w", err)
	}
	dRef, err := store.Put(physmath.FamilyDerivation, payload)
	if err != nil {
		return ReplayOutcome{}, fmt.Errorf("representation: mint derivation: %w", err)
	}
	if err := VerifyContentAddress(store, dRef, physmath.FamilyDerivation, payload); err != nil {
		return ReplayOutcome{}, err
	}
	// Replay re-executes every step from the premises and compares outputs.
	if err := physmath.Replay(store, d, nil); err != nil {
		return ReplayOutcome{}, fmt.Errorf("representation: replay failed: %w", err)
	}
	// Decode the exact rational value from the replayed result bytes.
	fam, resultPayload, err := store.Lookup(resultRef)
	if err != nil {
		return ReplayOutcome{}, fmt.Errorf("representation: lookup result: %w", err)
	}
	if fam != physmath.FamilyExpression {
		return ReplayOutcome{}, fmt.Errorf("representation: result family is %s/%s, want expression/1",
			fam.Namespace, fam.ArtifactType)
	}
	expr, err := physmath.DecodeExpr(artifact.NewDecoder(resultPayload))
	if err != nil {
		return ReplayOutcome{}, fmt.Errorf("representation: decode result: %w", err)
	}
	if expr.Kind != physmath.TagRational {
		return ReplayOutcome{}, fmt.Errorf("representation: result root kind is %v, want Rational", expr.Kind)
	}
	value := new(big.Rat).SetFrac(expr.Numerator.ToBig(), expr.Denominator.ToBig())
	return ReplayOutcome{DerivationRef: dRef, ResultRef: resultRef, Value: value}, nil
}
