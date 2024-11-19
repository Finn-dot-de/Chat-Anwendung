document.addEventListener("DOMContentLoaded", function () {
    // Container-Elemente
    const registerContainer = document.getElementById("registrierungs-container");
    const loginContainer = document.getElementById("login-container");
    const chatContainer = document.querySelector(".chat-container");
  
    // Umschalt-Links
    const switchToLogin = document.getElementById("switch-to-login");
    const switchToRegister = document.getElementById("switch-to-register");
  
    // Registrierung Event
    document
      .getElementById("username-submit")
      .addEventListener("click", async () => {
        const username = document.getElementById("username-input").value;
        const password = document.getElementById("password-input").value;
  
        if (!username || !password) {
          alert("Bitte gib einen Namen und ein Passwort ein!");
          return;
        }
  
        const response = await fetch("/api/create/user", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ username, password }),
        });
  
        if (response.ok) {
          localStorage.setItem("username", username);
          switchToChat();
        } else {
          alert("Fehler beim Speichern des Benutzernamens.");
        }
      });
  
    // Login Event
    document.getElementById("login-submit").addEventListener("click", async () => {
      const username = document.getElementById("login-username").value;
      const password = document.getElementById("login-password").value;
  
      if (!username || !password) {
        alert("Bitte gib deinen Benutzernamen und dein Passwort ein!");
        return;
      }
  
      const response = await fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });
  
      if (response.ok) {
        const data = await response.json();
        localStorage.setItem("username", data.username);
        alert("Login erfolgreich!");
        switchToChat();
      } else {
        alert("Login fehlgeschlagen. Benutzername oder Passwort ist falsch.");
      }
    });
  
    // Nachricht senden
    document.getElementById("send-button").addEventListener("click", async () => {
      const message = document.getElementById("chat-input").value;
      const username = localStorage.getItem("username");
  
      if (!message) {
        alert("Bitte gib eine Nachricht ein!");
        return;
      }
  
      await fetch("/api/message", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, content: message }),
      });
  
      document.getElementById("chat-input").value = "";
    });
  
    // Umschalten zwischen Login und Registrierung
    switchToLogin.addEventListener("click", (e) => {
      e.preventDefault();
      showLogin();
    });
  
    switchToRegister.addEventListener("click", (e) => {
      e.preventDefault();
      showRegister();
    });
  
    // Helper-Funktionen
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
    }
  
    // Initialzustand setzen
    showRegister(); // Standardmäßig Registrierung anzeigen
  });
  