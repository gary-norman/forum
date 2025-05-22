// INFO was a DOMContentLoaded function
export function listenToLikeDislike() {
  const postControls = document.querySelectorAll(".post-controls");
  const sidebar = document.querySelector(`.sidebar`);
  const sidebarProfile = document.querySelector(".sidebarProfile");
  let currentUserID;
  if (sidebarProfile) {
    currentUserID = sidebar
      .querySelector("h3")
      .getAttribute("data-current-user-ID");
  }

  postControls.forEach((singlePostControl) => {
    const likeCountElement = singlePostControl.querySelector(".btn-likes");
    const dislikeCountElement =
      singlePostControl.querySelector(".btn-dislikes");

    const likeButton = likeCountElement.closest("button");
    const dislikeButton = dislikeCountElement.closest("button");

    const likeID = likeCountElement.getAttribute("data-like-id");
    const dislikeID = dislikeCountElement.getAttribute("data-dislike-id");
    const postID = singlePostControl
      .closest(".card")
      .getAttribute("data-post-id");
    const commentID = singlePostControl
      .closest(".card")
      .getAttribute("data-comment-id");
    const channelID = singlePostControl
      .closest(".card")
      .getAttribute("data-channel-id");

    if (postControls.classList.contains("user")) {
      likeButton.addEventListener("click", function (event) {
        console.clear();
        let likeCount = parseInt(likeCountElement.textContent, 10);
        let dislikeCount = parseInt(dislikeCountElement.textContent, 10);

        if (likeButton.classList.contains("active")) {
          // Decrement the like count
          likeCountElement.textContent = `${likeCount - 1}`;
          // Remove the 'active' class to the like button
          likeButton.classList.remove("active");
        } else if (likeButton.classList.contains("inactive")) {
          // Increment the like count
          likeCountElement.textContent = `${likeCount + 1}`;
          // Add the 'active' class from the like button
          likeButton.classList.add("active");

          if (dislikeButton.classList.contains("active")) {
            // Decrement the like count
            dislikeCountElement.textContent = `${dislikeCount - 1}`;
            // Remove the 'active' class to the like button
            dislikeButton.classList.remove("active");
          }
        }
        let postData = checkData(
          commentID,
          postID,
          currentUserID,
          channelID,
          "like",
        );

        console.table(postData);

        fetchData(postData, "like");
      });
    }

    dislikeButton.addEventListener("click", function (event) {
      console.clear();
      let likeCount = parseInt(likeCountElement.textContent, 10);
      let dislikeCount = parseInt(dislikeCountElement.textContent, 10);

      if (dislikeButton.classList.contains("active")) {
        // Decrement the like count
        dislikeCountElement.textContent = `${dislikeCount - 1}`;
        // Remove the 'active' class to the like button
        dislikeButton.classList.remove("active");
      } else {
        // Increment the like count
        dislikeCountElement.textContent = `${dislikeCount + 1}`;
        // Add the 'active' class from the like button
        dislikeButton.classList.add("active");

        if (likeButton.classList.contains("active")) {
          // Decrement the like count
          likeCountElement.textContent = `${likeCount - 1}`;
          // Remove the 'active' class to the like button
          likeButton.classList.remove("active");
        }
      }

      let postData = checkData(
        commentID,
        postID,
        currentUserID,
        channelID,
        "dislike",
      );

      console.table(postData);

      fetchData(postData, "dislike");
    });
  });
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
