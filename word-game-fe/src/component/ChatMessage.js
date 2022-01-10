

export default ({author, message})=>{
    if(author === "system")
        return <div className="system message">{message}</div>
    return <div className="message"><b>{author}</b>: {message}</div>
}