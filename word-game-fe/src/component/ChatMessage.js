

export default ({author, message, members})=>{
    message.replace(/<@([0-9]*)>/g, (match, p1)=>{
        return members.find((m)=>m.userId === p1) || p1
    })
    if(author === "system")
        return <div className="system message">{message}</div>
    return <div className="message"><b>{author}</b>: {message}</div>
}