import socketManager from './websocket.js';
import {terminal} from './terminal.js';

const {term, fitAddon} = terminal;

document.addEventListener('DOMContentLoaded', () => {
    const terminalElement = document.getElementById('terminal');

    // delay so i can open dev tools
    setTimeout(() => {
        term.open(terminalElement);
        socketManager.connect();

        fitAddon.fit();
        if (socketManager.getSocket()?.readyState === WebSocket.OPEN) {
            socketManager.sendResize(term.cols, term.rows);
        } else {
            setTimeout(() => socketManager.sendResize(term.cols, term.rows), 100);
        }

        registerHandlers();
        window.addEventListener('resize', fitAndNotify);
    }, 1000);


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