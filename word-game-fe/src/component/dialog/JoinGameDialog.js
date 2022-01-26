import {Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle} from "@mui/material";
import React from "react";
import {useNavigate} from "react-router-dom";

export default ({open, startSpectating, joinGame})=>{

    const navigate = useNavigate();

    const backToLobby = ()=>{
        navigate("..")
    }

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