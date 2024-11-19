document.getElementById('username-submit').addEventListener('click', async () => {
    const username = document.getElementById('username-input').value;
    const password = document.getElementById('password-input').value;

    if (!username || !password) {
        alert('Bitte gib einen Namen und ein Passwort ein!');
        return;
    }

    const response = await fetch('/api/user', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
    });

    if (response.ok) {
        localStorage.setItem('username', username);
        document.getElementById('username-container').style.display = 'none';
        document.querySelector('.chat-container').style.display = 'block';
    } else {
        alert('Fehler beim Speichern des Benutzernamens.');
    }
});

document.getElementById('send-button').addEventListener('click', async () => {
    const message = document.getElementById('chat-input').value;
    const username = localStorage.getItem('username');

    if (!message) {
        alert('Bitte gib eine Nachricht ein!');
        return;
    }

    await fetch('/api/message', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, content: message }),
    });

    document.getElementById('chat-input').value = '';
});

document.getElementById('login-submit').addEventListener('click', async () => {
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;

    if (!username || !password) {
        alert('Bitte gib deinen Benutzernamen und dein Passwort ein!');
        return;
    }

    const response = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
    });

    if (response.ok) {
        const data = await response.json();
        localStorage.setItem('username', data.username);
        alert('Login erfolgreich!');
        document.getElementById('login-container').style.display = 'none';
        document.querySelector('.chat-container').style.display = 'block';
    } else {
        alert('Login fehlgeschlagen. Benutzername oder Passwort ist falsch.');
    }
});

