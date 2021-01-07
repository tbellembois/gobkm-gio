package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/tbellembois/gobkm-gio/globals"

	. "github.com/tbellembois/gobkm/types"
)

func SearchBookmark(s string) ([]Bookmark, error) {

	var bs []Bookmark

	var client http.Client
	req, err := http.NewRequest("GET", globals.ServerURL+"/searchBookmarks/?search="+s, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Basic "+globals.B64Auth)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &bs)
	if err != nil {
		panic(err)
	}

	return bs, nil
}

func StarBookmark(b Bookmark, s bool) error {

	body, err := json.Marshal(b)
	if err != nil {
		return err
	}

	var client http.Client
	req, err := http.NewRequest("POST", globals.ServerURL+"/starBookmark/?id=-"+strconv.Itoa(b.Id)+"&star="+strconv.FormatBool(s), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Basic "+globals.B64Auth)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	return nil
}

func AddBookmark(b Bookmark) (Bookmark, error) {

	var newB Bookmark

	body, err := json.Marshal(b)
	if err != nil {
		return Bookmark{}, err
	}

	var client http.Client
	req, err := http.NewRequest("POST", globals.ServerURL+"/addBookmark/", bytes.NewBuffer(body))
	if err != nil {
		return Bookmark{}, err
	}

	req.Header.Add("Authorization", "Basic "+globals.B64Auth)
	resp, err := client.Do(req)
	if err != nil {
		return Bookmark{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Bookmark{}, fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(rbody, &newB)
	if err != nil {
		panic(err)
	}

	return newB, nil
}

func AddFolder(f Folder) (Folder, error) {

	var newF Folder

	body, err := json.Marshal(f)
	if err != nil {
		return Folder{}, err
	}

	var client http.Client
	req, err := http.NewRequest("POST", globals.ServerURL+"/addFolder/", bytes.NewBuffer(body))
	if err != nil {
		return Folder{}, err
	}

	req.Header.Add("Authorization", "Basic "+globals.B64Auth)
	resp, err := client.Do(req)
	if err != nil {
		return Folder{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Folder{}, fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(rbody, &newF)
	if err != nil {
		panic(err)
	}

	return newF, nil
}

func DeleteFolder(id int) error {

	var client http.Client
	req, err := http.NewRequest("POST", globals.ServerURL+"/deleteFolder/?id="+strconv.Itoa(id), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Basic "+globals.B64Auth)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	return nil
}

func DeleteBookmark(id int) error {

	var client http.Client
	req, err := http.NewRequest("POST", globals.ServerURL+"/deleteBookmark/?id=-"+strconv.Itoa(id), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Basic "+globals.B64Auth)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	return nil
}

func SaveBookmark(b Bookmark) error {

	b.Id = -b.Id

	body, err := json.Marshal(b)
	if err != nil {
		return err
	}

	var client http.Client
	req, err := http.NewRequest("POST", globals.ServerURL+"/updateBookmark/", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Basic "+globals.B64Auth)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	return nil
}

func SaveFolder(f Folder) error {

	body, err := json.Marshal(f)
	if err != nil {
		return err
	}

	var client http.Client
	req, err := http.NewRequest("POST", globals.ServerURL+"/updateFolder/", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Basic "+globals.B64Auth)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	return nil
}

func GetStars() ([]Bookmark, error) {

	var stars []Bookmark

	var client http.Client
	req, err := http.NewRequest("GET", globals.ServerURL+"/getStars/", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Basic "+globals.B64Auth)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &stars)
	if err != nil {
		return nil, err
	}

	return stars, nil

}

func GetTags() ([]Tag, error) {

	var tags []Tag

	var client http.Client
	req, err := http.NewRequest("GET", globals.ServerURL+"/getTags/", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Basic "+globals.B64Auth)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil

}

func GetNode(folderId int) (Folder, error) {

	var tree Folder

	var client http.Client
	req, err := http.NewRequest("GET", globals.ServerURL+"/getFolderChildren/?id="+strconv.Itoa(folderId), nil)
	if err != nil {
		return Folder{}, err
	}

	req.Header.Add("Authorization", "Basic "+globals.B64Auth)
	resp, err := client.Do(req)
	if err != nil {
		return Folder{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Folder{}, err
	}

	err = json.Unmarshal(body, &tree)
	if err != nil {
		return Folder{}, err
	}

	return tree, err

}
