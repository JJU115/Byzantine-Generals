/*
	BGenerals.go - A simulation of Byzantine fault tolerance via the Byzantine generals problem:
	The Byzantine Generals Problem - Leslie Lamport, Marshall Pease, Robert Shostak - ACM Transactions on Programming Languages and Systems 4, 3 (July 1982), 382-401 
	
	In this program each "lieutenant" is a goroutine among which they all send and receive messages to and from each other.
	The commander and each lieutenant have a 1/3 chance of being traitorous meaning they will send the opposite command
	to even numbered lieutenants.  

	At the end each lieutenant will report their final vote.

	Author: Justin Underhay
	Date of last modification: Nov 16, 2018
*/


package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"math/rand"
)

var recurse int
var numGenerals int
var commChannels []chan string
var finalStrings chan []string



func general(ID int, c_msg string, traitor bool) {

	var received = make([]string, 0)
	received = append(received, c_msg)

	for i := 1; i < numGenerals; i++ {
		if i != ID {
			if !traitor || i%2 != 0 {
				commChannels[i] <- c_msg + strconv.Itoa(ID)
			} else { 
				if strings.Contains(c_msg, "A") {
					commChannels[i] <- strings.Replace(c_msg, "A", "R", 1) + strconv.Itoa(ID)
				} else {
					commChannels[i] <- strings.Replace(c_msg, "R", "A", 1) + strconv.Itoa(ID)	
				}	
			}
		}
	}

	var path []int

	for {
		select {
		case next_msg := <-commChannels[ID]:
			received = append(received, next_msg)
			fmt.Println("Lieutenant", ID, "received", next_msg, "now at", received)
			path = transpose(next_msg[2:])
			if len(path) < recurse+1 {
				for j := 1; j < numGenerals; j++ {
					if !find(path, j) && j != ID {
						if !traitor || j%2 != 0 {
							commChannels[j] <- next_msg + strconv.Itoa(ID)
						} else { 
							if strings.Contains(next_msg, "A") {
								commChannels[j] <- strings.Replace(next_msg, "A", "R", 1) + strconv.Itoa(ID)
							} else {
								commChannels[j] <- strings.Replace(next_msg, "R", "A", 1) + strconv.Itoa(ID)	
							}	
						}	
					}
				}
			}

		case <-time.After(2 * time.Second):
			finalStrings <- append(received, strconv.Itoa(ID))
			break
		}
	}

	fmt.Println("Lieutenant",ID,"Done")
}



func find(S []int, s int) bool {
	for _, j := range S {
		if j == s {
			return true
		}
	}

	return false
}



func transpose(P string) []int {

	var intPath = make([]int, len(P))
	var err error

	S := strings.Split(P, "")

	for i := 0; i < len(P); i++ {
		intPath[i], err = strconv.Atoi(S[i])

		if err != nil {
			fmt.Println("Internal error, quitting...")
			os.Exit(-1)
		}
	}

	return intPath
}



func getVote(msg string, msgMap []string) string {

	if len(msg[2:]) == recurse+1 {
		return msg[:1]
	}

	var votes = make([]string, 0)
	
	for _,j := range msgMap {
		if len(j) > len(msg) && j[2:2+len(msg[2:])] == msg[2:] {
			votes = append(votes, getVote(j, msgMap))
		}
	}

	A := 0
	R := 0

	if msg[:1] == "A" {
		A++
	} else {
		R++
	}

	for _,j := range votes {
		if j == "A" {
			A++
		} else {
			R++
		}
	}

	if A > R {
		return "A"
	} else {
		return "R"
	}
}



func main() {

	if len(os.Args) != 4 {
		fmt.Println("Invocation error. Usage: BGenerals <recursion_level> <num_generals> <A | R>")
		os.Exit(-1)
	}

	var err error

	recurse, err = strconv.Atoi(os.Args[1])

	if err != nil {
		fmt.Println("Conversion error, quitting...")
		os.Exit(-1)
	}

	numGenerals, err = strconv.Atoi(os.Args[2])

	if err != nil || recurse >= numGenerals-1 {
		fmt.Println("Error on numGenerals input: Ensure that recurse < numGenerals-1")
		os.Exit(-1)
	}

	commOrder := os.Args[3]
	commChannels = make([]chan string, numGenerals)
	finalStrings = make(chan []string, numGenerals)

	for i := 0; i < numGenerals; i++ {
		commChannels[i] = make(chan string, numGenerals*numGenerals*recurse)
	}

	var commTraitor bool = false
	var ctOrder string

	if rand.Intn(4)%3 == 0 {
		commTraitor = true
		fmt.Println("The commander is a traitor!")
	}

	if commOrder == "A" {
		ctOrder = "R"
	} else {
		ctOrder = "A"
	}

	for j := 1; j < numGenerals; j++ {

		if commTraitor && j%2 == 0 {
			if rand.Intn(10)%3 != 0 {
				go general(j, ctOrder+" "+strconv.Itoa(0), false)
			} else {
				fmt.Println("Lieutenant",j,"is a traitor!")
				go general(j, ctOrder+" "+strconv.Itoa(0), true)
			}	
		} else {
			if rand.Intn(10)%3 != 0 {
				go general(j, commOrder+" "+strconv.Itoa(0), false)
			} else {
				fmt.Println("Lieutenant",j,"is a traitor!")
				go general(j, commOrder+" "+strconv.Itoa(0), true)
			}
		}		
	}

	for j := 1; j < numGenerals; j++ {
		msgs := <-finalStrings
		fmt.Println("Lieutenant",msgs[len(msgs)-1],"votes",getVote(msgs[0], msgs))
	}

}
