package axion

import (
	"math"
	"math/big"
)

// ---------------------------------------------------------------------------
// Problem definition records (demo-local, not PhysMath/PhysLib types)
// ---------------------------------------------------------------------------

// Assumption records one physical assumption.
type Assumption struct {
	ID        string
	Statement string
}

// Assumptions returns the documented problem assumptions.
func Assumptions() []Assumption {
	return []Assumption{
		{ID: "medium", Statement: "Conventional solid-state dielectric, non-magnetic, zero temperature."},
		{ID: "coupling", Statement: "Axion couples only via L = -(g/4) a F Fdual; no other axion interactions."},
		{ID: "units", Statement: "Natural units hbar = c = 1; all quantities in eV."},
		{ID: "response", Statement: "Dielectric response given by Im[-1/eps_L(w,q)] = (pi/2) wp delta(w - wp) (A62)."},
		{ID: "photon", Statement: "Outgoing transverse photon with |q| = n w in the given-n approximation."},
		{ID: "recoil", Statement: "Recoil on the medium neglected (infinite-medium limit)."},
		{ID: "tzero", Statement: "Zero-temperature limit: no thermal population of plasmons; conversion requires Ea > wp."},
		{ID: "nrlimit", Statement: "Non-relativistic axion: v_a = 1e-3, kinematic truncation at O(v^2)."},
	}
}

// Source records one literature reference.
type Source struct {
	Cite        string
	Establishes string
	Fidelity    string
}

// Sources returns the documented source references.
func Sources() []Source {
	return []Source{
		{Cite: "arXiv:2610.16733v1", Establishes: "problem definition only", Fidelity: "problem-definition"},
		{Cite: "Caputo-Millar-Vitagliano, arXiv:2005.00078", Establishes: "self-energy to rate method; narrow-width delta pattern", Fidelity: "derivation-method"},
		{Cite: "Knapen-Kozaczuk-Lin, arXiv:2101.08275", Establishes: "density response to loss-function identity; loss function as inclusive excitation rate", Fidelity: "derivation-method"},
	}
}

// ---------------------------------------------------------------------------
// Parameters
// ---------------------------------------------------------------------------

// Params holds the exact-rational benchmark parameters.
type Params struct {
	G      *big.Rat // axion-photon coupling in eV^-1
	Wp     *big.Rat // plasma frequency in eV
	N      *big.Rat // refractive index (dimensionless)
	Va     *big.Rat // axion velocity (dimensionless)
	Mass25 *big.Rat // 25 eV case
	Mass15 *big.Rat // 15 eV case
}

func ratFromString(s string) *big.Rat {
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		panic("axion: bad rational literal " + s)
	}
	return r
}

// BenchmarkParams returns the exact-decimal benchmark parameters.
// g = 1e-10 GeV^-1 = 1e-10 * 1e-9 eV^-1 = 1e-19 eV^-1.
func BenchmarkParams() Params {
	return Params{
		G:      ratFromString("1e-19"),
		Wp:     ratFromString("16.6"),
		N:      ratFromString("1.5"),
		Va:     ratFromString("0.001"),
		Mass25: ratFromString("25"),
		Mass15: ratFromString("15"),
	}
}

// ---------------------------------------------------------------------------
// Rate evaluation
//
// Derived closed form (README.md, hand derivation with independent
// checks and an assembly audit): for on-shell axion (Ea, ka) with
// Ea > wp > 0, n > 0, ka >= 0,
//
//   Gamma = g^2 wp qg^2 ka^2 J(k) / (16 pi Ea wg),
//
// wg = Ea-wp, qg = n wg, k = ka/qg,
// J(k) = 2pi[(1+k^2)/2k^2 - ((1-k^2)^2/4k^3) ln((1+k)/|1-k|)].
// For Ea <= wp the channel is kinematically closed and the rate is
// exactly 0 (T=0 detailed balance: the medium can only absorb energy).
//
// Photon-normalization convention (documented in README.md "Assembly
// audit"): the outgoing transverse photon's cut is normalized as the
// free massless 2 pi delta(q^2) with the 1/(2|q|) Jacobian, with only
// the on-shell relation replaced by the given-n dispersion w = |q|/n.
// The stated model does not fix the in-medium spectral-density
// normalization; the residual dependence is one power of n and is
// recorded as the manifest limitation "prefactor_normalization".
//
// NR specialization (J(0) = 8pi/3 exactly):
//
//   Gamma_NR = g^2 wp qg^2 ka^2 / (6 Ea wg)   [exact rational]
//
// with Ea = m(1+v^2/2), ka = m v. The full NR value is
// Gamma_NR * C, C = J(k)/(8pi/3), evaluated in float64 with a
// documented remainder bound |C-1| < 1e-5 at benchmark kinematics.
// ---------------------------------------------------------------------------

