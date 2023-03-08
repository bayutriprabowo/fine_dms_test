package repository


type UserRepository interface {
	Select()
	Delete()
	Update()
	Create()
}
