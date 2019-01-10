package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type infos struct {
	op        string
	address   string
	obcode    string
	loc       string
	label     string
	parameter string
	x         string
	siz       int
}

var results []infos

var r = map[string]string{
	"ADD":    "18",
	"ADDF":   "58",
	"ADDR":   "90",
	"AND":    "40",
	"CLEAR":  "B4",
	"COMP":   "28",
	"COMPF":  "88",
	"COMPR":  "A0",
	"DIV":    "24",
	"DIVF":   "64",
	"DIVR":   "9C",
	"FIX":    "C4",
	"FLOAT":  "C0",
	"HIO":    "F4",
	"J":      "3C",
	"JEQ":    "30",
	"JGT":    "34",
	"JLT":    "38",
	"JSUB":   "48",
	"LDA":    "00",
	"LDB":    "68",
	"LDCH":   "50",
	"LDF":    "70",
	"LDL":    "08",
	"LDS":    "6C",
	"LDT":    "74",
	"LDX":    "04",
	"LPS":    "D0",
	"MUL":    "20",
	"MULF":   "60",
	"MULR":   "98",
	"NORM":   "C8",
	"OR":     "44",
	"RD":     "D8",
	"RMO":    "AC",
	"RSUB":   "4C",
	"SHIFTL": "A4",
	"SHIFTR": "A8",
	"SIO":    "F0",
	"SSK":    "EC",
	"STA":    "0C",
	"STB":    "78",
	"STCH":   "54",
	"STF":    "80",
	"STI":    "D4",
	"STL":    "14",
	"STS":    "7C",
	"STSW":   "E8",
	"STT":    "84",
	"STX":    "10",
	"SUB":    "1C",
	"SUBF":   "5C",
	"SUBR":   "94",
	"SVC":    "B0",
	"TD":     "E0",
	"TIO":    "F8",
	"TIX":    "2C",
	"TIXR":   "B8",
	"WD":     "DC",
}

func main() {
	pass1()
	pass2()

	printob()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Input op code: ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		ops := queryOP(strings.ToUpper(text))
		for _, k := range ops {
			v := results[k]
			fmt.Printf("%5s %20s %20s %20s %20s\n", v.loc, v.label, v.op, v.parameter, v.obcode)
		}
	}

}

func pass1() {
	content, err := ioutil.ReadFile("file.txt")
	if err != nil {
		//Do something
	}
	lines := strings.Fields(string(content))
	loc_num := 0
	for k := 0; k < len(lines); k++ {

		var info infos
		info.address = ""
		info.label = ""
		info.loc = ""
		info.obcode = ""
		info.op = ""
		info.parameter = ""
		info.siz = 0
		info.x = "0"

		if lines[k] == "END" {
			info.op = "END"
			k = k + 1
			info.parameter = lines[k]
			//info.loc = tenTohex(loc_num)
			results = append(results, info)
			break
		} else if isop(lines[k]) {
			info.op = lines[k] //optab(lines[k])
			info.loc = tenTohex(loc_num)
			loc_num += 3
			if lines[k] != "RSUB" {
				k = k + 1
				info.parameter = lines[k]
			}
			results = append(results, info)
		} else {
			info.label = lines[k]
			k = k + 1
			if lines[k] == "START" {
				info.op = "START"
				k = k + 1
				info.loc = lines[k]
				info.parameter = lines[k]
				loc_num = hexToten(lines[k])
				results = append(results, info)
			} else if lines[k] == "BYTE" || lines[k] == "RESW" || lines[k] == "RESB" || lines[k] == "WORD" {
				info.op = lines[k]
				if lines[k] == "WORD" {
					var x int
					info.loc = tenTohex(loc_num)
					loc_num += 3
					info.op = lines[k]
					k = k + 1
					info.parameter = lines[k]
					x, _ = strconv.Atoi(lines[k])
					info.siz = x
					results = append(results, info)
				} else if lines[k] == "BYTE" {
					info.loc = tenTohex(loc_num)
					info.op = lines[k]
					k = k + 1
					x := lines[k]
					if x[0] == 'C' {
						loc_num += 3
					} else if x[0] == 'X' {
						loc_num += 1
					}
					info.parameter = lines[k]
					results = append(results, info)
				} else if lines[k] == "RESW" {
					info.loc = tenTohex(loc_num)
					loc_num += 3
					info.op = lines[k]
					k = k + 1
					info.parameter = lines[k]
					results = append(results, info)
				} else if lines[k] == "RESB" {
					var x int
					info.op = lines[k]
					k = k + 1
					info.parameter = lines[k]
					x, _ = strconv.Atoi(lines[k])
					info.loc = tenTohex(loc_num)
					loc_num += x
					info.siz = x
					results = append(results, info)
				}
			} else if isop(lines[k]) {

				info.op = lines[k] //optab(lines[k])
				info.loc = tenTohex(loc_num)
				loc_num += 3
				k = k + 1
				info.parameter = lines[k]

				results = append(results, info)
			}
		}
	}
}

