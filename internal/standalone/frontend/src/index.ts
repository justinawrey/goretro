/* @refresh reload */

// import { render } from "solid-js/web";

import './tailwind.css'
import trace from './trace'

const NUM_TEXELS_WIDTH = 256 * 4
const NUM_TEXELS_HEIGHT = 256 * 4
const REAL_TEXEL_WIDTH = 10 / 4
const REAL_TEXEL_HEIGHT = 10 / 4

const canvas = document.getElementById('render-target') as HTMLCanvasElement
const ctx = canvas.getContext('2d')

if (ctx == null) {
    console.error('canvas unsupported')
} else {
    draw(ctx)
}

function draw(ctx: CanvasRenderingContext2D): void {
    let toggle = true

    for (let i = 0; i < NUM_TEXELS_WIDTH; i++) {
        for (let j = 0; j < NUM_TEXELS_HEIGHT; j++) {
            ctx.fillStyle = toggle ? 'white' : 'black'
            toggle = !toggle

            ctx.fillRect(
                i * REAL_TEXEL_WIDTH,
                j * REAL_TEXEL_HEIGHT,
                REAL_TEXEL_WIDTH,
                REAL_TEXEL_HEIGHT,
            )
        }
    }
}

// debug!
// eslint-disable-next-line @typescript-eslint/no-explicit-any
;(window as any).forceDraw = () => {
    trace(() => draw(ctx as CanvasRenderingContext2D))
}

// render(() => <></>, document.getElementById("root") as HTMLElement);
