// Modul: Utils - Hilfsfunktionen und globale Variablen
const Utils = (() => {
  const chatMessages = document.getElementById("chat-messages");

  function appendMessageToChat(message) {
    const currentUsername = sessionStorage.getItem("user");
    const messageData = JSON.parse(message);

    // Nachrichtenelement erstellen
    const messageElement = document.createElement("div");
    messageElement.className = "message";

    if (messageData.username === currentUsername) {
      messageElement.classList.add("sent");
    } else {
      messageElement.classList.add("received");
    }

    const messageContent = document.createElement("p");
    const usernameElement = document.createElement("span");
    usernameElement.className = "username-in-message";
    usernameElement.textContent =
      currentUsername === messageData.username ? "Du" : messageData.username;
    messageContent.appendChild(usernameElement);

    const messageText = document.createTextNode(`\n${messageData.content}`);
    messageContent.appendChild(messageText);
    messageElement.appendChild(messageContent);

    const timestamp = document.createElement("span");
    timestamp.className = "timestamp";
    timestamp.textContent = new Date(messageData.timestamp).toLocaleTimeString(
      "de-DE",
      { hour: "2-digit", minute: "2-digit" }
    );
    messageElement.appendChild(timestamp);

    chatMessages.appendChild(messageElement);
    chatMessages.scrollTop = chatMessages.scrollHeight;
  }

  return { appendMessageToChat };
})();

// Modul: API - Serveranfragen
const API = (() => {
  async function registerUser(username, password) {
    const response = await fetch("/api/create/user", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });
    return response;
  }

  async function loginUser(username, password) {
    const response = await fetch("/api/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });
    return response;
  }

  async function sendMessage(username, message) {
    const response = await fetch("/api/new/message", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, content: message }),
    });
    return response;
  }

  return { registerUser, loginUser, sendMessage };
})();

// Modul: SSE - Server-Sent Events
const SSE = (() => {
  function connectToEvents() {
    const eventSource = new EventSource("/api/events");

    eventSource.onmessage = function (event) {
      const message = event.data;
      Utils.appendMessageToChat(message);
    };

    eventSource.onerror = function () {
      eventSource.close();
      setTimeout(connectToEvents, 5000); // Neuer Versuch nach 5 Sekunden
    };
  }

  return { connectToEvents };
})();

// Modul: UI - Benutzeroberfläche und Interaktionen
const UI = (() => {
  const registerContainer = document.getElementById("registrierungs-container");
  const loginContainer = document.getElementById("login-container");
  const chatContainer = document.querySelector(".chat-container");

  function showRegister() {
    loginContainer.style.display = "none";
    registerContainer.style.display = "block";
  }

  function showLogin() {
    registerContainer.style.display = "none";
    loginContainer.style.display = "block";
  }

  function switchToChat() {
    registerContainer.style.display = "none";
    loginContainer.style.display = "none";
    chatContainer.style.display = "flex";
    SSE.connectToEvents();
  }

  function adjustChatContainer() {
    const screenWidth = window.innerWidth;

    if (screenWidth < 480) {
      chatContainer.style.width = "100vw";
      chatContainer.style.height = "100vh";
      chatContainer.style.borderRadius = "0";
    } else if (screenWidth < 768) {
      chatContainer.style.width = "90vw";
      chatContainer.style.height = "90vh";
    } else {
      chatContainer.style.width = "60vw";
      chatContainer.style.height = "80vh";
    }
  }

  return { showRegister, showLogin, switchToChat, adjustChatContainer };
})();

// Modul: EventHandler - Event-Listener für Benutzeraktionen
const EventHandler = (() => {
  function initialize() {
    // Umschalt-Links
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
          sessionStorage.setItem("user", username);
          UI.switchToChat();
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
          alert("Login erfolgreich!");
          UI.switchToChat();
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
        document.getElementById("chat-input").value = "";
      });

    // Fenstergröße überwachen
    window.addEventListener("resize", UI.adjustChatContainer);
    document.addEventListener("DOMContentLoaded", UI.adjustChatContainer);

    // Standardmäßig Registrierungsformular anzeigen
    UI.showRegister();
  }

  return { initialize };
})();

// Anwendung starten
document.addEventListener("DOMContentLoaded", EventHandler.initialize);
