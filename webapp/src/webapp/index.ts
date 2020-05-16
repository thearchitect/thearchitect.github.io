import {Terminal} from 'xterm'
import {AttachAddon} from './xterm-addon-attach';
import {FitAddon} from "./xterm-addon-fit";
import {AutoWebSocketBase} from "./websocket";
import {Unicode11Addon} from 'xterm-addon-unicode11';
import {WebLinksAddon} from 'xterm-addon-web-links';

import './xterm.styl'
import './styles.styl'

// const websocketEP = `${(location.protocol === 'https:' ? 'wss' : 'ws')}://${location.host}/ws`
const websocketEP = `${(location.protocol === 'https:' ? 'wss' : 'ws')}://thearchitect.themake.rs/`

interface Size {
    w: number
    h: number
}

class WS extends AutoWebSocketBase {
    constructor(terminal: Terminal) {
        super({
            endpoint: websocketEP,
            active: false,
            reconnectInterval: 1000,
        })

        this.terminal = terminal
    }

    private readonly terminal: Terminal

    private attach: (AttachAddon | null) = null
    private _socket: (WebSocket | null) = null

    private _size: Size = {w: 80, h: 25}

    private sendSize() {
        if (this._socket != null) {
            this._socket.send('\x1b]' + JSON.stringify(this._size))
        }
    }

    public set size(size: Size) {
        this._size = size
        this.sendSize()
    }

    protected connected(): void {
        console.log('connected')

        if (this.attach != null) {
            this.attach.dispose()
        }

        this._socket = this.socket

        if (this._socket != null) {
            this.attach = new AttachAddon(this._socket)
            this.attach.activate(this.terminal)
        } else {
            throw new Error()
        }

        this.sendSize()
    }

    protected disconnected(evt: CloseEvent): void {
        console.log('disconnected', evt)

        this._socket = null

        if (this.attach != null) {
            this.attach.dispose()
            this.attach = null
        }
    }
}


// fontFamily: '"Fixedsys Excelsior 3.01", monospace',
//
// regular: 16px (12pt)
// medium: 32px (24pt)
// large: 48px (36pt)
// larger: 64px (48pt)
// huge: 96px (72pt)
//
// font-weight:400;

const letsGetPartyStarted = () => {
    const terminal = new Terminal({
        convertEol: true,
        rows: 80,
        cols: 40,
        fontFamily: '"Fixedsys Excelsior 3.01", "Inconsolata", monospace',
        fontSize: 16,
        fontWeight: '400',
        fontWeightBold: '600',
        cursorBlink: true,
        cursorStyle: 'block',
        bellStyle: 'sound',
        letterSpacing: 0,
        // windowsMode: true,

        rendererType: 'canvas', // dom / canvas

        windowOptions: {
            // setWinSizePixels: true,
            refreshWin: true,
            // setWinSizeChars: true,
            // maximizeWin: true,
            // fullscreenWin: true,
            getWinSizePixels: true,
            getCellSizePixels: true,
            getWinSizeChars: true,
            pushTitle: true,
            popTitle: true,
            setWinLines: false,
        },
    })

    const fitAddon = new FitAddon()
    fitAddon.activate(terminal)

    terminal.loadAddon(new WebLinksAddon());

    {
        const tel = window.document.getElementById('terminal')
        if (tel != null) {
            terminal.open(tel)
        } else {
            throw new Error('no terminal host element')
        }
    }

    const unicode11Addon = new Unicode11Addon();
    unicode11Addon.activate(terminal)
    terminal.unicode.activeVersion = '11';


    const ws = new WS(terminal)
    ws.active = true

    const fit = () => {
        const sz = fitAddon.fit()
        if (sz) {
            console.log('new-size', sz[0], sz[1])
            ws.size = {w: sz[0], h: sz[1]}
        }
    }

    window.onresize = (_: UIEvent) => {
        fit()
    }

    setTimeout(fit, 1)
    // setTimeout(fit, 500)
    setInterval(fit, 10000)
}

Object.defineProperty(window, 'letsGetPartyStarted', {
    writable: false,
    configurable: false,
    enumerable: false,
    value: letsGetPartyStarted
})

letsGetPartyStarted()
