SELECT 
    m.message_id,
    u.username AS sender,
    m.content,
    m.timestamp
FROM 
    messages m
JOIN 
    users u ON m.sender_id = u.user_id
ORDER BY 
    m.timestamp;
