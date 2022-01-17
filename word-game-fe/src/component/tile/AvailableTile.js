import React, {useState} from "react";
import "../../css/tile.css";

export default ({index, children, style, className = "tile empty", onTileDropped, source = "main"})=>{

    const [draggedOver, setDraggedOver] = useState(false);

    const dragEnter = (ev)=>{
        console.log("Drag");
        setDraggedOver(true);
        ev.preventDefault();
    }

    const dragLeave = (ev)=>{
        setDraggedOver(false);
        ev.preventDefault();
    }

    // This is required to tell the browser this element is droppable
    const dragOver = (ev)=>{
        ev.preventDefault();
    }

    const onDrop = (ev)=>{
        if(!ev.dataTransfer.getData("source"))return;
        onTileDropped(ev.dataTransfer.getData("source"), parseInt(ev.dataTransfer.getData("index")), source, index)
        ev.preventDefault();
    }

    return <td className={className+(draggedOver ? " drag" : "")}
               onDragOver={dragOver}
               onDragEnter={dragEnter}
               onDragLeave={dragLeave}
               onDrop={onDrop}
               data-index={index}
               style={style}>
        {index}
        {/*{children}*/}
    </td>;
}