import '../css/chat.css'
import {IconButton, InputAdornment, TextField} from "@mui/material";
import ChatMessage from "./ChatMessage";
import SendIcon from '@mui/icons-material/Send';
export default ({messages})=>{
    return <div id="chat">
        <div id="chatHistory">
            {messages.map((m)=><ChatMessage {...m}/>)}
        </div>
        <div id="chatControls">
            <TextField placeholder="Message..." sx={{width: "100%"}}  InputProps={{endAdornment: <InputAdornment position="end">
                    <IconButton edge="end" color="primary">
                        <SendIcon/>
                    </IconButton>
                </InputAdornment>}}/>
        </div>

    </div>
}