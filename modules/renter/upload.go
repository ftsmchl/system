package renter

import (
	"fmt"
	"os"
)

func (r *Renter) Upload(source string) {

	sourceInfo, err := os.Stat(source)

	if err != nil {
		fmt.Println("something went wrong with stat in our file ", err)
		return
	}

	if sourceInfo.IsDir() {
		fmt.Println("The source given is a directory ..")
		return
	}

}
