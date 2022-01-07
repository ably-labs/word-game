import {Button, Card, CardActionArea, CardActions, CardContent, CardMedia, Typography} from "@mui/material";

export default ({lobby})=>{
    return  <Card sx={{ maxWidth: 345 }}>
        <CardActionArea>
            <CardMedia
                component="img"
                height="140"
                image={lobby.thumbnail}
                alt={lobby.name}
                sx={{objectFit: "contain", background: "grey"}}
            />
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    {lobby.name}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    {lobby.playerCount}/{lobby.maxPlayers} players.<br/>
                    Game Type: {lobby.gameType}<br/>
                    <b>{lobby.status}</b>
                </Typography>
            </CardContent>
        </CardActionArea>
        <CardActions>
            <Button size="small" color="primary">
                {lobby.joinable ? "Join" : "Spectate"}
            </Button>
        </CardActions>
    </Card>
}
