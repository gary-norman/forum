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

        shareButton.addEventListener('click', (e) => {
            getModalPos(shareButton, shareModal, window)
        });

        scrollWindow.addEventListener('scroll', (e) => {
            getModalPos(shareButton, shareModal, scrollWindow)
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

function getModalPos(shareButton, shareModal, windowPos) {
    //get button position
    const buttonPos = shareButton.getBoundingClientRect();


    //set modal styling
    shareModal.style.position = 'absolute';
    shareModal.style.top = `${buttonPos.bottom - 16}px`;
    shareModal.style.left = `${buttonPos.left  - 20}px`;
}


