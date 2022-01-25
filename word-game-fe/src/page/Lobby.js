import GameWindow from "../component/GameWindow";
import Chat from "../component/Chat";

import "../css/lobby.css";
import MemberList from "../component/MemberList";
import {useParams} from "react-router-dom";


export default ({realtime, user})=>{
    const params = useParams();
    const lobbyId = params.id;
    return <div id="lobby">
        <GameWindow realtime={realtime} lobbyId={lobbyId} user={user}/>
        <Chat realtime={realtime} lobbyId={lobbyId}/>
        <MemberList realtime={realtime} lobbyId={lobbyId} user={user}/>
    </div>
}