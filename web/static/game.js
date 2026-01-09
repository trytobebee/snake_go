// Snake Game - Web Version
// WebSocket connection and rendering

class SnakeGameClient {
    constructor() {
        this.canvas = document.getElementById('gameCanvas');
        this.ctx = this.canvas.getContext('2d');
        this.ws = null;
        this.gameState = null;

        // Canvas settings
        this.cellSize = 20;
        this.boardWidth = 25;
        this.boardHeight = 25;

        this.canvas.width = this.boardWidth * this.cellSize;
        this.canvas.height = this.boardHeight * this.cellSize;

        // UI elements
        this.scoreEl = document.getElementById('score');
        this.speedEl = document.getElementById('speed');
        this.eatenEl = document.getElementById('eaten');
        this.boostIndicator = document.getElementById('boostIndicator');
        this.messageDisplay = document.getElementById('messageDisplay');
        this.gameOverlay = document.getElementById('gameOverlay');
        this.overlayTitle = document.getElementById('overlayTitle');
        this.overlayMessage = document.getElementById('overlayMessage');
        this.connectionStatus = document.getElementById('connectionStatus');

        this.setupWebSocket();
        this.setupKeyboard();
        this.render();
    }

    setupWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsURL = `${protocol}//${window.location.host}/ws`;

        this.ws = new WebSocket(wsURL);

        this.ws.onopen = () => {
            console.log('WebSocket connected');
            this.updateConnectionStatus('connected');
        };

        this.ws.onmessage = (event) => {
            this.gameState = JSON.parse(event.data);
            this.updateUI();
            this.render();
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.updateConnectionStatus('disconnected');
        };

