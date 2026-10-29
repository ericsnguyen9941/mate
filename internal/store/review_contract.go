package store

import (
	"fmt"
	"strconv"
	"strings"
)

// A record without a contract keeps the historical report rules. New reviews
// opt into explicit verification gaps; old reports are never silently upgraded.
const (
	MetaReviewContract      = "review_contract"
	ReviewContract          = "2"
	VerificationGapsSection = "Verification gaps"
)

// ValidateReviewPass checks explicit structure, not the truth of the evidence.
// The reviewer owns semantic judgment; a contradictory pass cannot authorize
// a merge merely because its first Verdict line says pass.
func ValidateReviewPass(report string) error {
	sections := reviewSections(report)
	for _, name := range []string{VerdictSection, FindingsSection, VerificationGapsSection} {
		bodies := sections[strings.ToLower(name)]
		if len(bodies) != 1 {
			return fmt.Errorf("expected exactly one ## %s section", name)
		}
		body := strings.TrimSpace(bodies[0])
		if name == VerdictSection {
			if strings.SplitN(body, "\n", 2)[0] != "pass" {
				return fmt.Errorf("## Verdict is not pass")
			}
		} else if body != "none" {
			return fmt.Errorf("a pass requires ## %s to say none", name)
		}
	}
	return nil
}

// ValidatePriorFindingClosures prevents duplicate, unknown or malformed marks
// from disappearing into ReviewPriorFindings' legacy map representation.
func ValidatePriorFindingClosures(report string, count int) error {
	sections := reviewSections(report)[strings.ToLower(PriorFindingsSection)]
	if len(sections) != 1 {
		return fmt.Errorf("expected exactly one ## %s section", PriorFindingsSection)
	}
	if count == 0 {
		if strings.TrimSpace(sections[0]) != "none" {
			return fmt.Errorf("no prior findings: ## %s must say none", PriorFindingsSection)
		}
		return nil
	}
	seen := map[int]bool{}
	for _, item := range reportItems(report, PriorFindingsSection) {
		m := priorFindingPattern.FindStringSubmatch(item)
		if m == nil {
			return fmt.Errorf("invalid prior finding mark %q", item)
		}
		n, err := strconv.Atoi(m[1])
		if err != nil || n < 1 || n > count || seen[n] {
			return fmt.Errorf("duplicate or unexpected prior finding F%s", m[1])
		}
		if !strings.EqualFold(m[2], FindingClosed) {
			return fmt.Errorf("prior finding F%d is still open", n)
		}
		seen[n] = true
	}
	if len(seen) != count {
		return fmt.Errorf("expected %d closed prior findings, got %d", count, len(seen))
	}
	return nil
}

// Keep every occurrence, including empty sections, so a second heading cannot
// conceal contradictory evidence. Fenced examples are not report structure.
func reviewSections(report string) map[string][]string {
	out := map[string][]string{}
	name, fence := "", ""
	var body strings.Builder
	flush := func() {
		if name != "" {
			out[name] = append(out[name], body.String())
		}
		body.Reset()
	}
	for _, line := range strings.Split(report, "\n") {
		trim := strings.TrimSpace(line)
		if fence != "" {
			if strings.HasPrefix(trim, fence) && strings.Trim(trim, string(fence[0])) == "" {
				fence = ""
			}
		} else if strings.HasPrefix(trim, "```") || strings.HasPrefix(trim, "~~~") {
			n := 0
			for n < len(trim) && trim[n] == trim[0] {
				n++
			}
			fence = trim[:n]
		} else if strings.HasPrefix(trim, "## ") || strings.HasPrefix(trim, "# ") {
			flush()
			name = ""
			if strings.HasPrefix(trim, "## ") {
				name = strings.ToLower(strings.TrimSpace(trim[3:]))
			}
			continue
		}
		body.WriteString(line + "\n")
	}
	flush()
	return out
}
