package main

import (
	"fmt"
	"io"
	"os"

	"github.com/cryptix/wav"
)

/*
TODO:
  * STEP 1: just copy file (don't modify)
  * STEP 2: modify sample values

  * extract simpleRead into a function that converts wav file into byte[][]
	* each index in the byte[][] will be a byte[] of length 2, representing the 16-bit integer of the sample
  * convert to int32
	* iterate over byte[][]	
		* convert each int 16 to int 32
		* return an int32[]
  * process samples
	* randomly modify each sample by +/- 1%
	  * convert to float32
	  * randomly decide sign
	  * randomly decide delta value
	  * apply delta value
	  * cast modified float value to int32

  * extract simpleSweep into a function
	* accepts int32[]
	* writes file to disk
	  * takes # samples per second from source file
	  * bits per sample is 32, b/c wav library only includes 32-bit output


*/

func main() {
	if len(os.Args) != 2 {
		fmt.Println("%i", len(os.Args))
		fmt.Println("%A", os.Args)
		fmt.Fprintf(os.Stderr, "Usage: simpleRead <file.wav>\n")
		os.Exit(1)
	}
	testInfo, err := os.Stat(os.Args[1])
	checkErr(err)

	testWav, err := os.Open(os.Args[1])
	checkErr(err)

	wavReader, err := wav.NewReader(testWav, testInfo.Size())
	checkErr(err)

	fmt.Println("Hello, wav")
	fmt.Println(wavReader)

sampleLoop:
	for {
		s, err := wavReader.ReadRawSample()
		if err == io.EOF {
			break sampleLoop
		} else if err != nil {
			panic(err)
		}

		fmt.Printf("Sample: <%v>\n", s)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
