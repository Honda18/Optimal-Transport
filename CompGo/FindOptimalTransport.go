package main

import(
	"os"
	"math"
	"fmt"
  "sync"
  "strconv"
  "log"
  "strings"
)

var(
	minTrail [][]int = make([][]int,0)
	trail  [][]int = make([][]int,0)
	routesRow int
	routesCol int
	routes [][]*Route
    wg sync.WaitGroup
)

//Runs one round of the stepping stone method and returns the minimum marginal cost obtained from all the path calculations
func optimize() int{
	minMarginalCost:= math.MaxInt32
	trail = make([][]int,0)
	for j:=0; j<routesCol; j++{
		for i:=0; i<routesRow; i++{
      wg.Add(1)
      go func(x,y int){
        route:= routes[x][y]
        mCost:= math.MaxInt32
        if route.unused(){
          mCost, trail = marginalCost(y,x)
        }

        if mCost<minMarginalCost {
          minMarginalCost = mCost
          minTrail=trail
        }
        wg.Done()
      }(i,j)
		}
	}
  wg.Wait()
	return minMarginalCost

} 

//A wrapper function(to be used concurrently) that recieves an empty cell(by indices i,j) and returns its closed path as well as its marginal cost
func marginalCost(col,row int) (int, [][]int){
  trail:=getClosedPath(col,row)
  mCost:=getPathCost(trail)
  return mCost,trail
}

//Given a path it calculates the marginal cost
func getPathCost(trail [][]int) int{
	cost:=0
	sign:=1
	for i:=0; i<len(trail)-1; i++{
		col:=trail[i][1]
		row:=trail[i][0]
		route:=routes[row][col]
		cost+=route.getCost()*sign
		sign= sign*-1
		}
	return cost 
}

//Updates the allocation of the units based on the newly obtained path with minimum cost
func updatAllocation(){
	min:= math.MaxInt32
	for i:=1; i<len(minTrail); i+=2{
		indices := minTrail[i]
		route := routes[indices[0]][indices[1]]
		min= int(math.Min(float64(min), float64(route.getUnitsAllocated())))
	}

	for i:= 0; i<len(minTrail)-1; i++{
		indices:= minTrail[i]
		route := routes[indices[0]][indices[1]]
		if i%2==0{
			route.allocateToRoute(min)
		} else{
			route.decrementRoute(min)
		}
	}
}

// Returns the closed path of an empty cel given its indices 
func getClosedPath(terminalCol int, terminalRow int) ([][]int){
	trail:=make([][]int,0)
	flags:= make([][]int,0)
	for i:=0; i<routesRow; i++{
		arr:=make([]int,routesCol)
		flags=append(flags,arr)
	}
	col:=terminalCol
	row:=terminalRow
	size:=1
	sign:=-1
	trail= append(trail, []int{terminalRow, terminalCol})
	for trail[size-1][0]!=terminalRow || trail[size-1][1]!=terminalCol ||size==1{
		emptyCell:=true
		for emptyCell{
			if sign==-1{
				col = (col+1)%routesCol
			} else{
				row = (row+1)%routesRow
			}
			route:= routes[row][col]
			if(col==terminalCol&&row==terminalRow){
				trail=append(trail,[]int{row,col})
				size++
				break
			}
			checkPointCol:=trail[size-1][1]
			checkPointRow:=trail[size-1][0]
			if col==checkPointCol&&row==checkPointRow {
				indices:=trail[size-1]
				trail=trail[:size-1]
				flags[indices[0]][indices[1]]=0
				size--
				sign=sign*-1
				break
			}
			if !route.unused()&&(flags[row][col]==0) {
				emptyCell=false
				size++
				trail = append(trail, []int{row,col})
				sign=sign*-1
				flags[row][col]=1
			}
		}
	}

	for i:=0; i<len(trail); i++{
		indices:=trail[i]
		flags[indices[0]][indices[1]]=0
	} 

	return trail
}

//Prints the unit allocation of the routes as a grid
func printAllocation(){
	for j:=0; j<routesRow; j++{
		for i:=0; i<routesCol; i++{
			units:=routes[j][i].getUnitsAllocated()
			fmt.Print(units)
			fmt.Print(" ")
		}
		fmt.Println()
	}
}



//A setter function for the routes, to be used to set the routes after the input has been handled
func setRoutes(newRoutes [][]*Route){
	routes = newRoutes
	routesRow = len(routes)
	routesCol = len(routes[0])
}

// A function that updates the routes with the initial solution grid obtained by input handling
func setInitialSolution(grid [][]string){
  numRow:= len(grid)
  numCol:= len(grid[0])-2
  for i:=0; i<numRow; i++{
    for j:=0; j<numCol; j++{
      str:=grid[i][j+1]
      num, err := strconv.Atoi(str)
      if err==nil {
        routes[i][j].allocateToRoute(num)
      }
    }
  }
}

//A function that writes the optimal solution to a file called Solutions.txt, the inputs are templates extracted from the initial input files the writing format
func writeSolToFile(initialLine []string, grid [][]string, demand []string){
  f, err := os.Create("solution.txt")
    if err != nil {
        fmt.Println(err)
                f.Close()
        return
    }
    
    for i:=0; i<len(grid); i++{
      for j:=1; j<len(grid[0])-1; j++{
        grid[i][j]=strconv.Itoa(routes[i][j-1].getUnitsAllocated())
      }
    }

    str:=strings.Join(initialLine, " ")
    fmt.Fprintln(f,str)
    for _, v := range grid {
        str= strings.Join(v, " ")
        fmt.Fprintln(f, str)
        if err != nil {
            fmt.Println(err)
            return
        }
    }
    str = strings.Join(demand," ")
    fmt.Fprintln(f, str)
    err = f.Close()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Optimal configuration written succesfully to Solutions.txt")
}

func main(){
	if(len(os.Args)<=2){
		log.Fatal("Insufficient number of input files!")
	}
	fileName:=os.Args[1]
	_,grid,demand:=compileStringGrid(fileName)
	setRoutes(parseStringGrid(grid,demand))
    fileName2:=os.Args[2]
    initialLine,intialsGrid,demands:=compileStringGrid(fileName2)
    setInitialSolution(intialsGrid)
    fmt.Println("Initial Solution: ")
    printAllocation()
    fmt.Println()
	minCost:=optimize()
	for minCost<0 {
		updatAllocation()
		minCost=optimize()
	}
	fmt.Println()
	fmt.Println("Optimal Solution: ")
	printAllocation()
    writeSolToFile(initialLine,intialsGrid,demands)
}