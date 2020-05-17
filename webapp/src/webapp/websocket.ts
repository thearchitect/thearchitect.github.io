export interface AutoWebSocketBaseOptions {
    endpoint: string
    active?: boolean
    reconnectInterval?: number
}

export abstract class AutoWebSocketBase {
    protected constructor(o: AutoWebSocketBaseOptions) {
        this._endpoint = o.endpoint;

        if (o.active != undefined) {
            this.active = o.active
        }

        if (o.reconnectInterval != undefined) {
            this._reconnectInterval = o.reconnectInterval
        }
    }

    private readonly _endpoint: string
    private readonly _reconnectInterval: number = 1000

    private _reconnectTimeoutID: (number | null) = null
    private __socket: (WebSocket | null) = null

    private _connected: boolean = false
    private _shouldBeActive: boolean = false

    private get _isActive(): boolean {
        return this.__socket != null || this._reconnectTimeoutID != null
    }

    public get active(): boolean {
        return this._shouldBeActive
    }

    public set active(active: boolean) {
        this._shouldBeActive = active
        this.triggerFSM()
    }

    public get socket(): (WebSocket|null) {
        return this.__socket
    }

    private triggerFSM() {
        const scheduleReconnect = () => {
            deactivate()
            this._reconnectTimeoutID = setTimeout(() => {
                if (this._shouldBeActive) {
                    activate()
                }
            }, this._reconnectInterval)
        }

        const activate = () => {
            deactivate()

            console.warn('autows', 'connecting')
            this.__socket = new WebSocket(this._endpoint)
            this.__socket.onclose = (ev: CloseEvent) => {
                console.warn('autows', 'onclose', ev)
                scheduleReconnect()
                if (this._connected) {
                    this._connected = false
                    this.disconnected(ev)
                }
            }
            this.__socket.onerror = (ev: Event) => {
                console.warn('autows', 'onerror', ev)
            }
            this.__socket.onopen = (ev: Event): any => {
                console.warn('autows', 'onopen', ev)
                this._connected = true
                this.connected()
            }
        }

        const deactivate = () => {
            if (this.__socket != null) {
                this.__socket.close()
                this.__socket = null
            }

            if (this._reconnectTimeoutID != null) {
                clearTimeout(this._reconnectTimeoutID)
                this._reconnectTimeoutID = null
            }
        }

        if (this._shouldBeActive && !this._isActive) {
            activate()
        } else if (this._isActive && !this._shouldBeActive) {
            deactivate()
        }
    }

    protected abstract connected(): void

    protected abstract disconnected(ev: CloseEvent): void
}
