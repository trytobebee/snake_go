import { SnakeGameClient } from './game.js';

class ReplayClient extends SnakeGameClient {
    constructor() {
        super();
        this.isPaused = false;
        this.replayWS = null;

        // Hide unrelated UI (Mode/Diff/Controls are inherently hidden/overridden by style or html)

        // Get filename from URL
        const params = new URLSearchParams(window.location.search);
        this.filename = params.get('file');

        if (this.filename) {
            document.getElementById('filename-display').textContent = this.filename;
            this.setupReplayWS();
        } else {
            document.getElementById('filename-display').textContent = "No file selected";
        }

        // Controls
        const pauseBtn = document.getElementById('btn-pause');
        if (pauseBtn) {
            pauseBtn.onclick = () => {
                this.isPaused = !this.isPaused;
                if (this.replayWS && this.replayWS.readyState === WebSocket.OPEN) {
                    this.replayWS.send(JSON.stringify({ command: this.isPaused ? "pause" : "resume" }));
                }
            };
        }
    }

    // Override standard setup
    setupWebSocket() { /* No-op to prevent connecting to live server */ }
    setupKeyboard() { /* No-op */ }
    setupMobileControls() { /* No-op */ }
    setupDifficulty() { /* No-op */ }
    setupMode() { /* No-op */ }
    setupAutoPlay() { /* No-op */ }

    setupReplayWS() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        // Connect to replay endpoint
        this.replayWS = new WebSocket(protocol + '//' + window.location.host + '/ws/replay?file=' + encodeURIComponent(this.filename));

        this.replayWS.onopen = () => console.log("Replay connected");

        this.replayWS.onmessage = (evt) => {
            const msg = JSON.parse(evt.data);

            if (msg.type === 'config') {
                this.handleConfig(msg.config);
                // Force UI update and hide overlay on config load
                this.updateUI();
                if (this.gameOverlay) this.gameOverlay.style.display = 'none';

            } else if (msg.type === 'state') {
                this.gameState = msg.state;
                this.updateUI();

                // Force hide overlay during playback unless game over
                if (this.gameOverlay && !this.gameState.gameOver) {
                    this.gameOverlay.style.display = 'none';
                }

                // Update extra replay info
                if (msg.meta) {
                    const stepEl = document.getElementById('replayInfo');
                    const intentEl = document.getElementById('aiIntent');
                    if (stepEl) stepEl.textContent = "Step: " + msg.meta.step;
                    if (intentEl) intentEl.textContent = "INTENT: " + (msg.meta.intent || "IDLE");
                }
            }
        };

        this.replayWS.onerror = (err) => console.error("Replay WS Error:", err);
        this.replayWS.onclose = () => console.log("Replay WS Closed");
    }
}

// Init
document.addEventListener('DOMContentLoaded', () => {
    new ReplayClient();
});
