import {modifyActivePage, newContentLoaded} from "./main.js";
import {pages} from "./share.js";

export function navigateToPage(dest, entity) {
    const link = entity.getAttribute(`data-${dest}-id`);
    // console.log("link: ", link);
    const page = dest + "Page";

    if (!link) {
        console.error(`${dest} ID is missing`);
        return;
    }

    modifyActivePage(dest);
    changePage(page);
    fetchData(dest, link);
}

export function changePage(page) {
    // console.log("page: ", page);
    let pageId;
    if (typeof page != "string") {
        pageId = page.id;
    } else {
        const pageSlice = page.slice(0, -4);
        pageId = pageSlice + "-page";
    }
    pages.forEach((element) => {
        // console.log("page", page);
        // console.log("pageId: ", pageId);
        // console.log("element.id: ", element.id);
        if (element.id === pageId) {
            element.classList.add("active-feed");
            // console.log("set", element.id, "to active-feed");
        } else {
            element.classList.remove("active-feed");
        }
    });
}

function fetchData(entity, Id) {
    // console.log(`Fetching ${entity} data for ID:`, Id);

    fetch(`/${entity}s/${Id}`)
        .then((response) => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then((data) => {
            const content = data[`${entity}sHTML`];
            const target = document.getElementById(`${entity}-page`);
            if (target && content) {
                target.innerHTML = content;

                // Update URL for SPA routing
                // window.history.pushState({}, "", `/${entity}s/${Id}`);

                // Dispatch custom event
                document.dispatchEvent(newContentLoaded);
            } else {
                console.warn("Target element or content missing");
            }
        })
        .catch((error) => console.error(`Error fetching ${entity} data:`, error));
}



