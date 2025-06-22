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
            console.log('Received from server:', event.data);
            term.write(event.data);
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
    },

    sendResize(cols, rows) {
        const buffer = new Uint8Array(5);
        buffer[0] = 0x01; // custom protocol header for resize

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
