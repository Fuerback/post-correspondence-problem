package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Domino struct {
	top    string
	bottom string
}

type Solutions struct {
	indices  []int
	dominos  []Domino
	diff     string
	diffSide string
}

func main() {
	start := time.Now()
	d := []Domino{}
	mapDominos := [][]Domino{}
	var count int
	//values := make(chan string)
	//defer close(values)

	//scanner := bufio.NewScanner(os.Stdin)
	file, err := os.Open("./sample2.in")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		n, err := strconv.Atoi(text)
		if err == nil {
			count = n
		} else {
			stringSlices := strings.Split(text, " ")
			if stringSlices[0] != stringSlices[1] {
				domino := Domino{top: stringSlices[0], bottom: stringSlices[1]}
				d = append(d, domino)
			}
			count--

			if count == 0 {
				mapDominos = append(mapDominos, d)
				d = []Domino{}
			}
		}
	}

	for index, dms := range mapDominos {
		solvePCP(index, dms)
	}
	duration := time.Since(start)
	fmt.Println(duration)
}

func isSolvable(dominos []Domino) bool {
	var hasPrefix bool
	var hasSufix bool

	var topLonger bool
	var bottomLonger bool

	for _, d := range dominos {
		if hasPrefix && hasSufix && topLonger && bottomLonger {
			break
		}
		if strings.HasPrefix(d.top, d.bottom) || strings.HasPrefix(d.bottom, d.top) && !hasPrefix {
			hasPrefix = true
		}
		if strings.HasSuffix(d.top, d.bottom) || strings.HasSuffix(d.bottom, d.top) && !hasSufix {
			hasSufix = true
		}
		if len(d.top) > len(d.bottom) && !topLonger {
			topLonger = true
		}
		if len(d.bottom) > len(d.top) && !bottomLonger {
			bottomLonger = true
		}
	}

	return hasPrefix && hasSufix && bottomLonger && topLonger
}

func solvePCP(index int, dominos []Domino) {
	if isSolvable(dominos) {
		solutions := []Solutions{}

		for i := 0; i < len(dominos); i++ {
			if len(dominos[i].top) == len(dominos[i].bottom) {
				if dominos[i].top != dominos[i].bottom {
					continue
				}
			}
			if len(dominos[i].top) < len(dominos[i].bottom) {
				if !strings.HasPrefix(dominos[i].bottom, dominos[i].top) {
					continue
				}
			}
			if len(dominos[i].bottom) < len(dominos[i].top) {
				if !strings.HasPrefix(dominos[i].top, dominos[i].bottom) {
					continue
				}
			}

			s := NewSolution(dominos, []int{i})
			solutions = append(solutions, *s)
		}

		validSolutions := getValidSolutions(solutions)

		// improve performance
		depth := 70
		for i := 0; i < depth; i++ {
			if len(validSolutions) > 0 || len(solutions) == 0 {
				break
			}
			validSolutions, solutions = getSolutions(validSolutions, solutions, dominos)
		}

		if len(validSolutions) == 0 {
			fmt.Printf("Case %d: %s\n", index+1, "IMPOSSIBLE")
		} else {
			fmt.Printf("Case %d: %s\n", index+1, getResult(validSolutions))
		}
	} else {
		fmt.Printf("Case %d: %s\n", index+1, "IMPOSSIBLE")
	}
}

