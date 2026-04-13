# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This repo has two distinct applications:

1. **Cribbage hand analysis tool** — a Go backend that pre-computes scoring statistics for all possible cribbage hands, paired with a vanilla JS/HTML frontend for interactive exploration (`tools/`, `web/`, `scores/`)
2. **Cribbage AI engine** — a pluggable strategy simulation framework for running AI opponents against each other to compare algorithms (`engine/`, `strategy/`, `sim/`)

## Go Commands

```bash
# Run the AI simulation (default 1000 games per matchup)
go run ./sim
go run ./sim 5000

# Build/run hand analysis tools
go run ./tools/score -- as 2h 3d 4c 5s
go run ./tools/gen       # Generates all pre-computed hand summaries (slow, writes to scores/)

# Run tests (currently none exist)
go test ./...
```

## AI Engine Architecture

### Interfaces (`strategy/strategy.go`)

Two independent pluggable interfaces — mix and match freely:
- `Discarder` — `Discard(hand Cards, dealerCrib bool) (keep, crib Cards)`
- `Pegger` — `Play(hand Cards, state PeggingState) Card`
- `Strategy` — composes both; `NewStrategy(name, discarder, pegger)` creates a named `Player`

### Game Engine (`engine/`)

`engine.RunGame(p0, p1, dealer, rng)` runs a full game to 121 points, returning `GameResult` with wins, pegged/hand/crib points, and diff per player. Pegging logic lives in `engine/pegging.go`; nibs (Jack cut) and hand/crib scoring in `engine/engine.go`.

### Discard Strategies (`strategy/discard/`)

All summary-based strategies share a `SummaryCache` (map of all C(52,4)=270,725 four-card hand averages, built once at startup in ~5s). `MaxAvgDiff` additionally uses a `TwoCribCache` (all C(52,2)=1,326 two-card average crib scores, built in parallel at startup).

| Strategy | Description |
|---|---|
| `MaxAvg` | Keeps 4 cards with highest average hand score across all cuts |
| `MaxAvgDiff` | Maximizes `avg4 ± avg2`: adds avg crib value if dealer, subtracts if not |
| `MaxMin/MaxMedian/MaxMax/MaxMode` | Like MaxAvg but optimize different summary stats |
| `MinValue/MaxValue/Random` | Naive baselines — commented out in sim |

### Pegging Strategies (`strategy/peg/`)

| Strategy | Description |
|---|---|
| `MaxNext` | Plays card that scores most points immediately |
| `MaxSetup` | Takes points if available (like MaxNext), otherwise 1-ply lookahead for best future score |
| `MinValue/MaxValue/Random` | Naive baselines — commented out in sim |

### Simulation (`sim/main.go`)

Runs every combination of active discard × peg strategies against each other. Uses mirror optimization (only upper triangle of matchup matrix, O(n²/2)), concurrent goroutines per matchup. Output: per-matchup stat table (wins, pegged/hand/crib/total/diff) + ranked summary.

When adding or comparing strategies: comment out weak ones in `allCombinations` rather than deleting. Active strategies are the uncommented entries in the `discards` and `pegs` slices.

### Caches

Both caches are built at sim startup and injected into strategies:
- `discard.NewSummaryCache()` — ~5s, all 4-card hand averages
- `discard.NewTwoCribCache()` — parallel build (1,326 goroutines), all 2-card average crib scores

`ScorePeggingPlay` lives in `game/pegging.go` (not engine) so pegging strategies can use it without a circular import.

## Hand Analysis Tool Architecture

### Data Flow

1. `tools/gen` pre-computes scores for all 270,725 possible 4-card hands across 48 cuts → writes to `scores/` (gitignored, JSON + binary formats)
2. Frontend JS (`web/js/`) reads these pre-computed summaries to render statistics without a running server
3. For 6-card analysis, the frontend (or `game/counts/sixhands.go`) breaks all 15 possible 4-card sub-hands to find optimal discards

### Key Packages

- **`game/`** — Core cribbage logic: Card (ID = face×4 + suit, range 0–51), deck, hand scoring (`CountCards`), pegging scoring (`ScorePeggingPlay`), combinatorics helpers (`ChoseFour`, `ChoseTwo`, `ChooseFourWithRemaining`)
- **`game/math/`** — Binomial coefficients, `CombinationIndex`/`IndexToCombination`
- **`game/counts/`** — `FourSummary` and `SixHands` structs with custom binary marshaling; stats include avg, min/max, median, mode, std dev
- **`util/`** — File path constants and JSON/binary I/O helpers

### Frontend

Three HTML entry points: `index.html` (landing), `four/index.html`, `six/index.html`. No build step — pure static files. `web/js/combinatorics.js` mirrors the Go math package client-side.

### Card Representation

Cards are 2-char strings (`as`, `2h`) and integer IDs (0–51). The `game/` package has a global card cache initialized at `init()` time.

## Legacy Code

`/java/` contains an older Java implementation — not actively developed.
