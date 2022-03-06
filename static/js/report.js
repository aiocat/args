var captchaElement = document.getElementById("captcha");

document.getElementById("send").onclick = () => {
    captchaElement.execute();
};

captchaElement.addEventListener("verified", async ({ token }) => {
    let response = await fetch("/reports", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            hcaptcha: token,
            id: "{{.Id}}",
        }),
    });

    if (response.ok) {
        window.location.href = "/";
    } else {
        let responseJson = await response.json();
        document.getElementById("error").innerText = responseJson.error;
    }
});

captchaElement.addEventListener("error", ({ error }) => {
    document.getElementById("error").innerText = error;
});