# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Cribbage tools in one Go module (`mattlove.dev/crib`):

1. **Hand analyzer** (`site/`): a static HTML/JS site, live at mattlove.dev/crib, that scores 4-card hands across all cuts and ranks the 15 discard options of a 6-card hand. All scoring runs in the browser (`site/web/js/cards.js`); there is no server or pre-computed data.
2. **AI engine** (`game/`, `engine/`, `strategy/`, `sim/`): pluggable discard/pegging strategies and a simulator that plays them against each other to compare algorithms.
3. **Human-vs-AI game** (`play/`): a local Go web server that uses the engine. Not deployed.

See README.md for the user-facing overview.

## Commands

```bash
# Simulate every active strategy combination (default 1000 games per matchup)
go run ./sim
go run ./sim 5000

# Local human-vs-AI game at http://localhost:8080 (~10s startup building caches)
go run ./play

# Serve the hand analyzer locally at http://localhost:8000
python3 -m http.server -d site 8000

# Checks (both run in CI before every deploy)
go vet ./...
go run ./scripts/parity | node scripts/parity/check.js

# Score a hand from the command line (legacy dev tool)
go run ./legacy/tools/score 5c 5d 5h js         # 4 cards: stats across all 48 cuts
go run ./legacy/tools/score 5c 5d 5h js 6c 7c   # 6 cards: all 15 keep/discard options, best avg first
```

There are no Go tests yet (`go test ./...` finds none).

## Deployment

`.github/workflows/pages.yml` deploys **only `site/`** to GitHub Pages (Pages source: "GitHub Actions") on pushes to `main` that touch `site/`, `game/`, `scripts/parity/`, or the workflow. A `check` job (`go vet` + parity) must pass first. Every push to `main` that matches those paths is a production deploy.

## Scoring Code Exists Twice: Go and JS

`site/web/js/cards.js` re-implements the Go card model, `CountCards`, per-cut scoring, `makeSummaries` and `makeSixHands` so the site can run without a server. **Any change to scoring or statistics in `game/` or `game/counts/` must be mirrored in `cards.js`**, and vice versa. `scripts/parity` compares all 9 stats on every 4-card hand (270,725) plus all discards of 2,000 fixed-seed 6-card hands, and exits non-zero on any difference.

Conventions both sides follow: per-cut scores sorted numerically; Mode ties go to the lower score; Avg, ModeP and StdDev rounded to 2 decimals.

## Card Representation

- Integer IDs 0–51: `id = face*4 + suit`, face 0–12 (ace..king), suit 0–3 (clubs, diamonds, hearts, spades).
- 2-char strings: face `a23456789tjqk` + suit `cdhs`, e.g. `5h`, `tc`, `js`. The site's URL hashes use these (`four/#5c5d5hjs`).
- `game/` has a global card table built in `init()`; use `game.CardById`, `game.CardByIdString`, `game.NewDeck()`.

## AI Engine Architecture

### Interfaces (`strategy/strategy.go`)

Two independent pluggable interfaces, mixed and matched freely:
- `Discarder`: `Discard(hand Cards, dealerCrib bool) (keep, crib Cards)`
- `Pegger`: `Play(hand Cards, state PeggingState) Card`. `hand` is the unplayed cards; the engine guarantees at least one legal play. **Never append to `state.Series` directly**: it aliases the engine's slice. Copy it or use a full slice expression (`s[:len(s):len(s)]`), as `MaxNext`/`MaxSetup` do.
- `Strategy` composes both; `NewStrategy(name, discarder, pegger)` creates a named `Player`.

### Game Engine (`engine/`)

- `engine.RunGame(p0, p1, dealer, rng)` plays a full game to 121 and returns a `GameResult` (winner, scores, pegged/hand/crib points per player). Results are deterministic for a given rng seed and deterministic strategies.
- `engine.Pegging` (`engine/pegging.go`) holds the pegging state and rules and is **shared by the simulator and `play/`**. Drive it by calling `Resolve()` until it returns no award (it handles passing the turn, go points, resets and the last-card point, stopping after each award so callers can check for 121), stop if `Done()`, otherwise get a card from the `Current` player and call `Play(card)`. Change pegging rules here, not in `play/`.
- Nibs (cut Jack) and hand/crib scoring order (non-dealer hand, dealer hand, crib) are in `engine/engine.go`. `play/main.go` has its own copy of this round flow.
- `game.ScorePeggingPlay` lives in `game/pegging.go` (not `engine/`) so strategies can use it without a circular import.

