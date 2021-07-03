package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/clipboard"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/tbellembois/gobkm-gio/datastores"
	"github.com/tbellembois/gobkm-gio/globals"
	. "github.com/tbellembois/gobkm/types"
	"golang.org/x/exp/shiny/materialdesign/icons"

	// Android permissions
	_ "net"

	_ "gioui.org/app/permission/storage"
)

type (
	D = layout.Dimensions
	C = layout.Context
)

var (
	colorFolder      = color.NRGBA{66, 139, 202, 0xff}
	colorBookmark    = color.NRGBA{92, 184, 92, 0xff}
	colorCut         = color.NRGBA{238, 238, 238, 0xff}
	colorPaste       = color.NRGBA{153, 153, 153, 0xff}
	colorTags        = color.NRGBA{217, 83, 79, 0xff}
	colorStars       = color.NRGBA{217, 83, 79, 0xff}
	colorToggleMenus = color.NRGBA{217, 83, 79, 0xff}
	colorPreferences = color.NRGBA{153, 153, 153, 0xff}
	colorBreadcrum   = color.NRGBA{255, 255, 255, 0xff}
	colorHistory     = color.NRGBA{204, 204, 204, 0xff}
	colorSearch      = color.NRGBA{153, 153, 153, 0xff}
	colorNewFolder   = color.NRGBA{153, 153, 153, 0xff}
	colorNewBookmark = color.NRGBA{153, 153, 153, 0xff}
	colorError       = color.NRGBA{255, 0, 0, 0xff}
)

var (
	dstore datastores.Datastore

	w   *app.Window
	gtx layout.Context

	stars         []Bookmark // favorite bookmarks
	tags          []Tag      // bookmark tags
	currentFolder Folder     // currently displayed folder
	currentError  []string   // list of returned errors

	tooglestarClicked,
	toogletagClicked,
	tooglemenuClicked,
	parameterClicked bool // keeps a click on the corresponding buttons
	folderMenuClicked,
	folderCutClicked,
	folderParentCutClicked,
	folderPasteClicked int // keeps the id of the folder on which the action if performed
	bookmarkMenuClicked,
	bookmarkCutClicked,
	bookmarkParentCutClicked,
	bookmarkPasteClicked int // keeps the id of the bookmark on which the action if performed
	deleteConfirmed bool

	folderEdited   Folder            // the edited folder = Folder{} if no folder is currently edited
	bookmarkEdited Bookmark          // the edited bookmark = Bookmark{} if no bookmark is currently edited
	breadcrum      []Folder          // breadcrum folders
	history        []Folder          // folder history
	chF            chan (Folder)     // Folder retrieval channel
	chS            chan ([]Bookmark) // Stars retrieval channel
	chT            chan ([]Tag)      // Tags retrieval channel
	chFEdit        chan (Folder)     // Folder edit channel
	chBEdit        chan (Bookmark)   // Bookmark edit channel
	chBAdd         chan (Bookmark)   // Bookmark add channel
	chFAdd         chan (Folder)     // Folder add channel
	chBDel         chan (string)     // Bookmark del channel
	chFDel         chan (string)     // Folder del channel
	chBStar        chan (string)     // Star bookmark channel
	chFMove        chan (Folder)     // Folder move channel
	chBMove        chan (Bookmark)   // Bookmark move channel
	chEDel         chan (string)     // Error deletion channel
	gui            GUI               // struct containing the GUI elements
	searchResults  []Bookmark        // list of search results
)

func init() {

	dstore = datastores.NewDatastore()

	var (
		e error
		p datastores.Preferences
	)

	if e = dstore.InitDatastore(); e != nil {
		addError(e)
	}
	if p, e = dstore.LoadPreferences(); e != nil {
		addError(e)
	}

	globals.ServerURL = p.ServerURL
	globals.ServerLogin = p.ServerUsername
	globals.ServerPassword = p.ServerPassword
	globals.HistorySize = p.HistorySize

	globals.B64Auth = basicAuth(string(globals.ServerLogin), string(globals.ServerPassword))

	chF = make(chan Folder)
	chS = make(chan []Bookmark)
	chT = make(chan []Tag)
	chFEdit = make(chan Folder)
	chBEdit = make(chan Bookmark)
	chBAdd = make(chan Bookmark)
	chFAdd = make(chan Folder)
	chFDel = make(chan string)
	chBDel = make(chan string)
	chBStar = make(chan string)
	chFMove = make(chan Folder)
	chBMove = make(chan Bookmark)
	chEDel = make(chan string)

	gui = NewGUI(material.NewTheme(gofont.Collection()))

	folderEdited = Folder{}
	bookmarkEdited = Bookmark{}
	folderCutClicked = -1
	folderParentCutClicked = -1
	folderPasteClicked = -1
	bookmarkCutClicked = -1
	bookmarkParentCutClicked = -1
	bookmarkPasteClicked = -1
}

func addError(e error) {
	currentError = append(currentError, e.Error())
	gui.Errors = append(gui.Errors, new(widget.Clickable))
}

