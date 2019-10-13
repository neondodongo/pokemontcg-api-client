package client

import "io/ioutil"

type Page struct{
	Title string
	Body []byte
}

func (p *Page)Save() (e error){
	filename := "templates/" + p.Title + ".html"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func LoadPage(title string) (*Page, error) {
	filename := "templates/" + title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}