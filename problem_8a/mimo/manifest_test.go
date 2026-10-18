package mimo

// manifest_test.go is this package's whole conformance-enforcement story
// (PROTOCOL 11.3 recorded limitation: no shared Type-B checker exists).
//
// It performs, from the manifest.json contract (schema
// phys-demo/benchmark-solution-manifest/1):
//
//   - strict bidirectional schema validation (unknown keys rejected,
//     required keys and enums enforced per PROTOCOL 11.2);
//   - infrastructure_use adoption proofs (PROTOCOL 8.1/8.2): every recorded
//     test exists, reaches the claimed API symbol through the call graph of
//     this module, and every claimed code_path file/symbol exists;
//   - entry-point exercise proofs over protocol/derivation_test.go
//     (PROTOCOL 12.1/12.2, static call and non-vacuity checks);
//   - the category firewall (ADR 3): formalization claims require
//     formalization_refs;
//   - the computed conformance status (PROTOCOL 11.3): CONFORMANT requires
//     schema validity, exercised entry points, evidence at a stated level,
//     zero unresolved mandatory gaps, and no firewall violation. Status is
//     computed here, never declared by the manifest.
//
// What this test establishes: schema validity, honest bookkeeping, adoption
// traceability, and rejection of deliberate corruptions. It never
// establishes physical truth or the correctness of solver-local numbers
// (PROTOCOL 12, "What tests never establish").

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const (
	manifestFile          = "manifest.json"
	schemaIDExact         = "phys-demo/benchmark-solution-manifest/1"
	packageTypeExact      = "benchmark-solution"
	executedBySolverLocal = "solver-local"
	// minProtocolAssertions is the non-vacuity floor for
	// protocol/derivation_test.go (PROTOCOL 12.2): the protocol test must
	// contain at least this many test-failure assertions so that a stub
	// package can never count as verification.
	minProtocolAssertions = 8
)

// closedRuleRegistry is the exact nine-rule list of PROTOCOL 5 (live registry
// phys-math rules.go/families.go). Any executed_rules entry outside this set
// is fabricated and MUST fail validation (PROTOCOL 12.4).
var closedRuleRegistry = map[string]bool{
	"EXPR-CANON-001":       true,
	"EXPR-SIMPLIFY-001":    true,
	"EXPR-SUBSTITUTE-001":  true,
	"RAT-NORMALIZE-001":    true,
	"DIFFERENTIATE-001":    true,
	"TOTAL-DERIVATIVE-001": true,
	"EXPR-COMBINE-001":     true,
	"SELECT-BRANCH-001":    true,
	"EXPR-APPLY-001":       true,
}

// evidenceLevels is the PROTOCOL 10 evidence-level vocabulary (six levels).
var evidenceLevels = map[string]bool{
	"Represented":            true,
	"Structurally validated": true,
	"Executed":               true,
	"Replayed/verified":      true,
	"Physically interpreted": true,
	"Empirically supported":  true,
}

// supportStates is the PROTOCOL 5 capability vocabulary plus UNSUPPORTED, the
// four values allowed for derivations[].support by PROTOCOL 11.2.
var supportStates = map[string]bool{
	"SUPPORTED":             true,
	"PARTIALLY SUPPORTED":   true,
	"REPRESENTATIONAL ONLY": true,
	"UNSUPPORTED":           true,
}

var libraryNames = map[string]bool{"phys-artifact": true, "phys-math": true, "phys-lib": true}
var gapSeverities = map[string]bool{"mandatory": true, "advisory": true}

var (
	packageIDPattern     = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	apiSymbolPattern     = regexp.MustCompile(`^[a-z-]+\.[A-Za-z0-9_]+$`)
	assumptionIDPattern  = regexp.MustCompile(`^[a-z0-9_]+$`)
	formalizationPattern = regexp.MustCompile(`(?i)formal`)
)

// --- manifest schema structs (exactly the keys of PROTOCOL 11.2; the
// strict decoder below rejects anything else) ---

type problemBlock struct {
	Benchmark        string   `json:"benchmark"`
	Location         string   `json:"location"`
	Fidelity         string   `json:"fidelity"`
	RequestedOutputs []string `json:"requested_outputs"`
}

type parameterEntry struct {
	Symbol     string `json:"symbol"`
	Value      string `json:"value"`
	UnitState  string `json:"unit_state"`
	Unit       string `json:"unit"`
	Convention string `json:"convention"`
	Note       string `json:"note"`
}

type assumptionEntry struct {
	ID        string `json:"id"`
	Statement string `json:"statement"`
}

type capabilityAudit struct {
	Supported            []string `json:"supported"`
	Partial              []string `json:"partial"`
	RepresentationalOnly []string `json:"representational_only"`
	Deferred             []string `json:"deferred"`
	Unavailable          []string `json:"unavailable"`
}

