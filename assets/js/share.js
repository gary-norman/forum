const link = encodeURI(window.location.href);
export const data = {
  homePage: document.getElementById("home-page"),
  userPage: document.getElementById("user-page"),
  channelPage: document.getElementById("channel-page"),
  postPage: document.getElementById("post-page"),
};
export const pages = [
  data["homePage"],
  data["userPage"],
  data["channelPage"],
  data["postPage"],
];

pages.forEach((page) => {
  console.log("type: ", typeof page);
  console.log("page: ", page);
  console.log("classList: ", page.classList);
});

// TODO Add logic that positions the modal above the button if there's not enough space under
const commentMsg = encodeURIComponent(
  "Hey, I found this comment, you need to see it!",
);
const postMsg = encodeURIComponent(
  "Hey, I found this post. I think you may like it?",
);
const commentTitle = encodeURIComponent("Comment from User ??? Here");
const postTitle = encodeURIComponent("Post Title Here");
let activityButtons;
let scrollWindow;

console.log(data["homePage"]);

export function selectActiveFeed() {
  const activePage = pages.find((page) =>
    page.classList.contains("active-feed"),
  );

  switch (activePage) {
    case data["home-page"]:
      console.log("ON HOME PAGE BITCH")
      // scrollWindow = homePage.querySelector(`#home-feed`);
      // console.log("scrollWindow:", scrollWindow)
      break;
    case data["user-page"]:
      console.log("ON USER PAGE BITCH")
      const userFeeds = Array.from(
        document.querySelectorAll('[id^="activity-feed-"]'),
      );
      const activeFeed = userFeeds.find((feed) =>
        feed.classList.contains("collapsible-expanded"),
      );

      scrollWindow = activeFeed;
      // console.log(scrollWindow)
      break;
    case data["channel-page"]:
      console.log("ON CHANNEL PAGE BITCH")

      scrollWindow = data["channel-page"].querySelector(`#channel-feed`);
      // console.log(scrollWindow)
      break;
    default:
      console.log(`No active feed.`);
      break;
  }
}

// INFO was a DOMContentLoaded function
export function listenToShare() {
  let postID, commentID, channelID, msg, title;
  selectActiveFeed();
  const buttonControls = document.querySelectorAll('[class$="-controls"]');

  // console.log("buttonControls: ", buttonControls);

  buttonControls.forEach((singleControl) => {
    postID = singleControl.getAttribute("data-post-id");
    commentID = singleControl.getAttribute("data-comment-id");
    channelID = singleControl.getAttribute("data-channel-id");

    // console.log("buttonControls:", buttonControls);
    // console.log("PostID:", postID);
    // console.log("CommentID:", commentID);
    // console.log("ChannelID:", channelID);
    let shareModal;
    let shareButton;

    if (postID === null) {
      shareModal = singleControl.querySelector(
        `[id^="share-container-channel-${channelID}-comment-${commentID}"]`,
      );
      shareButton = singleControl.querySelector(
        `[id^="share-button-channel-${channelID}-comment-${commentID}"]`,
      );
      // console.log("comment ID ✔️")
    } else if (commentID === null) {
      shareModal = singleControl.querySelector(
        `[id^="share-container-channel-${channelID}-post-${postID}"]`,
      );
      shareButton = singleControl.querySelector(
        `[id^="share-button-channel-${channelID}-post-${postID}"]`,
      );
      // console.log("post ID ✔️")
    }

    // console.log(shareModal);
    // console.log(shareButton);

    //get all components needed
    if (shareModal != null) {
      let label = shareModal.querySelector("label");
      let icon = shareModal.querySelector("button");
      let input = shareModal.querySelector("input");
      // console.log("shareButton:", shareButton)
      // console.log("shareModal:", shareModal)

      // Listen for the 'toggle' event on the modal (native popover event)
      shareModal.addEventListener("toggle", () => {
        setTimeout(() => {
          toggleButtonActive(shareModal, shareButton);
        }, 100);
      });

      label.addEventListener("click", async () => {
        try {
          await updateCopyStatus(input, label, icon);
          console.log("Copy operation completed");
        } catch (err) {
          console.error("Error during copy operation:", err);
        }
      });

      icon.addEventListener("click", async () => {
        try {
          await updateCopyStatus(input, label, icon);
          console.log("Copy operation completed");
        } catch (err) {
          console.error("Error during copy operation:", err);
        }
      });
    }

    if (shareButton != null) {
      shareButton.addEventListener("click", (e) => {
        // setTimeout(() => {
        getModalPos(shareButton, shareModal);
        attachScrollListener(shareButton, shareModal);
        // }, 200);
      });
    }

    const activityBar = document.getElementById("activity-bar");
    if (activityBar) {
      activityButtons = activityBar.querySelectorAll("button");

      activityButtons.forEach((button) => {
        button.addEventListener("click", () => {
          //set timeout due to the animation timeout when switching between activity feeds in main.js
          setTimeout(() => {
            // console.log('Before selectActiveFeed:', scrollWindow);
            selectActiveFeed();
            // console.log('After selectActiveFeed:', scrollWindow);
            attachScrollListener(shareButton, shareModal);
          }, 500);
        });
      });
    }

    // TODO check api's of sites and fix title/message
    //if postID present, make msg = postMsg, if commentID present, make msg = commentMsg
    if (postID !== 0 && commentID == null) {
      msg = postMsg;
      title = postTitle;
    } else if (commentID !== 0 && postID === null) {
      msg = commentMsg;
      title = commentTitle;
    }
    // console.log("PostID:", postID)
    // console.log("CommentID:", commentID)
    // console.log("message:", msg)
    // console.log("title:", title)

    if (singleControl) {
      const fb = singleControl.querySelector(".facebook");
      fb.href = `https://www.facebook.com/share.php?u=${link}`;
      // TODO check https works
      const twitter = singleControl.querySelector(".twitter");
      twitter.href = `https://twitter.com/share?&url=${link}&text=${msg}&hashtags=javascript,programming`;
      // TODO check https works
      const linkedIn = singleControl.querySelector(".linkedin");
      linkedIn.href = `https://www.linkedin.com/sharing/share-offsite/?url=${link}`;
      // TODO check https works
      const reddit = singleControl.querySelector(".reddit");
      reddit.href = `https://www.reddit.com/submit?url=${link}&title=${title}`;
      // TODO check https works
      const whatsapp = singleControl.querySelector(".whatsapp");
      whatsapp.href = `https://api.whatsapp.com/send?text=${msg}: ${link}`;
      // TODO check https works
      const telegram = singleControl.querySelector(".telegram");
      telegram.href = `https://telegram.me/share/url?url=${link}&text=${msg}`;
      // TODO check https works
    }
  });
}

