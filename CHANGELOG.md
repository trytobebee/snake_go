# Changelog

All notable changes to this project will be documented in this file.

## [3.0.0] - 2026-01-13

### 🚀 Added
- **Intelligent AI (Auto-Play Mode)**:
  - Implemented Flood Fill algorithm for spatial awareness and trap avoidance.
  - Added utility-based greedy logic to prioritize high-value foods (corners/edges).
  - Added time-awareness: AI calculates if it can reach food before expiry and auto-triggers Boost.
- **Web Audio System**: 
  - Real-time sound synthesis for move, eat, boost, and game over using Web Audio API (zero external assets).
- **Mobile Experience & Haptics**: 
  - Integrated Vibration API for Android devices (haptic feedback on turns, eating, and crashes).
  - Added "Long Press to Boost" for mobile touch controls.
- **Infrastructure**:
  - Organized documentation into a dedicated `docs/` directory.
  - Implemented IP-based connection limiting to prevent duplicate game instances.

### ⚡ Changed
- Refined mobile UI layout to prevent Cumulative Layout Shift (CLS) when Boost is active.
- Improved difficulty/auto panel to fit better on small screens.
- Standardized coordinate system and direction validation to prevent 180-degree self-collisions.

---

## [2.0.0] - 2026-01-09
### Added
- Web-based version with High Performance Canvas rendering.
- Glassmorphism UI design.
- WebSocket-based real-time state synchronization.

---

## [1.0.0] - 2026-01-07
### Added
- Original Go terminal implementation.
- Multiple food types and scoring.
