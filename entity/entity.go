package entity

type User struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Password  string
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
