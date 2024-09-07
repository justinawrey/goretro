import { SetButton } from '../wailsjs/go/app/WebviewInputDriver'
import { app } from '../wailsjs/go/models'

// TODO: customizable?
const keymap: Record<app.Button, string> = {
    [app.Button.UP]: 'w',
    [app.Button.RIGHT]: 'd',
    [app.Button.DOWN]: 's',
    [app.Button.LEFT]: 'a',
    [app.Button.A]: 'j',
    [app.Button.B]: 'k',
    [app.Button.START]: '1',
    [app.Button.SELECT]: '2',
}

function handleKeypress(joypad: app.Joypad, key: string, to: boolean): void {
    let button: keyof typeof keymap
    for (button in keymap) {
        if (keymap[button] === key) {
            SetButton(joypad, button, to)
            return
        }
    }
}

// TODO: secondary joypad keybindings?
// TODO: unsubscribe?
window.addEventListener('keydown', (e) =>
    handleKeypress(app.Joypad.PRIMARY, e.key, true),
)

window.addEventListener('keyup', (e) =>
    handleKeypress(app.Joypad.PRIMARY, e.key, false),
)
