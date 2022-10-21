package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

var (
	WordsNum int
)

func GetAllWords() []string {
	file, err := os.Open("allWords.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	wr := bytes.Buffer{}
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		wr.WriteString(sc.Text())
	}

	words := strings.Split(wr.String(), ",")

	WordsNum = len(words)

	return words
}

func CountMostPopularLetters(words *[][]rune, mask *[5]rune) (*map[rune]int, *[5]map[rune]int) {
	OverallPopularity := make(map[rune]int)
	var IndexPopularity [5]map[rune]int
	for i := range IndexPopularity {
		IndexPopularity[i] = make(map[rune]int)
	}

	processLetter := func(letter rune, ind int) {
		if (*mask)[ind] != letter {
			OverallPopularity[letter] += 1
			IndexPopularity[ind][letter] += 1
		}
	}

	for wordInd := range *words {
		word := (*words)[wordInd]
		for i := range word {
			processLetter(word[i], i)
		}
	}

	return &OverallPopularity, &IndexPopularity
}

//func sum(a interface{}) int {
//	_sum := func() int {
//		res := 0
//		for i := range a {
//			res += a[i]
//		}
//		return res
//	}
//
//	switch v := a.(type) {
//	case []int:
//		return _sum
//	}
//}

func PrintPopularity(popularityArrs ...interface{}) {
	for popularityArrInd := range popularityArrs {

		switch popularityArr := popularityArrs[popularityArrInd].(type) {
		case *map[rune]int:
			for key, val := range *popularityArr {
				fmt.Printf("Overall popularity of letter %c is %v\n", key, val)
			}
		case *[5]map[rune]int:
			for letter := 'а'; letter <= 'я'; letter++ {

				summ := 0 /// Заменить на функцию sum
				for i := 0; i < 5; i++ {
					summ += (*popularityArr)[i][letter]
				}

				if summ != 0 {
					fmt.Printf("%c: ", letter)
					for i := 0; i < 5; i++ {
						fmt.Print((*popularityArr)[i][letter], " ")
					}
					fmt.Println()
				}
			}
		default:
			fmt.Printf("%v is not popularity arr\n", popularityArr)
		}
	}
}

func ConvertStringArrayToRuneArray(stringArr *[]string) *[][]rune {
	runeArr := make([][]rune, 0, len(*stringArr))
	for i := range *stringArr {
		runeArr = append(runeArr, []rune((*stringArr)[i]))
	}
	return &runeArr
}

func WordLooksLike(word, presentLetters, blockedLetters *[]rune, mask *[5]rune) bool {
	for i := range *mask {
		letter := (*mask)[i]
		if letter != '*' && letter != (*word)[i] {
			return false
		}
	}

	for i := range *presentLetters {
		if !strings.Contains(string(*word), string((*presentLetters)[i])) {
			return false
		}
	}

	for i := range *blockedLetters {
		if strings.Contains(string(*word), string((*blockedLetters)[i])) {
			return false
		}
	}

	return true
}

func GetSimilarWords(mask *[5]rune, presentLetters, blockedLetters *[]rune, words interface{}) *[][]rune {
	var runeWords *[][]rune
	switch v := words.(type) {
	case *[]string:
		runeWords = ConvertStringArrayToRuneArray(v)
	case []string:
		runeWords = ConvertStringArrayToRuneArray(&v)
	case [][]rune:
		runeWords = &v
	case *[][]rune:
		runeWords = v
	default:
		panic(fmt.Sprintf("Invalid words array %T", v))
	}

	similarWords := make([][]rune, 0, WordsNum)
	for i := range *runeWords {
		if WordLooksLike(&(*runeWords)[i], presentLetters, blockedLetters, mask) {
			similarWords = append(similarWords, (*runeWords)[i])
		}
	}

	return &similarWords
}

func GetBestWord(words *[][]rune, mask *[5]rune, presentLetters, blockedLetters *[]rune) []rune {
	if len(*words) == 1 {
		return (*words)[0]
	}

	OverallPopilarity, IndexPopularity := CountMostPopularLetters(words, mask)

	PrintPopularity(OverallPopilarity, IndexPopularity)

	GetWordPriority := func(word []rune) int {
		res := 0

		for i := range word {
			if !strings.Contains(string(*presentLetters), string(word[i])) {
				res += (*OverallPopilarity)[word[i]]
			}
			if (*mask)[i] == '*' {
				res += IndexPopularity[i][word[i]]
			}
		}

		return res
	}

	AllWordsStr := GetAllWords()
	AllWords := ConvertStringArrayToRuneArray(&AllWordsStr)
	maxPrioruty, maxWord := GetWordPriority((*AllWords)[0]), (*AllWords)[0]
	for i := 1; i < len(*AllWords); i++ {
		priority := GetWordPriority((*AllWords)[i])
		//fmt.Printf("Priority of word %s is %v\n", string((*AllWords)[i]), priority)
		if priority > maxPrioruty {
			maxPrioruty, maxWord = priority, (*AllWords)[i]
		}
	}

	return maxWord
}

func main() {
	var answer string

	mask := [5]rune{'*', '*', '*', '*', '*'}
	presentLetters := make([]rune, 0, 5)
	blockedLetters := make([]rune, 0, 40)

	step := 1

	for answer != "ok" {
		words := GetSimilarWords(&mask, &presentLetters, &blockedLetters, GetAllWords())

		for i := range *words {
			fmt.Println(string((*words)[i]))
		}

		bestWord := GetBestWord(words, &mask, &presentLetters, &blockedLetters)
		fmt.Println(step, ") ", string(bestWord))
		step += 1

		fmt.Scan(&answer)

		if answer == "ok" {
			return
		}

		if strings.Count(string(mask[:]), "*") == 0 {
			fmt.Println("My answer ", string(mask[:]))
			return
		}

		for i := range answer {
			if answer[i] == 'y' {
				mask[i] = bestWord[i]
			} else if answer[i] == 'w' {
				presentLetters = append(presentLetters, bestWord[i])
			} else if answer[i] == 'g' {
				if !strings.Contains(string(mask[:]), string(bestWord[i])) &&
					!strings.Contains(string(presentLetters), string(bestWord[i])) {
					blockedLetters = append(blockedLetters, bestWord[i])
				} else {
					// What here?
				}
			} else {
				fmt.Println("wrong format")
			}
		}
	}

}
