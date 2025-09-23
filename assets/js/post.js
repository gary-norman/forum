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
  createPostForm.addEventListener("submit", async function (e) {
    e.preventDefault();
    const form = e.target;
    const formData = new FormData(form); // Collect form data
    const channels = formData.getAll("post_channel_list"); // get an array of selected channels
    const title = formData.get("title")?.trim();
    const content = formData.get("content")?.trim();

    if (channels.length < 2) {
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
    if (!title) {
      showInlineNotification(notifierPost, "", "enter a title", false, "dummy");
      inputTitle.focus();
      return;
    }
    if (!content) {
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

    try {
      const response = await fetch("/posts/create", {
        method: "POST",
        body: formData,
        cache: "no-store",
      });

      const data = await response.json();

      if (!response.ok || data.message === "post creation failed!") {
        showInlineNotification(notifierPost, "", data.message, false, "dummy");
        return;
      }
      showInlineNotification(notifierPost, "", data.message, true, "dummy");
      setTimeout(
        function () {
          popoverTargetElement = this;
          popoverTargetAction = "hide";
        }.bind(this),
        2000,
      );
    } catch (error) {
      console.custom.error("Error during post creation:", error);
      showMainNotification("An error occurred during post creation.");
    }
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
