import React from "react";
import Grid from "@mui/material/Grid";
import LobbyPreview from "../component/LobbyPreview";
import {useEffect, useState} from "react";
import defAxios from "../Http";
import LobbySkeleton from "../component/LobbySkeleton";
import {Box, Button, Card, CardActions, CardContent, Container, Fab, Paper, Typography} from "@mui/material";
import AddIcon from '@mui/icons-material/Add';

export default class LobbyList extends React.Component {

    state = {
        lobbies: null,
        joined: null,
    }

    constructor(props){
        super(props);
        this.onMessage = this.onMessage.bind(this);
    }

    async componentDidMount(){
        const {data: lobbies} = await defAxios.get("lobby");
        const {data: joined} = await defAxios.get("lobby/joined");
        let joinMap = {};
        for(let i = 0; i < joined.length; i++){
            joinMap[joined[i].lobbyId] = joined[i]
        }
        this.setState({lobbies, joined: joinMap});
        const channel = this.props.realtime.channels.get("lobby-list");
        console.log("Subscribing to channel");
        channel.subscribe(this.onMessage);
    }

    onMessage(message){
        switch(message.name){
            case "lobbyAdd":
                console.log("Lobby add", message);
                this.setState({lobbies: this.state.lobbies.concat(message.data)})
                break;
        }
    }

    render(){
        let inner;
        if(this.state.lobbies === null){
            inner = Array(5).fill(1).map((a,i)=><Grid item xs={2} key={`lsk${i}`}>
                <LobbySkeleton/>
            </Grid>)

        }else if(this.state.lobbies.length > 0){
            inner = this.state.lobbies.map((lobby, i)=><Grid item xs={2} key={`lobby-${i}`}>
                <LobbyPreview lobby={lobby} joined={this.state.joined[lobby.id]}/>
            </Grid>)
        }else{
            inner = <Grid item xs={6}>
                <Card sx={{maxWidth: 345}}>
                    <CardContent>
                        <Typography variant="h6">No public lobbies yet.</Typography>
                        <Typography variant="body1">Why not create one?</Typography>
                    </CardContent>
                    <CardActions>
                        <Button color="primary" href={"lobby/new"}>Create Lobby</Button>
                    </CardActions>
                </Card>
            </Grid>
        }

        return  <Container >
            <Typography variant={"h4"}>Word Game</Typography>
            <Grid container spacing={2}>{inner}</Grid>
            <Fab color="primary" aria-label="add" sx={{position: "fixed", bottom: "3vh", right: "3vw"}} href={"lobby/new"}>
                <AddIcon />
            </Fab>
        </Container>
    }
}