func getSolutions(validSolutions []Solutions, solutions []Solutions, dominos []Domino) ([]Solutions, []Solutions) {

	newSolutions := []Solutions{}

	for j := 0; j < len(solutions); j++ {

		if solutions[j].diffSide == "x" {
			for k := 0; k < len(dominos); k++ {
				var pref string
				if len(solutions[j].diff) < len(dominos[k].top) {
					pref = solutions[j].diff
				} else {
					pref = solutions[j].diff[0:len(dominos[k].top)]
				}

				if strings.HasPrefix(dominos[k].top, pref) {
					newX := solutions[j].getTop() + dominos[k].top
					newY := solutions[j].getBottom() + dominos[k].bottom
					if len(newX) == len(newY) {
						if newX != newY {
							continue
						}
					}
					if len(newX) > len(newY) {
						if !strings.HasPrefix(newX, newY) {
							continue
						}
					}
					if len(newY) > len(newX) {
						if !strings.HasPrefix(newY, newX) {
							continue
						}
					}

					oldIndices := []int{}
					for p := 0; p < len(solutions[j].indices); p++ {
						oldIndices = append(oldIndices, solutions[j].indices[p])
					}
					oldIndices = append(oldIndices, k)
					s := NewSolution(dominos, oldIndices)
					newSolutions = append(newSolutions, *s)
				}
			}
		}

		if solutions[j].diffSide == "y" {
			for k := 0; k < len(dominos); k++ {
				var pref string
				if len(solutions[j].diff) < len(dominos[k].bottom) {
					pref = solutions[j].diff
				} else {
					pref = solutions[j].diff[0:len(dominos[k].bottom)]
				}
				if strings.HasPrefix(dominos[k].bottom, pref) {
					newX := solutions[j].getTop() + dominos[k].top
					newY := solutions[j].getBottom() + dominos[k].bottom

					if len(newX) == len(newY) {
						if newX != newY {
							continue
						}
					}
					if len(newX) > len(newY) {
						if !strings.HasPrefix(newX, newY) {
							continue
						}
					}
					if len(newY) > len(newX) {
						if !strings.HasPrefix(newY, newX) {
							continue
						}
					}

					oldIndices := []int{}
					for p := 0; p < len(solutions[j].indices); p++ {
						oldIndices = append(oldIndices, solutions[j].indices[p])
					}
					oldIndices = append(oldIndices, k)
					s := NewSolution(dominos, oldIndices)
					newSolutions = append(newSolutions, *s)
				}
			}
		}
	}
	validSolutions = getValidSolutions(newSolutions)

	solutions = newSolutions
	return validSolutions, solutions
}

func getResult(s []Solutions) string {
	var finalSolution string

	for _, result := range s {
		if len(finalSolution) == 0 || len(result.getTop()) < len(finalSolution) {
			finalSolution = result.getTop()
		}
		if len(finalSolution) > 0 && len(result.getTop()) == len(finalSolution) && strings.Compare(result.getTop(), finalSolution) < 0 {
			finalSolution = result.getTop()
		}
	}

	return finalSolution
}

func getValidSolutions(s []Solutions) []Solutions {
	validSolutions := []Solutions{}
	for i := 0; i < len(s); i++ {
		if s[i].isValidSolution() {
			validSolutions = append(validSolutions, s[i])
		}
	}
	return validSolutions
}

func NewSolution(d []Domino, i []int) *Solutions {
	s := &Solutions{
		dominos: d,
		indices: i,
	}
	s.updateDiffs()
	return s
}

func (s *Solutions) addIndice(i int) {
	s.indices = append(s.indices, i)
	s.updateDiffs()
}

func (s *Solutions) updateDiffs() {
	var x = ""
	var y = ""

	for i := 0; i < len(s.indices); i++ {
		x += s.dominos[s.indices[i]].top
		y += s.dominos[s.indices[i]].bottom
	}

	if len(x) > len(y) {
		s.diffSide = "y"
		lenY := len(y)
		lenX := len(x)
		s.diff = x[lenY:lenX]
	} else if len(x) < len(y) {
		s.diffSide = "x"
		lenY := len(y)
		lenX := len(x)
		s.diff = y[lenX:lenY]
	} else {
		s.diffSide = ""
		s.diff = ""
	}
}

func (s *Solutions) isValidSolution() bool {
	var x = ""
	var y = ""

	for i := 0; i < len(s.indices); i++ {
		x += s.dominos[s.indices[i]].top
		y += s.dominos[s.indices[i]].bottom
	}

	return len(x) == len(y)
}

func (s *Solutions) getTop() string {
	var x = ""

	for i := 0; i < len(s.indices); i++ {
		x += s.dominos[s.indices[i]].top
	}

	return x
}

func (s *Solutions) getBottom() string {
	var y = ""

	for i := 0; i < len(s.indices); i++ {
		y += s.dominos[s.indices[i]].bottom
	}

	return y
}
