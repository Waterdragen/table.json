use std::fs::File;
use std::ptr::write;
///
/// Generates table.json which maps 3-strokes to trigram types
/// - finger combo is a synonym of nstrokes
///

use indexmap::IndexMap;
use fxhash::FxBuildHasher;

const FINGERS: [&'static str; 10] = ["LP", "LR", "LM", "LI", "LT", "RT", "RI", "RM", "RR", "RP"];
const BAD_RED_MAP: [bool; 10] = [true, true, true, false, false, false, false, true, true, true];
type FxIndexMap<K, V> = IndexMap<K, V, FxBuildHasher>;

fn get_finger_combo_str(finger0: usize, finger1: usize, finger2: usize) -> String {
    let mut s = String::with_capacity(6);
    s.push_str(FINGERS[finger0]);
    s.push_str(FINGERS[finger1]);
    s.push_str(FINGERS[finger2]);
    s
}

fn write_trigram_table() {
    let mut table: FxIndexMap<String, String> = FxIndexMap::default();

    (0..10).for_each(|finger0| {
        (0..10).for_each(|finger1| {
            (0..10).for_each(|finger2| {
                let hand0 = finger0 >= 5;
                let hand1 = finger1 >= 5;
                let hand2 = finger2 >= 5;
                let finger_combo_str = get_finger_combo_str(finger0, finger1, finger2);

                // check same finger
                if finger0 == finger1 || finger1 == finger2 {
                    table.insert(finger_combo_str, String::from(if finger0 == finger2 {"sft"} else {"sfb"}));
                    return;
                }

                // alternates
                if hand0 != hand1 && hand1 != hand2 {
                    table.insert(finger_combo_str, String::from(if finger0 == finger2 {"alt-sfs"} else {"alt"}));
                    return;
                }

                // red or oneh
                if hand0 == hand1 && hand1 == hand2 {
                    // oneh
                    let towards_left = finger0 > finger1 && finger1 > finger2;
                    if towards_left || finger0 < finger1 && finger1 < finger2 {
                        table.insert(finger_combo_str, String::from(if towards_left == hand0 {"inoneh"} else {"outoneh"}));
                        return;
                    }
                    // red
                    let is_bad = BAD_RED_MAP[finger0] && BAD_RED_MAP[finger1] && BAD_RED_MAP[finger2];
                    let is_sfs = finger0 == finger2;
                    let mut red_type = "red";
                    if is_sfs {
                        red_type = if is_bad {"bad-red-sfs"} else {"red-sfs"};
                    } else if is_bad {
                        red_type = "bad-red";
                    }
                    table.insert(finger_combo_str, String::from(red_type));
                    return;
                }

                // rolls
                let (roll0, roll1) = if hand0 == hand1 {(finger0, finger1)} else {(finger1, finger2)};
                table.insert(finger_combo_str, String::from(if (roll0 > roll1) == hand1 {"inroll"} else {"outroll"}));
            });
        });
    });

    let file = File::create("table.json").unwrap();
    serde_json::to_writer_pretty(file, &table).unwrap()
}

fn main() {
    write_trigram_table();
}
