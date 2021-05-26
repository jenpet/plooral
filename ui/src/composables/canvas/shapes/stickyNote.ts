import Konva from "konva";
import {BusEvent, getBus} from "@/bus";
import { v4 as uuid } from "uuid";
import { ref, Ref, UnwrapRef, watch } from 'vue'

const rectFontSize = 18
const rectFontFamily = 'Patrick-Hand'

const shapeType = 'stickyNote'
const bus = getBus('canvas')

export class StickyNote extends Konva.Group {
    base: Konva.Rect
    text: Konva.Text
    inputContainer: HTMLElement | undefined
    input: HTMLElement | undefined
    editMode: Ref<UnwrapRef<boolean>> = ref(false)

    constructor(cfg : StickyNoteConfig) {
        super(toGroupConfig(cfg));
        this.base = new Konva.Rect(toRectConfig(cfg))
        this.text = new Konva.Text(toTextConfig(cfg))
        this.setupEditElement()
        this.add(this.base)
        this.add(this.text)
        this.setupListeners()
        // if the id is set via the config (i.e. is a recreation via an update) use the given id, otherwise create it
        this.id((typeof cfg.id !== 'undefined' ? cfg.id : uuid()))
        if (typeof cfg.editMode !== 'undefined') {
            this.editMode.value = cfg.editMode
        }
    }

    private showEditMode() {
        // at first find the text's position relative to the stage
        // TODO: Maybe this could be simplified via events in case the stage container changes its position or gets resized
        const stage = this.getStage()
        if (stage === null) {
            console.error('Stage of sticky note is null. Cannot enable show edit mode.')
            return
        }
        const textPosition = this.text.getAbsolutePosition()
        const areaPosition = {
            x: stage.container().offsetLeft + textPosition.x,
            y: stage.container().offsetTop + textPosition.y,
        };

        // hide the text layer first
        this.text.hide()
        this.input!.innerText = this.text.text()
        this.inputContainer!.style.display = 'flex'
        this.inputContainer!.style.top = areaPosition.y + 'px'
        this.inputContainer!.style.left = areaPosition.x + 'px'
        this.inputContainer!.style.height = this.text.height()*this.scaleY() - this.text.padding() * 2 +'px'
        this.inputContainer!.style.width = this.text.width()*this.scaleX() - this.text.padding() * 2 + 'px';
        this.inputContainer!.style.fontSize = this.text.fontSize()*this.scaleX() + 'px'
        this.inputContainer!.style.lineHeight = String(this.text.lineHeight())
        this.inputContainer!.style.fontFamily = this.text.fontFamily()
        this.inputContainer!.style.textAlign = this.text.align()
        this.inputContainer!.style.transformOrigin = '0 0'
        // add rotation to the input in case the shape was transformed / rotated
        const rotation = this.rotation()
        if (rotation) {
            this.inputContainer!.style.transform = 'rotateZ(' + rotation + 'deg)'
        }
        this.input!.focus()
        bus.emit('canvas:shape:editing', this.id())
    }

    private hideEditMode() {
        if (this.input !== undefined && String(this.text.text) !== this.input.innerText) {
            this.text.text(this.input.innerText.trim())
            bus.emit('canvas:shape:updated', this.toJSON())
        }
        this.text.show()
        this.inputContainer!.style.display = 'none'
        bus.emit('canvas:shape:lostfocus', this.toJSON())
        bus.emit('canvas:shape:select', this)
    }