        this.ws.onclose = () => {
            console.log('WebSocket closed');
            this.updateConnectionStatus('disconnected');
            setTimeout(() => this.setupWebSocket(), 3000);
        };
    }

    setupKeyboard() {
        document.addEventListener('keydown', (e) => {
            if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return;

            const key = e.key.toLowerCase();
            let action = '';

            // Direction keys
            if (key === 'arrowup' || key === 'w') action = 'up';
            else if (key === 'arrowdown' || key === 's') action = 'down';
            else if (key === 'arrowleft' || key === 'a') action = 'left';
            else if (key === 'arrowright' || key === 'd') action = 'right';
            else if (key === ' ' || key === 'p') action = 'pause';
            else if (key === 'r') action = 'restart';
            else if (key === 'q') action = 'quit';

            if (action) {
                e.preventDefault();
                this.ws.send(JSON.stringify({ action }));
            }
        });
    }

    updateUI() {
        if (!this.gameState) return;

        // Update stats
        this.scoreEl.textContent = this.gameState.score || 0;
        this.speedEl.textContent = (this.gameState.eatingSpeed || 0).toFixed(2);
        this.eatenEl.textContent = this.gameState.foodEaten || 0;

        // Boost indicator
        if (this.gameState.boosting) {
            this.boostIndicator.classList.add('active');
        } else {
            this.boostIndicator.classList.remove('active');
        }

        // Message display
        if (this.gameState.message) {
            this.messageDisplay.textContent = this.gameState.message;
            this.messageDisplay.classList.add('show');
            setTimeout(() => {
                this.messageDisplay.classList.remove('show');
            }, 3000);
        } else {
            this.messageDisplay.textContent = '';
        }

        // Game overlay
        if (this.gameState.gameOver) {
            this.showOverlay('💀 GAME OVER!', 'Press R to restart');
        } else if (this.gameState.paused) {
            this.showOverlay('⏸️ PAUSED', 'Press P to continue');
        } else {
            this.hideOverlay();
        }
    }

    showOverlay(title, message) {
        this.overlayTitle.textContent = title;
        this.overlayMessage.textContent = message;
        this.gameOverlay.classList.add('show');
    }

    hideOverlay() {
        this.gameOverlay.classList.remove('show');
    }

    updateConnectionStatus(status) {
        this.connectionStatus.className = `connection-status ${status}`;
        const statusText = this.connectionStatus.querySelector('.status-text');

        if (status === 'connected') {
            statusText.textContent = 'Connected';
        } else if (status === 'disconnected') {
            statusText.textContent = 'Disconnected';
        } else {
            statusText.textContent = 'Connecting...';
        }
    }

    render() {
        // Clear canvas
        this.ctx.fillStyle = '#1a1a2e';
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);

        if (!this.gameState) {
            // Show waiting message
            this.ctx.fillStyle = '#fff';
            this.ctx.font = '20px Inter, sans-serif';
            this.ctx.textAlign = 'center';
            this.ctx.fillText('Connecting to server...', this.canvas.width / 2, this.canvas.height / 2);
            return;
        }

        // Draw walls
        this.ctx.fillStyle = '#4a5568';
        for (let x = 0; x < this.boardWidth; x++) {
            this.drawCell(x, 0);
            this.drawCell(x, this.boardHeight - 1);
        }
        for (let y = 0; y < this.boardHeight; y++) {
            this.drawCell(0, y);
            this.drawCell(this.boardWidth - 1, y);
        }

        // Draw foods
        if (this.gameState.foods) {
            this.gameState.foods.forEach(food => {
                this.drawFood(food);
            });
        }

        // Draw snake
        if (this.gameState.snake) {
            this.gameState.snake.forEach((segment, index) => {
                if (index === 0) {
                    // Head
                    this.ctx.fillStyle = '#48bb78';
                    this.drawCell(segment.x, segment.y);
                    // Add eyes
                    this.ctx.fillStyle = '#000';
                    const centerX = segment.x * this.cellSize + this.cellSize / 2;
                    const centerY = segment.y * this.cellSize + this.cellSize / 2;
                    this.ctx.beginPath();
                    this.ctx.arc(centerX - 4, centerY - 4, 2, 0, Math.PI * 2);
                    this.ctx.arc(centerX + 4, centerY - 4, 2, 0, Math.PI * 2);
                    this.ctx.fill();
                } else {
                    // Body
                    this.ctx.fillStyle = '#68d391';
                    this.drawCell(segment.x, segment.y);
                }
            });
        }

        // Draw crash point
        if (this.gameState.gameOver && this.gameState.crashPoint) {
            this.ctx.font = `${this.cellSize}px sans-serif`;
            this.ctx.textAlign = 'center';
            this.ctx.textBaseline = 'middle';
            this.ctx.fillText('💥',
                this.gameState.crashPoint.x * this.cellSize + this.cellSize / 2,
                this.gameState.crashPoint.y * this.cellSize + this.cellSize / 2
            );
        }

        requestAnimationFrame(() => this.render());
    }

    drawCell(x, y) {
        this.ctx.fillRect(
            x * this.cellSize + 1,
            y * this.cellSize + 1,
            this.cellSize - 2,
            this.cellSize - 2
        );
    }

    drawFood(food) {
        const colors = {
            0: '#9f7aea', // Purple
            1: '#4299e1', // Blue
            2: '#ed8936', // Orange
            3: '#f56565'  // Red
        };

        const emojis = {
            0: '🟣',
            1: '🔵',
            2: '🟠',
            3: '🔴'
        };

        // Draw food as emoji
        this.ctx.font = `${this.cellSize - 4}px sans-serif`;
        this.ctx.textAlign = 'center';
        this.ctx.textBaseline = 'middle';
        this.ctx.fillText(
            emojis[food.foodType] || '🟣',
            food.pos.x * this.cellSize + this.cellSize / 2,
            food.pos.y * this.cellSize + this.cellSize / 2
        );

        // Draw timer if food has remainingSeconds
        if (food.remainingSeconds > 0 && food.remainingSeconds <= 5) {
            const timerEmojis = ['', '①', '②', '③', '④', '⑤'];
            this.ctx.font = `${this.cellSize / 2}px sans-serif`;
            this.ctx.fillText(
                timerEmojis[food.remainingSeconds],
                (food.pos.x + 1) * this.cellSize - 4,
                food.pos.y * this.cellSize + this.cellSize / 2
            );
        }
    }
}

// Initialize game when page loads
document.addEventListener('DOMContentLoaded', () => {
    new SnakeGameClient();
});
