import socketManager from './websocket.js';
import {terminal} from './terminal.js';
// Import for side effects: sets up settings menu event listeners
import "./settings/settingsMenu.js";
import showAsciiArt from "./welcome.js";

const {term, fitAddon} = terminal;

document.addEventListener('DOMContentLoaded', () => {
    const terminalElement = document.getElementById('terminal');

    setTimeout(() => {
        term.open(terminalElement);
        showAsciiArt(term);
        socketManager.connect();

        fitAddon.fit();
        if (socketManager.getSocket()?.readyState === WebSocket.OPEN) {
            socketManager.sendResize(term.cols, term.rows);
        } else {
            setTimeout(() => socketManager.sendResize(term.cols, term.rows), 100);
        }

        registerHandlers();
        window.addEventListener('resize', fitAndNotify);
    }, 0);


});

function fitAndNotify() {
    fitAddon.fit()
    const socket = socketManager.getSocket();
    if (socket && socket.readyState === WebSocket.OPEN) {
        socketManager.sendResize(term.cols, term.rows);
    }
}

function registerHandlers() {
    term.onData((data) => {
        const socket = socketManager.getSocket();
        if (!socket || socket.readyState !== WebSocket.OPEN) {
            return;
        }

        sendInput(data);
    });
}

function sendInput(data) {
    const encoder = new TextEncoder();
    const inputBytes = encoder.encode(data);
    const buffer = new Uint8Array(1 + inputBytes.length);
    buffer[0] = 0x00; // input tag
    buffer.set(inputBytes, 1);
    socketManager.getSocket().send(buffer);
}

document.getElementById('endButton').addEventListener('click', () => {
    const ok = confirm("End session?");
    if (!ok) return;

    const socket = socketManager.getSocket();
    if (socket && socket.readyState === WebSocket.OPEN) {
        socketManager.sendClose();
        term.write("\r\n\x1b[31mSession ended.\x1b[0m\r\n");
    } else {
        console.warn("WebSocket not connected.");
    }
});

document.getElementById('sudoButton').addEventListener('click', () => {
    const user = prompt("Enter user to switch to:", "root");
    if (!user) return;

    const encoder = new TextEncoder();
    const payload = encoder.encode(user);

    const buffer = new Uint8Array(1 + payload.length);
    buffer[0] = 0x11;  // Sudo header
    buffer.set(payload, 1);

    const socket = socketManager.getSocket();
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(buffer);
        console.log(`Sent sudo request for user: ${user}`);
    } else {
        console.warn("WebSocket not connected.");
    }
});

const scriptsButton = document.getElementById('scriptsDropdownButton');
const scriptsDropdown = document.getElementById('scriptsDropdown');

scriptsButton.addEventListener('click', () => {
    scriptsDropdown.parentElement.classList.toggle('show');
});

document.addEventListener('click', (event) => {
    if (!scriptsButton.contains(event.target) && !scriptsDropdown.contains(event.target)) {
        scriptsDropdown.parentElement.classList.remove('show');
    }
});

document.getElementById('keiGPTButton').addEventListener('click', () => {
    const keiGPTPanel = document.getElementById('keiGPTPanel');
    keiGPTPanel.classList.toggle('show');
    if (keiGPTPanel.classList.contains('show')) {
        keiGPTPanel.scrollIntoView({behavior: 'smooth'});
    }
});

document.getElementById('closeKeiGPT').addEventListener('click', () => {
    document.getElementById('keiGPTPanel').classList.remove('show');
});

// Dummy terminal observer (simulate GPT suggestions)
term.onData((input) => {
    const keiGPTOutput = document.getElementById('keiGPTOutput');
    const p = document.createElement('p');
    p.textContent = `👀 KeiGPT saw: ${input}`;
    keiGPTOutput.appendChild(p);
    keiGPTOutput.scrollTop = keiGPTOutput.scrollHeight;
});

