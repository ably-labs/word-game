import React from 'react';
import '../css/chat.css'
import {IconButton, InputAdornment, TextField} from "@mui/material";
import ChatMessage from "./ChatMessage";
import SendIcon from '@mui/icons-material/Send';
import defAxios from "../Http";

export default class Chat extends React.Component {

    state = {
        draft: "",
        messages: [],
    }


    anchor = React.createRef();

    constructor(props){
        super(props);
        this.onMessage = this.onMessage.bind(this);
        this.sendChat = this.sendChat.bind(this);
    }

    async componentDidMount(){
        if(!this.props.channel)return;
        let {data: messages} = await defAxios.get(`chat/${this.props.lobbyId}`)
        this.setState({messages}, ()=>this.anchor.current.scrollIntoView());
        console.log("Subscribing to channel");
        this.props.channel.subscribe(this.onMessage);
    }

    onMessage(message){
        switch(message.name){
            case "message":
                this.setState({messages: this.state.messages.concat(message.data)}, ()=>this.anchor.current.scrollIntoView())
                break;
        }
    }


    async sendChat(){
        if(this.state.draft.length === 0)return;
        await defAxios.post(`chat/${this.props.lobbyId}`, {message: this.state.draft});
        this.setState({draft: ""});
    }

    render(){
        return <div id="chat">
            <div id="chatHistory">
                {this.state.messages.map((m,i)=><ChatMessage {...m} key={`message-${i}`} members={this.props.members}/>)}
                <div ref={this.anchor}/>
            </div>
            <div id="chatControls">
                <TextField placeholder="Message..."
                           value={this.state.draft}
                           onChange={(e)=>this.setState({draft: e.target.value})}
                           sx={{width: "100%"}}
                           onKeyDown={(e)=>e.key === "Enter" && this.sendChat()}
                           InputProps={{endAdornment: <InputAdornment position="end">
                        <IconButton edge="end" color="primary" onClick={this.sendChat}>
                            <SendIcon/>
                        </IconButton>
                    </InputAdornment>}}/>
            </div>

        </div>
    }
}