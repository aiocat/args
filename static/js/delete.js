var secretElement = document.getElementById("secret");

document.getElementById("send").onclick = async () => {
    let response = await fetch(`/arguments/${secretElement.value}`, {
        method: "DELETE",
    });

    if (response.ok) {
        window.location.href = "/";
    } else {
        let responseJson = await response.json();
        document.getElementById("error").innerText = responseJson.error;
    }
};