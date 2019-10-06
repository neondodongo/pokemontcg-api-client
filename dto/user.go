package dto

type User struct{
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Email string `json:"email" bson:"email"`
	Title string `json:"title" bson:"title"`
	Rank string `json:"rank" bson:"rank"`
	CardCol []ColItem `json:"cardCol" bson:"cardCol"`
}

type ColItem struct{
	CardID string `json:"cardId" bson:"cardId"`
	Count int `json:"count" bson:"count"`
}

func (user User) InitUser(un, pw, em string){
	user.Username = un
	user.Email = em
	user.Password = pw
	user.CardCol = make([]ColItem, 0)
	user.Rank = "bitch ass"
	user.Title = "loser"
	return
}