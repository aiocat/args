let argumentStorage = localStorage.getItem("arguments");

if (!argumentStorage) {
    argumentStorage = JSON.stringify([]);
    localStorage.setItem("arguments", argumentStorage);
}

let argumentStorageParsed = JSON.parse(argumentStorage);

// Load all arguments
argumentStorageParsed.forEach((arg) => {
    let argLink = document.createElement("a");
    argLink.href = `/arguments/${arg.id}`;
    argLink.innerText = arg.title;

    let deleteBtn = document.createElement("button");
    deleteBtn.innerText = "Delete";
    deleteBtn.onclick = () => {
        argumentStorageParsed = argumentStorageParsed.filter(
            (i) => i.id !== arg.id
        );
        localStorage.setItem("arguments", JSON.stringify(argumentStorageParsed));
        window.location.reload();
    };

    document.body.append(argLink, deleteBtn);
});
