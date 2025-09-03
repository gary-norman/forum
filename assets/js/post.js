import { navigateToPage } from "./fetch_and_navigate.js";

// const linkPosts = document.querySelectorAll('[id^="link-post-"]');
// const linkPostPosts = document.querySelectorAll('[id^="link-post-post"]');
// const linkPostAuthors = document.querySelectorAll('[id^="link-post-author"]');
// const linkPostChannels = document.querySelectorAll('[id^="link-post-channel"]');
//
// // INFO was a DOMContentLoaded function
// export function listenToNavigationLinks() {
//   linkPosts.forEach((link) => {
//     link.addEventListener("click", (e) => {
//       e.preventDefault();
//       const fullID = e.target.id;
//       const prefix = "link-post-";
//       const entity = fullID.slice(prefix.length);
//       console.log("clicked ", entity);
//       navigateToPage("page", entity);
//     });
//   });

// post channel selection dropdown

document.addEventListener("DOMContentLoaded", () => {
  listenToDropdowns();
});
// INFO was a DOMContentLoaded function
export function listenToDropdowns() {
  const dropdownToggle = document.querySelector(".dropdown-toggle");
  const wrapperDropdown = document.querySelector(".wrapper-dropdown");

  if (dropdownToggle) {
    dropdownToggle.addEventListener("click", () => {
      console.log("clicked dropdownToggle");
      const isActive = !dropdownToggle.classList.contains("active");
      dropdownToggle.classList.toggle("active");

      if (isActive) {
        // dropdownToggle.style.background = "var(--clr-accent--2)";
        wrapperDropdown.style.display = "block";
      } else {
        // dropdownToggle.style.background = "";
        wrapperDropdown.style.display = "none";
      }
    });
  }
}
