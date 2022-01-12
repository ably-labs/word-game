import {
    Button,
    Dialog, DialogActions, DialogContent,
    DialogTitle, TextField,
} from "@mui/material";
import {useEffect, useState} from "react";
import defAxios from "../Http";


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


    const handleRegister = async ()=>{
        // TODO: Error handling
        const {data} = await defAxios.post("auth/register/start", {nickname});
        data.publicKey.challenge = Uint8Array.from(atob(data.publicKey.challenge), c => c.charCodeAt(0))
        data.publicKey.user.id = Uint8Array.from(atob(data.publicKey.user.id), c => c.charCodeAt(0))
        let credential = await navigator.credentials.create(data);
        console.log(credential);
        let {data: validationResponse} = await defAxios.post("auth/register/confirm", {
            id: credential.id,
            rawId: ArrayBufferToBase64(credential.rawId),
            response: {
                attestationObject: ArrayBufferToBase64(credential.response.attestationObject),
                clientDataJSON:  ArrayBufferToBase64(credential.response.clientDataJSON)
            },
            type: credential.type,
        })
        console.log(validationResponse);
    }


    return <Dialog onClose={onClose} open={open}>
        <DialogTitle>Login or Register</DialogTitle>
        <DialogContent>
            <TextField label="Nickname"
                       type="text"
                       size="small"
                       value={nickname}
                       {...error}
                       onChange={(e)=>setNickname(e.target.value)}/>
        </DialogContent>
        <DialogActions>
            <Button>Login</Button>
            <Button onClick={handleRegister}>Register</Button>
        </DialogActions>
    </Dialog>
}


function ArrayBufferToBase64(arr){
    const base64 = btoa(String.fromCharCode(...new Uint8Array(arr))).replace(/\//g, '_').replace(/\+/g, '-');
    const ind = base64.indexOf("=");
    if(ind === -1)return base64;
    return base64.substring(0, base64.indexOf("="));
}