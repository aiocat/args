var titleElement = document.getElementById("title");
var captchaElement = document.getElementById("captcha");

document.getElementById("send").onclick = () => {
    if (titleElement.value.length < 4) {
        document.getElementById("error").innerText = "Title is too short";
        return;
    } else if (titleElement.value.length > 120) {
        document.getElementById("error").innerText = "Title is too long";
        return;
    }
    
    captchaElement.execute();
};

captchaElement.addEventListener("verified", async ({ token }) => {
    let response = await fetch("/arguments", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            hcaptcha: token,
            title: titleElement.value,
        }),
    });

    let responseJson = await response.json();

    if (response.ok) {
        let notification = document.createElement("div");
        notification.className = "notification";

        let title = document.createElement("h1");
        title.innerText = "Success";

        let linkToArg = document.createElement("a");
        linkToArg.innerText = "Click here to jump into discussion";
        linkToArg.href = `/arguments/${responseJson.id}`;

        let description = document.createElement("h2");
        description.innerText =
            "Save this secret key to remove your discussion later.";

        let secret = document.createElement("p");
        secret.innerText = responseJson.secret;

        let remove = document.createElement("button");
        remove.innerText = "Ok!";
        remove.onclick = () => {
            notification.remove();
            window.location.href = `/arguments/${responseJson.id}`;
        };

        notification.append(title, description, linkToArg, secret, remove);
        document.body.appendChild(notification);
    } else {
        document.getElementById("error").innerText = responseJson.error;
    }
});

captchaElement.addEventListener("error", ({ error }) => {
    document.getElementById("error").innerText = error;
});