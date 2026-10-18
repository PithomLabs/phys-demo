// Package protocol contains the executable protocol test for the
// Problem 8(a) gemini benchmark solution. Every test documents what it
// establishes and what it does not (protocol §12 item 8).
package protocol

import (
	"math"
	"math/big"
	"testing"

	"github.com/PithomLabs/phys-demo/problem_8a/gemini/axion"
)

// ---------------------------------------------------------------------------
// Entry-point exercise (protocol §12 item 1)
// ---------------------------------------------------------------------------

func TestOrderedResults(t *testing.T) {
	// Exercises the primary entry point axion.OrderedResults.
	// Establishes: the function returns a 2-element array with the
	// correct mass labels, open/closed status, and rate values within
	// tolerance. Does NOT establish: physical truth of the rate formula.
	p := axion.BenchmarkParams()
	results := axion.OrderedResults(p)

	// Case 1: m_a = 25 eV (open channel)
	if !results[0].Open {
		t.Fatal("case 1 (25 eV) should be open")
	}
	if results[0].MassEV != "25" {
		t.Fatalf("case 1 mass = %q, want 25", results[0].MassEV)
	}
	// Golden value: ~1.31e-41 eV (independently derived)
	if math.Abs(results[0].RateEV/1.30725026e-41-1) > 1e-6 {
		t.Fatalf("case 1 rate = %e, want ~1.31e-41", results[0].RateEV)
	}

	// Case 2: m_a = 15 eV (closed channel: Ea < wp)
	if results[1].Open {
		t.Fatal("case 2 (15 eV) should be closed")
	}
	if results[1].MassEV != "15" {
		t.Fatalf("case 2 mass = %q, want 15", results[1].MassEV)
	}
	if results[1].RateEV != 0 {
		t.Fatalf("case 2 rate = %e, want 0", results[1].RateEV)
	}
}

func TestSkeletonNR(t *testing.T) {
	// Exercises the secondary entry point axion.SkeletonNR.
	// Establishes: exact-rational arithmetic produces a positive
	// rational for the open channel and zero for the closed channel.
	p := axion.BenchmarkParams()

	s25, ok25 := axion.SkeletonNR(p, p.Mass25)
	if !ok25 {
		t.Fatal("skeleton should be open for 25 eV")
	}
	if s25.Sign() <= 0 {
		t.Fatalf("skeleton(25) = %s, want positive", s25.RatString())
	}

	s15, ok15 := axion.SkeletonNR(p, p.Mass15)
	if ok15 {
		t.Fatal("skeleton should be closed for 15 eV")
	}
	if s15.Sign() != 0 {
		t.Fatalf("skeleton(15) = %s, want 0", s15.RatString())
	}
}

// ---------------------------------------------------------------------------
// Physics tests (solver-local, documented)
// ---------------------------------------------------------------------------

func TestJClosedSpecialValues(t *testing.T) {
	// J(0) = 8*pi/3 exactly (NR limit).
	// J(1) = 2*pi (removable singularity).
	// Establishes: closed-form evaluation reproduces known limits.
	if got := axion.JClosed(0); math.Abs(got-8*math.Pi/3) > 1e-14 {
		t.Fatalf("J(0) = %v, want 8pi/3 = %v", got, 8*math.Pi/3)
	}
	if got := axion.JClosed(1); math.Abs(got-2*math.Pi) > 1e-14 {
		t.Fatalf("J(1) = %v, want 2pi = %v", got, 2*math.Pi)
	}
}

func TestJClosedKappaOneLimit(t *testing.T) {
	// Both-sided approach to kappa=1 converges to J(1) = 2*pi.
	// Establishes: removable singularity is handled correctly.
	for _, eps := range []float64{1e-3, 1e-6, 1e-9} {
		below := axion.JClosed(1 - eps)
		above := axion.JClosed(1 + eps)
		target := 2 * math.Pi
		if math.Abs(below/target-1) > 10*eps {
			t.Fatalf("J(1-%e) = %v, too far from 2pi", eps, below)
		}
		if math.Abs(above/target-1) > 10*eps {
			t.Fatalf("J(1+%e) = %v, too far from 2pi", eps, above)
		}
	}
}

