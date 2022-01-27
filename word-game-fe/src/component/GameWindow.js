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
import InvalidPlacementDialog from "./dialog/InvalidPlacementDialog";
import BlankTileDialog from "./dialog/BlankTileDialog";
import LobbyDeleted from "./dialog/LobbyDeletedDialog";
import LobbyDeletedDialog from "./dialog/LobbyDeletedDialog";


const dialog = {
    ONE_PLAYER: "one_player",
    JOIN: "join",
    SWAP_TILES: "swap_tiles",
    BLANK_TILE: "blank_tiles",
    INVALID_PLACEMENT: "invalid_placement",
    LOBBY_DELETED: "lobby_deleted",
}

class GameWindow extends React.Component {
    channel;
    state = {
        openDialog: null,
        debug: false,
        showOnePlayerWarning: false,
        showJoinWarning: false,
        boards: {},
        blankData: {},
        placementError: "",
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
        this.clearDialog = this.clearDialog.bind(this);
        this.setBlankTile = this.setBlankTile.bind(this);
        this.resetGame = this.resetGame.bind(this);
        this.endGame = this.endGame.bind(this);
    }

    componentDidMount(){
        if(!this.props.lobby.id || !this.props.channel)return;
        if(this.props.user)this.initLobby();
        this.props.channel.subscribe(this.onMessage)
    }

    async initLobby(){
        const {data: boards} = await defAxios.get(`game/${this.props.lobby.id}/boards`);
        if(boards.deck) {
            boards.swap = {
                width: boards.deck.width,
                height: boards.deck.height,
                squares: boards.deck.squares.map(()=>({tile: null}))
            }
        }
        this.setState({boards})
    }

    async fetchBoards(){
        const {data: boards} = await defAxios.get(`game/${this.props.lobby.id}/boards`);
        this.setState({boards});
    }

    componentWillUnmount() {
        this.props.channel.unsubscribe(this.onMessage);
    }

    componentDidUpdate(prevProps, prevState, snapshot){
        if(this.props.user && prevProps.user?.id !== this.props.user.id){
            this.initLobby();
        }

        if(prevProps.lobby !== this.props.lobby){
            this.fetchBoards();
            this.props.channel.subscribe(this.onMessage)
        }
    }

    onMessage(message){
        switch(message.name){
            case "moveTile":
                if(this.isTurn())return;
                const {move, tile} = message.data;
                this.setState((state)=>{
                    if(move.to !== "deck")
                        state.boards[move.to].squares[move.toIndex].tile = tile;
                    if(move.from !== "deck")
                        state.boards[move.from].squares[move.fromIndex].tile = null;
                    return {boards: state.boards}
                })
                break;
            case "lobbyDeleted":
                this.setState({openDialog: dialog.LOBBY_DELETED})
        }
    }

    isTurn(){
        return this.props.lobby.state === "inGame" && this.props.lobby.playerTurnId === this.props.user.id
    }
    render() {
        switch(this.props.lobby?.state){
            case "roundOver":
                return this.renderGameOver();
            case "inGame":
                return this.renderInGame();
            case "waiting":
                return this.renderWaiting();
            case undefined:
                return this.renderLoading();
            default:
                return <div>Unknown state {this.props.lobby.state}</div>
        }

    }


    renderGameOver(){
        return <Box sx={{flexGrow: 1}}>
            <Typography align="center" variant="h5">Game Over...</Typography>
            <Typography align="center"><b>{this.props.members.find((m)=>m.id ===this.props.lobby.playerTurnId)?.user.name}</b> is the winner!</Typography>
            <Typography align="center">Start a new game or end the game</Typography>
            {this.props.lobby.creatorId === this.props.user.id ? <Typography align="center">
                <Button onClick={this.resetGame}>Continue</Button>
                <Button onClick={this.endGame}>Quit</Button>
            </Typography> : "Waiting for the lobby owner to make a choice."}


        </Box>
    }

    renderLoading(){
        return <Box sx={{flexGrow: 1}}>
            <Typography align="center">Loading...</Typography>
            <JoinGameDialog open={this.state.openDialog === dialog.JOIN} joinGame={this.joinGame} startSpectating={this.startSpectating}/>
        </Box>
    }

