import {
    Dialog,
    DialogTitle,
    Divider, IconButton, InputAdornment,
    List,
    ListItem,
    ListItemAvatar, ListItemButton,
    ListItemText, TextField,
} from "@mui/material";
import PersonIcon from '@mui/icons-material/Person';
import ArrowCircleRightIcon from '@mui/icons-material/ArrowCircleRight';
import GoogleIcon from "./GoogleIcon";
import {
    Link,
    useHref,
} from "react-router-dom";
import {useEffect, useState} from "react";


// Constants for nickname validation
const MIN_NICKNAME_LENGTH = 3;
const MAX_NICKNAME_LENGTH = 32;
const NICKNAME_REGEX = /^[A-Za-z0-9_-]*$/


export default ({onClose, open, onSignIn})=>{

    const [nickname, setNickname] = useState("");
    const [error, setError] = useState({error: false, helperText: ""});

    useEffect(()=>{
        if(nickname.length > MAX_NICKNAME_LENGTH) {
            setError({
                error: true,
                helperText: `Must be between ${MIN_NICKNAME_LENGTH} and ${MAX_NICKNAME_LENGTH} chars`
            })
        }else if(nickname.length && !NICKNAME_REGEX.test(nickname)){
            setError({
                error: true,
                helperText: `Must only contain alphanumeric and _ or -`
            })
        }else if(error.error){
            setError({error: false});
        }
    }, [nickname])

    return <Dialog onClose={onClose} open={open}>
        <DialogTitle>Login</DialogTitle>
        <List sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper' }}>
            <ListItem>
                <ListItemButton disabled onClick={()=>window.location.href="https://google.com"} >
                    <ListItemAvatar >
                        <GoogleIcon fontSize="large" />
                    </ListItemAvatar>
                    <ListItemText primary="Login with Google (TODO)"/>
                </ListItemButton>
            </ListItem>
            <Divider variant="inset" component="li" textAlign="left" sx={{fontVariantCaps: "all-small-caps"}}>or</Divider>
            <ListItem alignItems="flex-start">
                <ListItemAvatar>
                    <PersonIcon/>
                </ListItemAvatar>
                <ListItemText
                    primary="Continue as Guest"
                    secondary={
                    <TextField label="Nickname"
                        type="text"
                        size="small"
                        value={nickname}
                        {...error}
                        onChange={(e)=>setNickname(e.target.value)}
                        InputProps={{endAdornment: <InputAdornment position="end">
                        <IconButton edge="end" color="primary" disabled={nickname.length < MIN_NICKNAME_LENGTH || error.error} onClick={()=>onSignIn(nickname)}>
                            <ArrowCircleRightIcon/>
                        </IconButton>
                    </InputAdornment>}}/>}
                />
            </ListItem>
        </List>
    </Dialog>
}