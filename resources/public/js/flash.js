const flashMessage = document.getElementById('flashBox').innerHTML
const flashText = document.getElementById('flashText')
if (flashMessage !== "") {
    const flashJson = JSON.parse(flashMessage);
    const alertBox = document.getElementById('alertBox');
    alertBox.classList.add(`alert-${flashJson.type}`)
    flashText.innerHTML = flashJson.message
}

const flashContainer = document.getElementById('flashContainer')
if (flashText?.innerHTML === "") {
    flashContainer.style.display = 'none'
}