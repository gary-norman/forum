function populate() {
    while(true) {
        // document bottom
        let windowRelativeBottom = document.documentElement.getBoundingClientRect().bottom;

        // if the user hasn't scrolled far enough (>100px to the end)
        if (windowRelativeBottom > document.documentElement.clientHeight + 100) break;

        // let's add more data
        document.body.insertAdjacentHTML("beforeend", `<div>stuff</div>`);
    }
}