func TestThresholdGate(t *testing.T) {
	// Case 2 (15 eV): Ea = 15(1 + v^2/2) < wp = 16.6 eV, so
	// channel is kinematically closed. Establishes: the threshold
	// gate enforces Ea > wp correctly.
	p := axion.BenchmarkParams()
	Ea, _ := axion.NRKinematics(p.Mass15, p.Va)
	if Ea.Cmp(p.Wp) > 0 {
		t.Fatalf("Ea(15) = %s > wp = %s: should be below threshold", Ea.RatString(), p.Wp.RatString())
	}
	if axion.Allowed(Ea, p.Wp) {
		t.Fatal("Allowed should be false for 15 eV axion")
	}

	// Case 1 (25 eV): above threshold.
	Ea25, _ := axion.NRKinematics(p.Mass25, p.Va)
	if Ea25.Cmp(p.Wp) <= 0 {
		t.Fatalf("Ea(25) = %s <= wp: should be above threshold", Ea25.RatString())
	}
}

func TestSkeletonSignFlipUnderNegWp(t *testing.T) {
	// The exact-rational skeleton is linear in wp; negating wp
	// flips the skeleton sign (Ea - wp stays positive when wp < 0,
	// and the skeleton prefactor flips). The physical rate path
	// (GammaNR) rejects wp < 0 as out-of-domain. This confirms
	// the algebra is correct and the domain guard is necessary.
	p := axion.BenchmarkParams()
	neg := axion.BenchmarkParams()
	neg.Wp = new(big.Rat).Neg(p.Wp)
	sPos, okPos := axion.SkeletonNR(p, p.Mass25)
	sNeg, okNeg := axion.SkeletonNR(neg, p.Mass25)
	if !okPos || !okNeg {
		t.Fatal("skeleton gate unexpectedly closed")
	}
	if sPos.Sign() <= 0 || sNeg.Sign() >= 0 {
		t.Fatalf("skeleton signs = %v, %v, want positive, negative", sPos.Sign(), sNeg.Sign())
	}
	if r, ok := axion.GammaNR(neg, p.Mass25); ok || !math.IsNaN(r) {
		t.Fatalf("rate path must refuse wp < 0: got %v (ok=%v), want NaN,false", r, ok)
	}
}

func TestUnitSlipSensitivity(t *testing.T) {
	// Using g in GeV^-1 as if it were eV^-1 (1e9 slip) changes the
	// rate by exactly 1e18; the golden detects it.
	// Establishes: the g -> eV^-1 conversion is correct and testable.
	p := axion.BenchmarkParams()
	rGood, _ := axion.GammaNR(p, p.Mass25)
	slip := axion.BenchmarkParams()
	slip.G = new(big.Rat).Mul(p.G, big.NewRat(1000000000, 1))
	rSlip, _ := axion.GammaNR(slip, p.Mass25)
	if math.Abs(rSlip/rGood/1e18-1) > 1e-12 {
		t.Fatalf("unit slip changed rate by %v, want 1e18", rSlip/rGood)
	}
	if math.Abs(rGood/1.30725026e-41-1) < 1e-6 {
		return // golden intact
	}
	t.Fatal("golden itself drifted")
}

func TestRelativisticConsistency(t *testing.T) {
	// The general closed form evaluated on the exact on-shell NR
	// pair must agree with the NR specialization up to the
	// documented O(v^2) truncation of the NR kinematics
	// (ka = m v vs m gamma v); bound 5e-6 with margin.
	// Establishes: the general and NR formulas are self-consistent.
	p := axion.BenchmarkParams()
	va, _ := p.Va.Float64()
	m := 25.0
	gamma := 1 / math.Sqrt(1-va*va)
	ea, ka := m*gamma, m*gamma*va
	g, _ := p.G.Float64()
	wp, _ := p.Wp.Float64()
	n, _ := p.N.Float64()
	gen, ok := axion.GammaGeneral(g, wp, n, ea, ka)
	if !ok {
		t.Fatal("general form closed on open channel")
	}
	nr, _ := axion.GammaNR(p, p.Mass25)
	if math.Abs(gen/nr-1) > 5e-6 {
		t.Fatalf("general vs NR mismatch: %.12e vs %.12e", gen, nr)
	}
	// Ultra-relativistic regime evaluates without error and stays
	// positive on an open channel.
	ur, ok := axion.GammaGeneral(g, wp, n, 1e6, math.Sqrt(1e12-625))
	if !ok || ur <= 0 {
		t.Fatalf("UR evaluation failed: %v %v", ur, ok)
	}
}

