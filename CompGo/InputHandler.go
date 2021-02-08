package main

import(
	"os"
	"bufio"
	"log"
	"strings"
	"strconv"
	"fmt"
)


//Counts the number of lines in  file and returns it
func countLines(fileName string) int{
	file, err := os.Open(fileName)
	if err!=nil{
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
	    lineCount++
	}
	return lineCount
}
//Reads a grid from a file and splits into 3 string grids (first Line, middle part, and last line) and returns them 
func compileStringGrid(fileName string)([]string, [][]string, []string){
	file, err := os.Open(fileName)
  demand:=make([]string,0)
	if err!=nil{
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(file)
	numOfLines:= countLines(fileName)-2
	fmt.Println(numOfLines)
	fileScanner.Scan()
	line:=fileScanner.Text()
  line=strings.TrimSpace(line)
  firstLine:=strings.Split(line, " ")
	numOfCols:= len(strings.Split(line," "))
	grid := make([][]string,numOfLines)
	i:=0
	for fileScanner.Scan(){
		line = fileScanner.Text()
		lineArr := strings.Split(line, " ")
		if i==numOfLines{
			demand= lineArr
		} else{
			if(len(lineArr)!=numOfCols){
				//Error
			}
			grid[i]=lineArr
		}
		i++
	}
	return firstLine, grid, demand
}

//Parses the string grids into a grid of factory-warehouse routes, intialized to the specifications given in the file and returns it
func parseStringGrid(grid [][]string, demand []string) [][]*Route{
	numOfLines:=len(grid)
	numOfCols:=len(grid[0])
	factories:= make([]*Factory, numOfLines)
	warehouses:= make([]*Warehouse, numOfCols-2)
	for i:=0; i<numOfLines; i++{
		str:= grid[i][numOfCols-1]
		str = strings.TrimSpace(str)
		num,_:=strconv.Atoi(str)
		factories[i] = newFactory(num)
	}

	for i:=0; i<numOfCols-2; i++{
		str:= demand[i+1]
		str = strings.TrimSpace(str)
		num,_:=strconv.Atoi(str)
		warehouses[i] = newWarehouse(num)
	}
	cols:= len(warehouses)
	rows:= len(factories)
	routes:=make([][]*Route, rows)
	for i:=0; i<rows; i++{
		routes[i]=make([]*Route, cols)
	}
	for j:=0; j<cols; j++{
		for i:=0; i<rows; i++{
			str:= grid[i][j+1]
			str= strings.TrimSpace(str)
			num,_:=strconv.Atoi(str)
			route:= newRoute(factories[i], warehouses[j], num)
      routes[i][j]=route
		}
	}
	return routes
}
