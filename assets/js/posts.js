import {
  navigateToPage,
  navigateToChannel,
  navigateToPost,
  navigateToUser,
} from "./main.js";

const linkPosts = document.querySelectorAll('[id^="link-post-"]');
const linkPostPosts = document.querySelectorAll('[id^="link-post-post"]');
const linkPostAuthors = document.querySelectorAll('[id^="link-post-author"]');
const linkPostChannels = document.querySelectorAll('[id^="link-post-channel"]');

// INFO was a DOMContentLoaded function
export function listenToNavigationLinks() {
  linkPosts.forEach((link) => {
    link.addEventListener("click", (e) => {
      e.preventDefault();
      const fullID = e.target.id;
      const prefix = "link-post-";
      const entity = fullID.slice(prefix.length);
      console.log("clicked ", entity);
      navigateToPage(entity);
    });
  });

  // linkPostPosts.forEach((post) => {
  //   post.addEventListener("click", (e) => {
  //     e.preventDefault();
  //     console.log("clicked post");
  //     navigateToPost(post);
  //   });
  // });
  //
  // linkPostAuthors.forEach((author) => {
  //   author.addEventListener("click", (e) => {
  //     e.preventDefault();
  //     console.log("clicked author");
  //     navigateToUser(author);
  //   });
  // });
  //
  // linkPostChannels.forEach((channel) => {
  //   channel.addEventListener("click", (e) => {
  //     e.preventDefault();
  //     console.log("clicked channel");
  //     navigateToChannel(channel);
  //   });
  // });
}
