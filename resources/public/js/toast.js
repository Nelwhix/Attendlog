function showToast({type, message}) {
    const flashContainer = document.getElementById('flashContainer')
    const flashMessage = document.getElementById('flashBox').innerHTML
    const flashText = document.getElementById('flashText')
    const alertBox = document.getElementById('alertBox');

    flashText.innerHTML = message
    alertBox.classList.add(`alert-${type}`)
    flashContainer.style.display = 'block'
}