type derivationEntry struct {
	ID            string   `json:"id"`
	Support       string   `json:"support"`
	ExecutedRules []string `json:"executed_rules"`
	Reason        string   `json:"reason"`
}

type infraUseEntry struct {
	Concern     string `json:"concern"`
	Library     string `json:"library"`
	APISymbol   string `json:"api_symbol"`
	CodePath    string `json:"code_path"`
	Test        string `json:"test"`
	Result      string `json:"result"`
	Limitations string `json:"limitations"`
}

type deviationEntry struct {
	CapabilitySought string `json:"capability_sought"`
	APIExamined      string `json:"api_examined"`
	Reason           string `json:"reason"`
	LocalAlternative string `json:"local_alternative"`
}

type solverLocalEntry struct {
	ID                string `json:"id"`
	Algorithm         string `json:"algorithm"`
	Conventions       string `json:"conventions"`
	Inputs            string `json:"inputs"`
	Outputs           string `json:"outputs"`
	IndependentChecks string `json:"independent_checks"`
	ExecutedBy        string `json:"executed_by"`
}

type resultEntry struct {
	ID            string `json:"id"`
	Value         string `json:"value"`
	EvidenceLevel string `json:"evidence_level"`
	Evidence      string `json:"evidence"`
}

type gapEntry struct {
	GapID       string `json:"gap_id"`
	Requirement string `json:"requirement"`
	Reason      string `json:"reason"`
	Severity    string `json:"severity"`
}

type formalizationRefEntry struct {
	Path        string `json:"path"`
	ArtifactRef string `json:"artifact_ref"`
}

type entryPointsBlock struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
}

type manifest struct {
	SchemaID                string                  `json:"schema_id"`
	PackageType             string                  `json:"package_type"`
	PackageID               string                  `json:"package_id"`
	Problem                 problemBlock            `json:"problem"`
	Parameters              []parameterEntry        `json:"parameters"`
	Regime                  string                  `json:"regime"`
	Conventions             map[string]string       `json:"conventions"`
	Assumptions             []assumptionEntry       `json:"assumptions"`
	CapabilityAudit         capabilityAudit         `json:"capability_audit"`
	Derivations             []derivationEntry       `json:"derivations"`
	InfrastructureUse       []infraUseEntry         `json:"infrastructure_use"`
	Deviations              []deviationEntry        `json:"deviations"`
	SolverLocalComputations []solverLocalEntry      `json:"solver_local_computations"`
	Results                 []resultEntry           `json:"results"`
	Limitations             []string                `json:"limitations"`
	OpenObligations         []string                `json:"open_obligations"`
	DeclaredGaps            []gapEntry              `json:"declared_gaps"`
	FormalizationRefs       []formalizationRefEntry `json:"formalization_refs"`
	EntryPoints             entryPointsBlock        `json:"entry_points"`
}

// requiredTopLevelKeys are the PROTOCOL 11.2 top-level keys; formalization_refs
// is the single optional key.
var requiredTopLevelKeys = []string{
	"schema_id", "package_type", "package_id", "problem", "parameters", "regime",
	"conventions", "assumptions", "capability_audit", "derivations",
	"infrastructure_use", "deviations", "solver_local_computations", "results",
	"limitations", "open_obligations", "declared_gaps", "entry_points",
}

var optionalTopLevelKeys = map[string]bool{"formalization_refs": true}

// loadManifestRaw reads manifest.json bytes.
func loadManifestRaw(t *testing.T) []byte {
	t.Helper()
	raw, err := os.ReadFile(manifestFile)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	return raw
}

// decodeManifestStrict parses with unknown fields rejected (PROTOCOL 11.2
// strict JSON parsing).
func decodeManifestStrict(raw []byte) (manifest, map[string]json.RawMessage, error) {
	var m manifest
	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		return manifest{}, nil, fmt.Errorf("strict decode: %w", err)
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		return manifest{}, nil, fmt.Errorf("top-level parse: %w", err)
	}
	return m, top, nil
}

func loadManifest(t *testing.T) manifest {
	t.Helper()
	m, _, err := decodeManifestStrict(loadManifestRaw(t))
	if err != nil {
		t.Fatalf("decode manifest: %v", err)
	}
	return m
}

