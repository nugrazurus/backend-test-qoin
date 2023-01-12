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
	Dice  []int
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
	round := 1
	players := make(Players)
	if len(players) < 1 {
		players.initPlayer(m, n)
	}

	remain := players.evaluate()

	// jika pemain > 1 lanjut ke ronde selanjutnya
	for remain > 1 {
		round += 1
		fmt.Println("=======================")
		fmt.Printf("Giliran %d lempar dadu:\n", round)
		players.nextRound()
		remain = players.evaluate()
	}

	// mendapatkan pemain terakhir dan pemenang
	lastPlayer, winner := players.result()
	fmt.Println("=======================")
	fmt.Printf("Game berakhir karena hanya pemain #%d yang memiliki dadu.\n", lastPlayer)
	fmt.Printf("Game dimenangkan oleh pemain #%d karena memiliki poin lebih banyak dari pemain lainnya.\n", winner)

}

// Fungsi inisialisasi pemain
func (p Players) initPlayer(m, n int) {
	fmt.Println("Giliran 1 lempar dadu:")

	// Mulai permainan
	for i := 1; i <= m; i++ {
		dice := rollDice(n)
		p[i] = Player{
			Score: 0,
			Dice:  dice,
		}

		// Print detail player
		playerName := fmt.Sprintf("Pemain #%d", i)
		fmt.Printf("%10s (%d): %v \n", playerName, p[i].Score, p[i].Dice)
	}
}

// Fungsi untuk ronde selanjutnya
func (p Players) nextRound() {
	for i := 1; i <= len(p); i++ {
		player := p[i]
		if len(player.Dice) >= 1 {
			dice := rollDice(len(player.Dice))
			player.Dice = dice
		}
		p[i] = player
		playerName := fmt.Sprintf("Pemain #%d", i)
		fmt.Printf("%10s (%d): %v \n", playerName, p[i].Score, p[i].Dice)
	}
}

// Fungsi untuk evaluasi
func (p Players) evaluate() int {
	fmt.Println("Setelah evaluasi")
	remainPlayer := 0

	for i := 1; i <= len(p); i++ {
		player := p[i]
		for j := 0; j < len(player.Dice); j++ {
			if player.Dice[j] == 1 {
				if i > 1 {
					if entry, ok := p[i-1]; ok {
						entry.Dice = append(entry.Dice, 1)
						p[i-1] = entry
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
		p[i] = player
	}

	// Print detail setelah evaluasi
	for i := 1; i <= len(p); i++ {
		player := p[i]
		playerName := fmt.Sprintf("Pemain #%d", i)
		if len(player.Dice) > 0 {
			remainPlayer += 1
		}
		fmt.Printf("%10s (%d): %v \n", playerName, player.Score, player.Dice)
	}

	return remainPlayer
}

// Fungsi untuk mengecek pemenang
func (p Players) result() (int, int) {
	var winner, lastPlayer int
	score := 0
	for i := 1; i <= len(p); i++ {
		if len(p[i].Dice) > 0 {
			lastPlayer = i
		}
		if p[i].Score > score {
			score = p[i].Score
			winner = i
		}
	}
	return lastPlayer, winner
}
