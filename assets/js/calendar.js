import {addGlobalEventListener} from "./navigation.js";

const twoBodyContainer = document.querySelector('.two-body-container');
let isSelecting = false;
const dateStarting = document.querySelector('[data-date-starting]');
const dateEnding = document.querySelector('[data-date-ending]');
let yearButton = document.querySelector('[data-date="year"]');
let displayContainers = twoBodyContainer.querySelectorAll('.month-header');
let dayContainers = twoBodyContainer.querySelectorAll('.days-container');
let dateContainers = twoBodyContainer.querySelectorAll('.dates-container');


const date = new Date();
const month = date.getMonth();
const year = date.getFullYear();

let dateLeft = new Date(year, month, 1);
let dateRight = new Date(year, month + 1, 1);

let yearLeft = dateLeft.getFullYear();
let monthLeft = dateLeft.getMonth();

let yearRight = dateRight.getFullYear();
let monthRight = dateRight.getMonth();


displayCalendars();

addGlobalEventListener(
    "click",
    'button',
    e => {previousCalendar()},
    document.querySelector('.btn-previous').parentElement)

addGlobalEventListener(
    "click",
    'button',
    e => {nextCalendar()},
    document.querySelector('.btn-next').parentElement)


function displayCalendars() {
    // console.log("yearLeft", yearLeft)
    // console.log("monthLeft", monthLeft)
    // console.log("yearRight", yearRight)
    // console.log("monthRight", monthRight)
    //get the first day and index of the month
    const firstDayLeft = new Date(yearLeft, monthLeft, 1);
    let firstDayIndexLeft = firstDayLeft.getDay();
    // shift the indexes, so that Monday is index 0 (not 1)
    firstDayIndexLeft = (firstDayIndexLeft === 0) ? 6 : firstDayIndexLeft - 1;

    // do the same for right calendar
    const firstDayRight = new Date(yearRight, monthRight, 1);
    let firstDayIndexRight = firstDayRight.getDay();
    firstDayIndexRight = (firstDayIndexRight === 0) ? 6 : firstDayIndexRight - 1;

    // console.log("firstDayLeft", firstDayLeft)
    // console.log("firstDayLeftIndex", firstDayIndexLeft)
    //
    // console.log("firstDayRight", firstDayRight)
    // console.log("firstDayRightIndex", firstDayIndexRight)
    const lastDayLeft = new Date(yearLeft, monthLeft + 1, 0);
    const numberOfDaysLeft = lastDayLeft.getDate();

    const lastDayRight = new Date(yearRight, monthRight + 1, 0);
    const numberOfDaysRight = lastDayRight.getDate();

    console.log("dateLeft", dateLeft)
    console.log("dateRight", dateRight)
    let formattedDateLeft = dateLeft.toLocaleString("en-GB", {
        month: "long",
        year: "numeric",
    });

    let formattedDateRight = dateRight.toLocaleString("en-GB", {
        month: "long",
        year: "numeric",
    });

    displayContainers[0].textContent = `${formattedDateLeft}`;
    displayContainers[1].textContent = `${formattedDateRight}`;

    dayContainers.forEach(container => {
        for (let i = 1; i <= 7; i++) {
            const dayCell = document.createElement('div');
            dayCell.className = 'date-cell day-name';
            switch (i) {
                case 1 :
                    dayCell.textContent = "M";
                    break;
                case 2:
                    dayCell.textContent = "T";
                    break;
                case 3:
                    dayCell.textContent = "W";
                    break;
                case 4:
                    dayCell.textContent = "T";
                    break;
                case 5:
                    dayCell.textContent = "F";
                    break;
                case 6:
                    dayCell.textContent = "S";
                    break;
                case 7:
                    dayCell.textContent = "S";
                    break;
            }
            container.appendChild(dayCell);
        }
    })

    dateContainers.forEach(container => {
        if (container === dateContainers[0]) {
            // populate days of the week that belong to the previous month
            for (let x = 1; x <= firstDayIndexLeft; x++) {
                populateEmptyDateCell(container);
            }
            for (let i = 1; i <= numberOfDaysLeft; i++) {
                populateDateCell(container, i)
            }
        } else if (container === dateContainers[1]){
            for (let x = 1; x <= firstDayIndexRight; x++) {
                populateEmptyDateCell(container);
            }
            for (let i = 1; i <= numberOfDaysRight; i++) {
                populateDateCell(container, i)
            }
        }
    })
    yearButton.textContent = yearLeft;
}