// checkTopLevelKeys enforces the bidirectional exact key set of PROTOCOL 11.2.
func checkTopLevelKeys(top map[string]json.RawMessage) []string {
	var errs []string
	allowed := map[string]bool{}
	for _, k := range requiredTopLevelKeys {
		allowed[k] = true
	}
	for k := range optionalTopLevelKeys {
		allowed[k] = true
	}
	for k := range top {
		if !allowed[k] {
			errs = append(errs, fmt.Sprintf("unknown top-level key %q rejected", k))
		}
	}
	for _, k := range requiredTopLevelKeys {
		if _, ok := top[k]; !ok {
			errs = append(errs, fmt.Sprintf("required top-level key %q missing", k))
		}
	}
	return errs
}

// validateManifest enforces every schema rule of PROTOCOL 11.2 that does not
// require reading other files; testNames (when non-nil) is the set of
// executable Test functions in this package for infrastructure_use[].test.
func validateManifest(m manifest, testNames map[string]bool) error {
	var errs []string
	add := func(format string, args ...interface{}) {
		errs = append(errs, fmt.Sprintf(format, args...))
	}

	if m.SchemaID != schemaIDExact {
		add("schema_id %q, want %q", m.SchemaID, schemaIDExact)
	}
	if m.PackageType != packageTypeExact {
		add("package_type %q, want %q", m.PackageType, packageTypeExact)
	}
	if m.PackageID == "" || !packageIDPattern.MatchString(m.PackageID) {
		add("package_id %q is not a stable kebab id", m.PackageID)
	}

	if m.Problem.Benchmark == "" || m.Problem.Location == "" || m.Problem.Fidelity == "" {
		add("problem requires nonempty benchmark, location, fidelity")
	}
	if len(m.Problem.RequestedOutputs) == 0 {
		add("problem.requested_outputs must be nonempty")
	}
	for i, o := range m.Problem.RequestedOutputs {
		if strings.TrimSpace(o) == "" {
			add("problem.requested_outputs[%d] empty", i)
		}
	}

	if len(m.Parameters) == 0 {
		add("parameters must be nonempty")
	}
	for i, p := range m.Parameters {
		if p.Symbol == "" || p.Value == "" {
			add("parameters[%d] requires symbol and value", i)
		}
		switch p.UnitState {
		case "explicit":
			if p.Unit == "" {
				add("parameters[%d] unit_state explicit requires unit", i)
			}
			if p.Convention != "" || p.Note != "" {
				add("parameters[%d] explicit must not carry convention/note", i)
			}
		case "dimensionless":
			if p.Unit != "" || p.Convention != "" || p.Note != "" {
				add("parameters[%d] dimensionless must not carry unit/convention/note", i)
			}
		case "natural-units":
			if p.Convention == "" {
				add("parameters[%d] unit_state natural-units requires convention", i)
			}
			if p.Unit != "" || p.Note != "" {
				add("parameters[%d] natural-units must not carry unit/note", i)
			}
		case "unresolved":
			if p.Note == "" {
				add("parameters[%d] unit_state unresolved requires note", i)
			}
			if p.Unit != "" || p.Convention != "" {
				add("parameters[%d] unresolved must not carry unit/convention", i)
			}
		default:
			add("parameters[%d] invalid unit_state %q", i, p.UnitState)
		}
	}

	if strings.TrimSpace(m.Regime) == "" {
		add("regime must be nonempty")
	}
	if len(m.Conventions) == 0 {
		add("conventions must record the interaction/dispersion/measure conventions")
	}
	for k, v := range m.Conventions {
		if strings.TrimSpace(k) == "" || strings.TrimSpace(v) == "" {
			add("conventions[%q] must be nonempty", k)
		}
	}

	if len(m.Assumptions) == 0 {
		add("assumptions must be nonempty")
	}
	seenIDs := map[string]bool{}
	for i, a := range m.Assumptions {
		if a.ID == "" || a.Statement == "" {
			add("assumptions[%d] requires id and statement", i)
		}
		if !assumptionIDPattern.MatchString(a.ID) {
			add("assumptions[%d] id %q not kebab/snake", i, a.ID)
		}
		if seenIDs[a.ID] {
			add("assumptions duplicate id %q", a.ID)
		}
		seenIDs[a.ID] = true
	}

	for name, list := range map[string][]string{
		"supported": m.CapabilityAudit.Supported, "partial": m.CapabilityAudit.Partial,
		"representational_only": m.CapabilityAudit.RepresentationalOnly,
		"deferred":              m.CapabilityAudit.Deferred, "unavailable": m.CapabilityAudit.Unavailable,
	} {
		if len(list) == 0 {
			add("capability_audit.%s must record at least one entry (PROTOCOL 5)", name)
		}
		for i, e := range list {
			if strings.TrimSpace(e) == "" {
				add("capability_audit.%s[%d] empty", name, i)
			}
		}
	}

	if len(m.Derivations) == 0 {
		add("derivations must be nonempty")
	}
	seenDeriv := map[string]bool{}
	for i, d := range m.Derivations {
		if d.ID == "" || seenDeriv[d.ID] {
			add("derivations[%d] has empty or duplicate id %q", i, d.ID)
		}
		seenDeriv[d.ID] = true
		if !supportStates[d.Support] {
			add("derivations[%d] invalid support %q", i, d.Support)
		}
		if strings.TrimSpace(d.Reason) == "" {
			add("derivations[%d] requires a reason", i)
		}
		for _, r := range d.ExecutedRules {
			if !closedRuleRegistry[r] {
				add("derivations[%d] names rule %q outside the closed registry (PROTOCOL 12.4)", i, r)
			}
		}
		if d.Support == "SUPPORTED" && len(d.ExecutedRules) == 0 {
			add("derivations[%d] SUPPORTED but names no executed rule", i)
		}
		if d.Support != "SUPPORTED" && len(d.ExecutedRules) != 0 {
			add("derivations[%d] %q claims executed rules; only SUPPORTED may", i, d.Support)
		}
	}

	if len(m.InfrastructureUse) == 0 && len(m.Deviations) == 0 {
		add("infrastructure_use is empty and deviations is empty (PROTOCOL 8.1: empty requires justified deviations)")
	}
	// PROTOCOL 8.2.1: every capability audit entry marked [reused] must map to
	// an infrastructure_use entry, and vice versa.
	reusedApis := map[string]bool{}
	for _, iu := range m.InfrastructureUse {
		if iu.APISymbol != "" {
			reusedApis[iu.APISymbol] = true
		}
	}
	for _, s := range m.CapabilityAudit.Supported {
		if !strings.Contains(s, "[reused]") {
			continue
		}
		matched := false
		for api := range reusedApis {
			if strings.Contains(s, api) {
				matched = true
				break
			}
		}
		if !matched {
			add("capability_audit entry %q marked [reused] with no infrastructure_use api_symbol (PROTOCOL 8.2.1)", s)
		}
	}
	for i, iu := range m.InfrastructureUse {
		if iu.Concern == "" || iu.Result == "" || iu.Limitations == "" {
			add("infrastructure_use[%d] requires concern, result, and honest limitations", i)
		}
		if !libraryNames[iu.Library] {
			add("infrastructure_use[%d] library %q not in {phys-artifact, phys-math, phys-lib}", i, iu.Library)
		}
		if !apiSymbolPattern.MatchString(iu.APISymbol) {
			add("infrastructure_use[%d] api_symbol %q not package-qualified", i, iu.APISymbol)
		}
		if !strings.Contains(iu.CodePath, ":") {
			add("infrastructure_use[%d] code_path %q must be File:symbol", i, iu.CodePath)
		}
		if iu.Test == "" {
			add("infrastructure_use[%d] requires test", i)
		} else if testNames != nil && !testNames[iu.Test] {
			add("infrastructure_use[%d] test %q does not exist in this package (PROTOCOL 8.1/11.2)", i, iu.Test)
		}
		matched := false
		for _, s := range m.CapabilityAudit.Supported {
			if strings.Contains(s, "[reused]") && strings.Contains(s, iu.APISymbol) {
				matched = true
				break
			}
		}
		if !matched {
			add("infrastructure_use[%d] api_symbol %q is not claimed in capability_audit.supported as [reused] (PROTOCOL 8.2.1)", i, iu.APISymbol)
		}
	}

	for i, d := range m.Deviations {
		if d.CapabilitySought == "" || d.APIExamined == "" || d.Reason == "" || d.LocalAlternative == "" {
			add("deviations[%d] incomplete: all four PROTOCOL 8.3 fields required", i)
		}
	}

	if len(m.SolverLocalComputations) == 0 {
		add("solver_local_computations must be nonempty for a package with local numerics")
	}
	seenSL := map[string]bool{}
	for i, s := range m.SolverLocalComputations {
		if s.ID == "" || seenSL[s.ID] {
			add("solver_local_computations[%d] empty or duplicate id %q", i, s.ID)
		}
		seenSL[s.ID] = true
		if s.Algorithm == "" || s.Conventions == "" || s.Inputs == "" || s.Outputs == "" || s.IndependentChecks == "" {
			add("solver_local_computations[%d] incomplete (PROTOCOL 9.2 record fields)", i)
		}
		if s.ExecutedBy != executedBySolverLocal {
			add("solver_local_computations[%d] executed_by %q, want %q", i, s.ExecutedBy, executedBySolverLocal)
		}
	}

	if len(m.Results) == 0 {
		add("results must be nonempty")
	}
	seenRes := map[string]bool{}
	for i, r := range m.Results {
		if r.ID == "" || seenRes[r.ID] {
			add("results[%d] empty or duplicate id %q", i, r.ID)
		}
		seenRes[r.ID] = true
		if strings.TrimSpace(r.Value) == "" {
			add("results[%d] requires a value", i)
		}
		if !evidenceLevels[r.EvidenceLevel] {
			add("results[%d] evidence_level %q not in the PROTOCOL 10 vocabulary", i, r.EvidenceLevel)
		}
		if strings.TrimSpace(r.Evidence) == "" {
			add("results[%d] requires evidence", i)
		}
	}

	if len(m.Limitations) == 0 {
		add("limitations must be nonempty")
	}
	for i, l := range m.Limitations {
		if strings.TrimSpace(l) == "" {
			add("limitations[%d] empty", i)
		}
	}
	for i, o := range m.OpenObligations {
		if strings.TrimSpace(o) == "" {
			add("open_obligations[%d] empty", i)
		}
	}
	for i, g := range m.DeclaredGaps {
		if g.GapID == "" || g.Requirement == "" || g.Reason == "" {
			add("declared_gaps[%d] requires gap_id, requirement, reason", i)
		}
		if !gapSeverities[g.Severity] {
			add("declared_gaps[%d] invalid severity %q", i, g.Severity)
		}
	}
	for i, f := range m.FormalizationRefs {
		if strings.TrimSpace(f.Path) == "" {
			add("formalization_refs[%d] requires path", i)
		}
	}

	if m.EntryPoints.Primary == "" || m.EntryPoints.Secondary == "" {
		add("entry_points requires primary and secondary (PROTOCOL 12.1)")
	}
	if m.EntryPoints.Primary != "" && m.EntryPoints.Primary == m.EntryPoints.Secondary {
		add("entry_points primary and secondary must differ")
	}

	// Category firewall (ADR 3 / PROTOCOL 11.2).
	if formalizationClaim(m) && len(m.FormalizationRefs) == 0 {
		add("category firewall: results/derivations claim formalization status without formalization_refs")
	}

	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("manifest validation failed:\n  - %s", strings.Join(errs, "\n  - "))
}

