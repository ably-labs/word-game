import GameWindow from "../component/GameWindow";
import Chat from "../component/Chat";

import "../css/lobby.css";
import MemberList from "../component/MemberList";
import {useParams} from "react-router-dom";
import {useEffect, useRef, useState} from "react";
import defAxios from "../Http";
import JoinGameDialog from "../component/dialog/JoinGameDialog";


export default ({realtime, user})=>{
    const params = useParams();
    const lobbyId = params.id;


    const [lobby, setLobby] = useState({});
    const [joining, setJoining] = useState(true);
    const [members, setMembers] = useState([]);
    const [memberPresence, setMemberPresence] = useState({});
    const [channel, setChannel] = useState(null);
    const memberPresenceRef = useRef();

    memberPresenceRef.current = memberPresence;

    const onMessage = (message)=>{
        switch(message.name){
            case "memberAdd":
                setMembers(members.concat(message.data));
                setLobby({
                    ...lobby,
                    currentPlayers: lobby.currentPlayers+1
                })
                break;
            case "lobbyUpdate":
                setLobby(message.data);
                break;
            case "scoreUpdate":
                setMembers(members.map((m)=>m.id === message.data.id ? {...m, score: message.data.score} : m));
                break;
        }
    }


    const updatePresence = ()=>{
        channel.presence.get((err, presenceMembers)=>{
            let presence = {};
            for(let i = 0; i < presenceMembers.length; i++){
                const pm = presenceMembers[i];
                pm.state = pm.action === "present" ? "online" : "offline"
                presence[pm.clientId] = pm;
            }
            setMemberPresence(presence);
        });
    }

    const updateMemberState = (state)=>{
        return (member)=>{
            setMemberPresence({
                ...memberPresenceRef.current,
                [member.clientId]: {...member, state}
            })
        }
    }


    const fetchLobby = async ()=>{
        await defAxios.get(`lobby/${lobbyId}`).then(({data: lobby})=>{
            setLobby(lobby);
            setJoining(false);
        }).catch(()=>setJoining(true));
        const {data: members} = await defAxios.get(`lobby/${lobbyId}/member`);
        setMembers(members);
    }

    useEffect(()=>{
        if(!user)return;
        if(!channel){
            setChannel(realtime.channels.get(`lobby-${lobbyId}`));
        }
        fetchLobby();
    }, [user, joining])


    useEffect(()=>{
        if(!channel)return;
        updatePresence();
    }, [channel, members]);

    useEffect(()=>{
        if(!channel)return;
        channel.subscribe(onMessage);
        return function cleanup(){
            channel.unsubscribe(onMessage);
        }
    }, [channel, lobby, members])




    useEffect(()=>{
        if(!channel || !user?.id)return;
        channel.presence.enterClient(""+user.id);
        updatePresence();
        const enterEvent = updateMemberState("online");
        const leaveEvent = updateMemberState("offline")
        channel.presence.subscribe("enter", enterEvent);
        channel.presence.subscribe("leave", leaveEvent);
        return function cleanup(){
            channel.presence.unsubscribe(enterEvent);
            channel.presence.unsubscribe(leaveEvent);
        }
    }, [channel, user])

    const joinGame = ()=>{
        defAxios.put(`lobby/${lobbyId}/member`, {type: "player"}).then(()=>setJoining(false));
    }

    const startSpectating = ()=>{
        defAxios.put(`lobby/${lobbyId}/member`, {type: "spectator"}).then(()=>setJoining(false));
    }

    if(joining)
        return <JoinGameDialog open={true} joinGame={joinGame} startSpectating={startSpectating}/>


    return <div id="lobby">
        <GameWindow realtime={realtime} lobby={lobby} members={members} user={user} channel={channel}/>
        <Chat realtime={realtime} lobbyId={lobbyId} members={members} channel={channel}/>
        <MemberList realtime={realtime} lobby={lobby} members={members} memberPresence={memberPresence} user={user} channel={channel}/>
    </div>
}