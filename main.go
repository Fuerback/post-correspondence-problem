package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Domino struct {
	top    string
	bottom string
}

type Diff struct {
	DiffCompare int
	DiffDomino  string
}

type Result struct {
	PotentialResult int
	CurrentDiff     Diff
}

type PCP struct {
	Dominos      []Domino
	SavedResult  []Result
	SavedDominos []Domino
}

var count int

func main() {
	var pcp PCP

	// read input
	//scanner := bufio.NewScanner(os.Stdin)

	file, err := os.Open("./sample.in")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if count == 0 {
			updateCount(text)
			pcp = PCP{}
			pcp.SavedDominos = pcp.Dominos
		} else {
			stringSlices := strings.Split(text, " ")
			if stringSlices[0] != stringSlices[1] {
				domino := Domino{top: stringSlices[0], bottom: stringSlices[1]}
				pcp.Dominos = append(pcp.Dominos, domino)
			}
			count--
			if count == 0 {
				if pcp.isUnsovable() {
					fmt.Println("IMPOSSIBLE")
				} else {
					pcp.SavedDominos = pcp.Dominos
					if err := pcp.recursiveSolve(); err != nil {
						fmt.Println("err:", err)
					} else {
						fmt.Println(pcp.GetResult())
					}
				}
			}
		}
	}
}

func (p *PCP) recursiveSolve() error {
	if p.isResultReach() {
		return nil
	}

	for index, dom := range p.Dominos {
		if p.IsDominoValid(dom) {
			p.ApplyDomino(index)
			var err error
			if len(p.GetString(0)) > 100 {
				return errors.New("Too long")
			}
			err = p.recursiveSolve()
			if err == nil {
				return nil
			}

		}
	}
	return errors.New("Don't have result")
}

func (p *PCP) IsDominoValid(inputDomino Domino) bool {
	strTop := p.GetString(0)
	strBottom := p.GetString(1)

	tempA := strTop + inputDomino.top
	tempB := strBottom + inputDomino.bottom

	prefix, exist := getSubsetPrefix(tempA, tempB)
	if !exist {
		return false
	}

	return tempA == prefix || tempB == prefix
}

func (p *PCP) ApplyDomino(dominoIndex int) error {
	newDom := p.Dominos[dominoIndex]

	if p.IsDominoValid(newDom) {
		newRet := Result{}
		newRet.PotentialResult = dominoIndex
		if newDiff, err := p.CheckDiff(newDom); err == nil {
			newRet.CurrentDiff = newDiff
			p.SavedResult = append(p.SavedResult, newRet)
			return nil
		}

		return errors.New("Diff error on apply Domino")
	}
	return errors.New("Domino not valid in apply Domino")
}

func (p *PCP) CheckDiff(dom Domino) (Diff, error) {
	strTop := p.GetString(0) + dom.top
	strBottom := p.GetString(1) + dom.bottom
	retDiff := Diff{}
	retDiff.DiffCompare = strings.Compare(strTop, strBottom)

	if retDiff.DiffCompare == 0 {
		return retDiff, nil
	}

	if retDiff.DiffCompare == 1 {
		retDiff.DiffDomino = strings.TrimPrefix(strTop, strBottom)
	} else {
		retDiff.DiffDomino = strings.TrimPrefix(strBottom, strTop)
	}

	return retDiff, nil
}

func updateCount(text string) {
	n, err := strconv.Atoi(text)
	if err != nil {
		fmt.Println("count not updated")
	}
	count = n
}

func (p *PCP) isUnsovable() bool {
	var hasTopLongerThanBottom bool
	var hasBottomLongerThanTop bool

	for _, domino := range p.Dominos {
		if len(domino.top) > len(domino.bottom) {
			hasTopLongerThanBottom = true
			break
		}
	}

	for _, domino := range p.Dominos {
		if len(domino.bottom) > len(domino.top) {
			hasBottomLongerThanTop = true
			break
		}
	}

	if hasTopLongerThanBottom && hasBottomLongerThanTop {
		return false
	}

	return true
}

func getSubsetPrefix(str1, str2 string) (string, bool) {
	findSubset := false
	for i := 0; i < len(str1) && i < len(str2); i++ {
		if str1[i] != str2[i] {
			retStr := str1[:i]
			return retStr, findSubset
		}
		findSubset = true
	}

	if len(str1) > len(str2) {
		return str2, findSubset
	} else if len(str1) == len(str2) {
		return str1, str1 == str2
	}

	return str1, findSubset
}

func (p *PCP) GetResult() string {
	var finalResult string

	for _, result := range p.SavedResult {
		finalResult += p.SavedDominos[result.PotentialResult].top
	}
	return finalResult
}

func (p *PCP) GetString(index int) string {
	var dominosString string

	for _, result := range p.SavedResult {
		if index == 0 {
			dominosString += p.SavedDominos[result.PotentialResult].top
		} else {
			dominosString += p.SavedDominos[result.PotentialResult].bottom
		}
	}
	return dominosString
}

func (p *PCP) isResultReach() bool {
	if len(p.GetString(0)) == 0 && len(p.GetString(1)) == 0 {
		return false
	}
	return p.GetString(0) == p.GetString(1)
}

func (p *PCP) isCyclicResult() bool {
	if len(p.SavedResult) == 0 {
		return false
	}

	checkingRet := p.SavedResult[len(p.SavedResult)-1]
	for i := 0; i < len(p.SavedResult)-1; i++ {
		ret := p.SavedResult[i]
		//Find save result list has the same
		if ret.PotentialResult == checkingRet.PotentialResult && ret.CurrentDiff == checkingRet.CurrentDiff {
			return true
		}
	}
	return false
}

func (p *PCP) GetCurrentResult() []int {
	var retInt []int
	for i := 0; i < len(p.SavedResult); i++ {
		retInt = append(retInt, p.SavedResult[i].PotentialResult)
	}
	return retInt
}
