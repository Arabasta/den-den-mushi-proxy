package algo

type Type string

const (
	AlgorithmRandom         Type = "random"
	AlgorithmPowerOf2Random Type = "power_of_2_random"
)

var Rithms = map[Type]Rithm{
	AlgorithmRandom:         NewRandom(),
	AlgorithmPowerOf2Random: NewPowerOf2Random(),
}
