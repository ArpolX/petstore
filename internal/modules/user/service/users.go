package service

type ServiceUser struct {
	UserName  string
	FirstName string
	LastName  string
	Email     string
	Password  string
	Phone     string
}

type ServiceUserArray struct {
	UserArray []ServiceUser
}
