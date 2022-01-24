import {Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle} from "@mui/material";
import React from "react";

export default ({open, startSpectating, joinGame})=>{

    const backToLobby = ()=>{
        
    }

    // this.state.openDialog === dialog.JOIN}
    return  <Dialog
        open={open}
        onClose={backToLobby}>
        <DialogTitle id="alert-dialog-title">
            Join Game
        </DialogTitle>
        <DialogContent>
            <DialogContentText>
                Would you like to join this game?
            </DialogContentText>
        </DialogContent>
        <DialogActions>
            <Button onClick={backToLobby}>Back</Button>
            <Button onClick={startSpectating}>Spectate</Button>
            <Button onClick={joinGame} autoFocus>Join</Button>
        </DialogActions>
    </Dialog>
}