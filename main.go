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

func readWavFile() []int {
	testInfo, err := os.Stat(os.Args[1])
	checkErr(err)

	testWav, err := os.Open(os.Args[1])
	checkErr(err)

	wavReader, err := wav.NewReader(testWav, testInfo.Size())
	checkErr(err)

	fmt.Println("Hello, wav")
	fmt.Println(wavReader)

	samples := []int{}
	for {
		sample, err := wavReader.ReadSample()
		
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		samples = append(samples, int(sample))
	}

	return samples
}

func getModifiedSample(sample int, minPctDiff int, maxPctDiff int) int {
	fSample := float64(sample)
	fMinPctDiff := float64(minPctDiff)
	fMaxPctDiff := float64(maxPctDiff)

	fSignMultiplier := -1.0 //we will always decrease the sample. Originally, this was randomly determined to be 1 or -1, but I found that increasing samples can increase distortion too much

	randRange := fMaxPctDiff - fMinPctDiff
	randValue := rand.Float64() * randRange
	fPctDiff := (fMinPctDiff + randValue) / 100.0
	fDiff := fSample * fPctDiff * fSignMultiplier

	retVal := fSample + fDiff
	return int(retVal)
}

func getThickenedSample(sample int) int {
	minPct := 0
	maxPct := 2

	modSample1 := getModifiedSample(sample, minPct, maxPct) / 4.0
	modSample2 := getModifiedSample(sample, minPct, maxPct) / 4.0
	modSample3 := getModifiedSample(sample, minPct, maxPct) / 4.0
	modSample4 := getModifiedSample(sample, minPct, maxPct) / 4.0

	newSample := modSample1 + modSample2 + modSample3 + modSample4

	return newSample
}

func writeFile(samples []int, fileName string) {
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
		curSampleAsInt32 := int32(samples[i])

		newSample := getThickenedSample(int(curSampleAsInt32))
		newSample *= 2 //for some reason, we have to double everything.

		err = writer.WriteInt32(int32(newSample))
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
