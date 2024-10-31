package main

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"sync"
)

var (
	Byt = int(16)
	// proSample = make([][]byte, 100)
	echo = int(1e8)
	wg   sync.WaitGroup
)

func main() {
	wg.Add(8)
	go CollistionTest1000(85, echo) // collision 327334   echo 1e8
	go CollistionTest1000(65, echo) // collision 465434   echo 1e8
	go CollistionTest1000(45, echo) // collision 688424   echo 1e8
	go CollistionTest1000(25, echo) // collision 1163696   echo 1e8
	go CollistionTest100(85, echo)  // collision 3182748   echo 1e8
	go CollistionTest100(65, echo)  // collision 4310604   echo 1e8
	go CollistionTest100(45, echo)  // collision 6186124   echo 1e8
	go CollistionTest100(25, echo)  // collision 10994389  echo 1e8
	wg.Wait()
}

// p ranges from 1 to 100 for the probability of not repeating, and echo for the number of executions
// Sampling from 1000 samples
// Cnt 3bytes
func CollistionTest1000(p int, echo int) {
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
				}
			} else {
				count[1]++
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

// p ranges from 1 to 100 for the probability of not repeating, and echo for the number of executions
// Sampling from 100 samples
// Cnt 3bytes
func CollistionTest100(p int, echo int) {
	UK := RS(8) //UK 8bytes
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
	for i := 0; i < echo; i++ { //每次1e7
		h := sha1.New()
		h.Write(getAnd3(UK, count, GetSample100(proSample)))
		t := h.Sum(nil)
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
func getAnd(a []byte, b []byte) []byte {
	la := len(a)
	lb := len(b)
	c := make([]byte, la+lb)
	for i := 0; i < la; i++ {
		c[i] = a[i]
	}
	for i := la; i < la+lb; i++ {
		c[i] = b[i-la]
	}
	return c
}
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

// Sampling from 100 samples
func GetSample100(proSample [][]byte) []byte {
	p := RS(1)
	return proSample[int(p[0])%100]
}

// Sampling from 1000 samples
func GetSample1000(proSample [][]byte) []byte {
	p := RS(2)
	return proSample[int(int(p[1])*256+int(p[0]))%1000]
}
func RS(Byt int) []byte {
	r := make([]byte, Byt)
	rand.Read(r)
	return r
}
func getOxr(a []byte, b []byte) []byte {
	la := len(a)
	lb := len(b)
	if la > lb {
		c := make([]byte, la)
		i := 0
		t := 0
		for ; i < la; i++ {
			if i%lb == 0 {
				t = 0
			}
			c[i] = a[i] ^ b[t]
			t++
		}
		return c
	}
	c := make([]byte, lb)
	i := 0
	t := 0
	for ; i < lb; i++ {
		if i%la == 0 {
			t = 0
		}
		c[i] = b[i] ^ a[t]
		t++
	}
	return c
}
func getOxr3(a []byte, b []byte, c []byte) []byte {
	la := len(a)
	lb := len(b)
	lc := len(c)
	if la >= lb && la >= lc {
		cc := make([]byte, la)
		i := 0
		t := 0
		tt := 0
		for ; i < la; i++ {
			if i%lb == 0 {
				t = 0
			}
			if i%lc == 0 {
				tt = 0
			}
			cc[i] = a[i] ^ b[t] ^ c[tt]
			t++
			tt++
		}
		return cc
	}
	if lb >= la && lb >= lc {
		cc := make([]byte, lb)
		i := 0
		t := 0
		tt := 0
		for ; i < lb; i++ {
			if i%la == 0 {
				t = 0
			}
			if i%lc == 0 {
				tt = 0
			}
			cc[i] = b[i] ^ a[t] ^ c[tt]
			t++
			tt++
		}
		return cc
	}
	cc := make([]byte, lc)
	i := 0
	t := 0
	tt := 0
	for ; i < lc; i++ {
		if i%la == 0 {
			t = 0
		}
		if i%lb == 0 {
			tt = 0
		}
		cc[i] = c[i] ^ a[t] ^ b[tt]
		t++
		tt++
	}
	return cc
}
