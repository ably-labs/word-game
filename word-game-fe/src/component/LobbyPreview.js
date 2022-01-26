import {Button, Card, CardActionArea, CardActions, CardContent, CardMedia, Typography} from "@mui/material";
import defAxios from "../Http";
import { useNavigate } from "react-router-dom";


const stateMap = {
    waiting: "Waiting for Players...",
    inGame: "In Game",
    roundOver: "Post-game"
}

export default ({lobby, joined})=>{

    const navigate = useNavigate();

    const joinLobby = (type)=>{
        return async ()=>{
            await defAxios.put(`lobby/${lobby.id}/member`, {type})
            navigate(`/lobby/${lobby.id}`);
        }

    }

    const makeButtons = ()=>{
        if(joined){
            return <Button size="small" color="primary" onClick={()=>navigate(`/lobby/${lobby.id}`)}>
                Resume
            </Button>

        }
        return <><Button size="small" onClick={joinLobby("spectator")}>
            Spectate
        </Button>
        <Button size="small" color="primary" disabled={!lobby.joinable || lobby.state === "inGame"} onClick={joinLobby("player")}>
            Join
        </Button></>
    }

    return  <Card sx={{ maxWidth: 345}}>
        <CardActionArea>
            <CardMedia
                component="img"
                height="140"
                image={`${process.env.REACT_APP_BACKEND_BASE_URL}/lobby/${lobby.id}/thumbnail`}
                alt={lobby.name}
                sx={{objectFit: "contain", background: "grey"}}
            />
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    {lobby.name}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    {lobby.currentPlayers}/{lobby.maxPlayers} players.<br/>
                    Game Type: {lobby.gameType.name}<br/>
                    <b>{stateMap[lobby.state]}</b>
                </Typography>
            </CardContent>
        </CardActionArea>
        <CardActions>
            {makeButtons()}
        </CardActions>
    </Card>
}