func displayHistory(gtx layout.Context) D {
	c := []layout.FlexChild{}
	for i := range history {
		b := gui.HistoryListItem[i]
		t := history[i].Title
		id := history[i].Id

		c = append(c, layout.Rigid(func(gtx C) D {
			return gui.InsetTwo.Layout(gtx, func(gtx C) D {
				if b.Clicked() {
					getRemoteNode(id)

					folderMenuClicked = -1
					bookmarkMenuClicked = -1
				}

				ic, _ := widget.NewIcon(icons.AVAVTimer)
				return iconAndTextButton{gui.Theme}.Layout(gtx,
					b,
					ic,
					colorHistory,
					t)
			})
		}))
	}

	return gui.HistoryList.Layout(gtx, c...)
}

func displayError(gtx layout.Context) D {
	if currentError != nil {

		return gui.ErrorList.Layout(gtx, len(currentError), func(gtx layout.Context, i int) D {

			if gui.Errors[i].Clicked() {
				go func() {
					chEDel <- "del"
				}()
			}

			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return textButton{gui.Theme}.Layout(gtx, gui.Errors[i], colorError, currentError[i])
				}),
			)
		})

	}
	return D{}
}

// func displayTest(gtx layout.Context) D {
// 	s := layout.Stack{Alignment: layout.E}
// 	return s.Layout(gtx,
// 		layout.Stacked(
// 			func(gtx C) D {
// 				return material.Label(gui.Theme, unit.Dp(36), "AAAAAA").Layout(gtx)
// 			},
// 		),
// 		layout.Stacked(
// 			func(gtx C) D {
// 				return material.Label(gui.Theme, unit.Dp(36), "BBB").Layout(gtx)
// 			},
// 		),
// 		layout.Stacked(
// 			func(gtx C) D {
// 				return material.Label(gui.Theme, unit.Dp(36), "CCCCCCCC").Layout(gtx)
// 			},
// 		),
// 		layout.Stacked(
// 			func(gtx C) D {
// 				return material.Label(gui.Theme, unit.Dp(36), "DD").Layout(gtx)
// 			},
// 		),
// 	)
// }

func displayTopBar(gtx layout.Context) D {
	return gui.TopBar.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return gui.InsetTwo.Layout(gtx, func(gtx C) D {
				if gui.AddBookmarkBtn.Clicked() {
					w.ReadClipboard()
				}
				ic, _ := widget.NewIcon(icons.ActionBookmark)
				return iconAndTextButton{gui.Theme}.Layout(gtx,
					gui.AddBookmarkBtn,
					ic,
					colorNewBookmark,
					"")
			})
		}),
		layout.Rigid(func(gtx C) D {
			return gui.InsetTwo.Layout(gtx, func(gtx C) D {
				if gui.AddFolderBtn.Clicked() {
					addRemoteFolder(Folder{Title: "new folder", Parent: &Folder{Id: currentFolder.Id}})
				}
				ic, _ := widget.NewIcon(icons.FileCreateNewFolder)
				return iconAndTextButton{gui.Theme}.Layout(gtx,
					gui.AddFolderBtn,
					ic,
					colorNewFolder,
					"")
			})
		}),
		layout.Rigid(func(gtx C) D {
			return gui.InsetTwo.Layout(gtx, func(gtx C) D {
				if folderParentCutClicked != currentFolder.Id && bookmarkParentCutClicked != currentFolder.Id && (folderCutClicked != -1 || bookmarkCutClicked != -1) {
					if gui.PasteBtn.Clicked() {
						if folderCutClicked != -1 {
							f := Folder{Id: folderCutClicked, Parent: &currentFolder}
							//moveFolder(f)
							saveRemoteFolder(f)
							folderCutClicked = -1
						} else {
							b := Bookmark{Id: bookmarkCutClicked, Folder: &currentFolder}
							//moveBookmark(b)
							saveRemoteBookmark(b)
							bookmarkCutClicked = -1
						}
					}

					ic, _ := widget.NewIcon(icons.ContentContentPaste)
					return iconAndTextButton{gui.Theme}.Layout(gtx,
						gui.PasteBtn,
						ic,
						colorPaste,
						"")
				}
				return D{}
			})
		}),
	)
}

func displayBottomBar(gtx layout.Context) D {
	return gui.BottomBar.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return gui.InsetTwo.Layout(gtx, func(gtx C) D {
				if gui.ToogleStarBtn.Clicked() {
					if tooglestarClicked {
						tooglestarClicked = false
					} else {
						tooglestarClicked = true
					}
				}
				ic, _ := widget.NewIcon(icons.ActionStars)
				return iconAndTextButton{gui.Theme}.Layout(gtx,
					gui.ToogleStarBtn,
					ic,
					colorStars,
					"")
			})
		}),
		layout.Rigid(func(gtx C) D {
			return gui.InsetTwo.Layout(gtx, func(gtx C) D {
				if gui.ToogleTagBtn.Clicked() {
					if toogletagClicked {
						toogletagClicked = false
					} else {
						toogletagClicked = true
					}
				}
				ic, _ := widget.NewIcon(icons.ActionLabel)
				return iconAndTextButton{gui.Theme}.Layout(gtx,
					gui.ToogleTagBtn,
					ic,
					colorTags,
					"")
			})
		}),
		layout.Rigid(func(gtx C) D {
			return gui.InsetTwo.Layout(gtx, func(gtx C) D {
				if gui.ToogleMenuBtn.Clicked() {
					if tooglemenuClicked {
						tooglemenuClicked = false
					} else {
						tooglemenuClicked = true
					}
				}
				ic, _ := widget.NewIcon(icons.NavigationMenu)
				return iconAndTextButton{gui.Theme}.Layout(gtx,
					gui.ToogleMenuBtn,
					ic,
					colorToggleMenus,
					"")
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return material.Label(gui.Theme, unit.Dp(12), "").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return gui.InsetTwo.Layout(gtx, func(gtx C) D {
				if gui.ParameterBtn.Clicked() {
					if parameterClicked {
						parameterClicked = false
					} else {
						parameterClicked = true
					}
				}
				gui.ServerURL.SetText(globals.ServerURL)
				gui.ServerUsername.SetText(globals.ServerLogin)
				gui.ServerPassword.SetText(globals.ServerPassword)
				gui.HistorySize.SetText(strconv.Itoa(globals.HistorySize))

				ic, _ := widget.NewIcon(icons.ActionSettings)
				return iconAndTextButton{gui.Theme}.Layout(gtx,
					gui.ParameterBtn,
					ic,
					colorPreferences,
					"")
			})
		}),
	)
}

