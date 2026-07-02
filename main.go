package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
)

type Difficulty struct {
	name        string
	maxNumber   int
	maxAttempts int
}

type Result struct {
	Date     string `json:"date"`
	Outcome  string `json:"outcome"`
	Attempts int    `json:"attempts"`
	Level    string `json:"level"`
}

func readGuess() (int, error) {
	var input string
	fmt.Scan(&input)
	return strconv.Atoi(input)
}

func chooseDifficulty() Difficulty {
	easy := Difficulty{name: "Easy", maxNumber: 50, maxAttempts: 15}
	medium := Difficulty{name: "Medium", maxNumber: 100, maxAttempts: 10}
	hard := Difficulty{name: "Hard", maxNumber: 200, maxAttempts: 5}

	fmt.Println("Выберите сложность:")
	fmt.Println("1 - Easy (1-50, 15 попыток)")
	fmt.Println("2 - Medium (1-100, 10 попыток)")
	fmt.Println("3 - Hard (1-200, 5 попыток)")

	for {
		fmt.Print("Ваш выбор: ")
		var choice string
		fmt.Scan(&choice)

		switch choice {
		case "1":
			return easy
		case "2":
			return medium
		case "3":
			return hard
		default:
			fmt.Println("Нет такого варианта, введите 1, 2 или 3.")
		}
	}
}

func compare(guess, secret int) string {
	if guess < secret {
		return "Секретное число больше👆"
	}
	return "Секретное число меньше👇"
}

func hint(guess, secret int) string {
	diff := guess - secret
	if diff < 0 {
		diff = -diff
	}

	switch {
	case diff <= 5:
		return "🔥 Горячо"
	case diff <= 15:
		return "🙂 Тепло"
	default:
		return "❄️ Холодно"
	}
}

func printHistory(history []int) {
	fmt.Print("Введённые числа: ")
	for i, g := range history {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(g)
	}
	fmt.Println()
}

func askPlayAgain() bool {
	fmt.Print("Сыграть ещё раз? (y/n): ")
	var answer string
	fmt.Scan(&answer)
	return answer == "y" || answer == "Y"
}

func loadResults() []Result {
	data, err := os.ReadFile("results.json")
	if err != nil {
		return nil
	}

	var results []Result
	err = json.Unmarshal(data, &results)
	if err != nil {
		return nil
	}
	return results
}

func saveResult(result Result) {
	results := loadResults()
	results = append(results, result)

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println("Не удалось подготовить результат к сохранению:", err)
		return
	}

	err = os.WriteFile("results.json", data, 0644)
	if err != nil {
		fmt.Println("Не удалось записать файл:", err)
	}
}

func playGame(level Difficulty) Result {
	fmt.Printf("Игра 'Угадай число' - от 1 до %d началась!\n", level.maxNumber)
	fmt.Printf("Угадайте число за %d попыток!\n", level.maxAttempts)

	secret := rand.Intn(level.maxNumber) + 1
	var history []int

	won := false
	attempt := 1
	for attempt <= level.maxAttempts {
		color.Yellow("Попытка #%d", attempt)
		fmt.Print("Введите число: ")

		guess, err := readGuess()
		if err != nil {
			fmt.Println("Это не число! Введите целое число.")
			continue
		}

		history = append(history, guess)

		if guess == secret {
			color.Green("Вы угадали!🙌")
			won = true
			break
		}
		fmt.Println(compare(guess, secret))
		fmt.Println(hint(guess, secret))
		printHistory(history)

		attempt++
	}

	if !won {
		color.Red("Вы проиграли!😢")
		fmt.Printf("Секретное число было: %d\n", secret)
	}
	fmt.Println("Игра закончена!")

	outcome := "проигрыш"
	if won {
		outcome = "победа"
	}
	return Result{
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Outcome:  outcome,
		Attempts: len(history),
		Level:    level.name,
	}
}

func main() {
	for {
		level := chooseDifficulty()
		result := playGame(level)
		saveResult(result)
		if !askPlayAgain() {
			break
		}
	}
	fmt.Println("Спасибо за игру! 👋")
}
