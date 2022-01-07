import {
    Route,
    Routes,
} from "react-router-dom";
import LobbyList from "./page/LobbyList";
import LoginManager from "./component/LoginManager";

function App() {
    return (
        <>
            <Routes>
                <Route path="/" element={<LobbyList/>}/>
            </Routes>
            <LoginManager open={true}/>
        </>
    );
}

export default App;
