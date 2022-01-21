import '../../css/tile.css';
import AvailableTile from "./AvailableTile";


let tuples = ["","", "double", "triple", "quadruple", "quintuple", "sextuple", "septuple", "octuple", "nonuple", "decuple"]

export default ({text, type, wordMultiplier, letterMultiplier, style, index, onTileDropped, debug})=>{
    // Support for old board
    if(text && type){
        return <AvailableTile className={`tile bonus ${type}`} style={style} index={index} onTileDropped={onTileDropped} debug={debug}>
            <div className="bonusText">{text}</div>
        </AvailableTile>
    }
    let numMultiplier = wordMultiplier || letterMultiplier
    type = wordMultiplier ? "word" : "letter";
    let multiplier =  tuples[numMultiplier] || "x"+numMultiplier

    return <AvailableTile className={`tile bonus ${multiplier}-${type}`} style={style} index={index} onTileDropped={onTileDropped} debug={debug}>
        <div className="bonusText">{multiplier} {type}</div>
    </AvailableTile>

}