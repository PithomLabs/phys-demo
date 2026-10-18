// Package protocol verifies the Problem 8(a) demo artifact: exact
// parameter substitution, independent goldens with anti-goldens,
// threshold behavior, ordering, and honest evidence labeling. A
// green suite means the arithmetic reproduces from the stated
// formula with independently checked values — never that the
// physics is proved or that any engine executed the derivation.
package protocol

import (
	"math"
	"math/big"
	"testing"

	axion "github.com/PithomLabs/phys-demo/problem_8a/muse-deepseek/axion"
)

func TestSkeletonExactRational(t *testing.T) {
	// Independent recomputation of the NR-limit skeleton from the
	// closed form, grouped differently from the implementation
	// (single fraction, direct decimal literals): exact big.Rat
	// equality, no tolerance.
	p := axion.BenchmarkParams()
	got, ok := axion.SkeletonNR(p, p.Mass25)
	if !ok {
		t.Fatal("25 eV channel closed")
	}
	g, _ := new(big.Rat).SetString("1e-19")
	g2 := new(big.Rat).Mul(g, g)
	m := big.NewRat(25, 1)
	v := big.NewRat(1, 1000)
	v2 := new(big.Rat).Mul(v, v)
	ea := new(big.Rat).Add(m, new(big.Rat).Mul(m, new(big.Rat).Mul(v2, big.NewRat(1, 2))))
	ka := new(big.Rat).Mul(m, v)
	wp := big.NewRat(166, 10)
	wg := new(big.Rat).Sub(ea, wp)
	qg := new(big.Rat).Mul(big.NewRat(3, 2), wg)
	num := new(big.Rat).Mul(g2, wp)
	num.Mul(num, new(big.Rat).Mul(qg, qg))
	num.Mul(num, new(big.Rat).Mul(ka, ka))
	den := new(big.Rat).Mul(big.NewRat(6, 1), new(big.Rat).Mul(ea, wg))
	want := new(big.Rat).Quo(num, den)
	if got.Cmp(want) != 0 {
		t.Fatalf("skeleton mismatch:\n got %s\nwant %s", got.RatString(), want.RatString())
	}
}

func TestHeadlineNumbers(t *testing.T) {
	// Required output order [(1),(2)]: 25 eV open, 15 eV closed.
	p := axion.BenchmarkParams()
	res := axion.OrderedResults(p)
	if res[0].MassEV != "25" || res[1].MassEV != "15" {
		t.Fatalf("output ordering wrong: %+v", res)
	}
	if !res[0].Open {
		t.Fatal("25 eV channel must be open")
	}
	if res[1].Open || res[1].RateEV != 0 {
		t.Fatalf("15 eV channel must be closed with rate 0, got %+v", res[1])
	}
	// 1.30725026e-41 eV from the exact skeleton times the bounded
	// correction; tolerance 1e-6 relative covers float64 rounding
	// with three orders of margin over the 7.9e-7 kappa correction.
	want := 1.30725026e-41
	if math.Abs(res[0].RateEV/want-1) > 1e-6 {
		t.Fatalf("Gamma(25 eV) = %.8e, want %.8e", res[0].RateEV, want)
	}
	if res[0].RateEV <= 0 {
		t.Fatal("physical rate must be positive")
	}
	// Correction bound actually honored at benchmark kinematics.
	c, _ := axion.CorrectionNR(p, p.Mass25)
	if math.Abs(c-1) > 1e-5 {
		t.Fatalf("finite-kappa correction %v exceeds bound", c)
	}
}

func TestAngularFactorQuadrature(t *testing.T) {
	// Independent Simpson check of the closed-form J(kappa) at
	// representative values including the UR value 2/3, the
	// removable-singularity value kappa = 1 (with the single
	// removable point evaluated by its limit), and the benchmark
	// kinematics kappa ~ 0.001984.
	simpson := func(k float64) float64 {
		const N = 20000
		h := 2.0 / N
		s := 0.0
		for i := 0; i <= N; i++ {
			x := -1.0 + float64(i)*h
			den := 1 + k*k - 2*k*x
			f := 0.0
			if den == 0 {
				// Removable point x = (1+k^2)/(2k), inside [-1,1]
				// only at |k| = 1; limit of (1-x^2)/den is x/k.
				f = x / k
			} else {
				f = (1 - x*x) / den
			}
			w := 1.0
			if i > 0 && i < N {
				w = 4.0
				if i%2 == 0 {
					w = 2.0
				}
			}
			s += w * f
		}
		return 2 * math.Pi * s * h / 3
	}
	// Benchmark kinematics (exact on-shell NR pair, benchmark medium).
	mb, vb, wpb, nb := 25.0, 1e-3, 16.6, 1.5
	Eab := mb * (1 + vb*vb/2)
	kab := mb * vb
	kbench := kab / (nb * (Eab - wpb))
	for _, k := range []float64{0.01, 0.1, 0.5, 2.0 / 3.0, 1.0, 1.5, kbench} {
		got, want := axion.JClosed(k), simpson(k)
		if math.Abs(got/want-1) > 1e-6 {
			t.Fatalf("J(%v): closed %v vs quadrature %v", k, got, want)
		}
	}
	if math.Abs(axion.JClosed(0)/(8*math.Pi/3)-1) > 1e-15 {
		t.Fatal("J(0) != 8pi/3")
	}
}

