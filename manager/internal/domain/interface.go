package domain

type IHashUseCase interface {
	CrackHash(hash string, maxLength int) (HashData, error)
	GetStatus(requestId string) (HashData, error)
	UpdateRequest(requestId string, data []string) error
}

type IRepository interface {
	SendTaskToWorker(workerURL string, task WorkerTask) error
}
