import React from 'react'
import LetterTile from "./tile/LetterTile";
import BonusTile from "./tile/BonusTile";
import EmptyTile from "./tile/AvailableTile";
import '../css/board.css';

import KeyboardDoubleArrowDownIcon from '@mui/icons-material/KeyboardDoubleArrowDown';
import ShuffleIcon from '@mui/icons-material/Shuffle';
import SwapVertIcon from '@mui/icons-material/SwapVert';
import {IconButton} from "@mui/material";
import defAxios from "../Http";

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
        mainSize: {
            width: 15,
            height: 15,
        },
        deckSize: {
            width: 9,
            height: 1, // TODO: Deck width doesn't do anything yet
        },
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
        this.getDeck();
        this.getBoard();
    }

    async getDeck(){
        const {data: {squares: deck, width, height}} = await defAxios.get(`game/${this.props.lobbyId}/deck`);
        this.setState({
            deckSize: {width, height},
            boards: {
                main: this.state.boards.main,
                deck
            }
        })
    }


    async getBoard(){
        const {data: {squares: main, width, height}} = await defAxios.get(`game/${this.props.lobbyId}/board`);
        this.setState({
            mainSize: {width, height},
            boards: {
                main: main,
                deck: this.state.boards.deck
            }
        })
    }

    renderBoard(width, height, board, className, source){
        const rows = [];
        for(let y = 0; y < height; y++){
            let cols = [];
            for(let x = 0; x < width; x++){
                const i = (y*width)+x;
                const square = board[i];
                cols.push(this.drawSquare(square, i, source))
            }
            rows.push(<tr key={`board-${source}-${y}`}>{cols}</tr>)
        }
        return <table className={className}>
            <tbody>
            {rows}
            </tbody>
        </table>
    }

    render() {
        return <div id="gameWindow">
            {this.renderBoard(this.state.mainSize.width, this.state.mainSize.height, this.state.boards.main, "board", "main")}
            <div id="boardControls">
                <IconButton title="Recall" onClick={this.recallTiles}><KeyboardDoubleArrowDownIcon/></IconButton>
                <IconButton title="Shuffle" onClick={this.shuffleTiles}><ShuffleIcon/></IconButton>
                <IconButton title="Swap" onClick={this.swapTiles}><SwapVertIcon/></IconButton>
            </div>
            {this.renderBoard(this.state.deckSize.width, this.state.deckSize.height, this.state.boards.deck, "deck", "deck")}
        </div>
    }

    drawSquare(square, i, source){
        const key = `${source}${i}`
        if(square?.tile)
            return <LetterTile {...square.tile} index={i} source={source} key={key}/>
        else if(square?.bonus)
            return <BonusTile {...square.bonus} index={i} onTileDropped={this.handleTileDrop} source={source} key={key}/>
        else
            return <EmptyTile index={i} onTileDropped={this.handleTileDrop} source={source} key={key}/>
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