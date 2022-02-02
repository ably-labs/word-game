import {Button, Card, CardActionArea, CardActions, CardContent, CardMedia, Typography} from "@mui/material";
import defAxios from "../Http";
import { useNavigate } from "react-router-dom";
import LockIcon from '@mui/icons-material/Lock';
import NewReleasesIcon from '@mui/icons-material/NewReleases';
const stateMap = {
    waiting: "Waiting for Players...",
    inGame: "In Game",
    roundOver: "Post-game"
}

export default ({lobby, joined, user})=>{

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
                <Typography gutterBottom variant="h6" component="div">
                    {lobby.private ? <LockIcon fontSize="small"/> : ""}{lobby.name}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    {lobby.currentPlayers}/{lobby.maxPlayers} players.<br/>
                    Game Type: {lobby.gameType.name}<br/>
                    <b>{stateMap[lobby.state]}</b>
                    {lobby.state === "inGame" && lobby.playerTurnId === user?.id ? <Typography color="error"  style={{
                        display: 'flex',
                        alignItems: 'center',
                        flexWrap: 'wrap',
                    }}><NewReleasesIcon/>&nbsp;Your Turn!</Typography> : ""}
                </Typography>
            </CardContent>
        </CardActionArea>
        <CardActions>
            {makeButtons()}
        </CardActions>
    </Card>
}
