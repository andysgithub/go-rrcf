package rrcf

import (
	"encoding/json"
	"io/ioutil"
)

// SaveTree saves a tree as json data to the specified file
func SaveTree(tree RCTree, filename string) {
	treeJSON, _ := json.MarshalIndent(tree, "", " ")
	ioutil.WriteFile(filename, treeJSON, 0644)
}

// SaveForest saves a forest as json data to the specified file
func SaveForest(forest []RCTree, filename string) {
	forestJSON, _ := json.MarshalIndent(forest, "", " ")
	ioutil.WriteFile(filename, forestJSON, 0644)
}

// LoadTree -
func LoadTree() {

}

// LoadForest -
func LoadForest() {

}
