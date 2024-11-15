-- Tabellen löschen, falls sie existieren
DROP TABLE IF EXISTS messages CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Tabelle für Benutzer
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabelle für Nachrichten
CREATE TABLE messages (
    message_id SERIAL PRIMARY KEY,
    sender_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Beispielbenutzer hinzufügen
INSERT INTO users (username) VALUES 
    ('Alice'),
    ('Bob'),
    ('Charlie'),
    ('David'),
    ('Eve');

-- Beispielnachrichten hinzufügen
INSERT INTO messages (sender_id, content) VALUES
    (1, 'Hallo zusammen! Wie geht es euch allen?'),
    (2, 'Hi Alice! Mir geht es gut, danke! Was gibt’s Neues?'),
    (3, 'Hey zusammen, ich bin gerade im Urlaub und genieße die Sonne!'),
    (4, 'Klingt super, Charlie! Wo bist du gerade?'),
    (1, 'Ich plane ein neues Projekt, wäre cool, wenn ihr dabei seid!'),
    (5, 'Hey Leute! Habt ihr die Neuigkeiten über das neue Update gehört?'),
    (2, 'Ja, das Update sieht spannend aus. Ich freue mich darauf!'),
    (3, 'Oh, ich habe es noch nicht gesehen. Was sind die neuen Features?'),
    (1, 'Ich glaube, es gibt eine verbesserte Chat-Funktion und neue Emojis!'),
    (4, 'Neue Emojis? Das klingt gut, freue mich darauf! 😊');
