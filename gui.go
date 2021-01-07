package main

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type GUI struct {
	Theme *material.Theme

	MainPage                                  *layout.Flex
	MainWidgets                               []layout.FlexChild
	BottomBar, TopBar, SearchBar, MainContent *layout.Flex
	CurrentFolderList                         *layout.List

	InsetTen layout.Inset
	InsetTwo layout.Inset

	ToogleStarBtn                      *widget.Clickable
	ToogleTagBtn                       *widget.Clickable
	ToogleMenuBtn                      *widget.Clickable
	ParameterBtn                       *widget.Clickable
	ParameterSaveBtn                   *widget.Clickable
	ParameterCancelBtn                 *widget.Clickable
	CutBtn, EditBtn, StarBtn, PasteBtn *widget.Clickable
	EditCancelBtn, EditSaveBtn         *widget.Clickable
	AddBookmarkBtn                     *widget.Clickable
	AddFolderBtn                       *widget.Clickable
	DeleteBookmarkBtn                  *widget.Clickable
	DeleteFolderBtn                    *widget.Clickable

	EditTags         map[string]*widget.Bool
	SearchResultList *layout.List

	TagList          *layout.List
	StarList         *layout.List
	BookmarkList     *layout.List
	FolderList       *layout.List
	TagListItem      map[string]*widget.Clickable
	StarListItem     map[string]*widget.Clickable
	BookmarkListItem map[string]*widget.Clickable
	FolderListItem   map[string]*widget.Clickable

	FolderMenu       map[string]*widget.Clickable
	BookmarkMenu     map[string]*widget.Clickable
	StarLocateInTree map[string]*widget.Clickable

	BreadcrumList     *layout.List
	HistoryList       *layout.Flex
	HistoryListItem   []*widget.Clickable
	BreadcrumListItem []*widget.Clickable

	ServerURL      *widget.Editor
	ServerUsername *widget.Editor
	ServerPassword *widget.Editor

	HistorySize *widget.Editor

	SearchButton *widget.Clickable
	SearchForm   *widget.Editor

	FolderName,
	BookmarkName,
	BookmarkURL *widget.Editor
	TagName *widget.Editor
	AddTag  *widget.Clickable

	Errors    []*widget.Clickable
	ErrorList *layout.List
}

func NewGUI(th *material.Theme) GUI {
	stars := map[string]*widget.Clickable{}
	tags := map[string]*widget.Clickable{}
	edittags := map[string]*widget.Bool{}
	currentnodeb := map[string]*widget.Clickable{}
	currentnodef := map[string]*widget.Clickable{}
	currentnodefmenu := map[string]*widget.Clickable{}
	currentnodebmenu := map[string]*widget.Clickable{}
	starlocateintree := map[string]*widget.Clickable{}
	return GUI{
		AddBookmarkBtn:     new(widget.Clickable),
		AddFolderBtn:       new(widget.Clickable),
		AddTag:             new(widget.Clickable),
		BookmarkList:       &layout.List{Axis: layout.Vertical},
		BookmarkListItem:   currentnodeb,
		BookmarkMenu:       currentnodebmenu,
		BookmarkName:       &widget.Editor{SingleLine: true},
		BookmarkURL:        &widget.Editor{SingleLine: true},
		BottomBar:          &layout.Flex{Axis: layout.Horizontal},
		BreadcrumList:      &layout.List{Axis: layout.Horizontal},
		CurrentFolderList:  &layout.List{Axis: layout.Vertical},
		DeleteBookmarkBtn:  new(widget.Clickable),
		DeleteFolderBtn:    new(widget.Clickable),
		EditCancelBtn:      new(widget.Clickable),
		EditSaveBtn:        new(widget.Clickable),
		EditTags:           edittags,
		ErrorList:          &layout.List{Axis: layout.Vertical},
		FolderList:         &layout.List{Axis: layout.Vertical},
		FolderListItem:     currentnodef,
		FolderMenu:         currentnodefmenu,
		FolderName:         &widget.Editor{SingleLine: true},
		HistoryList:        &layout.Flex{Axis: layout.Vertical},
		HistorySize:        &widget.Editor{SingleLine: true},
		InsetTen:           layout.UniformInset(unit.Dp(10)),
		InsetTwo:           layout.UniformInset(unit.Dp(2)),
		MainContent:        &layout.Flex{Axis: layout.Vertical},
		MainPage:           &layout.Flex{Axis: layout.Vertical},
		ParameterBtn:       new(widget.Clickable),
		ParameterCancelBtn: new(widget.Clickable),
		ParameterSaveBtn:   new(widget.Clickable),
		SearchBar:          &layout.Flex{Axis: layout.Horizontal},
		SearchButton:       new(widget.Clickable),
		SearchForm:         &widget.Editor{SingleLine: true, Submit: true},
		SearchResultList:   &layout.List{Axis: layout.Vertical},
		ServerPassword:     &widget.Editor{SingleLine: true, Mask: '*'},
		ServerURL:          &widget.Editor{SingleLine: true},
		ServerUsername:     &widget.Editor{SingleLine: true},
		StarList:           &layout.List{Axis: layout.Vertical},
		StarListItem:       stars,
		StarLocateInTree:   starlocateintree,
		TagList:            &layout.List{Axis: layout.Vertical},
		TagListItem:        tags,
		TagName:            &widget.Editor{SingleLine: true},
		Theme:              th,
		ToogleMenuBtn:      new(widget.Clickable),
		ToogleStarBtn:      new(widget.Clickable),
		ToogleTagBtn:       new(widget.Clickable),
		TopBar:             &layout.Flex{Axis: layout.Horizontal},
	}
}
