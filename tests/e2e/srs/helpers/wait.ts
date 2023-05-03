export const wait = async (seconds: number) =>
  new Promise((r) => {
    setTimeout(() => r(true), 1000 * seconds)
  })