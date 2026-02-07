# Changelog

All notable changes to the **Snake Go Premium** project will be documented in this file.

## [6.0.0] - 2026-02-07

### ⚔️ Combat System Evolution
- **Scatter Shot (🌟)**: New power-up allowing 45-degree diagonal fire. Each shot launches 3 projectiles in a trident pattern.
- **Rapid Fire (⚡)**: New power-up that doubles bullet velocity and halves firing cooldown.
- **Diagonal Physics**: Implemented high-precision diagonal movement and collision detection for projectiles.

### 💰 Economy & Props
- **Weighted Spawn Engine**: Introduced a tiered probability system for items.
  - **Small Chests (💰)** & **Money Bags**: High-frequency drops for score padding.
  - **Power-ups (🛡️, ⚡, 🌟, etc.)**: Balanced mid-tier drops.
  - **Big Chests (👑)**: Rare high-value rewards (+120 pts).
- **Consolidated Prop Logic**: Merged diverse item handlers into a unified `Prop` package in `pkg/game/prop.go`.

### 🛡️ Core Robustness & UX
- **Broadcast Optimization**: Refactored `broadcastSessionCount` to be non-blocking. This prevents server "hangs" when a client has a poor connection during login/logout.
- **Atomic Auth Sequence**: Redesigned the "Kick Old Session" logic to avoid deadlocks in the client tracking mutex.
- **Asset Cache Busting**: Incremented Web versioning to **v2.5** to force browsers to fetch the latest rendering logic for new prop emojis.
- **Graceful Error Handling**: Added `EffectNone` constant to prevent potential logic holes or UI artifacts with future item additions.

---

## [5.1.0] - 2026-02-02

### 🚀 Major Improvements
- **Network Protocol Upgrade (Protobuf)**: Transitioned from JSON to Protobuf for all WebSocket communications. 
  - Reduced per-frame bandwidth by **~80%** (from 2.6KB down to 0.5KB for typical states).
  - Improved serialization/deserialization speed by ~6x.
- **Scalability Controls**:
  - Implemented a **500-player maximum limit** per server instance to prevent resource exhaustion.
  - Added "Server Full" graceful rejection logic with user-facing alerts.
- **Single Source of Truth (SSoT)**: Refactored `.proto` management. The server now dynamically serves the master `.proto` file to the frontend, eliminating synchronization bugs between backend and client.

### 📩 Feedback System
- **In-game Feedback**: Added a premium "💬 Feedback" floating button with a Glassmorphism UI for users to submit suggestions and bug reports.
- **Admin Dashboard**: New endpoint at `/admin/feedback` for developers to review user submissions in a clean, dark-themed interface.
- **Feishu (Lark) Integration**: Real-time push notifications to Feishu via Webhooks.
  - Supports **Interactive Card** format with headers, markdown, and formatted timestamps.

### 🛡️ Security & Ops
- **Environment Configuration**: Added `.env` file support using `godotenv`. sensitive keys like Feishu Webhook URLs are now kept out of the source code.
- **Hotkey Conflict Resolution**: Fixed a bug where typing in feedback or login fields would trigger game controls (e.g., restarting or pausing).

### 🐛 Bug Fixes
- **Player Count Persistence**: Fixed a race condition where refreshing the page could lead to ghost sessions and an ever-increasing player count.

---