func displayBreadcrum(gtx layout.Context) D {
	return gui.BreadcrumList.Layout(gtx, len(breadcrum), func(gtx layout.Context, i int) D {
		if gui.BreadcrumListItem[i].Clicked() {
			getRemoteNode(breadcrum[i].Id)

			folderMenuClicked = -1
		}
		return textButton{gui.Theme}.Layout(gtx,
			gui.BreadcrumListItem[i],
			colorBreadcrum,
			breadcrum[i].Title)
	})
}

func displayBookmarkMenu(gtx layout.Context, b *Bookmark) D {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			if gui.StarBtn.Clicked() {
				starBookmark(*b, !b.Starred)
			}

			var ic *widget.Icon
			if b.Starred {
				ic, _ = widget.NewIcon(icons.ToggleStar)
			} else {
				ic, _ = widget.NewIcon(icons.ToggleStarBorder)
			}

			return iconAndTextButton{gui.Theme}.Layout(gtx,
				gui.StarBtn,
				ic,
				colorBookmark,
				"")
		}),
		layout.Rigid(func(gtx C) D {
			if gui.CutBtn.Clicked() {
				bookmarkCutClicked = b.Id
				bookmarkParentCutClicked = b.Folder.Id
				folderCutClicked = -1
				folderParentCutClicked = -1
				bookmarkMenuClicked = -1

				gui.PasteBtn = new(widget.Clickable)
			}
			if gui.EditBtn.Clicked() {
				bookmarkEdited = *b
				gui.BookmarkName.SetText(b.Title)
				gui.BookmarkURL.SetText(b.URL)
				for i := range b.Tags {
					gui.EditTags[fmt.Sprintf("%d:%s", b.Tags[i].Id, b.Tags[i].Name)].Value = true
				}
				bookmarkMenuClicked = -1
			}

			ic, _ := widget.NewIcon(icons.ContentContentCut)
			return iconAndTextButton{gui.Theme}.Layout(gtx,
				gui.CutBtn,
				ic,
				colorBookmark,
				"")
		}),
		layout.Rigid(func(gtx C) D {
			if gui.DeleteBookmarkBtn.Clicked() {
				if deleteConfirmed {
					deleteRemoteBookmark(*b)
					deleteConfirmed = false
				} else {
					deleteConfirmed = true
				}
			}

			var ic *widget.Icon
			if deleteConfirmed {
				ic, _ = widget.NewIcon(icons.ActionCheckCircle)
			} else {
				ic, _ = widget.NewIcon(icons.ActionDelete)
			}
			return iconAndTextButton{gui.Theme}.Layout(gtx,
				gui.DeleteBookmarkBtn,
				ic,
				colorBookmark,
				"")
		}),
		layout.Rigid(func(gtx C) D {
			ic, _ := widget.NewIcon(icons.ContentCreate)
			return iconAndTextButton{gui.Theme}.Layout(gtx,
				gui.EditBtn,
				ic,
				colorBookmark,
				"")
		}),
	)
}

func displayFolderMenu(gtx layout.Context, f *Folder) D {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			if gui.CutBtn.Clicked() {
				folderCutClicked = f.Id
				folderParentCutClicked = f.Parent.Id
				bookmarkCutClicked = -1
				bookmarkParentCutClicked = -1
				folderMenuClicked = -1

				gui.PasteBtn = new(widget.Clickable)
			}
			if gui.PasteBtn.Clicked() {
				folderPasteClicked = f.Id
				folderCutClicked = -1
				folderParentCutClicked = -1
				bookmarkCutClicked = -1
				bookmarkParentCutClicked = -1
				folderMenuClicked = -1
			}
			if gui.EditBtn.Clicked() {
				folderEdited = *f
				gui.FolderName.SetText(f.Title)
				folderMenuClicked = -1
			}
			ic, _ := widget.NewIcon(icons.ContentContentCut)
			return iconAndTextButton{gui.Theme}.Layout(gtx,
				gui.CutBtn,
				ic,
				colorFolder,
				"")
		}),
		layout.Rigid(func(gtx C) D {
			if gui.DeleteFolderBtn.Clicked() {
				if deleteConfirmed {
					deleteRemoteFolder(f.Id)
					deleteConfirmed = false
				} else {
					deleteConfirmed = true
				}
			}

			var ic *widget.Icon
			if deleteConfirmed {
				ic, _ = widget.NewIcon(icons.ActionCheckCircle)
			} else {
				ic, _ = widget.NewIcon(icons.ActionDelete)
			}
			return iconAndTextButton{gui.Theme}.Layout(gtx,
				gui.DeleteFolderBtn,
				ic,
				colorFolder,
				"")
		}),
		layout.Rigid(func(gtx C) D {
			ic, _ := widget.NewIcon(icons.ContentCreate)
			return iconAndTextButton{gui.Theme}.Layout(gtx,
				gui.EditBtn,
				ic,
				colorFolder,
				"")
		}),
	)
}

