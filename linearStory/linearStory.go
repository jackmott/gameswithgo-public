// Optional homework:
// add a function that will insert a new page, after a given page
// add a function that will delete a page

package main

import "fmt"

type storyPage struct {
	text     string
	nextPage *storyPage
}

func (page *storyPage) playStory() {
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

// Delete

func main() {

	page1 := storyPage{"You are standing in an open field west of a white house.", nil}
	page1.addToEnd("You climb into the attic, it is pitch black, you can't see a thing!")
	page1.addToEnd("You are eaten by a Grue")

	page1.addAfter("Testing AddAfter")
	page1.playStory()

}
