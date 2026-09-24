// Card picker shared by the four and six pages. Builds the 52-card grid and the
// selected-card slots, mirrors the selection in the URL hash (e.g. #2d5c5d5h), and
// once handSize cards are picked, hides the grid and calls onComplete(hand) with the
// picked cards as Card objects sorted by id.
function setupCardPicker({ handSize, onComplete }) {
    // one CSS class per card id, positioned on the cards.png sprite sheet
    var styleInnerHTML = ""
    for (let i = 0; i < NumCards; i++) {
        const face = Math.floor(i / 4)
        const suit = i % 4
        const x = face * (-148) + (face / 3)
        const y = suit * (-230)
        styleInnerHTML += `.card${i} {
        background: url('../web/img/cards.png') ${x}px ${y}px;
    }`
    }
    const style = document.createElement('style')
    style.innerHTML = styleInnerHTML
    document.getElementsByTagName('head')[0].appendChild(style)

    // selectable cards: one row per suit, faces ace to king
    const selectableCards = document.getElementById("selectableCards")
    var selectableCardsHTML = ""
    for (let suit = 0; suit < 4; suit++) {
        for (let face = 0; face < 13; face++) {
            const i = face * 4 + suit
            selectableCardsHTML += `<div class="card selectable card${i}"></div>`
        }
        selectableCardsHTML += `<br>`
    }
    selectableCards.innerHTML = selectableCardsHTML

    // selected cards: empty slots filled in as cards are picked
    const selectedCards = document.getElementById("selectedCards")
    selectedCards.innerHTML = `<div class="card selected"></div>`.repeat(handSize)

    function cardIdOf(el) {
        for (const cl of el.classList) {
            if (cl.startsWith("card") && cl.length > 4) {
                return Number(cl.substring(4))
            }
        }
    }

    function drawSelectedCards() {
        const picked = selectableCards.getElementsByClassName("selected")
        const slots = selectedCards.getElementsByClassName("selected")
        const ids = Array.from(picked, cardIdOf).sort((a, b) => a - b)

        for (let i = 0; i < slots.length; i++) {
            slots[i].className = "card selected"
            if (i < ids.length) {
                slots[i].classList.add(`card${ids[i]}`)
            }
        }

        history.replaceState(null, "", ids.reduce((acc, id) => acc + indexToString[id], "#"))

        if (ids.length == handSize) {
            selectableCards.style.display = "none"
            onComplete(ids.map(id => Card(id)))
        }
    }

    for (const card of selectableCards.getElementsByClassName("selectable")) {
        card.addEventListener("click", () => {
            card.classList.toggle("selected")
            drawSelectedCards()
        })
    }

    // restore a selection from the URL hash, ignoring unknown codes and extra cards
    const hash = window.location.hash.substring(1)
    let restored = 0
    for (let i = 0; i + 2 <= hash.length && restored < handSize; i += 2) {
        const id = stringToIndex[hash.substring(i, i + 2)]
        if (id !== undefined) {
            selectableCards.getElementsByClassName(`card${id}`)[0].classList.add("selected")
            restored++
        }
    }
    drawSelectedCards()
}
