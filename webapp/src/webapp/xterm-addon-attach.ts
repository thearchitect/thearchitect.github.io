/**
 * Copyright (c) 2014, 2019 The xterm.js authors. All rights reserved.
 * @license MIT
 *
 * Implements the attach method, that attaches the terminal to a WebSocket stream.
 */

import { Terminal, IDisposable, ITerminalAddon } from 'xterm';

interface IAttachOptions {
    bidirectional?: boolean;
}

// FIXME https://stackoverflow.com/questions/11575558/is-it-possible-to-send-resize-pty-command-with-telnetlib

export class AttachAddon implements ITerminalAddon {
    private _socket: WebSocket;
    private _bidirectional: boolean;
    private _disposables: IDisposable[] = [];

    constructor(socket: WebSocket, options?: IAttachOptions) {
        this._socket = socket;
        // always set binary type to arraybuffer, we do not handle blobs
        this._socket.binaryType = 'arraybuffer';
        this._bidirectional = (options && options.bidirectional === false) ? false : true;

        // setInterval(() => {
        //     this._socket.send(new Uint8Array([255, 250, 31, 0, 40, 0, 40, 255, 240]))
        // }, 2000)
    }

    public activate(terminal: Terminal): void {
        this._disposables.push(
            addSocketListener(this._socket, 'message', ev => {
                const data: ArrayBuffer | string = ev.data;
                terminal.write(typeof data === 'string' ? data : new Uint8Array(data));
            })
        );

        if (this._bidirectional) {
            this._disposables.push(terminal.onData(data => this._sendData(data)));
            this._disposables.push(terminal.onBinary(data => this._sendBinary(data)));
        }

        this._disposables.push(addSocketListener(this._socket, 'close', () => this.dispose()));
        this._disposables.push(addSocketListener(this._socket, 'error', () => this.dispose()));
    }

    public dispose(): void {
        this._disposables.forEach(d => d.dispose());
    }

    private _sendData(data: string): void {
        // TODO: do something better than just swallowing
        // the data if the socket is not in a working condition
        if (this._socket.readyState !== 1) {
            console.warn('the socket is not in a working condition:', this._socket.readyState)
            return;
        }
        this._socket.send(data);
    }

    private _sendBinary(data: string): void {
        if (this._socket.readyState !== 1) {
            return;
        }
        const buffer = new Uint8Array(data.length);
        for (let i = 0; i < data.length; ++i) {
            buffer[i] = data.charCodeAt(i) & 255;
        }
        this._socket.send(buffer);
    }
}

function addSocketListener<K extends keyof WebSocketEventMap>(socket: WebSocket, type: K, handler: (this: WebSocket, ev: WebSocketEventMap[K]) => any): IDisposable {
    socket.addEventListener(type, handler);
    return {
        dispose: () => {
            if (!handler) {
                // Already disposed
                return;
            }
            socket.removeEventListener(type, handler);
        }
    };
}