func displayCurrentFolderSubfolders(gtx layout.Context, i int) D {
	return gui.InsetTwo.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			// folder button
			layout.Flexed(1, func(gtx C) D {
				if gui.FolderListItem[currentFolder.Folders[i].Title].Clicked() {
					getRemoteNode(currentFolder.Folders[i].Id)

					folderMenuClicked = -1
					bookmarkMenuClicked = -1

					isInHistory := false
					for _, f := range history {
						if f.Id == currentFolder.Folders[i].Id {
							isInHistory = true
						}
					}

					if !isInHistory {
						if len(history) > 0 && len(history) >= globals.HistorySize {
							history = history[1:]
						}
						history = append(history, *currentFolder.Folders[i])
						gui.HistoryListItem = append(gui.HistoryListItem, new(widget.Clickable))
					}
				}

				var c color.NRGBA
				if folderCutClicked == currentFolder.Folders[i].Id {
					c = colorCut
				} else {
					c = colorFolder
				}
				return textButton{gui.Theme}.Layout(gtx,
					gui.FolderListItem[currentFolder.Folders[i].Title],
					c,
					currentFolder.Folders[i].Title)
			}),
			// folder menu
			layout.Rigid(func(gtx C) D {
				if tooglemenuClicked {
					if gui.FolderMenu[currentFolder.Folders[i].Title].Clicked() {
						folderMenuClicked = currentFolder.Folders[i].Id
						gui.CutBtn = new(widget.Clickable)
						gui.PasteBtn = new(widget.Clickable)
						gui.EditBtn = new(widget.Clickable)

						bookmarkMenuClicked = -1
					}
					if folderMenuClicked == currentFolder.Folders[i].Id {
						return displayFolderMenu(gtx, currentFolder.Folders[i])
					}
					ic, _ := widget.NewIcon(icons.NavigationMenu)
					return iconAndTextButton{gui.Theme}.Layout(gtx,
						gui.FolderMenu[currentFolder.Folders[i].Title],
						ic,
						colorFolder,
						"")
				}
				return D{}
			}),
		)
	})
}

func displayCurrentFolderBookmarks(gtx layout.Context, i int) D {
	return gui.InsetTwo.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			// bookmark button
			layout.Flexed(1, func(gtx C) D {

				if gui.BookmarkListItem[currentFolder.Bookmarks[i].Title].Clicked() {
					go func() {
						err := openURL(currentFolder.Bookmarks[i].URL)
						if err != nil {
							addError(err)
						}
					}()
				}

				var c color.NRGBA
				if bookmarkCutClicked == currentFolder.Bookmarks[i].Id {
					c = colorCut
				} else {
					c = colorBookmark
				}

				var favicon image.Image
				var faviconOp paint.ImageOp
				var err error
				img := currentFolder.Bookmarks[i].Favicon
				if len(img) == 0 {
					img = globals.NO_ICON
				} else {
					img = img[22:]
				}
				ior := base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(img)))
				favicon, err = png.Decode(ior)
				if err != nil {
					panic(err)
				}
				faviconOp = paint.NewImageOp(favicon)

				var tags []string
				for _, t := range currentFolder.Bookmarks[i].Tags {
					tags = append(tags, t.Name)
				}

				return imageAndTextAndTagsButton{gui.Theme}.Layout(gtx,
					gui.BookmarkListItem[currentFolder.Bookmarks[i].Title],
					&faviconOp,
					c,
					currentFolder.Bookmarks[i].Title,
					tags)
			}),
			// bookmark menu
			layout.Rigid(func(gtx C) D {
				if tooglemenuClicked {
					if gui.BookmarkMenu[currentFolder.Bookmarks[i].Title].Clicked() {
						bookmarkMenuClicked = currentFolder.Bookmarks[i].Id
						gui.CutBtn = new(widget.Clickable)
						gui.EditBtn = new(widget.Clickable)
						gui.PasteBtn = nil
						gui.StarBtn = new(widget.Clickable)

						folderMenuClicked = -1
					}
					if bookmarkMenuClicked == currentFolder.Bookmarks[i].Id {
						return displayBookmarkMenu(gtx, currentFolder.Bookmarks[i])
					}
					ic, _ := widget.NewIcon(icons.NavigationMenu)
					return iconAndTextButton{gui.Theme}.Layout(gtx,
						gui.BookmarkMenu[currentFolder.Bookmarks[i].Title],
						ic,
						colorBookmark,
						"")
				}
				return D{}
			}),
		)
	})
}

