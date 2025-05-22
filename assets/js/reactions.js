export function listenToLikeDislike() {
  const postControls = document.querySelectorAll(".post-controls");
  const sidebarProfile = document.querySelector(".sidebarProfile");

  if (!sidebarProfile) return;

  const currentUserID = document
    .querySelector(".sidebar h3")
    ?.getAttribute("data-current-user-ID");

  if (!currentUserID) return;

  postControls.forEach((singlePostControl) => {
    const likeCountElement = singlePostControl.querySelector(".btn-likes");
    const dislikeCountElement =
      singlePostControl.querySelector(".btn-dislikes");

    if (!likeCountElement || !dislikeCountElement) return;

    const likeButton = likeCountElement.closest("button");
    const dislikeButton = dislikeCountElement.closest("button");

    const postCard = singlePostControl.closest(".card");
    if (!likeButton || !dislikeButton || !postCard) return;

    const postID = postCard.getAttribute("data-post-id");
    const commentID = postCard.getAttribute("data-comment-id");
    const channelID = postCard.getAttribute("data-channel-id");

    likeButton.addEventListener("click", () => {
      handleReaction({
        type: "like",
        likeButton,
        dislikeButton,
        likeCountElement,
        dislikeCountElement,
        postID,
        commentID,
        channelID,
        currentUserID,
      });
    });

    dislikeButton.addEventListener("click", () => {
      handleReaction({
        type: "dislike",
        likeButton,
        dislikeButton,
        likeCountElement,
        dislikeCountElement,
        postID,
        commentID,
        channelID,
        currentUserID,
      });
    });
  });
}

function handleReaction({
  type,
  likeButton,
  dislikeButton,
  likeCountElement,
  dislikeCountElement,
  postID,
  commentID,
  channelID,
  currentUserID,
}) {
  console.clear();

  let likeCount = parseInt(likeCountElement.textContent, 10);
  let dislikeCount = parseInt(dislikeCountElement.textContent, 10);

  const isLike = type === "like";
  const mainButton = isLike ? likeButton : dislikeButton;
  const altButton = isLike ? dislikeButton : likeButton;
  const mainCountElement = isLike ? likeCountElement : dislikeCountElement;
  const altCountElement = isLike ? dislikeCountElement : likeCountElement;

  let mainCount = isLike ? likeCount : dislikeCount;
  let altCount = isLike ? dislikeCount : likeCount;

  const isActive = mainButton.classList.contains("active");

  if (isActive) {
    mainCountElement.textContent = `${mainCount - 1}`;
    mainButton.classList.remove("active");
  } else {
    mainCountElement.textContent = `${mainCount + 1}`;
    mainButton.classList.add("active");

    if (altButton.classList.contains("active")) {
      altCountElement.textContent = `${altCount - 1}`;
      altButton.classList.remove("active");
    }
  }

  const postData = checkData(commentID, postID, currentUserID, channelID, type);

  console.table(postData);
  fetchData(postData, type);
}

export function listenToNoUser() {
  const btns = document.querySelectorAll(".btn-action.nouser");
  const timeOut = 3000; // 3 seconds

  btns.forEach((btn) => {
    btn.addEventListener("click", (e) => {
      e.preventDefault();

      // Remove tooltip from all others
      document
        .querySelectorAll(".btn-action.nouser.show-tooltip")
        .forEach((el) => el.classList.remove("show-tooltip"));

      // Add tooltip to this one
      btn.classList.add("show-tooltip");

      setTimeout(() => {
        btn.classList.remove("show-tooltip");
      }, timeOut);
    });
  });

  // Optional: click anywhere else to hide tooltips
  document.addEventListener("click", (e) => {
    if (!e.target.closest(".btn-action.nouser")) {
      document
        .querySelectorAll(".btn-action.nouser.show-tooltip")
        .forEach((el) => el.classList.remove("show-tooltip"));
    }
  });
}

function checkData(commentID, postID, reactionAuthorID, channelID, likeStatus) {
  const isLike = likeStatus === "like";
  const isDislike = likeStatus === "dislike";

  if (!isLike && !isDislike) {
    console.error("Invalid likeStatus:", likeStatus);
    return null;
  }

  if (!reactionAuthorID) {
    console.error("Missing or invalid reactionAuthorID:", reactionAuthorID);
    return null;
  }

  const postData = {
    liked: isLike,
    disliked: isDislike,
    authorId: reactionAuthorID,
    ...(postID ? { reactedPostId: Number(postID) } : {}),
    ...(commentID ? { reactedCommentId: Number(commentID) } : {}),
  };

  if (postID) {
    console.info(`${isLike ? "liked" : "disliked"} post: ${postID}`);
  } else if (commentID) {
    console.info(`${isLike ? "liked" : "disliked"} comment: ${commentID}`);
  }
  return postData;
}

function fetchData(postData, likeString) {
  fetch("/store-reaction", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(postData),
  })
    .then(async (response) => {
      const text = await response.text();

      if (!response.ok) {
        throw new Error(
          `HTTP Error: ${response.status} ${response.statusText}. Response: ${text}`,
        );
      }

      // Try to parse JSON if there's any content
      let data = null;
      if (text) {
        try {
          data = JSON.parse(text);
        } catch (err) {
          console.warn("Response is not valid JSON:", text);
          data = text;
        }
      }

      return data;
    })
    .then((data) => {
      console.log(`${likeString} updated (in js):`, data ?? "(no response)");
    })
    .catch((error) => {
      console.error("Error updating like:", error.message);
    });
}
