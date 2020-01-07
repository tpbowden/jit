package database

type Author struct {
	name  string
	email string
}

func NewAuthor(name string, email string) Author {
	return Author{
		name:  name,
		email: email,
	}
}
