import {
    SetButtonA,
    SetButtonB,
    SetButtonDown,
    SetButtonLeft,
    SetButtonRight,
    SetButtonSelect,
    SetButtonStart,
    SetButtonUp,
} from '../wailsjs/go/main/App'

// TODO: make sure these key mappings actually work
function handleKeypress(key: string, to: boolean): void {
    switch (key) {
        case 'w':
            SetButtonUp(to)
            return
        case 'a':
            SetButtonLeft(to)
            return
        case 's':
            SetButtonDown(to)
            return
        case 'd':
            SetButtonRight(to)
            return
        case 'j':
            SetButtonA(to)
            return
        case 'k':
            SetButtonB(to)
            return
        case '1':
            SetButtonStart(to)
            return
        case '2':
            SetButtonSelect(to)
            return
    }
}

window.addEventListener('keydown', (e) => handleKeypress(e.key, true))
window.addEventListener('keyup', (e) => handleKeypress(e.key, false))
