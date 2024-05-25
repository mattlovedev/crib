class Card implements Comparable<Card> {

    private static final int NUM_SUITS = 4;

    private static final int NUM_FACES = 13;

    private static final int NUM_CARDS = NUM_SUITS * NUM_FACES;

    /* used explicity for checking for nobs */
    public static final int JACK = 10;

    private static final Card[] cards = new Card[NUM_CARDS];

    static {
        for (int i = 0; i < NUM_CARDS; i++) {
            cards[i] = new Card(i);
        }
    }

    private static final String faces = "a23456789tjqk";
    private static final String suits = "hcds";

    private int id;
    private int face;
    private int suit;
    private int value;

    private Card(int id) {
        this.id = id;
        face = id / NUM_SUITS;
        suit = id % NUM_SUITS;
        value = face + 1;
        /* all face cards worth 10 */
        if (value > 10) {
            value = 10;
        }
    }

    static Card getInstance(int id) {
        return cards[id];
    }

    static Card getInstance(String id) {
        int face = faces.indexOf(id.charAt(0));
        int suit = suits.indexOf(id.charAt(1));
        /* cards are sorted by face then suit */
        return cards[face * NUM_SUITS + suit];
    }

    static Card getInstance(int face, int suit) {
        /* cards are sorted by face then suit */
        return cards[face * NUM_SUITS + suit];
    }

    boolean sameFace(Card other) {
        return face == other.face;
    }

    boolean sameSuit(Card other) {
        return suit == other.suit;
    }

    int getFace() {
        return face;
    }

    int getSuit() {
        return suit;
    }

    int getValue() {
        return value;
    }

    @Override
    public String toString() {
        return faces.charAt(face) + "" + suits.charAt(suit);
    }

    @Override
    public int compareTo(Card other) {
        return id - other.id;
    }

}
