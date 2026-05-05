package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/datumbrain/npy"
)

func main() {
	//r, w := io.Pipe()
	var frames uint
	var joints uint

	flag.UintVar(&frames, "f", 0, "Frame count")
	flag.UintVar(&joints, "j", 0, "Joint count")
	flag.Parse()

	if frames == 0 {
		fmt.Println("No frame count provided. Use the flag -f")
		os.Exit(1)
	}

	if joints == 0 {
		fmt.Println("No joint count provided. Use the flag -f")
		os.Exit(1)
	}

	// Read NPZ file
	npzFile, err := npy.ReadNPZFile("output.npz")
	if err != nil {
		log.Fatalf("Failed to read NPZ file: %v", err)
	}

	jointsPos := getArray[float32](npzFile, "posed_joints")
	rootPositions := getArray[float32](npzFile, "root_positions")
	globalRotMats := getArray[float32](npzFile, "global_rot_mats")

	checkDataLen("posed_joints", len(jointsPos.Data), int(3*frames*joints))
	checkDataLen("root_positions", len(rootPositions.Data), int(3*frames))
	checkDataLen("global_rot_mats", len(globalRotMats.Data), int(3*3*frames*joints))

	// Print data
	//fmt.Printf("Matrix: %v %d\n", jointsPos.Data, len(jointsPos.Data))
	//fmt.Printf("Matrix: %v %d\n", rootPositions.Data, len(rootPositions.Data))
	//fmt.Printf("Matrix: %v %d\n", globalRotMats.Data, len(globalRotMats.Data))

	j, _ := json.MarshalIndent(&jointsPos, "", "\t")
	os.WriteFile("posed_joints.json", j, 0666)
}

func checkDataLen(name string, expected int, got int) {
	if expected != got {
		log.Fatalf("Mismatching %s length, expected %d, got %d", name, expected, got)
	}
}

func getArray[T any](npz *npy.NPZFile, name string) *npy.Array[T] {
	value, ok := npy.Get[T](npz, name)
	if !ok {
		log.Fatalf("array %s not found in NPZ file", name)
	}
	return value
}
