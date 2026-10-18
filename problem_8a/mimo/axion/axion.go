package axion

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
)

// Params holds the exact-rational problem parameters in eV natural units
// (hbar = c = 1). The values are independent solver-local constructions;
// the representation package constructs the same values separately through
// physmath.NewRational and the protocol tests require agreement.
type Params struct {
	GEVInv *big.Rat // axion-photon coupling, converted to eV^-1
	OmegaP *big.Rat // plasmon pole frequency, eV
	N      *big.Rat // refractive index (dimensionless)
	VA     *big.Rat // axion velocity (dimensionless)
}

// GeVInvToEVInv converts a coupling expressed in GeV^-1 to eV^-1 exactly:
// 1 GeV = 10^9 eV, hence 1 GeV^-1 = 10^-9 eV^-1. This is a solver-local
// conversion (PROTOCOL §9.1/§9.2); unit semantics are not a Phys capability.
func GeVInvToEVInv(gGeVInv *big.Rat) *big.Rat {
	factor := new(big.Rat).SetFrac(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil))
	return new(big.Rat).Mul(gGeVInv, factor)
}

// ExactParameters returns fresh copies of the exact problem parameters:
// g = 10^-10 GeV^-1 = 10^-19 eV^-1, omega_p = 16.6 eV = 83/5 eV,
// n = 3/2, v_a = 1/1000.
func ExactParameters() Params {
	return Params{
		GEVInv: GeVInvToEVInv(new(big.Rat).SetFrac(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil))),
		OmegaP: new(big.Rat).SetFrac(big.NewInt(83), big.NewInt(5)),
		N:      new(big.Rat).SetFrac(big.NewInt(3), big.NewInt(2)),
		VA:     new(big.Rat).SetFrac(big.NewInt(1), big.NewInt(1000)),
	}
}

// GammaNRExact evaluates the non-relativistic rate exactly over the rationals
// for mass m (eV), including the threshold gate.
//
// Threshold convention: the physical condition is E_a = m (1 + v_a^2/2) >
// omega_p, while the closed NR prefactor requires omega_gamma = m - omega_p >
// 0 at leading order (it contains the factor m - omega_p). The gate is
// therefore the exact rational test m > omega_p, which implies the physical
// condition. The razor-thin window omega_p/(1+v_a^2/2) < m <= omega_p — width
// O(v_a^2 * omega_p) ~ 8 meV, where the exact rate is positive but the
// leading-order prefactor is not yet valid — is excluded by the same
// condition; no requested case lies inside it. Kinetic corrections to the
// prefactor are O(v_a^2) ~ 5e-7 relative, below displayed precision (tested
// in protocol.TestConventionInsensitivity).
func GammaNRExact(m *big.Rat) *big.Rat {
	p := ExactParameters()
	// Theta at leading NR order: gate on m > omega_p so that the factor
	// (m - omega_p) is positive and the rate cannot go negative.
	if m.Cmp(p.OmegaP) <= 0 {
		return new(big.Rat)
	}
	// Gamma_NR = g^2 omega_p n^2 (m - omega_p) m v_a^2 / 6.
	g2 := new(big.Rat).Mul(new(big.Rat).Set(p.GEVInv), new(big.Rat).Set(p.GEVInv))
	n2 := new(big.Rat).Mul(new(big.Rat).Set(p.N), new(big.Rat).Set(p.N))
	v2 := new(big.Rat).Mul(new(big.Rat).Set(p.VA), new(big.Rat).Set(p.VA))
	omegaGap := new(big.Rat).Sub(m, p.OmegaP)

	acc := new(big.Rat).Set(g2)
	acc.Mul(acc, p.OmegaP)
	acc.Mul(acc, n2)
	acc.Mul(acc, omegaGap)
	acc.Mul(acc, m)
	acc.Mul(acc, v2)
	acc.Quo(acc, big.NewRat(6, 1))
	return acc
}

