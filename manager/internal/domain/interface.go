package domain

type IHashUseCase interface {
	CrackHash(requestId *string, hash string, maxLength int) (HashData, error)
	GetStatus(requestId string) (HashData, int, error)
	UpdateRequest(requestId string, data []string) error
}

type IHTTPRepository interface {
	SendTaskToWorker(workerURL string, task WorkerTask) error
}

type ISQLRepository interface {
	Save(requestID string, hash string, maxLength int) error
	Remove(requestId string) error
	GetAll() ([]CrackHash, error)
}
