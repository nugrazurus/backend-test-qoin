package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var player, dice int

	// Input jumlah pemain
	fmt.Print("Masukkan jumlah pemain: ")
	_, err := fmt.Scanf("%d\n", &player)
	if err != nil {
		log.Fatal("Nilai yang dimasukkan bukan angka", err)
	}

	// Input jumlah dadu
	fmt.Print("Masukkan jumlad Dadu: ")
	_, err = fmt.Scanf("%d", &dice)
	if err != nil {
		log.Fatal("Nilai yang dimasukkan bukan angka", err)
	}

	// Jalankan permainan
	fmt.Printf("Pemain = %d, Dadu = %d \n", player, dice)
	fmt.Println("=======================")
	play(player, dice)
}

type Player struct {
	Score int
	Dice []int
}

type Players map[int]Player
var dice = []int{1, 2, 3, 4, 5, 6}
var onlyOnce sync.Once

// Fungsi untuk lempar dadu n sebagai banyaknya dadu
func rollDice(n int) []int {
	var dices []int
	if n < 1 {
		return []int{}
	}
	
	onlyOnce.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})
	
	for i := 0; i < n; i++ {
		dices = append(dices, dice[rand.Intn(len(dice))])	
	}
	
	return dices
}

// m sebagai jumlah pemain dan n sebagai jumlah dadu
func play(m, n int) {
	players := make(Players)
	fmt.Println("Giliran 1 lempar dadu:")

	// Mulai permainan
	for i := 1; i <= m; i++ {
		dice := rollDice(n)
		players[i] = Player{
			Score: 0,
			Dice: dice,
		}

		// Print detail player
		playerName := fmt.Sprintf("Pemain #%d", i)
		fmt.Printf("%10s (%d): %v \n", playerName, players[i].Score, players[i].Dice)
	}

	remain := evaluate(players)
	fmt.Println(remain)
}

// Fungsi untuk evaluasi
func evaluate(players Players) int {
	fmt.Println("Setelah evaluasi")
	remainPlayer := 0

	for i := 1; i <= len(players); i++ {
		player := players[i]
		for j := 0; j < len(player.Dice); j++ {
			if player.Dice[j] == 1 {
				if i > 1 {
					if entry, ok := players[i - 1]; ok {
						entry.Dice = append(entry.Dice, 1)
						players[i -1] = entry
					}
					player.Dice = append(player.Dice[:j], player.Dice[j+1:]...)
					j--
				}
			} else if player.Dice[j] == 6 {
				player.Score += 1
				player.Dice = append(player.Dice[:j], player.Dice[j+1:]...)
				j--
			}
		}
		players[i] = player
	}

	// Print detail setelah evaluasi
	for i := 1; i <= len(players); i++ {
		player := players[i]
		playerName := fmt.Sprintf("Pemain #%d", i)
		if len(player.Dice) > 0 {
			remainPlayer += 1
		}
		fmt.Printf("%10s (%d): %v \n", playerName, player.Score, player.Dice)
	}

	return remainPlayer
}