// J returns the dimensionless angular factor J(kappa) in closed form,
// evaluated in float64 (solver-local). The logarithm is evaluated through
// math.Atanh (J is even; for kappa > 1, atanh(1/kappa) yields the same log),
// which avoids the catastrophic cancellation of the naive
// log((1+k)/(1-k)) form at small kappa. Endpoint conventions: J(0) is
// defined by its limit 8 pi/3, and J(1) — where the closed form is the
// indeterminate 0 * infinity — is defined by its limit 2 pi.
func J(kappa float64) float64 {
	k := math.Abs(kappa)
	switch {
	case k == 0:
		return JZero()
	case math.Abs(1-k) < 1e-15:
		return JOne()
	}
	lnRatio := 2 * math.Atanh(k)
	if k > 1 {
		lnRatio = 2 * math.Atanh(1/k)
	}
	term1 := (1 + k*k) / (2 * k * k)
	term2 := (1 - k*k) * (1 - k*k) / (4 * k * k * k) * lnRatio
	return 2 * math.Pi * (term1 - term2)
}

// JZero returns the exact kappa -> 0 limit of J, 8 pi/3 (hand-derived series
// in the package README; checked numerically in
// protocol.TestJSeriesAndLimits).
func JZero() float64 { return 8 * math.Pi / 3 }

// JOne returns the exact kappa -> 1 limit of J, 2 pi.
func JOne() float64 { return 2 * math.Pi }

// JSeries returns the small-kappa expansion J = 8 pi/3 - (8 pi/15) kappa^2
// + O(kappa^4), used as an independent check of the closed form.
func JSeries(kappa float64) float64 {
	return JZero() - (8*math.Pi/15)*kappa*kappa
}

// GammaNRGeneralFloat evaluates the general (relativistic-form) rate at
// float64 precision with the full closed-form J(kappa) and full kinematics
// E_a = m (1 + v_a^2/2). It is the independent numerical path that the
// exact-rational result is cross-checked against (PROTOCOL §9.2).
func GammaNRGeneralFloat(mEV float64) float64 {
	const (
		g  = 1e-19 // eV^-1
		wp = 16.6  // eV
		n  = 1.5
		v  = 1e-3
	)
	ea := mEV * (1 + v*v/2)
	wg := ea - wp
	if wg <= 0 {
		return 0
	}
	q := n * wg
	k := mEV * v
	kappa := k / q
	return g * g * wp * q * q * k * k * J(kappa) / (16 * math.Pi * ea * wg)
}

// GammaNRVariantFloat evaluates the NR closed form with full kinematics
// (E_a = m(1+v_a^2/2) in omega_gamma and in the prefactor, J(0) = 8 pi/3
// folded to the 1/6). It isolates the stated E_a convention for the
// convention-insensitivity test.
func GammaNRVariantFloat(mEV float64) float64 {
	const (
		g  = 1e-19
		wp = 16.6
		n  = 1.5
		v  = 1e-3
	)
	ea := mEV * (1 + v*v/2)
	wg := ea - wp
	if wg <= 0 {
		return 0
	}
	return g * g * wp * n * n * wg * ea * v * v / 6
}

// FormatSci renders an exact rational in eV as a 3-significant-figure
// scientific-notation string (display convention only; the exact value is
// the rational itself).
func FormatSci(r *big.Rat) string {
	f, _ := new(big.Float).SetPrec(64).SetRat(r).Float64()
	return strconv.FormatFloat(f, 'e', 2, 64)
}

// Result is one requested numerical output, in problem order.
type Result struct {
	Case           string   // "(1)" or "(2)", problem.md output order
	MassEV         *big.Rat // axion mass, eV
	GammaEV        *big.Rat // exact conversion rate, eV
	Display        string   // 3-s.f. scientific display, eV
	AboveThreshold bool     // gate m > omega_p (implies physical E_a > omega_p)
}

// OrderedResults returns the two requested numerical outputs in problem.md
// order: (1) m_a = 25 eV, (2) m_a = 15 eV, both at v_a = 10^-3. This is the
// manifest's primary entry point.
func OrderedResults() ([]Result, error) {
	cases := []int64{25, 15}
	out := make([]Result, 0, len(cases))
	p := ExactParameters()
	for i, mEV := range cases {
		if mEV <= 0 {
			return nil, fmt.Errorf("axion: case %d mass must be positive", i+1)
		}
		m := big.NewRat(mEV, 1)
		gamma := GammaNRExact(m)
		display := "0 eV"
		if gamma.Sign() != 0 {
			display = FormatSci(gamma) + " eV"
		}
		out = append(out, Result{
			Case:           fmt.Sprintf("(%d)", i+1),
			MassEV:         m,
			GammaEV:        gamma,
			Display:        display,
			AboveThreshold: m.Cmp(p.OmegaP) > 0,
		})
	}
	return out, nil
}
