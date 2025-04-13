package domain

type Task struct {
	RequestId  string
	Hash       string
	Alphabet   string
	MaxLength  int
	PartNumber int
	PartCount  int
}
