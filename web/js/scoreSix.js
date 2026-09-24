function setHands(hands) {
    var handsHTML = "<hr>"
    hands.forEach(hand => {
        handsHTML += `<div class="handRow">`
        handsHTML += `<div class="cards">`
        hand.Hand.forEach(card => {
            handsHTML += `<div class="card card${card}"></div>`
        })
        handsHTML += `</div>` // cards
        handsHTML += `<div class="stats">`
        handsHTML += `<div>`
        handsHTML += `<h3>Avg:` + hand.Summary.Avg + `</h3>`
        handsHTML += `<h3>Below Avg:` + hand.Summary.BelowAvg + `</h3>`
        handsHTML += `<h3>Above Avg:` + hand.Summary.AboveAvg + `</h3>`
        handsHTML += `</div>`
        handsHTML += `<div>`
        handsHTML += `<h3>Std Dev:` + hand.Summary.StdDev + `</h3>`
        handsHTML += `<h3>Mode:` + hand.Summary.Mode + `</h3>`
        handsHTML += `<h3>Mode Pct:` + hand.Summary.ModeP + `</h3>`
        handsHTML += `</div>`
        handsHTML += `<div>`
        handsHTML += `<h3>Min:` + hand.Summary.Min + `</h3>`
        handsHTML += `<h3>Median:` + hand.Summary.Median + `</h3>`
        handsHTML += `<h3>Max:` + hand.Summary.Max + `</h3>`
        handsHTML += `</div>`
        handsHTML += `</div>` // stats
        handsHTML += `</div>` // handRow
        handsHTML += `<hr>`
    })
    document.getElementById("hands").innerHTML = handsHTML
}

setupCardPicker({
    handSize: 6,
    onComplete: hand => {
        setHands(makeSixHands(hand))
        document.getElementById("hands").style.display = "block"
    }
})
