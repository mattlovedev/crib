'use strict'

// Card sprite: face = id/4, suit = id%4
// x = face * -148 + face/3, y = suit * -230  (from scoreFour.js)
// We render at 50% scale: 74x115 px display
const SCALE = 0.5
const CARD_W = 148, CARD_H = 230

function cardBgStyle(id) {
  const face = Math.floor(id / 4)
  const suit = id % 4
  const x = (face * -148 + face / 3) * SCALE
  const y = suit * -230 * SCALE
  const bw = Math.round(148 * 13 * SCALE) + 'px'
  return `background:url('/web/img/cards.png') ${x}px ${y}px; background-size:${bw} auto;`
}

function cardEl(id, classes = []) {
  const div = document.createElement('div')
  div.className = ['card', ...classes].join(' ')
  if (id >= 0) {
    div.style.cssText = cardBgStyle(id)
    div.dataset.id = id
  } else {
    div.classList.add('face-down')
  }
  return div
}

function faceDown() {
  const div = document.createElement('div')
  div.className = 'card face-down'
  return div
}

// State
let state = null
let selectedDiscard = new Set()

// DOM refs
const elScoreAI = document.getElementById('score-ai')
const elScoreHuman = document.getElementById('score-human')
const elDealerAI = document.getElementById('dealer-ai')
const elDealerHuman = document.getElementById('dealer-human')
const elBarAI = document.getElementById('score-bar-ai')
const elBarHuman = document.getElementById('score-bar-human')

const elAIHand = document.getElementById('ai-hand')
const elAILabel = document.getElementById('ai-label')
const elHumanHand = document.getElementById('human-hand')
const elHumanLabel = document.getElementById('human-label')
const elCutCard = document.getElementById('cut-card')
const elCribArea = document.getElementById('crib-area')
const elCribLabel = document.getElementById('crib-label')
const elCribCards = document.getElementById('crib-cards')
const elPegArea = document.getElementById('peg-area')
const elPegCount = document.getElementById('peg-count')
const elPegSeries = document.getElementById('peg-series')

const btnDiscard = document.getElementById('btn-discard')
const btnGo = document.getElementById('btn-go')
const btnNext = document.getElementById('btn-next')
const btnNew = document.getElementById('btn-new')

const elLogEntries = document.getElementById('log-entries')

const elOverlay = document.getElementById('overlay')
const elOverlayTitle = document.getElementById('overlay-title')
const elRoundScores = document.getElementById('round-scores')
const btnOverlayNext = document.getElementById('btn-overlay-next')
const btnOverlayNew = document.getElementById('btn-overlay-new')

function cardName(id) {
  const faces = ['A','2','3','4','5','6','7','8','9','T','J','Q','K']
  const suits = ['♣','♦','♥','♠']
  return faces[Math.floor(id/4)] + suits[id%4]
}

// Build CSS classes for card sprite once
const styleEl = document.createElement('style')
document.head.appendChild(styleEl)

async function api(path, body) {
  const res = await fetch(path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  })
  return res.json()
}

function renderScores(s) {
  elScoreAI.textContent = s.scores[1]
  elScoreHuman.textContent = s.scores[0]
  elBarAI.style.width = Math.min(100, s.scores[1] / 121 * 100) + '%'
  elBarHuman.style.width = Math.min(100, s.scores[0] / 121 * 100) + '%'

  const humanDealer = s.dealer === 0
  elDealerHuman.classList.toggle('visible', humanDealer)
  elDealerAI.classList.toggle('visible', !humanDealer)
}

function renderLog(s) {
  elLogEntries.innerHTML = ''
  const entries = s.log || []
  for (let i = entries.length - 1; i >= 0; i--) {
    const div = document.createElement('div')
    div.className = 'log-entry'
    div.textContent = entries[i]
    elLogEntries.appendChild(div)
  }
}

function renderDiscard(s) {
  elAILabel.textContent = "AI's Hand"
  elAIHand.innerHTML = ''
  for (let i = 0; i < 6; i++) elAIHand.appendChild(faceDown())

  elHumanLabel.textContent = 'Your Hand — pick 2 to discard'
  elHumanHand.innerHTML = ''
  selectedDiscard.clear()

  s.human_hand.forEach(id => {
    const el = cardEl(id, ['selectable'])
    el.addEventListener('click', () => toggleDiscard(el, id))
    elHumanHand.appendChild(el)
  })

  elCutCard.innerHTML = ''
  elCribArea.style.display = 'none'
  elPegArea.style.display = 'none'

  btnDiscard.style.display = 'inline-block'
  btnDiscard.disabled = true
  btnDiscard.textContent = 'Discard Selected (0/2)'
  btnGo.style.display = 'none'
  btnNext.style.display = 'none'
}

function toggleDiscard(el, id) {
  if (selectedDiscard.has(id)) {
    selectedDiscard.delete(id)
    el.classList.remove('selected')
  } else if (selectedDiscard.size < 2) {
    selectedDiscard.add(id)
    el.classList.add('selected')
  }
  const n = selectedDiscard.size
  btnDiscard.disabled = n !== 2
  btnDiscard.textContent = `Discard Selected (${n}/2)`
}