func TestDimensionalScaling(t *testing.T) {
	// With g fixed, scaling all energies by lam scales Gamma by
	// lam^3 (Gamma has mass dimension 1, g^2 has -2). Independent
	// dimensional-consistency check of the assembled formula.
	// Establishes: the formula has the correct mass dimension.
	p := axion.BenchmarkParams()
	g, _ := p.G.Float64()
	wp, _ := p.Wp.Float64()
	n, _ := p.N.Float64()
	r1, ok1 := axion.GammaGeneral(g, wp, n, 25.0, 0.025)
	r2, ok2 := axion.GammaGeneral(g, 2*wp, n, 50.0, 0.05)
	if !ok1 || !ok2 {
		t.Fatal("channels must be open")
	}
	if math.Abs(r2/r1/8-1) > 1e-9 {
		t.Fatalf("scaling violated: ratio %v, want 8", r2/r1)
	}
}

func TestDomainGuards(t *testing.T) {
	// Narrow domain contract: kappa is a ratio of momentum
	// magnitudes (kappa >= 0, finite), and the rate path never
	// returns ok == true with a NaN or otherwise non-plausible rate.
	// Establishes: domain guards correctly reject unphysical inputs.
	for _, bad := range []float64{-1e-9, -0.5, -2, math.NaN(), math.Inf(1), math.Inf(-1)} {
		if got := axion.JClosed(bad); !math.IsNaN(got) {
			t.Fatalf("JClosed(%v) = %v, want NaN", bad, got)
		}
	}
	p := axion.BenchmarkParams()
	g, _ := p.G.Float64()
	wp, _ := p.Wp.Float64()
	n, _ := p.N.Float64()
	for _, c := range []struct {
		name             string
		g, wp, n, ea, ka float64
	}{
		{"negative kappa (ka < 0)", g, wp, n, 25, -1e-3},
		{"negative wp", g, -wp, n, 25, 0.025},
		{"zero wp", g, 0, n, 25, 0.025},
		{"zero n", g, wp, 0, 25, 0.025},
		{"negative n", g, wp, -n, 25, 0.025},
		{"NaN g", math.NaN(), wp, n, 25, 0.025},
		{"Inf Ea", g, wp, n, math.Inf(1), 0.025},
		{"NaN ka", g, wp, n, 25, math.NaN()},
	} {
		r, ok := axion.GammaGeneral(c.g, c.wp, c.n, c.ea, c.ka)
		if ok {
			t.Fatalf("%s: ok=true with rate %v", c.name, r)
		}
		if !math.IsNaN(r) {
			t.Fatalf("%s: rate %v is not NaN", c.name, r)
		}
	}
	// The PNL-typed path refuses unphysical medium records and
	// non-positive masses the same way.
	neg := axion.BenchmarkParams()
	neg.Wp = new(big.Rat).Neg(neg.Wp)
	zeroN := axion.BenchmarkParams()
	zeroN.N = big.NewRat(0, 1)
	for _, c := range []struct {
		name string
		p    axion.Params
		m    *big.Rat
	}{
		{"wp < 0", neg, neg.Mass25},
		{"n = 0", zeroN, zeroN.Mass25},
		{"m < 0", axion.BenchmarkParams(), big.NewRat(-25, 1)},
		{"m = 0", axion.BenchmarkParams(), big.NewRat(0, 1)},
		{"nil m", axion.BenchmarkParams(), nil},
	} {
		if r, ok := axion.GammaNR(c.p, c.m); ok || !math.IsNaN(r) {
			t.Fatalf("GammaNR(%s) = %v (ok=%v), want NaN,false", c.name, r, ok)
		}
		if cr, ok := axion.CorrectionNR(c.p, c.m); ok || !math.IsNaN(cr) {
			t.Fatalf("CorrectionNR(%s) = %v (ok=%v), want NaN,false", c.name, cr, ok)
		}
	}
	// The kinematic gate is not an error: exactly zero, ok=false.
	if r, ok := axion.GammaGeneral(g, wp, n, 16.0, 0.025); ok || r != 0 {
		t.Fatalf("sub-pole channel = %v (ok=%v), want 0,false", r, ok)
	}
	// kappa = 0 (rest axion) is in domain: open channel, zero rate.
	if r, ok := axion.GammaGeneral(g, wp, n, 25, 0); !ok || r != 0 {
		t.Fatalf("rest axion = %v (ok=%v), want 0,true", r, ok)
	}
}

