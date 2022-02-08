import {
    Box,
    Button,
    Checkbox,
    Container,
    FormControlLabel,
    FormGroup,
    InputLabel, MenuItem, Select,
    TextField,
    Typography
} from "@mui/material";
import {useEffect, useState} from "react";
import defAxios from "../Http";
import {useNavigate} from "react-router-dom";

export default ()=>{


    const [name, setName] = useState("");
    const [isPrivate, setPrivate] = useState(false)
    const [error, setError] = useState(false);
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();
    const [gameTypes, setGameTypes] = useState([]);
    const [gameTypeIndex, setGameTypeIndex] = useState(1);
    const [gameTypeData, setGameTypeData] = useState({})
    const [customGameType, setCustomGameType] = useState(false);


    useEffect(async ()=>{
        const {data: gameTypes} = await defAxios.get("/lobby/types");
        setGameTypes(gameTypes);
        setGameTypeData(gameTypes[0]);
    }, [])

    const createLobby = async ()=>{
        setLoading(true);
        let gameType = gameTypeIndex;
        if(customGameType){
            let {data: gameTypeObj} = await defAxios.post("/lobby/types", {
                ...gameTypeData,
                id: null,
                name: `custom-${new Date().getTime()}`
            })
            gameType = gameTypeObj.id;
        }
        let result = await defAxios.post("/lobby", {
            name, private: isPrivate, gameType
        });
        navigate(`/lobby/${result.data.id}`);
        console.log(result);
    }

    const changeGameType = async (e)=>{
        let index = e.target.value;
        if(index === -1){
            setCustomGameType(true);
        }else {
            setCustomGameType(false);
            let gameType = gameTypes.find((gt) => gt.id === index);
            setGameTypeIndex(gameType.id);
            setGameTypeData(gameType);
        }
    }

    return <Container>
        <Typography variant={"h4"}>New Lobby</Typography>
        <FormGroup>
            <TextField required label="Lobby Name" variant="standard" value={name} onChange={(e)=>setName(e.target.value)}/>
        </FormGroup>
        <FormGroup>
            <FormControlLabel control={<Checkbox checked={isPrivate} onChange={(e)=>setPrivate(e.target.checked)} />} label="Private Lobby" />
        </FormGroup>
        {/*<FormGroup>*/}
        {/*    <InputLabel>Game Type</InputLabel>*/}
        {/*    <Select value={customGameType ? -1 : gameTypeIndex} label="Game Type" onChange={changeGameType}>*/}
        {/*        {gameTypes.map((gt)=> <MenuItem value={gt.id}>{gt.name}</MenuItem>)}*/}
        {/*        <MenuItem value={-1}>Custom</MenuItem>*/}
        {/*    </Select>*/}
        {/*    <hr/>*/}
        {/*    <Typography variant="h5">Game Properties</Typography>*/}
        {/*    <FormControlLabel disabled={!customGameType} control={<Checkbox value={gameTypeData.enableBlankTiles} />} label="Enable Blank Tiles" />*/}
        {/*    <FormControlLabel disabled={!customGameType} control={<Checkbox value={gameTypeData.startAnywhere}/>} label="Start Anywhere" />*/}
        {/*</FormGroup>*/}
        {/*<FormGroup>*/}
        {/*    <Box>*/}
        {/*        <Typography>Board Size</Typography>*/}
        {/*        <TextField disabled={!customGameType} margin="normal" label="Width" type="number" value={gameTypeData.boardWidth}/>*/}
        {/*        <TextField disabled={!customGameType} margin="normal" label="Height" type="number" value={gameTypeData.boardHeight}/>*/}
        {/*    </Box>*/}
        {/*    <TextField disabled={!customGameType} margin="normal" label="Bag Size" type="number" value={gameTypeData.tileBagSize}/>*/}
        {/*    <TextField disabled={!customGameType} margin="normal" label="Deck Size" type="number" value={gameTypeData.playerDeckSize}/>*/}
        {/*    <TextField disabled={!customGameType} margin="normal" label="Tile Count" type="number" value={gameTypeData.playerTileCount}/>*/}
        {/*    <Select disabled={!customGameType} margin="normal" value={gameTypeData.bonusTilePattern || "regular"} label="Game Type">*/}
        {/*        <MenuItem value="none">None</MenuItem>*/}
        {/*        <MenuItem value="regular">Regular</MenuItem>*/}
        {/*        <MenuItem value="stripe">Stripe</MenuItem>*/}
        {/*        <MenuItem value="border">Border</MenuItem>*/}
        {/*    </Select>*/}
        {/*</FormGroup>*/}
        <Button onClick={createLobby} disabled={loading}>Create</Button>
    </Container>
}