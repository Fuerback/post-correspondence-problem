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

type Instance struct {
	SavedResult  []Result
	SavedDominos []Domino
}

type PCP struct {
	Dominos []Domino
}

var count int

func main() {
	pcp := PCP{}

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
		var process Instance
		if count == 0 {
			updateCount(text)
			pcp.Dominos = []Domino{}
			process := Instance{}
			process.SavedDominos = pcp.Dominos
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
					process.SavedDominos = pcp.Dominos
					if retInst, err := pcp.recursiveSolve(process); err != nil {
						fmt.Println("err:", err)
					} else {
						fmt.Println(retInst.GetResult())
					}
				}
			}
		}
	}
}

func (p *PCP) recursiveSolve(cur Instance) (Instance, error) {
	if cur.isCyclicResult() {
		return cur, errors.New(" Cyclic Result .....")
	}

	if cur.isResultReach() {
		return cur, nil
	}

	for index, dom := range p.Dominos {
		if p.IsDominoValid(cur, dom) {
			cur, _ = p.ApplyDomino(cur, index)
			var err error
			if len(cur.GetString(0)) > 100 {
				return cur, errors.New("Too long")
			}
			cur, err = p.recursiveSolve(cur)
			if err == nil {
				return cur, nil
			}

		}
	}
	return cur, errors.New("Don't have result")
}

func (p *PCP) IsDominoValid(curState Instance, inputDomino Domino) bool {
	strTop := curState.GetString(0)
	strBottom := curState.GetString(1)

	tempA := strTop + inputDomino.top
	tempB := strBottom + inputDomino.bottom

	prefix, exist := getSubsetPrefix(tempA, tempB)
	if !exist {
		return false
	}

	return tempA == prefix || tempB == prefix
}

func (p *PCP) ApplyDomino(curState Instance, dominoIndex int) (Instance, error) {
	newDom := p.Dominos[dominoIndex]

	if p.IsDominoValid(curState, newDom) {
		newRet := Result{}
		newRet.PotentialResult = dominoIndex
		if newDiff, err := p.CheckDiff(curState, newDom); err == nil {
			newRet.CurrentDiff = newDiff
			curState.SavedResult = append(curState.SavedResult, newRet)
			return curState, nil
		}

		return Instance{}, errors.New("Diff error on apply Domino")
	}
	return Instance{}, errors.New("Domino not valid in apply Domino")
}

func (p *PCP) CheckDiff(curState Instance, dom Domino) (Diff, error) {
	strTop := curState.GetString(0) + dom.top
	strBottom := curState.GetString(1) + dom.bottom
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

// Instance

func (c *Instance) GetResult() string {
	var finalResult string

	for _, result := range c.SavedResult {
		finalResult += c.SavedDominos[result.PotentialResult].top
	}
	return finalResult
}

func (c *Instance) GetString(index int) string {
	var dominosString string

	for _, result := range c.SavedResult {
		if index == 0 {
			dominosString += c.SavedDominos[result.PotentialResult].top
		} else {
			dominosString += c.SavedDominos[result.PotentialResult].bottom
		}
	}
	return dominosString
}

func (c *Instance) isResultReach() bool {
	if len(c.GetString(0)) == 0 && len(c.GetString(1)) == 0 {
		return false
	}
	return c.GetString(0) == c.GetString(1)
}

func (c *Instance) isCyclicResult() bool {
	if len(c.SavedResult) == 0 {
		return false
	}

	checkingRet := c.SavedResult[len(c.SavedResult)-1]
	for i := 0; i < len(c.SavedResult)-1; i++ {
		ret := c.SavedResult[i]
		//Find save result list has the same
		if ret.PotentialResult == checkingRet.PotentialResult && ret.CurrentDiff == checkingRet.CurrentDiff {
			return true
		}
	}
	return false
}

func (c *Instance) GetCurrentResult() []int {
	var retInt []int
	for i := 0; i < len(c.SavedResult); i++ {
		retInt = append(retInt, c.SavedResult[i].PotentialResult)
	}
	return retInt
}
