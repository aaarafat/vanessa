package rsu

import (
	"fmt"
	"log"
)

type Position struct {
	Lat float64
	Lng float64
}

type ObstaclesTable struct {
	table map[string] uint8

}
func NewObstaclesTable() *ObstaclesTable {
	return &ObstaclesTable{
		table : make(map[string] uint8),
	}
}

func (OTable *ObstaclesTable) Set(x float64, y float64,clear uint8) {
	
	key := fmt.Sprintf("%f,%f",x,y)
	if(clear==1){
		delete(OTable.table,key)
	}else if(clear==0){
		OTable.table[key] = clear
	}
	
}

// Return the table as list of coordinates as a pair of x,y
func (OTable *ObstaclesTable) GetTable() ([]Position,int) {
	var table [] Position
	var count int
	for k , v := range OTable.table {
		if v==1{
			count++
			var x,y float64
			fmt.Sscanf(k,"%f,%f",&x,&y)
			pos := &Position{x,y}
			table = append(table,*pos)
		}
	}
	return table,count
}

func (OTable *ObstaclesTable) Print() {
	log.Println("Printing Obstacle Table")
	for k , v := range OTable.table {
		log.Println(k,v)
	}
}

