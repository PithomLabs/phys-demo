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

type manifest struct {
	SchemaVersion string `json:"schema_version"`
	ArtifactType  string `json:"artifact_type"`
	DerivationID  string `json:"derivation_id"`
	FrameworkID   string `json:"framework_id"`
	Name          string `json:"name"`
	EntryPoints   struct {
		Primary   string `json:"primary"`
		Secondary string `json:"secondary"`
	} `json:"entry_points"`
	ProblemSource struct {
		Benchmark string `json:"benchmark"`
		Location  string `json:"location"`
		Fidelity  string `json:"fidelity"`
	} `json:"problem_source"`
	Regime       string               `json:"regime"`
	Conventions  map[string]string    `json:"conventions"`
	Claims       []manifestClaim      `json:"claims"`
	Obligations  []manifestObligation `json:"obligations"`
	Assumptions  []string             `json:"assumptions"`
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
	if m.SchemaVersion != "1" || m.ArtifactType != "derivation" {
		t.Fatalf("manifest header = %q %q", m.SchemaVersion, m.ArtifactType)
	}
	if m.DerivationID != "axion_dielectric_8a" {
		t.Fatalf("derivation_id = %q", m.DerivationID)
	}
	if m.EntryPoints.Primary != "axion.OrderedResults" || m.EntryPoints.Secondary != "axion.SkeletonNR" {
		t.Fatalf("entry points = %+v", m.EntryPoints)
	}
	if m.ProblemSource.Fidelity != "problem-definition only; solutions appendix not consulted" {
		t.Fatal("problem-source fidelity overstates the benchmark role")
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
