package main

const indexHTML = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>Four Relay Shield</title>

    <!-- Bootstrap CSS -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/css/bootstrap.min.css" rel="stylesheet"
        integrity="sha384-giJF6kkoqNQ00vy+HMDP7azOuL0xtbfIcaT9wjKHr8RbDVddVHyTfAAsrekwKmP1" crossorigin="anonymous">

    <script>
        const relayNormalClass = 'btn btn-secondary btn-lg';
        const relayActiveClass = 'btn btn-success btn-lg';

        const relay01Id = 'btn-relay-01';
        const relay02Id = 'btn-relay-02';
        const relay03Id = 'btn-relay-03';
        const relay04Id = 'btn-relay-04';

        let websocketInstance;

        window.onload = () => {
			let loc = window.location, new_uri;
			if (loc.protocol === "https:") {
				new_uri = "wss:";
			} else {
				new_uri = "ws:";
			}
			new_uri += "//" + loc.host;
			new_uri += loc.pathname + "ws";

            websocketInstance = new WebSocket(new_uri);
            websocketInstance.onopen = (evt) => this.onOpen(evt);
            websocketInstance.onmessage = (message) => this.onMessage(message);
        }

        function onOpen(evt) {
            console.log('Connected to relay via websocket');
        }

        function onMessage(message) {
            let dataFromRelay = JSON.parse(message.data);
            let topic = dataFromRelay[0];
            let data = dataFromRelay[1];

            switch (topic) {
                case 'INITIAL_RELAY_STATUS':
                    setIconsStyles(data.relay1Status, data.relay2Status, data.relay3Status, data.relay4Status)
                    break;
                case 'NOTIFY_RELAY_STATUS':
                    setIconStyle(data.relayNumber, data.relayNewValue);
                    break;
                default:
                    break;
            }
        }

        function setIconsStyles(relay1Status, relay2Status, relay3Status, relay4Status) {
            if (relay1Status) {
                let buttonRelay1 = document.getElementById(relay01Id);
                buttonRelay1.className = relayActiveClass;
            }

            if (relay2Status) {
                let buttonRelay2 = document.getElementById(relay02Id);
                buttonRelay2.className = relayActiveClass;
            }

            if (relay3Status) {
                let buttonRelay3 = document.getElementById(relay03Id);
                buttonRelay3.className = relayActiveClass;
            }

            if (relay4Status) {
                let buttonRelay4 = document.getElementById(relay04Id);
                buttonRelay4.className = relayActiveClass;
            }
        }

        function setIconStyle(relayNumber, relayValue) {
            let idButtonToChange;
            switch (relayNumber) {
                case 1:
                    idButtonToChange = relay01Id;
                    break;
                case 2:
                    idButtonToChange = relay02Id;
                    break;
                case 3:
                    idButtonToChange = relay03Id;
                    break;
                case 4:
                    idButtonToChange = relay04Id;
                    break;
                default:
                    break;
            }

            let buttonRelay = document.getElementById(idButtonToChange);

            if (relayValue) {
                buttonRelay.className = relayActiveClass;
            } else {
                buttonRelay.className = relayNormalClass;
            }
        }

        function relayButtonClick(relayNumber) {
            let idButton;
            switch (relayNumber) {
                case 1:
                    idButton = relay01Id;
                    break;
                case 2:
                    idButton = relay02Id;
                    break;
                case 3:
                    idButton = relay03Id;
                    break;
                case 4:
                    idButton = relay04Id;
                    break;
                default:
                    break;
            }

            const buttonRelay = document.getElementById(idButton);
            const value = buttonRelay.className !== relayActiveClass;

            const message = ["SET_VALUE_RELAY_COMMAND",
                { "relayNumber": relayNumber, "relayValue": value }
            ];

            websocketInstance.send(JSON.stringify(message));
        }
    </script>
</head>

<body>
    <div class="container">
        <h1>Four Relay Shield</h1>

        <button id="btn-relay-01" type="button" class="btn btn-secondary btn-lg" onclick="relayButtonClick(1);">RELAY
            #1</button>
        <button id="btn-relay-02" type="button" class="btn btn-secondary btn-lg" onclick="relayButtonClick(2);">RELAY
            #2</button>
        <button id="btn-relay-03" type="button" class="btn btn-secondary btn-lg" onclick="relayButtonClick(3);">RELAY
            #3</button>
        <button id="btn-relay-04" type="button" class="btn btn-secondary btn-lg" onclick="relayButtonClick(4);">RELAY
            #4</button>
    </div>
</body>

</html>

`
