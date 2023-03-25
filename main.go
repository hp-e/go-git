package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
	// Set the path to the Git repository
	repoPath := "C:\\dev\\projects\\stages\\cloud"

	// Open the Git repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("Error opening repository: %s", err)
	}

	// Get the most recent tag reference
	// tagRefs, err := repo.Tags()
	// if err != nil {
	// 	log.Fatalf("Error getting tags: %s", err)
	// }
	// var recentTagRef *plumbing.Reference

	// err = tagRefs.ForEach(func(ref *plumbing.Reference) error {
	// 	// fmt.Printf("Name: %s, h:%v, trg; %v\n", ref.Name(), ref.Hash().String(), ref.Target())
	// 	// if recentTagRef == nil || recentTagRef.Target().IsTag() .Hash().Before(ref.Target().Hash()) {
	// 	recentTagRef = ref
	// 	// }
	// 	return nil
	// })

	// if err != nil {
	// 	log.Fatalf("Error finding recent tag: %s", err)
	// }
	// fmt.Println(recentTagRef.Name(), recentTagRef.Hash().String())

	var tag *object.Tag
	// var tagDate *time.Time
	tagos, _ := repo.TagObjects()
	_ = tagos.ForEach(func(t *object.Tag) error {
		// fmt.Println(t.Hash)
		tag = t
		// tagDate = &t.Tagger.When
		return nil
	})
	// if recentTagRef != nil {

	// t := time.Now().AddDate(-1, 0, 0)
	// Get the commit history from the recent tag up to HEAD
	commitIter, err := repo.Log(&git.LogOptions{
		// From:  tag.Hash,
		Order: git.LogOrderCommitterTime,
		Since: &tag.Tagger.When,
	})
	if err != nil {
		log.Fatalf("Error getting commit history: %v", err)
	}

	fmt.Printf("Commits since %v\n\n", tag.Tagger.When)
	// Print the commit messages
	var messages []string
	authors := make(map[string]author)
	err = commitIter.ForEach(func(c *object.Commit) error {

		aut, ok := authors[c.Author.Name]
		if !ok {
			authors[c.Author.Name] = author{
				name:    c.Author.Name,
				first:   c.Committer.When,
				last:    c.Committer.When,
				commits: 1,
			}
		} else {
			if aut.last.Before(c.Committer.When) {
				aut.last = c.Committer.When
			}
			if aut.first.After(c.Committer.When) {
				aut.first = c.Committer.When
			}
			aut.commits++
			authors[c.Author.Name] = aut
		}

		// authors[c.Author.Name] = authors[c.Author.Name] + 1
		messages = append(messages, c.Message)
		return nil
	})
	if err != nil {
		log.Fatalf("Error iterating through commits: %s", err)
	}

	// Print the commit messages
	for _, message := range messages {
		if strings.HasPrefix(message, "fix:") {
			fmt.Println("FIX")
			fmt.Println(message)
		}
	}

	for au, cnt := range authors {
		fmt.Printf("Author: %s, commits: %d, first: %v, last: %v\n", au, cnt.commits, cnt.first, cnt.last)
	}
	// }
}

type author struct {
	commits int
	first   time.Time
	last    time.Time
	name    string
}
