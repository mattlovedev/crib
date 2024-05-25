import java.util.ArrayList;
import java.util.Collections;

class Deck {

    private final ArrayList<Card> cards;

    Deck() {
        cards = new ArrayList<>();
        for (int i = 0; i < 52; i++) {
            cards.add(Card.getInstance(i));
        }
    }

    /* used only for copying a deck that may already be missing cards */
    private Deck(ArrayList<Card> cards) {
        this.cards = new ArrayList<>(cards);
    }

    void shuffle() {
        Collections.shuffle(cards);
    }

    Card removeCard(String id) {
        Card card = Card.getInstance(id);
        cards.remove(card);
        return card;
    }

    void removeCard(Card card) {
        cards.remove(card);
    }

    ArrayList<Card> availableCards() {
        return new ArrayList<Card>(cards);
    }

    Deck copy() {
        return new Deck(cards);
    }

}
