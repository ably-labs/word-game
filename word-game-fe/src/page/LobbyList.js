import Grid from "@mui/material/Grid";
import LobbyPreview from "../component/LobbyPreview";
import {useEffect, useState} from "react";
import defAxios from "../Http";
import LobbySkeleton from "../component/LobbySkeleton";
import {Box, Button, Card, CardActions, CardContent, Container, Fab, Paper, Typography} from "@mui/material";
import AddIcon from '@mui/icons-material/Add';

export default ({user, realtime})=>{

    const [lobbies, setLobbies] = useState(null);

    useEffect(async ()=>{
        const {data: lobbies} = await defAxios.get("lobby");
        setLobbies(lobbies);
    }, [user]);

    let inner;
    if(lobbies === null){
        inner = Array(5).fill(1).map(()=><Grid item xs={2}>
                <LobbySkeleton/>
            </Grid>)

    }else if(lobbies.length > 0){
        inner = lobbies.map((lobby, i)=><Grid item xs={2}>
            <LobbyPreview lobby={lobby} key={`lobby-${i}`}/>
        </Grid>)
    }else{
        inner = <Grid item xs={6}>
            <Card sx={{maxWidth: 345}}>
                <CardContent>
                    <Typography variant="h6">No public lobbies yet.</Typography>
                    <Typography variant="body1">Why not create one?</Typography>
                </CardContent>
                <CardActions>
                    <Button color="primary">Create Lobby</Button>
                </CardActions>
            </Card>
        </Grid>
    }

    return  <Container >
        <Typography variant={"h4"}>Word Game</Typography>
        <Grid container spacing={2}>{inner}</Grid>
        <Fab color="primary" aria-label="add" sx={{position: "fixed", bottom: "3vh", right: "3vw"}}>
            <AddIcon />
        </Fab>
    </Container>
}