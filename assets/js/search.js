const searchInput = document.querySelector("[data-search]");
const userCardContainer = document.querySelector("[data-result-user-cards-container]")
const channelCardContainer = document.querySelector("[data-result-channel-cards-container]")


let users = [];
let channels = [];

console.log(searchInput, userCardContainer, channelCardContainer)

searchInput.addEventListener("input", e => {
    const value = e.target.value.toLowerCase();

    console.log(value);
    users.forEach(user => {
        const isVisible = user.username.toLowerCase().includes(value);
        user.element.classList.toggle("hide", !isVisible)
    })
    channels.forEach(channel => {
        const isVisible = channel.name.toLowerCase().includes(value);
        channel.element.classList.toggle("hide", !isVisible)
    })
})

fetch()
    .then((response) => {
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    })
    .then((data) => {
        users = data.map(user => {
            const card = userCardContainer.content.cloneNode(true).children[0];
            const avatar = card.querySelector("[data-result-user-avatar]");
            const name = card.querySelector("[data-result-user-name]");
            const imageAttr = avatar.getAttribute("data-image-user");

            name.textContent = user.username;
            avatar.style.background = `url('${imageAttr}') no-repeat center`;

            userCardContainer.append(card)
            return { username: user.username, avatar: user.avatar, element: card }
        })

        channels = data.map(channel => {
            const card = channelCardContainer.content.cloneNode(true).children[0];
            const avatar = card.querySelector("[data-result-channel-avatar]");
            const name = card.querySelector("[data-result-channel-name]");
            const imageAttr = avatar.getAttribute("data-image-channel");

            name.textContent = channel.name;
            avatar.style.background = `url('${imageAttr}') no-repeat center`;

            channelCardContainer.append(card)
            return { name: channel.name, avatar: channel.avatar, element: card }
        })
    })
    .catch((error) => console.error(`Error fetching ${entity} data:`, error));