func pass2() {
	for k := range results {
		var temp string

		if results[k].op == "START" {
			results[k].obcode = ""
		} else if optab(results[k].op) == "4C" {
			results[k].obcode = "4C0000"
		} else if results[k].op == "RESW" || results[k].op == "RESB" {
			results[k].obcode = ""
		} else if results[k].op == "WORD" {
			results[k].obcode = tenTohex(results[k].siz)
			for i := len(results[k].obcode); i < 6; i++ {
				results[k].obcode = "0" + results[k].obcode
			}

		} else if results[k].parameter != "" {
			results[k].obcode += optab(results[k].op)
			temp = results[k].parameter[len(results[k].parameter)-2 : len(results[k].parameter)]

			if temp == ",X" {
				results[k].x = "1"
				results[k].parameter = results[k].parameter

				for j := 0; j < len(results); j++ {
					if results[k].parameter == results[j].label {
						results[k].obcode += binaryTohex(results[k].x + hexTobinary(results[j].loc)[1:16])
					}

				}

			} else if results[k].op == "END" {
				results[k].obcode = ""
			} else if results[k].op == "BYTE" {
				if results[k].parameter[0] == 'C' {
					results[k].obcode = ASCII(results[k].parameter[2]) + ASCII(results[k].parameter[3]) + ASCII(results[k].parameter[4])

				}
				if results[k].parameter[0] == 'X' {
					results[k].obcode = results[k].parameter[2:4]
				}
			} else {
				for j := 0; j < len(results); j++ {
					if results[k].parameter == results[j].label {
						results[k].obcode += binaryTohex(results[k].x + hexTobinary(results[j].loc)[1:16])
					}
				}
			}
		}
	}
}

func optab(mnem string) string {

	return r[mnem]
}

func isop(lib string) bool {
	for k := range r {
		if k == lib {
			return true
		}
	}
	return false
}

func tenTohex(ten int) string {
	hexnum := fmt.Sprintf("%X", ten)

	return hexnum
}

func fixHexnum(num int, hexnum string) string {
	s := strconv.Itoa(num)
	return s + hexnum
}

func hexToten(hex string) int {
	if s, err := strconv.ParseInt(hex, 16, 32); err == nil {
		return int(s)
	}
	return 0
}

func binaryTohex(bin string) string {
	fixbin := bin
	var cou int

	cou = 4 - (len(bin) % 4)

	if cou == 4 {
		cou = 0
	}

	for i := 0; i < cou; i++ {
		fixbin = "0" + fixbin
	}

	oct, _ := strconv.ParseInt(fixbin, 2, 64)

	return fmt.Sprintf("%X", int(oct))
}

func hexTobinary(hex string) string {
	result := ""
	for i := 0; i < len(hex); i++ {
		oct, _ := strconv.ParseInt(string(hex[i]), 16, 64)

		t := fmt.Sprintf("%b", int64(oct))

		if len(t) < 4 {
			for j := 4 - len(t); j > 0; j-- {
				t = "0" + t
			}
		}

		result += t
	}

	return result
}

func ASCII(li byte) string {

	return fmt.Sprintf("%X", li)
}

func printob() {
	fmt.Printf("%5s %20s %20s %20s %20s\n", "Loc", "Label", "Op", "operands", "Object Code\n")

	for _, v := range results {
		fmt.Printf("%5s %20s %20s %20s %20s\n", v.loc, v.label, v.op, v.parameter, v.obcode)
	}
}

func queryOP(op string) []int {
	var x []int
	for k, v := range results {
		if v.op == op {
			x = append(x, k)
		}
	}

	return x
}
