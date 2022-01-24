import React from "react";
import LetterTile from "./tile/LetterTile";
import BonusTile from "./tile/BonusTile";
import EmptyTile from "./tile/AvailableTile";

export default ({board, name, debug, handleTileDrop})=>{
    if(!board)return <div>Board {name} not found</div>
    const rows = [];
    for(let y = 0; y < board.height; y++){
        let cols = [];
        for(let x = 0; x < board.width; x++){
            const i = (y*board.width)+x;
            const square = board.squares[i];
            cols.push(<Square square={square} i={i} source={name} debug={debug} handleTileDrop={handleTileDrop}/>)
        }
        rows.push(<tr key={`board-${name}-${y}`}>{cols}</tr>)
    }
    return <table className={`board ${name}`}>
        <tbody>
        {rows}
        </tbody>
    </table>
}


function Square({square, i, source, handleTileDrop, debug}){
    const key = `${source}${i}`
    if(square?.tile)
        return <LetterTile {...square.tile} index={i} source={source} key={key} debug={debug}/>

    if(square?.bonus)
        return <BonusTile {...square.bonus} index={i} onTileDropped={handleTileDrop} source={source} key={key} debug={debug}/>

    return <EmptyTile index={i} onTileDropped={handleTileDrop} source={source} key={key} debug={debug}/>
}