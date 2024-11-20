// Modul: Utils - Hilfsfunktionen und globale Variablen
const Utils = (() => {
  // Referenz auf den Chat-Nachrichten-Container
  const chatMessages = document.getElementById("chat-messages");

  /**
   * Fügt eine Nachricht in den Chat ein.
   * @param {string} message - Die Nachricht als JSON-String (z.B. von Server-Sent Events).
   */
  function appendMessageToChat(message) {
    // Der Benutzername des aktuell eingeloggten Nutzers (aus SessionStorage).
    const currentUsername = sessionStorage.getItem("user");
    // Nachrichtendaten als JSON-Objekt umwandeln.
    const messageData = JSON.parse(message);

    // Ein neues Nachrichten-Element erstellen.
    const messageElement = document.createElement("div");
    messageElement.className = "message";

    // Falls die Nachricht vom aktuellen Benutzer stammt, füge die Klasse 'sent' hinzu.
    // Andernfalls füge die Klasse 'received' hinzu.
    if (messageData.username === currentUsername) {
      messageElement.classList.add("sent");
    } else {
      messageElement.classList.add("received");
    }

    // Nachrichtentext erstellen
    const messageContent = document.createElement("p");

    // Benutzername (oder "Du" für eigene Nachrichten) anzeigen
    const usernameElement = document.createElement("span");
    usernameElement.className = "username-in-message";
    usernameElement.textContent =
      currentUsername === messageData.username ? "Du" : messageData.username;
    messageContent.appendChild(usernameElement);

    // Füge den eigentlichen Nachrichtentext hinzu
    const messageText = document.createTextNode(`\n${messageData.content}`);
    messageContent.appendChild(messageText);
    messageElement.appendChild(messageContent);

    // Zeitstempel der Nachricht hinzufügen
    const timestamp = document.createElement("span");
    timestamp.className = "timestamp";
    timestamp.textContent = new Date(messageData.timestamp).toLocaleTimeString(
      "de-DE",
      { hour: "2-digit", minute: "2-digit" }
    );
    messageElement.appendChild(timestamp);

    // Füge das erstellte Nachrichten-Element in den Chat-Container ein
    chatMessages.appendChild(messageElement);

    // Automatisches Scrollen nach unten, um die neueste Nachricht anzuzeigen
    chatMessages.scrollTop = chatMessages.scrollHeight;
  }

  // Exponiere die Funktion appendMessageToChat, damit sie außerhalb des Moduls verwendet werden kann
  return { appendMessageToChat };
})();

// Modul: API - Kommunikation mit dem Server
const API = (() => {
  /**
   * Registriert einen neuen Benutzer.
   * @param {string} username - Benutzername.
   * @param {string} password - Passwort.
   * @returns {Promise<Response>} Die Antwort vom Server.
   */
  async function registerUser(username, password) {
    return await fetch("/api/create/user", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });
  }

  /**
   * Loggt einen Benutzer ein.
   * @param {string} username - Benutzername.
   * @param {string} password - Passwort.
   * @returns {Promise<Response>} Die Antwort vom Server.
   */
  async function loginUser(username, password) {
    return await fetch("/api/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });
  }

  /**
   * Sendet eine neue Nachricht an den Server.
   * @param {string} username - Benutzername des Senders.
   * @param {string} message - Inhalt der Nachricht.
   */
  async function sendMessage(username, message) {
    return await fetch("/api/new/message", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, content: message }),
    });
  }

  // Exponiere die Funktionen für externe Nutzung
  return { registerUser, loginUser, sendMessage };
})();

// Modul: SSE - Server-Sent Events
const SSE = (() => {
  /**
   * Stellt eine Verbindung zu Server-Sent Events (SSE) her.
   * Diese Verbindung wird verwendet, um Echtzeit-Nachrichten vom Server zu empfangen.
   */
  function connectToEvents() {
    const eventSource = new EventSource("/api/events");

    // Event: Nachricht empfangen
    eventSource.onmessage = function (event) {
      const message = event.data; // Die empfangene Nachricht (als JSON-String).
      Utils.appendMessageToChat(message); // Nachricht in den Chat einfügen.
    };

    // Event: Fehler bei der Verbindung
    eventSource.onerror = function () {
      eventSource.close(); // Verbindung schließen
      setTimeout(connectToEvents, 5000); // Nach 5 Sekunden erneut versuchen
    };
  }

  // Exponiere die Funktion connectToEvents
  return { connectToEvents };
})();

