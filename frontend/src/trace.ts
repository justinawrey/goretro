function trace(label: string, fn: () => void): void {
    const start = new Date().getTime()
    fn()
    const elapsed = new Date().getTime() - start
    console.log(`${label}: ${elapsed}ms`)
}

export default trace
