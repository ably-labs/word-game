import {Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography} from "@mui/material";
import React, {useState} from "react";
import Board from "../Board";

export default ({open, clearDialog, error})=>{
        return <Dialog
        open={open}
        onClose={clearDialog}>
        <DialogTitle>Invalid Placement</DialogTitle>
        <DialogContent>
            {error}
        </DialogContent>
        <DialogActions>
            <Button onClick={clearDialog} autoFocus>Okay</Button>
        </DialogActions>
    </Dialog>
}