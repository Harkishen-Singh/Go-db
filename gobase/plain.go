package gobase

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"os/exec"
	"log"
)

const separator = ".##."
//EnvironmentPath ...
var EnvironmentPath = os.Getenv("GOPATH") + "/src/github.com/Harkishen-Singh/GO-file-database/"

func filter(text string) (string, bool) {

	if strings.Contains(text, separator) {
		return strings.Replace(text, separator, " ", -1), true
	}
	return "", false

}

func checkExistingDir(name string) bool {

	_, err := exec.Command("ls", name).Output()
	if err != nil {
		return false
	}
	return true
}

func makeDir(name string) bool {

	name = EnvironmentPath + "warehouse/" + name
	/**
		* returns the number of directories contained in the gien path [string]
	*/
	dirOccurence := func(path string) uint16 {
		var occ uint16
		for i:= 0; i< len(path); i++{
			if path[i] == '/' {
				occ++
			}
		}
		return occ
	}

	/**
		* creates new directories on name [string]
	*/
	createDirectoriesChain := func(name string) bool {
		if !checkExistingDir(name) {
			_, err := exec.Command("mkdir", "-p", name).Output()
			if err != nil {
				fmt.Println(err)
			}
			return true
		}

		return false
	}

	/**
		* controls the process of creation of chained directories on path [string]
	*/
	directoryController := func(path string) bool {
		tempPath := path + "/"
		createDirectoriesChain(tempPath)
		return true
	}

	if dirOccurence(name) == uint16(1) {
		if !checkExistingDir(name) {
			_, err := exec.Command("mkdir", name).Output()
			if err != nil {
				panic(err)
			}
		}
	} else {
		directoryController(name)
	}
	return true
}

func collectionStatus(collectionPath string) bool {

	warehouse() // check warehouse availability
	res := strings.LastIndex(collectionPath, "/")
	var subPath string
	if res != -1 {
		subPath = collectionPath[0: res]

		if !makeDir(subPath) {
			os.Exit(50005)
		}
	}
	expAddr := EnvironmentPath + "warehouse/" + subPath
	result, err := exec.Command("ls", expAddr).Output()
	if err != nil {
		fmt.Println("Path not found!")
		return false
	}
	var existingFiles = strings.Split(string(result), "\n")
	var checkStatus bool

	for _, element := range existingFiles {
		if element == collectionPath[res+1:] + ".data" {
			checkStatus = true
			break
		} else {
			checkStatus = false
		}
	}

	return checkStatus


}

func createCollection(address string) bool {

	address = EnvironmentPath + "warehouse/" + address + ".data"
	err := ioutil.WriteFile(address, []byte("default"), 0777)
	if err != nil {
		fmt.Println("Error in createCollection Address: "+address)
		fmt.Println(err)
		return false
	}
	return true

}

//Retrive ...
func Retrive(address string) (string, bool) {

	var documentAvailable = collectionStatus(address)
	var data string
	address = EnvironmentPath + "warehouse/" + address + ".data"
	if documentAvailable {
		openfile, err := ioutil.ReadFile(address)
		if err != nil {
			return "ERROR", false
		}
		data = string(openfile)
		return data, true
	}
	return "DOCUMENT_UNAVAILABLE", false

}

//CollectionsAvailable ...
func CollectionsAvailable(address string) ([]string, bool) {

	var existingCollections []string
	if address != "/" {
		path := EnvironmentPath + "warehouse/" + address
		response, err := exec.Command("ls", path).Output()
		if err != nil {
			fmt.Println("Error while looking for Collections, at Address: "+address)
			log.Fatal(err)
		}
		existingCollections = strings.Split(string(response), "\n")
	} else {
		response, err := exec.Command("ls", "warehouse/").Output()
		if err != nil {
			fmt.Println("Error while looking for Collections, at Address: "+address)
			log.Fatal(err)
		}
		existingCollections = strings.Split(string(response), "\n")
	}
	return existingCollections, true

}

//Save ...
func Save(path string, data string) bool {

	exists := collectionStatus(path)
	if exists == false {
		fmt.Println("No Collection existing at the specified datapath. Creating one ...")
		createCollection(path)
	}
	var address = EnvironmentPath + "warehouse/" + path + ".data"
	fmt.Println(address)
	file, err := os.OpenFile(address, os.O_WRONLY, 0600)
	file.Seek(0, 0)
	file.Truncate(0)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(data)
	if err != nil {
		fmt.Println("Error occured while writing the following data:")
		fmt.Println("\n" + data)
		fmt.Println("Address: warehouse/"+address)
		fmt.Println(err)
		return false
	}
	return true

}

//Delete ...
func Delete(path string) bool {

	path = EnvironmentPath + "warehouse/" + path
	_, err := exec.Command("rm", "-R", path).Output()
	if err != nil {
		_, err2 := exec.Command("rm", "-R", path + ".data").Output()
		if err2 != nil {
			panic("Deletion not possible ''Type 2''. Collection doesnot exists in Address: " + path)
		}
	}
	return true

}

func warehouse() {

	db := "warehouse"
	resp, err := exec.Command("ls", EnvironmentPath).Output()
	if err != nil {
		panic(err)
	}
	var checkWarehouse bool
	var respStringArr = strings.Split(string(resp), "\n")
	for _, ele := range respStringArr {
		if ele == "warehouse" {
			checkWarehouse = true
			break
		}
	}
	if !checkWarehouse {
		exec.Command("mkdir", db).Output()
	}

}