// formalizationClaim reports whether any results[]/derivations[] string claims
// formalization status (Spec 1.4/1.5 vocabulary, ADR 3 firewall).
func formalizationClaim(m manifest) bool {
	for _, r := range m.Results {
		if formalizationPattern.MatchString(r.ID) || formalizationPattern.MatchString(r.Value) ||
			formalizationPattern.MatchString(r.Evidence) {
			return true
		}
	}
	for _, d := range m.Derivations {
		if formalizationPattern.MatchString(d.ID) || formalizationPattern.MatchString(d.Reason) {
			return true
		}
	}
	return false
}

// --- source-level proofs (entry points, adoption reachability, code paths) ---

// parseFiles parses the Go files of a directory; testFiles selects *_test.go.
func parseFiles(t *testing.T, dir string, testFiles bool) map[string]*ast.File {
	t.Helper()
	pattern := "*.go"
	if testFiles {
		pattern = "*_test.go"
	}
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		t.Fatalf("glob %s: %v", dir, err)
	}
	fset := token.NewFileSet()
	out := map[string]*ast.File{}
	for _, path := range matches {
		f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		out[path] = f
	}
	if len(out) == 0 {
		t.Fatalf("no Go files found in %s", dir)
	}
	return out
}

// testFuncNames returns the Test functions defined by this package (root
// directory), for infrastructure_use[].test existence checks.
func testFuncNames(t *testing.T) map[string]bool {
	t.Helper()
	names := map[string]bool{}
	for _, f := range parseFiles(t, ".", true) {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if ok && fn.Recv == nil && strings.HasPrefix(fn.Name.Name, "Test") {
				names[fn.Name.Name] = true
			}
		}
	}
	return names
}

