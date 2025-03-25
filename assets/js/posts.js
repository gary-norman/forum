const linkPostPost = document.querySelector("#link-post-post");
const linkPostAuthor = document.querySelector("#link-post-author");
const linkPostChannel = document.querySelector("#link-post-channel");

document.addEventListener("DOMContentLoaded", () => {
  linkPostPost.addEventListener("click", () => {
    const postId = linkPostPost.dataset.postId;
    const postUrl = `/posts/${postId}`;
    window.location.href = postUrl;
  });

  linkPostAuthor.addEventListener("click", () => {
    const authorId = linkPostAuthor.dataset.authorId;
    const authorUrl = `/users/${authorId}`;
    window.location.href = authorUrl;
  });

  linkPostChannel.addEventListener("click", () => {
    const channelId = linkPostChannel.dataset.channelId;
    const channelUrl = `/channels/${channelId}`;
    window.location.href = channelUrl;
  });
});
