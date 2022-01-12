import Grid from "@mui/material/Grid";
import LobbyPreview from "../component/LobbyPreview";
import {useEffect, useState} from "react";
import defAxios from "../Http";


let games = (new Array(10)).fill(1).map((_, i)=>({
    playerCount: Math.round(Math.random()*8),
    maxPlayers: 8,
    name: `Lobby #${i}`,
    gameType: "Standard",
    status: Math.random() < 0.6 ? "Waiting for players..." : "In Game",
    joinable: Math.random() < 0.3,
    thumbnail: "https://placekitten.com/256/256"
}))

export default ({user, realtime})=>{

    const [lobbies, setLobbies] = useState([]);

    useEffect(async ()=>{
        const {data: lobbies} = await defAxios.get("lobby");
        setLobbies(lobbies)
    }, [user]);

    return <Grid container spacing={2}>
        {lobbies.map((lobby)=><Grid item xs={2}>
            <LobbyPreview lobby={lobby}/>
        </Grid>)}
    </Grid>
}