import {Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, Typography} from "@mui/material";
import React from "react";


const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ".split("");

export default ({open, setBlankTile, clearDialog})=>{

    const selectTile = (letter)=>{
        return ()=>{
            setBlankTile(letter);
        }
    }

    return <Dialog
        open={open}
        onClose={clearDialog}>
        <DialogTitle id="alert-dialog-title">
            Blank Tile
        </DialogTitle>
        <DialogContent>
            <DialogContentText>
                <Typography>Choose a letter.</Typography>
                {letters.map((l)=><div className="tile selector" onClick={selectTile(l)}>
                    <div className="letter">{l}</div>
                </div>)}
            </DialogContentText>
        </DialogContent>
        <DialogActions>
            <Button onClick={clearDialog}>Cancel</Button>
        </DialogActions>
    </Dialog>
}