const linkPostPosts = document.querySelectorAll('[id^="link-post-post"]');
const linkPostAuthors = document.querySelectorAll('[id^="link-post-author"]');
const linkPostChannels = document.querySelectorAll('[id^="link-post-channel"]');

document.addEventListener("DOMContentLoaded", () => {
  linkPostPosts.addEventListener("click", () => {
    const postId = linkPostPost.dataset.postid;
    const postUrl = `/posts/${postId}`;
    window.location.href = postUrl;
  });

  linkPostAuthors.addEventListener("click", () => {
    const authorId = linkPostAuthor.dataset.authorid;
    const authorUrl = `/users/${authorId}`;
    window.location.href = authorUrl;
  });

  linkPostChannels.addEventListener("click", () => {
    const channelId = linkPostChannel.dataset.channelid;
    const channelUrl = `/channels/${channelId}`;
    window.location.href = channelUrl;
  });
});
