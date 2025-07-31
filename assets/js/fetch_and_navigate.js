import { setActivePage, newContentLoaded } from "./main.js";
import { pages, selectActiveFeed } from "./share.js";
const expect =
  "background-color: rgb(108 207 93); color: #000000; font-weight: bold; padding: .1rem; border-radius: 1rem;";
const standard =
  "background-color: transparent; color: #ffffff; font-weight: normal;";
const warn =
    "background-color: #000000; color: #e3c144; font-weight: bold; border: 1px solid #e3c144; padding: 5px; border-radius: 5px;";


export function navigateToPage(dest, entity) {
  // dest is a string - "channel"
  // entity is the template element:work
  console.info("%cdest:%o", expect, dest);
  console.info("%centity:%o", expect, entity);

  let link;
  if (dest === "home") {
    link = entity.getAttribute(`data-dest`);
  } else {
    link = entity.getAttribute(`data-${dest}-id`);
  }
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

    if (element.id === pageId) {
      element.classList.add("active-feed");
    } else {
      // console.log("%cpage: ",warn, page);
      // TODO need to modify home-page template to populate by injection
      // TODO when injected, the content can be cleared

      // clear the content of the page previously active
      element.innerHTML = "";

      // TODO with injected home-page, this can be removed
      // need to clear content of the calendar dropdown, or the calendar won't work on otehr pages after navigation

      element.classList.remove("active-feed");
    }
  });
}

export function fetchData(entity, Id) {
  console.log(`%cFetching ${entity} data for ID:`, expect, Id);


  if (Id === "home") {
    console.log("home");
    history.pushState({}, "", `/`);
    fetch(`/`)
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
            document.dispatchEvent(newContentLoaded);
          } else {
            console.warn("Target element or content missing");
          }
        })
        .catch((error) => console.error(`Error fetching ${entity} data:`, error));
  } else {
    console.log("other");
    history.pushState({}, "", `/${entity}s/${Id}`);
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
            document.dispatchEvent(newContentLoaded);
          } else {
            console.warn("Target element or content missing");
          }
        })
        .catch((error) => console.error(`Error fetching ${entity} data:`, error));
  }
}
