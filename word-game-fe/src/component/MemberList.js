import '../css/memberList.css'

import React from 'react';
import defAxios from "../Http";
import {List, ListItem, ListItemText} from "@mui/material";

export default class MemberList extends React.Component {



    channel;
    interval;
    state = {
        members: []
    }


    constructor(props){
        super(props);
        this.onMessage = this.onMessage.bind(this);
        this.updatePresence = this.updatePresence.bind(this);
    }

    componentDidMount(){
        // The user ID is not immediately available
        if(!this.props.user?.id)return;
        this.setupPresence();
    }


    componentDidUpdate(lastProps, lastState, snapshot){
        if(lastProps.user === this.props.user || this.channel)return;
        this.setupPresence();
    }

    async setupPresence(){
        const {data: members} = await defAxios.get(`lobby/${this.props.lobbyId}/member`)
        this.setState({members});
        this.channel = this.props.realtime.channels.get(`lobby-${this.props.lobbyId}`);

        this.channel.presence.enterClient(""+this.props.user.id);

        this.updatePresence();

        // Update presence every 5 seconds to keep it in sync
        if(!this.interval)
            this.interval = setInterval(this.updatePresence, 5000)

        this.channel.presence.subscribe("enter", this.updateMemberState("online"));
        this.channel.presence.subscribe("leave", this.updateMemberState("offline"));

        this.channel.subscribe(this.onMessage);
    }

    onMessage(message) {
        console.log(message);
        switch (message.name) {
            case "memberAdd":
                this.setState((state)=>{
                    state.members.push(message.data)
                    return state;
                });
                break;
        }
    }

    updatePresence(){
        this.channel.presence.get((err, members)=>{
            let enrichedMembers = this.state.members.map((member)=>{
                // noinspection EqualityComparisonWithCoercionJS
                let presence = members.find((m)=>m.clientId == member.id)
                member.activity = {state: presence ? "online" : "offline"};
                return member;
            })
            this.setState({members: enrichedMembers})
        });
    }

    componentWillUnmount() {
        this.channel.presence.unsubscribe("enter", this.updateMemberState("online"));
        this.channel.presence.unsubscribe("leave", this.updateMemberState("offline"));
        clearInterval(this.interval);
    }

    updateMemberState(state){
        return (member)=>{
            this.setState(({members})=>{
                // noinspection EqualityComparisonWithCoercionJS
                const ind = members.findIndex((m)=> m.id == member.clientId);
                if(ind === -1)return console.log("Couldn't find", member);
                members[ind].activity = {state};
                return {members};
            })
        }
    }

    render(){
        return <div id="memberList">
            <List>
                {this.state.members.map((member)=> <ListItem key={`presence-${member.id}`}>
                    <ListItemText
                        primary={member.user.name+(member.type === "spectator" ? " (Spectator)" : "")}
                        secondary={<span className={`activity ${member.activity?.state || "offline"}`}>{member.activity?.state || "Offline"}</span>}
                    />
                </ListItem>)}
            </List>
        </div>
    }
}

