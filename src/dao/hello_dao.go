package dao

type HelloDAO struct{}

func NewHelloDAO() *HelloDAO {
	return &HelloDAO{}
}

func (d *HelloDAO) RetrieveMessage() string {
	return "Hello, World!"
}