### Discard Strategies (`strategy/discard/`)

Summary-based strategies share a `SummaryCache` (all C(52,4) = 270,725 four-card hand summaries, ~5s to build). `MaxAvgDiff` also uses a `TwoCribCache` (average crib value of all 1,326 two-card discards, built in parallel). Both are built once at startup by `sim` and `play` and injected into strategies.

| Strategy | Description |
|---|---|
| `MaxAvg` | Keeps the 4 cards with the highest average hand score across all cuts |
| `MaxAvgDiff` | Maximizes `avg4 ± avg2`: adds the discards' average crib value when dealer, subtracts it otherwise |
| `MaxMin/MaxMedian/MaxMax/MaxMode` | Like MaxAvg but optimize other summary stats |
| `MinValue/MaxValue/Random` | Naive baselines |

### Pegging Strategies (`strategy/peg/`)

| Strategy | Description |
|---|---|
| `MaxNext` | Plays the card that scores the most points immediately |
| `MaxSetup` | Takes points if available (like MaxNext), otherwise 1-ply lookahead for the best score on its own next play |
| `MinValue/MaxValue/Random` | Naive baselines |

### Simulation (`sim/main.go`)

Runs every combination of active discard × peg strategies against each other. Only the upper triangle of the matchup matrix is simulated (mirror cells reuse swapped results), one goroutine per matchup, each with its own seeded rng. Prints a per-matchup table (wins, pegged/hand/crib/total/diff) and a ranked summary.

When adding or comparing strategies, comment out weak ones in `allCombinations` rather than deleting them. Active strategies are the uncommented entries in the `discards` and `pegs` slices; currently MaxAvg and MaxAvgDiff discards × MaxNext and MaxSetup pegging.

## Human-vs-AI Game (`play/`)

`play/main.go` is an HTTP server with one global in-memory game (`/api/new`, `/api/discard`, `/api/peg`, `/api/next`); the AI is `MaxAvgDiff`/`MaxSetup`. It serves `/play/` from `play/` and everything else (card sprites, favicon, and the analyzer pages) from `site/`, so it must be run from the repo root. The frontend (`play/play.js`) uses root-absolute paths like `/web/img/cards.png`, so it only works behind this server, which is why it isn't deployed.

## Hand Analyzer (`site/`)

- Pages: `site/index.html` (landing), `site/four/`, `site/six/`, all sharing `site/web/css/scoreFour.css` and the `site/web/img/cards.png` sprite sheet.
- `site/web/js/cards.js`: card model and all scoring/statistics (the JS mirror of `game/`).
- `site/web/js/picker.js`: `setupCardPicker({ handSize, onComplete })`, the shared 52-card picker that syncs the selection with the URL hash.
- `site/web/js/scoreFour.js` / `scoreSix.js`: page-specific rendering.
- Plain scripts, no modules or build step; asset paths are relative (`../web/...`) so the site works under the `/crib/` prefix.

## Other Packages

- **`game/`**: cards, deck (`NewDeck`, `RemainingDeck`), hand scoring (`CountCards`; `cut` may be nil), pegging scoring, combination helpers (`ChooseTwo`, `ChooseFour`, `ChooseFourWithRemaining`).
- **`game/counts/`**: `FourSummary` (avg, min/median/max, mode, modeP, below/above avg, std dev, cuts per score) via `MakeSummaries*`, and `MakeSixHands` for the 15 discard options.
- **`game/math/`**: binomial helpers plus combination indexing (vendored from gonum).
- **`util/`**, and the JSON/binary file I/O and binary marshaling in `game/counts/`, exist only for the legacy tools.

## Legacy Code (`legacy/`)

Kept for reference, not actively developed:
- `legacy/java/`: the original Java implementation of the scoring logic (not built).
- `legacy/tools/`: the old pre-compute pipeline (`gen` wrote summaries to a gitignored `scores/` that the site used to fetch; `dump`/`read`/`id` inspected them) plus `score`, still handy for scoring a hand from the command line. They are Go packages in the module, so they must keep compiling.
