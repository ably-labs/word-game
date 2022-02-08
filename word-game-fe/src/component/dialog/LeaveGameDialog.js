import {Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle} from "@mui/material";
import React from "react";
import {useNavigate} from "react-router-dom";

export default ({open, leaveGame, cancel, isOwner, lobby})=>{

    const getText = ()=>{
        if(isOwner)return "Leaving this lobby will end the game for everyone else.";
        if(lobby.state === "inGame")return "You won't be able to re-join";
        return "";
    }

    return  <Dialog
        open={open}
        onClose={cancel}>
        <DialogTitle id="alert-dialog-title">
            Leave Game
        </DialogTitle>
        <DialogContent>
            <DialogContentText>
                Are you sure you want to leave?<br/>
                {getText()}
            </DialogContentText>
        </DialogContent>
        <DialogActions>
            <Button onClick={cancel}>Cancel</Button>
            <Button onClick={leaveGame} autoFocus>Leave</Button>
        </DialogActions>
    </Dialog>
}