import java.util.ArrayList;
import java.util.Collections;
import java.util.Comparator;
import java.util.HashMap;
import java.util.Map;

public class Crib {

    public static void main(String[] args) {
        if (args.length != 4 && args.length != 5 && args.length != 6) {
            System.out.println("usage: java Crib <c1> <c2> <c3> <c4> [<cut>|<c5> <c6>]");
            System.exit(1);
        }
        Deck deck = new Deck();
        Card c1 = deck.removeCard(args[0]);
        Card c2 = deck.removeCard(args[1]);
        Card c3 = deck.removeCard(args[2]);
        Card c4 = deck.removeCard(args[3]);

        if (args.length == 5) {
            Card cut = deck.removeCard(args[4]);
            Hand hand = new Hand(c1, c2, c3, c4);
            hand = new Hand(hand, cut);
            System.out.println(hand);
        } else if (args.length == 4) {
            Hand hand = new Hand(c1, c2, c3, c4);
            Play play = new Play(hand, deck, true);
        } else if (args.length == 6) {
            Card c5 = deck.removeCard(args[4]);
            Card c6 = deck.removeCard(args[5]);
            ArrayList<Card> cards = new ArrayList<>();
            cards.add(c1);
            cards.add(c2);
            cards.add(c3);
            cards.add(c4);
            cards.add(c5);
            cards.add(c6);
            //ArrayList<Play> plays = new ArrayList<>();
            for (int i = 0; i < cards.size() - 1; i++) {
                for (int j = i + 1; j < cards.size(); j++) {
                    ArrayList<Card> newcards = new ArrayList<>(cards);
                    Card card2 = newcards.remove(j);
                    Card card1 = newcards.remove(i);
                    Hand hand = new Hand(newcards);
                    Play play = new Play(hand, deck);
                    //plays.add(play);
                }
            }
            //Hand hand = new Hand(c1, c2, c3, c4);
            //Play play = new Play(hand, deck, c5, c6);
        }
    }

}
