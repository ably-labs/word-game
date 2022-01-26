import {Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography} from "@mui/material";
import React, {useState} from "react";
import Board from "../Board";
import {useNavigate} from "react-router-dom";

export default ({open})=>{


    const navigate = useNavigate();

    const backToLobby = ()=>{
        navigate("..")
    }

    return <Dialog
        open={open}
        onClose={backToLobby}>
        <DialogTitle>Lobby Deleted</DialogTitle>
        <DialogContent>
            The lobby has been deleted by the host.
        </DialogContent>
        <DialogActions>
            <Button onClick={backToLobby} autoFocus>Okay</Button>
        </DialogActions>
    </Dialog>
}