"""
Generates table.json which maps 3-strokes to trigram types
- finger combo is a synonym of nstrokes
"""

import json
import itertools

FINGERS = ["LP", "LR", "LM", "LI", "LT", "RT", "RI", "RM", "RR", "RP"]

def get_finger_combo_str(finger_combo: tuple[int, int, int]) -> str:
    return "".join(FINGERS[i] for i in finger_combo)

def write_trigram_table():
    table: dict[str, str] = {}

    bad_red_map = [1, 1, 1, 0, 0, 0, 0, 1, 1, 1]
    finger_combos = itertools.product(range(10), repeat=3)
    finger_combo: tuple[int, int, int]
    for finger_combo in finger_combos:
        finger0 = finger_combo[0]
        finger1 = finger_combo[1]
        finger2 = finger_combo[2]
        hand0 = finger0 >= 5
        hand1 = finger1 >= 5
        hand2 = finger2 >= 5
        finger_combo_str = get_finger_combo_str(finger_combo)

        # check same finger
        if finger0 == finger1 or finger1 == finger2:
            table[finger_combo_str] = "sft" if finger0 == finger2 else "sfb"
            continue

        # alternates
        if hand0 != hand1 and hand1 != hand2:
            table[finger_combo_str] = "alt-sfs" if finger0 == finger2 else "alt"
            continue

        # red or oneh
        if hand0 == hand1 == hand2:
            # oneh
            if (towards_left := finger0 > finger1 > finger2) or finger0 < finger1 < finger2:
                table[finger_combo_str] = "inoneh" if towards_left == hand0 else "outoneh"
                continue
            # red
            is_bad = (bad_red_map[finger0] + bad_red_map[finger1] + bad_red_map[finger2]) == 3
            is_sfs = finger0 == finger2
            red_type = "red"
            if is_sfs:
                red_type = "bad-red-sfs" if is_bad else "red-sfs"
            elif is_bad:
                red_type = "bad-red"
            table[finger_combo_str] = red_type
            continue

        # rolls
        (roll0, roll1) = (finger0, finger1) if hand0 == hand1 else (finger1, finger2)
        table[finger_combo_str] = "inroll" if (roll0 > roll1) == hand1 else "outroll"

    with open("table.json", 'w') as f:
        json.dump(table, f, indent=4)

if __name__ == '__main__':
    write_trigram_table()