// funcBodies indexes every package-level function reachable by name within
// this module's packages (root tests, axion, representation): name -> bodies.
func funcBodies(t *testing.T) map[string][]*ast.BlockStmt {
	t.Helper()
	idx := map[string][]*ast.BlockStmt{}
	collect := func(files map[string]*ast.File) {
		for _, f := range files {
			for _, decl := range f.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if ok && fn.Recv == nil && fn.Body != nil {
					idx[fn.Name.Name] = append(idx[fn.Name.Name], fn.Body)
				}
			}
		}
	}
	collect(parseFiles(t, ".", true))
	collect(parseFiles(t, "axion", false))
	collect(parseFiles(t, "representation", false))
	return idx
}

// apiReachable reports whether the test function exercises shortSymbol
// directly or through a transitive call into this module's packages
// (PROTOCOL 8.2.2: the claimed API must actually run under that test).
func apiReachable(t *testing.T, bodies map[string][]*ast.BlockStmt, testName, shortSymbol string) bool {
	t.Helper()
	bodiesOf, ok := bodies[testName]
	if !ok {
		return false
	}
	seen := map[string]bool{testName: true}
	queue := append([]*ast.BlockStmt{}, bodiesOf...)
	var fails bool
	walk := func(b *ast.BlockStmt) {
		ast.Inspect(b, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.SelectorExpr:
				if node.Sel.Name == shortSymbol {
					fails = true
					return false
				}
			case *ast.Ident:
				if sub, ok := bodies[node.Name]; ok && !seen[node.Name] {
					seen[node.Name] = true
					queue = append(queue, sub...)
				}
			}
			return true
		})
	}
	for len(queue) > 0 && !fails {
		b := queue[0]
		queue = queue[1:]
		walk(b)
	}
	return fails
}

