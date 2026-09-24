function setSummaryStats(stats) {
    const mappings = [
        { field: "average", value: "Avg" },
        { field: "belowAverage", value: "BelowAvg" },
        { field: "aboveAverage", value: "AboveAvg" },
        { field: "stdDev", value: "StdDev" },
        { field: "mode", value: "Mode" },
        { field: "modeP", value: "ModeP" },
        { field: "min", value: "Min" },
        { field: "median", value: "Median" },
        { field: "max", value: "Max" }
    ]

    mappings.forEach(({ field, value }) => {
        document.getElementById(field).innerHTML = stats[value]
    })
}

function setSummaryCounts(counts) {
    var countsHTML = ""
    for (const [count, values] of Object.entries(counts)) {
        countsHTML += `<div class="countsRow">`
        countsHTML += `<div class="countsHeader">${count} (${values.length}):</div>`
        values.forEach(card => {
            countsHTML += `<div class="card card${card}"></div>`
        })
        countsHTML += `</div>` // countsRow
    }
    document.getElementById("counts").innerHTML = countsHTML
}

setupCardPicker({
    handSize: 4,
    onComplete: hand => {
        const summaries = makeSummaries(hand)
        setSummaryStats(summaries)
        setSummaryCounts(summaries.Counts)
        document.getElementById("summary").style.display = "block"
    }
})