function renderPeg(s) {
  elAILabel.textContent = `AI's Hand (${s.ai_hand_count} cards)`
  elAIHand.innerHTML = ''
  for (let i = 0; i < s.ai_hand_count; i++) elAIHand.appendChild(faceDown())
  // show AI played cards next to face-down
  if (s.ai_played && s.ai_played.length > 0) {
    const sep = document.createElement('span')
    sep.style.cssText = 'opacity:0.4;align-self:center;margin:0 4px;font-size:12px;'
    sep.textContent = 'played:'
    elAIHand.appendChild(sep)
    s.ai_played.forEach(id => elAIHand.appendChild(cardEl(id)))
  }

  elHumanLabel.textContent = 'Your Hand — click to play'
  elHumanHand.innerHTML = ''

  const count = s.peg_count
  s.human_hand.forEach(id => {
    const val = Math.floor(id / 4) + 1 > 10 ? 10 : Math.floor(id / 4) + 1
    const legal = count + val <= 31
    const classes = legal ? ['playable'] : ['illegal']
    const el = cardEl(id, classes)
    if (legal) {
      el.addEventListener('click', () => playCard(id))
    }
    elHumanHand.appendChild(el)
  })

  if (s.cut >= 0) {
    elCutCard.innerHTML = ''
    elCutCard.appendChild(cardEl(s.cut))
  }

  elPegArea.style.display = 'block'
  elPegCount.textContent = s.peg_count

  elPegSeries.innerHTML = ''
  ;(s.peg_series || []).forEach(id => elPegSeries.appendChild(cardEl(id)))

  elCribArea.style.display = 'none'

  btnDiscard.style.display = 'none'
  btnNext.style.display = 'none'

  // Show go button only if human can't play
  const canPlay = s.can_play
  btnGo.style.display = (!canPlay && s.whose_turn === 'human') ? 'inline-block' : 'none'
}

function renderScore(s) {
  elAILabel.textContent = "AI's Hand"
  elAIHand.innerHTML = ''
  s.ai_hand.forEach(id => elAIHand.appendChild(cardEl(id)))

  elHumanLabel.textContent = 'Your Hand'
  elHumanHand.innerHTML = ''
  s.human_hand.forEach(id => elHumanHand.appendChild(cardEl(id)))

  if (s.cut >= 0) {
    elCutCard.innerHTML = ''
    elCutCard.appendChild(cardEl(s.cut))
  }

  elCribArea.style.display = 'block'
  elCribLabel.textContent = `Crib (${s.dealer === 0 ? 'Yours' : "AI's"})`
  elCribCards.innerHTML = ''
  s.crib.forEach(id => elCribCards.appendChild(cardEl(id)))

  elPegArea.style.display = 'none'

  btnDiscard.style.display = 'none'
  btnGo.style.display = 'none'
  btnNext.style.display = 'inline-block'

  // Show overlay with round summary
  elOverlayTitle.textContent = 'Round Over'
  const dealer = s.dealer === 0 ? 'Your' : "AI's"
  elRoundScores.innerHTML = `
    Your hand: <strong>${s.human_hand_score}</strong> pts<br>
    AI hand: <strong>${s.ai_hand_score}</strong> pts<br>
    ${dealer} crib: <strong>${s.crib_score}</strong> pts<br>
    <br>Score — You: <strong>${s.scores[0]}</strong> | AI: <strong>${s.scores[1]}</strong>
  `
  btnOverlayNext.style.display = 'inline-block'
  btnOverlayNew.style.display = 'inline-block'
  elOverlay.style.display = 'flex'
}

function renderGameOver(s) {
  elAIHand.innerHTML = ''
  s.ai_hand.forEach(id => elAIHand.appendChild(cardEl(id)))
  elHumanHand.innerHTML = ''
  s.human_hand.forEach(id => elHumanHand.appendChild(cardEl(id)))

  elCribArea.style.display = 'block'
  elCribLabel.textContent = `Crib (${s.dealer === 0 ? 'Yours' : "AI's"})`
  elCribCards.innerHTML = ''
  s.crib.forEach(id => elCribCards.appendChild(cardEl(id)))

  elPegArea.style.display = 'none'
  btnDiscard.style.display = 'none'
  btnGo.style.display = 'none'
  btnNext.style.display = 'none'

  elOverlayTitle.textContent = s.winner === 'human' ? '🎉 You Win!' : 'AI Wins'
  elRoundScores.innerHTML = `Final — You: <strong>${s.scores[0]}</strong> | AI: <strong>${s.scores[1]}</strong>`
  btnOverlayNext.style.display = 'none'
  btnOverlayNew.style.display = 'inline-block'
  elOverlay.style.display = 'flex'
}

function sortedHand(ids) {
  return [...ids].sort((a, b) => Math.floor(a / 4) - Math.floor(b / 4))
}

function render(s) {
  s = { ...s, human_hand: sortedHand(s.human_hand) }
  state = s
  if (s.error) {
    console.error('API error:', s.error)
    return
  }
  renderScores(s)
  renderLog(s)

  switch (s.phase) {
    case 'discard':   renderDiscard(s); break
    case 'peg':       renderPeg(s); break
    case 'score':     renderScore(s); break
    case 'game_over': renderGameOver(s); break
  }
}

async function newGame() {
  elOverlay.style.display = 'none'
  selectedDiscard.clear()
  render(await api('/api/new', {}))
}

async function doDiscard() {
  if (selectedDiscard.size !== 2) return
  const cards = [...selectedDiscard]
  render(await api('/api/discard', { cards }))
}

async function playCard(id) {
  render(await api('/api/peg', { card: id }))
}

async function sayGo() {
  render(await api('/api/peg', { card: -1 }))
}

async function nextRound() {
  elOverlay.style.display = 'none'
  render(await api('/api/next', {}))
}

btnDiscard.addEventListener('click', doDiscard)
btnGo.addEventListener('click', sayGo)
btnNext.addEventListener('click', nextRound)
btnNew.addEventListener('click', newGame)
btnOverlayNext.addEventListener('click', nextRound)
btnOverlayNew.addEventListener('click', newGame)

// Start
newGame()
