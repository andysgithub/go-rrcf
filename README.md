# go-rrcf

## Robust random cut forest implemented in Go

Based on the original Python rrcf project at https://github.com/kLabUM/rrcf.

This Go version features an easy to use API suitable for implementing as a web service if required. As such the forest initialisation returns a token to be used in subsequent calls. The token identifies the user's forest data, allowing multiple users and forests to be supported.

## The RRCF algorithm

The Robust Random Cut Forest (RRCF) algorithm is an ensemble method for detecting outliers in streaming data. RRCF offers a number of features that many competing anomaly detection algorithms lack:

- Designed to handle streaming data
- Performs well on high-dimensional data
- Reduces the influence of irrelevant dimensions
- Gracefully handles duplicates and near-duplicates that could otherwise mask the presence of outliers
- Features an anomaly-scoring algorithm with a clear underlying statistical meaning

## Robust random cut trees

A robust random cut tree (RRCT) is a binary search tree that can be used to detect outliers in a point set. Each tree can be instantiated from a point set, and points can be added and removed dynamically.

## Batch anomaly detection

A batch file containing multi-dimensional data can be read as a csv file. This is used to initialise the forest, with a given number of trees and leaves on each. A map of anomaly scores is then produced by calling ScoreForest with the returned token:

```go
import (
    "github.com/andysgithub/go-rrcf/utils"
)
    // Get random 3D data with anomalies
    points, _ := utils.ReadFromCsv("data/random3D.csv")

    // Construct a random forest
    token := InitForest(100, 256, points, 0)

    // Compute average anomaly score
    scores := ScoreForest(token)
```

### Test results

The resulting scores can then be used to produce a set of data points for saving as a csv file for further analysis. An example for this is included in main.go, and results can be found in the results/batch folder.

The first plot shows the source data, with outliers occupying the central region of the plot. The results of the anomaly detection are shown in the second plot, with all outliers detected above an anomaly score of 60.

![Image](https://github.com/andysgithub/go-rrcf/raw/master/results/batch/plot.png) 

## Streaming anomaly detection

For use with streaming time series data, the forest is initialised with empty trees by passing nil for the data parameter. Note that the final parameter in InitForest must be > 0, to allow shingling (overlapping of data items) to take place. This is because the sine data is single-valued, not multi-dimensional.

Anomaly scores are then collected by calling UpdateForest for each data point received:

```go
import (
    "github.com/andysgithub/go-rrcf/utils"
)
    // Get sine function data with anomalies
    points, _ := utils.ReadFromCsv("data/sine.csv")

    // Construct a forest of empty trees
    token := InitForest(40, 256, nil, 3)

    // Create a map to store the anomaly score of each point
    scores := make(map[int]float64)

    // For each streamed data point
    for sampleIndex, point := range points {
        // Update the forest with this point and record the average score
        scores[sampleIndex] = UpdateForest(token, sampleIndex, point)
    }
```

### Test results

An example to detect anomalies from streamed data is also in main.go. Results can be found in the results/streaming folder.

Anomalous data is injected into the sine wave function shown in the first plot. These outliers are cleary signalled in the anomaly score output.

![Image](https://github.com/andysgithub/go-rrcf/raw/master/results/streaming/plot.png) 

## Anomaly detection after training

To improve the initial anomaly scores for streaming, the forest can be trained with data having no outliers before real-world data is introduced. Note that this training data will need to be streamed, not presented as a batch file:

```go
import (
    "github.com/andysgithub/go-rrcf/utils"
)
    // Get sine function data for training
    points, _ := utils.ReadFromCsv("data/training.csv")

    // Construct a forest of empty trees
    token := InitForest(40, 256, nil, 3)

    // For each training data point
    for sampleIndex, point := range points {
        // Update the forest with this point
        UpdateForest(token, sampleIndex, point)
    }
    lastIndex := len(points)
```

Streamed data can then be presented to the newly-trained forest, as in the previous example:

```go
    // Get sine function data with anomalies
    points, _ = utils.ReadFromCsv("data/sine.csv")

    // Create a map to store the anomaly score of each point
    scores := make(map[int]float64)

    // For each streamed data point
    for sampleIndex, point := range points {
        // Update the forest with this point and record the average score
        scores[sampleIndex] = UpdateForest(token, lastIndex+sampleIndex, point)
    }
```

### Test results

An example can be found in main.go, with results in the results/training folder. The anomaly scores can be seen to be more clearly defined, compared to the results from an untrained forest.

![Image](https://github.com/andysgithub/go-rrcf/raw/master/results/training/plot.png) 