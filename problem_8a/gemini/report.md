---

### 1. Final Numerical Answers

Evaluating the non-relativistic axion-to-photon conversion rate $\Gamma_{\mathrm{NR}}$ in scientific form in $\mathrm{eV}$ units in the required order $[(1), (2)]$:

$$
\boxed{
\begin{aligned}
\text{1. } m_a = 25\,\mathrm{eV}: & \quad 1.31 \times 10^{-41}\,\mathrm{eV} \\
\text{2. } m_a = 15\,\mathrm{eV}: & \quad 0\,\mathrm{eV} \quad (\text{kinematically closed, } E_a < \omega_p)
\end{aligned}
}
$$

Ordered pair:
$$[\, 1.31 \times 10^{-41}\,\mathrm{eV},\; 0\,]$$

---

### 2. Package Architecture & Deliverables

The solution is implemented as a standalone, self-contained Type-B benchmark-solution package compliant with `phys-demo/PROTOCOL.md`:

| File | Role / Purpose |
| :--- | :--- |
| [`go.mod`](file:///home/chaschel/Desktop/biz/problem_8a_gemini/go.mod) / [`go.work`](file:///home/chaschel/Desktop/biz/problem_8a_gemini/go.work) | Module `github.com/PithomLabs/phys-demo/problem_8a/gemini` (Go 1.25.7, stdlib-only). |
| [`manifest.json`](file:///home/chaschel/Desktop/biz/problem_8a_gemini/manifest.json) | Type-B benchmark manifest (`schema_id: phys-demo/benchmark-solution-manifest/1`) with honest claims, derivation steps, unit states, solver-local computation records, and explicit justified deviation records. |
| [`manifest_test.go`](file:///home/chaschel/Desktop/biz/problem_8a_gemini/manifest_test.go) | Strict manifest contract validator testing claim honesty, deviation records, parameter units, solver-local attribution, and negative mutation drill tests. |
| [`axion/doc.go`](file:///home/chaschel/Desktop/biz/problem_8a_gemini/axion/doc.go) | Package documentation delineating demo-local solver scope from Phys core infrastructure. |
| [`axion/axion.go`](file:///home/chaschel/Desktop/biz/problem_8a_gemini/axion/axion.go) | Core physics implementation with exact-rational arithmetic (`math/big.Rat`) for parameter substitution and skeleton evaluation, closed-form angular factor $J(\kappa)$, domain guards, and declared entry points `OrderedResults` and `SkeletonNR`. |
| [`protocol/derivation_test.go`](file:///home/chaschel/Desktop/biz/problem_8a_gemini/protocol/derivation_test.go) | Executable protocol test suite exercising entry points, golden values, threshold gate, unit-slip sensitivity ($10^{18}$ shift), relativistic consistency, dimensional scaling, Simpson quadrature prefactor assembly, and domain guards. |
| [`README.md`](file:///home/chaschel/Desktop/biz/problem_8a_gemini/README.md) | Documentation containing the complete first-principles derivation, assembly audit, capability audit table, executed-vs-represented boundary, independent checks, growth-gate note, and explicit model limitations. |

---

### 3. Verification

All tests pass cleanly under race detection and zero cache:
```bash
$ go test -count=1 -race ./...
ok  	github.com/PithomLabs/phys-demo/problem_8a/gemini	1.030s
ok  	github.com/PithomLabs/phys-demo/problem_8a/gemini/protocol	1.017s
$ go vet ./...
# clean, no warnings
```
