package domain

type Status string

const (
	InProgressStatus   Status = "IN_PROGRESS"
	ReadyStatus        Status = "READY"
	ErrorStatus        Status = "ERROR"
	PartialReadyStatus Status = "PARTIAL_READY"
)

type HashData struct {
	RequestId        string
	Status           Status
	Data             []string
	CompletedWorkers int
}

type WorkerTask struct {
	RequestId  string
	Hash       string
	Alphabet   string
	MaxLength  int
	PartNumber int
	PartCount  int
}
