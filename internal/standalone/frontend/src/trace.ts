function trace(fn: () => void): void {
    const start = new Date().getTime()
    fn()
    const elapsed = new Date().getTime() - start
    console.log('elapsed time:', elapsed)
}

export default trace
