import {updateProfileImages} from "./update_UI_elements.js";

const searchInput = document.querySelector("[data-search]");
const userCardContainer = document.querySelector("[data-result-user-cards-container]");
const channelCardContainer = document.querySelector("[data-result-channel-cards-container]");
const postCardContainer = document.querySelector("[data-result-post-cards-container]");
const userCardTemplate = document.querySelector("[data-result-user-cards-template]");
const channelCardTemplate = document.querySelector("[data-result-channel-cards-template]");
const postCardTemplate = document.querySelector("[data-result-post-cards-template]");
const userResultsContainer = document.getElementById("results-users");
const channelResultsContainer = document.getElementById("results-channels");
const postResultsContainer = document.getElementById("results-posts");
const resultsContainer = document.getElementById("search-results-page");
const dividers = document.querySelectorAll(`hr`);
const prefix = "noimage"


let isFocus = false;
let isValue = false;


let users = [];
let channels = [];
let posts = [];

document.addEventListener("click", (e) => {
    if (e.target !== searchInput) {
        resultsContainer.classList.add("hide");
    } else if (e.target === searchInput && isValue ){
        resultsContainer.classList.remove("hide");
    }
})

searchInput.addEventListener("input", e => {
    const value = e.target.value.toLowerCase();
    let anyUserVisible, anyChannelVisible, anyPostVisible = false;
    let anyVisible = (anyUserVisible || anyChannelVisible || anyPostVisible);
    isValue = value !== "";
    isFocus = false;



    searchInput.addEventListener("focus", () => {
        isFocus = true;
    })

    resultsContainer.classList.toggle("hide", !(isFocus || isValue));


    dividers.forEach(divider => {
        divider.classList.toggle("hide", !anyVisible)
    })

    users.forEach(user => {
        const isVisible = user.username.toLowerCase().includes(value);
        user.element.classList.toggle("hide", !isVisible)
        if (isVisible) {
            anyUserVisible = true;
            anyVisible = (anyUserVisible || anyChannelVisible || anyPostVisible);
        }
    });

    channels.forEach(channel => {
        const isVisible = channel.name.toLowerCase().includes(value);
        channel.element.classList.toggle("hide", !isVisible)
        if (isVisible) {
            anyChannelVisible = true;
            anyVisible = (anyUserVisible || anyChannelVisible || anyPostVisible);
        }
    })


    posts.forEach(post => {
        const isVisible = post.title.toLowerCase().includes(value) || post.content.toLowerCase().includes(value);
        post.element.classList.toggle("hide", !isVisible)
        if (isVisible) {
            anyPostVisible = true;
            anyVisible = (anyUserVisible || anyChannelVisible || anyPostVisible);
        }
    })


    // Hide the parent container if no users are visible
    if (userResultsContainer) {
        userResultsContainer.classList.toggle("hide", !anyUserVisible);
    }

    // Hide the parent container if no channels are visible
    if (channelResultsContainer) {
        channelResultsContainer.classList.toggle("hide", !anyChannelVisible);
    }

    // Hide the parent container if no posts are visible
    if (postResultsContainer) {
        postResultsContainer.classList.toggle("hide", !anyPostVisible);
    }

    let calcHeight = `${Math.max(7, calculateVisibleChildrenHeight(resultsContainer))}rem`;
    if (calcHeight === "7.8rem") {
        calcHeight = "7rem"
    }

    if (!anyVisible) {
        //Prepare No results paragraph
        const noResults = document.createElement('p');
        noResults.textContent = "Search input doesn't match any items!";
        noResults.classList.add('no-result');
        noResults.style.padding = "2.4rem 2.4rem"
        noResults.style.textAlign = 'center';

        const noResultsGroup = resultsContainer.querySelectorAll(".no-result");

        if (noResultsGroup.length === 0) {
            resultsContainer.appendChild(noResults);
        }

        resultsContainer.style.padding = "0"
    } else {
        const noResults = resultsContainer.querySelectorAll(".no-result");
        if (noResults) {
            noResults.forEach(paragraph => {
                resultsContainer.removeChild(paragraph);
            })
        }
    }

    resultsContainer.style.height = calcHeight;
})

fetch("/search")
    .then((response) => {
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    })
    .then((data) => {
        users = data.users.map(user => {
            const card = userCardTemplate.content.cloneNode(true).children[0];
            const avatar = card.querySelector("[data-result-user-avatar]");
            const name = card.querySelector("[data-result-user-name]");
            const imageAttr = "db/userdata/images/user-images/";
            const id = user.id;

            name.textContent = user.username;
            card.setAttribute("data-user-id", id)
            if (user.avatar.startsWith(prefix) {
                avatar.classList.add("profile-pic--empty")
                avatar.classList.remove("profile-pic")
                avatar.setAttribute("data-name-user", user.name);
                // console.log("added a placeholder image for channel: ", channel.name);
            } else {
                avatar.setAttribute("data-image-user", `${imageAttr + user.avatar}` );
            }

            userCardContainer.append(card)
            return { username: user.username, avatar: user.avatar, element: card }
        })

        channels = data.channels.map(channel => {
            const card = channelCardTemplate.content.cloneNode(true).children[0];
            const avatar = card.querySelector("[data-result-channel-avatar]");
            const name = card.querySelector("[data-result-channel-name]");
            const imageAttr = "db/userdata/images/channel-images/";
            const id = channel.id;

            name.textContent = "/" + channel.name;
            card.setAttribute("data-channel-id", id);
            if (channel.avatar.startsWith(prefix) {
                avatar.classList.add("profile-pic--empty")
                avatar.classList.remove("profile-pic")
                avatar.setAttribute("data-name-channel", channel.name);
                // console.log("added a placeholder image for channel: ", channel.name);
            } else {
                avatar.setAttribute("data-image-channel", `${imageAttr + channel.avatar}` );
            }

            channelCardContainer.append(card)
            return { name: channel.name, avatar: channel.avatar, element: card }
        })

        posts = data.posts.map(post => {
            const card = postCardTemplate.content.cloneNode(true).children[0];
            const title = card.querySelector("[data-result-post-title]");
            const content = card.querySelector("[data-result-post-content]");
            const id = post.id;

            title.textContent = post.title;
            content.textContent = post.content;
            card.setAttribute("data-post-id", id)

            postCardContainer.append(card)
            return { title: post.title, content: post.content, element: card }
        })

        updateProfileImages();
    })
    .catch((error) => console.error(`Error fetching response data:`, error));


function containsElement(parentChildren, elementToCheck) {
    for (let i = 0; i < parentChildren.length; i++) {
        if (parentChildren[i] === elementToCheck) {
            return true;
        }
    }
    return false;
}

function calculateVisibleChildrenHeight(resultsContainer) {
    if (!resultsContainer) {
        console.error(`Container with ID "${resultsContainer}" not found.`);
        return 0;
    }

    let totalVisibleHeight = 0;
    const children = resultsContainer.children;

    for (let i = 0; i < children.length; i++) {
        const element = children[i];
        const isVisible = element.classList.contains("hide");

        if (!isVisible) {
            // Check if the element is visible
            const style = window.getComputedStyle(element);

            // Get the element's height (including padding and border)
            totalVisibleHeight += element.offsetHeight;

            // If you want to include margins as well, you can add:
            const marginTop = parseInt(style.marginTop) || 0;
            const marginBottom = parseInt(style.marginBottom) || 0;
            totalVisibleHeight += marginTop + marginBottom;
        }
    }

    return Math.min(320, totalVisibleHeight) / 10 + 1;
}
