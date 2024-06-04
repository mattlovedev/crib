// setup classes
var styleInnerHTML = ""
for (let i = 0; i < NumCards; i++) {
    const face = Math.floor(i / 4)
    const suit = i % 4
    //const x = (i % 13) * (-148) + (i % 13 / 3)
    //const y = Math.floor(i / 13) * (-230)
    const x = face * (-148) + (face / 3)
    const y = suit * (-230)
    styleInnerHTML += `.card${i} {
        background: url('../web/img/cards.png') ${x}px ${y}px;
    }`
}
var style = document.createElement('style')
style.innerHTML = styleInnerHTML
document.getElementsByTagName('head')[0].appendChild(style)

// setup selectable cards
var selectableCardsHTML = ""
//for (let i = 0; i < NumCards; i++) {
for (let suit = 0; suit < 4; suit++) {
    for (let face = 0; face < 13; face++) {
        const i = face * 4 + suit
        selectableCardsHTML += `<div class="card selectable card${i}"></div>`
    }
    selectableCardsHTML += `<br>`
}
document.getElementById("selectableCards").innerHTML = selectableCardsHTML

// setup selected cards
var selectedCardsHTML = ""
for (let i = 0; i < 6; i++) {
    selectedCardsHTML += `<div class="card selected"></div>`
}
document.getElementById("selectedCards").innerHTML = selectedCardsHTML

// util and game logic

function classListToCard(cl) {
    var selected = false
    var selectable = false
    var cardId = ""
    for (c of cl) {
        if (c == "selected") { selected = true }
        else if (c == "selectable") { selectable = true }
        else if (c.startsWith("card") && c.length > 4) { cardId = c }
    }
    return { selected, selectable, cardId }
}


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
            countsHTML += `<div class="card card${stringToIndex[card]}"></div>`
        })
        countsHTML += `</div>` // countsRow
    }
    document.getElementById("counts").innerHTML = countsHTML
}

function setHands(hands) {
    var handsHTML = ""
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

async function displaySummary() {
    const selectedCards = document.getElementById("selectedCards").getElementsByClassName("selected")
    const cards = [].slice.call(selectedCards)

    const cardIds = cards.map(card => {
        for (let cl of card.classList) {
            if (cl.startsWith("card") && cl.length > 4) {
                return parseInt(cl.substring(4))
            }
        }
    })

    const hand = cardIds.map(id => Card(id))

    const hands = makeSixHands(hand)

    console.log(hands)
    setHands(hands)

    document.getElementById("selectableCards").style.display = "none"
    document.getElementById("hands").style.display = "block"
}

function drawSelectedCards() {
    const selectedSelectableCards = document.getElementById("selectableCards").getElementsByClassName("selected")
    const selectedCards = document.getElementById("selectedCards").getElementsByClassName("selected")

    // first reset current selected cards
    for (let i = 0; i < selectedCards.length; i++) {
        selectedCards[i].className = "card selected"
    }

    // then fill in with latest selection

    const selectedCardIds = []

    for (let i = 0; i < selectedSelectableCards.length; i++) {
        const { cardId } = classListToCard(selectedSelectableCards[i].classList)
        selectedCardIds.push(Number(cardId.substring(4)))
        //selectedCards[i].classList.add(cardId)
    }

    selectedCardIds.sort((a,b) => a-b)

    for (let i = 0; i < selectedCardIds.length; i++) {
        selectedCards[i].classList.add(`card${selectedCardIds[i]}`)
    }

    if (selectedCardIds.length > 0) {
        history.replaceState(null, "", selectedCardIds.reduce((acc, curr) => acc + indexToString[curr] , "#"))
    } else {
        history.replaceState(null, "", "#")
    }

    if (selectedSelectableCards.length == 6) {
        displaySummary()
    }
}

function selectCard(e) {
    const { selected, selectable } = classListToCard(e.target.classList)
    if (!selectable) {
        return
    }
    if (!selected) {
        e.target.classList.add("selected")
    }
    if (selected) {
        e.target.classList.remove("selected")
    }
    drawSelectedCards()
}

const cards = document.getElementById("selection").getElementsByClassName("card")
Array.from(cards).forEach(card => {
    card.addEventListener("click", selectCard)
})

const hash = window.location.hash.substring(1)
for (let i = 0; i < hash.length; i+=2) {
    document.getElementsByClassName(`card${stringToIndex[hash.substring(i, i+2)]}`)[0].classList.add("selected")
}
drawSelectedCards()