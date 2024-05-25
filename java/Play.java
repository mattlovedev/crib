import java.util.ArrayList;
import java.util.Collections;
import java.util.Comparator;
import java.util.HashMap;
import java.util.Map;

class Play {

    private final Hand hand;
    private final Deck deck;
    private final boolean verbose;

    Play(Hand hand, Deck deck, boolean verbose) {
        this.hand = hand;
        this.deck = deck.copy();
        this.verbose = verbose;
        print();
    }

    Play(Hand hand, Deck deck) {
        this(hand, deck, false);
    }

    private void print() {
        //System.out.println(hand);
        System.out.print(hand + " ");
        if (verbose) System.out.println();
        ArrayList<Hand> hands = new ArrayList<>();
        for (Card cut : deck.availableCards()) {
            Hand h = new Hand(hand, cut);
            hands.add(h);
        }
        Collections.sort(hands);

        HashMap<Integer, Integer> map = new HashMap<>();
        HashMap<Integer, ArrayList<Card>> cuts = new HashMap<>();

        int min = 29, max = 0, mode;
        float median, mean = 0.0f;

        for (Hand h : hands) {
            if (h.getCount() < min) min = h.getCount();
            if (h.getCount() > max) max = h.getCount();
            mean += h.getCount();

            if (!map.containsKey(h.getCount()))
                map.put(h.getCount(), 0);
            map.put(h.getCount(), map.get(h.getCount()) + 1);

            if (!cuts.containsKey(h.getCount()))
                cuts.put(h.getCount(), new ArrayList<>());
            cuts.get(h.getCount()).add(h.getCut());
        }
        median = min + (float) ((max - min) / 2.0);
        mean /= hands.size();

        ArrayList<Map.Entry<Integer,Integer>> entryListCount = new ArrayList<>(map.entrySet());
        ArrayList<Map.Entry<Integer,Integer>> entryListOccur = new ArrayList<>(map.entrySet());

        Collections.sort(entryListCount, new Comparator<Map.Entry<Integer,Integer>>() {
            public int compare(Map.Entry<Integer,Integer> a, Map.Entry<Integer,Integer> b) {
                return a.getKey() - b.getKey();
            }
        });
        Collections.sort(entryListOccur, new Comparator<Map.Entry<Integer,Integer>>() {
            public int compare(Map.Entry<Integer,Integer> a, Map.Entry<Integer,Integer> b) {
                return b.getValue() - a.getValue();
            }
        });

        mode = entryListOccur.get(0).getKey();

        int belowMedian = 0, aboveMedian = 0, belowMean = 0, aboveMean = 0;

        for (Map.Entry<Integer,Integer> entry : entryListCount) {
            if (entry.getKey() < median) belowMedian += entry.getValue();
            else if (entry.getKey() > median) aboveMedian += entry.getValue();
            if (entry.getKey() < mean) belowMean += entry.getValue();
            else if (entry.getKey() > mean) aboveMean += entry.getValue();

            if (verbose) {
                System.out.printf("%2d: %4.1f%% %2d - ", entry.getKey(), entry.getValue() * 100.0 / hands.size(), entry.getValue());
                for (Card card : cuts.get(entry.getKey())) {
                    System.out.printf("%s ", card);
                }
                System.out.println();
            }
        }

        System.out.printf("min: %2d (+%2d) (%4.1f%%) max: %2d (+%2d) (%4.1f%%) mode: %2d (+%2d) (%4.1f%%) median: %4.1f (%4.1f%%/%4.1f%%) mean: %4.1f (%4.1f%%/%4.1f%%)\n",
            min, min - hand.getCount(), map.get(min) * 100.0 / hands.size(), max, max - hand.getCount(), map.get(max) * 100.0 / hands.size(),
            mode, mode - hand.getCount(), map.get(mode) * 100.0 / hands.size(), median, belowMedian * 100.0 / hands.size(), aboveMedian * 100.0 / hands.size(),
            mean, belowMean * 100.0 / hands.size(), aboveMean * 100.0 / hands.size());
    }

}
