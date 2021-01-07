package main

import (
	"github.com/tbellembois/gobkm-gio/api"
	. "github.com/tbellembois/gobkm/types"
)

func saveRemoteFolder(f Folder) {
	go func() {
		var (
			e error
		)
		if e = api.SaveFolder(f); e != nil {
			addError(e)
		}
		chFEdit <- folderEdited
	}()
}

func saveRemoteBookmark(b Bookmark) {
	go func() {
		var (
			e error
		)
		if e = api.SaveBookmark(b); e != nil {
			addError(e)
		}
		chBEdit <- bookmarkEdited
	}()
}

func addRemoteBookmark(b Bookmark) {
	go func() {
		var (
			e    error
			newB Bookmark
		)
		if newB, e = api.AddBookmark(b); e != nil {
			addError(e)
		}
		chBAdd <- newB
	}()
}

func starBookmark(b Bookmark, star bool) {
	go func() {
		var (
			e error
		)
		if e = api.StarBookmark(b, star); e != nil {
			addError(e)
		}
		chBStar <- "ok"
	}()
}

func deleteRemoteFolder(id int) {
	go func() {
		var (
			e error
		)
		if e = api.DeleteFolder(id); e != nil {
			addError(e)
		}
		chFDel <- "ok"
	}()
}

func deleteRemoteBookmark(b Bookmark) {
	go func() {
		var (
			e error
		)
		if e = api.DeleteBookmark(b.Id); e != nil {
			addError(e)
		}
		chBDel <- b.Title
	}()
}

func addRemoteFolder(f Folder) {
	go func() {
		var (
			e    error
			newF Folder
		)
		if newF, e = api.AddFolder(f); e != nil {
			addError(e)
		}
		chFAdd <- newF
	}()
}

func getRemoteTags() {
	go func() {
		var (
			t []Tag
			e error
		)
		if t, e = api.GetTags(); e != nil {
			addError(e)
		}
		chT <- t
	}()
}

func getRemoteStars() {
	go func() {
		var (
			b []Bookmark
			e error
		)
		if b, e = api.GetStars(); e != nil {
			addError(e)
		}
		chS <- b
	}()
}

func getRemoteNode(id int) {
	go func() {
		var (
			n Folder
			e error
		)
		if n, e = api.GetNode(id); e != nil {
			addError(e)
		}
		chF <- n
	}()
}

func doSearch(s string) {
	go func() {
		var (
			bs []Bookmark
			e  error
		)
		if bs, e = api.SearchBookmark(s); e != nil {
			addError(e)
		}
		searchResults = bs
	}()

}
