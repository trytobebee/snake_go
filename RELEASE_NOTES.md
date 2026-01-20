# 🧠 Release Notes - v4.0.0

**Release Date**: January 20, 2026

## 🎉 The Intelligence Evolution

This milestone marks the transition from heuristic algorithms to **Deep Reinforcement Learning**. The Snake game now features a "brain" powered by a Convolutional Neural Network (CNN) running on the high-performance **ONNX Runtime**.

## 🚀 Major Updates

### 1. Neural AI Integration
- **Deep Reinforcement Learning**: AI is no longer just pathfinding; it's making strategic decisions based on a 6-channel spatial grid (Players, AI, Food, Obstacles, Fireballs).
- **ONNX Runtime (Go Bindings)**: Native C++ acceleration for inference, ensuring that the AI can "think" as fast as the game can move.
- **Micro-Latency Optimization**: Designed a dedicated inference service that processes requests in **~1.3ms**, enabling real-time competitive play.

### 2. High-Performance Architecture
- **Global Inference Worker**: A singleton worker pattern that eliminates memory overhead and resource contention.
- **Asynchronous Task Queue**: Uses Go channels to pipeline AI requests, maximizing CPU cache locality and eliminating mutex locks during inference.
- **Singleton Model Instance**: One model serves all users, drastically reducing the server's memory footprint.

### 3. AI Safety & Robustness
- **Heuristic-Neural Hybrid**: If the neural network suggests a move that leads to immediate self-destruction, a safety interceptor overrides it using "Old School" heuristic logic.
- **Collision Look-ahead**: Real-time validation of every AI move through the game's physical engine.

### 4. Machine Learning Pipeline
- **JSONL Data Capture**: Seamlessly record game sessions into structured data for continuous training.
- **PyTorch to ONNX**: Full pipeline from training in Python to high-speed deployment in Go.

---

## 📊 Technical Comparison

| Metric | Heuristic AI (Legacy) | Deep Learning AI (v4.0.0) |
| :--- | :--- | :--- |
| **Logic Type** | Rule-based (Flood-fill/Greedy) | Neural Network (CNN + DQN) |
| **Inference Time** | <0.1ms | ~1.3ms (including queue overhead) |
| **Spatial Awareness** | Static distance calculations | 6-Channel Grid Analysis |
| **Strategic Vision** | Short-term (Greedy) | Long-term (Value-based Q-learning) |

---

## 🛠️ Infrastructure Improvements
- **Project Modularization**: Further separated core game logic from the AI inference service.
- **Dependency Clean-up**: Optimized `go.mod` and removed unused packages.
- **Multi-user Stability**: Fixed race conditions in ONNX initialization through `sync.Once` patterns.

## 📝 Full Changelog

### Added
- `pkg/game/ai_model.go`: The heart of the ONNX inference service.
- `ml/train.py`: Training script for the DQN model.
- High-performance `Predict` queue system.
- Collision safety interceptor.

### Changed
- Refactored `pkg/game/ai.go` to leverage the new global Predict service.
- Updated `pkg/game/game.go` to support global worker initialization.
- Optimized 16ms `BaseTick` synchronization.

---

**Download**: [Latest Release](https://github.com/trytobebee/snake_go/releases/latest)
**Documentation**: [README.md](./README.md)
**Architecture**: [docs/AI_ARCHITECTURE.md](./docs/AI_ARCHITECTURE.md)
