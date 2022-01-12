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
    const [isRegistering, setIsRegistering] = useState(false);
    const [isLoggingIn, setIsLoggingIn] = useState(false);

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
        setIsRegistering(true);
        setError({error: false, helperText: ""});
        try {
            const {data} = await defAxios.post("auth/register/start", {nickname});
            data.publicKey.challenge = Base64ToUint8Array(data.publicKey.challenge);
            data.publicKey.user.id = Base64ToUint8Array(data.publicKey.user.id);
            let credential = await navigator.credentials.create(data);
            console.log(credential);
            let {data: validationResponse, headers: {'set-cookie': cookie}} = await defAxios.post("auth/register/confirm", {
                id: credential.id,
                rawId: ArrayBufferToBase64(credential.rawId),
                response: {
                    attestationObject: ArrayBufferToBase64(credential.response.attestationObject),
                    clientDataJSON: ArrayBufferToBase64(credential.response.clientDataJSON)
                },
                type: credential.type,
            })
            onSignIn(validationResponse, cookie);
        }catch(e){
            setError({error: true, helperText: e.response?.data?.err || e.message});
        }finally {
            setIsRegistering(false);
        }
    }

    const handleLogin = async ()=>{
        setIsLoggingIn(true);
        setError({error: false, helperText: ""});
        try {
            const {data} = await defAxios.post("auth/login/start", {nickname});
            data.publicKey.challenge = Base64ToUint8Array(data.publicKey.challenge);
            data.publicKey.allowCredentials.forEach((cred) => cred.id = Base64ToUint8Array(cred.id))
            const assertion = await navigator.credentials.get(data);
            let {data: validationResponse, headers: {'set-cookie': cookie}} = await defAxios.post("auth/login/confirm", {
                id: assertion.id,
                rawId: ArrayBufferToBase64(assertion.rawId),
                response: {
                    authenticatorData: ArrayBufferToBase64(assertion.response.authenticatorData),
                    clientDataJSON: ArrayBufferToBase64(assertion.response.clientDataJSON),
                    signature: ArrayBufferToBase64(assertion.response.signature),
                    userHandle: ArrayBufferToBase64(assertion.response.userHandle)
                },
                type: assertion.type,
            });
            onSignIn(validationResponse, cookie);
        }catch(e){
            setError({error: true, helperText: e.response?.data?.err || e.message});
        }finally {
            setIsLoggingIn(false);
        }
    }

    const disabled = isRegistering || isLoggingIn || nickname.length === 0 || error.error;
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
            <Button onClick={handleLogin} disabled={disabled}>Login</Button>
            <Button onClick={handleRegister} disabled={disabled}>Register</Button>
        </DialogActions>
    </Dialog>
}


function ArrayBufferToBase64(arr){
    const base64 = btoa(String.fromCharCode(...new Uint8Array(arr))).replace(/\//g, '_').replace(/\+/g, '-');
    const ind = base64.indexOf("=");
    if(ind === -1)return base64;
    return base64.substring(0, base64.indexOf("="));
}

function Base64ToUint8Array(str){
    return Uint8Array.from(atob(str), c => c.charCodeAt(0))
}