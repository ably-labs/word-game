import '../css/memberList.css'

import React from 'react';

export default class MemberList extends React.Component {


    state = {
        members: []
    }

    componentDidMount(){
        const channel = this.props.realtime.channels.get(`chat-${this.props.lobbyId}`);
        channel.presence.enter();
        channel.presence.get((err, members)=>{
            if(!err) this.setState({members})
            else console.error(err)
        });

        channel.presence.on("enter", (member)=>{
            this.setState({members: this.state.members.concat(member)})
        })


    }

    render(){
        return <div id="memberList">
            {this.state.members.map((member)=><div>{JSON.stringify(member)}</div>)}
        </div>
    }
}

