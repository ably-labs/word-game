import React from 'react'
import LetterTile from "./tile/LetterTile";
import BonusTile from "./tile/BonusTile";
import EmptyTile from "./tile/AvailableTile";
import '../css/board.css';

import KeyboardDoubleArrowDownIcon from '@mui/icons-material/KeyboardDoubleArrowDown';
import ShuffleIcon from '@mui/icons-material/Shuffle';
import SwapVertIcon from '@mui/icons-material/SwapVert';
import {IconButton} from "@mui/material";

// Board layout for testing
const layout =
    "T--D---T---D--T" +
    "-W---L---L---W-" +
    "--W---D-D---W--" +
    "D--W---D---W--D" +
    "----W-----W----" +
    "-L---L---L---L-" +
    "--D---D-D---D--" +
    "T--D---S---D--T" +
    "--D---D-D---D--" +
    "-L---L---L---L-" +
    "----W-----W----" +
    "D--W---D---W--D" +
    "--W---D-D---W--" +
    "-W---L---L---W-" +
    "T--W---T---W--T"

const tileTypes = {
    "-": {},
    "T": {bonus: {text: "TRIPLE WORD", type: "triple-word"}},
    "D": {bonus: {text: "DOUBLE LETTER", type: "double-letter"}},
    "W": {bonus: {text: "DOUBLE WORD", type: "double-word"}},
    "L": {bonus: {text: "TRIPLE LETTER", type: "triple-letter"}},
    "S": {bonus: {text: "START", type: "double-word"}}
}
const GAME_BOARD_WIDTH = 15;
const GAME_BOARD_HEIGHT = 15;
const GAME_DECK_LENGTH = 9;

class GameWindow extends React.Component {

    state = {
        boards: {
            main: [],
            deck: []
        },
    }

    constructor(props){
        super(props);
        this.handleTileDrop = this.handleTileDrop.bind(this);
        this.recallTiles = this.recallTiles.bind(this);
        this.shuffleTiles = this.shuffleTiles.bind(this);
        this.swapTiles = this.swapTiles.bind(this);
    }

    componentDidMount(){
        this.constructBoard();
    }

    // TODO: The board data will eventually come from the server, but for testing this is a layout
    constructBoard(){
        let main = Array(GAME_BOARD_WIDTH * GAME_BOARD_HEIGHT).fill({});
        for(let i = 0; i < layout.length; i++){
            // Create a copy of the object
            main[i] = Object.create(tileTypes[layout[i]])
        }
        main[10].tile = {letter: "A", score: 1, draggable: false};
        let deck = Array(GAME_DECK_LENGTH).fill(0).map(()=>({}))
        deck[0].tile = {letter: "P", score: 1, draggable: true}
        deck[3].tile = {letter: "E", score: 1, draggable: true}
        deck[2].tile = {letter: "L", score: 1, draggable: true}
        deck[1].tile = {letter: "P", score: 1, draggable: true}
        this.setState({
            boards: {main, deck}
        });
    }

    render() {
        const rows = [];
        for(let y = 0; y < GAME_BOARD_HEIGHT; y++){
            let cols = [];
            for(let x = 0; x < GAME_BOARD_WIDTH; x++){
                const i = (y*GAME_BOARD_WIDTH)+x;
                const square = this.state.boards.main[i];
                cols.push(this.drawSquare(square, i, "main"))
            }
            rows.push(<tr>{cols}</tr>)
        }
        return <div id="gameWindow">
            <table className="board">
                {rows}
            </table>
            <div id="boardControls">
                <IconButton title="Recall" onClick={this.recallTiles}><KeyboardDoubleArrowDownIcon/></IconButton>
                <IconButton title="Shuffle" onClick={this.shuffleTiles}><ShuffleIcon/></IconButton>
                <IconButton title="Swap" onClick={this.swapTiles}><SwapVertIcon/></IconButton>
            </div>
            <table className="deck">
                <tr>{this.state.boards.deck.map((e, i)=>this.drawSquare(e,i, "deck"))}</tr>
            </table>
        </div>
    }

    drawSquare(square, i, source){
        if(square?.tile)
            return <LetterTile {...square.tile} index={i} source={source}/>
        else if(square?.bonus)
            return <BonusTile {...square.bonus} index={i} onTileDropped={this.handleTileDrop} source={source}/>
        else
            return <EmptyTile index={i} onTileDropped={this.handleTileDrop} source={source}/>
    }

    handleTileDrop(from, fromIndex, to, toIndex){
        console.log(`Moving ${from}#${fromIndex} -> ${to}#${toIndex}`);
        this.setState((state)=>{
            state.boards[to][toIndex].tile = state.boards[from][fromIndex].tile
            state.boards[from][fromIndex].tile = null;
            return {boards: state.boards}
        })
    }

    // Recalls all the tiles that are currently on the board but have not been played
    recallTiles(){
        this.setState((state)=>{
            for(let i = 0; i < state.boards.main.length; i++){
                const tile = state.boards.main[i].tile;
                if(!tile || !tile.draggable)continue;
                // Find the first empty square
                const newIndex = state.boards.deck.findIndex((s)=>!s.tile);
                state.boards.deck[newIndex].tile = state.boards.main[i].tile;
                state.boards.main[i].tile = null;
            }
            return {boards: state.boards};
        })
    }

    shuffleTiles(){
        this.setState((state)=>{
            // This shuffle does not take into account blank tiles, but it will do for now
            state.boards.deck.sort(()=>Math.random() > 0.5 ? 1 : -1)
            return {boards: state.boards}
        })
    }

    swapTiles(){
        // TODO
    }

}

export default GameWindow