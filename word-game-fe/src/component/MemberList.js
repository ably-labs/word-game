import '../css/memberList.css'

import React from 'react';
import {Box, Button, List, ListItem, ListItemText} from "@mui/material";


export default ({lobby, memberPresence, members, user, leaveLobby, skipUser})=>{
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
        <Box sx={{position: "absolute", bottom: 5}}>
            <Button onClick={leaveLobby}>Leave</Button>
            {user.id === lobby.creatorId ? <Button onClick={skipUser}>Skip User</Button> : ""}
        </Box>
    </div>
}

