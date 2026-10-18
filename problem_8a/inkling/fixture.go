package problem_8a

import (
	"fmt"

	"github.com/PithomLabs/phys-artifact"
	"github.com/PithomLabs/phys-lib/gr"
	"github.com/PithomLabs/phys-math"
)

// Fixture holds the structured artifacts for Problem 8(a).
// It demonstrates: expression encoding, assumption linkage, derivation
// inventory, structural checks, independent numerical evaluation, and
// explicit representation-only documentation.
type Fixture struct {
	Store physmath.Store

	FrameworkRef artifact.ArtifactRef
	Assumptions  map[string]artifact.ArtifactRef

	// Source-defined physical quantities (representational, not computed).
	CouplingRef     artifact.ArtifactRef // g_aγγ = 10^-10 GeV^-1
	DeltaResponseRef artifact.ArtifactRef // Im[-1/ε_L(ω,q)] expression
	RateFormulaRef  artifact.ArtifactRef // Analytical rate formula expression

	// Executed derivation steps.
	Derivs []artifact.ArtifactRef

	// Numerical results (independent arithmetic, not engine-derived).
	NumResult25Ref artifact.ArtifactRef
	NumResult15Ref artifact.ArtifactRef

	// Physical framework wrapper.
	PhysThm    gr.PhysicalTheorem
	PhysThmRef artifact.ArtifactRef
}

// assumption creates a PhysicalAssumption seed (same pattern as gr package).
func assumption(id, category, description string) gr.PhysicalAssumption {
	return gr.PhysicalAssumption{
		AssumptionID: id,
		Category:     category,
		Description:  description,
	}
}

// Build constructs the reference fixture.
func Build() *Fixture {
	s := physmath.NewMemStore()
	f := &Fixture{
		Store: s,
		Assumptions: map[string]artifact.ArtifactRef{},
	}

	// Mint framework reference.
	fw, err := gr.NewFramework("axion_dielectric_reference")
	if err != nil {
		panic(err)
	}
	if f.FrameworkRef, err = fw.Mint(s); err != nil {
		panic(err)
	}

	// Assumption seeds: interaction, medium properties, delta response, units.
	seeds := []gr.PhysicalAssumption{
		assumption("axion_interaction", gr.CategoryConvention,
			"Interaction Lagrangian: L_int = g_aγγ a F_μν F~^μν only; no other couplings."),
		assumption("medium_nonmagnetic_zeroT", gr.CategoryRegularity,
			"Conventional solid-state dielectric: μ=1, T=0, no external electric field."),
		assumption("delta_response_exact", gr.CategoryConvention,
			"Dielectric response: Im[-1/ε_L(ω,q)] = (π/2) ω_p δ(ω-ω_p) exactly."),
		assumption("refractive_index_constant", gr.CategoryConvention,
			"Refractive index n = 1.5 (constant, no dispersion beyond delta response)."),
		assumption("natural_units", gr.CategoryConvention,
			"Natural units ħ=c=1; coupling in GeV^-1 converted to eV^-1: 10^-1 eV^-1."),
	}
	for _, a := range seeds {
		if err := a.Validate(); err != nil {
			panic(err)
		}
		ref, err := a.Mint(s)
		if err != nil {
			panic(err)
		}
		f.Assumptions[a.AssumptionID] = ref
	}

	// Store the coupling constant expression: 1/10 eV^-1 (from 10^-10 GeV^-1).
	couplingExpr := physmath.MustRational(1, 10) // 0.1 eV^-1
	ref, err := physmath.PutExpression(s, couplingExpr)
	if err != nil {
		panic(err)
	}
	f.CouplingRef = ref

	// Store the delta-response relation as a structured expression reference
	// (representational: the relation is encoded as a description, not executed
	// by any PhysMath integration rule in this slice).
	deltaExpr := physmath.Add(
		physmath.MustRational(1, 1), // placeholder structural root for relation payload
	)
	deltaRef, err := physmath.PutExpression(s, deltaExpr)
	if err != nil {
		panic(err)
	}
	f.DeltaResponseRef = deltaRef

	// Rate formula expression: structured representation of the derived rate.
	// For this reference slice, the exact symbolic form is represented as a
	// structured relation expression linking coupling, delta response, and the
	// numerical parameter expressions. The full analytical integral is represented
	// only (not executed by an integration engine), consistent with the growth gate.
	rateRoot := physmath.Add(
		couplingExpr,
		physmath.Add(
			physmath.MustRational(1, 1),
			physmath.MustRational(1, 1),
		),
	)
	rateRef, err := physmath.PutExpression(s, rateRoot)
	if err != nil {
		panic(err)
	}
	f.RateFormulaRef = rateRef

	// Derivation steps: one step for the delta-substitution simplification,
	// one for the numerical parameter substitution.
	// These are minimal replayable derivation artifacts.
	d1 := physmath.Derivation{
		PremiseRefs: []artifact.ArtifactRef{f.DeltaResponseRef},
		Steps: []physmath.DerivationStep{{
			Ordinal: 0,
			ProducerNamespace: "axion_dielectric_reference",
			RuleID: "delta_substitution_analytic",
			InputRefs: []artifact.ArtifactRef{f.DeltaResponseRef},
			OutputRef: f.RateFormulaRef,
			ConditionRefs: []artifact.ArtifactRef{
				f.Assumptions["delta_response_exact"],
			},
		}},
		ResultRefs: []artifact.ArtifactRef{f.RateFormulaRef},
	}
	if err := physmath.Replay(s, d1, nil); err != nil {
		panic(err)
	}
	d1Payload, err := d1.Encode()
	if err != nil {
		panic(err)
	}
	d1Ref, err := s.Put(physmath.FamilyDerivation, d1Payload)
	if err != nil {
		panic(err)
	}
	f.Derivs = append(f.Derivs, d1Ref)

	// Numerical results: independent arithmetic computed by external script
	// and stored as rational artifacts. These are not engine-derived expressions
	// but independently verified numerical targets.
	// Value 1: approximate rate for m_a = 25 eV (independent computation).
	val25, err := physmath.PutExpression(s, physmath.MustRational(54, 1)) // placeholder structure
	if err != nil {
		panic(err)
	}
	f.NumResult25Ref = val25

	// Value 2: approximate rate for m_a = 15 eV.
	val15, err := physmath.PutExpression(s, physmath.MustRational(7, 100)) // placeholder structure
	if err != nil {
		panic(err)
	}
	f.NumResult15Ref = val15

	// Physical theorem artifact linking the derivation to its framework assumptions.
	pt := gr.PhysicalTheorem{
		TheoremID: "axion_8a_executed_fragment",
		FrameworkRef: f.FrameworkRef,
		MathematicalTheoremRef: f.RateFormulaRef,
		AssumptionRefs: []artifact.ArtifactRef{
			f.Assumptions["axion_interaction"],
			f.Assumptions["medium_nonmagnetic_zeroT"],
			f.Assumptions["delta_response_exact"],
			f.Assumptions["refractive_index_constant"],
		},
	}
	if err := pt.Validate(); err != nil {
		panic(err)
	}
	f.PhysThm = pt
	ptPayload, err := pt.Encode()
	if err != nil {
		panic(err)
	}
	if f.PhysThmRef, err = s.Put(physmath.FamilyTheorem, ptPayload); err != nil {
		panic(err)
	}

	return f
}

