// git_visual.go
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type GitVisual struct {
	repo     *git.Repository
	graph    bool
	all      bool
	decorate bool
	oneline  bool
	author   string
	since    string
	until    string
	maxCount int
	color    bool
}

func (g *GitVisual) openRepo(path string) error {
	var err error
	g.repo, err = git.PlainOpen(path)
	return err
}

func (g *GitVisual) getCommits() ([]*object.Commit, error) {
	var commits []*object.Commit
	var iter storer.RevIter
	var err error
	if g.all {
		iter, err = g.repo.Log(&git.LogOptions{All: true})
	} else {
		ref, err := g.repo.Head()
		if err != nil {
			return nil, err
		}
		iter, err = g.repo.Log(&git.LogOptions{From: ref.Hash()})
	}
	if err != nil {
		return nil, err
	}
	err = iter.ForEach(func(c *object.Commit) error {
		// Применяем фильтры
		if g.author != "" && !strings.Contains(c.Author.Name, g.author) {
			return nil
		}
		// Ограничение по дате можно добавить
		commits = append(commits, c)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if g.maxCount > 0 && len(commits) > g.maxCount {
		commits = commits[:g.maxCount]
	}
	return commits, nil
}

func (g *GitVisual) getRefs() (map[string][]string, error) {
	refs := make(map[string][]string)
	// Ветки
	branches, err := g.repo.Branches()
	if err != nil {
		return nil, err
	}
	err = branches.ForEach(func(ref *plumbing.Reference) error {
		hash := ref.Hash().String()[:7]
		refs[hash] = append(refs[hash], ref.Name().Short())
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Теги
	tags, err := g.repo.Tags()
	if err != nil {
		return nil, err
	}
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		hash := ref.Hash().String()[:7]
		refs[hash] = append(refs[hash], ref.Name().Short())
		return nil
	})
	if err != nil {
		return nil, err
	}
	// HEAD
	head, err := g.repo.Head()
	if err == nil {
		hash := head.Hash().String()[:7]
		refs[hash] = append(refs[hash], "HEAD")
	}
	return refs, nil
}

func (g *GitVisual) colorize(text, color string) string {
	if !g.color {
		return text
	}
	colors := map[string]string{
		"red":    "\033[91m",
		"green":  "\033[92m",
		"yellow": "\033[93m",
		"blue":   "\033[94m",
		"cyan":   "\033[96m",
		"reset":  "\033[0m",
	}
	return colors[color] + text + colors["reset"]
}

func (g *GitVisual) printCommit(c *object.Commit, prefix string, last bool, refs map[string][]string) {
	connector := "└── "
	if !last {
		connector = "├── "
	}
	hashShort := c.Hash.String()[:7]
	authorName := c.Author.Name
	date := c.Author.When.Format("2006-01-02 15:04")
	msg := strings.Split(c.Message, "\n")[0]

	refStr := ""
	if g.decorate {
		if r, ok := refs[hashShort]; ok {
			refStr = " (" + strings.Join(r, ", ") + ")"
		}
	}
	if g.oneline {
		fmt.Printf("%s%s%s%s%s\n", prefix, connector,
			g.colorize(hashShort, "yellow"),
			g.colorize(refStr, "blue"),
			" "+g.colorize(msg, "reset"))
	} else {
		fmt.Printf("%s%s%s%s\n", prefix, connector, g.colorize(hashShort, "yellow"), g.colorize(refStr, "blue"))
		fmt.Printf("%s    %s %s <%s>\n", prefix, g.colorize("Author:", "cyan"), authorName, c.Author.Email)
		fmt.Printf("%s    %s %s\n", prefix, g.colorize("Date:", "green"), date)
		fmt.Printf("%s    %s\n\n", prefix, g.colorize(msg, "reset"))
	}
}

func (g *GitVisual) showGraph() error {
	commits, err := g.getCommits()
	if err != nil {
		return err
	}
	refs, err := g.getRefs()
	if err != nil {
		return err
	}
	for i, c := range commits {
		last := i == len(commits)-1
		g.printCommit(c, "", last, refs)
	}
	return nil
}

func main() {
	var path, author, since, until string
	var graph, all, decorate, oneline, color bool
	var maxCount int
	flag.StringVar(&path, "path", ".", "Repository path")
	flag.BoolVar(&graph, "graph", true, "Show graph")
	flag.BoolVar(&all, "all", false, "Show all branches")
	flag.BoolVar(&decorate, "decorate", true, "Show refs")
	flag.BoolVar(&oneline, "oneline", false, "Compact output")
	flag.StringVar(&author, "author", "", "Filter by author")
	flag.StringVar(&since, "since", "", "Commits after date")
	flag.StringVar(&until, "until", "", "Commits before date")
	flag.IntVar(&maxCount, "n", 0, "Limit number of commits")
	flag.BoolVar(&color, "color", true, "Color output")
	flag.Parse()

	v := &GitVisual{
		graph:    graph,
		all:      all,
		decorate: decorate,
		oneline:  oneline,
		author:   author,
		since:    since,
		until:    until,
		maxCount: maxCount,
		color:    color,
	}
	if err := v.openRepo(path); err != nil {
		fmt.Fprintf(os.Stderr, "Error opening repo: %v\n", err)
		os.Exit(1)
	}
	if err := v.showGraph(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
