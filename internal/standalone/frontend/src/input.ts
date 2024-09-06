import { SetButton } from '../wailsjs/go/main/WebviewInputDriver'
import { main } from '../wailsjs/go/models'

// TODO: customizable?
const keymap: Record<main.Button, string> = {
    [main.Button.UP]: 'w',
    [main.Button.RIGHT]: 'd',
    [main.Button.DOWN]: 's',
    [main.Button.LEFT]: 'a',
    [main.Button.A]: 'j',
    [main.Button.B]: 'k',
    [main.Button.START]: '1',
    [main.Button.SELECT]: '2',
}

function handleKeypress(joypad: main.Joypad, key: string, to: boolean): void {
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
    handleKeypress(main.Joypad.PRIMARY, e.key, true),
)
window.addEventListener('keyup', (e) =>
    handleKeypress(main.Joypad.PRIMARY, e.key, false),
)