function attachScrollListener(shareButton, shareModal) {
  // console.log("attaching listener to:", scrollWindow)

  scrollWindow.addEventListener("scroll", (e) => {
    // console.log("scrolling on:", scrollWindow)
    scrollWindow.hasScrollListener = true; // Mark as attached
    getModalPos(shareButton, shareModal);

    // Parse the 'top' value and compare it with top of the mask / botton of screen
    const modalTop = parseInt(shareModal.style.top, 10); // Convert 'top' to a number

    if (modalTop <= 400 || modalTop >= window.innerHeight - 72) {
      // Close the popover
      shareModal.hidePopover();
    }
  });
}

function scrollToPost(postId) {
  selectActiveFeed();
  const container = scrollWindow;
  const post = scrollWindow.querySelector(`[data-post-id="${postId}"]`);

  if (post && container) {
    const containerRect = container.getBoundingClientRect();
    const postRect = post.getBoundingClientRect();

    // Calculate the position relative to the container
    const scrollOffset =
      postRect.top - containerRect.top + container.scrollTop - 40;
    // console.log("FIRST:");
    // console.log("containerRect:, ", containerRect);
    // console.log("postRect:, ", postRect);

    container.scrollTo({
      top: scrollOffset,
      behavior: "smooth", // Smooth scrolling animation
    });
    // console.log("THEN:");
    // console.log("containerRect:, ", containerRect);
    // console.log("postRect:, ", postRect);
    // console.log("container.scrollTop:, ", container.scrollTop);
    // console.log("scrollOffset", scrollOffset);

    post.classList.add("card-selected");

    setTimeout(() => {
      post.classList.remove("card-selected");
    }, 3000);
  } else {
    console.error("Post or container not found:", postId);
  }
}

const scrollButton = document.getElementById(`scroll-test`);
if (scrollButton) {
  scrollButton.addEventListener("click", () => {
    // Example: Scroll to post 3
    scrollToPost("3");
  });
}
const scrollButton1 = document.getElementById(`scroll-test-1`);
if (scrollButton1) {
  scrollButton1.addEventListener("click", () => {
    // Example: Scroll to post 3
    scrollToPost("3");
  });
}

function getModalPos(shareButton, shareModal) {
  //get button position
  const buttonPos = shareButton.getBoundingClientRect();

  //set modal styling
  shareModal.style.top = `${buttonPos.bottom - 16}px`;
  shareModal.style.left = `${buttonPos.left - 20}px`;
}

function toggleButtonActive(shareModal, shareButton) {
  if (shareModal.matches(":popover-open")) {
    shareButton.classList.add("active");
  } else {
    shareButton.classList.remove("active");
  }
}

function copyToClipboard(input) {
  input.focus();
  input.select();

  // Modern Clipboard API (recommended)
  if (navigator.clipboard) {
    return navigator.clipboard
      .writeText(input.value)
      .then(() => console.log("Text copied to clipboard!"))
      .catch((err) => {
        console.error("Failed to copy text:", err);
        throw err; // Ensure errors are propagated to the caller
      });
  } else {
    // Fallback for unsupported browsers
    console.warn("Clipboard API not supported, using execCommand as fallback");
    try {
      const success = document.execCommand("copy");
      if (!success) {
        throw new Error("Fallback copy failed");
      }
      console.log("Text copied to clipboard (fallback method)!");
      return Promise.resolve(); //
    } catch (err) {
      console.error("Fallback copy failed:", err);
      return Promise.reject(err); //
    }
  }
}

async function updateCopyStatus(input, label, icon) {
  const iconSpan = icon.querySelector(`span`);
  try {
    await copyToClipboard(input);

    // Update label text and styles
    label.textContent = "Copied!";
    label.style.color = "var(--color-hl-primary)";
    iconSpan.classList.remove("btn-copy");
    iconSpan.classList.add("btn-success");

    setTimeout(() => {
      label.textContent = "Copy URL link";
      label.style.color = "var(--color-fg-2)";
      iconSpan.classList.add("btn-copy");
      iconSpan.classList.remove("btn-success");
    }, 2000);
  } catch (err) {
    console.error("Failed to copy text:", err);
    label.textContent = "Failed to Copy";
    label.style.color = "var(--color-hl-red)";
    iconSpan.classList.remove("btn-copy");
    iconSpan.classList.add("btn-warning");

    setTimeout(() => {
      label.textContent = "Copy URL";
      label.style.color = "var(--color-fg-2)";
      iconSpan.classList.add("btn-copy");
      iconSpan.classList.remove("btn-warning");
    }, 2000);
  }
}