function previousCalendar(){
    dateContainers.forEach(container => {
        container.innerHTML = "";
    })
    dayContainers.forEach(container => {
        container.innerHTML = "";
    })
    // manage if the starting month was in left/ or right side
    if (monthRight <= 0) {
        console.warn("monthRight <= 0");
        monthRight = 10;
        monthLeft = 9;
        yearLeft = yearLeft - 1;
        yearRight = yearLeft;
        yearButton.textContent = yearLeft;
    } else if (monthLeft <= 0) {
        console.warn("monthLeft <= 0");
        monthRight = 11;
        monthLeft = 10;
        yearRight = yearLeft - 1;
        yearLeft = yearRight;
        yearButton.textContent = yearLeft;
    } else {
        console.warn("else");
        monthLeft = monthLeft - 2;
        monthRight = monthRight - 2;
    }
    dateRight.setFullYear(yearRight);
    dateRight.setMonth(monthRight);
    dateLeft.setFullYear(yearLeft);
    dateLeft.setMonth(monthLeft);
    displayCalendars();
}

function nextCalendar() {
    dateContainers.forEach(container => {
        container.innerHTML = "";
    })
    dayContainers.forEach(container => {
        container.innerHTML = "";
    })

    // manage if the starting month was in left/ or right side
    if (monthRight >= 11) {
        console.warn("monthRight >= 11");
        monthRight = 1;
        monthLeft = 0;
        yearLeft = yearLeft + 1;
        yearRight = yearLeft;
        yearButton.textContent = yearLeft;
    } else if (monthLeft >= 11) {
        console.warn("monthLeft >= 11");
        monthRight = 0;
        monthLeft = 11;
        yearRight = yearLeft + 1;
        yearButton.textContent = yearLeft;
    } else {
        monthLeft = monthLeft + 2;
        monthRight = monthRight + 2;
    }
    dateRight.setFullYear(yearRight);
    dateRight.setMonth(monthRight);
    dateLeft.setFullYear(yearLeft);
    dateLeft.setMonth(monthLeft);
    displayCalendars();
}

function processDateRange() {

    dateContainers.forEach(container => {

        const dateCells = container.querySelectorAll('.date-cell');

        console.log("running date script")


        dateCells.forEach((dateCell) => {
            dateCell.addEventListener('click', () => {
                if (isSelecting) {
                    if (dateCell.classList.contains('active')) {
                        dateCells.forEach(cell => {
                            cell.classList.remove('selected');
                            cell.classList.remove('active');
                        });
                    }
                    isSelecting = false;
                }
            });
        });

        dateCells.forEach(dateCell => {
            dateCell.addEventListener('click', () => {
                dateCell.classList.add('active');
                toggleDateTexts();
                isSelecting = true;
                // console.log("isSelecting: ", isSelecting);
            });


            // Hover in
            dateCell.addEventListener('mouseenter', () => {
                const activeIndex = Array.from(dateCells).findIndex(cell => cell.classList.contains('active'));
                const hoveredIndex = Array.from(dateCells).indexOf(dateCell);

                if (activeIndex !== -1) {
                    // console.log("activeIndex: ",activeIndex)

                    dateCells.forEach(cell => cell.classList.remove('selected'));

                    if (hoveredIndex > activeIndex) {
                        dateCells.forEach((cell, index) => {
                            if (index > activeIndex && index < hoveredIndex) {
                                cell.classList.add('selected');
                            }
                        });
                    }
                }

            });

        })
    })

}

function toggleDateTexts() {
    dateStarting.classList.toggle('hide');
    dateEnding.classList.toggle('hide');
}

function populateDateCell(container, i) {
    let dateCell = document.createElement("div");
    let currentDate;
    switch (container) {
        case dateContainers[0]:
            currentDate= new Date(yearLeft, monthLeft, i);
            break;
        case dateContainers[1]:
            currentDate= new Date(yearRight, monthRight, i);
    }

    dateCell.dataset.date = currentDate.toDateString();
    dateCell.className = 'date-cell';
    dateCell.innerHTML += i;
    container.appendChild(dateCell);
    if (
        currentDate.getFullYear() === new Date().getFullYear() &&
        currentDate.getMonth() === new Date().getMonth() &&
        currentDate.getDate() === new Date().getDate()
    ) {
        dateCell.classList.add("current-date");
    }
}

function populateEmptyDateCell(container){
    let dateCell = document.createElement("div");

    // TODO need to add day text for the days of the previous month
    dateCell.className = 'date-cell disabled';
    dateCell.innerHTML += "";
    container.appendChild(dateCell);
}