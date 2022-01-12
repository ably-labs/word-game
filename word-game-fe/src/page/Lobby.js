import GameWindow from "../component/GameWindow";
import Chat from "../component/Chat";

import "../css/lobby.css";
import MemberList from "../component/MemberList";
import {useParams} from "react-router-dom";

const messages = [
    {author: "peter", message: "Hello"},
    {author: "system", message: "This is a system message."},
    {author: "peter", message: "This message is really long and will probably wrap"},
]

export default ()=>{
    const params = useParams();
    console.log(params);
    return <div id="lobby">
        <GameWindow/>
        <Chat messages={messages}/>
        <MemberList/>
    </div>
}