func displayCurrentFolder(gtx layout.Context) D {
	return gui.MainContent.Layout(gtx,
		layout.Rigid(func(gtx C) D {

			nbF := len(currentFolder.Folders)
			nbB := len(currentFolder.Bookmarks)
			nbNodes := nbF + nbB

			return gui.CurrentFolderList.Layout(gtx, nbNodes, func(gtx layout.Context, i int) D {

				if i < nbF {
					return displayCurrentFolderSubfolders(gtx, i)
				} else {
					return displayCurrentFolderBookmarks(gtx, i-nbF)
				}

			},
			)
		},
		))
}

func displaySearch(gtx layout.Context) D {
	return gui.SearchBar.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return gui.InsetTen.Layout(gtx, func(gtx C) D {
				return material.Editor(gui.Theme, gui.SearchForm, "search").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			for _, e := range gui.SearchForm.Events() {
				if e, ok := e.(widget.SubmitEvent); ok {
					doSearch(e.Text)
				}
			}

			if gui.SearchButton.Clicked() {
				doSearch(gui.SearchForm.Text())
			}

			ic, _ := widget.NewIcon(icons.ActionSearch)
			return iconAndTextButton{gui.Theme}.Layout(gtx,
				gui.SearchButton,
				ic,
				colorSearch,
				"")
		}),
	)
}

func displayTags(tags []Tag, gtx layout.Context) D {
	if toogletagClicked {
		return gui.TagList.Layout(gtx, len(tags), func(gtx layout.Context, i int) D {
			return gui.InsetTwo.Layout(gtx, func(gtx C) D {
				if gui.TagListItem[tags[i].Name].Clicked() {
					doSearch(tags[i].Name)
				}
				ic, _ := widget.NewIcon(icons.ActionLabel)
				return iconAndTextButton{gui.Theme}.Layout(gtx,
					gui.TagListItem[tags[i].Name],
					ic,
					colorTags,
					tags[i].Name)
			})
		})
	}
	return D{}
}

func displayStars(stars []Bookmark, gtx layout.Context) D {
	if tooglestarClicked {

		return gui.StarList.Layout(gtx, len(stars), func(gtx layout.Context, i int) D {
			return gui.InsetTwo.Layout(gtx, func(gtx C) D {

				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx C) D {

						if gui.StarListItem[stars[i].Title].Clicked() {
							go func() {
								err := openURL(stars[i].URL)
								if err != nil {
									addError(err)
								}
							}()
						}
						ic, _ := widget.NewIcon(icons.ToggleStar)
						return iconAndTextButton{gui.Theme}.Layout(gtx,
							gui.StarListItem[stars[i].Title],
							ic,
							colorStars,
							stars[i].Title)

					}),
					layout.Rigid(func(gtx C) D {
						if gui.StarLocateInTree[stars[i].Title].Clicked() {
							getRemoteNode(stars[i].Folder.Id)
						}

						ic, _ := widget.NewIcon(icons.FileFolderOpen)
						return iconAndTextButton{gui.Theme}.Layout(gtx,
							gui.StarLocateInTree[stars[i].Title],
							ic,
							colorStars,
							"")
					}),
				)
			})
		})

	}
	return D{}
}

func displayEditFolder(gtx layout.Context) D {
	return gui.MainPage.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if gui.EditCancelBtn.Clicked() {
						folderEdited = Folder{}
						gui.FolderName.SetText("")
					}

					ic, _ := widget.NewIcon(icons.NavigationCancel)
					return gui.InsetTwo.Layout(gtx, func(gtx C) D {
						return iconAndTextButton{gui.Theme}.Layout(gtx,
							gui.EditCancelBtn,
							ic,
							colorFolder,
							"")
					})
				}),
				layout.Rigid(func(gtx C) D {
					if gui.EditSaveBtn.Clicked() {
						// we need to keep folderEditClicked untouched
						f := folderEdited
						f.Title = gui.FolderName.Text()
						f.Parent = nil
						saveRemoteFolder(f)
					}

					ic, _ := widget.NewIcon(icons.ContentSave)
					return gui.InsetTwo.Layout(gtx, func(gtx C) D {
						return iconAndTextButton{gui.Theme}.Layout(gtx,
							gui.EditSaveBtn,
							ic,
							colorFolder,
							"")
					})
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return gui.InsetTen.Layout(gtx, func(gtx C) D {
				return material.Editor(gui.Theme, gui.FolderName, "folder name").Layout(gtx)
			})
		}),
	)
}

