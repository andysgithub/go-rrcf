# go-rrcf

## Robust random cut forest implemented in Go

Based on the Python rrcf project at https://github.com/kLabUM/rrcf

## About

The *Robust Random Cut Forest* (RRCF) algorithm is an ensemble method for detecting outliers in streaming data. RRCF offers a number of features that many competing anomaly detection algorithms lack:

- Designed to handle streaming data
- Performs well on high-dimensional data
- Reduces the influence of irrelevant dimensions
- Gracefully handles duplicates and near-duplicates that could otherwise mask the presence of outliers
- Features an anomaly-scoring algorithm with a clear underlying statistical meaning

## Robust random cut trees

A robust random cut tree (RRCT) is a binary search tree that can be used to detect outliers in a point set. A RRCT can be instantiated from a point set. Points can also be added and removed from an RRCT.

### Creating the tree

```go
import (
	"github.com/andysgithub/go-rrcf/random"
)
    // A random cut tree can be instantiated from a point set (n x d)
    rnd := random.NewRandomState(0)
	X = rnd.Normal2D(100, 2)
	tree = rrcf.NewRCTree(X, nil, 9, nil)

    // A random cut tree can also be instantiated with no points
    tree = rrcf.NewRCTree(nil, nil, 0, nil)
```

### Inserting points

```go
rnd := random.NewRandomState(0)
tree = rrcf.NewRCTree(nil, nil, 9, nil)

for _, index := range indexes {
    x := rnd.Normal1D(2)
    leafNode, err := tree.InsertPoint(x, index, 0)
}
```

### Deleting points

```go
deletedNode := tree.ForgetPoint(2)
```

## Batch anomaly detection

An example to detect outliers in a batch setting is included in main.go. Results of this trial can be found in the results/batch folder.

![Image](/home/andy/go/src/github.com/andysgithub/go-rrcf/results/batch/plot.png) 

## Streaming anomaly detection

An example to detect anomalies in streaming time series data is also in main.go. Results can be found in the results/streaming folder.

![Image](/home/andy/go/src/github.com/andysgithub/go-rrcf/results/streaming/plot.png) 