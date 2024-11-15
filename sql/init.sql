-- Tabelle für Benutzer
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabelle für Nachrichten
CREATE TABLE messages (
    message_id SERIAL PRIMARY KEY,
    sender_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    receiver_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT FALSE
);

-- Beispielbenutzer hinzufügen
INSERT INTO users (username) VALUES ('alice'), ('bob');

-- Beispielnachrichten hinzufügen
INSERT INTO messages (sender_id, receiver_id, content)
VALUES
    (1, 2, 'Hallo Bob! Wie geht’s?'),
    (2, 1, 'Hi Alice! Alles gut, danke. Wie geht es dir?');
