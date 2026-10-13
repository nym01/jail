package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Run(code string) result {
	dir, err := os.MkdirTemp("", "jail-*")
	if err != nil {
		return result{Stderr: err.Error(), ExitCode: -1}
	}

	config := map[string]any{
		"ociVersion": "1.0.0",
		"process": map[string]any{
			"terminal": false,
			"user":     map[string]any{"uid": 0, "gid": 0},
			"args":     []string{"python3", "-c", code},
			"env":      []string{"PATH=/usr/bin:/bin"},
			"cwd":      "/",
		},
		"root":     map[string]any{"path": "/rootfs", "readonly": true},
		"hostname": "jail",
		"mounts": []map[string]any{
			{"destination": "/proc", "type": "proc", "source": "proc"},
		},
		"linux": map[string]any{
			"namespaces": []map[string]string{
				{"type": "pid"},
				{"type": "network"},
				{"type": "ipc"},
				{"type": "uts"},
				{"type": "mount"},
			},
		},
	}

	data, err := json.Marshal(config)
	if err != nil {
		return result{Stderr: err.Error(), ExitCode: -1}
	}

	err = os.WriteFile(filepath.Join(dir, "config.json"), data, 0644)
	if err != nil {
		return result{Stderr: err.Error(), ExitCode: -1}
	}

	return result{}
}
