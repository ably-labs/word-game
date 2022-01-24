import {Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle} from "@mui/material";
import React from "react";

export default ({open, clearDialog, startGame})=>{
    // this.state.openDialog === dialog.ONE_PLAYER
    return <Dialog
        open={open}
        onClose={clearDialog}>
        <DialogTitle id="alert-dialog-title">
            Really Start?
        </DialogTitle>
        <DialogContent>
            <DialogContentText>
                You are the only player currently in the lobby.
            </DialogContentText>
        </DialogContent>
        <DialogActions>
            <Button onClick={clearDialog}>Cancel</Button>
            <Button onClick={startGame} autoFocus>
                Start
            </Button>
        </DialogActions>
    </Dialog>
}