import { navigateToPage } from "./fetch_and_navigate.js";
import {
  showMainNotification,
  showInlineNotification,
} from "./notifications.js";
import {
  applyCustomTheme,
  showThemesClickable,
  pickTheme,
  themeList,
} from "./consoleThemes.js";

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

const notifierPost = document.querySelector("#post-popover-title");
const createPostForm = document.querySelector("#form-create-post");

if (createPostForm) {
  const btnSelectChannels = createPostForm.querySelector(".dropdown-toggle");
  const inputTitle = createPostForm.querySelector('input[name="title"]');
  const inputContent = createPostForm.querySelector('textarea[name="content"]');

  createPostForm.addEventListener("submit", function (e) {
    const formData = new FormData(createPostForm);
    const channels = formData.getAll("post_channel_list");
    const title = formData.get("title")?.trim();
    const content = formData.get("content")?.trim();

    // Validate channels
    if (channels.length < 2) {
      e.preventDefault();
      showInlineNotification(
        notifierPost,
        "",
        "select at least 2 channels",
        false,
        "dummy",
      );
      btnSelectChannels.classList.add("active");
      return;
    }

    // Validate title
    if (!title) {
      e.preventDefault();
      showInlineNotification(notifierPost, "", "enter a title", false, "dummy");
      inputTitle.focus();
      return;
    }

    // Validate content
    if (!content) {
      e.preventDefault();
      showInlineNotification(
        notifierPost,
        "",
        "enter some content",
        false,
        "dummy",
      );
      inputContent.focus();
      return;
    }

    // ✅ If we reach here, validation passed
    // Allow the form to submit normally → server handles it → browser redirects
  });
}

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
