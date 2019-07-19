package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"time"

	randy "github.com/Pallinder/go-randomdata"
	log "github.com/sirupsen/logrus"
)

var searchSpace = flag.Int("space", 100000, "size of search space")
var buffalos = flag.Int("buffalos", 500, "number of buffalos")

func main() {
	flag.Parse()
	rand.Seed(time.Now().Unix())

	numBuffalos := *buffalos
	numFalsePositives := int(float64(numBuffalos) * 0.05)
	searchSpace := *searchSpace

	// Generate
	log.Infof("generating %d items with %d buffalos", searchSpace, numBuffalos)
	values := make([]string, 0, searchSpace)
	for i := 0; i < searchSpace; i++ {
		val := randy.Paragraph()
		if i < numBuffalos {
			// Throw in a buffalo
			randLoc := rand.Intn(len(val))
			newVal := val[:randLoc] + " buffalo " + val[randLoc:]
			val = newVal
		} else if i < numBuffalos+numFalsePositives {
			// Throw in a buffalo false positive (no spaces)
			randLoc := rand.Intn(len(val))
			newVal := val[:randLoc] + "buffalo" + val[randLoc:]
			val = newVal
		}
		values = append(values, val)
	}

	// Shuffle and output
	log.Info("shuffling")
	rand.Shuffle(len(values), func(i, j int) { values[i], values[j] = values[j], values[i] })
	log.Info("writing")
	for _, val := range values {
		b, err := json.Marshal(Value{val})
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}
}

type Value struct {
	Value string `json:"value"`
}
