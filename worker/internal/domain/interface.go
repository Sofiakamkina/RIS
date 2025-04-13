package domain

type ITaskCrackerUseCase interface {
	CrackTask(task Task) error
}

type IRepository interface {
	SendCrackedTaskToManager(requestId string, words []string) error
}