// Modul: UI - Benutzeroberfläche
const UI = (() => {
  const registerContainer = document.getElementById("registrierungs-container");
  const loginContainer = document.getElementById("login-container");
  const chatContainer = document.querySelector(".chat-container");

  /**
   * Zeigt das Registrierungsformular an und versteckt andere Bereiche.
   */
  function showRegister() {
    loginContainer.style.display = "none";
    registerContainer.style.display = "block";
  }

  /**
   * Zeigt das Login-Formular an und versteckt andere Bereiche.
   */
  function showLogin() {
    registerContainer.style.display = "none";
    loginContainer.style.display = "block";
  }

  /**
   * Zeigt den Chat-Bereich an und versteckt andere Bereiche.
   */
  function switchToChat() {
    registerContainer.style.display = "none";
    loginContainer.style.display = "none";
    chatContainer.style.display = "flex";
    SSE.connectToEvents(); // Verbindung zu SSE herstellen
  }

  // Exponiere die Funktionen für externe Nutzung
  return { showRegister, showLogin, switchToChat };
})();

// Modul: EventHandler - Event-Listener
const EventHandler = (() => {
  /**
   * Initialisiert alle Event-Listener und setzt den Standardzustand der Seite.
   */
  function initialize() {
    // Umschalt-Links für Login und Registrierung
    document
      .getElementById("switch-to-login")
      .addEventListener("click", (e) => {
        e.preventDefault();
        UI.showLogin();
      });

    document
      .getElementById("switch-to-register")
      .addEventListener("click", (e) => {
        e.preventDefault();
        UI.showRegister();
      });

    // Registrierung
    document
      .getElementById("username-submit")
      .addEventListener("click", async () => {
        const username = document.getElementById("username-input").value;
        const password = document.getElementById("password-input").value;

        if (!username || !password) {
          alert("Bitte gib einen Namen und ein Passwort ein!");
          return;
        }

        const response = await API.registerUser(username, password);
        if (response.ok) {
          sessionStorage.setItem("user", username); // Benutzername speichern
          UI.switchToChat(); // Zum Chat wechseln
        } else {
          alert("Fehler bei der Registrierung.");
        }
      });

    // Login
    document
      .getElementById("login-submit")
      .addEventListener("click", async () => {
        const username = document.getElementById("login-username").value;
        const password = document.getElementById("login-password").value;

        if (!username || !password) {
          alert("Bitte gib deinen Benutzernamen und dein Passwort ein!");
          return;
        }

        const response = await API.loginUser(username, password);
        if (response.ok) {
          sessionStorage.setItem("user", username);
          UI.switchToChat(); // Zum Chat wechseln
        } else {
          alert("Login fehlgeschlagen.");
        }
      });

    // Nachricht senden
    document
      .getElementById("send-button")
      .addEventListener("click", async () => {
        const message = document.getElementById("chat-input").value;
        const username = sessionStorage.getItem("user");

        if (!message) {
          alert("Bitte gib eine Nachricht ein!");
          return;
        }

        await API.sendMessage(username, message);
        document.getElementById("chat-input").value = ""; // Eingabefeld leeren
      });

    // Anpassung der Chatgröße bei Fensteränderung
    window.addEventListener("resize", UI.adjustChatContainer);
    document.addEventListener("DOMContentLoaded", UI.adjustChatContainer);

    // Zeige standardmäßig das Registrierungsformular an
    UI.showRegister();
  }

  // Exponiere die Funktion initialize
  return { initialize };
})();

// Anwendung starten, wenn die Seite geladen ist
document.addEventListener("DOMContentLoaded", EventHandler.initialize);
