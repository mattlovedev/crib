import java.util.ArrayList;
import java.util.Collections;

class Hand implements Comparable<Hand> {

    private final ArrayList<Card> cards;
    private int count = 0;
    private final Card cut;

    Hand(ArrayList<Card> cards) {
        this.cards = cards;
        Collections.sort(cards);
        count = countCards(false);
        cut = null;
    }

    Hand(Card c1, Card c2, Card c3, Card c4) {
        cards = new ArrayList<>();
        cards.add(c1);
        cards.add(c2);
        cards.add(c3);
        cards.add(c4);
        Collections.sort(cards);
        count = countCards(false);
        cut = null;
    }

    Hand(Hand hand, Card cut) {
        this(hand, cut, false);
    }

    Hand(Hand hand, Card cut, boolean crib) {
        cards = new ArrayList<>(hand.cards);
        Card nobsJack = Card.getInstance(Card.JACK, cut.getSuit());
        if (cards.contains(nobsJack))
            count++;
        cards.add(cut);
        Collections.sort(cards);
        count += countCards(crib);
        this.cut = cut;
    }

    int getCount() {
        return count;
    }

    Card getCut() {
        return cut;
    }

    private int countCards(boolean crib) {
        int sum = 0;
        sum += countFifteens();
        sum += countPairs();
        sum += countRuns();
        if (crib && countFlush() == 5)
            sum += 5;
        else
            // TODO: bug - 4 hole cards flush of four without needing cut
            sum += countFlush();
        return sum;
    }

    private static final int[] masks = { 1, 2, 4, 8, 16, 32 };

    private int countFifteens() {
        int fifteens = 0;
        for (int i = 1; i < masks[cards.size()]; i++) { // 1 bit for every card position 0001 to 1111
            int sum = 0;
            for (int j = 0; j < cards.size(); j++) {
                if ((masks[j] & i) > 0) {
                    sum += cards.get(j).getValue();
                }
            }
            if (sum == 15) {
                fifteens += 2;
            }
        }
        return fifteens;
    }

    private int countPairs() {
        int pairs = 0;
        for (int i = 0; i < cards.size() - 1; i++)
            for (int j = i + 1; j < cards.size(); j++)
                if (cards.get(i).sameFace(cards.get(j)))
                    pairs += 2;
        return pairs;
    }

    private int countRuns() {
        ArrayList<Card> uniques = new ArrayList<>(cards), duplicates = new ArrayList<>();
        removeDuplicates(uniques, duplicates);
        
        for (int len = uniques.size(); len > 2; len--) {
            for (int i = 0; i <= uniques.size() - len; i++) {
                if (isStraight(uniques, i, len)) {
                    return len * (1 + isInStraight(uniques, i, len, duplicates));
                }
            }
        }
        return 0;
    }

    private static void removeDuplicates(ArrayList<Card> uniques, ArrayList<Card> duplicates) {
        for (int i = 0; i < uniques.size() - 1; i++) {
            while (i < uniques.size() - 1 && uniques.get(i).sameFace(uniques.get(i+1))) {
                duplicates.add(uniques.remove(i+1));
            }
        }
    }

    private static boolean isStraight(ArrayList<Card> cards, int start, int len) {
        for (int i = start; i < start + len - 1; i++) {
            if (cards.get(i+1).getFace() - cards.get(i).getFace() != 1) {
                return false;
            }
        }
        return true;
    }

    private static int isInStraight(ArrayList<Card> cards, int start, int len, ArrayList<Card> duplicates) {
        int count = 0;
        int[] dupes = new int[13];
        for (int i = 0; i < duplicates.size(); i++) {
            for (int j = start; j < start + len; j++) {
                if (duplicates.get(i).sameFace(cards.get(j))) {
                    dupes[duplicates.get(i).getFace()]++;
                    //count++;
                }
            }
        }
        boolean one_match = false;
        for (int i = 0; i < dupes.length; i++) {
            if (dupes[i] == 1) {
                if (one_match) {
                    count += 2;
                } else {
                    count++;
                    one_match = true;
                }
            } else if (dupes[i] == 2) {
                count += 2;
            }
        }
        return count;
    }

    private int countFlush() {
        for (int i = 1; i < cards.size(); i++)
            if (!cards.get(0).sameSuit(cards.get(i)))
                return 0;
        return cards.size();
    }

    @Override
    public String toString() {
        if (cards.size() == 5)
            return String.format("%s %s %s %s %s (%2d)", cards.get(0), cards.get(1), cards.get(2), cards.get(3), cards.get(4), count);
        return String.format("%s %s %s %s (%2d)", cards.get(0), cards.get(1), cards.get(2), cards.get(3), count);
    }

    @Override
    public int compareTo(Hand other) {
        //return count - other.count;
        return other.count - count;
    }

}
