import React from "react";
import Grid from "@mui/material/Grid";
import LobbyPreview from "../component/LobbyPreview";
import defAxios from "../Http";
import LobbySkeleton from "../component/LobbySkeleton";
import {Button, Card, CardActions, CardContent, Container, Fab, Paper, Typography} from "@mui/material";
import AddIcon from '@mui/icons-material/Add';

export default class LobbyList extends React.Component {

    state = {
        lobbies: null,
        joined: null,
    }

    channel;

    constructor(props){
        super(props);
        this.onMessage = this.onMessage.bind(this);
    }

    async componentDidMount(){
        const {data: lobbies} = await defAxios.get("lobby");
        this.setState({lobbies});
        this.channel = this.props.realtime.channels.get("lobby-list");
        console.log("Subscribing to channel");
        this.channel.subscribe(this.onMessage);
        if(this.props.user)
            return this.getJoinedLobbies();
    }


    async getJoinedLobbies(){
        const {data: joined} = await defAxios.get("lobby/joined");
        let joinMap = {};
        for(let i = 0; i < joined.length; i++){
            joinMap[joined[i].lobbyId] = joined[i]
        }
        this.setState({joined: joinMap});
    }

    componentWillUnmount() {
        this.channel.unsubscribe(this.onMessage)
    }


    componentDidUpdate(prevProps, prevState, ss){
        if(this.props.user !== prevProps.user){
            return this.getJoinedLobbies();
        }
    }

    onMessage(message){
        switch(message.name){
            case "lobbyAdd":
                console.log("Lobby add", message);
                this.setState({lobbies: this.state.lobbies.concat(message.data)})
                break;
            case "lobbyRemove":
                this.setState({lobbies: this.state.lobbies.filter((l)=>l.id !== message.data.id)})
                break;
            case "lobbyUpdate":
                this.setState((state)=>{
                    let lind = state.lobbies.findIndex((l)=>l.id === message.data.id);
                    if(lind > -1)
                        state.lobbies[lind] = message.data;
                    else
                        state.lobbies.push(message.data);
                    return {lobbies: state.lobbies};
                })
                break;
        }
    }


    mapLobby(joined){
        return (lobby, i)=><Grid item xs={2} key={`lobby-${i}`}>
            <LobbyPreview lobby={lobby} joined={joined} user={this.props.user}/>
        </Grid>
    }

    renderInner(){
        if(this.state.lobbies === null)
            return <Grid container spacing={2}>{Array(5).fill(1).map((a,i)=><Grid item xs={2} key={`lsk${i}`}>
                <LobbySkeleton/>
            </Grid>)}</Grid>

        if(this.state.lobbies.length === 0)
            return <Grid container spacing={2}>
                <Grid item xs={6}>
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
            </Grid>
        if(!this.state.joined)
            return <Grid container spacing={2}>{this.state.lobbies.map(this.mapLobby(false))}</Grid>

        return <>
            <Typography variant="h5">My Games</Typography>
            <Grid container spacing={2}>
                {this.state.lobbies.filter((l)=>this.state.joined[l.id]).map(this.mapLobby(true))}
            </Grid>
            <Typography variant="h5">New Games</Typography>
            <Grid container spacing={2}>
                {this.state.lobbies.filter((l)=>!this.state.joined[l.id]).map(this.mapLobby(false))}
            </Grid>
        </>
    }

    render(){
        return <Container >
            <Typography variant={"h4"}>Word Game</Typography>
            {this.renderInner()}
            <Fab color="primary" aria-label="add" sx={{position: "fixed", bottom: "3vh", right: "3vw"}} href={"lobby/new"}>
                <AddIcon />
            </Fab>
        </Container>
    }
}
