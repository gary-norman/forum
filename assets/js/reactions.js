
document.addEventListener('DOMContentLoaded', function () {

    const postControls = document.querySelectorAll('.post-controls');

    postControls.forEach(singlePostControl => {
        const likeCountElement = singlePostControl.querySelector('.btn-likes');
        const dislikeCountElement = singlePostControl.querySelector('.btn-dislikes');

        const likeButton = likeCountElement.closest('button');
        const dislikeButton = dislikeCountElement.closest('button');
        const sidebar = document.querySelector(`.sidebar`)

        const likeID = likeCountElement.getAttribute("data-like-id");
        const dislikeID = dislikeCountElement.getAttribute("data-dislike-id");
        const postID = singlePostControl.closest('.card').getAttribute('data-post-id');
        const commentID = singlePostControl.closest('.card').getAttribute('data-comment-id');
        const channelID = singlePostControl.closest('.card').getAttribute('data-channel-id');
        const reactionAuthorID = sidebar.querySelector('h3').getAttribute('data-current-user-ID');

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


            let postData1;

            if (postID === null) {
                // Send the updated like to the backend via POST request
                postData1 = {
                    liked: true,
                    disliked: false,
                    commented_comment_id: parseInt(commentID, 10),
                    author_id: parseInt(reactionAuthorID, 10),
                    channel_id: parseInt(channelID, 10),
                };
            } else if (commentID === null) {
                // Send the updated like to the backend via POST request
                postData1 = {
                    liked: true,
                    disliked: false,
                    reacted_post_id: parseInt(postID, 10),
                    author_id: parseInt(reactionAuthorID, 10),
                    channel_id: parseInt(channelID, 10),
                };
            }

            let postData = checkData(commentID, postID, reactionAuthorID, channelID, "like")

            // // Send the updated like to the backend via POST request
            // const postData = {
            //     liked: true,
            //     disliked: false,
            //     reacted_post_id: parseInt(postID, 10),
            //     author_id: parseInt(reactionAuthorID, 10),
            //     channel_id: parseInt(channelID, 10),
            // };

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
                reacted_comment_id: parseInt(commentID, 10),
                author_id: parseInt(reactionAuthorID, 10),
                channel_id: parseInt(channelID, 10),
            };
            console.log("liked comment: ", postData)

        } else if (commentID === null || commentID === 0) {
            // Send the updated like to the backend via POST request
            postData = {
                liked: true,
                disliked: false,
                reacted_post_id: parseInt(postID, 10),
                author_id: parseInt(reactionAuthorID, 10),
                channel_id: parseInt(channelID, 10),
            };
            console.log("liked post: ", postData)

        }
    } else if (likeStatus === "dislike") {
        if (postID === null || postID === 0) {
            // Send the updated like to the backend via POST request
            postData = {
                liked: false,
                disliked: true,
                reacted_comment_id: parseInt(commentID, 10),
                author_id: parseInt(reactionAuthorID, 10),
                channel_id: parseInt(channelID, 10),
            };
            console.log("disliked comment: ", postData)

        } else if (commentID === null || commentID === 0) {
            // Send the updated like to the backend via POST request
            postData = {
                liked: false,
                disliked: true,
                reacted_post_id: parseInt(postID, 10),
                author_id: parseInt(reactionAuthorID, 10),
                channel_id: parseInt(channelID, 10),
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





//
// // Function to handle like button click
// document.querySelectorAll('.post-controls').forEach( button => button.addEventListener('click', function(event) {
//
//     // Check if the clicked button is a like button
//     if (event.target.className.includes('like-btn')) {
//         console.log("clicked reaction button");
//         let button = event.target.closest('[class^="btn-"]');
//         let likeCountElement = button.querySelector('.btn-likes');
//         let likeCount = parseInt(likeCountElement.textContent, 10);
//         likeCountElement.textContent = likeCount + 1; // Increment the like count
//
//         // Get the post ID using the closest parent element with the data-post-id
//         let postID = button.closest('.card').getAttribute('data-post-id');
//
//         // Send the updated like to the backend via POST request
//         const postData = {
//             liked: true,
//             disliked: false,
//             postID: postID,
//             authorID: 123,  // Example: Replace with actual author ID
//             channelID: 456 // Example: Replace with actual channel ID
//         };
//
//         fetch('http://localhost:8989/store-reaction', {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json'
//             },
//             body: JSON.stringify(postData)
//         })
//             .then(response => response.json())
//             .then(data => {
//                 console.log('Like updated:', data);
//             })
//             .catch(error => {
//                 console.error('Error updating like:', error);
//             });
//     }
//
//     // Check if the clicked button is a dislike button
//     if (event.target.closest('.dislike-btn')) {
//         let button = event.target.closest('.dislike-btn');
//         let dislikeCountElement = button.querySelector('.btn-dislikes');
//         let dislikeCount = parseInt(dislikeCountElement.textContent, 10);
//         dislikeCountElement.textContent = dislikeCount + 1; // Increment the dislike count
//
//         // Get the post ID using the closest parent element with the data-post-id
//         let postID = button.closest('.card').getAttribute('data-post-id');
//
//         // Send the updated dislike to the backend via POST request
//         const postData = {
//             liked: false,
//             disliked: true,
//             postID: postID,
//             authorID: 123,  // Example: Replace with actual author ID
//             channelID: 456 // Example: Replace with actual channel ID
//         };
//
//         fetch('http://localhost:8989/store-reaction', {
//             method: 'POST',
//             headers: {
//                 'Content-Type': 'application/json'
//             },
//             body: JSON.stringify(postData)
//         })
//             .then(response => response.json())
//             .then(data => {
//                 console.log('Dislike updated:', data);
//             })
//             .catch(error => {
//                 console.error('Error updating dislike:', error);
//             });
//     }
// }));
//