// shortSymbol extracts the symbol from a package-qualified api_symbol.
func shortSymbol(api string) string {
	parts := strings.Split(api, ".")
	return parts[len(parts)-1]
}

// codePathExists verifies infrastructure_use.code_path: the file exists and
// declares the named symbol (PROTOCOL 8.1).
func codePathExists(t *testing.T, codePath string) error {
	t.Helper()
	idx := strings.LastIndex(codePath, ":")
	if idx < 0 {
		return fmt.Errorf("code_path %q is not File:symbol", codePath)
	}
	file, symbol := codePath[:idx], codePath[idx+1:]
	if _, err := os.Stat(file); err != nil {
		return fmt.Errorf("code_path file %s: %v", file, err)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, parser.SkipObjectResolution)
	if err != nil {
		return fmt.Errorf("parse %s: %v", file, err)
	}
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == symbol {
			return nil
		}
	}
	return fmt.Errorf("code_path symbol %s not declared in %s", symbol, file)
}

// entryPointBlockers statically verifies PROTOCOL 12.1/12.2 against
// protocol/derivation_test.go: each manifest entry point resolves to a
// package function actually called by a test there, and the protocol test is
// not a vacuous stub.
func entryPointBlockers(t *testing.T, m manifest) []string {
	t.Helper()
	var blockers []string
	files := parseFiles(t, "protocol", true)
	assertions := 0
	type selectorCall struct{ x, sel string }
	calls := map[selectorCall]bool{}
	for _, f := range files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				sel, ok := n.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				if x, ok := sel.X.(*ast.Ident); ok {
					calls[selectorCall{x.Name, sel.Sel.Name}] = true
					if x.Name == "t" || x.Name == "tb" {
						switch sel.Sel.Name {
						case "Fatal", "Fatalf", "Error", "Errorf":
							assertions++
						}
					}
				}
				return true
			})
		}
	}
	for _, ep := range []struct{ name, value string }{
		{"primary", m.EntryPoints.Primary},
		{"secondary", m.EntryPoints.Secondary},
	} {
		parts := strings.Split(ep.value, ".")
		if len(parts) != 2 {
			blockers = append(blockers, fmt.Sprintf("entry point %s %q is not Package.Function", ep.name, ep.value))
			continue
		}
		if !calls[selectorCall{parts[0], parts[1]}] {
			blockers = append(blockers, fmt.Sprintf("entry point %s %s is not called by protocol/derivation_test.go (PROTOCOL 12.1)", ep.name, ep.value))
		}
	}
	if assertions < minProtocolAssertions {
		blockers = append(blockers, fmt.Sprintf(
			"protocol/derivation_test.go has %d failure assertions, below the non-vacuity floor %d (PROTOCOL 12.2)",
			assertions, minProtocolAssertions))
	}
	return blockers
}

// conformance is the computed status of PROTOCOL 11.3.
type conformance struct {
	Status  string
	Reasons []string
}

// computeConformance computes CONFORMANT/NONCONFORMANT. CONFORMANT requires:
// schema valid, entry points exercised (extraBlockers), every claim carrying
// evidence at a stated level, zero unresolved mandatory gaps, and no category
// firewall violation. The status is never read from the manifest.
func computeConformance(m manifest, extraBlockers []string) conformance {
	var reasons []string
	if err := validateManifest(m, testFuncNamesNoFail()); err != nil {
		reasons = append(reasons, "schema invalid: "+err.Error())
	}
	reasons = append(reasons, extraBlockers...)
	for _, r := range m.Results {
		if !evidenceLevels[r.EvidenceLevel] || strings.TrimSpace(r.Evidence) == "" {
			reasons = append(reasons, fmt.Sprintf("result %s lacks evidence at a stated level", r.ID))
		}
	}
	for _, g := range m.DeclaredGaps {
		if g.Severity == "mandatory" {
			reasons = append(reasons, fmt.Sprintf("unresolved mandatory gap: %s (%s)", g.GapID, g.Requirement))
		}
	}
	if formalizationClaim(m) && len(m.FormalizationRefs) == 0 {
		reasons = append(reasons, "category firewall violation: formalization claim without formalization_refs")
	}
	if len(reasons) > 0 {
		return conformance{Status: "NONCONFORMANT", Reasons: reasons}
	}
	return conformance{Status: "CONFORMANT"}
}

