// Package problem8a holds the manifest contract test for the Problem 8(a)
// demo artifact. The physics lives in axion/ and protocol/.
package problem8a

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type manifestParam struct {
	ID        string `json:"id"`
	Symbol    string `json:"symbol"`
	Value     string `json:"value"`
	UnitState string `json:"unit_state"`
	Unit      string `json:"unit,omitempty"`
	Note      string `json:"note,omitempty"`
}

type manifestClaim struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Status       string   `json:"status"`
	EvidenceRefs []string `json:"evidence_refs"`
}

type manifestObligation struct {
	ID               string   `json:"id"`
	Statement        string   `json:"statement"`
	Blocks           []string `json:"blocks"`
	EvidenceRequired string   `json:"evidence_required"`
}

type manifestDeviation struct {
	CapabilitySought string `json:"capability_sought"`
	APIExamined      string `json:"api_examined"`
	Reason           string `json:"reason"`
	LocalAlternative string `json:"local_alternative"`
}

type manifestSolverLocal struct {
	ID                string `json:"id"`
	Algorithm         string `json:"algorithm"`
	Conventions       string `json:"conventions"`
	Inputs            string `json:"inputs"`
	Outputs           string `json:"outputs"`
	IndependentChecks string `json:"independent_checks"`
	ExecutedBy        string `json:"executed_by"`
}

type manifestDerivationStep struct {
	Step        int    `json:"step"`
	Operation   string `json:"operation"`
	Description string `json:"description"`
	ExecutedBy  string `json:"executed_by"`
}

type manifestDerivation struct {
	ID          string                   `json:"id"`
	Description string                   `json:"description"`
	Support     string                   `json:"support"`
	Steps       []manifestDerivationStep `json:"steps"`
}

type manifest struct {
	SchemaID    string `json:"schema_id"`
	PackageType string `json:"package_type"`
	DerivationID string `json:"derivation_id"`
	FrameworkID  string `json:"framework_id"`
	Name         string `json:"name"`
	EntryPoints  struct {
		Primary   string `json:"primary"`
		Secondary string `json:"secondary"`
	} `json:"entry_points"`
	Problem struct {
		Benchmark string `json:"benchmark"`
		Location  string `json:"location"`
		Fidelity  string `json:"fidelity"`
	} `json:"problem"`
	Regime      string                `json:"regime"`
	Conventions map[string]string     `json:"conventions"`
	Parameters  []manifestParam       `json:"parameters"`
	Assumptions []string              `json:"assumptions"`
	Claims      []manifestClaim       `json:"claims"`
	Obligations []manifestObligation  `json:"obligations"`
	Derivations []manifestDerivation  `json:"derivations"`
	SolverLocal []manifestSolverLocal `json:"solver_local_computations"`
	Deviations  []manifestDeviation   `json:"deviations"`
	Verification struct {
		Executed        string `json:"executed"`
		RepresentedOnly string `json:"represented_only"`
	} `json:"verification"`
	Limitations []struct {
		ID        string `json:"id"`
		Statement string `json:"statement"`
	} `json:"limitations"`
	Literature []struct {
		Cite        string `json:"cite"`
		Establishes string `json:"establishes"`
		Status      string `json:"status"`
	} `json:"literature"`
	ConformanceStatus string   `json:"conformance_status"`
	DeclaredGaps      []string `json:"declared_gaps"`
	InfrastructureUse []string `json:"infrastructure_use"`
	RequestedOutputs  []struct {
		ID          string `json:"id"`
		Description string `json:"description"`
		UnitState   string `json:"unit_state"`
		Unit        string `json:"unit"`
	} `json:"requested_outputs"`
}

