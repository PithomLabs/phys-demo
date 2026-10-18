// Package axion is a demo-local reference artifact for TheoPhysEval
// Problem 8(a) (axion-to-photon conversion in a dielectric medium),
// living under phys-demo/problem_8a/muse (Go module
// github.com/PithomLabs/phys-demo/problem_8a/muse).
//
// Scope discipline: this package records the problem, assumptions,
// sources, derivation units, and numerical evaluation for one
// benchmark problem. It introduces no shared contracts, no PhysMath
// rules or families, and no phys-lib domain package. Physical
// assumptions stay in these demo-local records (see the growth-gate
// note in README.md), never in the theory-neutral math layer.
//
// Evidence discipline: the analytical derivation (tensor contraction,
// delta-collapse, phase-space integration) is pen-and-paper work
// recorded in README.md with independent checks; it is NOT executed
// by PhysMath rules. What this package executes is exact-rational
// arithmetic for the parameter substitution plus float64 evaluation
// of the closed-form angular factor, each with independent test
// oracles. A green suite proves the arithmetic reproduces from the
// stated formula — never that the physics is true, and never that
// PhysMath executed the derivation.
package axion
