import React from 'react'
import LetterTile from "./tile/LetterTile";
import BonusTile from "./tile/BonusTile";
import EmptyTile from "./tile/AvailableTile";
import '../css/board.css';

import KeyboardDoubleArrowDownIcon from '@mui/icons-material/KeyboardDoubleArrowDown';
import ShuffleIcon from '@mui/icons-material/Shuffle';
import SwapVertIcon from '@mui/icons-material/SwapVert';
import {
    Box,
    Button,
    Dialog, DialogActions,
    DialogContent, DialogContentText,
    DialogTitle,
    IconButton,
    Typography
} from "@mui/material";
import defAxios from "../Http";
import GameInvite from "./GameInvite";
import JoinGameDialog from "./dialog/JoinGameDialog";
import OnePlayerDialog from "./dialog/OnePlayerDialog";
import Board from "./Board";
import SwapTilesDialog from "./dialog/SwapTilesDialog";


const dialog = {
    ONE_PLAYER: "one_player",
    JOIN: "join",
    SWAP_TILES: "swap_tiles",
    BLANK_TILE: "blank_tiles"
}

class GameWindow extends React.Component {
    channel;
    state = {
        openDialog: null,
        debug: false,
        showOnePlayerWarning: false,
        showJoinWarning: false,
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
        this.toggleDebug = this.toggleDebug.bind(this);
        this.startGame = this.startGame.bind(this);
        this.joinGame = this.joinGame.bind(this);
        this.startSpectating = this.startSpectating.bind(this);
        this.backToLobby = this.backToLobby.bind(this);
        this.joinGame = this.joinGame.bind(this);
        this.startSpectating.bind(this);
    }

    componentDidMount(){
       if(this.props.user)this.initLobby();
        this.channel = this.props.realtime.channels.get(`lobby-${this.props.lobbyId}`);
        console.log("Subscribing to messages");
        this.channel.subscribe(this.onMessage)
    }

    async initLobby(){
        const {data: lobby} = await defAxios.get(`lobby/${this.props.lobbyId}`, {validateStatus: ()=>true});
        if(lobby.err){
            this.setState({})
            return;
        }
        const {data: boards} = await defAxios.get(`game/${this.props.lobbyId}/boards`);
        if(boards.deck) {
            boards.swap = {
                width: boards.deck.width,
                height: boards.deck.height,
                squares: boards.deck.squares.map(()=>({tile: null}))
            }
        }
        this.setState({boards, lobby})
    }

    async fetchBoards(){
        const {data: boards} = await defAxios.get(`game/${this.props.lobbyId}/boards`);
        this.setState({boards});
    }

    componentWillUnmount() {
        this.channel.unsubscribe(this.onMessage);
    }

    componentDidUpdate(prevProps, prevState, snapshot){
        console.log("Component update");
        if(this.props.user && prevProps.user?.id !== this.props.user.id){
            this.initLobby();
        }
        console.log(prevProps.user, this.props.user);
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
    render() {
        switch(this.state.lobby?.state){
            case "inGame":
                return this.renderInGame();
            case "waiting":
                return this.renderWaiting();
            case undefined:
                return this.renderLoading();
            default:
                return <div>Unknown state {this.state.lobby.state}</div>
        }

    }

    renderLoading(){
        return <Box sx={{flexGrow: 1}}>
            <Typography align="center">Loading...</Typography>
            <JoinGameDialog open={this.state.openDialog === dialog.JOIN}/>
        </Box>
    }

    renderWaiting(){
        return <Box sx={{flexGrow: 1}}>
            <Typography align="center" variant="h5">Waiting for players ({this.state.lobby.currentPlayers}/{this.state.lobby.maxPlayers})</Typography>
            <Typography align="center">Invite players with this URL:</Typography>
            <Typography align="center">
                <GameInvite lobbyId={this.props.lobbyId}/>
            </Typography>
            <Typography align="center" variant={"body2"}>
                {this.state.lobby.creatorId === this.props.user.id ? <Button onClick={this.startGame}>Start</Button> : "Waiting for the lobby owner to start the game."}
            </Typography>
            <OnePlayerDialog open={this.state.openDialog === dialog.ONE_PLAYER} startGame={this.startGame} clearDialog={this.clearDialog}/>
        </Box>
    }

    renderInGame(){
        return <div id="gameWindow">
            <Board handleTileDrop={this.handleTileDrop} board={this.state.boards["main"]} name={"main"} debug={this.state.debug}/>
            <div id="boardControls">
                <div>
                    <Button disabled={this.isTurn()} onClick={this.play}>Play</Button>
                    <Button disabled={this.isTurn()} onClick={this.pass}>Pass</Button>
                    <Button onClick={this.toggleDebug}>Debug</Button>
                </div>
                <IconButton title="Recall" onClick={this.recallTiles}><KeyboardDoubleArrowDownIcon/></IconButton>
                <IconButton title="Shuffle" onClick={this.shuffleTiles}><ShuffleIcon/></IconButton>
                <IconButton title="Swap" onClick={this.swapTiles}><SwapVertIcon/></IconButton>
            </div>
            <Board handleTileDrop={this.handleTileDrop} board={this.state.boards["deck"]} name={"deck"} debug={this.state.debug}/>
            <SwapTilesDialog open={true} keepDeck={this.state.boards.deck} swapDeck={this.state.boards.swap} debug={this.state.debug} handleTileDrop={this.handleTileDrop}/>
        </div>
    }


    async handleTileDrop(from, fromIndex, to, toIndex){
        if(!this.isTurn() && (from === "board" || to === "board"))return;
        console.log(`Moving ${from}#${fromIndex} -> ${to}#${toIndex}`);
        if(to !== "swap" && from !== "swap") {
            let result = await defAxios.patch(`game/${this.props.lobbyId}/boards`, {from, fromIndex, to, toIndex})
            if (result.data.err) return console.log("Couldn't move tile", result.data.err);
        }
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

    toggleDebug(){
        this.setState({debug: !this.state.debug})
    }

    clearDialog(){
        this.setState({openDialog: null});
    }

    async startGame(){
        if(!this.state.showOnePlayerWarning && this.state.lobby.currentPlayers === 1){
            return this.setState({showOnePlayerWarning: true});
        }

        let {data: lobby} = await defAxios.patch(`lobby/${this.props.lobbyId}`, {
            state: "inGame",
        })

        return this.setState({showOnePlayerWarning: false, lobby});
    }

    backToLobby(){

    }

    startSpectating(){
    }

    async joinGame(){
        await defAxios.put(`lobby/${this.props.lobbyId}/member`, {type: "player"});
        this.initLobby();
    }

}

export default GameWindow