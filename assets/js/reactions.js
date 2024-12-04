// Function to handle like button click
document.querySelector('.button-row').addEventListener('click', function(event) {
    // Check if the clicked button is a like button
    if (event.target.closest('.like-btn')) {
        let button = event.target.closest('.like-btn');
        let likeCountElement = button.querySelector('.btn-likes');
        let likeCount = parseInt(likeCountElement.textContent, 10);
        likeCountElement.textContent = likeCount + 1; // Increment the like count

        // Get the post ID using the closest parent element with the data-post-id
        let postID = button.closest('.card').getAttribute('data-post-id');

        // Send the updated like to the backend via POST request
        const postData = {
            liked: true,
            disliked: false,
            postID: postID,
            authorID: 123,  // Example: Replace with actual author ID
            channelID: 456 // Example: Replace with actual channel ID
        };

        fetch('http://localhost:8989/store-reaction', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(postData)
        })
            .then(response => response.json())
            .then(data => {
                console.log('Like updated:', data);
            })
            .catch(error => {
                console.error('Error updating like:', error);
            });
    }

    // Check if the clicked button is a dislike button
    if (event.target.closest('.dislike-btn')) {
        let button = event.target.closest('.dislike-btn');
        let dislikeCountElement = button.querySelector('.btn-dislikes');
        let dislikeCount = parseInt(dislikeCountElement.textContent, 10);
        dislikeCountElement.textContent = dislikeCount + 1; // Increment the dislike count

        // Get the post ID using the closest parent element with the data-post-id
        let postID = button.closest('.card').getAttribute('data-post-id');

        // Send the updated dislike to the backend via POST request
        const postData = {
            liked: false,
            disliked: true,
            postID: postID,
            authorID: 123,  // Example: Replace with actual author ID
            channelID: 456 // Example: Replace with actual channel ID
        };

        fetch('http://localhost:8989/store-reaction', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(postData)
        })
            .then(response => response.json())
            .then(data => {
                console.log('Dislike updated:', data);
            })
            .catch(error => {
                console.error('Error updating dislike:', error);
            });
    }
});