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
    el("messages").style.display = "none";
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
    const options = { weekday: "long", day: "2-digit", month: "long" };
    const dateString = date.toLocaleDateString("en-US", options);
    const dayNumSuffix = getDayNumSuffix(date.getDate());
    return `${dateString.slice(0, -3)}${dayNumSuffix}, ${date.getHours()}:${date
        .getMinutes()
        .toString()
        .padStart(2, "0")}`;
}

function getDayNumSuffix(dayNum) {
    if (dayNum >= 11 && dayNum <= 13) return "th";
    switch (dayNum % 10) {
        case 1:
            return "st";
        case 2:
            return "nd";
        case 3:
            return "rd";
        default:
            return "th";
    }
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
            (msg) => `
                <div class="message-row">
                    <div class="message-info">
                        <span class="timestamp">${formatTimestamp(
                            msg.Timestamp
                        )}</span>
                        <strong>${msg.Sender}</strong>
                    </div>
                    <div class="message-content">${msg.Content}</div>
                </div>`
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

function sendNotification(title, options) {
    if (!("Notification" in window) || Notification.permission === "denied") {
        return;
    }

    if (Notification.permission === "default") {
        Notification.requestPermission().then((permission) => {
            if (permission === "granted") {
                new Notification(title, options);
            }
        });
    } else {
        new Notification(title, options);
    }
}

function showNotifications(data, prevMsgs) {
    const loggedInUser = el("logoutButton").dataset.nickname;

    if (loggedInUser) {
        const prevMsgArray = prevMsgs
            .split("</p>")
            .map((msg) => msg.replace("<p>", "").trim())
            .map((msg) => msg.replace("<strong>", "").replace("</strong>", ""));

        data.forEach((msg) => {
            const messageContent = `${msg.Sender}: ${msg.Content}`;
            const mentionRegex = new RegExp(`\\b${loggedInUser}\\b`, "i");

            if (
                mentionRegex.test(msg.Content) &&
                !prevMsgArray.includes(messageContent)
            ) {
                console.log(`Sending notification for: ${messageContent}`);
                sendNotification("New mention in Tiny Chat", {
                    body: messageContent,
                });
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

function requestNotificationPermission() {
    if (!("Notification" in window)) {
        console.log("This browser does not support notifications.");
        return;
    }

    Notification.requestPermission().then((permission) => {
        if (permission === "granted") {
            console.log("Notification permission granted.");
            // Send a test notification
            new Notification("Tiny Chat", {
                body: "Test notification",
            });
        } else {
            console.log("Notification permission denied.");
        }
    });
}

el("requestNotificationPermission").addEventListener(
    "click",
    requestNotificationPermission
);

function showLoginAndRegisterForms() {
    if (!localStorage.getItem("token")) {
        setFormVisibility(true);
        el("messages").style.display = "none";
    } else {
        setFormVisibility(false);
        el("messages").style.display = "block";
        startFetchingMessages();
    }
}

showLoginAndRegisterForms();