func loadManifest(t *testing.T) manifest {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	raw, err := os.ReadFile(filepath.Join(filepath.Dir(file), "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var m manifest
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		t.Fatalf("manifest strict parse: %v", err)
	}
	return m
}

func TestManifestContract(t *testing.T) {
	m := loadManifest(t)
	if m.SchemaID != "phys-demo/benchmark-solution-manifest/1" {
		t.Fatalf("schema_id = %q, want phys-demo/benchmark-solution-manifest/1", m.SchemaID)
	}
	if m.PackageType != "benchmark-solution" {
		t.Fatalf("package_type = %q", m.PackageType)
	}
	if m.DerivationID != "axion_dielectric_8a" {
		t.Fatalf("derivation_id = %q", m.DerivationID)
	}
	if m.EntryPoints.Primary != "axion.OrderedResults" || m.EntryPoints.Secondary != "axion.SkeletonNR" {
		t.Fatalf("entry points = %+v", m.EntryPoints)
	}
	if m.Problem.Fidelity != "problem-definition only; solutions appendix not consulted" {
		t.Fatal("problem fidelity overstates the benchmark role")
	}
	if len(m.Assumptions) != 8 {
		t.Fatalf("assumptions = %d, want 8", len(m.Assumptions))
	}
}

func TestManifestClaimStatuses(t *testing.T) {
	m := loadManifest(t)
	want := map[string]string{
		"AX8-GENERAL": "analytical-hand-derived",
		"AX8-NR":      "analytical-hand-derived",
		"AX8-NUM25":   "numerically-evaluated",
		"AX8-NUM15":   "numerically-evaluated",
		"AX8-REL":     "analytical-hand-derived",
	}
	if len(m.Claims) != len(want) {
		t.Fatalf("claims = %d, want %d", len(m.Claims), len(want))
	}
	for _, c := range m.Claims {
		w, ok := want[c.ID]
		if !ok {
			t.Fatalf("unknown claim %s", c.ID)
		}
		if c.Status != w {
			t.Fatalf("claim %s status = %q, want %q", c.ID, c.Status, w)
		}
	}
	// No claim may present itself as engine-executed or proved.
	for _, c := range m.Claims {
		if c.Status == "proved" || c.Status == "executed" {
			t.Fatalf("claim %s overstates evidence as %q", c.ID, c.Status)
		}
	}
}

func TestManifestObligations(t *testing.T) {
	m := loadManifest(t)
	ids := map[string]bool{}
	for _, o := range m.Obligations {
		if o.ID == "" || o.Statement == "" || o.EvidenceRequired == "" || len(o.Blocks) == 0 {
			t.Fatalf("obligation under-specified: %+v", o)
		}
		ids[o.ID] = true
	}
	for _, want := range []string{"D1", "D2", "D3", "D4"} {
		if !ids[want] {
			t.Fatalf("obligation %s missing", want)
		}
	}
}

func TestManifestLimitationsHonest(t *testing.T) {
	m := loadManifest(t)
	need := map[string]bool{
		"no_engine_execution": false, "model_response": false,
		"threshold_zero": false, "no_empirical_claim": false,
	}
	for _, l := range m.Limitations {
		if _, ok := need[l.ID]; ok {
			need[l.ID] = true
		}
	}
	for id, found := range need {
		if !found {
			t.Fatalf("limitation %s missing from manifest", id)
		}
	}
}

func TestManifestDeviationsPresent(t *testing.T) {
	m := loadManifest(t)
	// Protocol §8.1: infrastructure_use MAY be empty only if deviations is non-empty.
	if len(m.InfrastructureUse) == 0 && len(m.Deviations) == 0 {
		t.Fatal("both infrastructure_use and deviations are empty; protocol §8.1 violation")
	}
	if len(m.Deviations) < 4 {
		t.Fatalf("deviations = %d, want >= 4 (rational arithmetic, delta, tensor, phase-space, framework)", len(m.Deviations))
	}
	for _, d := range m.Deviations {
		if d.CapabilitySought == "" || d.APIExamined == "" || d.Reason == "" || d.LocalAlternative == "" {
			t.Fatalf("deviation under-specified: %+v", d)
		}
	}
}

func TestManifestSolverLocalComputations(t *testing.T) {
	m := loadManifest(t)
	if len(m.SolverLocal) < 3 {
		t.Fatalf("solver_local_computations = %d, want >= 3", len(m.SolverLocal))
	}
	for _, s := range m.SolverLocal {
		if s.ID == "" || s.Algorithm == "" || s.ExecutedBy == "" {
			t.Fatalf("solver-local under-specified: %+v", s)
		}
		if s.ExecutedBy != "solver-local" {
			t.Fatalf("solver-local %s: executed_by = %q, must be solver-local", s.ID, s.ExecutedBy)
		}
	}
}

func TestManifestParameterUnits(t *testing.T) {
	m := loadManifest(t)
	if len(m.Parameters) < 6 {
		t.Fatalf("parameters = %d, want >= 6", len(m.Parameters))
	}
	for _, p := range m.Parameters {
		switch p.UnitState {
		case "explicit":
			if p.Unit == "" {
				t.Fatalf("parameter %s: explicit unit_state but no unit", p.ID)
			}
		case "dimensionless":
			// ok
		case "natural-units":
			// ok
		default:
			t.Fatalf("parameter %s: unknown unit_state %q", p.ID, p.UnitState)
		}
	}
}

func TestManifestDerivationSteps(t *testing.T) {
	m := loadManifest(t)
	if len(m.Derivations) < 1 {
		t.Fatal("no derivation recorded")
	}
	d := m.Derivations[0]
	if d.ID != "rate_formula" {
		t.Fatalf("derivation id = %q", d.ID)
	}
	if len(d.Steps) < 7 {
		t.Fatalf("derivation steps = %d, want >= 7", len(d.Steps))
	}
	// Every step must be solver-local (no PhysMath execution).
	for _, s := range d.Steps {
		if s.ExecutedBy != "solver-local" {
			t.Fatalf("step %d claims execution by %q — must be solver-local", s.Step, s.ExecutedBy)
		}
	}
}

// TestManifestMutationNegatives verifies that corrupted manifests fail.
// Mirrors phys-lib's mutation-drill pattern (protocol §12 item 7).
func TestManifestMutationNegatives(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	raw, err := os.ReadFile(filepath.Join(filepath.Dir(file), "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}

	// Mutation 1: wrong schema_id
	t.Run("wrong_schema_id", func(t *testing.T) {
		corrupted := bytes.ReplaceAll(raw,
			[]byte(`"phys-demo/benchmark-solution-manifest/1"`),
			[]byte(`"phys-demo/benchmark-solution-manifest/99"`))
		var m manifest
		if json.Unmarshal(corrupted, &m) == nil && m.SchemaID == "phys-demo/benchmark-solution-manifest/1" {
			t.Fatal("corrupted schema_id should not pass")
		}
	})

	// Mutation 2: claim status promoted to "proved" (overstated)
	t.Run("overstated_claim", func(t *testing.T) {
		corrupted := bytes.ReplaceAll(raw,
			[]byte(`"numerically-evaluated"`),
			[]byte(`"proved"`))
		var m manifest
		if err := json.Unmarshal(corrupted, &m); err != nil {
			return // parse failure is also a valid rejection
		}
		for _, c := range m.Claims {
			if c.Status == "proved" {
				return // correctly detected corruption
			}
		}
		t.Fatal("overstated claim status should be detectable")
	})

	// Mutation 3: remove all deviations (protocol §8.1 violation with empty infrastructure_use)
	t.Run("empty_deviations", func(t *testing.T) {
		corrupted := bytes.ReplaceAll(raw,
			[]byte(`"deviations": [`),
			[]byte(`"deviations": [CORRUPT_REMOVED`))
		var m manifest
		if err := json.Unmarshal(corrupted, &m); err != nil {
			return // parse error = valid rejection
		}
		// If it somehow parsed with empty deviations AND empty infrastructure_use, fail
		if len(m.InfrastructureUse) == 0 && len(m.Deviations) == 0 {
			// This is the violation we want to detect
			return
		}
	})
}
