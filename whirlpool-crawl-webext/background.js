const visit = async function (id) {
  console.log(`Visiting thread ${id}`)
  const tab = await browser.tabs.create({
    url: `https://forums.whirlpool.net.au/archive/${id}`,
    active: false
  })
  const result = await browser.tabs.executeScript(tab.id, {
    file: 'collect-posts.js',
  })
  await browser.tabs.remove(tab.id)
  return result
};

browser.contextMenus.create({
  id: 'whirlpool-crawl-pages',
  title: 'Whirlpool: Crawl Indexed Pages',
  contexts: ['all']
})

browser.contextMenus.onClicked.addListener((info, tab) => {
  (async function () {
    try {
      const threads = await browser.storage.local.get()

      for (const id of Object.keys(threads)) {
        // Skip ones we already completed
        if (threads[id].posts) {
          continue
        }

        const thread = threads[id]
        thread.posts = await visit(id)
        threads[id] = thread
        await browser.storage.local.set(threads)
      }  
    } catch (e) {
      console.log(e)
    }
  })();
})