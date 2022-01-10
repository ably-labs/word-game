import '../../css/tile.css';
import AvailableTile from "./AvailableTile";

export default ({text, type, style, index, onTileDropped})=>{
    return <AvailableTile className={`tile bonus ${type}`} style={style} index={index} onTileDropped={onTileDropped}>
        <div className="bonusText">{text}</div>
    </AvailableTile>

}