// band-limit helpers over *big.Rat.
func ratMul(xs ...*big.Rat) *big.Rat {
	acc := new(big.Rat).SetInt64(1)
	for _, x := range xs {
		acc.Mul(acc, x)
	}
	return acc
}
func ratQuo(a, b *big.Rat) *big.Rat { return new(big.Rat).Quo(a, b) }

// NRKinematics returns (Ea, ka) under the documented O(v^2)
// non-relativistic truncation: p = m v, E = m + p^2/2m.
func NRKinematics(m, v *big.Rat) (Ea, ka *big.Rat) {
	v2 := ratMul(v, v)
	half := new(big.Rat).SetFrac64(1, 2)
	Ea = new(big.Rat).Add(m, ratMul(m, v2, half))
	ka = ratMul(m, v)
	return Ea, ka
}

// Allowed reports whether conversion is kinematically possible.
func Allowed(Ea, wp *big.Rat) bool { return Ea.Cmp(wp) > 0 }

// SkeletonNR evaluates the exact-rational NR-limit skeleton
// g^2 wp qg^2 ka^2 / (6 Ea wg). No pi, no log: fully exact.
// This is exact-rational algebra with only the kinematic gate
// Ea > wp applied: it performs no physical-domain validation
// (wp > 0, n > 0, m > 0 are enforced by GammaNR/CorrectionNR).
// A nil mass yields (nil, false).
func SkeletonNR(p Params, m *big.Rat) (*big.Rat, bool) {
	if m == nil {
		return nil, false
	}
	Ea, ka := NRKinematics(m, p.Va)
	if !Allowed(Ea, p.Wp) {
		return new(big.Rat), false
	}
	wg := new(big.Rat).Sub(Ea, p.Wp)
	qg := ratMul(p.N, wg)
	num := ratMul(p.G, p.G, p.Wp, qg, qg, ka, ka)
	den := ratMul(new(big.Rat).SetInt64(6), Ea, wg)
	return ratQuo(num, den), true
}

// JClosed evaluates the closed-form angular factor J(kappa) in
// float64. Pure numerical evaluation, not a PhysMath execution.
//
// Domain: kappa = ka/qg is a ratio of momentum magnitudes and is
// therefore >= 0 for this model. Negative or non-finite input is out
// of domain and returns NaN: no finite value may be presented as a
// plausible physical result. kappa == 0 gives 8*pi/3 exactly;
// kappa == 1 is a removable singularity of the closed form and takes
// its exact limit 2*pi.
func JClosed(kappa float64) float64 {
	if math.IsNaN(kappa) || math.IsInf(kappa, 0) || kappa < 0 {
		return math.NaN()
	}
	if kappa == 0 {
		return 8 * math.Pi / 3
	}
	if kappa == 1 {
		return 2 * math.Pi
	}
	k2 := kappa * kappa
	return 2 * math.Pi * ((1+k2)/(2*k2) -
		((1-k2)*(1-k2)/(4*k2*kappa))*math.Log((1+kappa)/math.Abs(1-kappa)))
}

