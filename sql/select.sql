SELECT 
    m.id,
    u.username AS sender,
    m.content,
    m.timestamp
FROM 
    messages m
JOIN 
    users u ON m.sender_id = u.id
ORDER BY 
    m.timestamp;
