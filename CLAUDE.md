# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **Cribbage hand analysis tool** — a Go backend that pre-computes scoring statistics for all possible cribbage hands, paired with a vanilla JS/HTML frontend for interactive exploration.

## Go Commands

```bash
# Build a specific tool
go build ./tools/score
go build ./tools/gen
go build ./tools/dump
go build ./tools/read
go build ./tools/id

# Run a tool directly
go run ./tools/score -- as 2h 3d 4c 5s
go run ./tools/gen       # Generates all pre-computed hand summaries (slow, writes to scores/)

# Run tests (currently none exist)
go test ./...
```

No Makefile or build scripts — each tool is a standalone binary built with `go build`.

## Architecture

### Data Flow

1. `tools/gen` pre-computes scores for all 270,725 possible 4-card hands across 48 cuts → writes to `scores/` (gitignored, JSON + binary formats)
2. Frontend JS (`web/js/`) reads these pre-computed summaries to render statistics without a running server
3. For 6-card analysis, the frontend (or `game/counts/sixhands.go`) breaks all 15 possible 4-card sub-hands to find optimal discards

### Key Go Packages

- **`game/`** — Core cribbage logic: Card (ID = face×4 + suit, range 0–51), deck, hand scoring (`CountCards` counts fifteens, pairs, runs, flushes, nobs)
- **`game/math/`** — Combinatorics: binomial coefficients, `CombinationIndex`/`IndexToCombination` for converting card sets to canonical indices
- **`game/counts/`** — `FourSummary` and `SixHands` structs with custom binary marshaling for compact storage; stats include avg, min/max, median, mode, std dev
- **`util/`** — File path constants (`util/files.go`) and JSON/binary I/O helpers

### Frontend

Three HTML entry points: `index.html` (landing), `four/index.html` (4-card scoring), `six/index.html` (6-card analysis). No build step — pure static files.

- `web/js/cards.js` — Card rendering via CSS sprite sheet (`web/img/cards.png`), selection state
- `web/js/combinatorics.js` — Client-side combination utilities (mirrors the Go math package)
- `web/js/scoreFour.js` / `scoreSix.js` — UI logic for each variant

### Data Storage

`scores/` (gitignored) holds pre-computed data split across hashed files:
- `scores/four/` — 17 files, keyed by SHA1 hash of hand ID
- `scores/six/` — 47 files, keyed similarly

File path constants live in `util/files.go`.

### Card Representation

Cards are represented as 2-char strings (`as` = ace of spades, `2h` = 2 of hearts) and integer IDs (0–51). The `game/` package has a global card cache initialized at `init()` time.

## Legacy Code

`/java/` contains an older Java implementation of the same game logic — not actively developed.
