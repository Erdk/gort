package rayengine

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

func readObj(filename string) (*bvhNode, error) {

	red := newLambertianRGB(0.85, 0.05, 0.05)

	objfile, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("can't open file")
	}
	defer objfile.Close()

	scanner := bufio.NewScanner(objfile)
	vertexList := []*Vec{}
	hlist := []hitable{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		switch fields[0] {
		case "v":
			// process vertex
			x, _ := strconv.ParseFloat(fields[1], 64)
			y, _ := strconv.ParseFloat(fields[2], 64)
			z, _ := strconv.ParseFloat(fields[3], 64)
			v := &Vec{x*39 + 260, y * 39, z*39 + 250}
			vertexList = append(vertexList, v)

		case "f":
			v1, _ := strconv.Atoi(strings.Split(fields[1], "/")[0])
			v2, _ := strconv.Atoi(strings.Split(fields[2], "/")[0])
			v3, _ := strconv.Atoi(strings.Split(fields[3], "/")[0])

			hlist = append(hlist, &triangle{
				vertexList[v1-1], vertexList[v2-1], vertexList[v3-1], red, true,
			})
		}
	}

	obj := bvhNodeInit(hlist, len(hlist), 0.0, 1.0)

	return obj, errors.New("can't load obj file")
}
