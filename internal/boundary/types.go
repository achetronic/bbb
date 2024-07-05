package boundary

// Person represents the data joint between an Account and User for a human
type Person struct {
	AccountId string
	UserId    string
	Email     string
	Subject   string
}

// Group represents TODO
type Group struct {
	Id      string
	Name    string
	Version uint32

	Members []string
}