func displayEditBookmark(gtx layout.Context) D {
	return gui.MainPage.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if gui.EditCancelBtn.Clicked() {
						bookmarkEdited = Bookmark{}
						gui.BookmarkName.SetText("")
						gui.BookmarkURL.SetText("")
						for i := range gui.EditTags {
							gui.EditTags[i].Value = false
						}
					}

					ic, _ := widget.NewIcon(icons.NavigationCancel)
					return gui.InsetTwo.Layout(gtx, func(gtx C) D {
						return iconAndTextButton{gui.Theme}.Layout(gtx,
							gui.EditCancelBtn,
							ic,
							colorBookmark,
							"")
					})
				}),
				layout.Rigid(func(gtx C) D {
					if gui.EditSaveBtn.Clicked() {
						// we need to keep bookmarkEditClicked untouched
						b := bookmarkEdited
						b.Title = gui.BookmarkName.Text()
						b.URL = gui.BookmarkURL.Text()
						b.Folder = nil
						// adding new tags
						for k, guit := range gui.EditTags {
							splitk := strings.Split(k, ":")
							tagId := splitk[0]
							tagName := splitk[1]
							if guit.Value {
								found := false
								for _, bt := range b.Tags {
									if bt.Name == tagName {
										found = true
										break
									}
								}
								if !found {
									var (
										id  int
										err error
									)
									if id, err = strconv.Atoi(tagId); err != nil {
										panic(err)
									}
									b.Tags = append(b.Tags, &Tag{Name: tagName, Id: id})
								}
							}
						}
						// removing former tags
						n := 0
						for _, bt := range b.Tags {
							found := false
							for k, t := range gui.EditTags {
								if t.Value {
									splitk := strings.Split(k, ":")
									//tagId := splitk[0]
									tagName := splitk[1]
									if tagName == bt.Name {
										found = true
										break
									}
								}
							}
							if found {
								b.Tags[n] = bt
								n++
							}
						}
						b.Tags = b.Tags[:n]
						saveRemoteBookmark(b)
					}

					ic, _ := widget.NewIcon(icons.ContentSave)
					return gui.InsetTwo.Layout(gtx, func(gtx C) D {
						return iconAndTextButton{gui.Theme}.Layout(gtx,
							gui.EditSaveBtn,
							ic,
							colorBookmark,
							"")
					})
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return gui.InsetTen.Layout(gtx, func(gtx C) D {
				return material.Editor(gui.Theme, gui.BookmarkName, "bookmark name").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return gui.InsetTen.Layout(gtx, func(gtx C) D {
				return material.Editor(gui.Theme, gui.BookmarkURL, "bookmark URL").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return gui.InsetTen.Layout(gtx, func(gtx C) D {
				list := &layout.List{
					Axis:      layout.Vertical,
					Alignment: layout.Middle,
				}
				return list.Layout(gtx, len(tags), func(gtx layout.Context, i int) D {
					return material.CheckBox(gui.Theme, gui.EditTags[fmt.Sprintf("%d:%s", tags[i].Id, tags[i].Name)], tags[i].Name).Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return gui.InsetTen.Layout(gtx, func(gtx C) D {
						return material.Editor(gui.Theme, gui.TagName, "add tag").Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					if gui.AddTag.Clicked() {
						gui.EditTags[fmt.Sprintf("%d:%s", -1, gui.TagName.Text())] = new(widget.Bool)
						tags = append(tags, Tag{Name: gui.TagName.Text(), Id: -1})

						gui.TagName.SetText("")
					}

					ic, _ := widget.NewIcon(icons.ContentAdd)
					return iconAndTextButton{gui.Theme}.Layout(gtx,
						gui.AddTag,
						ic,
						colorTags,
						"")
				}),
			)
		}),
	)
}

func displayEditParameters(gtx layout.Context) D {
	return gui.MainPage.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if gui.ParameterCancelBtn.Clicked() {
						parameterClicked = false
					}

					ic, _ := widget.NewIcon(icons.NavigationCancel)
					return gui.InsetTwo.Layout(gtx, func(gtx C) D {
						return iconAndTextButton{gui.Theme}.Layout(gtx,
							gui.ParameterCancelBtn,
							ic,
							colorPreferences,
							"")
					})
				}),
				layout.Rigid(func(gtx C) D {
					if gui.ParameterSaveBtn.Clicked() {
						parameterClicked = false

						go func() {
							var (
								hs  int
								err error
							)
							if hs, err = strconv.Atoi(gui.HistorySize.Text()); err != nil {
								addError(err)
							}
							p := datastores.Preferences{
								ServerURL:      gui.ServerURL.Text(),
								ServerUsername: gui.ServerUsername.Text(),
								ServerPassword: gui.ServerPassword.Text(),
								HistorySize:    hs,
							}

							globals.ServerURL = gui.ServerURL.Text()
							globals.ServerLogin = gui.ServerUsername.Text()
							globals.ServerPassword = gui.ServerPassword.Text()
							globals.HistorySize = hs

							globals.B64Auth = basicAuth(string(globals.ServerLogin), string(globals.ServerPassword))

							if err = dstore.SavePreferences(p); err != nil {
								addError(err)
							}

							getRemoteStars()
							getRemoteTags()
							getRemoteNode(1)
						}()
					}

					ic, _ := widget.NewIcon(icons.ContentSave)
					return gui.InsetTwo.Layout(gtx, func(gtx C) D {
						return iconAndTextButton{gui.Theme}.Layout(gtx,
							gui.ParameterSaveBtn,
							ic,
							colorPreferences,
							"")
					})
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				// layout.Rigid(func() {
				// 	//server URL label
				// 	material.Label(gui.Theme, unit.Dp(24), "server URL").Layout(gtx)
				// }),
				layout.Flexed(1, func(gtx C) D {
					//server URL
					return gui.InsetTen.Layout(gtx, func(gtx C) D {
						return material.Editor(gui.Theme, gui.ServerURL, "server URL").Layout(gtx)
					})
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return gui.InsetTen.Layout(gtx, func(gtx C) D {
				return material.Editor(gui.Theme, gui.ServerUsername, "username").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return gui.InsetTen.Layout(gtx, func(gtx C) D {
				return material.Editor(gui.Theme, gui.ServerPassword, "password").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return gui.InsetTen.Layout(gtx, func(gtx C) D {
				return material.Editor(gui.Theme, gui.HistorySize, "history size").Layout(gtx)
			})
		}),
	)
}

func displaySearchResults(gtx layout.Context, bs []Bookmark) D {
	return gui.MainPage.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if gui.EditCancelBtn.Clicked() {
						searchResults = nil
						gui.SearchForm.SetText("")
					}

					ic, _ := widget.NewIcon(icons.NavigationCancel)
					return gui.InsetTwo.Layout(gtx, func(gtx C) D {
						return iconAndTextButton{gui.Theme}.Layout(gtx,
							gui.EditCancelBtn,
							ic,
							colorBookmark,
							"")
					})
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return gui.SearchResultList.Layout(gtx, len(bs), func(gtx layout.Context, i int) D {
				return gui.InsetTwo.Layout(gtx, func(gtx C) D {
					ic, _ := widget.NewIcon(icons.ActionSearch)
					return iconAndTextButton{gui.Theme}.Layout(gtx,
						gui.BookmarkListItem[bs[i].Title],
						ic,
						colorBookmark,
						bs[i].Title)
				})
			})
		}),
	)

}

func main() {

	go func() {
		getRemoteStars()
		getRemoteTags()
		getRemoteNode(1)
	}()

	go func() {
		w = app.NewWindow(
			app.Title("GoBkm"),
		)

		sysinset := &layout.Inset{}
		resetSysinset := func(x system.Insets) {
			sysinset.Top = x.Top
			sysinset.Bottom = x.Bottom
			sysinset.Left = x.Left
			sysinset.Right = x.Right
		}

		var ops op.Ops
		for {
			select {
			case <-chEDel:
				currentError = nil
				gui.Errors = nil
			case <-chBStar:
				getRemoteStars()
				getRemoteNode(currentFolder.Id)
			case <-chFMove:
				getRemoteNode(currentFolder.Id)
			case <-chBMove:
				getRemoteNode(currentFolder.Id)
			case <-chFDel:
				// reloading the view on new folder
				// note that it may be useless as we may
				// be in another folder
				getRemoteNode(currentFolder.Id)
			case k := <-chBDel:
				// cleaning the gui bookmarks elements for the old bookmark
				delete(gui.BookmarkMenu, k)
				delete(gui.BookmarkListItem, k)

				for i, bo := range currentFolder.Bookmarks {
					if bo.Title == k {
						currentFolder.Bookmarks = append(currentFolder.Bookmarks[:i], currentFolder.Bookmarks[i+1:]...)
						break
					}
				}

				// reloading the view on new bookmark
				// note that it may be useless as we may
				// be in another folder
				getRemoteNode(currentFolder.Id)
			case <-chFAdd:
				// reloading the view on new folder
				// note that it may be useless as we may
				// be in another folder
				getRemoteNode(currentFolder.Id)
			case <-chBAdd:
				// reloading the view on new bookmark
				// note that it may be useless as we may
				// be in another folder
				getRemoteNode(currentFolder.Id)
			case b := <-chBEdit:
				// A bookmark has been modified.
				// - b contains the old bookmark.
				// - gui.BookmarkName contains the new bookmark name.

				// cleaning the gui bookmarks elements for the old bookmark
				delete(gui.BookmarkMenu, b.Title)
				delete(gui.BookmarkListItem, b.Title)

				// creating widgets for the new bookmark
				gui.BookmarkMenu[gui.BookmarkName.Text()] = new(widget.Clickable)
				gui.BookmarkListItem[gui.BookmarkName.Text()] = new(widget.Clickable)

				// updating layout
				for i, bo := range currentFolder.Bookmarks {
					if bo.Title == b.Title {
						currentFolder.Bookmarks[i].Title = gui.BookmarkName.Text()
						currentFolder.Bookmarks[i].URL = gui.BookmarkURL.Text()
					}
				}

				// disabling editing mode
				bookmarkEdited = Bookmark{}
				gui.BookmarkName.SetText("")
				gui.BookmarkURL.SetText("")
				for i := range gui.EditTags {
					gui.EditTags[i].Value = false
				}

				// lazily reloading the tags to update the gui variables
				getRemoteTags()
				// lazily reloading the current node to update the gui variables
				getRemoteNode(currentFolder.Id)
			case f := <-chFEdit:
				// A folder has been modified.
				// - f contains the old folder
				// - gui.FolderName contains the new folder name

				// cleaning the gui folders elements for the old folder
				delete(gui.FolderMenu, f.Title)
				delete(gui.FolderListItem, f.Title)

				// creating widgets for the new folder
				gui.FolderMenu[gui.FolderName.Text()] = new(widget.Clickable)
				gui.FolderListItem[gui.FolderName.Text()] = new(widget.Clickable)

				// updating layout
				for i, fo := range currentFolder.Folders {
					if fo.Title == f.Title {
						currentFolder.Folders[i].Title = gui.FolderName.Text()
					}
				}

				// disabling editing mode
				folderEdited = Folder{}
				gui.FolderName.SetText("")
			case b := <-chS:
				stars = b
				for i := range stars {
					gui.StarListItem[stars[i].Title] = new(widget.Clickable)
					gui.StarLocateInTree[stars[i].Title] = new(widget.Clickable)
				}
				displayStars(stars, gtx)
			case t := <-chT:
				tags = t
				for i := range t {
					gui.EditTags[fmt.Sprintf("%d:%s", tags[i].Id, tags[i].Name)] = new(widget.Bool)
					gui.TagListItem[tags[i].Name] = new(widget.Clickable)
				}
			case f := <-chF:
				currentFolder = f

				for i := range currentFolder.Bookmarks {
					gui.BookmarkListItem[currentFolder.Bookmarks[i].Title] = new(widget.Clickable)
					gui.BookmarkMenu[currentFolder.Bookmarks[i].Title] = new(widget.Clickable)
				}
				for i := range currentFolder.Folders {
					gui.FolderListItem[currentFolder.Folders[i].Title] = new(widget.Clickable)
					gui.FolderMenu[currentFolder.Folders[i].Title] = new(widget.Clickable)
				}

				breadcrum = nil
				gui.BreadcrumListItem = nil
				breadcrum = append(breadcrum, currentFolder)
				gui.BreadcrumListItem = append(gui.BreadcrumListItem, new(widget.Clickable))
				p := currentFolder.Parent
				for p != nil {
					breadcrum = append([]Folder{*p}, breadcrum...)
					gui.BreadcrumListItem = append(gui.BreadcrumListItem, new(widget.Clickable))
					p = p.Parent
				}
			case e := <-w.Events():
				switch e := e.(type) {
				case clipboard.Event:
					if !isValidUrl(e.Text) {
						addError(errors.New("not a valid URL"))
					} else {
						addRemoteBookmark(Bookmark{URL: e.Text, Title: e.Text, Folder: &Folder{Id: currentFolder.Id}})
					}
				case system.DestroyEvent:
				case system.FrameEvent:
					//gtx.Reset(e.Queue, e.Config, e.Size)
					gtx = layout.NewContext(&ops, e)

					if len(currentError) != 0 {
						for range currentError {
							gui.Errors = append(gui.Errors, new(widget.Clickable))
						}
					}

					if len(searchResults) != 0 {
						// TODO: improve comparison
						for i := range searchResults {
							gui.BookmarkListItem[searchResults[i].Title] = new(widget.Clickable)
						}

						gui.MainWidgets = []layout.FlexChild{
							layout.Rigid(
								func(gtx C) D {
									// Edit parameters form
									return displaySearchResults(gtx, searchResults)
								}),
						}
					} else if bookmarkEdited.Id != 0 {
						gui.MainWidgets = []layout.FlexChild{
							layout.Rigid(
								func(gtx C) D {
									// Edit bookmark form
									return displayEditBookmark(gtx)
								}),
						}
					} else if folderEdited.Id != 0 {
						gui.MainWidgets = []layout.FlexChild{
							layout.Rigid(
								func(gtx C) D {
									// Edit folder form
									return displayEditFolder(gtx)
								}),
						}
					} else if parameterClicked {
						gui.MainWidgets = []layout.FlexChild{
							layout.Rigid(
								func(gtx C) D {
									// Edit parameters form
									return displayEditParameters(gtx)
								}),
						}
					} else {
						gui.MainWidgets = []layout.FlexChild{
							// layout.Rigid(
							// 	func(gtx C) D {
							// 		// Test
							// 		return displayTest(gtx)
							// 	}),
							layout.Rigid(
								func(gtx C) D {
									// Errors
									return displayError(gtx)
								}),
							layout.Rigid(
								func(gtx C) D {
									// Toogle stars button
									return displayTopBar(gtx)
								}),
							layout.Rigid(
								func(gtx C) D {
									// Stars list
									return displayStars(stars, gtx)
								}),
							layout.Rigid(
								func(gtx C) D {
									// Search form
									return displaySearch(gtx)
								}),
							layout.Rigid(
								func(gtx C) D {
									// Breadcrum
									return displayBreadcrum(gtx)
								}),
							layout.Flexed(1,
								func(gtx C) D {
									// Bookmark list
									return displayCurrentFolder(gtx)
								}),
							layout.Flexed(0, func(gtx C) D {
								return material.Label(gui.Theme, unit.Dp(12), "").Layout(gtx)
							}),
							layout.Rigid(
								func(gtx C) D {
									// Folder history
									return displayHistory(gtx)
								}),
							layout.Rigid(
								func(gtx C) D {
									// Tags list
									return displayTags(tags, gtx)
								}),
							layout.Rigid(
								func(gtx C) D {
									// Bottom bar
									return displayBottomBar(gtx)
								}),
						}
					}

					resetSysinset(e.Insets)
					sysinset.Layout(gtx, func(gtx C) D {
						return gui.MainPage.Layout(gtx, gui.MainWidgets...)
					})
					e.Frame(&ops)
				}
			}
		}
	}()
	app.Main()

}
