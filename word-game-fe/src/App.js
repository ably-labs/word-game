import {
    Route,
    Routes,
} from "react-router-dom";
import LobbyList from "./page/LobbyList";
import LoginManager from "./component/LoginManager";
import {useEffect, useState} from "react";
import Lobby from "./page/Lobby";
import '../src/css/app.css';
import Ably from 'ably'
import defAxios from "./Http";

function App() {
    const [user, setUser] = useState({})
    const realtime =  new Ably.Realtime({
        authUrl: process.env.REACT_APP_BACKEND_BASE_URL+"/auth/token",
        authMethod: "GET",
        autoConnect: false,
        authHeaders: {
            // Ably does not require manually setting withCredentials,
            // so adding 'authorization' here is needed to force withCredentials
            // in order to use cookie auth
            authorization: "required",
        }
    });


    useEffect(()=>{
        defAxios.get("auth/me").then(({data: currentUser})=>{
            onSignIn(currentUser);
        }).catch(()=>null)
    }, [])

    const onSignIn = (user)=>{
        setUser(user);
        if(realtime.connection.state === "connected")return console.warn("Tried to connect when already connected");
        realtime.connection.connect()

    }

    realtime.connection.on("connected", ()=>{
        console.log("Connected!");
    })

    return (
        <>
            <Routes>
                <Route path="/" element={<LobbyList user={user} realtime={realtime}/>}/>
                <Route path="/lobby/:id" element={<Lobby/>}/>
            </Routes>
            <LoginManager open={!user.name} onSignIn={onSignIn}/>
        </>
    );
}

export default App;