func TestPrefactorAssembly(t *testing.T) {
	// Independent numerical audit of the assembled prefactor under
	// the documented phase-space convention (README "Assembly
	// audit"), from the cut/spectral-density template calibrated on
	// the vacuum width in the same notation.
	//
	// The angular integral I_ang = 2 ka^2 J(kappa) is evaluated here
	// by direct Simpson quadrature over the polar angle (independent
	// of axion.JClosed and the closed form). Agreement is at
	// quadrature accuracy.
	//
	// Establishes: the prefactor assembly is correct under the
	// documented convention. Does NOT establish: that the convention
	// is the uniquely correct physical normalization (see limitation
	// prefactor_normalization).
	p := axion.BenchmarkParams()
	g, _ := p.G.Float64()
	wp, _ := p.Wp.Float64()
	n, _ := p.N.Float64()
	mb, vb := 25.0, 1e-3
	Ea := mb * (1 + vb*vb/2)
	ka := mb * vb
	wg := Ea - wp
	qg := n * wg

	const N = 200000
	h := 2.0 / N
	s := 0.0
	for i := 0; i <= N; i++ {
		x := -1.0 + float64(i)*h
		f := 2 * qg * qg * ka * ka * (1 - x*x) / (ka*ka + qg*qg - 2*ka*qg*x)
		wgt := 1.0
		if i > 0 && i < N {
			wgt = 4.0
			if i%2 == 0 {
				wgt = 2.0
			}
		}
		s += wgt * f
	}
	Iang := 2 * math.Pi * s * h / 3 // = 2 ka^2 J(kappa)

	pi4 := math.Pi * math.Pi * math.Pi * math.Pi
	assembly := (1 / (2 * Ea)) * (1 / (32 * pi4)) *
		g * g * (math.Pi * math.Pi * wp / 2) * 4 * math.Pi * n * qg * Iang
	got, ok := axion.GammaGeneral(g, wp, n, Ea, ka)
	if !ok {
		t.Fatal("assembly channel closed")
	}
	if math.Abs(assembly/got-1) > 1e-9 {
		t.Fatalf("assembly %.12e vs closed form %.12e (ratio %.12f)", assembly, got, assembly/got)
	}
}

func TestAssumptionRegister(t *testing.T) {
	// Establishes: assumption and source records are present and
	// well-formed. Does NOT establish: their physical validity.
	as := axion.Assumptions()
	if len(as) != 8 {
		t.Fatalf("assumptions = %d, want 8", len(as))
	}
	seen := map[string]bool{}
	for _, a := range as {
		if a.ID == "" || a.Statement == "" {
			t.Fatalf("assumption under-specified: %+v", a)
		}
		if seen[a.ID] {
			t.Fatalf("duplicate assumption %s", a.ID)
		}
		seen[a.ID] = true
	}
	if len(axion.Sources()) < 3 {
		t.Fatal("source register incomplete")
	}
	if axion.Sources()[0].Fidelity != "problem-definition" {
		t.Fatal("benchmark must be tagged problem-definition, not derivation source")
	}
}

func TestCorrectionBound(t *testing.T) {
	// At benchmark kinematics (kappa ~ 2e-3), the correction factor
	// C = J(kappa)/(8pi/3) should be very close to 1. The documented
	// bound is |C-1| < 1e-5.
	// Establishes: the NR approximation is valid at benchmark parameters.
	p := axion.BenchmarkParams()
	c, ok := axion.CorrectionNR(p, p.Mass25)
	if !ok {
		t.Fatal("correction unexpectedly closed")
	}
	if math.Abs(c-1) > 1e-5 {
		t.Fatalf("correction factor = %v, |C-1| = %e > 1e-5", c, math.Abs(c-1))
	}
}

func TestResponseStrength(t *testing.T) {
	// The plasmon-pole model (A62) has strength (pi/2)*wp. The rate
	// prefactor inherits this normalization.
	// Establishes: the response strength agrees with the model.
	p := axion.BenchmarkParams()
	got := axion.ResponseStrength(p)
	want := math.Pi / 2 * 16.6
	if math.Abs(got/want-1) > 1e-14 {
		t.Fatalf("response strength = %v, want %v", got, want)
	}
}

func TestNilMassHandling(t *testing.T) {
	// SkeletonNR returns (nil, false) for nil mass.
	// Establishes: nil mass does not panic.
	p := axion.BenchmarkParams()
	s, ok := axion.SkeletonNR(p, nil)
	if ok || s != nil {
		t.Fatalf("SkeletonNR(nil) = %v, %v, want nil, false", s, ok)
	}
}
