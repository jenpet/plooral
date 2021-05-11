import Konva from "konva";
import {getBus, BusEvent as Event, BusEvent} from "@/bus";
import {StickyNote, StickyNoteConfig} from "./shapes/stickyNote"
import {ShapeDefintion} from "@/composables/canvas/shapes";

const bus = getBus('canvas')

const canvasWidth = 2480
const canvasHeight = 1754

export let stage : Konva.Stage
export let defaultLayer : Konva.Layer
let defaultTransformer : Konva.Transformer

function handleCanvasKeyPress(ev : KeyboardEvent) {
    // delete or backspace key
    if ((ev.code === 'Backspace' || ev.code === 'Delete')) {
        destroyTransformerNodes()
    } else if (ev.code === 'Enter') {
        if (defaultTransformer.nodes().length == 1) {
            bus.emit('canvas:shape:edit', defaultTransformer.nodes()[0])
        }
    }
}

function handleCanvasKeyUp(ev : KeyboardEvent) {
    // delete or backspace key
    if ((ev.code === 'Backspace' || ev.code === 'Delete')) {
        destroyTransformerNodes()
    }
}

function handleSelectShape(ev : BusEvent) {
    addTransformerNode(ev.content)
}

export function setup (el: string, width : number, height : number, content? : any) : void {
    setupStage(el, width, height, content)
    bus.on('canvas:shape:added', handleShapeAdded)
    bus.on('canvas:shape:updated', handleShapeUpdated)
    bus.on('canvas:shape:destroyed', handleShapeDestroyed)
    bus.on('canvas:transformable:clicked', handleTransformableClicked)
    bus.on('canvas:shape:select', handleSelectShape)
    bus.on('canvas:shape:editing', handleShapeEditing)
    bus.on('canvas:shape:lostfocus', () => stage.content.focus())
    bus.on('sidebar:action-menu:export', downloadImage)
    document.addEventListener('keypress', handleCanvasKeyPress)
    document.addEventListener('keyup', handleCanvasKeyUp)
}

function setupStage(el: string, width : number, height : number, content? : any) : void {
    defaultTransformer = new Konva.Transformer({
        padding: 2,
        anchorCornerRadius: 1,
        enabledAnchors: ['top-left', 'top-right', 'bottom-left', 'bottom-right'],
    })
    defaultLayer = new Konva.Layer()
    defaultLayer.setName("default")
    defaultLayer.add(defaultTransformer)
    if (!content) {
        stage = new Konva.Stage({
            container: el,
            width: canvasWidth,
            height: canvasHeight
        })
        stage.add(defaultLayer)
    } else {
        stage = new Konva.Stage({
            container: el,
            width: content.attrs.width,
            height: content.attrs.height
        })
        stage.add(defaultLayer)
        if (content.children.length == 1) {
            const inputDefaultLayer = content.children[0]
            inputDefaultLayer.children.forEach((c:Konva.Node) => {
                const node = nodeFromJSON(JSON.stringify(c))
                if (typeof node !== 'undefined') {
                    defaultLayer.add(node)
                }
            })
            defaultLayer.draw()
        }
    }
    // required to bind to keyboard events decently
    stage.container().tabIndex = 1

    stage.on('click tap', (e) => {
        if (e.target === stage) {
            deselectTransformerNodes()
            bus.emit('canvas:stage:clicked')
            return
        }
    })

    // drop actions
    stage.container().ondrop = (ev: DragEvent) => {
        ev.preventDefault()
        stage.setPointersPositions(ev)
        const data = ev.dataTransfer
        if (data) {
            const element = <ShapeDefintion>JSON.parse(data.getData("text/plain"))
            const stagePos = stage.getPointerPosition()
            element.x = stagePos!.x - (element.x ? element.x : 0)
            element.y = stagePos!.y - (element.y ? element.y : 0)
            addShape(element)
        }
    }

    stage.container().ondragover = (ev: DragEvent) => {
        ev.preventDefault()
    }
}

