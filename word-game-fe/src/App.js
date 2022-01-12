import {
    Route,
    Routes,
} from "react-router-dom";
import LobbyList from "./page/LobbyList";
import LoginManager from "./component/LoginManager";
import {useState} from "react";
import Lobby from "./page/Lobby";
import '../src/css/app.css';


function App() {
    const [user, setUser] = useState({})
    return (
        <>
            <Routes>
                <Route path="/" element={<LobbyList/>}/>
                <Route path="/lobby" element={<Lobby/>}/>
            </Routes>
            <LoginManager open={!user.name} onSignIn={(user)=>setUser(user)}/>
        </>
    );
}

export default App;