    renderWaiting(){
        let inner;
        if(this.props.lobby.currentPlayers < this.props.lobby.maxPlayers){
            inner = <>
                <Typography align="center" variant="h5">Waiting for players ({this.props.lobby.currentPlayers}/{this.props.lobby.maxPlayers})</Typography>
                <Typography align="center">Invite players with this URL:</Typography>
                <Typography align="center">
                    <GameInvite lobbyId={this.props.lobby.id}/>
                </Typography>
            </>
        }else{
            inner =  <>
                <Typography align="center" variant="h5">Lobby is Full</Typography>
                {this.props.lobby.creatorId === this.props.user.id  && <Typography align="center">Start the game when you're ready</Typography>}
            </>
        }
        return <Box sx={{flexGrow: 1}}>
            {inner}
            <Typography align="center" variant={"body2"}>
                {this.props.lobby.creatorId === this.props.user.id ? <Button onClick={this.startGame}>Start</Button> : "Waiting for the lobby owner to start the game."}
            </Typography>
            <OnePlayerDialog open={this.state.openDialog === dialog.ONE_PLAYER} startGame={this.startGame} clearDialog={this.clearDialog}/>
        </Box>
    }

    renderInGame(){
        return <div id="gameWindow">
            <Board handleTileDrop={this.handleTileDrop} board={this.state.boards["main"]} name={"main"} debug={this.state.debug}/>
            <div id="boardControls">
                <div>
                    <Button disabled={!this.isTurn()} onClick={this.play}>Play</Button>
                    <Button disabled={!this.isTurn()} onClick={this.pass}>Pass</Button>
                    <Button onClick={this.toggleDebug}>Debug</Button>
                </div>
                <IconButton title="Recall" onClick={this.recallTiles}><KeyboardDoubleArrowDownIcon/></IconButton>
                <IconButton title="Shuffle" onClick={this.shuffleTiles}><ShuffleIcon/></IconButton>
                <IconButton title="Swap" onClick={this.swapTiles}><SwapVertIcon/></IconButton>
            </div>
            {this.state.boards.deck && <Board handleTileDrop={this.handleTileDrop} board={this.state.boards.deck} name={"deck"} debug={this.state.debug}/>}
            <SwapTilesDialog open={this.state.openDialog === dialog.SWAP_TILES} keepDeck={this.state.boards.deck} swapDeck={this.state.boards.swap} debug={this.state.debug} handleTileDrop={this.handleTileDrop}/>
            <InvalidPlacementDialog open={this.state.openDialog === dialog.INVALID_PLACEMENT} clearDialog={this.clearDialog} error={this.state.placementError}/>
            <BlankTileDialog open={this.state.openDialog === dialog.BLANK_TILE} clearDialog={this.clearDialog} setBlankTile={this.setBlankTile}/>
            <LobbyDeletedDialog open={this.state.openDialog === dialog.LOBBY_DELETED}/>
        </div>
    }


    async handleTileDrop(from, fromIndex, to, toIndex){
        if(!this.isTurn() && (from === "main" || to === "main"))return console.log("Ignoring invalid turn");
        console.log(this.state.boards[from].squares[fromIndex].tile?.blank, to);
        if(to === "main" && this.state.boards[from].squares[fromIndex].tile?.blank){
            this.setState({openDialog: dialog.BLANK_TILE, blankData: {from, fromIndex, to, toIndex}})
            return
        }
        console.log(`Moving ${from}#${fromIndex} -> ${to}#${toIndex}`);
        if(to !== "swap" && from !== "swap") {
            defAxios.patch(`game/${this.props.lobby.id}/boards`, {from, fromIndex, to, toIndex}).catch(()=>{
                this.setState((state)=>{
                    state.boards[from].squares[fromIndex].tile = state.boards[to].squares[toIndex].tile
                    state.boards[to].squares[toIndex].tile = null;
                    return {boards: state.boards}
                })
            })
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
       let result = await defAxios.post(`game/${this.props.lobby.id}/boards`, null, {validateStatus: ()=>true});
       if(result.data?.err){
           console.log("Error time");
           result.data.err = result.data.err[0].toUpperCase()+result.data.err.substring(1);
           this.setState({openDialog: dialog.INVALID_PLACEMENT, placementError: result.data.err});
           return;
       }
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


    async setLobbyState(state){
        let {data: lobby} = await defAxios.patch(`lobby/${this.props.lobby.id}`, {state})
        return this.setState({openDialog: null, lobby});
    }

    async startGame(){
        if(!this.state.openDialog && this.props.lobby.currentPlayers === 1){
            return this.setState({openDialog: dialog.ONE_PLAYER});
        }
        return this.setLobbyState("inGame")
    }

    async resetGame(){
        return this.setLobbyState("waiting")
    }

    async endGame(){
        await defAxios.delete(`lobby/${this.props.lobby.id}`)
    }

    async setBlankTile(letter){
        await defAxios.patch(`game/${this.props.lobby.id}/boards`, {...this.state.blankData, letter})
        const {to, toIndex, from, fromIndex} = this.state.blankData;
        this.setState((state)=>{
            state.boards[to].squares[toIndex].tile = {
                ...state.boards[from].squares[fromIndex].tile,
                letter
            }
            state.boards[from].squares[fromIndex].tile = null;
            return {boards: state.boards, openDialog: null}
        })
    }

}

export default GameWindow