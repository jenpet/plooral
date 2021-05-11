import {ShapeDefintion} from "@/composables/canvas/shapes/definition";
import colorPalette from "@/composables/canvas/shapes/colors";

const shapes : ShapeDefintion[] = [
    {
        type: 'stickyNote',
        label: 'Sticky Note Square',
        class: 'sticky-note-square',
        style: {
            width: '6rem',
            height: '6rem'
        },
        width: 96,
        height: 96,
    },
    {
        type: 'stickyNote',
        label: 'Sticky Note Rectangle',
        class: 'sticky-note-rectangle',
        style: {
            width: '12rem',
            height: '6rem'
        },
        width: 192,
        height: 96,
    }
]

export function list() : ShapeDefintion[][] {
    const list : ShapeDefintion[][] = []
    shapes.forEach(s => {
        const sType: ShapeDefintion[] = []
        colorPalette.forEach(c => {
            // deep clone for styles
            const shape = JSON.parse(JSON.stringify(s))
            shape.style.backgroundColor = c.background
            shape.style.borderColor = c.border
            sType.push(shape)
        })
        list.push(sType)
    })
    return list
}