// testFuncNamesNoFail is the conformance-path variant (no testing.T in
// scope); it re-reads the same sources.
func testFuncNamesNoFail() map[string]bool {
	names := map[string]bool{}
	fset := token.NewFileSet()
	matches, _ := filepath.Glob("*_test.go")
	for _, path := range matches {
		f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			continue
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if ok && fn.Recv == nil && strings.HasPrefix(fn.Name.Name, "Test") {
				names[fn.Name.Name] = true
			}
		}
	}
	return names
}

// --- the tests ---

// TestManifestSchemaConformance validates the real manifest against every
// schema rule and computes the package status (PROTOCOL 11.3).
func TestManifestSchemaConformance(t *testing.T) {
	raw := loadManifestRaw(t)
	m, top, err := decodeManifestStrict(raw)
	if err != nil {
		t.Fatalf("strict decode: %v", err)
	}
	for _, e := range checkTopLevelKeys(top) {
		t.Errorf("top-level key rule: %s", e)
	}
	names := testFuncNames(t)
	if err := validateManifest(m, names); err != nil {
		t.Fatalf("validate: %v", err)
	}
	// Adoption proofs: code_path files/symbols exist; each recorded test
	// reaches the claimed API through this module's call graph.
	bodies := funcBodies(t)
	for _, iu := range m.InfrastructureUse {
		if err := codePathExists(t, iu.CodePath); err != nil {
			t.Errorf("infrastructure_use %q: %v", iu.Concern, err)
		}
		if !apiReachable(t, bodies, iu.Test, shortSymbol(iu.APISymbol)) {
			t.Errorf("infrastructure_use %q: test %s does not reach %s (PROTOCOL 8.2.2)",
				iu.Concern, iu.Test, iu.APISymbol)
		}
	}
	// Entry-point exercise and protocol-test non-vacuity.
	blockers := entryPointBlockers(t, m)
	for _, b := range blockers {
		t.Errorf("entry-point check: %s", b)
	}
	c := computeConformance(m, blockers)
	if c.Status != "CONFORMANT" {
		t.Fatalf("conformance computed NONCONFORMANT:\n  - %s", strings.Join(c.Reasons, "\n  - "))
	}
	t.Logf("conformance: CONFORMANT (schema valid, %d entry points exercised, %d infrastructure_use entries adopted, %d deviations, %d solver-local computations, %d advisory gaps, 0 mandatory gaps)",
		2, len(m.InfrastructureUse), len(m.Deviations), len(m.SolverLocalComputations), len(m.DeclaredGaps))
}

// --- deliberate-corruption negatives (PROTOCOL 12.7) ---

func mutateManifest(t *testing.T, f func(*manifest)) manifest {
	t.Helper()
	m := loadManifest(t)
	f(&m)
	return m
}

func mustReject(t *testing.T, m manifest, wantSubstring string) {
	t.Helper()
	err := validateManifest(m, testFuncNames(t))
	if err == nil {
		t.Fatalf("corrupted manifest accepted; want rejection containing %q", wantSubstring)
	}
	if !strings.Contains(err.Error(), wantSubstring) {
		t.Fatalf("rejection %q does not mention %q", err.Error(), wantSubstring)
	}
}

func TestManifestRejectsWrongSchemaID(t *testing.T) {
	mustReject(t, mutateManifest(t, func(m *manifest) {
		m.SchemaID = "phys-demo/benchmark-solution-manifest/2"
	}), "schema_id")
}

func TestManifestRejectsUnknownTopLevelKey(t *testing.T) {
	raw := loadManifestRaw(t)
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	top["sidecar_notes"] = json.RawMessage(`{"x":1}`)
	mutated, err := json.Marshal(top)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if _, _, err := decodeManifestStrict(mutated); err == nil {
		t.Fatal("unknown top-level key accepted by strict decode")
	}
}

func TestManifestRejectsMissingRequiredTopLevelKey(t *testing.T) {
	raw := loadManifestRaw(t)
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	delete(top, "regime")
	mutated, err := json.Marshal(top)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	_, parsed, err := decodeManifestStrict(mutated)
	if err != nil {
		t.Fatalf("strict decode of otherwise valid manifest: %v", err)
	}
	found := false
	for _, e := range checkTopLevelKeys(parsed) {
		if strings.Contains(e, `"regime"`) {
			found = true
		}
	}
	if !found {
		t.Fatal("missing required key regime not detected")
	}
}

