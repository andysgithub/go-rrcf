package rrcf

import (
	"encoding/json"
	"io/ioutil"
)

// SaveTree -
func SaveTree(tree RCTree, filename string) {
	treeJSON, _ := json.MarshalIndent(tree, "", " ")
	ioutil.WriteFile(filename, treeJSON, 0644)
}

// SaveForest -
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