function addShape(sd : ShapeDefintion) {
    let shape : Konva.Shape | Konva.Group
    switch (sd.type) {
        case 'stickyNote': {
            const cfg = {
                x: sd.x,
                y: sd.y,
                width: sd.width,
                height: sd.height,
                color: sd.style.backgroundColor,
                text: 'Text'
            }
            shape = new StickyNote(cfg)
            break
        }
        default: {
            console.error('Could not find correct shape type. No shape added.')
            return
        }
    }
    defaultLayer.add(shape)
    defaultLayer.draw()
    bus.emit("canvas:shape:added", shape)
}

// TODO: Maybe handle shape added and updated within the same function via `upserted`
// TODO: Ensure that all events which are received via the API are already parsed as JSON and typed
function handleShapeAdded(event : Event) {
    if (!bus.isOriginTo(event)) {
        const node = nodeFromJSON(event.content)
        defaultLayer.add(node!)
        defaultLayer.draw()
    } else {
        // shape / group was added by the user and should be immediately added / highlighted
        addTransformerNode(event.content)
    }
}

function handleShapeUpdated(event : Event) {
    if (!bus.isOriginTo(event)) {
        const node = nodeFromJSON(event.content)
        // destroy previous node
        destroyNode(node!.attrs.id)
        defaultLayer.add(node!)
        defaultLayer.draw()
    }
}

function handleShapeDestroyed(event : Event) {
    if (!bus.isOriginTo(event)) {
        const shapeData = JSON.parse(event.content)
        destroyNode(shapeData.attrs.id)
    }
}

function handleTransformableClicked(event : Event) {
    const node = event.content.node
    const multiSelect = event.content.multiSelect
    if (!node.hasName('transformable')) {
        return;
    }

    // do we pressed shift or ctrl?
    const isSelected = defaultTransformer.nodes().indexOf(node) >= 0;

    if (!multiSelect && !isSelected) {
        // if no key pressed and the node is not selected
        // select just one
        defaultTransformer.nodes([node]);
        defaultTransformer.moveToTop()
    } else if (multiSelect && isSelected) {
        // if we pressed keys and node was selected
        // we need to remove it from selection:
        const nodes = defaultTransformer.nodes().slice(); // use slice to have new copy of array
        // remove node from array
        nodes.splice(nodes.indexOf(node), 1);
        defaultTransformer.nodes(nodes);
    } else if (multiSelect && !isSelected) {
        // add the node into selection
        const nodes = defaultTransformer.nodes().concat([node]);
        defaultTransformer.nodes(nodes);
        defaultTransformer.moveToTop()
    }
    defaultLayer.draw();
}

function destroyTransformerNodes() : void {
    defaultTransformer.nodes().forEach(shape => {
        shape.destroy()
        bus.emit('canvas:shape:destroyed', shape.toJSON())
    })
    defaultTransformer.nodes([])
    defaultLayer.draw()
}

function deselectTransformerNodes() {
    defaultTransformer.nodes([])
    defaultLayer.draw()
}

function deselectTransformerNode(nodeId : string) {
    const trIdx = defaultTransformer.nodes().findIndex(value => value.attrs.id === nodeId)
    if (trIdx >= 0) {
        const nodes = defaultTransformer.nodes().slice(); // use slice to have new copy of array
        // remove node from array
        nodes.splice(trIdx, 1);
        defaultTransformer.nodes(nodes);
    }
}

function addTransformerNode(node: Konva.Group | Konva.Shape) : void {
    defaultTransformer.nodes([node])
    defaultLayer.draw()
}

function destroyNode(nodeId : string) {
    deselectTransformerNode(nodeId)
    defaultLayer.findOne(`#${nodeId}`).destroy()
    defaultLayer.draw()
}

function handleShapeEditing(event : Event) {
    deselectTransformerNode(event.content)
}

function downloadImage() {
    const dataURL = stage.toDataURL({ pixelRatio: 3 });
    const link = document.createElement('a');
    link.download = "plooral-export.png";
    link.href = dataURL;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}

export function stageJSON() : string {
    return stage.toJSON()
}

function nodeFromJSON(json : string) : Konva.Group | undefined {
    const node = Konva.Node.create(JSON.parse(json))
    if (node.hasName('stickyNote')) {
        const cfg = new StickyNoteConfig(node)
        return new StickyNote(cfg)
    } else {
        return undefined
    }
}
