package main

import "fmt"
import "sync"
import "os"
import "io/ioutil"

func main() {
	words,err:=ioutil.ReadFile(os.Args[1])
	if err!=nil{
		panic(err)
	}
	
	// build new list without words with repeat letters
	var noRepeatLetterWords []byte
	for i:=0;i<len(words);i+=6 {
		word := words[i : i+5]
		if word[0]!=word[1] && word[0]!=word[2]	&& word[0]!=word[3]	&& word[0]!=word[4] && word[1]!=word[2] && word[1]!=word[3] && word[1]!=word[4] && word[2]!=word[3] && word[2]!=word[4]	&& word[3]!=word[4]{
			noRepeatLetterWords=append(noRepeatLetterWords,words[i : i+6]...)
		}			
	}
	
	wordBits := make([]uint32, len(noRepeatLetterWords)/6)
	fmt.Printf("Any 5 words with no letters in common\tWord count:%d perms:%d\n",len(wordBits),len(wordBits)*(len(wordBits)-1)*(len(wordBits)-2)*(len(wordBits)-3)/2/3/4/5*(len(wordBits)-4))
	
	// compute binary map of letters in words
	const zeroBit = 97 // 65
	for i := range wordBits {
		word := noRepeatLetterWords[i*6 : i*6+5]
		wordBits[i] = 1 << (uint32(word[0]) - zeroBit)
		wordBits[i] |= 1 << (uint32(word[1]) -zeroBit)
		wordBits[i] |= 1 << (uint32(word[2]) - zeroBit)
		wordBits[i] |= 1 << (uint32(word[3]) - zeroBit)
		wordBits[i] |= 1 << (uint32(word[4]) - zeroBit)
	}
	
	// brute force permutations
	// use very fast OR between each binary map
	// since brute force can simply parallel
	var wg sync.WaitGroup
	for w1, b1 := range wordBits[:len(wordBits)-1] {
		wg.Add(1)
		// first level perms go routined 
		go func(w1 int,b1 uint32){
			for w2, b2 := range wordBits[w1+1:] {
				if b1 & b2 != 0 {
					continue
				}
				b12:=b1 | b2
				for w3, b3 := range wordBits[w2+w1+2:] {
					if b3 & b12 != 0{
						continue
					}
					b123:=b12 | b3
					for w4,b4 := range wordBits[w3+w2+w1+3:] {
						if b4 & b123 != 0{
							continue
						}
						for w5,b5 := range wordBits[w4+w3+w2+w1+4:] {
							if b5 & (b123 | b4) != 0 {
								continue
							}
							fmt.Printf("%s %s %s %s %s\n",
							noRepeatLetterWords[w1*6:w1*6+5],
							noRepeatLetterWords[(w1+w2+1)*6:(w1+w2+1)*6+5],
							noRepeatLetterWords[(w3+w1+w2+2)*6:(w3+w1+w2+2)*6+5],
							noRepeatLetterWords[(w4+w3+w2+w1+3)*6:(w4+w3+w2+w1+3)*6+5],
							noRepeatLetterWords[(w5+w4+w3+w2+w1+4)*6:(w5+w4+w3+w2+w1+4)*6+5])
						}
					}
				}
			}
			wg.Done()
		}(w1,b1)
	}
	wg.Wait()
}

/*  word source: "valid-wordle-words.txt"

Any 5 words with no letters in common	Word count:8322 perms:332227328551594944
glent jumby prick vozhd waqfs
bemix clunk grypt vozhd waqfs
fjord gucks nymph vibex waltz
blunk cimex grypt vozhd waqfs
brung cylix kempt vozhd waqfs
chunk fjord gymps vibex waltz
jumby pling treck vozhd waqfs
clipt jumby kreng vozhd waqfs
brick glent jumpy vozhd waqfs
bling jumpy treck vozhd waqfs
brung kempt vozhd waqfs xylic

*/
