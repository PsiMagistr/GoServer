// Матрица боев.
package game

var elementToIndex = map[string]int{
	"fire":  0,
	"earth": 1,
	"air":   2,
	"water": 3,
}

// Матрица Атака-Атака.
var atkVsAtk = [4][4]float64{
	{2.0, 1.0, 1.0, 0.5}, // Огонь против [О, З, Вз, Вд]
	{0.5, 2.0, 1.0, 1.0}, // Земля против [О, З, Вз, Вд]
	{1.0, 0.5, 2.0, 1.0}, // Воздух против [О, З, Вз, Вд]
	{1.0, 1.0, 0.5, 2.0}, // Вода против [О, З, Вз, Вд]
}

// Матрица Атака vs Щит
// [Атака][Щит]
var atkVsShield = [4][4]float64{
	{1.0, 1.0, 1.0, 0.0}, // Огонь атакует щиты [О, З, Вз, Вд]
	{0.0, 1.0, 1.0, 1.0}, // Земля атакует щиты [О, З, Вз, Вд]
	{1.0, 0.0, 1.0, 1.0}, // Воздух атакует щиты [О, З, Вз, Вд]
	{1.0, 1.0, 0.0, 1.0}, // Вода атакует щиты [О, З, Вз, Вд]
}

func GetAtkMult(attakerElement, defenderElement string, defenderIsShield bool) float64 {
	i := elementToIndex[attakerElement]
	j := elementToIndex[defenderElement]
	if defenderIsShield {
		return atkVsShield[i][j]
	}
	return atkVsAtk[i][j]
}