func TestJClosedKappaOneLimit(t *testing.T) {
	// The closed form is indeterminate at kappa = 1 (the divergent
	// log coefficient cancels the log argument's zero). The limit is
	// 2*pi; the exact branch must return it, and the one-sided
	// approaches must be finite and first-order close:
	// J(1 -+ delta) = 2pi(1 -+ delta) + O(delta^2 ln(1/delta)).
	if got := axion.JClosed(1); got != 2*math.Pi {
		t.Fatalf("J(1) = %v, want exactly 2pi", got)
	}
	for _, d := range []float64{1e-3, 1e-4, 1e-5, 1e-6} {
		for _, k := range []float64{1 - d, 1 + d} {
			got := axion.JClosed(k)
			if math.IsNaN(got) || math.IsInf(got, 0) {
				t.Fatalf("J(%v) = %v, want finite", k, got)
			}
			if math.Abs(got-2*math.Pi) > 4*math.Pi*d {
				t.Fatalf("J(%v) = %v off 2pi by %v (bound %v)", k, got, got-2*math.Pi, 4*math.Pi*d)
			}
		}
	}
}

func TestThresholdBoundary(t *testing.T) {
	// Exactly at the pole the channel closes (Ea == wp fails the
	// strict gate); just above, the rate turns on continuously
	// (~q^4), not discontinuously.
	p := axion.BenchmarkParams()
	// Exactly at the pole mass the kinetic energy still opens the
	// channel (Ea = m(1+v^2/2) > wp); the turn-on is continuous
	// (~q^4), while clearly below the pole it is exactly closed.
	at, _ := new(big.Rat).SetString("16.6")
	rAt, openAt := axion.GammaNR(p, at)
	if !openAt || rAt <= 0 {
		t.Fatal("pole-mass channel must be open via kinetic energy")
	}
	r25, _ := axion.GammaNR(p, p.Mass25)
	if !(rAt < r25*1e-6) {
		t.Fatal("near-threshold rate must turn on continuously, far below the 25 eV rate")
	}
	below, _ := new(big.Rat).SetString("16.5")
	if r, ok := axion.GammaNR(p, below); ok || r != 0 {
		t.Fatalf("sub-pole channel must be closed with rate 0, got %v %v", r, ok)
	}
}

func TestResponseSignSensitivity(t *testing.T) {
	// Physical part: the plasmon-pole model's response strength
	// (pi/2)*wp is strictly positive, and an open channel has a
	// strictly positive rate.
	p := axion.BenchmarkParams()
	wp, _ := p.Wp.Float64()
	if got := axion.ResponseStrength(p); got != math.Pi/2*wp || !(got > 0) {
		t.Fatalf("response strength = %v, want (pi/2)*wp > 0", got)
	}
	r, ok := axion.GammaNR(p, p.Mass25)
	if !ok || !(r > 0) {
		t.Fatalf("open-channel rate = %v (ok=%v), want strictly positive", r, ok)
	}
	// Unphysical signed-input algebraic sensitivity check: negating
	// wp is outside the model domain (the T=0 plasmon-pole response
	// requires wp > 0), and the guarded rate path must refuse it
	// rather than emit a negative "rate". The exact-rational
	// skeleton is pure algebra: its prefactor is linear in wp and
	// its gap Ea - wp stays positive when wp is negated (Ea > 0 >
	// wp), so the skeleton changes sign exactly. No positivity
	// validator exists or is claimed anywhere in this package.
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
	// the vacuum width in the same notation:
	//
	//   Gamma = (1/2Ea)(1/2) integral d^4p/(2pi)^4
	//             [2 pi delta(q^2)] [-2 Im D00(p)] g^2 |q x p|^2
	//   free-photon cut:  2 pi delta(q^2)|_+ = (pi/|q|) delta(Ea - p0 - |q|/n)
	//   medium weight:    -2 Im D00 = (2/|p|^2)(pi/2) wp delta(p0 - wp)
	//
	// Collapsing p0 (which removes the medium delta) and then the
	// radial q integral (integral q dq delta = n qg) leaves
	//
	//   Gamma = (1/2Ea) [1/(32 pi^4)] g^2 (pi^2 wp/2) 4 pi n qg I_ang,
	//   I_ang = integral dOmega (2/|p|^2) |q x p|^2 = 2 ka^2 J(kappa),
	//
	// with I_ang evaluated here by direct Simpson quadrature over the
	// polar angle (independent of axion.JClosed and of the closed
	// form). Agreement is at quadrature accuracy.
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
