package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"os"
	"sync"
)

var (
	Byt = int(16)
	// proSample = make([][]byte, 100)
	echo = int(1e4)
	wg   sync.WaitGroup
)

func main() {
	wg.Add(4)
	go GoF(85, echo)
	go GoF(65, echo)
	go GoF(45, echo)
	go GoF(25, echo)
	go GoFF(85, echo)
	go GoFF(65, echo)
	go GoFF(45, echo)
	go GoFF(25, echo)
	wg.Wait()
}

// The value of p ranges from 0 to 100, indicating the probability of non-repetition.
// echo represents how many operations are performed
// Sample from sample space size 1000 samples
func GoF(p int, echo int) {

	file, _ := os.OpenFile(fmt.Sprintf("1000sample%d", 100-p), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer file.Close()
	writer := bufio.NewWriter(file)

	UK := RS(8) //UK 8bytes
	proSample := make([][]byte, 1000)
	j := int(0)
	for ; j < p*10; j++ {
		proSample[j] = RS(4)
	}
	for ; j < 1000; j++ {
		t := RS(2)
		proSample[j] = proSample[int(int(t[1])*256+int(t[0]))%(p*10)]
	}
	count := []byte{1, 0, 0, 0}
	mp := make(map[string]bool)
	re := int(0)
	for i := 0; i < echo; i++ {
		h := sha1.New()
		h.Write(getAnd3(UK, count, GetSample1000(proSample)))
		t := h.Sum(nil)
		for data := range t {
			v := t[data]
			s := ByteToBit(int(v))
			writer.WriteString(s)
			writer.Flush()
		}

		_, ok := mp[string(t)]
		if ok {
			re++
			continue
		}
		mp[string(t)] = true
		if count[0] == 255 {
			count[0] = 0
			// count[1]++
			if count[1] == 255 {
				count[1] = 0
				if count[2] == 255 {
					count[2] = 0
					// 	count[3]++ //4bytes
				} else {
					count[2]++ //3bytes
				}
			} else {
				count[1]++ //2bytes
			}
		} else {
			count[0]++
		}
	}
	for n := range mp {
		delete(mp, n)
	}
	fmt.Printf("1000,Pro: %d Re: %d Echo: %d\n", 100-p, re, echo)
	wg.Done()
}

// The value of p ranges from 0 to 100, indicating the probability of non-repetition.
// echo represents how many operations are performed
// Sample from sample space size 100 samples
func GoFF(p int, echo int) {

	file, _ := os.OpenFile(fmt.Sprintf("100sample%d", 100-p), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer file.Close()

	writer := bufio.NewWriter(file)
	UK := RS(8)
	proSample := make([][]byte, 100)
	j := int(0)
	for ; j < p; j++ {
		proSample[j] = RS(4)
	}
	for ; j < 100; j++ {
		t := RS(1)
		proSample[j] = proSample[int(t[0])%p]
	}
	count := []byte{1, 0, 0, 0}
	mp := make(map[string]bool)
	re := int(0)
	for i := 0; i < echo; i++ {
		h := sha1.New()
		h.Write(getAnd3(UK, count, GetSample100(proSample)))
		t := h.Sum(nil)
		for data := range t {
			v := t[data]
			s := ByteToBit(int(v))
			writer.WriteString(s)
			writer.Flush()
		}
		_, ok := mp[string(t)]
		if ok {
			re++
			continue
		}
		mp[string(t)] = true
		if count[0] == 255 {
			count[0] = 0
			// count[1]++
			if count[1] == 255 {
				count[1] = 0
				if count[2] == 255 {
					count[2] = 0
					// 	count[3]++ //4bytes
				} else {
					count[2]++ //3bytes
				}
			} else {
				count[1]++ //2bytes
			}
		} else {
			count[0]++
		}
	}
	for n := range mp {
		delete(mp, n)
	}
	fmt.Printf("100,Pro: %d Re: %d Echo:%d\n", 100-p, re, echo)
	wg.Done()
}

// return a\|b\|c
func getAnd3(a []byte, b []byte, c []byte) []byte {
	la := len(a)
	lb := len(b)
	lc := len(c)
	cc := make([]byte, la+lb+lc)
	i := 0
	for ; i < la; i++ {
		cc[i] = a[i]
	}
	for ; i < la+lb; i++ {
		cc[i] = b[i-la]
	}
	for ; i < la+lb+lc; i++ {
		cc[i] = c[i-la-lb]
	}
	return cc
}

// Sample from sample space size 100 samples
func GetSample100(proSample [][]byte) []byte { //1-100 ->0-19
	p := RS(1)
	return proSample[int(p[0])%100]
}

// Sample from sample space size 1000 samples
func GetSample1000(proSample [][]byte) []byte { //1-100 ->0-19
	p := RS(2)
	return proSample[int(int(p[1])*256+int(p[0]))%1000]
}

// Return a random number of Byt bytes
func RS(Byt int) []byte {
	r := make([]byte, Byt)
	rand.Read(r)
	return r
}

// Byte to Bit ,return a string
func ByteToBit(b int) string {
	var s string
	for i := 0; i < 8; i++ {
		if b&1 == 1 {
			s += "1"
		} else {
			s += "0"
		}
		b /= 2
	}
	return s
}
