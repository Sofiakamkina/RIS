package domain

import (
	"crypto/md5"
	"encoding/hex"
	"log/slog"
)

type TaskCrackerUseCase struct {
	repository IRepository
	logger     *slog.Logger
}

func NewTaskCrackerUseCase(repository IRepository, logger *slog.Logger) *TaskCrackerUseCase {
	return &TaskCrackerUseCase{
		repository: repository,
		logger:     logger,
	}
}

func (t *TaskCrackerUseCase) CrackTask(task Task) error {
	go func() {
		alphabetLength := len(task.Alphabet)
		totalCombinations := pow(alphabetLength, task.MaxLength)

		startIndex := (totalCombinations / task.PartCount) * task.PartNumber
		endIndex := (totalCombinations / task.PartCount) * (task.PartNumber + 1)

		if task.PartNumber == task.PartCount-1 {
			endIndex = totalCombinations
		}

		t.logger.Info("starting task", "requestId", task.RequestId, "partNumber", task.PartNumber, "startIndex", startIndex, "endIndex", endIndex)

		var crackedWords []string
		for i := startIndex; i < endIndex; i++ {
			word := getWordByIndex(i, task.Alphabet, task.MaxLength)
			hash := computeMD5(word)

			if hash == task.Hash {
				crackedWords = append(crackedWords, word)
				t.logger.Info("found matching word", "word", word, "hash", hash)
			}
		}

		if err := t.repository.SendCrackedTaskToManager(task.RequestId, crackedWords); err != nil {
			t.logger.Error("failed to send cracked words to manager", "error", err)
		}

		t.logger.Info("task completed", "requestId", task.RequestId, "partNumber", task.PartNumber)
	}()

	return nil
}

func computeMD5(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func getWordByIndex(index int, alphabet string, maxLength int) string {
	word := ""
	alphabetLength := len(alphabet)

	for i := 0; i < maxLength; i++ {
		charIndex := index % alphabetLength
		word = string(alphabet[charIndex]) + word
		index = index / alphabetLength
	}

	return word
}

func pow(a, b int) int {
	result := 1

	for i := 0; i < b; i++ {
		result *= a
	}

	return result
}
