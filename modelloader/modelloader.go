package modelloader

import (
	"fmt"
	"github.com/go-skynet/go-llama.cpp"
)

type ModelLoader struct {
	Model   *llama.LLama
	Threads int
}

// NewModelLoader creates a new ModelLoader
func NewModelLoader(model string, threads, contextlength, gpulayers int) (*ModelLoader, error) {
	l, err := llama.New(model, llama.EnableF16Memory, llama.SetContext(contextlength), llama.SetGPULayers(gpulayers),
		llama.SetModelSeed(47))
	if err != nil {
		fmt.Println("Loading the model failed:", err.Error())
		return nil, fmt.Errorf("loading the model failed: %w", err)
	}
	// Return the loaded Model
	return &ModelLoader{
		Model:   l,
		Threads: threads,
	}, nil
}

// Predictor Predicts
func (ml *ModelLoader) Predictor(input string, stopwords []string, maxtokens, seed, topk int, topp float32) (string, error) {
	answer, err := ml.Model.Predict(input, llama.SetStopWords(stopwords...), llama.SetTokens(maxtokens))
	if err != nil {
		return "", fmt.Errorf("predicting the model failed: %w", err)
	}
	return answer, nil
}
