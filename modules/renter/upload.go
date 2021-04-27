package renter

import (
	"fmt"
	"github.com/ftsmchl/system/modules"
	"github.com/ftsmchl/system/modules/renter/renterfile"
	"os"
)

func (r *Renter) Upload(source string) error {

	//source := "/" + pathTrim

	//Check if the file is a directory
	sourceInfo, err := os.Stat(source)

	if err != nil {
		fmt.Println("something went wrong with stat in our file ", err)
		return err
	}

	if sourceInfo.IsDir() {
		fmt.Println("The source given is a directory ..")
		return err
	}

	//Check for read access
	file, err := os.Open(source)
	if err != nil {
		fmt.Println("something gone wrong while opening the file : ", err)
	}
	file.Close()

	ErasureCode, _ := modules.NewRSSubCode(defaultDataPieces, defaultParityPieces)

	r.mu.Lock()
	numContracts := len(r.storageContracts)
	r.mu.Unlock()

	requiredContracts := (ErasureCode.NumPieces() + ErasureCode.MinPieces()) / 2

	if numContracts < requiredContracts {
		return fmt.Errorf("not enough contracts to upload the file : got %v, needed %v", numContracts, requiredContracts)
	}

	rentfile := renterfile.New(source, ErasureCode, uint64(sourceInfo.Size()), sourceInfo.Mode())

	if sourceInfo.Size() == 0 {
		return nil
	}

	//connects a worker with a specific storage contract
	r.updateWorkerPool()

	r.buildAndPushChunks(rentfile)

	//fmt.Println(rentfile)

	r.uploadHeap.newUploads <- struct{}{}

	return nil
}
