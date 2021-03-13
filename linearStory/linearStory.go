package main

import (
	"fmt"
)

type storyPage struct {
	text     string
	nextPage *storyPage
}

func (page *storyPage) playStory() {
	// if page == nil {
	// 	fmt.Println("END")
	// 	return
	// }
	// scanner := bufio.NewScanner(os.Stdin)
	// fmt.Println((page.text))
	// scanner.Scan()
	// page.nextPage.playStory()

	for page != nil {
		fmt.Println(page.text)
		page = page.nextPage
	}
}

func (page *storyPage) addToEnd(text string) {
	for page.nextPage != nil {
		page = page.nextPage
	}
	page.nextPage = &storyPage{text, nil}
}

func (page *storyPage) addAfter(text string) {
	newPage := &storyPage{text, page.nextPage}
	page.nextPage = newPage
}

func main() {

	page1 := storyPage{"It was a dark and stormy night.", nil}
	page1.addToEnd("You are alone, and you need to find the scared helmet, before the bad guys do")
	page1.addToEnd("You see a troll ahead")

	page1.addAfter("Bla")
	page1.playStory()

}
