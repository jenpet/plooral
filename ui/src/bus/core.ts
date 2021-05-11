// TODO: Maybe it would be a nice abstraction when the API interaction hooks in smoothly with the overall emitter in the frontend.
// An adapter could publish all events to the API and take all events from the API and publish them on the bus
// source and target might be used to distinguish whether an event was triggered by a component on its own
import mitt, {Emitter, EventType} from 'mitt'
import {Handler, WildcardHandler} from "mitt/src/index";

import { v4 as uuid } from "uuid";

// central instance of the bus
const emitter : Emitter = mitt()

// unique identifier for this client
const clientId : string = uuid()

export function getBus(componentId : string) : ComponentBus {
    return new ComponentBus(emitter, componentId)
}

export class ComponentBus {
    emitter : Emitter
    componentId : string
    private readonly _clientId : string
    constructor(emitter : Emitter, componentId : string) {
        this.emitter = emitter
        this.componentId = componentId
        this._clientId = clientId
    }

    public emit<T>(type: EventType, content? : T) : void {
        const event = new BusEvent(this.getSourceId())
        event.content = content
        this.emitter.emit(type, event)
    }

    // simply forwards an event to the bus without altering it

    public forward(type: EventType, event: BusEvent) : void {
        this.emitter.emit(type, event)
    }
    public onAll(handler : WildcardHandler) {
        this.emitter.on('*', (t,e) => {
            handler(t,e)
        })
    }

    public on(type: EventType, handler: Handler): void {
        this.emitter.on(type, e => {
            handler(e)
        })
    }

    public isOriginTo(e: BusEvent) : boolean {
        return typeof e.source !== 'undefined' && this.getSourceId() === e.source
    }

    public isOriginClientTo(e: BusEvent) : boolean {
        return typeof e.source !== 'undefined' && this._clientId === e.source.split(':')[0]
    }

    private getSourceId() {
        return `${this._clientId}:${this.componentId}`
    }

    get clientId(): string {
        return this._clientId;
    }
}

// dedicated mixins from the mitt.Emitter itself
Object.assign(ComponentBus.prototype, emitter.off)
Object.assign(ComponentBus.prototype, emitter.all)

export class BusEvent {

    protected _source: string
    protected _content : any
    protected _timestamp : number

    public constructor(source : string, content? : any, target? : string) {
        this._source = source
        this._content = content
        this._timestamp = new Date().getTime()
    }

    get timestamp(): number {
        return this._timestamp;
    }

    set timestamp(value: number) {
        this._timestamp = value;
    }
    get content(): any {
        return this._content;
    }

    set content(value: any) {
        this._content = value;
    }

    get source(): string {
        return this._source;
    }

    set source(value: string) {
        this._source = value;
    }
}

// debugging
getBus("debuglog").onAll((type, event) => console.log(`event type '${String(type)}' on bus from ${(<BusEvent> event).source}`))