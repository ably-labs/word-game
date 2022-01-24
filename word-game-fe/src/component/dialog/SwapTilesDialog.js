import {Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography} from "@mui/material";
import React, {useState} from "react";
import Board from "../Board";

export default ({open, clearDialog, swapTiles, keepDeck, swapDeck, debug, handleTileDrop})=>{
        return <Dialog
        open={open}
        onClose={clearDialog}>
        <DialogTitle>Swap Tiles</DialogTitle>
        <DialogContent>
            <Typography>Keep</Typography>
            <Board board={keepDeck} name="deck" debug={debug} handleTileDrop={handleTileDrop}/>
            <Typography>Swap</Typography>
            <Board board={swapDeck} name="swap" debug={debug} handleTileDrop={handleTileDrop}/>
        </DialogContent>
        <DialogActions>
            <Button onClick={clearDialog}>Cancel</Button>
            <Button onClick={swapTiles} autoFocus>Confirm</Button>
        </DialogActions>
    </Dialog>
}