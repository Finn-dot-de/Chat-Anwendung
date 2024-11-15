const sendmessage = document.getElementById('send-button')

sendmessage.addEventListener('click', function() {
    const input = document.getElementById('message-input')
    const message = input.value
    const send = fetch('' + message)
})

function sendmessage() {
    const input = document.getElementById('message-input')
    const message = input.value
    const send = fetch('' + message)
}

function getmessage() {
    const send = fetch('')
    const message = send.json()
    return message
}

function displaymessage() {
    const message = getmessage()
    const display = document.getElementById('message-display')
    display.innerHTML = message
}

displaymessage()

