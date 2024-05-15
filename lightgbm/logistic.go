package lightgbm

import (
	"fmt"
)

type TransformLogistic struct{}

func (t *TransformLogistic) Transform(rawPredictions []float64, outputPredictions []float64, startIndex int) error {
	if len(rawPredictions) != 1 {
		return fmt.Errorf("expected len(rawPredictions) = 1 (got %d)", len(rawPredictions))
	}

	outputPredictions[startIndex] = Sigmoid(rawPredictions[0])
	return nil
}

func (t *TransformLogistic) NOutputGroups() int {
	return 1
}

func (t *TransformLogistic) Type() TransformType {
	return Logistic
}

func (t *TransformLogistic) Name() string {
	return Logistic.Name()
}
