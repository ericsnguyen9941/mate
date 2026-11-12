package notice

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// The hand-back question sends the text under its own question name with its
// own three labels, and reads the answer back the way the notice question
// does (CON-274).
func TestClassifyHandbackAsksItsOwnQuestion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			State     string
			Questions map[string]struct {
				Type     string
				Criteria map[string]string
			}
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
			return
		}
		q, ok := req.Questions["handback"]
		if !ok || q.Type != "choice" || len(q.Criteria) != 3 || q.Criteria[LabelRoutine] == "" || strings.Contains(req.State, "sk-secret") {
			t.Errorf("bad request: %+v", req)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"model": Model, "answers": map[string]any{"handback": map[string]any{
			"type": "choice", "choice": LabelRoutine, "confidence": .91,
			"probabilities": map[string]float64{LabelRoutine: .91, LabelNeedsJudgment: .08, LabelUnknown: .01},
		}}, "usage": map[string]int{"input_tokens": 300, "output_tokens": 5}})
	}))
	defer srv.Close()
	c := New("test-key")
	c.endpoint = srv.URL
	got, err := c.ClassifyHandback(context.Background(), "## Deviations from Build\nUpdated docs. token=sk-secret\n## Still open\nnone")
	if err != nil || got.Label != LabelRoutine || got.Confidence != .91 {
		t.Fatalf("got %+v, %v", got, err)
	}
}

// A document keeps its head, which is where a hand-back says what matters.
func TestPrepareTextKeepsTheHead(t *testing.T) {
	text := "## Deviations\nfirst\n" + strings.Repeat("é", MaxTextBytes)
	got := PrepareText(text)
	if !strings.HasPrefix(got, "## Deviations\nfirst") || len(got) > MaxTextBytes {
		t.Fatalf("PrepareText kept %d bytes starting %q", len(got), got[:20])
	}
}
