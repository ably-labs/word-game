import '../css/memberList.css'

import React from 'react';
import {List, ListItem, ListItemText} from "@mui/material";


export default ({lobby, memberPresence, members})=>{
    return <div id="memberList">
        <List>
            {members.filter((m)=>m.type !== "spectator" || memberPresence[m.id] && memberPresence[m.id].state !== "offline" ).map((member)=>{
                const isSpectator = member.type === "spectator";
                const isTurn = lobby.playerTurnId === member.user.id;
                const presence = memberPresence[member.id];
                let name = member.user.name;
                if(isSpectator)name += " (Spectator)"
                else if(isTurn)name += " (Current Turn)"
                return <ListItem key={`presence-${member.id}`}>
                    <ListItemText
                        primary={name}
                        secondary={<>
                            {!isSpectator && <div>{member.score.toLocaleString()} Points</div>}
                            <div className={`activity ${presence?.state || "offline"}`}>{presence?.state || "Offline"}</div>
                        </>}
                    />
                </ListItem>})}
        </List>
    </div>
}

