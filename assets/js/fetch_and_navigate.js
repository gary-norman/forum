import { setActivePage, newContentLoaded } from "./main.js";
import { pages } from "./share.js";
const expect =
  "background-color: rgb(108 207 93); color: #000000; font-weight: bold; padding: .1rem; border-radius: 1rem;";
const standard =
  "background-color: transparent; color: #ffffff; font-weight: normal;";

export function navigateToPage(dest, entity) {
  // dest is a string - "channel"
  // entity is the template element:work
  // console.info("%cdest:%o", expect, dest);
  // console.info("%centity:%o", expect, entity);

  const link = entity.getAttribute(`data-${dest}-id`);
  // console.log("link: ", link);
  const page = dest + "Page";

  if (!link) {
    console.error(`${dest} ID is missing`);
    return Promise.reject(`${dest} ID is missing`);
  }

  setActivePage(dest);
  changePage(page);
  return fetchData(dest, link);
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
      // console.info(
      //   "%cset%o%c to %cactive-feed",
      //   expect,
      //   element.id,
      //   standard,
      //   expect,
      // );
    } else {
      // TODO need to modify home-page template to populate by injection
      // TODO when injected, the content can be cleared
      // clear the content of the page previously active
      // element.innerHTML = "";

      // TODO with injected home-page, this can be removed
      // need to clear content of the calendar dropdown, or the calendar won't work on otehr pages after navigation

      element.classList.remove("active-feed");
    }
  });
}

function fetchData(entity, Id) {
  // console.log(`%cFetching ${entity} data for ID:`, expect, Id);

  history.pushState({}, "", `/${entity}s/${Id}`);
  return fetch(`/${entity}s/${Id}`)
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
