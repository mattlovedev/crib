// Compares the hand statistics computed by web/js/cards.js against Go's, read from stdin.
// Exits non-zero if any statistic differs. See main.go for the input format.
//
//   go run ./scripts/parity | node scripts/parity/check.js

const fs = require("fs")
const path = require("path")

// cards.js is a plain browser script; load its top-level declarations into this scope.
eval(fs.readFileSync(path.join(__dirname, "../../web/js/cards.js"), "utf8").replace(/^const /gm, "var "))

const statNames = ["Avg", "Min", "Median", "Max", "Mode", "ModeP", "BelowAvg", "AboveAvg", "StdDev"]
const codes = ids => ids.map(id => indexToString[id]).join("")

let checked = 0
let mismatches = 0
const report = (what, jsStats, goStats) => {
    checked++
    statNames.forEach((name, i) => {
        if (Number(jsStats[name]) !== Number(goStats[i])) {
            mismatches++
            if (mismatches <= 20) {
                console.log(`MISMATCH ${what} ${name}: js=${jsStats[name]} go=${goStats[i]}`)
            }
        }
    })
}

const sixCache = {}
for (const line of fs.readFileSync(0, "utf8").trim().split("\n")) {
    const f = line.split(" ")
    if (f[0] === "4") {
        const hand = f.slice(1, 5).map(Number)
        report(codes(hand), makeSummaries(hand.map(Card)), f.slice(5))
    } else if (f[0] === "6") {
        // rows may be ordered differently when averages tie, so match them by kept hand
        const six = f.slice(1, 7).map(Number)
        const kept = f.slice(7, 11).map(Number)
        const key = six.join(",")
        sixCache[key] ??= Object.fromEntries(makeSixHands(six.map(Card)).map(r => [r.Hand.join(","), r]))
        const row = sixCache[key][kept.join(",")]
        const what = `${codes(six)} keep ${codes(kept)}`
        if (!row || row.Crib.join(",") !== f.slice(11, 13).join(",")) {
            checked++
            mismatches++
            console.log(`MISMATCH ${what}: no matching JS row`)
            continue
        }
        report(what, row.Summary, f.slice(13))
    }
}

console.log(`checked ${checked} summaries (${statNames.length} stats each): ${mismatches} mismatches`)
process.exit(mismatches === 0 ? 0 : 1)
