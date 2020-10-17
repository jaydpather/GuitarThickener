package main

import (
	"fmt"
	"io"
	"os"
	"math/rand"
	"github.com/cryptix/wav"
)

/*
TODO:
  * MOST DEVS: main func loads file, modifies values and writes to file
	* write the whole thing in 1 go, minimal testing.
	  * just verify that it runs to the end, don't check for correct results

  * STEP 1: just copy file (don't modify)
    * extract simpleRead into a new function
    * get int32s - call ReadSample instead of ReadRawSample

    * extract simpleSweep into a function
	  * accepts int32[]
	  * writes file to disk
	    * takes # samples per second from source file
		* bits per sample is 32, b/c wav library only includes 32-bit output
		
	+ had to google slices: declaring, appending, len()
	+ wrote file, came out distorted
	  + complete guess: each byte[] in the 16 bits per sample is 1 byte per channel
		+ this means you need to cast to int16 to average the 2 values. (this is how you convert to mono)
	+ SOLUTION: converted file to mono using web app
	  + had to double the sample rate, not sure why	
	

* STEP 2: modify sample values
  * process samples
	* randomly modify each sample by +/- 1%
	  * convert to float32
	  * randomly decide sign
	  * randomly decide delta value
	  * apply delta value
	  * cast modified float value to int32


*/


func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func readWavFile() []int32 {
	testInfo, err := os.Stat(os.Args[1])
	checkErr(err)

	testWav, err := os.Open(os.Args[1])
	checkErr(err)

	wavReader, err := wav.NewReader(testWav, testInfo.Size())
	checkErr(err)

	fmt.Println("Hello, wav")
	fmt.Println(wavReader)

	samples := []int32{}
	for {
		sample, err := wavReader.ReadSample()
		
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		samples = append(samples, sample)
	}

	return samples
}

func getModifiedSample(sample int32, minPctDiff int32, maxPctDiff int32) int32 {
	fSample := float32(sample)
	fMinPctDiff := float32(minPctDiff)
	fMaxPctDiff := float32(maxPctDiff)

	fSignMultiplier := float32(-1.0)
	// if rand.Float32() >= 0.5 {
	// 	fSignMultiplier = 1.0
	// }

	randRange := fMaxPctDiff - fMinPctDiff
	randValue := rand.Float32() * randRange
	fPctDiff := (fMinPctDiff + randValue) / 100.0
	fDiff := fPctDiff * fSignMultiplier * fSample

	retVal := int32(fSample + fDiff)
	return retVal
}

func getThickenedSample(sample int32) int32 {
	minPct := int32(1)
	maxPct := int32(3)

	reducedSample := int32(float32(sample) / 4.0)

	modSample1 := getModifiedSample(reducedSample, minPct, maxPct) 
	modSample2 := getModifiedSample(reducedSample, minPct, maxPct) 
	modSample3 := getModifiedSample(reducedSample, minPct, maxPct) 
	modSample4 := getModifiedSample(reducedSample, minPct, maxPct) 

	totalSample := modSample1 + modSample2 + modSample3 + modSample4
	//fTotal := float32(totalSample)
	//fTotal /= 4.0

	return totalSample //int32(fTotal)
}

func writeFile(samples []int32, fileName string) {
	wavOut, err := os.Create(fileName)
	checkErr(err)
	defer wavOut.Close()

	meta := wav.File{
		Channels:        1,
		SampleRate:      32000, //for some reason, this needs to be double the input sample rate
		SignificantBits: 16, //not sure why this needs to be 16, b/c github.com/cryptix/wav only supports int32 output
	}

	writer, err := meta.NewWriter(wavOut)
	checkErr(err)
	defer writer.Close()

	for i := 0; i < len(samples); i++ {
		newSample := getThickenedSample(samples[i])
		newSample *= 2 //for some reason, we have to double everything.

		err = writer.WriteInt32(newSample)
		checkErr(err)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("%i", len(os.Args))
		fmt.Println("%A", os.Args)
		fmt.Fprintf(os.Stderr, "Usage: simpleRead <file.wav>\n")
		os.Exit(1)
	}

	samples := readWavFile()
	//fmt.Println("%i", len(samples))
	//fmt.Println("%A", samples)

	writeFile(samples, "output_" + os.Args[1]);
}