    private setupListeners() {
        // Konva events that should be listened to
        // the text is to upper most layer, if clicked or tapped emit that the whole group was clicked
        this.text.on('mousedown', e => {
            const multiSelect = e.evt.shiftKey || e.evt.metaKey;
            // TODO: Maybe work with generics when emitting events (e.g. BusEvent<T>; BusEvent<TransformableActive>)
            bus.emit('canvas:transformable:clicked', {node:this, multiSelect: multiSelect})
        })

        // enable edit mode when double clicked
        this.text.on('dblclick dbltap', () => this.editMode.value = true)

        // if the canvas stage was clicked hide the edit mode
        bus.on('canvas:stage:clicked', () => this.editMode.value = false)

        // emit an event that the shape was updated after dragging or transforming
        this.on('dragend transformend', ev => {
            bus.emit("canvas:shape:updated", ev.target.toJSON())
        })

        this.on('mouseenter', (() => {
            const stage = this.getStage()
            if (stage !== null) {
                stage.container().style.cursor = 'pointer'
            }
        }).bind(this))

        this.on('mouseleave', (() => {
            const stage = this.getStage()
            if (stage !== null) {
                stage.container().style.cursor = 'default'
            }
        }))

        bus.on('canvas:shape:edit', ((ev:BusEvent) => {
            if (ev.content.id() === this.id()) {
                this.editMode.value = true
            }
        }).bind(this))

        // watch changes on the edit mode property and show or hide the inputs
        watch(this.editMode, (newValue, oldValue) => {
            if (newValue && !oldValue) {
                this.showEditMode()
            }
            if (!newValue && oldValue) {
                this.hideEditMode()
            }
        })
    }

    private setupEditElement() {
        const canvas = document.getElementById('canvas')
        if (canvas === null) {
            console.error('could not get canvas or stage')
            return
        }
        this.inputContainer = document.createElement('div')
        this.inputContainer.setAttribute('class', 'centered-text-edit-container')
        this.input = document.createElement('div')
        this.input.setAttribute('class', 'text')
        this.input.setAttribute('contenteditable', 'true')
        this.inputContainer.appendChild(this.input)
        canvas.appendChild(this.inputContainer)
    }
}

export class StickyNoteConfig {
    id?: string
    x?: number
    y?: number
    width?: number
    height?: number
    color?: string
    text?: string
    scaleX?: number
    scaleY?: number
    skewX?: number
    rotation?: number
    editMode? = false

    constructor(node : Konva.Node) {
        if (node.hasName('stickyNote')) {
            this.id = node.id()
            this.x = node.x()
            this.y = node.y()
            this.width = node.width()
            this.height = node.height()
            this.scaleX = node.scaleX()
            this.scaleY = node.scaleY()
            this.skewX = node.skewX()
            this.rotation = node.rotation()
            const bg = node.children.toArray().find(
                n => n.hasName(shapeType + '.bg')
            )
            if (typeof bg != 'undefined') {
                this.color = bg.attrs.fill
            }
            const text = node.children.toArray().find(
                n => n.hasName(shapeType + '.text')
            )
            if (typeof bg != 'undefined') {
                this.text = text.attrs.text
            }
        }
    }
}


function toRectConfig(cfg? : StickyNoteConfig) : Konva.RectConfig {
    if (typeof cfg === 'undefined') {
        return {}
    }
    return {
        name: 'stickyNote.bg',
        width: cfg.width,
        height: cfg.height,
        listening: false,
        fill: cfg.color
    }
}

function toGroupConfig(cfg? : StickyNoteConfig) : Konva.ContainerConfig {
    if (typeof cfg === 'undefined') {
        return {}
    }
    return {
        name: 'stickyNote transformable',
        x: cfg.x,
        y: cfg.y,
        width: cfg.width,
        height: cfg.height,
        scaleX: cfg.scaleX,
        scaleY: cfg.scaleY,
        skewX: cfg.skewX,
        rotation: cfg.rotation,
        draggable: true
    }
}
function toTextConfig(cfg? : StickyNoteConfig) : Konva.TextConfig {
    if (typeof cfg === 'undefined') {
        return {}
    }
    return {
        name: 'stickyNote.text',
        text: cfg.text,
        fontSize: rectFontSize,
        fontFamily: rectFontFamily,
        fill: '#000',
        width: cfg.width,
        height: cfg.height,
        align: 'center',
        verticalAlign: 'middle'
    }
}