(async () => {
  const posts = document.querySelectorAll('.reply.reply-archived[data-replyid]')
  return Array.prototype.map.call(posts, (post) => {
    try {
      return {
        id: post.dataset.replyid,
        user: post.dataset.uname,
        text: post.querySelector('.replytext').innerText,
        date: post.querySelector('.date').innerText
      }
    } catch (e) {
      console.log(e)
    }
  })
})()