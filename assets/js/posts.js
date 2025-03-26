const linkPostPosts = document.querySelectorAll('[id^="link-post-post"]');
const linkPostAuthors = document.querySelectorAll('[id^="link-post-author"]');
const linkPostChannels = document.querySelectorAll('[id^="link-post-channel"]');

// TODO: change these to forEach

document.addEventListener("DOMContentLoaded", () => {
  linkPostPosts.forEach((post) => {
    post.addEventListener("click", () => {
      const postId = post.dataset.postid;
      const postUrl = `/posts/${postId}`;
      window.location.href = postUrl;
    });
  });

  linkPostAuthors.forEach((author) => {
    author.addEventListener("click", () => {
      const authorId = author.dataset.authorid;
      const authorUrl = `/users/${authorId}`;
      window.location.href = authorUrl;
    });
  });

  linkPostChannels.forEach((channel) => {
    channel.addEventListener("click", () => {
      const channelId = channel.dataset.channelid;
      const channelUrl = `/channels/${channelId}`;
      window.location.href = channelUrl;
    });
  });
});
