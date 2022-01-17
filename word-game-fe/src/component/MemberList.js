import '../css/memberList.css'

import React from 'react';
import defAxios from "../Http";
import {List, ListItem, ListItemText} from "@mui/material";

export default class MemberList extends React.Component {



    channel;
    state = {
        members: []
    }

    async componentDidMount(){
        const {data: members} = await defAxios.get(`lobby/${this.props.lobbyId}/member`)
        this.setState({members});
        this.channel = this.props.realtime.channels.get(`lobby-${this.props.lobbyId}`);
        this.channel.presence.enter();
        this.channel.presence.get((err, members)=>{
            console.log("Get members", err, members, this.state.members);
            let enrichedMembers = this.state.members.map((member)=>{
                // noinspection EqualityComparisonWithCoercionJS
                let presence = members.find((m)=>m.clientId == member.id)
                member.activity = {state: presence ? "online" : "offline"};
                return member;
            })
            this.setState({members: enrichedMembers})
        });

        this.channel.presence.subscribe("enter", this.updateMemberState("online"));
        this.channel.presence.subscribe("leave", this.updateMemberState("offline"));

        // TODO memberJoin/memberLeave events

    }

    componentWillUnmount() {
        this.channel.presence.unsubscribe("enter", this.updateMemberState("online"));
        this.channel.presence.unsubscribe("leave", this.updateMemberState("offline"));
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

