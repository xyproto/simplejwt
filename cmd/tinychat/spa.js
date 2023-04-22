const registerForm = document.getElementById("registerForm");
const loginForm = document.getElementById("loginForm");
const sendMessageForm = document.getElementById("sendMessageForm");
const logoutButton = document.getElementById("logoutButton");
const messages = document.getElementById("messages");

async function register(event) {
    event.preventDefault();

    const nickname = document.getElementById("registerNickname").value;
    const password = document.getElementById("registerPassword").value;
    const response = await fetch("/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ nickname, password }),
    });

    if (response.status === 201) {
        alert("Successfully registered!");
        registerForm.reset();
    } else {
        alert("Error: " + (await response.text()));
    }
}

async function login(event) {
    event.preventDefault();

    const nickname = document.getElementById("loginNickname").value;
    const password = document.getElementById("loginPassword").value;
    const response = await fetch("/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ nickname, password }),
    });

    if (response.status === 200) {
        const token = await response.text();
        localStorage.setItem("token", token);
        loginForm.reset();
        loginForm.style.display = "none";
        registerForm.style.display = "none";
        logoutButton.style.display = "block";
        sendMessageForm.style.display = "block";
        fetchMessages();
        startFetchingMessages();
    } else {
        alert("Error: " + (await response.text()));
    }
}

function handleKeyPress(event) {
    if (event.key === "Enter" && !event.shiftKey) {
        event.preventDefault();
        sendMessage(event);
    }
}

async function logout() {
    localStorage.removeItem("token");
    loginForm.style.display = "block";
    registerForm.style.display = "block";
    logoutButton.style.display = "none";
    sendMessageForm.style.display = "none";
    messages.innerHTML = ""; // Clear the messages
    stopFetchingMessages();
}

async function sendMessage(event) {
    event.preventDefault();

    const content = document.getElementById("messageContent").value;
    const token = localStorage.getItem("token");

    const response = await fetch("/send", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            Authorization: "Bearer " + token,
        },
        body: JSON.stringify({ content }),
    });

    if (response.status === 201) {
        sendMessageForm.reset();
    } else {
        alert("Error: " + (await response.text()));
    }
}

async function fetchMessages() {
    const token = localStorage.getItem("token");

    const response = await fetch("/messages", {
        method: "GET",
        headers: { Authorization: "Bearer " + token },
    });

    if (response.status === 200) {
        const data = await response.json();
        messages.innerHTML = data
            .map(
                (msg) => `<p><strong>${msg.Sender}:</strong> ${msg.Content}</p>`
            )
            .join("");
    } else {
        alert("Error: " + (await response.text()));
    }
}

let fetchInterval;

function startFetchingMessages() {
    fetchInterval = setInterval(fetchMessages, 1000);
}

function stopFetchingMessages() {
    clearInterval(fetchInterval);
}

registerForm.addEventListener("submit", register);
loginForm.addEventListener("submit", login);
sendMessageForm.addEventListener("submit", sendMessage);
logoutButton.addEventListener("click", logout);

document
    .getElementById("messageContent")
    .addEventListener("keypress", handleKeyPress);

if (localStorage.getItem("token")) {
    loginForm.style.display = "none";
    registerForm.style.display = "none";
    logoutButton.style.display = "block";
    sendMessageForm.style.display = "block";
    startFetchingMessages();
}
