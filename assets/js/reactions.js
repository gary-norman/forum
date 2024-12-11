
document.addEventListener('DOMContentLoaded', function () {

    const postControls = document.querySelectorAll('.post-controls');

    postControls.forEach(singlePostControl => {
        const likeCountElement = singlePostControl.querySelector('.btn-likes');
        const dislikeCountElement = singlePostControl.querySelector('.btn-dislikes');

        const likeButton = likeCountElement.closest('button');
        const dislikeButton = dislikeCountElement.closest('button');

        const likeID = likeCountElement.getAttribute("data-like-id");
        const dislikeID = dislikeCountElement.getAttribute("data-dislike-id");
        const postID = singlePostControl.closest('.card').getAttribute('data-post-id');
        const channelID = singlePostControl.closest('.card').getAttribute('data-channel-id');
        const authorID = singlePostControl.closest('.card').getAttribute('data-author-id');

        // console.log("likeID",likeID);
        // console.log("likeButton",likeButton, "\n")
        //
        // console.log("dislikeID",dislikeID)
        // console.log("dislike Button",dislikeButton, "\n")
        //
        // console.log("postID (js)",postID)

        likeButton.addEventListener('click', function(event) {
            let likeCount = parseInt(likeCountElement.textContent, 10);

            if (likeButton.classList.contains('reacted')) {
                // Decrement the like count
                likeCountElement.textContent = `${likeCount - 1}`;
                // Remove the 'reacted' class to the like button
                likeButton.classList.remove('reacted');
            } else {
                // Increment the like count
                likeCountElement.textContent = `${likeCount + 1}`;
                // Add the 'reacted' class from the like button
                likeButton.classList.add('reacted');
            }


            // Send the updated like to the backend via POST request
        const postData = {
            liked: true,
            disliked: false,
            reacted_post_id: parseInt(postID, 10),
            author_id: parseInt(authorID, 10),
            channel_id: parseInt(channelID, 10),
        };

        // console.log("postID",postData.reacted_post_id);
        // console.log("authorID",postData.author_id);
        // console.log("channelID",postData.channel_id);

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
                console.log('Like updated (in js):', data);
            })
            .catch(error => {
                console.error('Error updating like:', error);
            });
        });

        dislikeButton.addEventListener('click', function(event) {
            // console.log("clicked the dislike button");

            let dislikeCount = parseInt(dislikeCountElement.textContent, 10);

            // Decrement the like count
            dislikeCountElement.textContent = `${dislikeCount + 1}`;
        });
    });

});






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
