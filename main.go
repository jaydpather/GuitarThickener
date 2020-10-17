package main

import (
	"fmt"
	"io"
	"os"

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
		curSample, err := wavReader.ReadSample()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		curSampleAsInt := int(curSample)

		samples = append(samples, curSampleAsInt)
	}

	return samples
}

func writeFile(samples []int) {
	wavOut, err := os.Create("Test.wav")
	checkErr(err)
	defer wavOut.Close()

	meta := wav.File{
		Channels:        1,
		SampleRate:      44100,
		SignificantBits: 16, //hardcoded to 32 bits per sample, b/c that's all that's supported by github.com/cryptix/wav
	}

	writer, err := meta.NewWriter(wavOut)
	checkErr(err)
	defer writer.Close()

	for i := 0; i < len(samples); i++ {
		curSampleAsInt32 := int32(samples[i])
		//fmt.Println("curSampleAsInt32: %i", curSampleAsInt32)
		err = writer.WriteInt32(curSampleAsInt32)
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

	writeFile(samples);
}
