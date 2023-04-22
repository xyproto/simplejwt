const el = (id) => document.getElementById(id);

const setFormVisibility = (visible) => {
    const display = visible ? "block" : "none";
    el("registerForm").style.display = display;
    el("loginForm").style.display = display;
    el("logoutButton").style.display = visible ? "none" : "block";
    el("sendMessageForm").style.display = visible ? "none" : "block";
    el("messages").style.display = visible ? "none" : "block";
};

const fetchAPI = async (url, method, headers, body) => {
    const res = await fetch(url, { method, headers, body });
    if (res.ok) return res;

    const errMsg = await res.text();

    if (errMsg === "Invalid or expired token") {
        logout();
    } else {
        setStatusMessage(`Error: ${errMsg}`, true);
    }

    return null;
};

async function updateUI(fn) {
    try {
        const res = await fn();
        if (res) setStatusMessage("");
    } catch (err) {
        setStatusMessage(`Error: ${err.message}`, true);
    }
}

async function register(event) {
    event.preventDefault();
    const nickname = el("registerNickname").value;
    const password = el("registerPassword").value;
    await updateUI(async () => {
        const res = await fetchAPI(
            "/register",
            "POST",
            { "Content-Type": "application/json" },
            JSON.stringify({ nickname, password })
        );
        if (res) {
            setStatusMessage("Successfully registered!");
            el("registerForm").reset();
        }
    });
}

async function login(event) {
    event.preventDefault();
    const nickname = el("loginNickname").value;
    const password = el("loginPassword").value;
    await updateUI(async () => {
        const res = await fetchAPI(
            "/login",
            "POST",
            { "Content-Type": "application/json" },
            JSON.stringify({ nickname, password })
        );
        if (res) {
            const token = await res.text();
            localStorage.setItem("token", token);
            setFormVisibility(false);
            fetchMessages();
            startFetchingMessages();
            el("logoutButton").dataset.nickname = nickname;
        }
    });
}

function handleKeyPress(event) {
    if (event.key === "Enter" && !event.shiftKey) {
        event.preventDefault();
        sendMessage(event);
    }
}

function logout() {
    localStorage.removeItem("token");
    setFormVisibility(true);
    el("messages").innerHTML = "";
    stopFetchingMessages();
}

async function sendMessage(event) {
    event.preventDefault();
    const content = el("messageContent").value;
    const token = localStorage.getItem("token");
    await updateUI(async () => {
        const res = await fetchAPI(
            "/send",
            "POST",
            {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            JSON.stringify({ content })
        );
        if (res) el("sendMessageForm").reset();
    });
}

function formatTimestamp(timestamp) {
    const date = new Date(timestamp);
    return `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`;
}

async function fetchMessages() {
    const token = localStorage.getItem("token");
    const res = await fetchAPI("/messages", "GET", {
        Authorization: `Bearer ${token}`,
    });
    if (!res) return;
    const data = await res.json();
    const prevMsgs = el("messages").innerHTML;
    el("messages").innerHTML = data
        .map(
            (msg) =>
                `<p><strong>${msg.Sender} (${formatTimestamp(
                    msg.Timestamp
                )}):</strong> ${msg.Content}</p>`
        )
        .join("");
    el("messages").scrollTop = el("messages").scrollHeight;
    showNotifications(data, prevMsgs);
}

let fetchInterval;
function startFetchingMessages() {
    fetchInterval = setInterval(fetchMessages, 1000);
}
function stopFetchingMessages() {
    clearInterval(fetchInterval);
}

function setStatusMessage(message, isError = false) {
    const statusMessage = el("statusMessage");
    statusMessage.textContent = message;
    statusMessage.classList.toggle("error", isError);
    setTimeout(() => {
        statusMessage.textContent = "";
        statusMessage.classList.remove("error");
    }, 3000);
}

function showNotifications(data, prevMsgs) {
    if (!("Notification" in window) || Notification.permission === "denied")
        return;
    if (Notification.permission === "default") Notification.requestPermission();
    const loggedInUser = el("logoutButton").dataset.nickname;

    if (loggedInUser) {
        const prevMsgArray = prevMsgs
            .split("</p>")
            .map((msg) => msg.replace("<p>", "").trim())
            .map((msg) => msg.replace("<strong>", "").replace("</strong>", ""));

        data.forEach((msg) => {
            const messageContent = `${msg.Sender}: ${msg.Content}`;
            if (
                msg.Content.toLowerCase().includes(
                    loggedInUser.toLowerCase()
                ) &&
                !prevMsgArray.includes(messageContent)
            ) {
                console.log(`Sending notification for: ${messageContent}`);
                setTimeout(() => {
                    new Notification("New mention in Tiny Chat", {
                        body: messageContent,
                    });
                }, 100);
            } else {
                console.log(`No notification for: ${messageContent}`);
            }
        });
    }
}

el("registerForm").addEventListener("submit", register);
el("loginForm").addEventListener("submit", login);
el("sendMessageForm").addEventListener("submit", sendMessage);
el("logoutButton").addEventListener("click", logout);
el("messageContent").addEventListener("keypress", handleKeyPress);

if (localStorage.getItem("token")) {
    setFormVisibility(false);
    startFetchingMessages();
}

function showLoginAndRegisterForms() {
    if (!localStorage.getItem("token")) {
        setFormVisibility(true);
    } else {
        setFormVisibility(false);
        startFetchingMessages();
    }
}

showLoginAndRegisterForms();
