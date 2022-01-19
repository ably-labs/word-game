import React from 'react'
import LetterTile from "./tile/LetterTile";
import BonusTile from "./tile/BonusTile";
import EmptyTile from "./tile/AvailableTile";
import '../css/board.css';

import KeyboardDoubleArrowDownIcon from '@mui/icons-material/KeyboardDoubleArrowDown';
import ShuffleIcon from '@mui/icons-material/Shuffle';
import SwapVertIcon from '@mui/icons-material/SwapVert';
import {Button, IconButton} from "@mui/material";
import defAxios from "../Http";


class GameWindow extends React.Component {


    channel;
    state = {
        boards: {},
        lobby: {},
    }

    constructor(props){
        super(props);
        this.handleTileDrop = this.handleTileDrop.bind(this);
        this.recallTiles = this.recallTiles.bind(this);
        this.shuffleTiles = this.shuffleTiles.bind(this);
        this.swapTiles = this.swapTiles.bind(this);
        this.onMessage = this.onMessage.bind(this);
        this.play = this.play.bind(this);
        this.pass = this.pass.bind(this);
    }

    async componentDidMount(){
        const {data: lobby} = await defAxios.get(`lobby/${this.props.lobbyId}`);
        const {data: boards} = await defAxios.get(`game/${this.props.lobbyId}/boards`);
        this.setState({boards, lobby})
        this.channel = this.props.realtime.channels.get(`lobby-${this.props.lobbyId}`);
        console.log("Subscribing to messages");
        this.channel.on("attach", console.log);
        console.log(this.channel.state);
        this.channel.subscribe(this.onMessage)
    }

    async fetchBoards(){
        const {data: boards} = await defAxios.get(`game/${this.props.lobbyId}/boards`);
        this.setState({boards});
    }

    componentWillUnmount() {
        this.channel.unsubscribe(this.onMessage);
    }

    onMessage(message){
        console.log(message);
        switch(message.name){
            case "moveTile":
                const {move, tile} = message.data;
                console.log(move, tile);
                this.setState((state)=>{
                    if(move.to !== "deck")
                        state.boards[move.to].squares[move.toIndex].tile = tile;
                    if(move.from !== "deck")
                        state.boards[move.from].squares[move.fromIndex].tile = null;
                    return {boards: state.boards}
                })
                break;
        }
    }

    isTurn(){
        return this.state.lobby.playerTurnId !== this.props.user.id
    }

    renderBoard(name){
        const board = this.state.boards[name];
        if(!board)return <div>Board {name} not found</div>
        const rows = [];
        for(let y = 0; y < board.height; y++){
            let cols = [];
            for(let x = 0; x < board.width; x++){
                const i = (y*board.width)+x;
                const square = board.squares[i];
                cols.push(this.drawSquare(square, i, name))
            }
            rows.push(<tr key={`board-${name}-${y}`}>{cols}</tr>)
        }
        return <table className={`board ${name}`}>
            <tbody>
            {rows}
            </tbody>
        </table>
    }

    render() {
        return <div id="gameWindow">
            {this.renderBoard("main")}
            <div id="boardControls">
                <div>
                    <Button disabled={this.isTurn()} onClick={this.play}>Play</Button>
                    <Button disabled={this.isTurn()} onClick={this.pass}>Pass</Button>
                </div>
                <IconButton title="Recall" onClick={this.recallTiles}><KeyboardDoubleArrowDownIcon/></IconButton>
                <IconButton title="Shuffle" onClick={this.shuffleTiles}><ShuffleIcon/></IconButton>
                <IconButton title="Swap" onClick={this.swapTiles}><SwapVertIcon/></IconButton>
            </div>
            {this.renderBoard("deck")}
        </div>
    }

    drawSquare(square, i, source){
        const key = `${source}${i}`
        if(square?.tile) {
            return <LetterTile {...square.tile} index={i} source={source} key={key}/>
        }
        else if(square?.bonus)
            return <BonusTile {...square.bonus} index={i} onTileDropped={this.handleTileDrop} source={source} key={key}/>
        else
            return <EmptyTile index={i} onTileDropped={this.handleTileDrop} source={source} key={key}/>
    }

    async handleTileDrop(from, fromIndex, to, toIndex){
        if(!this.isTurn() && (from === "board" || to === "board"))return;
        console.log(`Moving ${from}#${fromIndex} -> ${to}#${toIndex}`);
        let result = await defAxios.patch(`game/${this.props.lobbyId}/boards`, {from, fromIndex, to, toIndex})
        if(result.data.err)return console.log("Couldn't move tile", result.data.err);
        this.setState((state)=>{
            state.boards[to].squares[toIndex].tile = state.boards[from].squares[fromIndex].tile
            state.boards[from].squares[fromIndex].tile = null;
            return {boards: state.boards}
        })
    }

    // Recalls all the tiles that are currently on the board but have not been played
    recallTiles(){
        this.setState((state)=>{
            for(let i = 0; i < state.boards.main.length; i++){
                const tile = state.boards.main.squares[i].tile;
                if(!tile || !tile.draggable)continue;
                // Find the first empty square
                const newIndex = state.boards.deck.findIndex((s)=>!s.tile);
                state.boards.deck.squares[newIndex].tile = state.boards.main.squares[i].tile;
                state.boards.main.squares[i].tile = null;
            }
            return {boards: state.boards};
        })
    }

    shuffleTiles(){
        this.setState((state)=>{
            // This shuffle does not take into account blank tiles, but it will do for now
            state.boards.deck.squares.sort(()=>Math.random() > 0.5 ? 1 : -1)
            return {boards: state.boards}
        })
    }

    async play(){
       let result = await defAxios.post(`game/${this.props.lobbyId}/boards`);
       console.log(result);
       await this.fetchBoards();
    }

    pass(){
        this.recallTiles();
        this.play();
    }

    swapTiles(){
        // TODO
    }

}

export default GameWindow