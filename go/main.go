/*
Generates table.json which maps 3-strokes to trigram types
- finger combo is a synonym of nstrokes
*/

package main

import (
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	"log"
	"os"
	"strings"
)

var FINGERS = []string{"LP", "LR", "LM", "LI", "LT", "RT", "RI", "RM", "RR", "RP"}
var BAD_RED_MAP = []bool{true, true, true, false, false, false, false, true, true, true}

func getFingerComboStr(finger0, finger1, finger2 uint) string {
	return FINGERS[finger0] + FINGERS[finger1] + FINGERS[finger2]
}

func ifElse(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

func assertOk(err error) {
	if err != nil {
		log.Fatalf("Assertion failed:\n%v", err)
	}
}

func marshalTable(table *orderedmap.OrderedMap[string, string]) []byte {
	sb := new(strings.Builder)
	sb.WriteString("{\n")
	lastIndex := table.Len() - 1
	for index, key := range table.Keys() {
		value, _ := table.Get(key)
		sb.WriteString(fmt.Sprintf("    \"%s\": \"%s\"", key, value))
		if index != lastIndex {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}
	sb.WriteString("}")
	return []byte(sb.String())
}

func writeTrigramTable() {
	table := orderedmap.NewOrderedMap[string, string]()

	for finger0 := uint(0); finger0 < 10; finger0++ {
		for finger1 := uint(0); finger1 < 10; finger1++ {
			for finger2 := uint(0); finger2 < 10; finger2++ {
				hand0 := finger0 >= 5
				hand1 := finger1 >= 5
				hand2 := finger2 >= 5
				fingerComboStr := getFingerComboStr(finger0, finger1, finger2)

				// check same finger
				if finger0 == finger1 || finger1 == finger2 {
					table.Set(fingerComboStr, ifElse(finger0 == finger2, "sft", "sfb"))
					continue
				}

				// alternates
				if hand0 != hand1 && hand1 != hand2 {
					table.Set(fingerComboStr, ifElse(finger0 == finger2, "alt-sfs", "alt"))
					continue
				}

				// red or oneh
				if hand0 == hand1 && hand1 == hand2 {
					// oneh
					towardsLeft := finger0 > finger1 && finger1 > finger2
					if towardsLeft || (finger0 < finger1 && finger1 < finger2) {
						table.Set(fingerComboStr, ifElse(towardsLeft == hand0, "inoneh", "outoneh"))
						continue
					}

					// red
					isBad := BAD_RED_MAP[finger0] && BAD_RED_MAP[finger1] && BAD_RED_MAP[finger2]
					isSfs := finger0 == finger2
					redType := "red"
					if isSfs {
						redType = ifElse(isBad, "bad-red-sfs", "red-sfs")
					} else if isBad {
						redType = "bad-red"
					}
					table.Set(fingerComboStr, redType)
					continue
				}

				// rolls
				roll0, roll1 := finger0, finger1
				if hand1 == hand2 {
					roll0, roll1 = finger1, finger2
				}
				table.Set(fingerComboStr, ifElse((roll0 > roll1) == hand1, "inroll", "outroll"))
			}
		}
	}

	jsonData := marshalTable(table)
	err := os.WriteFile("table.json", jsonData, os.ModePerm)
	assertOk(err)
}

func main() {
	writeTrigramTable()
}
