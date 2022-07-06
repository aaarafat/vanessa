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

// Return the table as list of coordinates
func (OTable *ObstacleTable) GetTable() ([]string,int) {
	var table []string
	for k , _ := range OTable.table {
		table = append(table,k)
	}
	return table , len(table)
}

func (OTable *ObstacleTable) Print() {
	log.Println("Printing Obstacle Table")
	for k , v := range OTable.table {
		log.Println(k,v)
	}
}

