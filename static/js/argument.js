var opinion = 0;

// Prepare opinions
function prepareOpinions() {
    let agree = document.getElementById("agree");
    let however = document.getElementById("however");
    let disagree = document.getElementById("disagree");

    agree.onclick = () => {
        agree.style.borderWidth = "3px";
        however.style.borderWidth = "1px";
        disagree.style.borderWidth = "1px";
        opinion = 1;
    };

    however.onclick = () => {
        agree.style.borderWidth = "1px";
        however.style.borderWidth = "3px";
        disagree.style.borderWidth = "1px";
        opinion = 2;
    };

    disagree.onclick = () => {
        agree.style.borderWidth = "1px";
        however.style.borderWidth = "1px";
        disagree.style.borderWidth = "3px";
        opinion = 3;
    };
}

// Parse argument
function parseArgument(argument_id) {
    let argument = document.getElementById("argument");
    let captchaElement = document.getElementById("captcha");

    document.getElementById("send").onclick = () => {
        if (opinion === 0) {
            document.getElementById("error").innerText = "Please select an opinion";
            return;
        }
        captchaElement.execute();
    };

    captchaElement.addEventListener("verified", async ({ token }) => {
        let response = await fetch(`/arguments/${argument_id}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                hcaptcha: token,
                opinion,
                argument: argument.value,
            }),
        });

        let responseJson = await response.json();

        if (response.ok) {
            let notification = document.createElement("div");
            notification.className = "notification";

            let title = document.createElement("h1");
            title.innerText = "Success";

            let description = document.createElement("h2");
            description.innerText =
                "Save this secret key to remove your argument later.";

            let secret = document.createElement("p");
            secret.innerText = responseJson.secret;

            let remove = document.createElement("button");
            remove.innerText = "Ok!";
            remove.onclick = () => {
                notification.remove();
                window.location.reload();
            };

            notification.append(title, description, secret, remove);
            document.body.appendChild(notification);
        } else {
            document.getElementById("error").innerText = responseJson.error;
        }

        argument.value = "";
        opinion = 0;
        agree.style.borderWidth = "1px";
        however.style.borderWidth = "1px";
        disagree.style.borderWidth = "1px";
    });

    captchaElement.addEventListener("error", ({ error }) => {
        document.getElementById("error").innerText = error;
    });
}

// Prepare save
function prepareSave(title, argument_id) {
    let argumentStorage = localStorage.getItem("arguments")
    let saveElement = document.getElementById("save")

    if (!argumentStorage) {
        argumentStorage = JSON.stringify([])
        localStorage.setItem("arguments", argumentStorage)
    }

    let argumentStorageParsed = JSON.parse(argumentStorage)
    
    if (argumentStorageParsed.filter(i => i.id === argument_id).length) {
        saveElement.innerText = "Saved"
        saveElement.disabled = true
    } else {
        saveElement.onclick = () => {
            argumentStorageParsed.push({"title": title, "id": argument_id})
    
            argumentStorage = JSON.stringify(argumentStorageParsed)
            localStorage.setItem("arguments", argumentStorage)

            saveElement.innerText = "Saved"
            saveElement.disabled = true
        }
    }
}

// Calculate replies
function calculateReplies() {
    let replies = document.getElementsByClassName("replies")[0];
    let accepts = 0;
    let howevers = 0;
    let declines = 0;
    let totalReplies = replies.childElementCount;

    for (let i = 0; i < totalReplies; i += 3) {
        let foundElem = replies.children[i];
        let opinion = parseInt(foundElem.className.split("-")[1]);

        if (opinion === 1) {
            accepts++;
        } else if (opinion === 2) {
            howevers++;
        } else {
            declines++;
        }
    }

    // Calculate percents
    document.getElementById("agree-percent").innerText =
        "%" + Math.round((accepts / (totalReplies / 3)) * 100);
    document.getElementById("however-percent").innerText =
        "%" + Math.round((howevers / (totalReplies / 3)) * 100);
    document.getElementById("disagree-percent").innerText =
        "%" + Math.round((declines / (totalReplies / 3)) * 100);
}

if (document.getElementsByClassName("replies")[0].childElementCount > 0) {
    calculateReplies();
}

prepareOpinions();