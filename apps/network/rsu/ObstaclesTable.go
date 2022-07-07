package rsu

import (
	"log"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)


type ObstaclesTable struct {
	table map[string] uint8

}
func NewObstaclesTable() *ObstaclesTable {
	return &ObstaclesTable{
		table : make(map[string] uint8),
	}
}

func (OTable *ObstaclesTable) Set(position Position,clear uint8) {
	
	key := string(position.Marshal())
	if(clear==1){
		delete(OTable.table,key)
	}else if(clear==0){
		OTable.table[key] = clear
	}
	
}

// Return the table as list of coordinates as a pair of x,y
func (OTable *ObstaclesTable) GetTable() []Position {
	var table [] Position
	for k  := range OTable.table {
			pos := UnmarshalPosition([]byte(k))
			table = append(table,pos)
	}
	return table
}

func (OTable *ObstaclesTable) Print() {
	log.Println("Printing Obstacle Table")
	for k , v := range OTable.table {
		pos := UnmarshalPosition([]byte(k))
		log.Println(pos.Lat,pos.Lng,v)
	}
}

