const REGEX_THREAD_ID = /t=(\d+)/;

(async function () {
  let threads = await browser.storage.local.get()
  threads = threads || {}

  const rows = document.querySelectorAll('#threads tbody tr:not(.sticky)')
  rows.forEach(row => {
    const link = row.querySelector('a.title')
    const threadID = link.href.match(REGEX_THREAD_ID)[1]
    let thread = threads[threadID]
    if (typeof thread !== 'object') {
      thread = {}
    }
    thread.title = link.innerText;
    threads[threadID] = thread;
  })

  await browser.storage.local.set(threads)
})();