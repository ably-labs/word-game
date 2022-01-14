import {Button, Checkbox, Container, FormControlLabel, FormGroup, TextField, Typography} from "@mui/material";
import {useState} from "react";
import defAxios from "../Http";
import {useNavigate} from "react-router-dom";

export default ()=>{


    const [name, setName] = useState("");
    const [isPrivate, setPrivate] = useState(false)
    const [error, setError] = useState(false);
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    const createLobby = async ()=>{
        setLoading(true);
        let result = await defAxios.post("/lobby", {
            name, private: isPrivate
        });
        navigate(`/lobby/${result.data.id}`);
        console.log(result);
    }

    return <Container>
        <Typography variant={"h4"}>New Lobby</Typography>
        <FormGroup>
            <TextField required label="Lobby Name" variant="standard" value={name} onChange={(e)=>setName(e.target.value)}/>
        </FormGroup>
        <FormGroup>
            <FormControlLabel control={<Checkbox checked={isPrivate} onChange={(e)=>setPrivate(e.target.checked)} />} label="Private Lobby" />
        </FormGroup>
        <Button onClick={createLobby} disabled={loading}>Create</Button>
    </Container>
}