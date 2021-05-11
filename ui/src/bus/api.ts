import {EventType} from "mitt";
import {ComponentBus, getBus, BusEvent} from "@/bus/core";
import config from "@/config";

const componentId = 'api'

// event types which are allowed to be exchanged with the api
const allowedEventTypes = [
    'canvas:shape:added',
    'canvas:shape:updated',
    'canvas:shape:destroyed'
]

export class APIAdapter {
    connection : WebSocket
    bus : ComponentBus
    orgSlug : string
    boardSlug : string

    constructor(orgSlug : string, boardSlug : string) {
        this.orgSlug = orgSlug
        this.boardSlug = boardSlug
        this.bus = getBus(componentId)
        this.connection = new WebSocket(`ws://${config.apiHost}/api/v1/orgs/${orgSlug}/boards/${boardSlug}/ws?clientId=${this.bus.clientId}`)

        this.connection.onmessage = this.forwardToBus.bind(this)

        this.bus.onAll(this.forwardToAPI.bind(this))
        this.connection.onopen = (e) => {
            console.log("Successfully connected to the echo websocket server...")
        }
    }

    private forwardToAPI(type: EventType, e?: any) : void {
        const event = <BusEvent>e
        // avoid that only valid events will be forwarded to the API and that an event that is on the bus
        // but not originating from it wont be sent back to the server again.
        if (!APIAdapter.isPermittedType(type) || !this.bus.isOriginClientTo(event)) {
            return
        }
        console.log(`Forwarding from bus to API with type: '${String(type)}'`)
        const msg = new WebSocketEvent(String(type), event.source, event.content)
        this.connection.send(msg.toJSON())
    }

    private forwardToBus(e : MessageEvent<WebSocketEvent>) : void {
        // type the WebSocketEvent back to a BusEvent to drop the type
        const event = <WebSocketEvent> JSON.parse(String(e.data))
        if (!APIAdapter.isPermittedType(event.type)) {
            return
        }
        if (!this.bus.isOriginClientTo(event)) {
            console.log(`Forwarding from API to bus with type: '${event.type}'`)
            this.bus.forward(event.type, <BusEvent> event)
        }
    }

    private static isPermittedType(type : EventType) : boolean {
        return String(type) !== '' && allowedEventTypes.includes(String(type))
    }
}

class WebSocketEvent extends BusEvent {
    protected _type : string

    public constructor(type : string, source = '', target = '', content = '') {
        super(source, target, content)
        this._type = type
    }

    public toJSON() : string {
        return JSON.stringify({type: this.type, source: this.source, content : this.content})
    }

    get type(): string {
        return this._type;
    }

    set type(value: string) {
        this._type = value;
    }
}