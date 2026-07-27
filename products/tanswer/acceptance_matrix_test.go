package tanswer

import (
	"os"
	"strings"
	"testing"
)

func TestAcceptanceMatrixDeclaresRequiredControls(t *testing.T) {
	raw, err := os.ReadFile("ACCEPTANCE_MATRIX.md")
	if err != nil {
		t.Fatalf("read final acceptance matrix: %v", err)
	}
	matrix := string(raw)
	for _, requirement := range []string{
		"R-01 Configuration contract",
		"R-02 Runtime discovery",
		"R-03 Command inventory",
		"R-04 Semantic protected writes",
		"R-05 Raw API safety",
		"R-06 Human and AI task journey",
		"R-07 Documentation and secret hygiene",
		"R-08 Negative controls",
	} {
		if !strings.Contains(matrix, requirement) {
			t.Errorf("acceptance matrix missing %q", requirement)
		}
	}
}
