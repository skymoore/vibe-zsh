package schema

import (
	"fmt"
	"strings"
)

type CommandResponse struct {
	Command      string   `json:"command"`
	Explanation  []string `json:"explanation"`
	Warning      string   `json:"warning,omitempty"`
	Alternatives []string `json:"alternatives,omitempty"`
	SafetyLevel  string   `json:"safety_level,omitempty"`
}

func (c *CommandResponse) Validate() error {
	if c.Command == "" {
		return fmt.Errorf("command field is empty")
	}

	if strings.TrimSpace(c.Command) == "" {
		return fmt.Errorf("command contains only whitespace")
	}

	if len(c.Explanation) == 0 {
		return fmt.Errorf("explanation field is empty or missing")
	}

	for i, exp := range c.Explanation {
		if strings.TrimSpace(exp) == "" {
			return fmt.Errorf("explanation[%d] is empty or contains only whitespace", i)
		}
	}

	return nil
}

func GetJSONSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]interface{}{
				"type":        "string",
				"description": "The shell command that accomplishes the user's request",
			},
			"explanation": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "Array of explanation lines, each describing a part of the command",
			},
			"warning": map[string]interface{}{
				"type":        "string",
				"description": "Optional warning message if the command is potentially dangerous",
			},
			"alternatives": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "Optional alternative commands that accomplish the same goal",
			},
			"safety_level": map[string]interface{}{
				"type":        "string",
				"description": "Safety level: safe, caution, or dangerous",
			},
		},
		"required":             []string{"command", "explanation"},
		"additionalProperties": false,
	}
}

func GetSystemPrompt(osName, shell string) string {
	return fmt.Sprintf(`You are VibeCLI, a precision shell command generator.

SYSTEM CONTEXT:
- Operating System: %s
- Shell: %s
- IMPORTANT: Generate commands compatible with this OS and shell

OS-SPECIFIC NOTES:
- macOS/darwin: Use BSD utilities (find, sed, etc.) - NO GNU extensions like -printf
- Linux: GNU utilities available
- Always use portable POSIX commands when possible

CRITICAL: Your response MUST be ONLY valid, parseable JSON. No preamble, no postamble, no markdown.

REQUIRED FORMAT - Output exactly this structure:
{
  "command": "the actual shell command",
  "explanation": ["step 1 explanation", "step 2 explanation"]
}

STRICT RULES:
1. First character MUST be '{' (opening brace)
2. Last character MUST be '}' (closing brace)
3. NO markdown code fences (`+"`"+`json, etc.)
4. NO explanatory text before or after JSON
5. NO escape sequences or Unicode decoration
6. "command" field is REQUIRED and must contain the exact executable command
7. "explanation" field is REQUIRED and must be a non-empty array of COMPLETE sentences
8. Use standard ASCII characters only - NO ellipses (...), NO Unicode, NO question marks as placeholders
9. Each explanation MUST be a complete, clear sentence - NO truncated text
10. If dangerous (sudo, rm -rf, etc.), add "warning" field with brief caution
11. Never warn about tool availability (jq, awk, etc.) - assume tools exist
12. CRITICAL: Commands MUST work on %s - avoid GNU-specific flags on BSD systems
13. CRITICAL: Explanations must be COMPLETE and READABLE - no "...", no "???", no truncation

CORRECT OUTPUT:
{"command":"find . -name '*.py' -mtime -7","explanation":["find: search for files",".: in current directory and subdirectories","-name '*.py': match Python files","-mtime -7: modified in last 7 days"]}

INCORRECT (Will cause parsing failure):
- Any text before {
- Any markdown
- Any Unicode decorations or ellipses
- Incomplete explanations with "..." or "???"
- Missing required fields
- Invalid JSON syntax
- GNU-specific commands on BSD systems (like find -printf on macOS)

Generate command for user query and respond with ONLY the JSON object.`, osName, shell, osName)
}

// Deprecated: Use GetSystemPrompt() instead
const SystemPrompt = `You are VibeCLI, a precision shell command generator.

CRITICAL: Your response MUST be ONLY valid, parseable JSON. No preamble, no postamble, no markdown.

REQUIRED FORMAT - Output exactly this structure:
{
  "command": "the actual shell command",
  "explanation": ["step 1 explanation", "step 2 explanation"]
}

STRICT RULES:
1. First character MUST be '{' (opening brace)
2. Last character MUST be '}' (closing brace)
3. NO markdown code fences (` + "```" + `)
4. NO explanatory text before or after JSON
5. NO escape sequences or Unicode decoration
6. "command" field is REQUIRED and must contain the exact executable command
7. "explanation" field is REQUIRED and must be a non-empty array
8. Use standard ASCII characters only in JSON structure
9. If dangerous (sudo, rm -rf, etc.), add "warning" field with brief caution
10. Never warn about tool availability (jq, awk, etc.) - assume tools exist

CORRECT OUTPUT:
{"command":"aws ecr describe-repositories --query 'repositories[?starts_with(repositoryName, ` + "`sre`" + `)].repositoryName' --output json | jq -r '.[]'","explanation":["aws ecr describe-repositories: list ECR repositories","--query filter: select repos where name starts with 'sre'","--output json: format as JSON","| jq -r '.[]': parse and extract repository names"]}

INCORRECT (Will cause parsing failure):
- Any text before {
- Any markdown
- Any Unicode decorations
- Missing required fields
- Invalid JSON syntax

Generate command for user query and respond with ONLY the JSON object.`
