package rsu

import (
	"fmt"
	"log"
)



type ObstacleTable struct {
	table map[string] uint8

}
func NewObstacleTable() *ObstacleTable {
	return &ObstacleTable{
		table : make(map[string] uint8),
	}
}

func (OTable *ObstacleTable) Set(x uint32, y uint32,clear uint8) {
	
	key := fmt.Sprintf("%d,%d",x,y)
	if(clear==1){
		delete(OTable.table,key)
	}else if(clear==0){
		OTable.table[key] = clear
	}
	
}

// Return the table as list of coordinates as a pair of x,y
func (OTable *ObstacleTable) GetTable() ([][]uint32,int) {
	var table [][]uint32
	var count int
	for k , v := range OTable.table {
		if v==1{
			count++
			var x,y uint32
			fmt.Sscanf(k,"%d,%d",&x,&y)
			table = append(table,[]uint32{x,y})
		}
	}
	return table,count
}

func (OTable *ObstacleTable) Print() {
	log.Println("Printing Obstacle Table")
	for k , v := range OTable.table {
		log.Println(k,v)
	}
}

