SELECT 
    m.message_id,
    u1.username AS sender,
    u2.username AS receiver,
    m.content,
    m.timestamp,
    m.is_read
FROM 
    messages m
JOIN 
    users u1 ON m.sender_id = u1.user_id
JOIN 
    users u2 ON m.receiver_id = u2.user_id
WHERE 
    (m.sender_id = 1 AND m.receiver_id = 2)
    OR (m.sender_id = 2 AND m.receiver_id = 1)
ORDER BY 
    m.timestamp;
