import {Card, CardActionArea, CardActions, CardContent, CardMedia, Skeleton, Typography} from "@mui/material";

export default ()=>{

    return  <Card sx={{ maxWidth: 345 }}>
        <CardActionArea>
            <CardMedia component={Skeleton} variant={"rectangle"} height={140}/>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    <Skeleton variant="text"/>
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    <Skeleton variant="text"/>
                    <Skeleton variant="text"/>
                    <Skeleton variant="text"/>
                </Typography>
            </CardContent>
        </CardActionArea>
        <CardActions sx={{justifyContent: "flex-end"}}>
            <Skeleton variant="rectangular" width={64} height={25} />
            <Skeleton variant="rectangular" width={64} height={25} />
        </CardActions>
    </Card>
}
