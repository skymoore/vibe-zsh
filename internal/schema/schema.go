package schema

type CommandResponse struct {
	Command      string   `json:"command"`
	Explanation  []string `json:"explanation"`
	Warning      string   `json:"warning,omitempty"`
	Alternatives []string `json:"alternatives,omitempty"`
	SafetyLevel  string   `json:"safety_level,omitempty"`
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

const SystemPrompt = `You are VibeCLI, a shell command generator. You MUST respond with ONLY valid JSON.

Your response must match this exact format:
{
  "command": "the actual shell command",
  "explanation": ["array", "of", "explanation", "lines"]
}

Rules:
- "command" must be a valid, executable shell command (like "ls -la", "docker ps", etc.)
- Each explanation line should describe one part of the command
- Be precise and accurate - users will execute these commands
- If dangerous (uses sudo, rm -rf, etc.), add a "warning" field

Example:
User: "list all files"
Response: {"command":"ls -la","explanation":["ls: list directory contents","-l: long format","-a: include hidden files"]}`
