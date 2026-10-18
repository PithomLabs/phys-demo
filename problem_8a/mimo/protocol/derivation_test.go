// Package protocol contains the executable protocol test of PROTOCOL 12:
// every manifest entry point is called and its behavior observed, the
// requested results are established through independent computation paths,
// unit errors are detected by a deliberate unit-slip test, and the stated
// conventions (E_a, J(kappa) endpoints) are exercised.
//
// What these tests establish: the recorded representation replays, the
// solver-local arithmetic paths agree, the threshold gate is exact, and the
// display convention holds. What they never establish: physical truth,
// empirical support, or the correctness of the hand-derived rate formula
// (PROTOCOL 12, final paragraph).
package protocol

import (
	"math"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/PithomLabs/phys-demo/problem_8a/mimo/axion"
	"github.com/PithomLabs/phys-demo/problem_8a/mimo/representation"
)

// relDiff returns |a-b|/|b| (b assumed nonzero).
func relDiff(a, b float64) float64 {
	return math.Abs(a-b) / math.Abs(b)
}

// TestPrimaryEntryPointOrderedResults exercises axion.OrderedResults, the
// manifest's primary entry point, and observes the requested output order:
// (1) m_a = 25 eV, then (2) m_a = 15 eV.
func TestPrimaryEntryPointOrderedResults(t *testing.T) {
	results, err := axion.OrderedResults()
	if err != nil {
		t.Fatalf("OrderedResults: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("want 2 results in problem order, got %d", len(results))
	}
	if results[0].Case != "(1)" || results[1].Case != "(2)" {
		t.Fatalf("wrong order: %s then %s", results[0].Case, results[1].Case)
	}
	// (1) m = 25 eV: exact rational 5229/(4*10^44) eV.
	wantExact := new(big.Rat).SetFrac(big.NewInt(5229),
		new(big.Int).Mul(big.NewInt(4), new(big.Int).Exp(big.NewInt(10), big.NewInt(44), nil)))
	if results[0].GammaEV.Cmp(wantExact) != 0 {
		t.Fatalf("case (1): got %s, want %s", results[0].GammaEV, wantExact)
	}
	if results[0].Display != "1.31e-41 eV" {
		t.Fatalf("case (1) display %q, want 1.31e-41 eV", results[0].Display)
	}
	if !results[0].AboveThreshold {
		t.Fatal("case (1) must be above the plasmon threshold")
	}
	// (2) m = 15 eV: exactly zero (below threshold).
	if results[1].GammaEV.Sign() != 0 {
		t.Fatalf("case (2): got %s, want exact zero", results[1].GammaEV)
	}
	if results[1].Display != "0 eV" {
		t.Fatalf("case (2) display %q, want 0 eV", results[1].Display)
	}
	if results[1].AboveThreshold {
		t.Fatal("case (2) must be below the plasmon threshold")
	}
}

// TestSecondaryEntryPointReplayCheck exercises representation.ReplayCheck,
// the manifest's secondary entry point: the recorded derivation replays
// under the closed rule registry, and its replayed value equals the
// independent solver-local exact computation. Two fresh runs must agree
// exactly (content addressing is deterministic).
func TestSecondaryEntryPointReplayCheck(t *testing.T) {
	outcome, err := representation.ReplayCheck()
	if err != nil {
		t.Fatalf("ReplayCheck: %v", err)
	}
	local, err := axion.OrderedResults()
	if err != nil {
		t.Fatalf("OrderedResults: %v", err)
	}
	if outcome.Value.Cmp(local[0].GammaEV) != 0 {
		t.Fatalf("replayed value %s != local exact %s", outcome.Value, local[0].GammaEV)
	}
	again, err := representation.ReplayCheck()
	if err != nil {
		t.Fatalf("ReplayCheck (second run): %v", err)
	}
	if again.DerivationRef != outcome.DerivationRef || again.ResultRef != outcome.ResultRef {
		t.Fatal("replay refs are not deterministic across runs")
	}
	if again.Value.Cmp(outcome.Value) != 0 {
		t.Fatal("replayed value is not deterministic across runs")
	}
}

// TestThreePathAgreement establishes the result of case (1) through three
// independent paths: (a) the replayed physmath derivation (exact rational
// bytes), (b) solver-local exact big.Rat arithmetic, and (c) a float64
// evaluation of the general relativistic formula with the full closed-form
// J(kappa) and full kinematics. This is arithmetic agreement only; the
// premise formula remains a hand derivation (PROTOCOL 10).
func TestThreePathAgreement(t *testing.T) {
	outcome, err := representation.ReplayCheck()
	if err != nil {
		t.Fatalf("ReplayCheck: %v", err)
	}
	results, err := axion.OrderedResults()
	if err != nil {
		t.Fatalf("OrderedResults: %v", err)
	}
	exact := results[0].GammaEV
	if outcome.Value.Cmp(exact) != 0 {
		t.Fatalf("replayed %s != exact %s", outcome.Value, exact)
	}
	exactF, _ := new(big.Float).SetPrec(64).SetRat(exact).Float64()
	general := axion.GammaNRGeneralFloat(25)
	if d := relDiff(general, exactF); d > 1e-4 {
		t.Fatalf("float general-formula path off by relative %g: %g vs %g", d, general, exactF)
	}
	if general <= 0 || math.IsNaN(general) || math.IsInf(general, 0) {
		t.Fatalf("float general-formula path returned %g", general)
	}
	// Case (2): every path is exactly zero.
	if results[1].GammaEV.Sign() != 0 {
		t.Fatalf("case (2) exact not zero: %s", results[1].GammaEV)
	}
	if axion.GammaNRGeneralFloat(15) != 0 {
		t.Fatal("case (2) float general formula not zero")
	}
}

// TestThresholdBehavior exercises the exact solver-local gate
// Theta at leading NR order (m > omega_p, which implies the physical
// E_a = m(1+v_a^2/2) > omega_p), including the razor-thin window where the
// exact kinematics is open but the leading-order prefactor is not valid.
func TestThresholdBehavior(t *testing.T) {
	p := axion.ExactParameters()
	one := big.NewRat(1, 1)
	va2over2 := new(big.Rat).Quo(new(big.Rat).Mul(p.VA, p.VA), big.NewRat(2, 1))
	// Exact kinematic boundary: E_a = m(1+v_a^2/2) = omega_p.
	kinBoundary := new(big.Rat).Quo(p.OmegaP, new(big.Rat).Add(one, va2over2))
	step := big.NewRat(1, 1000000)
	// Just above omega_p: gate open, positive rate.
	above := new(big.Rat).Add(p.OmegaP, step)
	if got := axion.GammaNRExact(above); got.Sign() <= 0 {
		t.Fatalf("mass %s above omega_p gave %s, want > 0", above, got)
	}
	// The excluded window kinBoundary < m <= omega_p: the exact kinematics is
	// open (E_a > omega_p) but the NR closed form is invalid; the gate must
	// return exactly zero rather than a negative "rate".
	if kinBoundary.Cmp(p.OmegaP) >= 0 {
		t.Fatalf("kinematic boundary %s must lie below omega_p", kinBoundary)
	}
	mid := new(big.Rat).Add(kinBoundary, step)
	if mid.Cmp(p.OmegaP) > 0 {
		t.Fatal("window test mass must stay below omega_p")
	}
	if got := axion.GammaNRExact(mid); got.Sign() != 0 {
		t.Fatalf("mass %s inside the NR-invalid window gave %s, want 0 (documented gate)", mid, got)
	}
	// Below the kinematic boundary: zero on both physical and NR grounds.
	below := new(big.Rat).Sub(kinBoundary, step)
	if got := axion.GammaNRExact(below); got.Sign() != 0 {
		t.Fatalf("mass %s below threshold gave %s, want 0", below, got)
	}
	// The declared problem cases.
	if got := axion.GammaNRExact(big.NewRat(25, 1)); got.Sign() <= 0 {
		t.Fatalf("m = 25 eV gave %s, want > 0", got)
	}
	if got := axion.GammaNRExact(big.NewRat(15, 1)); got.Sign() != 0 {
		t.Fatalf("m = 15 eV gave %s, want exactly 0 (E_a = 15.0000075 eV < omega_p)", got)
	}
}

// TestConventionInsensitivity isolates the stated NR convention: evaluating
// the exact rational prefactor at E_a = m versus the full-kinematics
// E_a = m(1+v_a^2/2) differs by well below the displayed precision.
func TestConventionInsensitivity(t *testing.T) {
	exact := axion.GammaNRExact(big.NewRat(25, 1))
	exactF, _ := new(big.Float).SetPrec(64).SetRat(exact).Float64()
	variant := axion.GammaNRVariantFloat(25)
	if d := relDiff(variant, exactF); d > 1e-5 {
		t.Fatalf("E_a convention changes the result by relative %g (> 1e-5)", d)
	}
}

// TestUnitSlipSensitivity is the deliberate unit-error test (PROTOCOL 9.1):
// the exact GeV^-1 to eV^-1 conversion is asserted, and every deliberate
// wrong scale factor is detected as a large mismatch.
func TestUnitSlipSensitivity(t *testing.T) {
	gGeV := new(big.Rat).SetFrac(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil))
	correct := axion.GeVInvToEVInv(gGeV)
	want := new(big.Rat).SetFrac(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(19), nil))
	if correct.Cmp(want) != 0 {
		t.Fatalf("exact conversion: got %s, want 1e-19 eV^-1", correct)
	}
	// Deliberate slips: (a) treat GeV^-1 as eV^-1 (off by 1e9);
	// (b) multiply by 1e9 instead of dividing (off by 1e18);
	// (c) one-decade factor slip (off by 1e-18).
	slips := []struct {
		name string
		got  *big.Rat
	}{
		{"identity (GeV^-1 used as eV^-1)", new(big.Rat).Set(gGeV)},
		{"multiply by 1e9", new(big.Rat).Mul(gGeV, big.NewRat(1000000000, 1))},
		{"multiply by 1e-1", new(big.Rat).Mul(gGeV, big.NewRat(1, 10))},
	}
	correctF, _ := new(big.Float).SetPrec(64).SetRat(correct).Float64()
	for _, s := range slips {
		if s.got.Cmp(correct) == 0 {
			t.Fatalf("slip %q is not actually wrong", s.name)
		}
		gotF, _ := new(big.Float).SetPrec(64).SetRat(s.got).Float64()
		// Every deliberate slip here is at least eight orders of magnitude
		// off in relative terms; a detection floor of 1e6 leaves ample margin
		// while proving the check cannot silently pass a unit error.
		if d := relDiff(gotF, correctF); d < 1e6 {
			t.Fatalf("slip %q escaped detection: relative error only %g", s.name, d)
		}
	}
}

