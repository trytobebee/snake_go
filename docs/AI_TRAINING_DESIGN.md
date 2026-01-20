# Snake AI Machine Learning Design (v2.0 - RL)

This document outlines the reinforcement learning design and architecture used for the Snake AI.

## 1. Input Representation (State S)
A multi-layered grid (tensor) representation of the 25x25 game board.
- **Shape**: (6, 25, 25)
- **Channels**:
    1. **Player Head**: Position of the AI-controlled snake's head.
    2. **Player Body**: Positions of all AI snake body segments.
    3. **Enemy Head**: Position of the human opponent's head.
    4. **Enemy Body**: Positions of all human opponent body segments.
    5. **Food**: Locations of all active food items (Standard, Bonus, etc.).
    6. **Hazards**: Locations of physical obstacles and active fireballs.

## 2. Output Representation (Action Q-values)
Four continuous values representing the expected long-term reward for each movement action.
- **Actions**: [0: UP, 1: DOWN, 2: LEFT, 3: RIGHT]
- **Policy**: `Action = argmax(Q_values)`

## 3. Training Logic: Deep Q-Learning (DQN)
The model is trained using **Offline Reinforcement Learning** from JSONL session recordings.

### Loss Function
**Mean Squared Error (MSE)** on the Temporal Difference (TD) error:
`Loss = MSE(Q(s, a), Reward + Gamma * max(Q_target(s_next)) * (1 - Done))`

### Key Hyperparameters
- **Gamma (Discount Factor)**: 0.95 (Ensures AI values long-term rewards like growth).
- **Learning Rate**: 0.0005 (Adam Optimizer).
- **Target Update Frequency**: Every 2 epochs (Ensures training stability).
- **Batch Size**: 64.

## 4. Reward Structure (Aligned with Game Score)
The training reward is derived directly from the delta of the game score:
- **Food Eaten**: **+10.0 to +40.0** (Varies by food rarity).
- **Survival Bonus**: **+0.1** (Per 16ms tick to encourage longevity).
- **Collision Penalty**: **-100.0** (High penalty for game over).
- **Combat Success**: **+50.0** (AI Headshot), **+20.0** (AI Body Hit).
- **Destruction**: **+10.0** (Per obstacle destroyed).
- **Incoming Hits**: **-30.0** (Headshot), **-10.0** (Body hit).

## 5. Inference & Deployment
- **Engine**: ONNX Runtime (Go Bindings).
- **Service**: **Centralized Global Worker Pattern**.
  - All game sessions push requests to a shared `PredictionQueue`.
  - A single dedicated worker processes requests to maximize CPU/Cache efficiency.
- **Latency**: ~1.3ms end-to-end.

## 6. Safety Layer (Physic Override)
To ensure robustness, a **post-inference safety check** is performed:
1. If the NN suggests a move that leads to immediate collision (Wall/Body), the move is **discarded**.
2. The system falls back to the **Heuristic AI (Flood-fill)** to find the safest alternative move.
3. This provides a safety net against model "hallucinations" in tight spaces.
