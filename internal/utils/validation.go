package utils

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"GoServer/internal/models"
)

func ValidateCharacter(char *models.Character) (*models.Character, error) {
	name := strings.TrimSpace(char.Name)
	if utf8.RuneCountInString(name) < 3 {
		return nil, fmt.Errorf("Имя персонажа не должно быть короче трех символов.")
	}
	if char.Gender != "male" && char.Gender != "female" {
		return nil, fmt.Errorf("Неверно указан пол персонажа.")
	}
	id, err := strconv.Atoi(char.AvatarID)
	if err != nil || id < 1 || id > 9 {
		return nil, fmt.Errorf("Выбрана несуществующая аватарка (должна быть от 1 до 9)")
	}
	char.Name = name
	return char, nil
}
