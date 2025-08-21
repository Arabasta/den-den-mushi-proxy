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

        const jwt2 = "kei"

        // connect to proxy
        let websocketUrl = `ws://${window.location.hostname}:45007/v1/ws?Authorization=${encodeURIComponent(jwt2)}`;
        console.log(websocketUrl);
        this.socket = new WebSocket(websocketUrl, ['X-Proxy-Session-Token', jwt]);

        this.socket.onmessage = (event) => {
            const buffer = new Uint8Array(event.data);
            const header = buffer[0];
            const payload = buffer.slice(1);

            switch (header) {
                case 0x01: // custom protocol header for output
                    term.write(new TextDecoder().decode(payload));
                    break;
                case 0x03: // blocked control char detected
                    showToast(controlCharToString(payload))
                    break;
                case 0x04: // blocked command detected
                    term.write(`\r\n\x1b[31m${new TextDecoder().decode(payload)} not allowed\x1b[0m\r\n`);
                    break;
                case 0x07: // pty session event
                    // note to frontend dev, should NOT use the same toast as blocked control char
                    showSideToast(new TextDecoder().decode(payload));
                    break;
                case 0x13: // pty normal close
                    term.write(`\r\n\x1b[32mSession ended normally, shutting down session.\r\nSessionId: ${new TextDecoder().decode(payload)}\x1b[0m\r\n`);
                    break;
                case 0x14: // pty error close
                    term.write(`\r\n\x1b[31mError reading from pty, shutting down session.\r\nSessionId: ${new TextDecoder().decode(payload)}\x1b[0m\r\n`);
                    break;
                case 0x15: // pty connection success
                    term.write(`\r\n\x1b[32mConnected to terminal session.\r\nSessionId: ${new TextDecoder().decode(payload)}\x1b[0m\r\n`);
                    break;
                case 0x16: // pty CR timeout
                    term.write(`\r\n\x1b[33mCR end time reached, bye bye.\r\nSessionId: ${new TextDecoder().decode(payload)}\x1b[0m\r\n`);
                    break;
                case 0x17: // CR timeout warning
                    // note to frontend dev, should NOT use the same toast as blocked control char
                    showToast("CR session will end in " + new TextDecoder().decode(payload) + " minutes");
                    break;
                case 0x19: // no active observer, packet dropped
                    showToast("No active observer...");
                    break;
                default:
                    console.warn("Unknown header:", header);
            }
        };

        this.socket.onopen = () => {
            console.log("WebSocket connection opened");
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

    sendClose() {
        const buffer = new Uint8Array(1);
        buffer[0] = 0x18;
        this.socket.send(buffer);
        this.socket.close();
    },

    getSocket() {
        return this.socket;
    }
};

export default socketManager;

// just whatever lol frontend no one cares
function showToast(message) {
    const toast = document.createElement('div');
    toast.innerText = "Blocked control character: " + message;
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
    setTimeout(() => toast.remove(), 1000);
}

function showSideToast(message, type = 'info') {
    const toast = document.createElement('div');
    toast.innerHTML = `
        <div class="toast-content">
            <div class="toast-message">${message}</div>
        </div>
    `;

    // Base styles
    toast.style.cssText = `
        position: fixed;
        right: 20px;
        top: 20px;
        width: 300px;
        padding: 16px;
        border-radius: 4px;
        box-shadow: 0 3px 10px rgba(0,0,0,0.2);
        z-index: 9999;
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
        display: flex;
        align-items: center;
        transform: translateX(120%);
        transition: transform 0.3s ease-out;
        opacity: 0;
    `;

    // Type-specific styles
    const typeStyles = {
        info: `background: #f0f7ff; color: #0066cc; border-left: 4px solid #0066cc;`,
        success: `background: #f0fff4; color: #0f9960; border-left: 4px solid #0f9960;`,
        warning: `background: #fffaf0; color: #d9822b; border-left: 4px solid #d9822b;`,
        error: `background: #fff0f0; color: #db3737; border-left: 4px solid #db3737;`
    };

    toast.style.cssText += typeStyles[type] || typeStyles.info;

    document.body.appendChild(toast);

    // Trigger the slide-in animation
    setTimeout(() => {
        toast.style.transform = 'translateX(0)';
        toast.style.opacity = '1';
    }, 10);

    // Auto-remove after delay
    setTimeout(() => {
        toast.style.transform = 'translateX(120%)';
        toast.style.opacity = '0';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

const controlCharToString = (payload) => {
    if (!payload?.length) return "Empty input";

    // Handle single-byte input
    if (payload.length === 1) {
        const code = payload[0];
        const namedControl = {
            3: "Ctrl+C",
            18: "Ctrl+R",
            21: "Ctrl+U",
            26: "Ctrl+Z"
        }[code];

        if (namedControl) return namedControl;

        // If printable ASCII (32–126), return as char
        if (code >= 32 && code <= 126) {
            return String.fromCharCode(code);
        }

        return `Control character (code ${code})`;
    }

    // Handle known arrow escape sequences
    if (payload.length === 3 && payload[0] === 27 && payload[1] === 91) {
        return {
            65: "Arrow Up",
            66: "Arrow Down",
            67: "Arrow Right",
            68: "Arrow Left",
            70: "End",
            72: "Home",

        }[payload[2]] ?? `Special key sequence (${payload.join(',')})`;
    }

    // Default fallback
    return `Unknown sequence (${payload.join(',')})`;
};
