import { setActivePage, newContentLoaded } from "./main.js";
import { pages } from "./share.js";

export function navigateToPage(dest, entity) {
  // dest is a string - "channel"
  // entity is the template element
  console.log("dest:", dest);
  console.log("entity:", entity);

  const link = entity.getAttribute(`data-${dest}-id`);
  // console.log("link: ", link);
  const page = dest + "Page";

  if (!link) {
    console.error(`${dest} ID is missing`);
    return;
  }

  setActivePage(dest);
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
      console.log("set", element.id, "to active-feed");
    } else {
      // TODO need to modify home-page template to populate by injection
      // TODO when injected, the content can be cleared
      // clear the content of the page previously active
      // element.textContent = "";

      // TODO with injected home-page, this can be removed
      // need to clear content of the calendar dropdown, or the calendar won't work on otehr pages after navigation


      element.classList.remove("active-feed");
    }
  });
}

function fetchData(entity, Id) {
  console.log(`Fetching ${entity} data for ID:`, Id);

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
      // console.log("target: ", target);
      // console.log("content: ", content);
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