func TestManifestRejectsMissingUnit(t *testing.T) {
	mustReject(t, mutateManifest(t, func(m *manifest) {
		m.Parameters[0].Unit = ""
	}), "explicit requires unit")
}

func TestManifestRejectsDimensionlessWithUnit(t *testing.T) {
	mustReject(t, mutateManifest(t, func(m *manifest) {
		m.Parameters[2].Unit = "eV"
	}), "dimensionless must not carry")
}

func TestManifestRejectsBrokenInfraTestRef(t *testing.T) {
	mustReject(t, mutateManifest(t, func(m *manifest) {
		m.InfrastructureUse[0].Test = "TestDoesNotExist"
	}), "does not exist")
}

func TestManifestRejectsForgottenReuseMarker(t *testing.T) {
	// PROTOCOL 8.2.1: an audit capability marked [reused] with no matching
	// infrastructure_use entry must fail.
	mustReject(t, mutateManifest(t, func(m *manifest) {
		m.CapabilityAudit.Supported = append(m.CapabilityAudit.Supported,
			"physmath.SolveSymbolically - magic [reused]")
	}), "8.2.1")
}

func TestManifestRejectsInventedRuleID(t *testing.T) {
	mustReject(t, mutateManifest(t, func(m *manifest) {
		m.Derivations[2].ExecutedRules = []string{"DELTA-SUBSTITUTION-ANALYTIC"}
	}), "closed registry")
}

func TestManifestRejectsBadEvidenceLevel(t *testing.T) {
	mustReject(t, mutateManifest(t, func(m *manifest) {
		m.Results[0].EvidenceLevel = "Verified"
	}), "evidence_level")
}

func TestManifestRejectsBadSupportState(t *testing.T) {
	mustReject(t, mutateManifest(t, func(m *manifest) {
		m.Derivations[0].Support = "EXECUTABLE"
	}), "invalid support")
}

func TestManifestRejectsBadExecutedBy(t *testing.T) {
	mustReject(t, mutateManifest(t, func(m *manifest) {
		m.SolverLocalComputations[0].ExecutedBy = "physmath"
	}), "executed_by")
}

func TestManifestRejectsIncompleteDeviation(t *testing.T) {
	mustReject(t, mutateManifest(t, func(m *manifest) {
		m.Deviations[0].Reason = ""
	}), "deviations[0] incomplete")
}

func TestManifestRejectsEmptyInfraUseWithoutDeviations(t *testing.T) {
	mustReject(t, mutateManifest(t, func(m *manifest) {
		m.InfrastructureUse = nil
		m.Deviations = nil
	}), "PROTOCOL 8.1")
}

func TestManifestRejectsFormalizationClaimWithoutRefs(t *testing.T) {
	m := mutateManifest(t, func(m *manifest) {
		m.Results[0].Evidence += " This is a formally represented theorem."
	})
	mustReject(t, m, "category firewall")
	// The computed status must also flip to NONCONFORMANT (ADR 3 evasion rule).
	c := computeConformance(m, nil)
	if c.Status != "NONCONFORMANT" {
		t.Fatalf("firewall violation computed %s, want NONCONFORMANT", c.Status)
	}
}

func TestConformanceFlipsOnMandatoryGap(t *testing.T) {
	m := mutateManifest(t, func(m *manifest) {
		m.DeclaredGaps[0].Severity = "mandatory"
	})
	if err := validateManifest(m, testFuncNames(t)); err != nil {
		t.Fatalf("manifest with mandatory gap should pass schema: %v", err)
	}
	c := computeConformance(m, nil)
	if c.Status != "NONCONFORMANT" {
		t.Fatalf("mandatory gap computed %s, want NONCONFORMANT (PROTOCOL 11.3/11.4)", c.Status)
	}
	joined := strings.Join(c.Reasons, " ")
	if !strings.Contains(joined, "GAP-NO-DISTRIBUTION-RULE") {
		t.Fatalf("mandatory-gap reason missing gap id: %s", joined)
	}
}

func TestConformanceFlipsOnUnexercisedEntryPoint(t *testing.T) {
	m := mutateManifest(t, func(m *manifest) {
		m.EntryPoints.Primary = "axion.Nope"
	})
	blockers := entryPointBlockers(t, m)
	if len(blockers) == 0 {
		t.Fatal("unresolvable entry point produced no blocker")
	}
	c := computeConformance(m, blockers)
	if c.Status != "NONCONFORMANT" {
		t.Fatalf("unexercised entry point computed %s, want NONCONFORMANT", c.Status)
	}
}
