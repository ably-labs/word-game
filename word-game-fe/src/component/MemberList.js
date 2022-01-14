import '../css/memberList.css'

import React from 'react';
import defAxios from "../Http";
import {List, ListItem, ListItemText} from "@mui/material";

export default class MemberList extends React.Component {


    state = {
        members: []
    }

    async componentDidMount(){
        const {data: members} = await defAxios.get(`lobby/${this.props.lobbyId}/member`)
        this.setState({members});
        const channel = this.props.realtime.channels.get(`lobby-${this.props.lobbyId}`);
        channel.presence.enter();
        channel.presence.get((err, members)=>{
            console.log(members);
            let enrichedMembers = this.state.members.map((member)=>{
                // noinspection EqualityComparisonWithCoercionJS
                let presence = members.find((m)=>m.clientId == member.id)
                member.activity = {state: presence ? "online" : "offline"};
                return member;
            })
            this.setState({members: enrichedMembers})
        });

        channel.presence.subscribe("enter", this.updateMemberState("online"));
        channel.presence.subscribe("leave", this.updateMemberState("offline"));

        channel.subscribe(this.onMessage)
    }

    onMessage(message){
        switch(message.name){

        }
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

