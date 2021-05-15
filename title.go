package main

import (
	sprite "github.com/pdevine/go-asciisprite"
)

/*
const titleText = `

XXXXXXX          XXXXXXXXX   XXXXXXX   XX  XXXXXXX XXXXXXX    XXXXXX   XXXXXX
XXXXXXX          XXXXXXXXX   XXXXXXX   XX  XXXXXXX XXXXXXX    XXXXXX   XXXXXX
      XX         XX      XX        XX      XX            XX  XX    XX  XX    XX
      XX         XX      XX        XX      XX            XX  XX    XX  XX    XX
XXXXXXX          XX  XX  XX  XXXXXXX   XX  XX      XXXXXXX   XX    XX  XX    XX
XXXXXXX          XX  XX  XX  XXXXXXX   XX  XX      XXXXXXX   XX    XX  XX    XX
XX    XX         XX  XX  XX  XX    XX  XX  XX            XX  XX    XX  XX    XX
XX    XX         XX  XX  XX  XX    XX  XX  XX            XX  XX    XX  XX    XX
XXXXXXX          XX  XX  XX  XXXXXXX   Xx   XX           XX   XXXXXX   XX    XX
XXXXXXX          XX  XX  XX  XXXXXXX   XX    XX          XX   XXXXXX   XX    XX
`
*/

const titleText = `

XXXXXXXXXXX
XXXXXXXXXXXX
XXXXXXXXXXXXX
          XXXX
           XXX
          XXXX
XXXXXXXXXXXX
XXXXXXXXXXXXX
XXXX      XXXX
XXXX       XXX
XXXX      XXXX
XXXXXXXXXXXX
XXXXXXXXXXX
`

const bigBombText = `

       XXXXXXX       
     XXXXXXXXXXX     
   XXXXXXXXXXXXXXX   
  XXXXXXXXXXXXXXXXX  
 XXXXXXXXXXXXXXXXXXX 
XXXXXXXXXXXXXXXXXXXXX
XXXXXXXXXXXXXXXXXXXXX
 XXXXXXXXXXXXXXXXXXX 
  XXXXXXXXXXXXXXXXX  
   XXXXXXXXXXXXXXX   
     XXXXXXXXXXX     
       XXXXXXX       



`


type Title struct {
	sprite.BaseSprite
}

func NewTitle() *Title {
	t := &Title{BaseSprite: sprite.BaseSprite{
		X:       5,
		Y:       11,
		Visible: true},
	}

	//surf := sprite.NewSurfaceFromString(titleText, true)
	surf := sprite.NewSurfaceFromString(bigBombText, true)
	t.BlockCostumes = []*sprite.Surface{&surf}

	return t
}
