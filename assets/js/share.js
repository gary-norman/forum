const link = encodeURI(window.location.href);
let msg;
let title;

const commentMsg = encodeURIComponent('Hey, I found this comment, you need to see it!');
const postMsg = encodeURIComponent('Hey, I found this post. I think you may like it?');
const  commentTitle = encodeURIComponent('Comment from User ??? Here');
const postTitle = encodeURIComponent('Post Title Here');
const scrollWindow = document.getElementById("activity-feed-activity");

document.addEventListener('DOMContentLoaded', function () {
    let postID;
    let commentID;
    const postControls = document.querySelectorAll(`.post-controls`);

    postControls.forEach(singlePostControl => {
        postID = singlePostControl.closest('.card').getAttribute('data-post-id');
        commentID = singlePostControl.closest('.card').getAttribute('data-comment-id');

        //get all components needed
        const shareModal = singlePostControl.querySelector(`#share-container-`+ postID);
        const shareButton = singlePostControl.querySelector(`#share-button-`+ postID);
        const label = shareModal.querySelector('label');
        const icon = shareModal.querySelector('button');
        const input = shareModal.querySelector('input');

        shareButton.addEventListener('click', (e) => {
            getModalPos(shareButton, shareModal, window)
        });

        // Listen for the 'toggle' event on the modal (native popover event)
        shareModal.addEventListener('toggle', () => {
            toggleButtonActive(shareModal, shareButton);
        });

        scrollWindow.addEventListener('scroll', (e) => {
            getModalPos(shareButton, shareModal, scrollWindow)

            // Parse the 'top' value and compare it with top of the mask / botton of screen
            const modalTop = parseInt(shareModal.style.top, 10); // Convert 'top' to a number

            if (modalTop <= 400 || modalTop >= window.innerHeight - 72) {
                // Close the popover
                shareModal.hidePopover();
            }
        });

        label.addEventListener('click', async () => {
            try {
                await updateCopyStatus(input, label, icon);
                console.log('Copy operation completed');
            } catch (err) {
                console.error('Error during copy operation:', err);
            }
        });


        icon.addEventListener('click', async () => {
            try {
                await updateCopyStatus(input, label, icon);
                console.log('Copy operation completed');
            } catch (err) {
                console.error('Error during copy operation:', err);
            }
        });


        //if postID present, make msg = postMsg, if commentID present, make msg = commentMsg
        if (postID !== 0 && commentID === 0) {
            msg = postMsg;
            title = postTitle;
        } else if (commentID !== 0 && postID === 0) {
            msg = commentMsg;
            title = commentTitle;
        }

        const fb = singlePostControl.querySelector('.facebook');
        fb.href = `https://www.facebook.com/share.php?u=${link}`;

        const twitter = singlePostControl.querySelector('.twitter');
        twitter.href = `http://twitter.com/share?&url=${link}&text=${msg}&hashtags=javascript,programming`;

        const linkedIn = singlePostControl.querySelector('.linkedin');
        linkedIn.href = `https://www.linkedin.com/sharing/share-offsite/?url=${link}`;

        const reddit = singlePostControl.querySelector('.reddit');
        reddit.href = `http://www.reddit.com/submit?url=${link}&title=${title}`;

        const whatsapp = singlePostControl.querySelector('.whatsapp');
        whatsapp.href = `https://api.whatsapp.com/send?text=${msg}: ${link}`;

        const telegram = singlePostControl.querySelector('.telegram');
        telegram.href = `https://telegram.me/share/url?url=${link}&text=${msg}`;
    });
});

function scrollToPost(postId) {
    const container = scrollWindow;
    const post = document.querySelector(`[data-post-id="${postId}"]`)


    if (post && container) {
        const containerRect = container.getBoundingClientRect();
        const postRect = post.getBoundingClientRect();

        // Calculate the position relative to the container
        const scrollOffset = postRect.top - containerRect.top + container.scrollTop;

        container.scrollTo({
            top: scrollOffset,
            behavior: 'smooth', // Smooth scrolling animation
        });

        post.classList.add('card-selected');

        setTimeout(() => {
            post.classList.remove('card-selected');
            }, 3000);
    } else {
        console.error('Post or container not found:', postId);
    }
}

const scrollButton = document.getElementById(`scroll-test`);
scrollButton.addEventListener('click', () => {
    // Example: Scroll to post 3
    scrollToPost('3');
})


function getModalPos(shareButton, shareModal, windowPos) {
    //get button position
    const buttonPos = shareButton.getBoundingClientRect();

    //set modal styling
    shareModal.style.position = 'absolute';
    shareModal.style.top = `${buttonPos.bottom - 16}px`;
    shareModal.style.left = `${buttonPos.left  - 20}px`;
}

function toggleButtonActive(shareModal, shareButton) {
    if (shareModal.matches(':popover-open')) {
        shareButton.classList.add('active');
    } else {
        shareButton.classList.remove('active');
    }
}

function copyToClipboard(input) {
    input.focus();
    input.select();

    // Modern Clipboard API (recommended)
    if (navigator.clipboard) {
        return navigator.clipboard.writeText(input.value)
            .then(() => console.log('Text copied to clipboard!'))
            .catch(err => {
                console.error('Failed to copy text:', err);
                throw err; // Ensure errors are propagated to the caller
            });
    } else {
        // Fallback for unsupported browsers
        console.warn('Clipboard API not supported, using execCommand as fallback');
        try {
            const success = document.execCommand('copy');
            if (!success) {
                throw new Error('Fallback copy failed');
            }
            console.log('Text copied to clipboard (fallback method)!');
            return Promise.resolve(); // ✅ Return a resolved Promise
        } catch (err) {
            console.error('Fallback copy failed:', err);
            return Promise.reject(err); // ✅ Return a rejected Promise
        }
    }
}

async function updateCopyStatus(input, label, icon) {
    const iconSpan = icon.querySelector(`span`);
    try {
        await copyToClipboard(input);

        // Update label text and styles
        label.textContent = 'Copied!';
        label.style.color = "var(--color-hl-primary)";
        iconSpan.classList.remove("btn-copy");
        iconSpan.classList.add("btn-success");

        setTimeout(() => {
            label.textContent = 'Copy URL link';
            label.style.color = "var(--color-fg-2)";
            iconSpan.classList.add("btn-copy");
            iconSpan.classList.remove("btn-success");
        }, 2000);
    } catch (err) {
        console.error('Failed to copy text:', err);
        label.textContent = 'Failed to Copy';
        label.style.color = "var(--color-hl-red)";
        iconSpan.classList.remove("btn-copy");
        iconSpan.classList.add("btn-warning");

        setTimeout(() => {
            label.textContent = 'Copy URL';
            label.style.color = "var(--color-fg-2)";
            iconSpan.classList.add("btn-copy");
            iconSpan.classList.remove("btn-warning");
        }, 2000);
    }
}
