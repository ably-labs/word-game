import {Button, Card, CardActionArea, CardActions, CardContent, CardMedia, Typography} from "@mui/material";


const stateMap = {
    waiting: "Waiting for Players...",
    inGame: "In Game",
    roundOver: "Post-game"
}

export default ({lobby})=>{

    const joinLobby = async ()=>{

    }

    const spectateLobby = async ()=>{

    }

    return  <Card sx={{ maxWidth: 345 }}>
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
            <Button size="small" >
                Spectate
            </Button>
            <Button size="small" color="primary" disabled={!lobby.joinable}>
                Join
            </Button>
        </CardActions>
    </Card>
}
