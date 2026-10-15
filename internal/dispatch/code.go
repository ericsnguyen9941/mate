package dispatch

import (
	"path"
	"strings"
)

// codeExtensions is what makes a changed file "code" for the small-change
// review (CON-74): a conservative list of source-file extensions, kept
// short on purpose. The question it answers is "could a few changed lines
// here change what the program does", so it names the languages a crew
// writes and the build and shell files that run them, and leaves out prose
// (.md, .txt), data and configuration (.json, .yaml, .toml, .csv), styling
// and markup (.css, .html, .svg): a 2-line edit to any of those is reviewed
// on the light profile. A change the Mate judges weightier is still sent
// to the full profile with --model/--effort.
var codeExtensions = map[string]bool{
	".go": true, ".mod": true, ".sum": true,
	".js": true, ".mjs": true, ".cjs": true, ".jsx": true, ".ts": true, ".tsx": true, ".vue": true, ".svelte": true,
	".py": true, ".rb": true, ".rs": true, ".java": true, ".kt": true, ".kts": true, ".scala": true, ".swift": true,
	".c": true, ".h": true, ".cc": true, ".cpp": true, ".cxx": true, ".hpp": true, ".m": true, ".mm": true,
	".cs": true, ".fs": true, ".php": true, ".pl": true, ".lua": true, ".dart": true, ".ex": true, ".exs": true,
	".erl": true, ".hs": true, ".clj": true, ".zig": true, ".nim": true,
	".sh": true, ".bash": true, ".zsh": true, ".fish": true, ".ps1": true,
	".sql": true, ".proto": true, ".graphql": true,
	".mk": true, ".cmake": true, ".gradle": true,
}

// codeBasenames are the extension-less files that are code by name.
var codeBasenames = map[string]bool{
	"makefile": true, "dockerfile": true, "justfile": true, "rakefile": true, "gemfile": true,
}

// IsCodePath reports whether a changed path looks like code, by
// codeExtensions and codeBasenames. A path with none of them is not code:
// the list errs towards "not code" only for files a few lines of which the
// light reviewer can be trusted with (prose, data, configuration), and the
// caller treats any code file as taking the change off the small profile.
func IsCodePath(p string) bool {
	base := strings.ToLower(path.Base(strings.TrimSpace(p)))
	if codeBasenames[base] {
		return true
	}
	return codeExtensions[path.Ext(base)]
}
