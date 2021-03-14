package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type choice struct {
	cmd      string
	desc     string
	nextNode *storyNode
}

type storyNode struct {
	text    string
	choices []*choice
}

func (node *storyNode) addChoice(cmd string, desc string, nextNode *storyNode) {
	choice := &choice{cmd, desc, nextNode}

	node.choices = append(node.choices, choice)
}

func (node *storyNode) render() {
	fmt.Println((node.text))
	if node.choices != nil {
		for _, choice := range node.choices {
			fmt.Println(choice.cmd, ":", choice.desc)
		}
	}
}

func (node *storyNode) executeCmd(cmd string) *storyNode {
	for _, choice := range node.choices {
		if strings.ToLower(choice.cmd) == strings.ToLower(cmd) {
			return choice.nextNode
		}
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
/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*
/*   Wizards and Warlocks   /*
/*/*/*/*/*/*/*/*/*/*/*/*/*/*/*
	`}

	firstRoom := storyNode{text: `
	You are in a large chamber, deep underground
	You see three passages leading out. North, south and east.
	`}

	darkRoom := storyNode{text: "Nothin'"}
	darkRoomLit := storyNode{text: "Sumthin', you can go still north"}
	grue := storyNode{text: "While stumbling around in darkness you get eaten by a grue"}
	trap := storyNode{text: "You fall down a huge hole and you die miserably"}
	treasure := storyNode{text: "Goooooold!"}

	start.addChoice("", "Hit enter to start", &firstRoom)

	firstRoom.addChoice("n", "Go North", &darkRoom)
	firstRoom.addChoice("s", "Go South", &darkRoom)
	firstRoom.addChoice("e", "Go East", &trap)

	darkRoom.addChoice("s", "Try to go back", &grue)
	darkRoom.addChoice("O", "Turn on your lantern", &darkRoomLit)

	darkRoomLit.addChoice("n", "go north", &treasure)
	darkRoomLit.addChoice("s", "go south", &firstRoom)

	start.play()

	fmt.Println()
	fmt.Println("FIN")
}
