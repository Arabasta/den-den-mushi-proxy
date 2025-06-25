import {terminal} from "./terminal.js";

const {term, fitAddon} = terminal;

const socketManager = {
    socket: null,

    connect() {
        if (this.socket?.readyState === WebSocket.OPEN) return;

        // get the jwt from url param
        // todo: dont do it like this lol
        const params = new URLSearchParams(location.search);
        const jwt = params.get('jwt_tmp');

        // connect to proxy
        let websocketUrl = `ws://${window.location.hostname}:45007/ws`;
        console.log(websocketUrl);
        this.socket = new WebSocket(websocketUrl, ['jwt', jwt]);

        this.socket.onmessage = (event) => {
            const buffer = new Uint8Array(event.data);
            const header = buffer[0];
            const payload = buffer.slice(1);

            switch (header) {
                case 0x01: // custom protocol header for output
                    term.write(new TextDecoder().decode(payload));
                    break;
                case 0x03: // blocked control char detected
                    showToast(new TextDecoder().decode(payload))
                    break;
                case 0x04: // blocked command detected
                    term.write(`\r\n\x1b[31m${new TextDecoder().decode(payload)}\x1b[0m\r\n`);
                    break;
                default:
                    console.warn("Unknown header:", header);
            }
        };

        this.socket.onopen = () => {
            console.log("WebSocket connection opened");
            term.write('\r\n\x1b[32mConnected to terminal session\x1b[0m\r\n');
            term.focus();
        };

        this.socket.onclose = ({code, reason}) => {
        };

        this.socket.onerror = (error) => {
            term.write(`\r\n\x1b[31mError: ${error.message}\x1b[0m\r\n`);
        };

        this.socket.binaryType = "arraybuffer";
    },

    sendResize(cols, rows) {
        const buffer = new Uint8Array(5);
        buffer[0] = 0x10; // custom protocol header for resize

        // backend expects big endian
        buffer[1] = (cols >> 8) & 0xff;
        buffer[2] = cols & 0xff;
        buffer[3] = (rows >> 8) & 0xff;
        buffer[4] = rows & 0xff;
        this.socket.send(buffer);
    },

    getSocket() {
        return this.socket;
    }
};

export default socketManager;

// just whatever lol frontend no one cares
function showToast(message) {
    const toast = document.createElement('div');
    toast.innerText = message;
    toast.style = `
        position: fixed;
        left: 50%;
        transform: translateX(-50%);
        background-color: #f87171;
        color: white;
        padding: 10px 10px;
        border-radius: 6px;
        box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        z-index: 9999;
        font-weight: bold;
    `;
    document.body.appendChild(toast);
    setTimeout(() => toast.remove(), 10000);
}
