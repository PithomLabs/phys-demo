// Package axion records the solver-local physics of TheoPhysEval Problem 8(a)
// (axion-to-photon conversion in a dielectric medium, arXiv:2610.16733v1,
// App. A eq. (A62)).
//
// Boundary discipline (phys-demo/PROTOCOL.md §9.2): this package is plain
// solver-local Go over math/big and math. It deliberately imports NO Phys
// library. Exact-rational representation of the same parameters through
// physmath and the replayed substitution derivation live in the sibling
// representation package; the protocol tests cross-check the two independent
// paths. Nothing computed here is attributed to PhysMath or PhysLib.
//
// Physics (independent derivation from problem.md (A62); full derivation in
// the package README):
//
//	Gamma(E_a,k_a) = g^2 * omega_p * q_gamma^2 * k_a^2 / (16 pi E_a omega_gamma)
//	                 * J(kappa) * Theta(E_a - omega_p),
//	q_gamma = n * omega_gamma, omega_gamma = E_a - omega_p, kappa = k_a/q_gamma,
//	J(kappa) = 2 pi [ (1+kappa^2)/(2 kappa^2)
//	                 - (1-kappa^2)^2/(4 kappa^3) ln((1+kappa)/|1-kappa|) ].
//
// Non-relativistic specialization (E_a = m at leading order, J(0) = 8 pi/3;
// the pi's cancel exactly):
//
//	Gamma_NR = g^2 * omega_p * n^2 * (m - omega_p) * m * v_a^2 / 6,
//
// which is pure rational arithmetic given the exact rational inputs.
package axion
