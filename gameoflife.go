package main

import (
	"bufio"
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/eiannone/keyboard"
	"math/rand"
	"os"
	"time"
)

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}


func main() {
	var (
		nums    [][]bool
		M       int
		N       int
		errFile error
	)
	nums, M, N, errFile = generateField()
	if errFile != nil {
		fmt.Println("Введеные данные корректны, введите снова:")
		for errFile != nil {
			nums, M, N, errFile = generateField()
		}
	}

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")
	go changeNums(nums)
	go printSec(nums, M, N)

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc {
			break
		}
	}
}

// Вызов генерации поля в зависимости от выбора пользователя
func generateField() ([][]bool, int, int, error){
	var con uint8
	fmt.Println("1. Загрузить случайное поле\n2. Загрузить поле из файла\nВаш выбор:")
	fmt.Scan(&con)

	var nums[][]bool
	var M, N int

	if con == 1 {
		nums, M, N = genRandomField()
	} else if con == 2 {
		var str string
		fmt.Println("Введите путь к файлу:")
		fmt.Scan(&str)
		nums, M, N = genFileField(str)
	}
	
	if M == 0 || N == 0 {
		return nums, M, N, errors.New("одно из данных нулевое")
	}
	return nums, M, N, nil
}

// Генерация поля из файла
func genFileField(str string) ([][]bool, int, int) {
	var M, N int
	var nums [][]bool

	readFile, err := os.Open(str)
	if err != nil {
		panic("Ошибка открытия файла состояния поля")
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	readFile.Close()

	for _, val := range lines {
		var temp []bool
		N = len(val)
		for i := range val {
			ibool := false
			if val[i] == '1' {
				ibool = true
			}
			temp = append(temp, ibool)
		}
		nums = append(nums, temp)
	}

	M = len(lines)
	return nums, M, N
}

// Случайная генерация поля
func genRandomField() ([][]bool, int, int){
	var M, N int
	var nums [][]bool

	fmt.Println("Введите размер игрового поля через пробел: ")
	fmt.Scan(&M, &N)

	for i := 0; i < M; i++ {
		var temp []bool
		for j := 0; j < N; j++ {
			condition := randBool()
			temp = append(temp, condition)
		}
		nums = append(nums, temp)
	}
	return nums, M, N
}

// Возврат случайного булевого числа
func randBool() bool {
	var src cryptoSource
	rnd := rand.New(src)
	con := rnd.Intn(2)
	if con == 1 {
		return true
	} else {
		return false
	}
}

// Постоянное изменение значения поля
func changeNums(nums [][]bool){
	for {
		time.Sleep(1 * time.Second) // Для постепенного изменения поля, иначе решение происходит слишком быстро (ненаглядно)
		for i := range nums {
			for j := range nums[i] {
				board1 := len(nums)-1
				board2 := len(nums[i])-1
				sumNeighbour := counterNeighbours(i, j,board1, board2, nums)

				if nums[i][j] == true && sumNeighbour != 2 && sumNeighbour != 3 {
					nums[i][j] = false
				}
				if nums[i][j] == false && sumNeighbour == 3 {
					nums[i][j] = true
				}

			}
		}
	}
}

// Подсчет всех соседей
func counterNeighbours(i, j, board1, board2 int, nums [][]bool) int8 {
	var sumNeighbour int8 = 0

	for h := i-1; h <= i+1; h++ {
		for  g := j-1; g <= j+1; g++ {
			if h == i && g ==j {
				continue
			}
			if h < 0 || g < 0 || h > board1 || g > board2 {
				continue
			}
			if nums[h][g] {
				sumNeighbour++
			}
		}
	}
	return sumNeighbour
}

// Печать состояния поля с разницей в секунду
func printSec(nums [][]bool, M int, N int){
	for {
		fmt.Print("Map:\n")

		for i := 0; i < M; i++ {
			for j := 0; j < N; j++ {
				fmt.Printf(" %v",Btoi(nums[i][j]))
			}
			fmt.Print("\n")
		}

		time.Sleep(1 * time.Second)
	}
}

func Btoi(b bool) int8 {
	if b {
		return 1
	}
	return 0
}