// INFO was a DOMContentLoaded function
import { activePage } from "./main.js";
import {showInlineNotification} from "./notifications.js";

let replyForm, notifierReply;

export function toggleReplyForm(target) {
  const postCard = target.closest(".card");
  const postForm = postCard.querySelector(".form-reply");
  const textArea = postForm.querySelector('[id^="comment-form-textarea-"]');

  postForm.classList.add("replying");
  textArea.value = "";
  textArea.focus();
  listenToReplyForm();
}


function listenToReplyForm() {

    replyForm = document.querySelector(".form-reply");
    console.log("replyForm:", replyForm)

    const replyContent = document.querySelector('textarea[name="content"]');
    console.log("replyContent:", replyContent)

  const formCard = replyForm.closest('.card');

    notifierReply = formCard.querySelector(".popover-title");
    console.log("notifierReply:", notifierReply)

    replyForm.addEventListener("submit", async function (e) {
      e.preventDefault();
      const form = e.target;
      const formData = new FormData(form); // Collect form data
      const content = formData.get("content")?.trim();

      if (!content || content.trim().length <1) {
        showInlineNotification(
            notifierReply,
            "",
            "enter some content",
            false,
            "invisible-notify",
        );
        replyContent.focus();

      }
    });
}