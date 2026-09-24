# crib

Cribbage tools: a hand analyzer that runs in the browser, live at
[mattlove.dev/crib](https://mattlove.dev/crib/), and a simulator for comparing
AI strategies, which also powers a local human-vs-AI game.

Cards are written as two characters, face then suit: faces `a23456789tjqk`,
suits `cdhs`. For example `5h` is the five of hearts and `tc` is the ten of clubs.

## Hand analyzer (`site/`)

- **Score Four:** pick 4 cards to see how the hand scores across every possible cut:
  average, min/median/max, mode, standard deviation, and which cuts give each score.
- **Score Six:** pick the 6 cards you were dealt to see all 15 ways to keep 4 and
  discard 2, ranked by average score.

The selection is kept in the URL, so hands can be linked directly, e.g.
[`/four/#5c5d5hjs`](https://mattlove.dev/crib/four/#5c5d5hjs).

It's plain static HTML/JS with no build step. To run it locally:

```bash
python3 -m http.server -d site 8000   # then open http://localhost:8000
```

Pushes to `main` that touch the site are deployed to GitHub Pages by
[`.github/workflows/pages.yml`](.github/workflows/pages.yml), after `go vet` and the
parity check below pass.

## AI strategy simulator (`sim/`)

```bash
go run ./sim          # 1000 games per matchup
go run ./sim 5000
```

Plays every combination of the active discard and pegging strategies against each
other and prints win rates and points per game (pegged, hand, crib, differential),
then a ranked summary. Strategies live in `strategy/discard/` and `strategy/peg/`;
the interfaces are in `strategy/strategy.go`. The game rules are in `engine/`.

## Play against the AI (`play/`)

```bash
go run ./play         # then open http://localhost:8080
```

A local web game against the best-performing strategy in the simulator so far (MaxAvgDiff discards,
MaxSetup pegging). It takes about 10 seconds to start while it builds its caches.
It needs the Go server, so it is not part of the deployed site.

## Checks

```bash
go vet ./...
go run ./scripts/parity | node scripts/parity/check.js
```

The site scores hands in JavaScript (`site/web/js/cards.js`), re-implementing the Go
scoring in `game/`. The parity check compares the two on every 4-card hand and a
sample of 6-card hands, and fails if they disagree.

## Layout

| Path | What |
|---|---|
| `site/` | The hand analyzer (deployed) |
| `game/` | Cards, hand and pegging scoring, hand statistics |
| `engine/` | Full-game rules: dealing, pegging, scoring |
| `strategy/` | AI discard and pegging strategies |
| `sim/` | Strategy-vs-strategy simulator |
| `play/` | Local human-vs-AI web game |
| `scripts/parity/` | JS vs Go scoring check |
| `legacy/` | The original Java implementation and old Go tools (including `legacy/tools/score`, a command-line hand scorer) |
