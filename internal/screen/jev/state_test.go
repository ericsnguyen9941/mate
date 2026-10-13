package jev

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/notice"
	"github.com/nguyenngocanh94/mate/internal/screen"
)

// The state question's labels are exactly mate's states: Jev classifies, it
// never adds a state.
func TestStateQuestionLabelsAreMatesStates(t *testing.T) {
	body, err := StateRequestBody("› 1. Review hooks")
	if err != nil {
		t.Fatal(err)
	}
	var req struct {
		Questions map[string]struct {
			Criteria map[string]string `json:"criteria"`
		} `json:"questions"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatal(err)
	}
	if len(req.Questions) != 1 {
		t.Fatalf("questions = %v, want the state question alone", req.Questions)
	}
	crit := req.Questions[string(AxisState)].Criteria
	if len(crit) != len(screen.PaneStates) {
		t.Fatalf("%d labels, want %d", len(crit), len(screen.PaneStates))
	}
	for _, s := range screen.PaneStates {
		if crit[string(s)] == "" {
			t.Errorf("state %s has no criterion", s)
		}
	}
}

// The measured three-question request is untouched.
func TestStateQuestionLeavesTheMeasuredRequestAlone(t *testing.T) {
	if StatePromptHash() == PromptHash() {
		t.Fatal("the state question must have its own prompt hash")
	}
	body, _ := RequestBody("")
	if strings.Contains(string(body), string(AxisState)) {
		t.Fatal("the measured request now asks the state question")
	}
}

func stateAnswer(choice string, conf float64) []byte {
	probs := map[string]float64{}
	rest := (1 - conf) / float64(len(screen.PaneStates)-1)
	for _, s := range screen.PaneStates {
		probs[string(s)] = rest
	}
	probs[choice] = conf
	b, _ := json.Marshal(map[string]any{"model": notice.Model, "answers": map[string]any{
		string(AxisState): map[string]any{"type": "choice", "choice": choice, "confidence": conf, "probabilities": probs}}})
	return b
}

func TestObserveStateReadsJevsChoice(t *testing.T) {
	client := notice.New("test-key").WithTransport(roundTrip(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}},
			Body: io.NopCloser(strings.NewReader(string(stateAnswer("dialog_hooks_review", 0.93))))}, nil
	}))
	r, err := New(client).(*Observer).ObserveState(context.Background(), nil, "Hooks need review\n› 1. Review hooks\n")
	if err != nil || r.State != screen.PaneHooksReview || r.Confidence != 0.93 || r.Source != Source {
		t.Fatalf("reading = %+v %v", r, err)
	}
}

func TestParseStateRefusesALabelOutsideMatesStates(t *testing.T) {
	bad := strings.Replace(string(stateAnswer("dialog_other", 0.9)), `"choice":"dialog_other"`, `"choice":"press_enter"`, 1)
	if _, err := ParseState([]byte(bad)); err == nil {
		t.Fatal("a choice outside the states was accepted")
	}
}
