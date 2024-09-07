/* @refresh reload */
import './tailwind.css'
import './index.css'

import './display'
import './input'

if (import.meta.env.DEV) {
    import('./debug')
}
