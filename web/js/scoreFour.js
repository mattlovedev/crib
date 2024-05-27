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
for (let i = 0; i < 4; i++) {
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

async function displaySummary() {
    //const selectedSelectableCards = document.getElementById("selectableCards").getElementsByClassName("selected")
    //const cards = [].slice.call(selectedSelectableCards)
    const selectedCards = document.getElementById("selectedCards").getElementsByClassName("selected")
    const cards = [].slice.call(selectedCards)
    const cardStrIds = cards.map(card => {
        for (let cl of card.classList) {
            if (cl.startsWith("card") && cl.length > 4) {
                return indexToString[parseInt(cl.substring(4))]
            }
        }
    }).reduce((acc, cur) => acc + cur, "")

    const hash = await window.crypto.subtle.digest("SHA-1", new TextEncoder().encode(cardStrIds))
    const intPrefix = new DataView(hash).getBigInt64(0, true)
    const prefix = ((intPrefix % 49n) + 49n) % 49n

    fetch(`../scores/four/four_summaries_${prefix}.json`)
    .then(res => res.text())
    .then(text => {
        const c = JSON.parse(text)
        const hand = c[cardStrIds]

        setSummaryStats(hand)
        setSummaryCounts(hand.Counts)

        document.getElementById("selectableCards").style.display = "none"
        document.getElementById("summary").style.display = "block"
    })
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

    if (selectedSelectableCards.length == 4) {
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
