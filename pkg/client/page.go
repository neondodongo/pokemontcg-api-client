package client

import "io/ioutil"

type UserFacingHTML struct{
	Title string
	Body []byte
}

func (p *UserFacingHTML)Save() (e error){
	filename := "templates/" + p.Title + ".html"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func LoadPage(title string) (*UserFacingHTML, error) {
	filename := "templates/" + title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &UserFacingHTML{Title: title, Body: body}, nil
}