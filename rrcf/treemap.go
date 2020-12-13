package rrcf

// // ToMap serializes an RCTree to a nested map
// // Returns a nested map representing all nodes in the RCTree
// func (rcTree RCTree) ToMap() *NodeObject {
// 	// Create empty map
// 	obj := NewNodeObject()

// 	// Create map to keep track of duplicates
// 	var duplicates *NodeObject

// 	for k, v := range rrcf.leaves {
// 		duplicates.SetDefault(v, []).append(k)
// 	}
// 	// Serialize tree to map
// 	rrcf.Serialize(self.root, obj, duplicates)
// 	// Return tree map
// 	return obj
// }

// // Serialise recursively serializes tree into a nested map
// func (rcTree RCTree) Serialise(node *Node, obj *NodeObject, duplicates *NodeObject) error {
// 	if node.isBranch {
// 		obj.nodeType = "Branch"
// 		obj.q = node.branch.q
// 		obj.p = node.branch.p
// 		obj.n = node.branch.n
// 		obj.b = node.branch.b
// 		obj.l = nil
// 		obj.r = nil
// 		if node.l {
// 			rrcf.Serialize(node.l, obj.l, duplicates)
// 		}
// 		if node.r {
// 			rrcf.Serialize(node.r, obj.r, duplicates)
// 		}
// 	} else if isinstance(node, Leaf) {
// 		obj.nodeType = "Leaf"
// 		obj.I = node.leaf.I
// 		obj.x = node.leaf.x
// 		obj.d = node.leaf.d
// 		obj.n = node.n
// 		obj.ixs = duplicates[node]
// 	} else {
// 		return errors.New("'node' must be Branch or Leaf instance")
// 	}
// 	return nil
// }

// // LoadMap deserializes a nested dict representing an RCTree and loads into the RCTree instance.
// // Note that this will delete all data in the current RCTree and replace it with the loaded data.
// func (rcTree RCTree) LoadMap(obj) {
// 	// Create anchor node
// 	anchor := Branch(q=None, p=None)
// 	// Create dictionary for restoring duplicates
// 	duplicates := {}
// 	// Deserialize json object
// 	rrcf.Deserialize(obj, anchor, duplicates)
// 	// Get root node
// 	root := anchor.l
// 	root.u := nil
// 	// Fill in leaves dict
// 	leaves = {}
// 	for k, v in duplicates.items() {
// 		for i in v {
// 			leaves[i] = k
// 		}
// 	}
// 	// Set root of tree to new root
// 	rrcf.root = root
// 	rrcf.leaves = leaves
// 	// Set number of dimensions based on first leaf
// 	rrcf.ndim = len(next(iter(leaves.values())).x)
// }

// // Deserialise recursively deserializes tree from a nested map
// func (rcTree RCTree) Deserialise(obj, node, duplicates, side="l") {
// 	if obj["type"] == "Branch" {
// 		q = obj["q"]
// 		p = obj["p"]
// 		n = np.int64(obj["n"])
// 		b = np.asarray(obj["b"])
// 		branch = Branch(q=q, p=p, n=n, b=b, u=node)
// 		setattr(node, side, branch)
// 		if "l" in obj {
// 			self._deserialize(obj["l"], branch, duplicates, side="l")
// 		}
// 		if "r" in obj {
// 			self._deserialize(obj["r"], branch, duplicates, side="r")
// 		}
// 	} else if obj["type"] == "Leaf" {
// 		i = obj["i"]
// 		x = np.asarray(obj["x"])
// 		d = obj["d"]
// 		n = np.int64(obj["n"])
// 		leaf = Leaf(i=i, x=x, d=d, n=n, u=node)
// 		setattr(node, side, leaf)
// 		duplicates[leaf] = obj["ixs"]
// 	} else {
// 		raise TypeError("'type' must be Branch or Leaf")
// 	}
// }

// // FromMap deserializes a nested map representing an RCTree
// // and creates a new RCTree instance from the loaded data
// func (rcTree RCTree) FromMap() {
// 	newinstance = cls()
// 	newinstance.LoadMap(obj)
// 	return newinstance
// }
