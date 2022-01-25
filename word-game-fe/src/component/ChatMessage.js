


export default ({author, message, members})=>{
    // noinspection EqualityComparisonWithCoercionJS
    message = message.replace(/<@([0-9]*)>/g, (match, p1)=>members.find((m)=>m.id == p1)?.user.name || p1)
    if(author === "system")
        return <div className="system message">{message}</div>
    return <div className="message"><b>{author}</b>: {message}</div>
}