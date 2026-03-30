let socket = null
let nickname = ''
let selectedRoom = null

const lobbyScreen = document.getElementById('lobby-screen')
const chatScreen = document.getElementById('chat-screen')
const nicknameInput = document.getElementById('nickname-input')
const roomList = document.getElementById('room-list')
const joinBtn = document.getElementById('join-btn')
const messagesDiv = document.getElementById('messages')
const messageInput = document.getElementById('message-input')
const sendBtn = document.getElementById('send-btn')
const roomNameSpan = document.getElementById('room-name')
const backBtn = document.getElementById('back-btn')

async function loadRooms() {
    const res = await fetch('/rooms')
    const rooms = await res.json()

    rooms.forEach(room => {
        const div = document.createElement('div')
        div.className = 'room-option'
        div.textContent = room.name 
        div.addEventListener('click', () => {
            document.querySelectorAll('.room-option').forEach(el => el.classList.remove('selected'))
            div.classList.add('selected')
            selectedRoom = room
        })
        roomList.appendChild(div)
    })
}

loadRooms()

joinBtn.addEventListener('click', () => {
    nickname = nicknameInput.value.trim()
    if (!nickname) {
        alert('Please enter a nickname')
        return
    }
    if (!selectedRoom) {
        alert('Please select a room')
        return
    }
    joinRoom()
})

function joinRoom() {
    lobbyScreen.classList.add('hidden')
    chatScreen.classList.remove('hidden')
    roomNameSpan.textContent = selectedRoom.name
    messagesDiv.innerHTML = ''

    socket = new WebSocket(`ws://localhost:8080/ws?room_id=${selectedRoom.id}&nickname=${nickname}`)

    socket.onmessage = (event) => {
        const msg = JSON.parse(event.data)
        appendMessage(msg)
    }

    socket.onclose = () => {
        appendSystemMessage('Disconnected from chat')
    }
}

backBtn.addEventListener('click', () => {
    if (socket) {
        socket.close()
        socket = null
    }
    chatScreen.classList.add('hidden')
    lobbyScreen.classList.remove('hidden')
})  

sendBtn.addEventListener('click', sendMessage)

messageInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') {
        sendMessage()
    }
})

function sendMessage() {
    const body = messageInput.value.trim()
    if (!body || !socket) return
    socket.send(JSON.stringify({ body: body }))
    messageInput.value = ''
}

function appendMessage(msg) {
    const isMine = msg.nickname === nickname

    const row = document.createElement('div')
    row.className = `bubble-row ${isMine ? 'mine' : 'theirs'}`

    if (!isMine) {
        const name = document.createElement('div')
        name.className = 'bubble-name'
        name.textContent = msg.nickname
        row.appendChild(name)
    }

    const bubble = document.createElement('div')
    bubble.className = 'bubble'
    bubble.textContent = msg.body
    row.appendChild(bubble)

    messagesDiv.appendChild(row)
    messagesDiv.scrollTop = messagesDiv.scrollHeight
}

function appendSystemMessage(text) {
    const div = document.createElement('div')
    div.style.cssText = 'text-align:center; color:#888; font-size:13px; padding:8px;'
    div.textContent = text
    messagesDiv.appendChild(div)
}