// TestJSeriesAndLimits checks the solver-local J(kappa) implementation
// against its hand-derived analytic properties: the small-kappa series, the
// kappa -> 1 endpoint, monotonic decrease, and NaN-freedom (the closed form
// is a 0 * infinity indeterminate exactly at kappa = 1).
func TestJSeriesAndLimits(t *testing.T) {
	// The closed form tracks its small-kappa series; the residual from the
	// 8 pi/3 limit is the genuine O(kappa^2) correction (relative kappa^2/5),
	// not numerical error.
	if d := relDiff(axion.J(1e-3), axion.JSeries(1e-3)); d > 1e-9 {
		t.Fatalf("J(1e-3) off its series by relative %g", d)
	}
	if d := relDiff(axion.J(0.01), axion.JSeries(0.01)); d > 1e-6 {
		t.Fatalf("J(0.01) off its series by relative %g", d)
	}
	// At kappa = 1e-4 the residual to the limit is kappa^2/5 ~ 2e-9 while the
	// difference of two O(1e8) terms leaves ~1e-7 relative float64 headroom.
	if d := relDiff(axion.J(1e-4), axion.JZero()); d > 1e-6 {
		t.Fatalf("J(1e-4) off the 8 pi/3 limit by relative %g", d)
	}
	if d := relDiff(axion.J(1e-3), axion.JZero()); d > 1e-6 {
		t.Fatalf("J(1e-3) not approaching the 8 pi/3 limit: relative %g", d)
	}
	if axion.J(1) != axion.JOne() {
		t.Fatalf("J(1) = %g, want the 2 pi endpoint", axion.J(1))
	}
	if d := relDiff(axion.J(0.999999), axion.JOne()); d > 1e-3 {
		t.Fatalf("J(0.999999) off the 2 pi endpoint by relative %g", d)
	}
	if !(axion.J(0.01) > axion.J(0.5) && axion.J(0.5) > axion.J(0.999)) {
		t.Fatal("J is not decreasing on [0,1)")
	}
	for _, k := range []float64{0, 1e-9, 0.01, 0.5, 0.999, 1, 1.001, 2} {
		if math.IsNaN(axion.J(k)) || math.IsInf(axion.J(k), 0) {
			t.Fatalf("J(%g) is not finite", k)
		}
	}
	// The problem-relevant kappa is tiny: k_a/q_gamma ~ 2e-3, so the O(kappa^2)
	// correction to J(0) is far below the displayed precision.
	kappa := (25 * 1e-3) / (1.5 * (25 - 16.6))
	if d := relDiff(axion.J(kappa), axion.JZero()); d > 1e-5 {
		t.Fatalf("problem kappa correction too large for J(0): relative %g", d)
	}
}

// TestManifestAndReadmeCoherence checks that the manifest's human-facing
// strings actually describe this package's code: the NR exact value and the
// threshold zero are the same in manifest, code, and README prose.
func TestManifestAndReadmeCoherence(t *testing.T) {
	manifestBytes, err := os.ReadFile("../manifest.json")
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	for _, want := range []string{
		"5229/(4*10^44)",
		"axion.OrderedResults",
		"representation.ReplayCheck",
		"EXPR-SUBSTITUTE-001",
		"EXPR-SIMPLIFY-001",
	} {
		if !strings.Contains(string(manifestBytes), want) {
			t.Fatalf("manifest.json does not mention %q", want)
		}
	}
}