// CorrectionNR returns J(k)/(8pi/3) for the NR kinematics. Same
// domain contract as GammaNR: out-of-domain or non-finite inputs
// return (NaN, false).
func CorrectionNR(p Params, m *big.Rat) (float64, bool) {
	if !validMedium(p) || m == nil || m.Sign() <= 0 {
		return math.NaN(), false
	}
	Ea, ka := NRKinematics(m, p.Va)
	if !Allowed(Ea, p.Wp) {
		return 0, false
	}
	wg := new(big.Rat).Sub(Ea, p.Wp)
	qg := ratMul(p.N, wg)
	kf, _ := new(big.Rat).Quo(ka, qg).Float64()
	return JClosed(kf) / (8 * math.Pi / 3), true
}

// validMedium reports whether the medium/parameter records are inside
// the documented physical domain of the model: all fields present,
// wp > 0, n > 0. (*big.Rat values are always finite, so finiteness
// needs no check here.) SkeletonNR is deliberately outside this
// contract: it is exact-rational algebra, linear in wp, and callers
// (GammaNR/CorrectionNR) enforce the physical domain.
func validMedium(p Params) bool {
	if p.G == nil || p.Wp == nil || p.N == nil || p.Va == nil {
		return false
	}
	return p.Wp.Sign() > 0 && p.N.Sign() > 0
}

// GammaNR returns the full NR rate in eV (float64) and whether the
// channel is open. Case (2) is closed: returns (0, false).
// Out-of-domain or non-finite inputs return (NaN, false): never a
// plausible numeric rate with ok == true.
func GammaNR(p Params, m *big.Rat) (float64, bool) {
	if !validMedium(p) || m == nil || m.Sign() <= 0 {
		return math.NaN(), false
	}
	s, ok := SkeletonNR(p, m)
	if !ok {
		return 0, false
	}
	c, ok := CorrectionNR(p, m)
	if !ok {
		return 0, false
	}
	sf, _ := s.Float64()
	return sf * c, true
}

// GammaGeneral evaluates the general (Ea, ka) closed form in float64.
// Caller supplies an on-shell pair; the threshold gate applies.
//
// Domain contract (narrow, documented): all inputs finite, wp > 0,
// n > 0, ka >= 0 (hence kappa = ka/qg >= 0), and Ea > wp for an open
// channel. Exactly at or below threshold the channel is kinematically
// closed and the function returns (0, false). Any out-of-domain or
// non-finite input, and any assembled rate that is not finite,
// returns (NaN, false): the function never returns ok == true with a
// NaN or otherwise non-plausible rate.
func GammaGeneral(g, wp, n, Ea, ka float64) (float64, bool) {
	for _, v := range []float64{g, wp, n, Ea, ka} {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return math.NaN(), false
		}
	}
	if wp <= 0 || n <= 0 || ka < 0 {
		return math.NaN(), false
	}
	if !(Ea > wp) {
		return 0, false
	}
	wg := Ea - wp
	qg := n * wg
	kappa := ka / qg
	J := JClosed(kappa)
	rate := g * g * wp * qg * qg * ka * ka * J / (16 * math.Pi * Ea * wg)
	if math.IsNaN(rate) || math.IsInf(rate, 0) {
		return math.NaN(), false
	}
	return rate, true
}

// ResponseStrength returns the (pi/2)*wp prefactor of the
// plasmon-pole response model (A62). Pinned by tests: the rate
// prefactor inherits this normalization.
func ResponseStrength(p Params) float64 {
	wp, _ := p.Wp.Float64()
	return math.Pi / 2 * wp
}

// RatString renders a *big.Rat in scientific form with the given
// significant digits (for the completion report, not a proof).
func RatString(r *big.Rat, digits int) string {
	f := new(big.Float).SetPrec(128).SetRat(r)
	return f.Text('e', digits)
}

// CheckResult bundles one mass case outcome in required order.
type CheckResult struct {
	MassEV string
	Open   bool
	RateEV float64
}

// OrderedResults returns [(1) 25 eV, (2) 15 eV] as required.
func OrderedResults(p Params) [2]CheckResult {
	out := [2]CheckResult{}
	for i, m := range []*big.Rat{p.Mass25, p.Mass15} {
		rate, ok := GammaNR(p, m)
		out[i] = CheckResult{MassEV: m.RatString(), Open: ok, RateEV: rate}
	}
	return out
}
