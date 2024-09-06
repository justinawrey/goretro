import { EventsOn } from '../wailsjs/runtime/runtime'
import { main } from '../wailsjs/go/models'

const NUM_TEXELS_WIDTH = 256
const NUM_TEXELS_HEIGHT = 256

const canvas = document.getElementById('render-target') as HTMLCanvasElement
const canvasCtx = canvas.getContext('2d', { alpha: false })

let imageData: ImageData
let ctx: CanvasRenderingContext2D

if (canvasCtx == null) {
    console.error('canvas unsupported')
} else {
    ctx = canvasCtx
    ctx.imageSmoothingEnabled = false

    // @ts-expect-error fuck
    ctx.mozImageSmoothingEnabled = false
    // @ts-expect-error fuck
    ctx.oImageSmoothingEnabled = false
    // @ts-expect-error fuck
    ctx.webkitImageSmoothingEnabled = false
    // @ts-expect-error fuck
    ctx.msImageSmoothingEnabled = false

    imageData = ctx.createImageData(NUM_TEXELS_WIDTH, NUM_TEXELS_HEIGHT)
}

EventsOn(main.RenderEvent.RENDER, draw)

function draw(data: number[]): void {
    for (let i = 0; i < imageData.data.length; i += 4) {
        imageData.data[i] = data[i]
        imageData.data[i + 1] = data[i + 1]
        imageData.data[i + 2] = data[i + 2]
        imageData.data[i + 3] = data[i + 3]
    }

    ctx.putImageData(imageData, 0, 0)
}
