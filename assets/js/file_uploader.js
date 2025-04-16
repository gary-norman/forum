function fileUploader() {
    // adapted from https://medium.com/@cwrworksite/drag-and-drop-file-upload-with-preview-using-javascript-cd85524e4a63
    // user
    const dropAreaUser = document.querySelector("#drop-zone--user-image");
    let dropButtonUser, inputUser, uploadedFileUser, file, filename;

    // post
    const dropAreaPost = document.querySelector("#drop-zone--post");
    const dropButtonPost = dropAreaPost.querySelector(".button");
    const inputPost = dropAreaPost.querySelector("input");
    const uploadedFilePost = document.querySelector("#uploaded-file--post");
    const dragText = document.querySelector(".dragText");
    const dragButton = document.querySelector(".button");

    if (dropAreaUser) {
        dropButtonUser = dropAreaUser.querySelector(".button");
        inputUser = dropAreaUser.querySelector("input");
        uploadedFileUser = document.querySelector("#uploaded-file--user-image");
    }

    if (inputUser) {
        inputUser.addEventListener("change", function () {
            file = this.files[0];
            dropAreaUser.classList.add("active");
        });

        // when file is inside drag area
        dropAreaUser.addEventListener("dragover", (event) => {
            event.preventDefault();
            dropAreaUser.classList.add("active");
            dragText.textContent = "release to Upload";
            dragButton.style.display = "none";
            // console.log('File is inside the drag area');
        });
        // when file leaves the drag area
        dropAreaUser.addEventListener("dragleave", () => {
            dropAreaUser.classList.remove("active");
            // console.log('File left the drag area');
            dragText.textContent = "drag your file here";
        });
        // when file is dropped
        dropAreaUser.addEventListener("drop", (event) => {
            event.preventDefault();
            dropAreaUser.classList.add("dropped");
            // console.log('File is dropped in drag area');
            file = event.dataTransfer.files[0]; // grab single file even if user selects multiple files
            // console.log(file);
            displayFile(uploadedFileUser, dropAreaUser);
        });
    }

    // get post image from manual click
    inputPost.addEventListener("change", function () {
        file = this.files[0];
        dropAreaUser.classList.add("active");
    });

    function displayFile(uploadedFile, dropArea) {
        let fileType = file.type;
        // console.log(fileType);
        let validExtensions = ["image/*"];
        if (validExtensions.includes(fileType)) {
            let fileReader = new FileReader();
            fileReader.onload = () => {
                uploadedFile.innerHTML = `<div class="dragText">uploaded</div>
        <div class="uploaded-file">${file.name}</div>`;
                dropArea.classList.add("dropped");
            };
            fileReader.readAsDataURL(file);
        } else {
            alert("This is not an Image File");
            dropArea.classList.remove("active");
            dragText.textContent = "Drag and drop your file, or";
            dragButton.style.display = "unset";
        }
    }
}


