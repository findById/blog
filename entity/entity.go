package entity

type User struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type Article struct {
	Id        int
	Title     string
	Content   string
	Timestamp string
	Status    int
	DelFlg    int
	UserId    int
}
