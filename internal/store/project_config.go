package store

// MateConfig pins the coordinator profile independently of the captain's CLI.
type MateConfig struct {
	Model  string `yaml:"model,omitempty"`
	Effort string `yaml:"effort,omitempty"`
	// RefreshContext defaults to 150000; -1 disables automatic refresh.
	RefreshContext int64 `yaml:"refresh_context,omitempty"`
	// Mode is the Mate's mode when a start names none: `headless` (CON-78,
	// the default for a harness with headless turns) or `pane`, the
	// harness's own TUI. It wins over the mode the last run recorded
	// (CON-220) and takes effect at the next start or restart.
	Mode string `yaml:"mode,omitempty"`
	// Headless is the first spelling of Mode: headless.
	Headless bool `yaml:"headless,omitempty"`
}

// ProjectConfig is `projects/<name>/project.yaml`.
type ProjectConfig struct {
	Mate MateConfig `yaml:"mate,omitempty"`
	// Repos are the git repositories the project owns, zero or more
	// (docs/mvp.md M9). A crew works in exactly one of them.
	Repos []RepoConfig `yaml:"repos"`
	// Budget is the optional token/cost ceiling mvp.md M5 task 27 checks at
	// the end of every observer poll. A nil Budget (the field absent from
	// project.yaml) means no limit is configured, not a limit of zero:
	// BudgetConfig's own fields carry that same "unset means no limit"
	// convention one level down.
	Budget *BudgetConfig `yaml:"budget,omitempty"`
	// Review is whether each task's work is independently reviewed before
	// the Mate lands it (CON-101): ReviewOn or ReviewOff. Empty means the
	// captain has not decided, and the Mate asks before it lands one.
	Review string `yaml:"review,omitempty"`
	// NextWork is what the Mate does when a task lands and nothing the
	// captain asked for is left (CON-120): NextWorkPick, the default when
	// empty, takes the tracker's next ready task without asking;
	// NextWorkAsk asks the captain which work comes next.
	NextWork string `yaml:"next_work,omitempty"`
	// Workflow is the project's own process as the admin task page checks
	// it (docs/design/task-page.md, CON-148): the steps the harness asks a
	// crew to take, each with the evidence that proves it was taken. Empty
	// means the project declares none and the page draws no checklist.
	Workflow []WorkflowStep `yaml:"workflow,omitempty"`
}

// WorkflowStep is one step of ProjectConfig.Workflow. Every pattern is a
// Go regular expression; a step is done when any one of them matches:
// Command against a shell command the crew ran, Read against a path it
// read (a file tool's target, or a command naming the path), Skill against
// a skill it loaded, Handback against a line of its hand-back. Files, when
// set, says which changed files make the step apply at all: a crew whose
// diff matches none of them has the step as "n/a", not "missed". Since,
// when set, is the day the repo started asking for the step (2026-10-28,
// or an RFC 3339 time): a task spawned before it has the step as "n/a".
// Roles names the crews the step is asked of (WorkflowRoleShip,
// WorkflowRoleScout, WorkflowRoleReview); empty means ship crews only,
// because a workflow is the way a change is shipped. A crew of another
// role has the step as "n/a" (CON-201).
type WorkflowStep struct {
	Step     string   `yaml:"step"`
	Command  string   `yaml:"command,omitempty"`
	Read     string   `yaml:"read,omitempty"`
	Skill    string   `yaml:"skill,omitempty"`
	Handback string   `yaml:"handback,omitempty"`
	Files    string   `yaml:"files,omitempty"`
	Since    string   `yaml:"since,omitempty"`
	Roles    []string `yaml:"roles,omitempty"`
}

// The values of WorkflowStep.Roles, the dashboard's task roles.
const (
	WorkflowRoleShip   = "ship"
	WorkflowRoleScout  = "scout"
	WorkflowRoleReview = "review"
)

// AskedOf reports whether the step is asked of a crew of role.
func (s WorkflowStep) AskedOf(role string) bool {
	if len(s.Roles) == 0 {
		return role == WorkflowRoleShip
	}
	for _, r := range s.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// The values of ProjectConfig.Review.
const (
	ReviewOn  = "on"
	ReviewOff = "off"
)

// ValidReview reports whether s is a value project.yaml's review may hold.
func ValidReview(s string) bool { return s == ReviewOn || s == ReviewOff }

// The values of ProjectConfig.NextWork.
const (
	NextWorkPick = "pick"
	NextWorkAsk  = "ask"
)

// ValidNextWork reports whether s is a value project.yaml's next_work may hold.
func ValidNextWork(s string) bool { return s == NextWorkPick || s == NextWorkAsk }

// NextWorkMode is NextWork with its default filled in: an empty or unknown
// value picks, so a project the captain never configured keeps moving.
func (c ProjectConfig) NextWorkMode() string {
	if c.NextWork == NextWorkAsk {
		return NextWorkAsk
	}
	return NextWorkPick
}

// BudgetConfig is `project.yaml`'s `budget:` block. Every field is optional
// and a zero value means "no limit on this dimension" - hand-editing the
// file to add one line is the whole interface for the MVP; `project add
// --budget-usd` is the one flag that also writes it, for the caller who
// wants a project-wide ceiling at creation time.
//
// Crossing a limit opens an incident (`budget`, mvp.md section 4b) on the
// crew (for the two crew-scoped fields) or on the project's Mate (for
// ProjectUSD); per the 2026-09-20 decision it never makes a crew `blocked`
// and it is never resolved, because going over budget is a fact about the
// task's past spend, not a condition that clears.
type BudgetConfig struct {
	// CrewTokens is the total token ceiling (all four buckets summed) for
	// one crew's whole task.
	CrewTokens int64 `yaml:"crew_tokens,omitempty"`
	// CrewUSD is the cost ceiling for one crew's whole task, using
	// `pricing.yaml`. A crew whose model has no price never crosses this,
	// because its cost is unknown rather than zero.
	CrewUSD float64 `yaml:"crew_usd,omitempty"`
	// ProjectUSD is the cost ceiling across every crew this project has
	// ever spawned, checked against the Mate rather than any one crew.
	ProjectUSD float64 `yaml:"project_usd,omitempty"`
}
