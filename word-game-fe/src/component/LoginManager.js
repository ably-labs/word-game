import {
    Dialog,
    DialogTitle,
    Divider, IconButton, InputAdornment,
    List,
    ListItem,
    ListItemAvatar, ListItemButton,
    ListItemText, TextField,
} from "@mui/material";
import PersonIcon from '@mui/icons-material/Person';
import ArrowCircleRightIcon from '@mui/icons-material/ArrowCircleRight';
import GoogleIcon from "./GoogleIcon";
import {
    Link,
    useHref,
} from "react-router-dom";


export default ({onClose, open})=>{
    return <Dialog onClose={onClose} open={open}>
        <DialogTitle>Login</DialogTitle>
        <List sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper' }}>
            <ListItem>
                <ListItemButton disabled onClick={()=>window.location.href="https://google.com"} >
                    <ListItemAvatar >
                        <GoogleIcon fontSize="large" />
                    </ListItemAvatar>
                    <ListItemText primary="Login with Google (TODO)"/>
                </ListItemButton>
            </ListItem>
            <Divider variant="inset" component="li" textAlign="left" sx={{fontVariantCaps: "all-small-caps"}}>or</Divider>
            <ListItem alignItems="flex-start">
                <ListItemAvatar>
                    <PersonIcon/>
                </ListItemAvatar>
                <ListItemText
                    primary="Continue as Guest"
                    secondary={<TextField label="Nickname" type="text" size="small" InputProps={{endAdornment: <InputAdornment position="end">
                        <IconButton edge="end" color="primary">
                            <ArrowCircleRightIcon color="primary"/>
                        </IconButton>
                    </InputAdornment>}}/>}
                />
            </ListItem>
        </List>
    </Dialog>
}