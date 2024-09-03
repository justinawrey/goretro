import { EventsEmit, EventsOn } from '../wailsjs/runtime/runtime'

const RECEIVE_INPUT_A = 'receive-input-a'
const RECEIVE_INPUT_B = 'receive-input-b'
const RECEIVE_INPUT_SELECT = 'receive-input-select'
const RECEIVE_INPUT_START = 'receive-input-start'
const RECEIVE_INPUT_UP = 'receive-input-up'
const RECEIVE_INPUT_RIGHT = 'receive-input-right'
const RECEIVE_INPUT_DOWN = 'receive-input-down'
const RECEIVE_INPUT_LEFT = 'receive-input-left'

const REQUEST_INPUT_A = 'request-input-a'
const REQUEST_INPUT_B = 'request-input-b'
const REQUEST_INPUT_SELECT = 'request-input-select'
const REQUEST_INPUT_START = 'request-input-start'
const REQUEST_INPUT_UP = 'request-input-up'
const REQUEST_INPUT_RIGHT = 'request-input-right'
const REQUEST_INPUT_DOWN = 'request-input-down'
const REQUEST_INPUT_LEFT = 'request-input-left'

EventsOn(REQUEST_INPUT_A, () => EventsEmit(RECEIVE_INPUT_A, 'a'))
EventsOn(REQUEST_INPUT_B, () => EventsEmit(RECEIVE_INPUT_B, 'b'))
EventsOn(REQUEST_INPUT_SELECT, () => EventsEmit(RECEIVE_INPUT_SELECT, 'select'))
EventsOn(REQUEST_INPUT_START, () => EventsEmit(RECEIVE_INPUT_START, 'start'))
EventsOn(REQUEST_INPUT_UP, () => EventsEmit(RECEIVE_INPUT_UP, 'up'))
EventsOn(REQUEST_INPUT_RIGHT, () => EventsEmit(RECEIVE_INPUT_RIGHT, 'right'))
EventsOn(REQUEST_INPUT_DOWN, () => EventsEmit(RECEIVE_INPUT_DOWN, 'down'))
EventsOn(REQUEST_INPUT_LEFT, () => EventsEmit(RECEIVE_INPUT_LEFT, 'left'))
