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
let isFocus = false;
let isValue = false;


let users = [];
let channels = [];
let posts = [];

document.addEventListener("click", (e) => {
    console.log("e.target: ", e.target)
    if (e.target !== searchInput) {
        resultsContainer.classList.add("hide");
    } else if (e.target === searchInput && isValue ){
        resultsContainer.classList.remove("hide");
    }
})

searchInput.addEventListener("input", e => {
    const value = e.target.value.toLowerCase();
    let anyUserVisible, anyChannelVisible, anyPostVisible = false;
    isValue = value !== "";
    isFocus = false;

    searchInput.addEventListener("focus", () => {
        isFocus = true;
    })

    resultsContainer.classList.toggle("hide", !(isFocus || isValue))

    users.forEach(user => {
        const isVisible = user.username.toLowerCase().includes(value);
        user.element.classList.toggle("hide", !isVisible)
        if (isVisible) {
            anyUserVisible = true;
        }
    });

    // Hide the parent container if no users are visible
    if (userResultsContainer) {
        userResultsContainer.classList.toggle("hide", !anyUserVisible);
    }

    channels.forEach(channel => {
        const isVisible = channel.name.toLowerCase().includes(value);
        channel.element.classList.toggle("hide", !isVisible)
        if (isVisible) {
            anyChannelVisible = true;
        }
    })

    // Hide the parent container if no channels are visible
    if (channelResultsContainer) {
        channelResultsContainer.classList.toggle("hide", !anyChannelVisible);
    }

    posts.forEach(post => {
        const isVisible = post.title.toLowerCase().includes(value) ||  post.content.toLowerCase().includes(value);
        post.element.classList.toggle("hide", !isVisible)
        if (isVisible) {
            anyPostVisible = true;
        }
    })

    // Hide the parent container if no posts are visible
    if (postResultsContainer) {
        console.log("posts container visible")
        postResultsContainer.classList.toggle("hide", !anyPostVisible);
    }

})

fetch("/search")
    .then((response) => {
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    })
    .then((data) => {
        console.table(data)
        users = data.users.map(user => {
            const card = userCardTemplate.content.cloneNode(true).children[0];
            const avatar = card.querySelector("[data-result-user-avatar]");
            const name = card.querySelector("[data-result-user-name]");
            const imageAttr = "db/userdata/images/user-images/";
            const id = user.id;

            name.textContent = user.username;
            avatar.style.background = `url('${imageAttr + user.avatar}') no-repeat center`;
            avatar.style.backgroundSize = "cover";
            card.setAttribute("data-user-id", id)

            userCardContainer.append(card)
            return { username: user.username, avatar: user.avatar, element: card }
        })

        channels = data.channels.map(channel => {
            const card = channelCardTemplate.content.cloneNode(true).children[0];
            const avatar = card.querySelector("[data-result-channel-avatar]");
            const name = card.querySelector("[data-result-channel-name]");
            const imageAttr = "db/userdata/images/channel-images/";
            const id = channel.id;

            name.textContent = channel.name;
            avatar.style.background = `url('${imageAttr + channel.avatar}') no-repeat center`;
            avatar.style.backgroundSize = "cover";
            card.setAttribute("data-channel-id", id)

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
    })
    .catch((error) => console.error(`Error fetching response data:`, error));
