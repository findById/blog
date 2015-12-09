package entity

type User struct {
	Id        int
	First_Name string
	Last_Name  string
	Email     string
	Password  string
	Status      string
	DelFlg      int
}

type Article struct {
	Id          int
	Title       string
	Content     string
	Keywords    string
	Description string
	Timestamp   string
	Status      string
	DelFlg      int
	UserId      int
	Lang        string
	Tag        string
}
