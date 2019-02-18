package handler

const (
	indexHTML = `
<!doctype html>
<html>
	<head>
		<meta charset="utf-8">
	</head>
	<body>
		<form action="/login" method="get">
			<input type="input" name="name" placeholder="Your favourite streamer name">
			<input type="submit" value="Login using Twitch">
		</form>
	</body>
</html>
`

	streamHTML = `
<!doctype html>
<html>
	<head>
		<meta charset="utf-8">
		<script>  
			window.addEventListener("load", function(evt) {
				var wsProto = document.location.protocol == 'https:'
					? 'wss:'
					: 'ws:';

				var ws = new WebSocket(wsProto + '//' + document.location.host +'/events');

				var $chat = document.getElementById("chat");
				var addToChat = function(message) {
					var d = document.createElement("div");
					d.innerHTML = message;
					$chat.appendChild(d);

						$chat.scrollTop = $chat.scrollHeight;
				};

				var $events = document.getElementById("events");
				var addToEvents = function(message) {
					if ($events.childNodes.length >= 10) {
						$events.removeChild($events.lastChild);
					}

					var d = document.createElement("div");
					d.innerHTML = message;
					$events.appendChild(d);

					$events.insertBefore(d, $events.firstChild);
				};

				ws.onmessage = function(e) {
					var message = JSON.parse(e.data);
					if (message.type == 'message') {
						addToChat(message.content);
					} else {
						addToEvents('[' + message.type + '] ' + message.content);
					}
				}

				ws.onopen = function(e) {
					console.log("You're connected");
				}

				ws.onclose = function(e) {
					console.log("You're disconnected");
					ws = null;
				}

				ws.onerror = function(e) {
					console.log(e.data);
				}
			});
		</script>
	</head>
	<body>
		<iframe
			src="https://player.twitch.tv/?channel=%s&muted=true"
			height="720"
			width="1280"
			frameborder="0"
			scrolling="no"
			allowfullscreen="true">
		</iframe>
		<div id="chat" style="overflow-y: scroll; height:720px; width:500px; display:inline-block; word-break: break-word;"></div>
		<div id="events" style="height:100px; width:500px; word-break: break-word;"></div>
	</body>
</html>
`
)