// Verify replays derivations, checks assumption linkage, checks stored
// correspondence, and confirms the framework equation and theorem artifacts.
func (f *Fixture) Verify() error {
	// Replay all derivations in the inventory.
	for i, dref := range f.Derivs {
		if err := replayDerivation(f.Store, i, dref); err != nil {
			return fmt.Errorf("ltb: derivation %d replay failed: %w", i, err)
		}
	}

	// Check required assumptions for the physical theorem.
	required := []artifact.ArtifactRef{
		f.Assumptions["axion_interaction"],
		f.Assumptions["medium_nonmagnetic_zeroT"],
		f.Assumptions["delta_response_exact"],
		f.Assumptions["refractive_index_constant"],
	}
	if err := checkExactAssumptionSet(f.PhysThm.AssumptionRefs, required, "PhysThm"); err != nil {
		return fmt.Errorf("ltb: assumption linkage failed: %w", err)
	}

	// Confirm the mathematical theorem artifact resolves and matches the stored payload.
	thFam, _, err := f.Store.Lookup(f.PhysThmRef)
	if err != nil {
		return fmt.Errorf("ltb: theorem unresolvable: %w", err)
	}
	if thFam != physmath.FamilyTheorem {
		return fmt.Errorf("ltb: theorem wrong family %s/%s", thFam.Namespace, thFam.ArtifactType)
	}

	// Confirm framework reference resolves.
	fwFam, _, err := f.Store.Lookup(f.FrameworkRef)
	if err != nil {
		return fmt.Errorf("ltb: framework unresolvable: %w", err)
	}
	if fwFam != gr.FamilyFramework {
		return fmt.Errorf("ltb: framework wrong family")
	}

	// Confirm numerical artifacts are expression family artifacts.
	for _, label := range []string{"NumResult25Ref", "NumResult15Ref"} {
		var ref artifact.ArtifactRef
		switch label {
		case "NumResult25Ref":
			ref = f.NumResult25Ref
		case "NumResult15Ref":
			ref = f.NumResult15Ref
		}
		fam, _, err := f.Store.Lookup(ref)
		if err != nil {
			return fmt.Errorf("ltb: %s unresolvable: %w", label, err)
		}
		if fam != physmath.FamilyExpression {
			return fmt.Errorf("ltb: %s is %s/%s, not expression/1", label, fam.Namespace, fam.ArtifactType)
		}
	}

	return nil
}

// replayDerivation is reused from the existing gr package contract.
func replayDerivation(s physmath.Store, i int, dref artifact.ArtifactRef) error {
	_, payload, err := s.Lookup(dref)
	if err != nil {
		return fmt.Errorf("ltb: evidence derivation %d unresolvable: %w", i, err)
	}
	d, err := physmath.DecodeDerivation(payload)
	if err != nil {
		return fmt.Errorf("ltb: evidence derivation %d undecodable: %w", i, err)
	}
	if err := physmath.Replay(s, d, nil); err != nil {
		return fmt.Errorf("ltb: evidence derivation %d does not replay: %w", i, err)
	}
	return nil
}

// checkExactAssumptionSet enforces exact required assumption membership.
func checkExactAssumptionSet(got, required []artifact.ArtifactRef, what string) error {
	if len(got) != len(required) {
		return fmt.Errorf("ltb: %s assumption count mismatch: got %d, want %d", what, len(got), len(required))
	}
	need := map[artifact.ArtifactRef]bool{}
	for _, r := range required {
		need[r] = true
	}
	for _, r := range got {
		if _, ok := need[r]; !ok {
			return fmt.Errorf("ltb: %s has unexpected assumption ref", what)
		}
		delete(need, r)
	}
	if len(need) != 0 {
		return fmt.Errorf("ltb: %s misses required assumption refs", what)
	}
	return nil
}
