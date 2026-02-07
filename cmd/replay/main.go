package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/trytobebee/snake_go/pkg/config"
	"github.com/trytobebee/snake_go/pkg/game"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ReplayServer handles serving replay UI and data
type ReplayServer struct {
	addr      string
	recordDir string
}

func main() {
	server := &ReplayServer{
		addr:      ":8081",
		recordDir: "records",
	}

	// Serve static files (REUSE existing web/static)
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve replay specific UI
	http.HandleFunc("/", server.handleIndex)
	http.HandleFunc("/view", server.handleView)

	// WebSocket for replay data
	http.HandleFunc("/ws/replay", server.handleReplayWS)

	fmt.Printf("📼 Snake Replay Tool starting on http://localhost%s\n", server.addr)
	log.Fatal(http.ListenAndServe(server.addr, nil))
}

type RecordFile struct {
	Name      string
	Size      int64
	Time      time.Time
	SessionID string
}

func (s *ReplayServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(s.recordDir)
	if err != nil {
		// Try creating dir if not exists
		os.MkdirAll(s.recordDir, 0755)
		files = []os.DirEntry{}
	}

	var records []RecordFile
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".jsonl" {
			info, _ := f.Info()
			// expecting format: game_{sessionID}_{timestamp}.jsonl
			parts := strings.Split(f.Name(), "_")
			sessID := ""
			if len(parts) >= 2 {
				sessID = parts[1]
			}
			records = append(records, RecordFile{
				Name:      f.Name(),
				Size:      info.Size(),
				Time:      info.ModTime(),
				SessionID: sessID,
			})
		}
	}

	// Sort by time desc
	sort.Slice(records, func(i, j int) bool {
		return records[i].Time.After(records[j].Time)
	})

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Snake Replays</title>
    <style>
        body { font-family: monospace; background: #1a202c; color: #fff; padding: 2rem; }
        h1 { color: #48bb78; }
        .file-list { display: grid; gap: 1rem; }
        .file-item { 
            background: #2d3748; padding: 1rem; border-radius: 8px; 
            display: flex; justify-content: space-between; align-items: center;
        }
        .file-item:hover { background: #4a5568; }
        a { color: #63b3ed; text-decoration: none; font-weight: bold; }
        .meta { color: #a0aec0; font-size: 0.9em; }
    </style>
</head>
<body>
    <h1>📼 Replay Library</h1>
    <div class="file-list">
        {{range .}}
        <div class="file-item">
            <div>
                <div class="name">{{.Name}}</div>
                <div class="meta">Session: {{.SessionID}} | Size: {{.Size}} bytes | {{.Time.Format "2006-01-02 15:04:05"}}</div>
            </div>
            <a href="/view?file={{.Name}}">WATCH REPLAY ▶</a>
        </div>
        {{else}}
        <p>No recordings found in ./records/</p>
        {{end}}
    </div>
</body>
</html>`

	t, _ := template.New("index").Parse(tmpl)
	t.Execute(w, records)
}

func (s *ReplayServer) handleView(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// Redirect to the static HTML page with the file parameter
	http.Redirect(w, r, "/static/replay.html?file="+filename, http.StatusFound)
}

// Websocket logic
func (s *ReplayServer) handleReplayWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	filename := r.URL.Query().Get("file")
	path := filepath.Join(s.recordDir, filename)
	file, err := os.Open(path)
	if err != nil {
		log.Println("Failed to open record:", err)
		return
	}
	defer file.Close()

	// Parse Record (First line to get config if possible, or use default)
	// For now, send default config immediately to initialize frontend
	defaultConfig := game.GameConfig{
		Width:            config.StandardWidth,
		Height:           config.StandardHeight,
		GameDuration:     int(config.GameDuration.Seconds()),
		FireballCooldown: int(config.FireballCooldown.Milliseconds()),
	}
	if err := conn.WriteJSON(struct {
		Type   string           `json:"type"`
		Config *game.GameConfig `json:"config"`
	}{
		Type:   "config",
		Config: &defaultConfig,
	}); err != nil {
		return
	}

	scanner := bufio.NewScanner(file)

	// Control vars
	paused := false

	// Read Loop for controls
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Simple command parsing
			var cmd struct {
				Command string `json:"command"`
			}
			json.Unmarshal(msg, &cmd)
			if cmd.Command == "pause" {
				paused = true
			}
			if cmd.Command == "resume" {
				paused = false
			}
		}
	}()

	// Stream Loop
	for scanner.Scan() {
		line := scanner.Bytes()

		// Parse Record
		var rec game.StepRecord
		if err := json.Unmarshal(line, &rec); err != nil {
			log.Println("JSON parse error:", err)
			continue
		}

		// Replay speed control
		for paused {
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond) // Fixed 10fps playback

		// Construct message (inject intent into meta)
		msg := struct {
			Type  string         `json:"type"`
			State game.GameState `json:"state"`
			Meta  map[string]any `json:"meta"`
		}{
			Type:  "state",
			State: rec.State,
			Meta: map[string]any{
				"step":   rec.StepID,
				"intent": rec.AIContext.Intent,
			},
		}

		if err := conn.WriteJSON(msg); err != nil {
			break
		}
	}
}
