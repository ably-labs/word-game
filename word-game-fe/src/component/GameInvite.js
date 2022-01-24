import React, {useRef} from "react";
import {IconButton, InputAdornment, TextField} from "@mui/material";
import {ContentCopy} from "@mui/icons-material";

export default ({lobbyId})=>{

    const ref = useRef();

    const copy = ()=>{
        ref.current.select();
        ref.current.setSelectionRange(0, 99999);
        navigator.clipboard.writeText(ref.current.value);
    }

    return <TextField value={`https://wg.unacc.eu/${lobbyId}`} inputRef={ref} InputProps={{endAdornment: <InputAdornment position="end">
            <IconButton edge="end" onClick={copy}>
                <ContentCopy/>
            </IconButton>
    </InputAdornment>}}/>
}