package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Deck []Card

var d Deck

var memLog bytes.Buffer

func Printf(text string, v ...interface{}) {
	text = fmt.Sprintf(text+"\n", v...)
	memLog.WriteString(text)
	fmt.Print(text)
}

func getText() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	memLog.WriteString(text)
	return strings.TrimSpace(text)
}

type Card struct {
	Front, Back string
	Wrong       int
}

func (c *Card) SetFront() {
	var front string
	Printf("The card:")
Loop:
	for {
		front = getText()
		for _, c := range d {
			if c.Front == front {
				Printf("The term %q already exists. Try again:", front)
				continue Loop
			}
		}
		c.Front = front
		return
	}
}

func (c *Card) SetBack() {
	var back string
	Printf("The definition of the card:")
Loop:
	for {
		back = getText()
		for _, c := range d {
			if c.Back == back {
				Printf("The definition %q already exists. Try again:", back)
				continue Loop
			}
		}
		c.Back = back
		return
	}
}

func (d *Deck) Add() {
	var c Card
	c.SetFront()
	c.SetBack()
	*d = append(*d, c)
	Printf("The pair (%q:%q) has been added", c.Front, c.Back)
}

func (d *Deck) Remove() {
	_d := *d
	Printf("Which card?")
	text := getText()
	for i, c := range *d {
		if c.Front == text {
			*d = append(_d[:i], _d[i+1:]...)
			Printf("The card has been removed")
			return
		}
	}
	Printf("Can't remove %q: there is no such card.", text)

}

func (d *Deck) Import(filepath string) {
	if filepath == "" {
		Printf("File name:")
		filepath = getText()
	}
	b, err := os.ReadFile(filepath)
	if err != nil {
		Printf("File not found")
		return
	}
	err = json.Unmarshal(b, d)
	if err != nil {
		log.Fatal(err)
	}
	Printf("%d cards have been loaded.", len(*d))
}

func (d *Deck) Export(filepath string) {
	if filepath == "" {
		Printf("File name")
		filepath = getText()

	}
	b, _ := json.Marshal(d)
	err := os.WriteFile(filepath, b, 0666)
	if err != nil {
		log.Fatal(err)
	}
	Printf("%d cards have been saved.", len(*d))
}

func (d *Deck) Ask() {
	rand.Seed(time.Now().UnixNano())
	_d := *d
	Printf("How many times to ask?")
	text := getText()
	nbrCards, _ := strconv.Atoi(text)
Loop:
	for i := 0; i < nbrCards; i++ {
		cardIndex := rand.Intn(len(*d))
		c := _d[cardIndex]
		Printf("Print the definition of %q:", c.Front)
		answer := getText()
		if answer == c.Back {
			Printf("Correct!")
			continue
		}
		_d[cardIndex].Wrong++
		for _, c2 := range *d {
			if answer == c2.Back && c != c2 {
				Printf("Wrong. The right answer is %q, but your definition is correct for %q.", c.Back, c2.Front)
				continue Loop
			}
		}
		Printf("Wrong. The right answer is %q.", c.Back)
	}
}

func (d *Deck) HardestCard() {
	// iterate over deck
	// find card with most cards wrong
	//if no cards that are wrong...
	var maxWrong int
	var wrongMap = make(map[int][]Card)
	for _, c := range *d {
		wrongMap[c.Wrong] = append(wrongMap[c.Wrong], c)
		if c.Wrong > maxWrong {
			maxWrong = c.Wrong
		}
	}
	if maxWrong == 0 {
		Printf("There are no cards with errors.")
		return
	}
	maxWrongSlice := wrongMap[maxWrong]
	if len(maxWrongSlice) == 1 {
		c := maxWrongSlice[0]
		Printf("The hardest card is %q. You have %d errors answering it.", c.Front, c.Wrong)
		return
	}
	var wrongTermsSlice []string
	for _, c := range maxWrongSlice {
		term := fmt.Sprintf("%q", c.Front)
		wrongTermsSlice = append(wrongTermsSlice, term)
	}
	wrongTerms := strings.Join(wrongTermsSlice, ", ")
	Printf("The hardest cards are %s", wrongTerms)
}

func (d *Deck) Log() {
	Printf("File name:")
	text := getText()
	err := os.WriteFile(text, memLog.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
	Printf("The log has been saved.")
}

func (d *Deck) ResetStats() {
	// don't think this will work
	_d := *d
	for i := range _d {
		_d[i].Wrong = 0
	}
	Printf("Card statistics have been reset.")
}

func main() {
	importFrom := flag.String("import_from", "", "json file to import from")
	exportTo := flag.String("export_to", "", "json file to export to")
	flag.Parse()
	if *importFrom != "" {
		d.Import(*importFrom)
	}
	for {
		Printf("Input the action (add, remove, import, export, ask, exit, log, hardest card, reset stats):")
		action := getText()
		switch action {
		case "add":
			d.Add()
		case "remove":
			d.Remove()
		case "import":
			d.Import("")
		case "export":
			d.Export("")
		case "ask":
			d.Ask()
		case "log":
			d.Log()
		case "hardest card":
			d.HardestCard()
		case "reset stats":
			d.ResetStats()
		case "exit":
			if *exportTo != "" {
				d.Export(*exportTo)
			}
			Printf("Bye bye!")
			os.Exit(0)
		default:
			Printf("That is not a valid option!")
		}
	}
}
