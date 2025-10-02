package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Question struct {
	text    string
	options []string
	answer  int
}

type Game struct {
	name      string
	points    int
	questions []Question
}

func (g *Game) Init() {
	fmt.Println("Seja bem-vindo ao quiz de Hunter x Hunter!")
	fmt.Print("Digite seu nome: ")
	fmt.Scanln(&g.name)
	fmt.Printf("Olá, %s! Vamos começar o quiz.\n", g.name)
}

func (g *Game) LoadQuestions(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("erro ao ler o CSV: %w", err)
	}

	for i, record := range records {
		if i == 0 {
			continue 
		}
		if len(record) < 6 {
			return fmt.Errorf("linha %d inválida no CSV", i+1)
		}

		ans, err := strconv.Atoi(record[5])
		if err != nil {
			return fmt.Errorf("erro ao converter resposta na linha %d: %w", i+1, err)
		}

		q := Question{
			text:    record[0],
			options: record[1:5],
			answer:  ans,
		}
		g.questions = append(g.questions, q)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(g.questions), func(i, j int) {
		g.questions[i], g.questions[j] = g.questions[j], g.questions[i]
	})

	return nil
}

func (g *Game) Start() {
	for i, q := range g.questions {
		fmt.Printf("\nPergunta %d: %s\n", i+1, q.text)
		for idx, opt := range q.options {
			fmt.Printf("%d) %s\n", idx+1, opt)
		}

		var resposta int
		fmt.Print("Sua resposta: ")
		fmt.Scanln(&resposta)

		if resposta == q.answer {
			fmt.Println("✔️ Correto!")
			g.points++
		} else {
			fmt.Printf("❌ Errado! Resposta correta: %d) %s\n", q.answer, q.options[q.answer-1])
		}
	}
	fmt.Printf("\nFim do quiz, %s! Sua pontuação: %d/%d\n", g.name, g.points, len(g.questions))
}

func SaveScore(name string, points int, total int) {
	f, err := os.OpenFile("ranking.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Erro ao salvar ranking:", err)
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	writer.Write([]string{name, strconv.Itoa(points), strconv.Itoa(total)})
}

func main() {
	game := &Game{}
	game.Init()

	if err := game.LoadQuestions("perguntas.csv"); err != nil {
		fmt.Println("Erro:", err)
		return
	}

	game.Start()

	SaveScore(game.name, game.points, len(game.questions))
	fmt.Println("Resultado salvo em ranking.csv ✅")
}
