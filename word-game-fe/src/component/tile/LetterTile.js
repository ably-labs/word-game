
import '../../css/tile.css';

export default ({letter, score, style, draggable = false, index, source = "main"})=>{

    const onDragStart = (ev)=>{
        ev.dataTransfer.setData("index", index);
        ev.dataTransfer.setData("source", source);
    }

    return <td className="tile" style={style} draggable={draggable} data-index={index} onDragStart={onDragStart}>
        <div className="letter">{letter.toUpperCase()}</div>
        <div className="score">{score}</div>
    </td>
}