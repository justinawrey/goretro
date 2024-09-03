import { RequestFrame } from '../wailsjs/go/main/App'
import trace from './trace'
import draw from './display'

// debug!
// eslint-disable-next-line @typescript-eslint/no-explicit-any
;(window as any).requestFrame = () =>
    trace('roundtrip js -> go -> js asking for data', () => {
        RequestFrame().then((data) => draw(data))
    })
