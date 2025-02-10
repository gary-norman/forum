
document.addEventListener('DOMContentLoaded', function () {

    const postControls = document.querySelectorAll('.post-controls');
    const sidebar = document.querySelector(`.sidebar`);
    const reactionAuthorID = sidebar.querySelector('h3').getAttribute('data-current-user-ID');

    postControls.forEach(singlePostControl => {
        const likeCountElement = singlePostControl.querySelector('.btn-likes');
        const dislikeCountElement = singlePostControl.querySelector('.btn-dislikes');

        const likeButton = likeCountElement.closest('button');
        const dislikeButton = dislikeCountElement.closest('button');


        const likeID = likeCountElement.getAttribute("data-like-id");
        const dislikeID = dislikeCountElement.getAttribute("data-dislike-id");
        const postID = singlePostControl.closest('.card').getAttribute('data-post-id');
        const commentID = singlePostControl.closest('.card').getAttribute('data-comment-id');
        const channelID = singlePostControl.closest('.card').getAttribute('data-channel-id');



        likeButton.addEventListener('click', function(event) {
            let likeCount = parseInt(likeCountElement.textContent, 10);
            let dislikeCount = parseInt(dislikeCountElement.textContent, 10);

            if (likeButton.classList.contains('active')) {
                // Decrement the like count
                likeCountElement.textContent = `${likeCount - 1}`;
                // Remove the 'active' class to the like button
                likeButton.classList.remove('active');
            } else {
                // Increment the like count
                likeCountElement.textContent = `${likeCount + 1}`;
                // Add the 'active' class from the like button
                likeButton.classList.add('active');

                if (dislikeButton.classList.contains('active')) {
                    // Decrement the like count
                    dislikeCountElement.textContent = `${dislikeCount - 1}`;
                    // Remove the 'active' class to the like button
                    dislikeButton.classList.remove('active');
                }
            }

            let postData = checkData(commentID, postID, reactionAuthorID, channelID, "like")


            console.log("postData: ", postData)

            fetchData(postData, "like");
        });

        dislikeButton.addEventListener('click', function(event) {
            let likeCount = parseInt(likeCountElement.textContent, 10);
            let dislikeCount = parseInt(dislikeCountElement.textContent, 10);

            if (dislikeButton.classList.contains('active')) {
                // Decrement the like count
                dislikeCountElement.textContent = `${dislikeCount - 1}`;
                // Remove the 'active' class to the like button
                dislikeButton.classList.remove('active');
            } else {
                // Increment the like count
                dislikeCountElement.textContent = `${dislikeCount + 1}`;
                // Add the 'active' class from the like button
                dislikeButton.classList.add('active');

                if (likeButton.classList.contains('active')) {
                    // Decrement the like count
                    likeCountElement.textContent = `${likeCount - 1}`;
                    // Remove the 'active' class to the like button
                    likeButton.classList.remove('active');
                }
            }

            let postData = checkData(commentID, postID, reactionAuthorID, channelID, "dislike")

            console.log("postData: ", postData)


            fetchData(postData, "dislike");
        });
    });

});

function checkData(commentID, postID, reactionAuthorID, channelID, likeStatus) {
    let postData;

    if (likeStatus === "like") {
        if (postID === null || postID === 0) {
            // Send the updated like to the backend via POST request
            postData = {
                liked: true,
                disliked: false,
                reactedCommentId: parseInt(commentID, 10),
                authorId: parseInt(reactionAuthorID, 10),
            };
            console.log("liked comment: ", postData)

        } else if (commentID === null || commentID === 0) {
            // Send the updated like to the backend via POST request
            postData = {
                liked: true,
                disliked: false,
                reactedPostId: parseInt(postID, 10),
                authorId: parseInt(reactionAuthorID, 10),
            };
            console.log("liked post: ", postData)

        }
    } else if (likeStatus === "dislike") {
        if (postID === null || postID === 0) {
            // Send the updated like to the backend via POST request
            postData = {
                liked: false,
                disliked: true,
                reactedCommentId: parseInt(commentID, 10),
                authorId: parseInt(reactionAuthorID, 10),
            };
            console.log("disliked comment: ", postData)

        } else if (commentID === null || commentID === 0) {
            // Send the updated like to the backend via POST request
            postData = {
                liked: false,
                disliked: true,
                reactedPostId: parseInt(postID, 10),
                authorId: parseInt(reactionAuthorID, 10),
            };
            console.log("disliked post: ", postData)

        }
    }

    return postData;
}

function fetchData(postData, likeString) {
    fetch('http://localhost:8989/store-reaction', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(postData)
    })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => {
                    throw new Error(`HTTP Error: ${response.status} ${response.statusText}. Response: ${text}`);
                });
            }
            return response.json();
        })
        .then(data => {
            console.log(`${likeString} updated (in js):`, data);
        })
        .catch(error => {
            console.error('Error updating like:', error);
        });
}