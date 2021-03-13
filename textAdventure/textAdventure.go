package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type choices struct {
	cmd        string
	desc       string
	nextNode   *storyNode
	nextChoice *choices
}

type storyNode struct {
	text    string
	choices *choices
}

func (node *storyNode) addChoice(cmd string, desc string, nextNode *storyNode) {
	choice := &choices{cmd, desc, nextNode, nil}

	if node.choices == nil {
		node.choices = choice
	} else {
		currentChoice := node.choices
		for currentChoice.nextChoice != nil {
			currentChoice = currentChoice.nextChoice
		}
		currentChoice.nextChoice = choice
	}
}

func (node *storyNode) render() {
	fmt.Println((node.text))
	currentChoice := node.choices
	for currentChoice != nil {
		fmt.Println(currentChoice.cmd, ":", currentChoice.desc)
		currentChoice = currentChoice.nextChoice
	}
}

func (node *storyNode) executeCmd(cmd string) *storyNode {
	currentChoice := node.choices
	for currentChoice != nil {
		if strings.ToLower(currentChoice.cmd) == strings.ToLower(cmd) {
			return currentChoice.nextNode
		}
		currentChoice = currentChoice.nextChoice
	}
	fmt.Println("Sorry, I didn't understand that")
	return node
}

var scanner *bufio.Scanner

func (node *storyNode) play() {
	node.render()
	if node.choices != nil {
		scanner.Scan()
		node.executeCmd(scanner.Text()).play()
	}
}

func main() {
	scanner = bufio.NewScanner((os.Stdin))

	start := storyNode{text: `
	You are in a large chamber, deep underground
	You see three passages leading out. North, south and east.
	`}

	darkRoom := storyNode{text: "Nothin'"}
	darkRoomLit := storyNode{text: "Sumthin', you can go still north"}
	grue := storyNode{text: "While stumbling around in darkness you get eaten by a grue"}
	trap := storyNode{text: "You fall down a huge hole and you die miserably"}
	treasure := storyNode{text: "Goooooold!"}

	start.addChoice("n", "Go North", &darkRoom)
	start.addChoice("s", "Go South", &darkRoom)
	start.addChoice("e", "Go East", &trap)

	darkRoom.addChoice("s", "Try to go back", &grue)
	darkRoom.addChoice("O", "Turn on your lantern", &darkRoomLit)

	darkRoomLit.addChoice("n", "go north", &treasure)
	darkRoomLit.addChoice("s", "go south", &start)

	start.play()

	fmt.Println()
	fmt.Println("FIN")
}
