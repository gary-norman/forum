import { navigateToChannel, navigateToPost, navigateToAuthor } from "./main.js";

const linkPostPosts = document.querySelectorAll('[id^="link-post-post"]');
const linkPostAuthors = document.querySelectorAll('[id^="link-post-author"]');
const linkPostChannels = document.querySelectorAll('[id^="link-post-channel"]');

// INFO was a DOMContentLoaded function
export function listenToNavigationLinks()  {
    linkPostPosts.forEach((post) => {
      post.addEventListener("click", (e) => {
        e.preventDefault();
        console.log("clicked post");
        navigateToPost(post);
      });
    });

    linkPostAuthors.forEach((author) => {
      author.addEventListener("click", (e) => {
        e.preventDefault();
        console.log("clicked author");
        navigateToAuthor(author);
      });
    });

    linkPostChannels.forEach((channel) => {
      channel.addEventListener("click", (e) => {
        e.preventDefault();
        console.log("clicked channel");
        navigateToChannel(channel);
      });
    });
  }

