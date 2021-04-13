package renter

import (
	"fmt"
	"os"
)

func (r *Renter) Upload(source string) {

	//Check if the file is a directory
	sourceInfo, err := os.Stat(source)

	if err != nil {
		fmt.Println("something went wrong with stat in our file ", err)
		return
	}

	if sourceInfo.IsDir() {
		fmt.Println("The source given is a directory ..")
		return
	}

	//Check for read access
	file, err := os.Open(source)
	if err != nil {
		fmt.Println("somethig gone wrong while opening the file : ", err)
	}
	file.Close()
}
